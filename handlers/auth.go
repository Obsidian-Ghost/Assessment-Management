package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"assessment-management-system/models"
	"assessment-management-system/services"
	"assessment-management-system/utils"
)

// AuthHandler handles authentication-related routes
type AuthHandler struct {
	authService       *services.AuthService
	userService       *services.UserService
	jwtSecret         string
	jwtExpiration     time.Duration
	refreshExpiration time.Duration
	validator         *validator.Validate
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	authService *services.AuthService,
	userService *services.UserService,
	jwtSecret string,
	jwtExpiration time.Duration,
) *AuthHandler {
	// Refresh token valid for 7 days
	refreshExpiration := time.Hour * 24 * 7

	return &AuthHandler{
		authService:       authService,
		userService:       userService,
		jwtSecret:         jwtSecret,
		jwtExpiration:     jwtExpiration,
		refreshExpiration: refreshExpiration,
		validator:         utils.NewValidator(),
	}
}

// HandleLogin handles user login
func (h *AuthHandler) HandleLogin(c echo.Context) error {
	var loginReq models.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	user, err := h.authService.Authenticate(c.Request().Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate token pair (access token + refresh token)
	tokenResponse, err := h.authService.CreateTokenPair(
		c.Request().Context(),
		user,
		h.jwtSecret,
		h.jwtExpiration,
		h.refreshExpiration,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate tokens: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  tokenResponse.AccessToken,
		"refresh_token": tokenResponse.RefreshToken,
		"expires_in":    tokenResponse.ExpiresIn,
		"user":          user,
	})
}

// HandleGetMe handles retrieving the current user's information
func (h *AuthHandler) HandleGetMe(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the full user information from the database
	fullUser, err := h.userService.GetUserByID(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user information")
	}

	organization, err := h.userService.GetUserOrganization(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization information")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":         fullUser,
		"organization": organization,
	})
}

// HandleChangePassword handles changing a user's password
func (h *AuthHandler) HandleChangePassword(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	if err := h.authService.ChangePassword(c.Request().Context(), user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to change password: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// HandleRefreshToken handles refreshing an access token using a refresh token
func (h *AuthHandler) HandleRefreshToken(c echo.Context) error {
	var req models.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Refresh the access token
	tokenResponse, err := h.authService.RefreshAccessToken(
		c.Request().Context(),
		req.RefreshToken,
		h.jwtSecret,
		h.jwtExpiration,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token: "+err.Error())
	}

	return c.JSON(http.StatusOK, tokenResponse)
}

// HandleRevokeToken handles revoking a refresh token
func (h *AuthHandler) HandleRevokeToken(c echo.Context) error {
	var req models.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get current user ID from context
	_, err := utils.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Revoke the refresh token
	if err := h.authService.RevokeRefreshToken(c.Request().Context(), req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke token: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Token revoked successfully"})
}

// HandleRevokeAllTokens handles revoking all refresh tokens for a user
func (h *AuthHandler) HandleRevokeAllTokens(c echo.Context) error {
	// Get current user ID from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Revoke all refresh tokens for the user
	if err := h.authService.RevokeAllUserRefreshTokens(c.Request().Context(), user.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke tokens: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All tokens revoked successfully"})
}
