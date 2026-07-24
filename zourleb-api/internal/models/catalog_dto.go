package models

import "time"

// TourFilter captures public catalog query parameters.
type TourFilter struct {
	RegionID   *uint
	CategoryID *uint
	Type       string
	Currency   string
	Featured   *bool
	PriceMin   *float64
	PriceMax   *float64
	DateFrom   *time.Time
	Search     string
	Sort       string // recent | price_asc | price_desc | featured
}

// TourCard is the list/grid representation (already localized).
type TourCard struct {
	ID         uint     `json:"id"`
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Cover      string   `json:"cover"`
	Type       string   `json:"type"`
	DurationDays int    `json:"duration_days"`
	Difficulty string   `json:"difficulty"`
	PriceFrom  float64  `json:"price_from"`
	Currency   string   `json:"currency"`
	Featured   bool     `json:"featured"`
	RegionID   *uint    `json:"region_id"`
	CategoryID *uint    `json:"category_id"`
	Agency     AgencyMini `json:"agency"`
	RatingAvg  float64  `json:"rating_avg"`
}

// AgencyMini is the compact agency reference embedded in cards.
type AgencyMini struct {
	ID       uint   `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Verified bool   `json:"verified"`
}

// TourDetail is the full localized tour representation.
type TourDetail struct {
	TourCard
	Description string             `json:"description"`
	Itinerary  string             `json:"itinerary"`
	Included   string             `json:"included"`
	Excluded   string             `json:"excluded"`
	MinAge     int                `json:"min_age"`
	MaxCapacity int               `json:"max_capacity"`
	Images     []TourImageDTO     `json:"images"`
	PastGallery []TourImageDTO    `json:"past_gallery"`
	Departures []TourDepartureDTO `json:"departures"`
	Prices     []TourPriceDTO     `json:"prices"`
}

type TourImageDTO struct {
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	IsCover   bool   `json:"is_cover"`
}

type TourDepartureDTO struct {
	ID         uint      `json:"id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Capacity   int       `json:"capacity"`
	SeatsLeft  int       `json:"seats_left"`
	Status     string    `json:"status"`
}

type TourPriceDTO struct {
	DepartureID  *uint   `json:"departure_id"`
	TravelerType string  `json:"traveler_type"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
}

// CategoryDTO / RegionDTO are localized lookup items.
type CategoryDTO struct {
	ID        uint   `json:"id"`
	Slug      string `json:"slug"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type RegionDTO struct {
	ID       uint   `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

// AgencyDetail is the public agency profile (localized).
type AgencyDetail struct {
	ID          uint    `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Logo        string  `json:"logo"`
	Cover       string  `json:"cover"`
	Description string  `json:"description"`
	About       string  `json:"about"`
	Phone       string  `json:"phone"`
	Website     string  `json:"website"`
	Verified    bool    `json:"verified"`
	RatingAvg   float64 `json:"rating_avg"`
	RatingCount int     `json:"rating_count"`
}

// HomePayload is the aggregated home-screen response.
type HomePayload struct {
	Banners      []BannerDTO   `json:"banners"`
	Featured     []TourCard    `json:"featured"`
	Categories   []CategoryDTO `json:"categories"`
	Regions      []RegionDTO   `json:"regions"`
	LastMinute   []TourCard    `json:"last_minute"`
	Modules      ModuleFlags   `json:"modules"`
}

type BannerDTO struct {
	ID     uint   `json:"id"`
	Image  string `json:"image"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	TourID *uint  `json:"tour_id"`
}

// ModuleFlags tells clients which optional modules are enabled.
type ModuleFlags struct {
	Shop    bool `json:"shop"`
	Reviews bool `json:"reviews"`
	Boost   bool `json:"boost"`
}
