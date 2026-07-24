package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// I18nHandler serves languages and UI translation bundles.
type I18nHandler struct {
	i18n *service.I18nService
}

func NewI18nHandler(i18n *service.I18nService) *I18nHandler {
	return &I18nHandler{i18n: i18n}
}

// Languages handles GET /languages.
func (h *I18nHandler) Languages(c echo.Context) error {
	langs, err := h.i18n.Languages()
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, langs)
}

// Bundle handles GET /i18n/:locale — the UI translation bundle.
func (h *I18nHandler) Bundle(c echo.Context) error {
	bundle, err := h.i18n.Bundle(c.Param("locale"))
	if err != nil {
		return response.Fail(c, err)
	}
	return response.OK(c, bundle)
}
