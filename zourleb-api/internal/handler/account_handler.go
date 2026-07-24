package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// AccountHandler exposes the authenticated user's own profile endpoints.
type AccountHandler struct {
	account *service.AccountService
}

func NewAccountHandler(account *service.AccountService) *AccountHandler {
	return &AccountHandler{account: account}
}

// Me handles GET /me.
func (h *AccountHandler) Me(c echo.Context) error {
	res, err := h.account.Me(middleware.UserID(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// Update handles PATCH /me.
func (h *AccountHandler) Update(c echo.Context) error {
	var req models.UpdateMeRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.account.Update(middleware.UserID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// RegisterDevice handles POST /me/devices.
func (h *AccountHandler) RegisterDevice(c echo.Context) error {
	var req models.RegisterDeviceRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	if err := h.account.RegisterDevice(middleware.UserID(c), req); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"registered": true})
}
