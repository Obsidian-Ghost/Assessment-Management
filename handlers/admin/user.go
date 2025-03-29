package admin

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/models"
	"assessment-management-system/services"
	"assessment-management-system/utils"
)

// UserHandler handles user-related routes for admin
type UserHandler struct {
	userService *services.UserService
	validator   *validator.Validate
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   utils.NewValidator(),
	}
}

// HandleCreateUser handles creating a new user
func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	user, err := h.userService.CreateUser(c.Request().Context(), admin.OrganizationID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// HandleGetAllUsers handles retrieving all users in an organization
func (h *UserHandler) HandleGetAllUsers(c echo.Context) error {
	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Extract filter parameters from query
	role := c.QueryParam("role")
	search := c.QueryParam("search")

	users, err := h.userService.GetUsersByOrganization(c.Request().Context(), admin.OrganizationID, role, search)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// HandleGetUserByID handles retrieving a user by ID
func (h *UserHandler) HandleGetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Ensure the user belongs to the admin's organization
	if user.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to user from another organization")
	}

	return c.JSON(http.StatusOK, user)
}

// HandleUpdateUser handles updating a user
func (h *UserHandler) HandleUpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the user belongs to the admin's organization
	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if user.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to user from another organization")
	}

	updatedUser, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// HandleDeleteUser handles deleting a user
func (h *UserHandler) HandleDeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the user belongs to the admin's organization
	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if user.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to user from another organization")
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// HandleBulkUploadUsers handles bulk uploading users
func (h *UserHandler) HandleBulkUploadUsers(c echo.Context) error {
	var req models.BulkUserUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	results, err := h.userService.BulkCreateUsers(c.Request().Context(), admin.OrganizationID, req.Users)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to bulk upload users: "+err.Error())
	}

	return c.JSON(http.StatusCreated, results)
}

// HandleGetTeacherStats handles retrieving teacher statistics
func (h *UserHandler) HandleGetTeacherStats(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Teacher ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	stats, err := h.userService.GetTeacherStats(c.Request().Context(), id, admin.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve teacher statistics: "+err.Error())
	}

	if stats == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Teacher not found")
	}

	return c.JSON(http.StatusOK, stats)
}

// HandleGetStudentStats handles retrieving student statistics
func (h *UserHandler) HandleGetStudentStats(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Student ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	stats, err := h.userService.GetStudentStats(c.Request().Context(), id, admin.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve student statistics: "+err.Error())
	}

	if stats == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Student not found")
	}

	return c.JSON(http.StatusOK, stats)
}
