package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// OTPRepository persists the phone-verification audit trail.
type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(v *models.PhoneVerification) error {
	return r.db.Create(v).Error
}

// CountSince returns how many OTPs were requested for a phone since t (daily cap).
func (r *OTPRepository) CountSince(phone string, since time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.PhoneVerification{}).
		Where("phone = ? AND created_at >= ?", phone, since).Count(&n).Error
	return n, err
}

// MarkVerified flags the most recent verification row for a phone as verified.
func (r *OTPRepository) MarkVerified(phone string) error {
	now := time.Now()
	return r.db.Model(&models.PhoneVerification{}).
		Where("phone = ? AND verified_at IS NULL", phone).
		Order("created_at desc").Limit(1).
		Update("verified_at", now).Error
}
