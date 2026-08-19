package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/payment"
	"github.com/zourleb/zourleb-api/pkg/phone"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// BookingService orchestrates the booking flow: pricing, phone-verification
// gating, seat reservation, listing, cancellation and payment initiation.
type BookingService struct {
	bookings *repository.BookingRepository
	tours    *repository.TourRepository
	users    *repository.UserRepository
	payments *repository.PaymentRepository
	payReg   *payment.Registry
	settings *SettingsService
	notify   *NotificationService
	defLang  string
}

func NewBookingService(b *repository.BookingRepository, t *repository.TourRepository, u *repository.UserRepository, pay *repository.PaymentRepository, reg *payment.Registry, s *SettingsService, n *NotificationService, defLang string) *BookingService {
	return &BookingService{bookings: b, tours: t, users: u, payments: pay, payReg: reg, settings: s, notify: n, defLang: defLang}
}

// Create books a tour departure for a user. When booking.require_otp is on, the
// contact phone must match the user's verified phone.
func (s *BookingService) Create(ctx context.Context, userID uint, req models.CreateBookingRequest, locale string) (*models.BookingResponse, error) {
	tour, err := s.tours.FindByID(req.TourID)
	if err != nil {
		return nil, response.ErrTourNotFound
	}
	if tour.Status != models.TourStatusPublished {
		return nil, response.ErrTourNotFound
	}

	contact, err := phone.Normalize(req.ContactPhone)
	if err != nil {
		return nil, response.ErrPhoneInvalid
	}

	// email-verification gate
	if s.settings.Bool(models.SettingBookingEmail, false) {
		u, err := s.users.FindByID(userID)
		if err != nil {
			return nil, response.ErrUnauthorized
		}
		if !u.EmailVerified() {
			return nil, response.ErrEmailNotVerified
		}
	}

	// phone-verification gate
	if s.settings.Bool(models.SettingBookingOTP, true) {
		u, err := s.users.FindByID(userID)
		if err != nil {
			return nil, response.ErrUnauthorized
		}
		if !u.PhoneVerified() || u.Phone != contact {
			return nil, response.ErrPhoneNotVerified
		}
	}

	// resolve departure (optional)
	var departure *models.TourDeparture
	if req.DepartureID != nil {
		departure, err = s.tours.FindDeparture(*req.DepartureID)
		if err != nil || departure.TourID != tour.ID {
			return nil, response.ErrBadRequest.WithMessage("Invalid departure.")
		}
		if departure.Status != models.DepartureOpen {
			return nil, response.ErrDepartureFull
		}
	}

	// price the booking
	subtotal, currency := s.price(tour, req)

	booking := &models.Booking{
		Code:                 newBookingCode(),
		UserID:               userID,
		TourID:               tour.ID,
		DepartureID:          req.DepartureID,
		Status:               models.BookingAwaitingPayment,
		TravelersCount:       len(req.Travelers),
		Subtotal:             subtotal,
		Currency:             currency,
		ContactPhone:         contact,
		ContactPhoneVerified: true,
		PaymentStatus:        models.PayStatusPending,
		Notes:                req.Notes,
	}
	for _, t := range req.Travelers {
		ph := ""
		if t.Phone != "" {
			if n, err := phone.Normalize(t.Phone); err == nil {
				ph = n
			}
		}
		booking.Travelers = append(booking.Travelers, models.BookingTraveler{
			FullName: t.FullName, TravelerType: t.TravelerType, Phone: ph,
		})
	}

	ok, err := s.bookings.CreateWithSeatHold(booking, departure)
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	if !ok {
		return nil, response.ErrDepartureFull
	}

	title := s.tours.TitleFor(tour.ID, locale, s.defLang)
	resp := models.NewBookingResponse(booking, title)
	if s.notify != nil {
		s.notify.Notify(ctx, userID, "booking.created", "Booking received", "Your booking "+booking.Code+" is awaiting payment.", map[string]any{"booking_code": booking.Code, "tour_title": title})
	}
	return &resp, nil
}

// List returns the user's bookings.
func (s *BookingService) List(userID uint, p pagination.Params, locale string) ([]models.BookingResponse, pagination.Meta, error) {
	rows, total, err := s.bookings.ListForUser(userID, p)
	if err != nil {
		return nil, pagination.Meta{}, response.ErrInternal.Wrap(err)
	}
	out := make([]models.BookingResponse, 0, len(rows))
	for i := range rows {
		title := s.tours.TitleFor(rows[i].TourID, locale, s.defLang)
		out = append(out, models.NewBookingResponse(&rows[i], title))
	}
	return out, pagination.NewMeta(p, total), nil
}

