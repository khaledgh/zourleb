package middleware

import "github.com/labstack/echo/v4"

// Context keys used across middleware and handlers.
const (
	CtxUserID    = "user_id"
	CtxUserEmail = "user_email"
	CtxLocale    = "locale"
	CtxRequestID = "request_id"
	CtxAgencyID  = "agency_id" // resolved agency scope for agency-portal routes
)

// UserID returns the authenticated user's id, or 0 if unauthenticated.
func UserID(c echo.Context) uint {
	if v, ok := c.Get(CtxUserID).(uint); ok {
		return v
	}
	return 0
}

// Locale returns the resolved request locale (always set by Locale middleware).
func Locale(c echo.Context) string {
	if v, ok := c.Get(CtxLocale).(string); ok && v != "" {
		return v
	}
	return "ar"
}

// RequestID returns the per-request correlation id.
func RequestID(c echo.Context) string {
	if v, ok := c.Get(CtxRequestID).(string); ok {
		return v
	}
	return ""
}

// AgencyID returns the resolved agency scope, or 0 if none.
func AgencyID(c echo.Context) uint {
	if v, ok := c.Get(CtxAgencyID).(uint); ok {
		return v
	}
	return 0
}
