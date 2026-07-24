package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// OTPHandler exposes phone-verification endpoints.
type OTPHandler struct {
	otp *service.OTPService
}

func NewOTPHandler(otp *service.OTPService) *OTPHandler {
	return &OTPHandler{otp: otp}
}

// Request handles POST /otp/request.
func (h *OTPHandler) Request(c echo.Context) error {
	var req models.OTPRequestRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.otp.Request(c.Request().Context(), req, c.RealIP()); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"sent": true})
}

// Verify handles POST /otp/verify. Auth is optional: when present, the verified
// phone is attached to the user.
func (h *OTPHandler) Verify(c echo.Context) error {
	var req models.OTPVerifyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.otp.Verify(c.Request().Context(), req, middleware.UserID(c)); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"verified": true})
}
