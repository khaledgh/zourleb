package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// EngagementHandler exposes favorites and reviews endpoints.
type EngagementHandler struct {
	svc *service.EngagementService
}

func NewEngagementHandler(svc *service.EngagementService) *EngagementHandler {
	return &EngagementHandler{svc: svc}
}

// AddFavorite handles POST /favorites/:tourId.
func (h *EngagementHandler) AddFavorite(c echo.Context) error {
	tourID, err := paramUint(c, "tourId")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.svc.AddFavorite(middleware.UserID(c), tourID); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"favorited": true})
}

// RemoveFavorite handles DELETE /favorites/:tourId.
func (h *EngagementHandler) RemoveFavorite(c echo.Context) error {
	tourID, err := paramUint(c, "tourId")
	if err != nil {
		return response.Fail(c, response.ErrBadRequest)
	}
	if err := h.svc.RemoveFavorite(middleware.UserID(c), tourID); err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string]bool{"removed": true})
}

// ListFavorites handles GET /favorites.
func (h *EngagementHandler) ListFavorites(c echo.Context) error {
	ids, err := h.svc.Favorites(middleware.UserID(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, map[string][]uint{"tour_ids": ids})
}

// CreateReview handles POST /reviews.
func (h *EngagementHandler) CreateReview(c echo.Context) error {
	var req models.CreateReviewRequest
	if err := bindAndValidate(c, &req); err != nil {
		return response.Fail(c, err)
	}
	res, err := h.svc.CreateReview(middleware.UserID(c), req)
	if err != nil {
		return response.Fail(c, err)
	}
	return response.Created(c, res)
}

// TourReviews handles GET /tours/:slug/reviews — but uses tour id via query for
// simplicity here; mounted as GET /reviews?tour_id=.
func (h *EngagementHandler) TourReviews(c echo.Context) error {
	tourID, _ := strconv.ParseUint(c.QueryParam("tour_id"), 10, 64)
	rows, meta, err := h.svc.TourReviews(uint(tourID), pagination.FromQuery(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, rows, meta)
}

// paramUint parses a uint path param.
func paramUint(c echo.Context, name string) (uint, error) {
	n, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