// Get returns one booking the user owns.
func (s *BookingService) Get(userID uint, code, locale string) (*models.BookingResponse, error) {
	b, err := s.bookings.FindByCodeForUser(code, userID)
	if err != nil {
		return nil, response.ErrBookingNotFound
	}
	title := s.tours.TitleFor(b.TourID, locale, s.defLang)
	resp := models.NewBookingResponse(b, title)
	return &resp, nil
}

// Pay initiates payment for a booking via the chosen provider. Manual/offline
// payments stay pending until an admin confirms; immediate providers settle now.
func (s *BookingService) Pay(ctx context.Context, userID uint, code string, req models.PayRequest) (*models.PayResponse, error) {
	b, err := s.bookings.FindByCodeForUser(code, userID)
	if err != nil {
		return nil, response.ErrBookingNotFound
	}
	if b.PaymentStatus == models.PayStatusPaid {
		return nil, response.ErrConflict.WithMessage("Booking is already paid.")
	}
	if existing, err := s.payments.FindByIdempotencyKey(req.IdempotencyKey); err == nil {
		return &models.PayResponse{PaymentID: existing.ID, Status: existing.Status}, nil
	}
	provider, err := s.payReg.Get(req.Provider)
	if err != nil {
		return nil, response.ErrBadRequest.WithMessage("Unknown payment provider.")
	}

	pay := &models.Payment{
		PayableType: models.PayableBooking, PayableID: b.ID,
		PayerType: models.PayerUser, PayerID: userID,
		Amount: b.Subtotal, Currency: b.Currency,
		Provider: provider.Name(), Status: models.PayStatusPending,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.payments.Create(pay); err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}

	res, err := provider.Charge(ctx, payment.ChargeRequest{
		Amount: b.Subtotal, Currency: b.Currency,
		Description: "Booking " + b.Code, IdempotencyKey: req.IdempotencyKey,
		Reference: b.Code, ReturnURL: req.ReturnURL,
	})
	if err != nil {
		return nil, response.ErrInternal.Wrap(err)
	}
	_ = s.payments.UpdateStatus(pay.ID, string(res.Status), res.ProviderRef)
	if res.Status == payment.StatusPaid {
		_ = s.bookings.MarkPaid(b.ID)
		if s.notify != nil {
			title := s.tours.TitleFor(b.TourID, s.defLang, s.defLang)
			s.notify.Notify(ctx, userID, "booking.paid", "Payment confirmed", "Your booking "+b.Code+" for "+title+" is confirmed.", map[string]any{"booking_code": b.Code, "tour_title": title})
		}
	}
	return &models.PayResponse{PaymentID: pay.ID, Status: string(res.Status), RedirectURL: res.RedirectURL}, nil
}

// Cancel cancels a pending/awaiting booking and releases seats.
func (s *BookingService) Cancel(userID uint, code string) error {
	b, err := s.bookings.FindByCodeForUser(code, userID)
	if err != nil {
		return response.ErrBookingNotFound
	}
	if b.Status == models.BookingCancelled || b.Status == models.BookingCompleted {
		return response.ErrConflict.WithMessage("Booking can no longer be cancelled.")
	}
	if err := s.bookings.Cancel(b); err != nil {
		return response.ErrInternal.Wrap(err)
	}
	return nil
}

// price computes subtotal by matching traveler types to tour prices, falling
// back to the tour's price_from when a specific price row is missing.
func (s *BookingService) price(tour *models.Tour, req models.CreateBookingRequest) (float64, string) {
	currency := tour.BaseCurrency
	if currency == "" {
		currency = s.settings.String(models.SettingDefaultCurrency, "USD")
	}

	// index prices by traveler type for the chosen currency
	byType := map[string]float64{}
	for _, pr := range tour.Prices {
		if pr.Currency != currency {
			continue
		}
		if pr.DepartureID == nil || (req.DepartureID != nil && *pr.DepartureID == *req.DepartureID) {
			byType[pr.TravelerType] = pr.Amount
		}
	}

	var subtotal float64
	for _, t := range req.Travelers {
		if amt, ok := byType[t.TravelerType]; ok {
			subtotal += amt
		} else {
			subtotal += tour.PriceFrom
		}
	}
	return subtotal, currency
}

func newBookingCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		b[i] = alphabet[n.Int64()]
	}
	return "ZB" + string(b)
}

var _ = fmt.Sprintf
var _ = time.Now
