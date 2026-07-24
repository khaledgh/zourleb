package models

import "time"

// OTP channels and purposes.
const (
	OTPChannelWhatsApp = "whatsapp"
	OTPChannelSMS      = "sms"

	OTPPurposeRegister = "register"
	OTPPurposeBooking  = "booking"
)

// PhoneVerification is the durable audit trail for OTP requests. Live codes
// are kept in Redis with a TTL; this row records the attempt history.
type PhoneVerification struct {
	Base
	Phone       string     `gorm:"size:24;index" json:"phone"`
	Channel     string     `gorm:"size:16" json:"channel"`
	Purpose     string     `gorm:"size:16" json:"purpose"`
	CodeHash    string     `gorm:"size:64" json:"-"`
	Attempts    int        `gorm:"default:0" json:"attempts"`
	MaxAttempts int        `gorm:"default:3" json:"max_attempts"`
	ExpiresAt   time.Time  `json:"expires_at"`
	VerifiedAt  *time.Time `json:"verified_at"`
	IP          string     `gorm:"size:64" json:"-"`
}
