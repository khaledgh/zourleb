package models

// CreateBoostRequest creates a boost for a tour using a package.
type CreateBoostRequest struct {
	TourID    uint `json:"tour_id" validate:"required"`
	PackageID uint `json:"package_id" validate:"required"`
}

// PayRequest initiates payment for a payable (boost or booking).
type PayRequest struct {
	Provider       string `json:"provider" validate:"omitempty"`
	IdempotencyKey string `json:"idempotency_key" validate:"required,max=64"`
	ReturnURL      string `json:"return_url" validate:"omitempty,max=512"`
}

// PayResponse is returned after initiating a payment.
type PayResponse struct {
	PaymentID   uint   `json:"payment_id"`
	Status      string `json:"status"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// SaveBoostPackageRequest is the admin package editor payload.
type SaveBoostPackageRequest struct {
	Name           string  `json:"name" validate:"required,max=120"`
	Placement      string  `json:"placement" validate:"required,oneof=home_banner featured_list category_top search_top"`
	DurationDays   int     `json:"duration_days" validate:"required,min=1"`
	Price          float64 `json:"price" validate:"required,min=0"`
	Currency       string  `json:"currency" validate:"required,oneof=USD LBP"`
	MaxActiveSlots int     `json:"max_active_slots" validate:"min=1"`
	Priority       int     `json:"priority"`
}
