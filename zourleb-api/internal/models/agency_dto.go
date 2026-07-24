package models

import "time"

// --- Agency profile ---

type UpdateAgencyRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=160"`
	Phone    *string `json:"phone" validate:"omitempty,lb_phone"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Website  *string `json:"website" validate:"omitempty,max=255"`
	RegionID *uint   `json:"region_id"`
	Logo     *string `json:"logo" validate:"omitempty,max=512"`
	Cover    *string `json:"cover" validate:"omitempty,max=512"`
	Translations []AgencyTranslationInput `json:"translations" validate:"omitempty,dive"`
}

type AgencyTranslationInput struct {
	Locale      string `json:"locale" validate:"required,max=8"`
	Description string `json:"description" validate:"omitempty,max=5000"`
	About       string `json:"about" validate:"omitempty,max=5000"`
}

// ApplyAgencyRequest is a user's request to register a new agency.
type ApplyAgencyRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=160"`
	Phone    string `json:"phone" validate:"required,lb_phone"`
	Email    string `json:"email" validate:"required,email"`
	RegionID *uint  `json:"region_id"`
}

// --- Tour authoring ---

type SaveTourRequest struct {
	CategoryID   *uint   `json:"category_id"`
	RegionID     *uint   `json:"region_id"`
	Type         string  `json:"type" validate:"required,oneof=day_trip multi_day"`
	DurationDays int     `json:"duration_days" validate:"min=1"`
	Difficulty   string  `json:"difficulty" validate:"omitempty,max=16"`
	MinAge       int     `json:"min_age" validate:"min=0"`
	MaxCapacity  int     `json:"max_capacity" validate:"min=0"`
	BaseCurrency string  `json:"base_currency" validate:"omitempty,oneof=USD LBP"`
	IsShowcase   bool    `json:"is_showcase"`
	StartsFromDate *time.Time `json:"starts_from_date"`
	Translations []TourTranslationInput `json:"translations" validate:"required,min=1,dive"`
}

type TourTranslationInput struct {
	Locale      string `json:"locale" validate:"required,max=8"`
	Title       string `json:"title" validate:"required,max=200"`
	Summary     string `json:"summary" validate:"omitempty,max=500"`
	Description string `json:"description" validate:"omitempty"`
	Itinerary   string `json:"itinerary" validate:"omitempty"`
	Included    string `json:"included" validate:"omitempty"`
	Excluded    string `json:"excluded" validate:"omitempty"`
}

type AddDepartureRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	Capacity  int       `json:"capacity" validate:"required,min=1"`
}

type AddPriceRequest struct {
	DepartureID  *uint   `json:"departure_id"`
	TravelerType string  `json:"traveler_type" validate:"required,oneof=adult child infant student"`
	Currency     string  `json:"currency" validate:"required,oneof=USD LBP"`
	Amount       float64 `json:"amount" validate:"required,min=0"`
}

type AddImageRequest struct {
	URL     string `json:"url" validate:"required,max=512"`
	IsCover bool   `json:"is_cover"`
	Source  string `json:"source" validate:"omitempty,oneof=current past_gallery"`
}

// PublishTourRequest toggles a tour between draft and published.
type PublishTourRequest struct {
	Publish bool `json:"publish"`
}
