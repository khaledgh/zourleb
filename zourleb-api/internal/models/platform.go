package models

import (
	"time"

	"gorm.io/datatypes"
)

// Review status values.
const (
	ReviewPending  = "pending"
	ReviewApproved = "approved"
	ReviewRejected = "rejected"
)

// Setting value types.
const (
	SettingString = "string"
	SettingBool   = "bool"
	SettingInt    = "int"
	SettingJSON   = "json"
)

// Well-known setting keys (feature flags & platform defaults).
const (
	SettingShopEnabled    = "shop.enabled"
	SettingReviewsEnabled = "reviews.enabled"
	SettingBoostEnabled   = "boost.enabled"
	SettingDefaultCurrency = "default_currency"
	SettingBookingOTP     = "booking.require_otp"
)

// Language is a supported UI/content language, managed by the super admin.
type Language struct {
	Base
	Code       string `gorm:"size:8;uniqueIndex" json:"code"`
	Name       string `gorm:"size:64" json:"name"`
	NativeName string `gorm:"size:64" json:"native_name"`
	IsRTL      bool   `gorm:"default:false" json:"is_rtl"`
	IsActive   bool   `gorm:"default:true;index" json:"is_active"`
	IsDefault  bool   `gorm:"default:false" json:"is_default"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

// Translation is a single UI string keyed by (locale, namespace, key).
type Translation struct {
	Base
	Locale    string `gorm:"size:8;uniqueIndex:uq_translation" json:"locale"`
	Namespace string `gorm:"size:64;uniqueIndex:uq_translation" json:"namespace"`
	Key       string `gorm:"size:128;uniqueIndex:uq_translation" json:"key"`
	Value     string `gorm:"type:text" json:"value"`
}

// Setting is a typed key/value used for feature flags and platform config.
type Setting struct {
	Key   string `gorm:"size:64;primaryKey" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Type  string `gorm:"size:16;default:string" json:"type"`
}

// Review is a verified-booking rating of a tour.
type Review struct {
	Base
	UserID    uint   `gorm:"index" json:"user_id"`
	TourID    uint   `gorm:"index" json:"tour_id"`
	BookingID *uint  `json:"booking_id"`
	Rating    int    `json:"rating"`
	Comment   string `gorm:"type:text" json:"comment"`
	Status    string `gorm:"size:16;default:pending;index" json:"status"`
}

// Favorite is a tourist's saved tour.
type Favorite struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	TourID    uint      `gorm:"primaryKey" json:"tour_id"`
	CreatedAt time.Time `json:"created_at"`
}

// DeviceToken maps a user to a OneSignal player id.
type DeviceToken struct {
	Base
	UserID          uint   `gorm:"index" json:"user_id"`
	OneSignalPlayer string `gorm:"size:128;uniqueIndex" json:"onesignal_player_id"`
	Platform        string `gorm:"size:16" json:"platform"` // ios/android/web
}

// Notification is an in-app message persisted alongside push delivery.
type Notification struct {
	Base
	UserID uint           `gorm:"index" json:"user_id"`
	Type   string         `gorm:"size:48" json:"type"`
	Title  string         `gorm:"size:200" json:"title"`
	Body   string         `gorm:"type:text" json:"body"`
	Data   datatypes.JSON `json:"data"`
	ReadAt *time.Time     `json:"read_at"`
}

// AuditLog records sensitive actions for accountability.
type AuditLog struct {
	Base
	ActorID    *uint          `gorm:"index" json:"actor_id"`
	Action     string         `gorm:"size:64;index" json:"action"`
	TargetType string         `gorm:"size:48" json:"target_type"`
	TargetID   uint           `json:"target_id"`
	Payload    datatypes.JSON `json:"payload"`
	IP         string         `gorm:"size:64" json:"ip"`
}
