package handler

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/pkg/response"
)

// HealthHandler exposes liveness and readiness probes.
type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Live handles GET /health — process is up.
func (h *HealthHandler) Live(c echo.Context) error {
	return response.OK(c, map[string]string{"status": "ok"})
}

// Ready handles GET /ready — dependencies (DB) are reachable.
func (h *HealthHandler) Ready(c echo.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return response.Fail(c, response.ErrInternal.Wrap(err))
	}
	if err := sqlDB.Ping(); err != nil {
		return response.Fail(c, response.ErrInternal.WithMessage("database unreachable"))
	}
	return response.OK(c, map[string]string{"status": "ready"})
}
