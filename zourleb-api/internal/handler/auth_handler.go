package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.auth.Register(c.Request().Context(), req, ac(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, res)
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.auth.Login(c.Request().Context(), req, ac(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// Google handles POST /auth/google.
func (h *AuthHandler) Google(c echo.Context) error {
	var req models.GoogleAuthRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.auth.Google(c.Request().Context(), req, ac(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req models.RefreshRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.auth.Refresh(c.Request().Context(), req.RefreshToken, ac(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// EmailVerify handles POST /auth/email/verify.
func (h *AuthHandler) EmailVerify(c echo.Context) error {
	var req models.EmailVerifyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.auth.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"verified": true})
}

// EmailResend handles POST /auth/email/resend for the authenticated user.
func (h *AuthHandler) EmailResend(c echo.Context) error {
	userID := middleware.UserID(c)
	if err := h.auth.RequestEmailVerification(c.Request().Context(), userID); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"sent": true})
}

// Logout handles POST /auth/logout.
func (h *AuthHandler) Logout(c echo.Context) error {
	var req models.RefreshRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.auth.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"logged_out": true})
}

// ac builds the auth context (UA + IP) from the request.
func ac(c echo.Context) service.AuthContext {
	return service.AuthContextFromRequest(c.Request().UserAgent(), c.RealIP())
}
