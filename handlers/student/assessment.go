package student

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/models"
	"assessment-management-system/services"
	"assessment-management-system/utils"
)

// AssessmentHandler handles assessment-related routes for students
type AssessmentHandler struct {
	assessmentService *services.AssessmentService
	courseService     *services.CourseService
	validator         *validator.Validate
}

// NewAssessmentHandler creates a new AssessmentHandler
func NewAssessmentHandler(
	assessmentService *services.AssessmentService,
	courseService *services.CourseService,
) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
		courseService:     courseService,
		validator:         utils.NewValidator(),
	}
}

// HandleGetCourseAssessments handles retrieving all assessments for a course
func (h *AssessmentHandler) HandleGetCourseAssessments(c echo.Context) error {
	courseID := c.Param("courseId")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), courseID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	assessments, err := h.assessmentService.GetAssessmentsByCourse(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessments: "+err.Error())
	}

	// Enrich with submission status for this student
	assessmentStatuses := make([]models.StudentAssessmentStatus, 0, len(assessments))
	for _, assessment := range assessments {
		status, err := h.assessmentService.GetStudentAssessmentStatus(c.Request().Context(), assessment.ID, student.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment status: "+err.Error())
		}
		assessmentStatuses = append(assessmentStatuses, *status)
	}

	return c.JSON(http.StatusOK, assessmentStatuses)
}

// HandleGetAssessmentByID handles retrieving an assessment by ID
func (h *AssessmentHandler) HandleGetAssessmentByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
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

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), assessment.CourseID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	// Get the student's submission and grade for this assessment, if any
	status, err := h.assessmentService.GetStudentAssessmentStatus(c.Request().Context(), id, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment status: "+err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

// HandleSubmitAssessment handles submitting an assessment
func (h *AssessmentHandler) HandleSubmitAssessment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	var req models.CreateSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
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

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), assessment.CourseID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	// Check if the assessment is past due
	if assessment.DueDate != nil && assessment.DueDate.Before(time.Now()) {
		return echo.NewHTTPError(http.StatusForbidden, "Assessment is past due")
	}

	// Check if the student has already submitted
	hasSubmitted, err := h.assessmentService.HasStudentSubmitted(c.Request().Context(), id, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check submission status: "+err.Error())
	}

	if hasSubmitted {
		return echo.NewHTTPError(http.StatusForbidden, "You have already submitted this assessment")
	}

	submission, err := h.assessmentService.SubmitAssessment(c.Request().Context(), id, student.ID, req.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit assessment: "+err.Error())
	}

	return c.JSON(http.StatusCreated, submission)
}

// HandleViewSubmission handles viewing a student's own submission
func (h *AssessmentHandler) HandleViewSubmission(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
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

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), assessment.CourseID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	submission, err := h.assessmentService.GetStudentSubmission(c.Request().Context(), id, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submission: "+err.Error())
	}

	if submission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "You have not submitted this assessment")
	}

	// Get the grade if available
	grade, err := h.assessmentService.GetSubmissionGrade(c.Request().Context(), submission.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve grade: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"submission": submission,
		"grade":      grade,
	})
}

// HandleViewGrade handles viewing a student's grade for an assessment
func (h *AssessmentHandler) HandleViewGrade(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
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

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), assessment.CourseID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	submission, err := h.assessmentService.GetStudentSubmission(c.Request().Context(), id, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submission: "+err.Error())
	}

	if submission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "You have not submitted this assessment")
	}

	grade, err := h.assessmentService.GetSubmissionGrade(c.Request().Context(), submission.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve grade: "+err.Error())
	}

	if grade == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Your submission has not been graded yet")
	}

	return c.JSON(http.StatusOK, grade)
}
