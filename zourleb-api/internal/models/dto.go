package models

// Request/response DTOs. Kept separate from GORM entities so DB internals never
// leak to the API surface.

// --- Auth ---

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=120"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Locale   string `json:"locale" validate:"omitempty,max=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" validate:"required"`
	Locale  string `json:"locale" validate:"omitempty,max=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type EmailVerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"` // access token lifetime, seconds
}

type UserResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	Avatar        string   `json:"avatar"`
	Phone         string   `json:"phone"`
	Locale        string   `json:"locale"`
	EmailVerified bool     `json:"email_verified"`
	PhoneVerified bool     `json:"phone_verified"`
	Roles         []string `json:"roles"`
}

// NewUserResponse maps a User (and its role keys) to the API shape.
func NewUserResponse(u *User, roleKeys []string) UserResponse {
	return UserResponse{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		Avatar:        u.Avatar,
		Phone:         u.Phone,
		Locale:        u.Locale,
		EmailVerified: u.EmailVerified(),
		PhoneVerified: u.PhoneVerified(),
		Roles:         roleKeys,
	}
}

// --- OTP ---

type OTPRequestRequest struct {
	Phone   string `json:"phone" validate:"required,lb_phone"`
	Channel string `json:"channel" validate:"required,oneof=whatsapp sms"`
	Purpose string `json:"purpose" validate:"required,oneof=register booking"`
}

type OTPVerifyRequest struct {
	Phone string `json:"phone" validate:"required,lb_phone"`
	Code  string `json:"code" validate:"required,min=4,max=8"`
}

// --- Account ---

type UpdateMeRequest struct {
	Name   *string `json:"name" validate:"omitempty,min=2,max=120"`
	Locale *string `json:"locale" validate:"omitempty,max=8"`
	Avatar *string `json:"avatar" validate:"omitempty,max=512"`
}

type RegisterDeviceRequest struct {
	OneSignalPlayerID string `json:"onesignal_player_id" validate:"required,max=128"`
	Platform          string `json:"platform" validate:"required,oneof=ios android web"`
}

// --- Agency ---

type InviteAgencyMemberRequest struct {
	Email   string `json:"email" validate:"required,email"`
	RoleKey string `json:"role_key" validate:"required,oneof=agency_staff"`
}
