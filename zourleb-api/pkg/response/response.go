package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Envelope is the unified API response shape returned by every endpoint.
//
//	{ "success": bool, "data": ..., "error": {...}, "meta": {...} }
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// APIError carries a stable machine-readable code plus a localized message.
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Fields  interface{} `json:"fields,omitempty"` // field-level validation errors
}

// OK writes a 200 success envelope.
func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

// OKMeta writes a 200 success envelope with pagination/meta.
func OKMeta(c echo.Context, data interface{}, meta interface{}) error {
	return c.JSON(http.StatusOK, Envelope{Success: true, Data: data, Meta: meta})
}

// Created writes a 201 success envelope.
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

// NoContent writes a 204.
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Fail writes an error envelope from an *Error (or coerces a generic error).
func Fail(c echo.Context, err error) error {
	e := AsError(err)
	return c.JSON(e.HTTPStatus, Envelope{
		Success: false,
		Error:   &APIError{Code: e.Code, Message: e.Message, Fields: e.Fields},
	})
}
