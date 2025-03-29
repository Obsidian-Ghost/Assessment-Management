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

// CourseHandler handles course-related routes for admin
type CourseHandler struct {
	courseService *services.CourseService
	validator     *validator.Validate
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		validator:     utils.NewValidator(),
	}
}

// HandleCreateCourse handles creating a new course
func (h *CourseHandler) HandleCreateCourse(c echo.Context) error {
	var req models.CreateCourseRequest
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

	course, err := h.courseService.CreateCourse(c.Request().Context(), admin.OrganizationID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create course: "+err.Error())
	}

	return c.JSON(http.StatusCreated, course)
}

// HandleGetAllCourses handles retrieving all courses in an organization
func (h *CourseHandler) HandleGetAllCourses(c echo.Context) error {
	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courses, err := h.courseService.GetCoursesByOrganization(c.Request().Context(), admin.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve courses: "+err.Error())
	}

	return c.JSON(http.StatusOK, courses)
}

// HandleGetCourseByID handles retrieving a course by ID
func (h *CourseHandler) HandleGetCourseByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	course, err := h.courseService.GetCourseByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	// Ensure the course belongs to the admin's organization
	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	return c.JSON(http.StatusOK, course)
}

// HandleUpdateCourse handles updating a course
func (h *CourseHandler) HandleUpdateCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	var req models.UpdateCourseRequest
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

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	updatedCourse, err := h.courseService.UpdateCourse(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update course: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedCourse)
}

// HandleDeleteCourse handles deleting a course
func (h *CourseHandler) HandleDeleteCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	if err := h.courseService.DeleteCourse(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete course: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// HandleAssignTeacher handles assigning a teacher to a course
func (h *CourseHandler) HandleAssignTeacher(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	var req models.AssignTeacherRequest
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

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	if err := h.courseService.AssignTeacherToCourse(c.Request().Context(), courseID, req.TeacherID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign teacher to course: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Teacher assigned successfully"})
}

// HandleRemoveTeacher handles removing a teacher from a course
func (h *CourseHandler) HandleRemoveTeacher(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	teacherID := c.Param("teacherId")
	if teacherID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Teacher ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	if err := h.courseService.RemoveTeacherFromCourse(c.Request().Context(), courseID, teacherID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove teacher from course: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Teacher removed successfully"})
}

// HandleGetCourseTeachers handles retrieving all teachers assigned to a course
func (h *CourseHandler) HandleGetCourseTeachers(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	teachers, err := h.courseService.GetCourseTeachers(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course teachers: "+err.Error())
	}

	return c.JSON(http.StatusOK, teachers)
}

// HandleToggleEnrollment handles toggling a course's enrollment status
func (h *CourseHandler) HandleToggleEnrollment(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	var req struct {
		EnrollmentOpen bool `json:"enrollment_open"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	if err := h.courseService.ToggleCourseEnrollment(c.Request().Context(), courseID, req.EnrollmentOpen); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to toggle course enrollment: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Course enrollment status updated successfully"})
}

// HandleManageStudentEnrollment handles adding or removing students from a course
func (h *CourseHandler) HandleManageStudentEnrollment(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	var req models.EnrollStudentRequest
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

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	action := c.QueryParam("action")
	if action == "enroll" {
		if err := h.courseService.EnrollStudentInCourse(c.Request().Context(), courseID, req.StudentID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enroll student in course: "+err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Student enrolled successfully"})
	} else if action == "unenroll" {
		if err := h.courseService.UnenrollStudentFromCourse(c.Request().Context(), courseID, req.StudentID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unenroll student from course: "+err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Student unenrolled successfully"})
	}

	return echo.NewHTTPError(http.StatusBadRequest, "Invalid action, must be 'enroll' or 'unenroll'")
}

// HandleBulkEnrollStudents handles enrolling multiple students in a course
func (h *CourseHandler) HandleBulkEnrollStudents(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	var req models.BulkEnrollmentRequest
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

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	results, err := h.courseService.BulkEnrollStudents(c.Request().Context(), courseID, req.StudentIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to bulk enroll students: "+err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// HandleGetCourseStudents handles retrieving all students enrolled in a course
func (h *CourseHandler) HandleGetCourseStudents(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the admin's organization ID from the token
	admin, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Verify that the course belongs to the admin's organization
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if course.OrganizationID != admin.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied to course from another organization")
	}

	students, err := h.courseService.GetCourseStudents(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course students: "+err.Error())
	}

	return c.JSON(http.StatusOK, students)
}
