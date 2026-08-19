package service

import (
	"context"
	"time"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/payment"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// BoostService implements the ad-boost lifecycle: create → pay → activate.
type BoostService struct {
	boosts   *repository.BoostRepository
	payments *repository.PaymentRepository
	agencies *repository.AgencyRepository
	payReg   *payment.Registry
	notify   *NotificationService
}

func NewBoostService(b *repository.BoostRepository, p *repository.PaymentRepository, a *repository.AgencyRepository, reg *payment.Registry, n *NotificationService) *BoostService {
	return &BoostService{boosts: b, payments: p, agencies: a, payReg: reg, notify: n}
}

// Packages lists available boost packages.
func (s *BoostService) Packages() ([]models.BoostPackage, error) {
	rows, err := s.boosts.Packages()
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

// Create makes a pending-payment boost for the agency's tour.
func (s *BoostService) Create(agencyID uint, req models.CreateBoostRequest) (*models.Boost, error) {
	pkg, err := s.boosts.FindPackage(req.PackageID)
	if err != nil {
		return nil, response.ErrNotFound.WithMessage("Boost package not found.")
	}
	b := &models.Boost{
		AgencyID:  agencyID,
		TourID:    &req.TourID,
		PackageID: pkg.ID,
		Placement: pkg.Placement,
		Status:    models.BoostPendingPayment,
	}
	if err := s.boosts.Create(b); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return b, nil
}

// List returns the agency's boosts.
func (s *BoostService) List(agencyID uint) ([]models.Boost, error) {
	rows, err := s.boosts.ListForAgency(agencyID)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	return rows, nil
}

// Pay initiates payment for a boost. On synchronous success (e.g. a provider
// that settles immediately) the boost is activated. For manual/offline the
// payment stays pending until an admin confirms (see AdminService.ConfirmPayment).
func (s *BoostService) Pay(ctx context.Context, agencyID, boostID uint, req models.PayRequest) (*models.PayResponse, error) {
	b, err := s.boosts.FindByID(boostID)
	if err != nil || b.AgencyID != agencyID {
		return nil, response.ErrNotFound.WithMessage("Boost not found.")
	}
	if b.Status != models.BoostPendingPayment {
		return nil, response.ErrConflict.WithMessage("Boost is not awaiting payment.")
	}
	pkg, err := s.boosts.FindPackage(b.PackageID)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	// idempotency: reuse an existing payment for this key
	if existing, err := s.payments.FindByIdempotencyKey(req.IdempotencyKey); err == nil {
		return &models.PayResponse{PaymentID: existing.ID, Status: existing.Status}, nil
	}

	provider, err := s.payReg.Get(req.Provider)
	if err != nil {
		return nil, response.ErrBadRequest.WithMessage("Unknown payment provider.")
	}

	pay := &models.Payment{
		PayableType: models.PayableBoost, PayableID: b.ID,
		PayerType: models.PayerAgency, PayerID: agencyID,
		Amount: pkg.Price, Currency: pkg.Currency,
		Provider: provider.Name(), Status: models.PayStatusPending,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.payments.Create(pay); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	res, err := provider.Charge(ctx, payment.ChargeRequest{
		Amount: pkg.Price, Currency: pkg.Currency,
		Description: "Boost: " + pkg.Name, IdempotencyKey: req.IdempotencyKey,
		Reference: "boost_" + itoa(b.ID), ReturnURL: req.ReturnURL,
	})
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	_ = s.payments.UpdateStatus(pay.ID, string(res.Status), res.ProviderRef)
	if res.Status == payment.StatusPaid {
		if err := s.activate(b, pkg, pay.ID); err != nil {
			return nil, err
		}
	}
	return &models.PayResponse{
		PaymentID: pay.ID, Status: string(res.Status), RedirectURL: res.RedirectURL,
	}, nil
}

// Activate transitions a boost to active for [now, now+duration]. Exposed for
// the admin confirm-payment path.
func (s *BoostService) Activate(boostID, paymentID uint) error {
	b, err := s.boosts.FindByID(boostID)
	if err != nil {
		return response.ErrNotFound
	}
	pkg, err := s.boosts.FindPackage(b.PackageID)
	if err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return s.activate(b, pkg, paymentID)
}

func (s *BoostService) activate(b *models.Boost, pkg *models.BoostPackage, paymentID uint) error {
	now := time.Now()
	end := now.AddDate(0, 0, pkg.DurationDays)
	b.Status = models.BoostActive
	b.StartsAt = &now
	b.EndsAt = &end
	b.PaymentID = &paymentID
	if err := s.boosts.Save(b); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	if s.notify != nil {
		if a, err := s.agencies.FindByID(b.AgencyID); err == nil && a.CreatedBy != 0 {
			s.notify.Notify(context.Background(), a.CreatedBy, "boost.activated", "Boost active", "Your boost for placement "+pkg.Placement+" is now live.", map[string]any{"boost_id": b.ID, "placement": pkg.Placement, "duration_days": pkg.DurationDays})
		}
	}
	return nil
}

// TrackImpression increments a boost's impression counter.
func (s *BoostService) TrackImpression(boostID uint) error {
	return wrapInternal(s.boosts.IncrImpressions(boostID))
}

// TrackClick increments a boost's click counter.
func (s *BoostService) TrackClick(boostID uint) error {
	return wrapInternal(s.boosts.IncrClicks(boostID))
}

func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
