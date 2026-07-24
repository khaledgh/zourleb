package models

import "time"

// Booking status values.
const (
	BookingPending         = "pending"
	BookingAwaitingPayment = "awaiting_payment"
	BookingConfirmed       = "confirmed"
	BookingCancelled       = "cancelled"
	BookingCompleted       = "completed"
)

// Payment status values (shared across booking/order/boost ledgers).
const (
	PayStatusPending  = "pending"
	PayStatusPaid     = "paid"
	PayStatusFailed   = "failed"
	PayStatusRefunded = "refunded"
)

// Booking is a tourist's reservation of a tour departure.
type Booking struct {
	SoftBase
	Code                 string  `gorm:"size:16;uniqueIndex" json:"code"`
	UserID               uint    `gorm:"index" json:"user_id"`
	TourID               uint    `gorm:"index" json:"tour_id"`
	DepartureID          *uint   `gorm:"index" json:"departure_id"`
	Status               string  `gorm:"size:20;default:pending;index" json:"status"`
	TravelersCount       int     `gorm:"default:1" json:"travelers_count"`
	Subtotal             float64 `gorm:"type:decimal(14,2)" json:"subtotal"`
	Currency             string  `gorm:"size:3" json:"currency"`
	ContactPhone         string  `gorm:"size:24" json:"contact_phone"`
	ContactPhoneVerified bool    `gorm:"default:false" json:"contact_phone_verified"`
	PaymentStatus        string  `gorm:"size:16;default:pending" json:"payment_status"`
	Notes                string  `gorm:"type:text" json:"notes"`

	Travelers []BookingTraveler `gorm:"foreignKey:BookingID" json:"travelers,omitempty"`
	Payments  []BookingPayment  `gorm:"foreignKey:BookingID" json:"payments,omitempty"`
}

type BookingTraveler struct {
	Base
	BookingID      uint       `gorm:"index" json:"booking_id"`
	FullName       string     `gorm:"size:160" json:"full_name"`
	TravelerType   string     `gorm:"size:16;default:adult" json:"traveler_type"`
	Phone          string     `gorm:"size:24" json:"phone"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
}

type BookingPayment struct {
	Base
	BookingID      uint    `gorm:"index" json:"booking_id"`
	Provider       string  `gorm:"size:32" json:"provider"`
	ProviderRef    string  `gorm:"size:128" json:"provider_ref"`
	Amount         float64 `gorm:"type:decimal(14,2)" json:"amount"`
	Currency       string  `gorm:"size:3" json:"currency"`
	Status         string  `gorm:"size:16;default:pending" json:"status"`
	IdempotencyKey string  `gorm:"size:64;uniqueIndex" json:"-"`
}
