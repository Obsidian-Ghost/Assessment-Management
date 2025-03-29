package teacher

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/services"
)

// CourseHandler handles course-related routes for teachers
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// HandleGetAssignedCourses handles retrieving all courses assigned to a teacher
func (h *CourseHandler) HandleGetAssignedCourses(c echo.Context) error {
	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courses, err := h.courseService.GetTeacherCourses(c.Request().Context(), teacher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve assigned courses: "+err.Error())
	}

	return c.JSON(http.StatusOK, courses)
}

// HandleGetCourseByID handles retrieving a course by ID (teacher must be assigned to the course)
func (h *CourseHandler) HandleGetCourseByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the teacher's ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the teacher is assigned to the course
	isAssigned, err := h.courseService.IsTeacherAssignedToCourse(c.Request().Context(), id, teacher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check course assignment: "+err.Error())
	}

	if !isAssigned {
		return echo.NewHTTPError(http.StatusForbidden, "You are not assigned to this course")
	}

	course, err := h.courseService.GetCourseByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	return c.JSON(http.StatusOK, course)
}

// HandleGetCourseStudents handles retrieving all students enrolled in a course
func (h *CourseHandler) HandleGetCourseStudents(c echo.Context) error {
	courseID := c.Param("id")
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

	students, err := h.courseService.GetCourseStudents(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course students: "+err.Error())
	}

	return c.JSON(http.StatusOK, students)
}

// HandleGetOrganizationDetails handles retrieving organization details
func (h *CourseHandler) HandleGetOrganizationDetails(c echo.Context) error {
	// Get the teacher's organization ID from the token
	teacher, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	organization, err := h.courseService.GetOrganizationByID(c.Request().Context(), teacher.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization details: "+err.Error())
	}

	if organization == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, organization)
}
