package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"assessment-management-system/models"
)

// RoleBasedAccessControl returns a middleware function for RBAC
func RoleBasedAccessControl(roles ...models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the user from the context
			user, err := GetUserFromContext(c)
			if err != nil {
				return err
			}

			// Check if the user has one of the required roles
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to access this resource")
			}

			// Continue with the next handler
			return next(c)
		}
	}
}

// AdminOnly is a shorthand for RoleBasedAccessControl with only admin role
func AdminOnly() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleAdmin)
}

// TeacherOrAdmin is a shorthand for RoleBasedAccessControl with teacher and admin roles
func TeacherOrAdmin() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleTeacher, models.RoleAdmin)
}

// TeacherOnly is a shorthand for RoleBasedAccessControl with only teacher role
func TeacherOnly() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleTeacher)
}

// StudentOnly is a shorthand for RoleBasedAccessControl with only student role
func StudentOnly() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleStudent)
}

// StudentOrTeacher is a shorthand for RoleBasedAccessControl with student and teacher roles
func StudentOrTeacher() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleStudent, models.RoleTeacher)
}

// AllRoles allows any authenticated user regardless of role
func AllRoles() echo.MiddlewareFunc {
	return RoleBasedAccessControl(models.RoleAdmin, models.RoleTeacher, models.RoleStudent)
}
