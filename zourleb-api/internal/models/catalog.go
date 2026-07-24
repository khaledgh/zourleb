package models

import "time"

// Tour status values.
const (
	TourStatusDraft     = "draft"
	TourStatusPending   = "pending"
	TourStatusPublished = "published"
	TourStatusArchived  = "archived"
)

// Tour types.
const (
	TourTypeDayTrip  = "day_trip"
	TourTypeMultiDay = "multi_day"
)

// Tour image sources.
const (
	ImageSourceCurrent = "current"
	ImageSourcePast    = "past_gallery"
)

// Departure status values.
const (
	DepartureOpen      = "open"
	DepartureFull      = "full"
	DepartureCancelled = "cancelled"
)

// Region is a Lebanese governorate/area; self-referential via ParentID.
type Region struct {
	Base
	Slug     string `gorm:"size:120;uniqueIndex" json:"slug"`
	ParentID *uint  `gorm:"index" json:"parent_id"`

	Translations []RegionTranslation `gorm:"foreignKey:RegionID" json:"translations,omitempty"`
}

type RegionTranslation struct {
	Base
	RegionID uint   `gorm:"index;uniqueIndex:uq_region_locale" json:"region_id"`
	Locale   string `gorm:"size:8;uniqueIndex:uq_region_locale" json:"locale"`
	Name     string `gorm:"size:120" json:"name"`
}

// Category groups tours (adventure, religious, ski, beach…).
type Category struct {
	Base
	Slug      string `gorm:"size:120;uniqueIndex" json:"slug"`
	Icon      string `gorm:"size:120" json:"icon"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`

	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID" json:"translations,omitempty"`
}

type CategoryTranslation struct {
	Base
	CategoryID uint   `gorm:"index;uniqueIndex:uq_category_locale" json:"category_id"`
	Locale     string `gorm:"size:8;uniqueIndex:uq_category_locale" json:"locale"`
	Name       string `gorm:"size:120" json:"name"`
}

// Tour is a publishable trip offering owned by an agency.
type Tour struct {
	SoftBase
	AgencyID      uint       `gorm:"index" json:"agency_id"`
	Slug          string     `gorm:"size:191;uniqueIndex" json:"slug"`
	CategoryID    *uint      `gorm:"index" json:"category_id"`
	RegionID      *uint      `gorm:"index" json:"region_id"`
	Type          string     `gorm:"size:16;default:day_trip" json:"type"`
	DurationDays  int        `gorm:"default:1" json:"duration_days"`
	Difficulty    string     `gorm:"size:16" json:"difficulty"`
	MinAge        int        `json:"min_age"`
	MaxCapacity   int        `json:"max_capacity"`
	BaseCurrency  string     `gorm:"size:3;default:USD" json:"base_currency"`
	PriceFrom     float64    `gorm:"type:decimal(12,2)" json:"price_from"`
	Status        string     `gorm:"size:16;default:draft;index" json:"status"`
	IsShowcase    bool       `gorm:"default:false;index" json:"is_showcase"` // past-tour portfolio
	Featured      bool       `gorm:"default:false;index" json:"featured"`
	StartsFromDate *time.Time `json:"starts_from_date"`
	CreatedBy     uint       `json:"created_by"`

	Translations []TourTranslation `gorm:"foreignKey:TourID" json:"translations,omitempty"`
	Images       []TourImage       `gorm:"foreignKey:TourID" json:"images,omitempty"`
	Departures   []TourDeparture   `gorm:"foreignKey:TourID" json:"departures,omitempty"`
	Prices       []TourPrice       `gorm:"foreignKey:TourID" json:"prices,omitempty"`
}

type TourTranslation struct {
	Base
	TourID          uint   `gorm:"index;uniqueIndex:uq_tour_locale" json:"tour_id"`
	Locale          string `gorm:"size:8;uniqueIndex:uq_tour_locale" json:"locale"`
	Title           string `gorm:"size:200" json:"title"`
	Summary         string `gorm:"size:500" json:"summary"`
	Description     string `gorm:"type:text" json:"description"`
	Itinerary       string `gorm:"type:longtext" json:"itinerary"` // JSON or long text
	Included        string `gorm:"type:text" json:"included"`
	Excluded        string `gorm:"type:text" json:"excluded"`
	MetaTitle       string `gorm:"size:200" json:"meta_title"`
	MetaDescription string `gorm:"size:300" json:"meta_description"`
}

type TourImage struct {
	Base
	TourID    uint   `gorm:"index" json:"tour_id"`
	URL       string `gorm:"size:512" json:"url"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsCover   bool   `gorm:"default:false" json:"is_cover"`
	Source    string `gorm:"size:16;default:current" json:"source"` // current | past_gallery
}

type TourDeparture struct {
	Base
	TourID      uint      `gorm:"index" json:"tour_id"`
	StartDate   time.Time `gorm:"index" json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Capacity    int       `json:"capacity"`
	SeatsBooked int       `gorm:"default:0" json:"seats_booked"`
	Status      string    `gorm:"size:16;default:open" json:"status"`
}

func (d *TourDeparture) SeatsLeft() int { return d.Capacity - d.SeatsBooked }

// TourPrice supports dual-currency and per-traveler-type pricing.
// DepartureID null means the price applies to all departures.
type TourPrice struct {
	Base
	TourID       uint    `gorm:"index" json:"tour_id"`
	DepartureID  *uint   `gorm:"index" json:"departure_id"`
	TravelerType string  `gorm:"size:16;default:adult" json:"traveler_type"` // adult/child/infant/student
	Currency     string  `gorm:"size:3" json:"currency"`                     // USD/LBP
	Amount       float64 `gorm:"type:decimal(14,2)" json:"amount"`
}
