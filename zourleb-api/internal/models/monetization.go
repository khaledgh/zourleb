package models

import "time"

// Boost placement slots.
const (
	PlacementHomeBanner   = "home_banner"
	PlacementFeaturedList = "featured_list"
	PlacementCategoryTop  = "category_top"
	PlacementSearchTop    = "search_top"
)

// Boost status values.
const (
	BoostPendingPayment = "pending_payment"
	BoostActive         = "active"
	BoostExpired        = "expired"
	BoostRejected       = "rejected"
)

// Banner types.
const (
	BannerBoosted = "boosted"
	BannerManual  = "manual"
)

// Payable / payer polymorphic types for the payments ledger.
const (
	PayableBoost        = "boost"
	PayableOrder        = "order"
	PayableBooking      = "booking"
	PayableSubscription = "subscription"

	PayerAgency = "agency"
	PayerUser   = "user"
)

// BoostPackage is an admin-defined paid placement product.
type BoostPackage struct {
	Base
	Name           string  `gorm:"size:120" json:"name"`
	Placement      string  `gorm:"size:32" json:"placement"`
	DurationDays   int     `json:"duration_days"`
	Price          float64 `gorm:"type:decimal(12,2)" json:"price"`
	Currency       string  `gorm:"size:3;default:USD" json:"currency"`
	MaxActiveSlots int     `gorm:"default:5" json:"max_active_slots"`
	Priority       int     `gorm:"default:0" json:"priority"`
}

// Boost is a purchased placement of a tour (or product) for a date range.
type Boost struct {
	Base
	AgencyID    uint       `gorm:"index" json:"agency_id"`
	TourID      *uint      `gorm:"index" json:"tour_id"`
	ProductID   *uint      `gorm:"index" json:"product_id"`
	PackageID   uint       `json:"package_id"`
	Placement   string     `gorm:"size:32;index" json:"placement"`
	Status      string     `gorm:"size:20;default:pending_payment;index" json:"status"`
	StartsAt    *time.Time `gorm:"index" json:"starts_at"`
	EndsAt      *time.Time `gorm:"index" json:"ends_at"`
	Impressions int64      `gorm:"default:0" json:"impressions"`
	Clicks      int64      `gorm:"default:0" json:"clicks"`
	PaymentID   *uint      `json:"payment_id"`
}

// Banner powers the home carousel. Boosted banners are auto-injected.
type Banner struct {
	Base
	Type      string     `gorm:"size:16;default:manual" json:"type"`
	TourID    *uint      `gorm:"index" json:"tour_id"`
	Image     string     `gorm:"size:512" json:"image"`
	Title     string     `gorm:"size:200" json:"title"`
	Link      string     `gorm:"size:512" json:"link"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	Active    bool       `gorm:"default:true;index" json:"active"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
}

// Payment is the polymorphic money ledger across all payable types.
type Payment struct {
	Base
	PayableType    string  `gorm:"size:20;index" json:"payable_type"`
	PayableID      uint    `gorm:"index" json:"payable_id"`
	PayerType      string  `gorm:"size:16" json:"payer_type"`
	PayerID        uint    `gorm:"index" json:"payer_id"`
	Amount         float64 `gorm:"type:decimal(14,2)" json:"amount"`
	Currency       string  `gorm:"size:3" json:"currency"`
	Provider       string  `gorm:"size:32" json:"provider"`
	ProviderRef    string  `gorm:"size:128" json:"provider_ref"`
	Status         string  `gorm:"size:16;default:pending;index" json:"status"`
	IdempotencyKey string  `gorm:"size:64;uniqueIndex" json:"-"`
}
