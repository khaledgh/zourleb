package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/service"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// RequireFeature blocks a route group when its settings flag is disabled,
// returning FEATURE_DISABLED so clients can hide the module.
func RequireFeature(settings *service.SettingsService, key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settings.Bool(key, false) {
				return response.Fail(c, response.ErrFeatureOff)
			}
			return next(c)
		}
	}
}
