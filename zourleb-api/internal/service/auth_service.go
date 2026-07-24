package service

import (
	"context"
	"strings"
	"time"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/hash"
	"github.com/zourleb/zourleb-api/pkg/oauth"
	"github.com/zourleb/zourleb-api/pkg/response"
	"github.com/zourleb/zourleb-api/pkg/token"
)

// AuthService implements registration, login, token issuance and Google sign-in.
type AuthService struct {
	users  *repository.UserRepository
	tokens *token.Manager
	google *oauth.GoogleVerifier
	cfg    *config.Config
}

func NewAuthService(users *repository.UserRepository, tm *token.Manager, gv *oauth.GoogleVerifier, cfg *config.Config) *AuthService {
	return &AuthService{users: users, tokens: tm, google: gv, cfg: cfg}
}

type AuthContext struct {
	UserAgent string
	IP        string
}

// Register creates an email/password account and assigns the tourist role.
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest, ac AuthContext) (*models.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	if _, err := s.users.FindByEmail(email); err == nil {
		return nil, response.ErrEmailTaken
	} else if !repository.IsNotFound(err) {
		return nil, response.ErrInternal.Wrap(err)
	}

	pw, err := hash.Password(req.Password)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	locale := req.Locale
	if locale == "" {
		locale = s.cfg.App.DefaultLang
	}

	u := &models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: &pw,
		Locale:       locale,
		Status:       models.UserStatusActive,
	}
	if err := s.users.Create(u); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	s.assignTourist(u.ID)

	return s.issue(u, ac)
}

// Login authenticates an email/password user.
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest, ac AuthContext) (*models.AuthResponse, error) {
	u, err := s.users.FindByEmail(normalizeEmail(req.Email))
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, response.ErrInvalidCredentials
		}
		return nil, response.ErrInternal.Wrap(err)
	}
	if u.PasswordHash == nil || !hash.CheckPassword(*u.PasswordHash, req.Password) {
		return nil, response.ErrInvalidCredentials
	}
	if !u.IsActive() {
		return nil, response.ErrAccountBlocked
	}
	_ = s.users.TouchLogin(u.ID)
	return s.issue(u, ac)
}

// Google verifies a Google ID token and links or creates the user.
func (s *AuthService) Google(ctx context.Context, req models.GoogleAuthRequest, ac AuthContext) (*models.AuthResponse, error) {
	if !s.google.Configured() {
		return nil, response.ErrFeatureOff.WithMessage("Google sign-in is not configured.")
	}
	profile, err := s.google.Verify(ctx, req.IDToken)
	if err != nil {
		return nil, response.ErrTokenInvalid.Wrap(err)
	}

	// Try by google_id, then by email (link existing account).
	u, err := s.users.FindByGoogleID(profile.Sub)
	if repository.IsNotFound(err) {
		u, err = s.users.FindByEmail(normalizeEmail(profile.Email))
	}

	switch {
	case err == nil: // existing user — ensure google_id is linked
		if u.GoogleID == nil {
			gid := profile.Sub
			u.GoogleID = &gid
		}
		if u.EmailVerifiedAt == nil && profile.EmailVerified {
			now := time.Now()
			u.EmailVerifiedAt = &now
		}
		if err := s.users.Update(u); err != nil {
			return nil, response.ErrInternal.Wrap(err)
		}
	case repository.IsNotFound(err): // new user
		locale := req.Locale
		if locale == "" {
			locale = s.cfg.App.DefaultLang
		}
		gid := profile.Sub
		now := time.Now()
		u = &models.User{
			Name:     profile.Name,
			Email:    normalizeEmail(profile.Email),
			GoogleID: &gid,
			Avatar:   profile.Picture,
			Locale:   locale,
			Status:   models.UserStatusActive,
		}
		if profile.EmailVerified {
			u.EmailVerifiedAt = &now
		}
		if err := s.users.Create(u); err != nil {
			return nil, response.ErrInternal.Wrap(err)
		}
		s.assignTourist(u.ID)
	default:
		return nil, response.ErrInternal.Wrap(err)
	}

	if !u.IsActive() {
		return nil, response.ErrAccountBlocked
	}
	_ = s.users.TouchLogin(u.ID)
	return s.issue(u, ac)
}

// Refresh rotates a refresh token, returning a new token pair.
func (s *AuthService) Refresh(ctx context.Context, raw string, ac AuthContext) (*models.AuthResponse, error) {
	rt, err := s.users.FindRefreshByHash(hash.SHA256(raw))
	if err != nil {
		return nil, response.ErrTokenInvalid
	}
	if !rt.Active(time.Now()) {
		return nil, response.ErrTokenInvalid
	}
	// rotate: revoke the presented token
	if err := s.users.RevokeRefreshToken(rt.ID); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	u, err := s.users.FindByID(rt.UserID)
	if err != nil {
		return nil, response.ErrTokenInvalid
	}
	return s.issue(u, ac)
}

// Logout revokes a specific refresh token.
func (s *AuthService) Logout(ctx context.Context, raw string) error {
	rt, err := s.users.FindRefreshByHash(hash.SHA256(raw))
	if err != nil {
		return nil // idempotent
	}
	return s.users.RevokeRefreshToken(rt.ID)
}

// issue builds the access+refresh pair and the auth response.
func (s *AuthService) issue(u *models.User, ac AuthContext) (*models.AuthResponse, error) {
	access, exp, err := s.tokens.IssueAccess(u.ID, u.Email, u.Locale)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	rawRefresh, err := token.NewRefreshToken()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	rt := &models.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash.SHA256(rawRefresh),
		ExpiresAt: time.Now().Add(s.tokens.RefreshTTL()),
		UserAgent: ac.UserAgent,
		IP:        ac.IP,
	}
	if err := s.users.CreateRefreshToken(rt); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	roles, _ := s.users.RoleKeys(u.ID)
	return &models.AuthResponse{
		User:         models.NewUserResponse(u, roles),
		AccessToken:  access,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(time.Until(exp).Seconds()),
	}, nil
}

func (s *AuthService) assignTourist(userID uint) {
	if roleID, err := s.users.RoleIDByKey(models.RoleTourist); err == nil {
		_ = s.users.AssignRole(userID, roleID, nil)
	}
}

func normalizeEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}

// AuthContextFromRequest is a small helper for handlers.
func AuthContextFromRequest(userAgent, ip string) AuthContext {
	return AuthContext{UserAgent: userAgent, IP: ip}
}
