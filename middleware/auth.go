package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"assessment-management-system/models"
	"assessment-management-system/utils"
)

// AuthMiddleware returns a middleware function for JWT authentication
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header missing")
			}

			// Check if the Authorization header is in the correct format
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			// Extract the token
			tokenString := parts[1]
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Empty token")
			}

			// Validate the token
			claims, err := utils.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Create a user object from the claims
			user := &models.User{
				ID:             claims.UserID,
				OrganizationID: claims.OrganizationID,
				Email:          claims.Email,
				FirstName:      claims.FirstName,
				LastName:       claims.LastName,
				Role:           claims.Role,
			}

			// Store the user in the context
			c.Set("user", user)

			// Continue with the next handler
			return next(c)
		}
	}
}

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(c echo.Context) (*models.User, error) {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	return user, nil
}
