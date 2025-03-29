package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/services"
)

// AssessmentHandler handles assessment-related routes for admin (read-only)
type AssessmentHandler struct {
	assessmentService *services.AssessmentService
	courseService     *services.CourseService
}

// NewAssessmentHandler creates a new AssessmentHandler
func NewAssessmentHandler(
	assessmentService *services.AssessmentService,
	courseService *services.CourseService,
) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
		courseService:     courseService,
	}
}

// HandleGetAllAssessments handles retrieving all assessments for the admin's organization
func (h *AssessmentHandler) HandleGetAllAssessments(c echo.Context) error {
	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.QueryParam("courseId")
	assessments, err := h.assessmentService.GetAssessmentsByOrganization(c.Request().Context(), admin.OrganizationID, courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessments: "+err.Error())
	}

	return c.JSON(http.StatusOK, assessments)
}

// HandleGetAssessmentByID handles retrieving an assessment by ID
func (h *AssessmentHandler) HandleGetAssessmentByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	if assessment == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assessment not found")
	}

	// Verify that the assessment belongs to a course in the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), assessment.CourseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to assessment from another organization")
	}

	return c.JSON(http.StatusOK, assessment)
}

// HandleGetAssessmentSubmissions handles retrieving submissions for an assessment
func (h *AssessmentHandler) HandleGetAssessmentSubmissions(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	if assessment == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assessment not found")
	}

	// Verify that the assessment belongs to a course in the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), assessment.CourseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to assessment from another organization")
	}

	submissions, err := h.assessmentService.GetAssessmentSubmissions(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submissions: "+err.Error())
	}

	return c.JSON(http.StatusOK, submissions)
}

// HandleGetSubmissionGrades handles retrieving grades for a submission
func (h *AssessmentHandler) HandleGetSubmissionGrades(c echo.Context) error {
	submissionID := c.Param("submissionId")
	if submissionID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Submission ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the submission to check if it belongs to the admin's organization
	submission, err := h.assessmentService.GetSubmissionByID(c.Request().Context(), submissionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submission: "+err.Error())
	}

	if submission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
	}

	// Get the assessment to check organization
	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), submission.AssessmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	// Get the course to check organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), assessment.CourseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to submission from another organization")
	}

	grade, err := h.assessmentService.GetSubmissionGrade(c.Request().Context(), submissionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve grade: "+err.Error())
	}

	if grade == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Grade not found")
	}

	return c.JSON(http.StatusOK, grade)
}
