package models

import "time"

// User status values.
const (
	UserStatusActive  = "active"
	UserStatusBlocked = "blocked"
)

// Well-known role keys.
const (
	RoleSuperAdmin  = "super_admin"
	RoleAgencyOwner = "agency_owner"
	RoleAgencyStaff = "agency_staff"
	RoleTourist     = "tourist"
)

// User is a platform account. Password may be null for OAuth-only accounts.
type User struct {
	SoftBase
	Name            string     `gorm:"size:120" json:"name"`
	Email           string     `gorm:"size:191;uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PasswordHash    *string    `gorm:"size:255" json:"-"`
	GoogleID        *string    `gorm:"size:128;index" json:"-"`
	Avatar          string     `gorm:"size:512" json:"avatar"`
	Phone           string     `gorm:"size:24;index" json:"phone"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	Locale          string     `gorm:"size:8;default:ar" json:"locale"`
	Status          string     `gorm:"size:16;default:active" json:"status"`
	LastLoginAt     *time.Time `json:"last_login_at"`

	Roles []UserRole `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) IsActive() bool        { return u.Status == UserStatusActive }
func (u *User) EmailVerified() bool   { return u.EmailVerifiedAt != nil }
func (u *User) PhoneVerified() bool   { return u.PhoneVerifiedAt != nil }

// Role groups a set of permissions under a stable key.
type Role struct {
	Base
	Key         string       `gorm:"size:32;uniqueIndex" json:"key"`
	Name        string       `gorm:"size:64" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission is a single capability such as "tour.create".
type Permission struct {
	Base
	Key string `gorm:"size:64;uniqueIndex" json:"key"`
}

// RolePermission is the explicit join (also covered by many2many above) so it
// can be migrated/seeded directly.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}

// UserRole assigns a role to a user, optionally scoped to one agency. Global
// roles (super_admin, tourist) leave AgencyID null.
type UserRole struct {
	Base
	UserID   uint  `gorm:"index;uniqueIndex:uq_user_role_agency" json:"user_id"`
	RoleID   uint  `gorm:"index;uniqueIndex:uq_user_role_agency" json:"role_id"`
	AgencyID *uint `gorm:"index;uniqueIndex:uq_user_role_agency" json:"agency_id"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}

// RefreshToken stores the hash of a rotating refresh token.
type RefreshToken struct {
	Base
	UserID    uint       `gorm:"index" json:"user_id"`
	TokenHash string     `gorm:"size:64;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	UserAgent string     `gorm:"size:255" json:"-"`
	IP        string     `gorm:"size:64" json:"-"`
}

func (t *RefreshToken) Active(now time.Time) bool {
	return t.RevokedAt == nil && t.ExpiresAt.After(now)
}
