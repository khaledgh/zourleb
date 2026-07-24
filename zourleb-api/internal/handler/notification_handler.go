package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// NotificationHandler exposes in-app notification endpoints.
type NotificationHandler struct {
	notify *service.NotificationService
}

func NewNotificationHandler(notify *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notify: notify}
}

// List handles GET /notifications.
func (h *NotificationHandler) List(c echo.Context) error {
	rows, meta, err := h.notify.List(middleware.UserID(c), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

// MarkRead handles POST /notifications/:id/read.
func (h *NotificationHandler) MarkRead(c echo.Context) error {
	id, err := paramUint(c, "id")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.notify.MarkRead(middleware.UserID(c), id); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"read": true})
}
