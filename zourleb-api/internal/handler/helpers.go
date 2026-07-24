// Package handler holds Echo HTTP handlers: parse/validate, call service, format.
package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/pkg/response"
)

// bindAndValidate binds the request body into dst and runs validation,
// returning a domain error suitable for response.Fail on failure.
func bindAndValidate(c echo.Context, dst interface{}) error {
	if err := c.Bind(dst); err != nil {
		return response.ErrBadRequest.WithMessage("Malformed request body.")
	}
	if err := c.Validate(dst); err != nil {
		return err
	}
	return nil
}
