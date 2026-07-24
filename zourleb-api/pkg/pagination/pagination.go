// Package pagination provides request parsing and a response meta block.
package pagination

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// Params holds parsed page/per_page values.
type Params struct {
	Page    int
	PerPage int
}

// Offset is the SQL offset for the current page.
func (p Params) Offset() int { return (p.Page - 1) * p.PerPage }

// Limit is the SQL limit for the current page.
func (p Params) Limit() int { return p.PerPage }

// FromQuery parses ?page & ?per_page with clamping and sane defaults.
func FromQuery(c echo.Context) Params {
	page := atoiDefault(c.QueryParam("page"), defaultPage)
	if page < 1 {
		page = defaultPage
	}
	per := atoiDefault(c.QueryParam("per_page"), defaultPerPage)
	if per < 1 {
		per = defaultPerPage
	}
	if per > maxPerPage {
		per = maxPerPage
	}
	return Params{Page: page, PerPage: per}
}

// Meta is the pagination block placed in the response envelope's meta field.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta builds a Meta from params and a total count.
func NewMeta(p Params, total int64) Meta {
	pages := 0
	if p.PerPage > 0 {
		pages = int((total + int64(p.PerPage) - 1) / int64(p.PerPage))
	}
	return Meta{Page: p.Page, PerPage: p.PerPage, Total: total, TotalPages: pages}
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
