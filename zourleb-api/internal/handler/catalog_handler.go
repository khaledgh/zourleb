package handler

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/middleware"
	"github.com/zourleb/zourleb-api/internal/models"
	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/pagination"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// CatalogHandler serves the public catalog and home aggregation.
type CatalogHandler struct {
	catalog *service.CatalogService
}

func NewCatalogHandler(catalog *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalog: catalog}
}

// ListTours handles GET /tours.
func (h *CatalogHandler) ListTours(c echo.Context) error {
	f := parseTourFilter(c)
	p := pagination.FromQuery(c)
	cards, meta, err := h.catalog.ListTours(f, p, middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OKMeta(c, cards, meta)
}

// GetTour handles GET /tours/:slug.
func (h *CatalogHandler) GetTour(c echo.Context) error {
	detail, err := h.catalog.TourDetail(c.Param("slug"), middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, detail)
}

// Home handles GET /home.
func (h *CatalogHandler) Home(c echo.Context) error {
	payload, err := h.catalog.Home(middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, payload)
}

// Categories handles GET /categories.
func (h *CatalogHandler) Categories(c echo.Context) error {
	rows, err := h.catalog.Categories(middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

// Regions handles GET /regions.
func (h *CatalogHandler) Regions(c echo.Context) error {
	rows, err := h.catalog.Regions(middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, rows)
}

// GetAgency handles GET /agencies/:slug.
func (h *CatalogHandler) GetAgency(c echo.Context) error {
	a, err := h.catalog.AgencyDetail(c.Param("slug"), middleware.Locale(c))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, a)
}

// parseTourFilter extracts filter params from the query string.
func parseTourFilter(c echo.Context) models.TourFilter {
	f := models.TourFilter{
		Type:     c.QueryParam("type"),
		Currency: c.QueryParam("currency"),
		Search:   c.QueryParam("q"),
		Sort:     c.QueryParam("sort"),
	}
	if v := c.QueryParam("region_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(n)
			f.RegionID = &id
		}
	}
	if v := c.QueryParam("category_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(n)
			f.CategoryID = &id
		}
	}
	if v := c.QueryParam("featured"); v != "" {
		b := v == "true" || v == "1"
		f.Featured = &b
	}
	if v := c.QueryParam("price_min"); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceMin = &n
		}
	}
	if v := c.QueryParam("price_max"); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceMax = &n
		}
	}
	if v := c.QueryParam("date_from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			f.DateFrom = &t
		}
	}
	return f
}
