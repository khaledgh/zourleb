package response

import (
	"errors"
	"net/http"
)

// Error is a domain error carrying a stable code, an HTTP status, and a
// (server-side localized) message. Services return these; handlers pass them
// to response.Error which renders the envelope.
type Error struct {
	Code       string
	Message    string
	HTTPStatus int
	Fields     interface{}
	cause      error
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) WithFields(f interface{}) *Error {
	clone := *e
	clone.Fields = f
	return &clone
}

func (e *Error) WithMessage(msg string) *Error {
	clone := *e
	clone.Message = msg
	return &clone
}

func (e *Error) Wrap(cause error) *Error {
	clone := *e
	clone.cause = cause
	return &clone
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil {
		return e.Code + ": " + e.cause.Error()
	}
	return e.Code + ": " + e.Message
}

// New constructs a reusable error template.
func New(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, HTTPStatus: status}
}

// AsError coerces any error into a *Error, defaulting to INTERNAL.
func AsError(err error) *Error {
	if err == nil {
		return ErrInternal
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return ErrInternal.Wrap(err)
}

// Catalog of stable error codes. Messages are English defaults; the i18n layer
// may localize them based on Accept-Language before the response is written.
var (
	ErrInternal       = New("INTERNAL", "Something went wrong.", http.StatusInternalServerError)
	ErrValidation     = New("VALIDATION_FAILED", "The request is invalid.", http.StatusUnprocessableEntity)
	ErrBadRequest     = New("BAD_REQUEST", "The request could not be processed.", http.StatusBadRequest)
	ErrUnauthorized   = New("UNAUTHORIZED", "Authentication is required.", http.StatusUnauthorized)
	ErrForbidden      = New("FORBIDDEN", "You do not have permission to do that.", http.StatusForbidden)
	ErrNotFound       = New("NOT_FOUND", "The requested resource was not found.", http.StatusNotFound)
	ErrConflict       = New("CONFLICT", "The resource already exists.", http.StatusConflict)
	ErrTooManyRequest = New("RATE_LIMITED", "Too many requests. Please slow down.", http.StatusTooManyRequests)
	ErrFeatureOff     = New("FEATURE_DISABLED", "This feature is not available.", http.StatusNotFound)

	// Auth
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "Email or password is incorrect.", http.StatusUnauthorized)
	ErrEmailTaken         = New("EMAIL_TAKEN", "That email is already registered.", http.StatusConflict)
	ErrTokenInvalid       = New("TOKEN_INVALID", "The token is invalid or expired.", http.StatusUnauthorized)
	ErrEmailNotVerified   = New("EMAIL_NOT_VERIFIED", "Please verify your email first.", http.StatusForbidden)
	ErrAccountBlocked     = New("ACCOUNT_BLOCKED", "This account is blocked.", http.StatusForbidden)

	// OTP / phone
	ErrOTPInvalid    = New("OTP_INVALID", "The code is incorrect.", http.StatusUnprocessableEntity)
	ErrOTPExpired    = New("OTP_EXPIRED", "The code has expired.", http.StatusUnprocessableEntity)
	ErrOTPMaxAttempt = New("OTP_MAX_ATTEMPTS", "Too many attempts. Request a new code.", http.StatusTooManyRequests)
	ErrOTPCooldown   = New("OTP_COOLDOWN", "Please wait before requesting another code.", http.StatusTooManyRequests)
	ErrPhoneInvalid  = New("PHONE_INVALID", "The phone number is invalid.", http.StatusUnprocessableEntity)
	ErrPhoneNotVerified = New("PHONE_NOT_VERIFIED", "The contact phone must be verified.", http.StatusForbidden)

	// Domain
	ErrTourNotFound    = New("TOUR_NOT_FOUND", "Tour not found.", http.StatusNotFound)
	ErrAgencyNotFound  = New("AGENCY_NOT_FOUND", "Agency not found.", http.StatusNotFound)
	ErrAgencyPending   = New("AGENCY_PENDING", "Your agency is awaiting approval.", http.StatusForbidden)
	ErrBookingNotFound = New("BOOKING_NOT_FOUND", "Booking not found.", http.StatusNotFound)
	ErrDepartureFull   = New("DEPARTURE_FULL", "This departure is fully booked.", http.StatusConflict)
)
