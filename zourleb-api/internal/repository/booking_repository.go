package repository

import (
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/pkg/pagination"
)

// BookingRepository persists bookings and their travelers/payments, with the
// seat-reservation transaction kept short.
type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateWithSeatHold creates a booking and atomically reserves departure seats.
// Returns ErrDepartureFull-equivalent via the returned bool=false when seats
// are insufficient. The caller maps that to a domain error.
func (r *BookingRepository) CreateWithSeatHold(b *models.Booking, departure *models.TourDeparture) (bool, error) {
	ok := true
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if departure != nil {
			// atomic conditional update: only succeeds if seats remain
			res := tx.Model(&models.TourDeparture{}).
				Where("id = ? AND capacity - seats_booked >= ?", departure.ID, b.TravelersCount).
				UpdateColumn("seats_booked", gorm.Expr("seats_booked + ?", b.TravelersCount))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				ok = false
				return nil // not enough seats; abort create without error
			}
		}
		return tx.Create(b).Error
	})
	return ok, err
}

// FindByCodeForUser loads a booking the user owns, with travelers.
func (r *BookingRepository) FindByCodeForUser(code string, userID uint) (*models.Booking, error) {
	var b models.Booking
	err := r.db.Preload("Travelers").Preload("Payments").
		Where("code = ? AND user_id = ?", code, userID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListForUser returns the user's bookings, paginated.
func (r *BookingRepository) ListForUser(userID uint, p pagination.Params) ([]models.Booking, int64, error) {
	q := r.db.Model(&models.Booking{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Booking
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}

// Cancel transitions a booking to cancelled and releases held seats.
func (r *BookingRepository) Cancel(b *models.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if b.DepartureID != nil {
			if err := tx.Model(&models.TourDeparture{}).
				Where("id = ?", *b.DepartureID).
				UpdateColumn("seats_booked", gorm.Expr("GREATEST(seats_booked - ?, 0)", b.TravelersCount)).
				Error; err != nil {
				return err
			}
		}
		return tx.Model(b).Updates(map[string]any{"status": models.BookingCancelled}).Error
	})
}

// MarkPaid marks a booking confirmed and paid after successful settlement.
func (r *BookingRepository) MarkPaid(id uint) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).
		Updates(map[string]any{
			"payment_status": models.PayStatusPaid,
			"status":         models.BookingConfirmed,
		}).Error
}

// ListForAgency returns bookings across an agency's tours (agency portal).
func (r *BookingRepository) ListForAgency(agencyID uint, p pagination.Params) ([]models.Booking, int64, error) {
	sub := r.db.Model(&models.Tour{}).Select("id").Where("agency_id = ?", agencyID)
	q := r.db.Model(&models.Booking{}).Where("tour_id IN (?)", sub)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Booking
	err := q.Order("created_at desc").Offset(p.Offset()).Limit(p.Limit()).Find(&rows).Error
	return rows, total, err
}
