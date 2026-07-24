package service

import (
	"context"
	"fmt"
	"time"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/cache"
	"github.com/zourleb/zourleb-api/pkg/hash"
	"github.com/zourleb/zourleb-api/pkg/otp"
	"github.com/zourleb/zourleb-api/pkg/phone"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// OTPService issues and verifies phone OTP codes. Live codes live in the cache
// (Redis/in-memory) with a TTL; the DB keeps the audit trail.
type OTPService struct {
	repo  *repository.OTPRepository
	users *repository.UserRepository
	cache cache.Store
	disp  *otp.Dispatcher
	cfg   config.OTPConfig
}

func NewOTPService(repo *repository.OTPRepository, users *repository.UserRepository, c cache.Store, disp *otp.Dispatcher, cfg config.OTPConfig) *OTPService {
	return &OTPService{repo: repo, users: users, cache: c, disp: disp, cfg: cfg}
}

type otpState struct {
	CodeHash string
	Attempts int
}

func codeKey(p string) string     { return "otp:code:" + p }
func cooldownKey(p string) string { return "otp:cooldown:" + p }
func attemptKey(p string) string  { return "otp:attempts:" + p }

// Request generates and dispatches an OTP, enforcing cooldown and daily cap.
func (s *OTPService) Request(ctx context.Context, req models.OTPRequestRequest, ip string) error {
	e164, err := phone.Normalize(req.Phone)
	if err != nil {
		return response.ErrPhoneInvalid
	}

	// resend cooldown
	if _, hit, _ := s.cache.Get(ctx, cooldownKey(e164)); hit {
		return response.ErrOTPCooldown
	}
	// daily cap (audit-table based)
	since := time.Now().Add(-24 * time.Hour)
	if n, err := s.repo.CountSince(e164, since); err == nil && int(n) >= s.cfg.DailyCap {
		return response.ErrTooManyRequest.WithMessage("Daily verification limit reached.")
	}

	code, err := otp.Generate(s.cfg.Length)
	if err != nil {
		return response.ErrInternal.Wrap(err)
	}
	codeHash := hash.SHA256(code)

	if err := s.cache.Set(ctx, codeKey(e164), codeHash, s.cfg.TTL); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	_ = s.cache.Set(ctx, attemptKey(e164), "0", s.cfg.TTL)
	_ = s.cache.Set(ctx, cooldownKey(e164), "1", s.cfg.ResendWindow)

	// audit row (no plaintext stored)
	_ = s.repo.Create(&models.PhoneVerification{
		Phone:       e164,
		Channel:     req.Channel,
		Purpose:     req.Purpose,
		CodeHash:    codeHash,
		MaxAttempts: s.cfg.MaxAttempts,
		ExpiresAt:   time.Now().Add(s.cfg.TTL),
		IP:          ip,
	})

	// dispatch (external call — outside any DB transaction)
	sent, err := s.disp.Send(ctx, req.Channel, e164, code)
	if err != nil {
		return response.ErrInternal.Wrap(err)
	}
	if !sent {
		return response.ErrBadRequest.WithMessage(fmt.Sprintf("Channel %q is unavailable.", req.Channel))
	}
	return nil
}

// Verify checks a submitted code, enforcing expiry and max attempts. On success
// it records phone_verified_at on the user (if authenticated) and audit trail.
func (s *OTPService) Verify(ctx context.Context, req models.OTPVerifyRequest, userID uint) error {
	e164, err := phone.Normalize(req.Phone)
	if err != nil {
		return response.ErrPhoneInvalid
	}

	stored, hit, err := s.cache.Get(ctx, codeKey(e164))
	if err != nil {
		return response.ErrInternal.Wrap(err)
	}
	if !hit {
		return response.ErrOTPExpired
	}

	// attempts guard
	attempts, _ := s.cache.Incr(ctx, attemptKey(e164), s.cfg.TTL)
	if int(attempts) > s.cfg.MaxAttempts {
		_ = s.cache.Del(ctx, codeKey(e164))
		return response.ErrOTPMaxAttempt
	}

	if stored != hash.SHA256(req.Code) {
		return response.ErrOTPInvalid
	}

	// success: clear state, mark verified
	_ = s.cache.Del(ctx, codeKey(e164))
	_ = s.cache.Del(ctx, attemptKey(e164))
	_ = s.repo.MarkVerified(e164)

	if userID != 0 {
		if u, err := s.users.FindByID(userID); err == nil {
			now := time.Now()
			u.Phone = e164
			u.PhoneVerifiedAt = &now
			_ = s.users.Update(u)
		}
	}
	return nil
}

// IsVerified reports whether a phone currently has a verified user record.
func (s *OTPService) NormalizedVerifiedForUser(userID uint, rawPhone string) (bool, string) {
	e164, err := phone.Normalize(rawPhone)
	if err != nil {
		return false, ""
	}
	u, err := s.users.FindByID(userID)
	if err != nil {
		return false, e164
	}
	return u.PhoneVerified() && u.Phone == e164, e164
}
