package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// BoostHandler exposes the agency-facing boost endpoints.
type BoostHandler struct {
	boosts *service.BoostService
}

func NewBoostHandler(boosts *service.BoostService) *BoostHandler {
	return &BoostHandler{boosts: boosts}
}

// Packages handles GET /agency/boost-packages.
func (h *BoostHandler) Packages(c echo.Context) error {
	rows, err := h.boosts.Packages()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

// Create handles POST /agency/boosts.
func (h *BoostHandler) Create(c echo.Context) error {
	var req models.CreateBoostRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	b, err := h.boosts.Create(middleware.AgencyID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, b)
}

// List handles GET /agency/boosts.
func (h *BoostHandler) List(c echo.Context) error {
	rows, err := h.boosts.List(middleware.AgencyID(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

// Pay handles POST /agency/boosts/:id/pay.
func (h *BoostHandler) Pay(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	var req models.PayRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.boosts.Pay(c.Request().Context(), middleware.AgencyID(c), id, req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}
