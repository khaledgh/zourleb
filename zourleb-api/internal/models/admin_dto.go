package models

// --- Agencies ---

type ApproveAgencyRequest struct {
	Approve bool `json:"approve"`
}

// --- Languages ---

type SaveLanguageRequest struct {
	Code       string `json:"code" validate:"required,max=8"`
	Name       string `json:"name" validate:"required,max=64"`
	NativeName string `json:"native_name" validate:"required,max=64"`
	IsRTL      bool   `json:"is_rtl"`
	IsActive   bool   `json:"is_active"`
	SortOrder  int    `json:"sort_order"`
}

// --- Translations (UI strings) ---

type SaveTranslationRequest struct {
	Locale    string `json:"locale" validate:"required,max=8"`
	Namespace string `json:"namespace" validate:"required,max=64"`
	Key       string `json:"key" validate:"required,max=128"`
	Value     string `json:"value" validate:"required"`
}

// BulkTranslationsRequest imports many UI strings at once.
type BulkTranslationsRequest struct {
	Items []SaveTranslationRequest `json:"items" validate:"required,min=1,dive"`
}

// --- Settings ---

type SaveSettingRequest struct {
	Key   string `json:"key" validate:"required,max=64"`
	Value string `json:"value"`
	Type  string `json:"type" validate:"omitempty,oneof=string bool int json"`
}

// --- Banners ---

type SaveBannerRequest struct {
	Type      string `json:"type" validate:"omitempty,oneof=manual boosted"`
	TourID    *uint  `json:"tour_id"`
	Image     string `json:"image" validate:"required,max=512"`
	Title     string `json:"title" validate:"omitempty,max=200"`
	Link      string `json:"link" validate:"omitempty,max=512"`
	SortOrder int    `json:"sort_order"`
	Active    bool   `json:"active"`
}

// --- Categories & Regions ---

type SaveCategoryRequest struct {
	Slug         string                 `json:"slug" validate:"required,max=120"`
	Icon         string                 `json:"icon" validate:"omitempty,max=120"`
	SortOrder    int                    `json:"sort_order"`
	Translations []LocaleNameInput      `json:"translations" validate:"required,min=1,dive"`
}

type SaveRegionRequest struct {
	Slug         string            `json:"slug" validate:"required,max=120"`
	ParentID     *uint             `json:"parent_id"`
	Translations []LocaleNameInput `json:"translations" validate:"required,min=1,dive"`
}

type LocaleNameInput struct {
	Locale string `json:"locale" validate:"required,max=8"`
	Name   string `json:"name" validate:"required,max=120"`
}

// --- Reviews ---

type ModerateReviewRequest struct {
	Status string `json:"status" validate:"required,oneof=approved rejected"`
}

// --- Users ---

type SetUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active blocked"`
}

// --- Payments ---

type ConfirmPaymentRequest struct {
	PaymentID uint `json:"payment_id" validate:"required"`
}

// AnalyticsSummary is the platform dashboard payload.
type AnalyticsSummary struct {
	Users     int64              `json:"users"`
	Agencies  int64              `json:"agencies"`
	Tours     int64              `json:"tours"`
	Bookings  int64              `json:"bookings"`
	ActiveBoosts int64           `json:"active_boosts"`
	Revenue   map[string]float64 `json:"revenue"`
}
