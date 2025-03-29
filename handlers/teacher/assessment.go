package teacher

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/models"
	"assessment-management-system/services"
	"assessment-management-system/utils"
)

// AssessmentHandler handles assessment-related routes for teachers
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

// HandleCreateAssessment handles creating a new assessment
func (h *AssessmentHandler) HandleCreateAssessment(c echo.Context) error {
	var req models.CreateAssessmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the teacher is assigned to the course
	isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), req.CourseID, teacher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
	}

	if !isAssigned {
		return echo.NewHTTPError(http.StatusForbidden, "You are not assigned to this course")
	}

	assessment, err := h.assessmentService.CreateAssessment(c.Request().Context(), teacher.ID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create assessment: "+err.Error())
	}

	return c.JSON(http.StatusCreated, assessment)
}

// HandleGetAssessments handles retrieving all assessments for a course
func (h *AssessmentHandler) HandleGetAssessments(c echo.Context) error {
	courseID := c.Param("courseId")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the teacher is assigned to the course
	isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), courseID, teacher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
	}

	if !isAssigned {
		return echo.NewHTTPError(http.StatusForbidden, "You are not assigned to this course")
	}

	assessments, err := h.assessmentService.GetAssessmentsByCourse(c.Request().Context(), courseID)
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

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
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

	// Check if the teacher created the assessment or is assigned to the course
	if assessment.TeacherID != teacher.ID {
		isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), assessment.CourseID, teacher.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
		}

		if !isAssigned {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view this assessment")
		}
	}

	return c.JSON(http.StatusOK, assessment)
}

// HandleUpdateAssessment handles updating an assessment
func (h *AssessmentHandler) HandleUpdateAssessment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	var req models.UpdateAssessmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the assessment to verify ownership
	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	if assessment == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assessment not found")
	}

	// Only the teacher who created the assessment can update it
	if assessment.TeacherID != teacher.ID {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update this assessment")
	}

	updatedAssessment, err := h.assessmentService.UpdateAssessment(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update assessment: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedAssessment)
}

// HandleDeleteAssessment handles deleting an assessment
func (h *AssessmentHandler) HandleDeleteAssessment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the assessment to verify ownership
	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	if assessment == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assessment not found")
	}

	// Only the teacher who created the assessment can delete it
	if assessment.TeacherID != teacher.ID {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to delete this assessment")
	}

	if err := h.assessmentService.DeleteAssessment(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete assessment: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// HandleGetSubmissions handles retrieving all submissions for an assessment
func (h *AssessmentHandler) HandleGetSubmissions(c echo.Context) error {
	assessmentID := c.Param("id")
	if assessmentID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Assessment ID is required")
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the assessment to verify permission
	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), assessmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	if assessment == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Assessment not found")
	}

	// Check if the teacher is assigned to the course or created the assessment
	if assessment.TeacherID != teacher.ID {
		isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), assessment.CourseID, teacher.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
		}

		if !isAssigned {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view submissions for this assessment")
		}
	}

	submissions, err := h.assessmentService.GetAssessmentSubmissions(c.Request().Context(), assessmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submissions: "+err.Error())
	}

	return c.JSON(http.StatusOK, submissions)
}

// HandleGradeSubmission handles grading a submission
func (h *AssessmentHandler) HandleGradeSubmission(c echo.Context) error {
	submissionID := c.Param("submissionId")
	if submissionID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Submission ID is required")
	}

	var req models.GradeSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(h.validator, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get the submission to verify permission
	submission, err := h.assessmentService.GetSubmissionByID(c.Request().Context(), submissionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve submission: "+err.Error())
	}

	if submission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
	}

	// Get the assessment to verify permission
	assessment, err := h.assessmentService.GetAssessmentByID(c.Request().Context(), submission.AssessmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assessment: "+err.Error())
	}

	// Check if the teacher created the assessment or is assigned to the course
	if assessment.TeacherID != teacher.ID {
		isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), assessment.CourseID, teacher.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
		}

		if !isAssigned {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to grade this submission")
		}
	}

	// Validate that the score doesn't exceed the maximum score
	if req.Score > float64(assessment.MaxScore) {
		return echo.NewHTTPError(http.StatusBadRequest, "Score cannot exceed the maximum score for this assessment")
	}

	grade, err := h.assessmentService.GradeSubmission(c.Request().Context(), submissionID, teacher.ID, req.Score, req.Feedback)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to grade submission: "+err.Error())
	}

	return c.JSON(http.StatusOK, grade)
}
