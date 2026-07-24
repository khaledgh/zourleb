package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/service"
)

// RequestIDMiddleware assigns a correlation id per request (honoring an
// inbound X-Request-ID) and echoes it in the response header.
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Request().Header.Get("X-Request-ID")
			if rid == "" {
				rid = uuid.NewString()
			}
			c.Set(CtxRequestID, rid)
			c.Response().Header().Set("X-Request-ID", rid)
			return next(c)
		}
	}
}

// LocaleMiddleware resolves the request locale from ?locale or Accept-Language,
// validating it against active languages, and stores it on the context.
func LocaleMiddleware(i18n *service.I18nService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requested := c.QueryParam("locale")
			if requested == "" {
				requested = parseAcceptLanguage(c.Request().Header.Get("Accept-Language"))
			}
			c.Set(CtxLocale, i18n.ResolveLocale(requested))
			return next(c)
		}
	}
}

// parseAcceptLanguage returns the primary language subtag of the first entry,
// e.g. "fr-CA,fr;q=0.9,en;q=0.8" -> "fr".
func parseAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}
	first := strings.Split(header, ",")[0]
	first = strings.Split(first, ";")[0]
	first = strings.TrimSpace(first)
	if i := strings.Index(first, "-"); i > 0 {
		return strings.ToLower(first[:i])
	}
	return strings.ToLower(first)
}
