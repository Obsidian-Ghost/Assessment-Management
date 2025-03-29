package admin

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"assessment-management-system/models"
	"assessment-management-system/services"
	"assessment-management-system/utils"
)

// OrganizationHandler handles organization-related routes for admin
type OrganizationHandler struct {
	organizationService *services.OrganizationService
	validator           *validator.Validate
}

// NewOrganizationHandler creates a new OrganizationHandler
func NewOrganizationHandler(organizationService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
		validator:           utils.NewValidator(),
	}
}

// HandleCreateOrganization handles creating a new organization
func (h *OrganizationHandler) HandleCreateOrganization(c echo.Context) error {
	var req models.CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	organization, err := h.organizationService.CreateOrganization(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create organization: "+err.Error())
	}

	return c.JSON(http.StatusCreated, organization)
}

// HandleGetAllOrganizations handles retrieving all organizations
func (h *OrganizationHandler) HandleGetAllOrganizations(c echo.Context) error {
	organizations, err := h.organizationService.GetAllOrganizations(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organizations: "+err.Error())
	}

	return c.JSON(http.StatusOK, organizations)
}

// HandleGetOrganizationByID handles retrieving an organization by ID
func (h *OrganizationHandler) HandleGetOrganizationByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization ID is required")
	}

	organization, err := h.organizationService.GetOrganizationByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization: "+err.Error())
	}

	if organization == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, organization)
}

// HandleUpdateOrganization handles updating an organization
func (h *OrganizationHandler) HandleUpdateOrganization(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization ID is required")
	}

	var req models.UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	organization, err := h.organizationService.UpdateOrganization(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update organization: "+err.Error())
	}

	if organization == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, organization)
}

// HandleDeleteOrganization handles deleting an organization
func (h *OrganizationHandler) HandleDeleteOrganization(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization ID is required")
	}

	if err := h.organizationService.DeleteOrganization(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete organization: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// HandleGetOrganizationStats handles retrieving organization statistics
func (h *OrganizationHandler) HandleGetOrganizationStats(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization ID is required")
	}

	stats, err := h.organizationService.GetOrganizationStats(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization statistics: "+err.Error())
	}

	if stats == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, stats)
}
