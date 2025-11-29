package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/models"
)

// RoleMiddleware checks if the user has one of the required roles
func RoleMiddleware(allowedRoles ...models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(models.UserRole)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "user role not found")
			}

			// Check if user's role is in the allowed roles
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}
