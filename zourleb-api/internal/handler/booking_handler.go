package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// BookingHandler exposes the tourist booking endpoints.
type BookingHandler struct {
	bookings *service.BookingService
}

func NewBookingHandler(bookings *service.BookingService) *BookingHandler {
	return &BookingHandler{bookings: bookings}
}

// Create handles POST /bookings.
func (h *BookingHandler) Create(c echo.Context) error {
	var req models.CreateBookingRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.bookings.Create(c.Request().Context(), middleware.UserID(c), req, middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, res)
}

// List handles GET /bookings.
func (h *BookingHandler) List(c echo.Context) error {
	p := pagination.FromQuery(c)
	rows, meta, err := h.bookings.List(middleware.UserID(c), p, middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

// Get handles GET /bookings/:code.
func (h *BookingHandler) Get(c echo.Context) error {
	res, err := h.bookings.Get(middleware.UserID(c), c.Param("code"), middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// Pay handles POST /bookings/:code/pay.
func (h *BookingHandler) Pay(c echo.Context) error {
	var req models.PayRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.bookings.Pay(c.Request().Context(), middleware.UserID(c), c.Param("code"), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, res)
}

// Cancel handles POST /bookings/:code/cancel.
func (h *BookingHandler) Cancel(c echo.Context) error {
	if err := h.bookings.Cancel(middleware.UserID(c), c.Param("code")); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"cancelled": true})
}
