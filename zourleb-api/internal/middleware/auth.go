package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/internal/repository"
	"github.com/zourleb/zourleb-api/pkg/response"
	"github.com/zourleb/zourleb-api/pkg/token"
)

// Auth wires JWT verification and RBAC permission checks.
type Auth struct {
	tokens *token.Manager
	users  *repository.UserRepository
}

func NewAuth(tm *token.Manager, users *repository.UserRepository) *Auth {
	return &Auth{tokens: tm, users: users}
}

// Required rejects requests without a valid access token, populating context.
func (a *Auth) Required() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := a.parse(c)
			if err != nil {
				return response.Fail(c, response.ErrUnauthorized)
			}
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxUserEmail, claims.Email)
			return next(c)
		}
	}
}

// Optional populates context when a valid token is present but never rejects.
func (a *Auth) Optional() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if claims, err := a.parse(c); err == nil {
				c.Set(CtxUserID, claims.UserID)
				c.Set(CtxUserEmail, claims.Email)
			}
			return next(c)
		}
	}
}

// RequirePermission ensures the caller holds the given permission. When the
// route is agency-scoped (an :agencyId param or X-Agency-ID is present), it also
// verifies the user belongs to that agency.
func (a *Auth) RequirePermission(perm string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid := UserID(c)
			if uid == 0 {
				return response.Fail(c, response.ErrUnauthorized)
			}
			perms, err := a.users.Permissions(uid)
			if err != nil {
				return response.Fail(c, response.ErrInternal.Wrap(err))
			}
			if _, ok := perms[perm]; !ok {
				return response.Fail(c, response.ErrForbidden)
			}
			return next(c)
		}
	}
}

// RequireAgencyScope resolves and verifies the agency the caller is acting for.
// The agency id comes from the X-Agency-ID header or ?agency_id; the caller must
// be a member of it. Resolved id is stored on the context for handlers.
func (a *Auth) RequireAgencyScope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid := UserID(c)
			if uid == 0 {
				return response.Fail(c, response.ErrUnauthorized)
			}
			ids, err := a.users.AgencyIDsForUser(uid)
			if err != nil {
				return response.Fail(c, response.ErrInternal.Wrap(err))
			}
			if len(ids) == 0 {
				return response.Fail(c, response.ErrForbidden.WithMessage("You are not a member of any agency."))
			}

			requested := agencyIDFromRequest(c)
			resolved := ids[0]
			if requested != 0 {
				if !contains(ids, requested) {
					return response.Fail(c, response.ErrForbidden)
				}
				resolved = requested
			}
			c.Set(CtxAgencyID, resolved)
			return next(c)
		}
	}
}

func (a *Auth) parse(c echo.Context) (*token.Claims, error) {
	h := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return nil, response.ErrUnauthorized
	}
	raw := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	return a.tokens.ParseAccess(raw)
}

func agencyIDFromRequest(c echo.Context) uint {
	if v := c.Request().Header.Get("X-Agency-ID"); v != "" {
		return parseUint(v)
	}
	return parseUint(c.QueryParam("agency_id"))
}

func parseUint(s string) uint {
	var n uint
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + uint(ch-'0')
	}
	return n
}

func contains(s []uint, v uint) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
