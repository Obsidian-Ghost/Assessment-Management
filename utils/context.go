package utils

import (
	"errors"

	"github.com/labstack/echo/v4"

	"assessment-management-system/models"
)

// ContextUserKey is the key used to store and retrieve the user from the context
const ContextUserKey = "user"

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(c echo.Context) (*models.User, error) {
	user, ok := c.Get(ContextUserKey).(*models.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
