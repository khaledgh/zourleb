// Package models holds GORM entities and shared model types.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Base is embedded in entities that need timestamps and an auto ID.
type Base struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SoftBase adds soft-delete semantics on top of Base.
type SoftBase struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AllModels returns every entity for AutoMigrate. Keep this in sync as the
// schema grows; ordering matters only for FK constraints GORM can't infer.
func AllModels() []interface{} {
	return []interface{}{
		// identity & access
		&User{}, &Role{}, &Permission{}, &RolePermission{}, &UserRole{},
		&RefreshToken{},
		// agencies
		&Agency{}, &AgencyTranslation{}, &AgencyMember{}, &SubscriptionTier{},
		// catalog
		&Region{}, &RegionTranslation{},
		&Category{}, &CategoryTranslation{},
		&Tour{}, &TourTranslation{}, &TourImage{}, &TourDeparture{}, &TourPrice{},
		// bookings
		&Booking{}, &BookingTraveler{}, &BookingPayment{},
		// phone verification
		&PhoneVerification{},
		// shop
		&Product{}, &ProductTranslation{}, &ProductImage{}, &Order{}, &OrderItem{},
		// monetization
		&BoostPackage{}, &Boost{}, &Banner{}, &Payment{},
		// platform
		&Language{}, &Translation{}, &Setting{}, &Review{}, &Favorite{},
		&DeviceToken{}, &Notification{}, &AuditLog{},
	}
}
