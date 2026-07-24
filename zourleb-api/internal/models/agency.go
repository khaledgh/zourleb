package models

import "time"

// Agency status values.
const (
	AgencyStatusPending   = "pending"
	AgencyStatusApproved  = "approved"
	AgencyStatusSuspended = "suspended"
)

// Agency is a travel agency tenant.
type Agency struct {
	SoftBase
	Slug            string  `gorm:"size:160;uniqueIndex" json:"slug"`
	Name            string  `gorm:"size:160" json:"name"`
	Logo            string  `gorm:"size:512" json:"logo"`
	Cover           string  `gorm:"size:512" json:"cover"`
	Phone           string  `gorm:"size:24" json:"phone"`
	Email           string  `gorm:"size:191" json:"email"`
	RegionID        *uint   `gorm:"index" json:"region_id"`
	Website         string  `gorm:"size:255" json:"website"`
	Verified        bool    `gorm:"default:false" json:"verified"`
	Status          string  `gorm:"size:16;default:pending;index" json:"status"`
	SubscriptionTier string `gorm:"size:32;default:basic" json:"subscription_tier"`
	CommissionRate  float64 `gorm:"type:decimal(5,2);default:10.00" json:"commission_rate"`
	RatingAvg       float64 `gorm:"type:decimal(3,2);default:0" json:"rating_avg"`
	RatingCount     int     `gorm:"default:0" json:"rating_count"`
	CreatedBy       uint    `json:"created_by"`

	Translations []AgencyTranslation `gorm:"foreignKey:AgencyID" json:"translations,omitempty"`
	Members      []AgencyMember      `gorm:"foreignKey:AgencyID" json:"-"`
}

func (a *Agency) IsApproved() bool { return a.Status == AgencyStatusApproved }

// AgencyTranslation holds locale-specific descriptive copy.
type AgencyTranslation struct {
	Base
	AgencyID    uint   `gorm:"index;uniqueIndex:uq_agency_locale" json:"agency_id"`
	Locale      string `gorm:"size:8;uniqueIndex:uq_agency_locale" json:"locale"`
	Description string `gorm:"type:text" json:"description"`
	About       string `gorm:"type:text" json:"about"`
}

// AgencyMember links a user to an agency with a role (staff invites).
type AgencyMember struct {
	Base
	AgencyID  uint       `gorm:"index;uniqueIndex:uq_agency_user" json:"agency_id"`
	UserID    uint       `gorm:"index;uniqueIndex:uq_agency_user" json:"user_id"`
	RoleID    uint       `json:"role_id"`
	InvitedAt *time.Time `json:"invited_at"`
	JoinedAt  *time.Time `json:"joined_at"`
}
