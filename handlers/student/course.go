package student

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"assessment-management-system/middleware"
	"assessment-management-system/services"
)

// CourseHandler handles course-related routes for students
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// HandleGetEnrolledCourses handles retrieving all courses a student is enrolled in
func (h *CourseHandler) HandleGetEnrolledCourses(c echo.Context) error {
	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courses, err := h.courseService.GetStudentCourses(c.Request().Context(), student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve enrolled courses: "+err.Error())
	}

	return c.JSON(http.StatusOK, courses)
}

// HandleGetCourseByID handles retrieving a course by ID (student must be enrolled)
func (h *CourseHandler) HandleGetCourseByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the student is enrolled in the course
	isEnrolled, err := h.courseService.IsStudentEnrolledInCourse(c.Request().Context(), id, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check enrollment: "+err.Error())
	}

	if !isEnrolled {
		return echo.NewHTTPError(http.StatusForbidden, "You are not enrolled in this course")
	}

	course, err := h.courseService.GetCourseByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	// Get course teachers
	teachers, err := h.courseService.GetCourseTeachers(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course teachers: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"course":   course,
		"teachers": teachers,
	})
}

// HandleGetAvailableCourses handles retrieving available courses for enrollment
func (h *CourseHandler) HandleGetAvailableCourses(c echo.Context) error {
	// Get the student's organization ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courses, err := h.courseService.GetAvailableCourses(c.Request().Context(), student.OrganizationID, student.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve available courses: "+err.Error())
	}

	return c.JSON(http.StatusOK, courses)
}

// HandleEnrollInCourse handles enrolling a student in a course
func (h *CourseHandler) HandleEnrollInCourse(c echo.Context) error {
	courseID := c.Param("id")
	if courseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Course ID is required")
	}

	// Get the student's ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Check if the course exists and is open for enrollment
	course, err := h.courseService.GetCourseByID(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve course: "+err.Error())
	}

	if course == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Course not found")
	}

	if !course.EnrollmentOpen {
		return echo.NewHTTPError(http.StatusForbidden, "Course is not open for enrollment")
	}

	// Check if the course belongs to the student's organization
	if course.OrganizationID != student.OrganizationID {
		return echo.NewHTTPError(http.StatusForbidden, "You cannot enroll in courses from other organizations")
	}

	// Check if the course has at least one teacher assigned
	hasTeachers, err := h.courseService.CourseHasTeachers(c.Request().Context(), courseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if course has teachers: "+err.Error())
	}

	if !hasTeachers {
		return echo.NewHTTPError(http.StatusForbidden, "You cannot enroll in a course without teachers")
	}

	// Enroll the student
	if err := h.courseService.EnrollStudentInCourse(c.Request().Context(), courseID, student.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enroll in course: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully enrolled in course"})
}

// HandleGetOrganizationDetails handles retrieving organization details
func (h *CourseHandler) HandleGetOrganizationDetails(c echo.Context) error {
	// Get the student's organization ID from the token
	student, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	organization, err := h.courseService.GetOrganizationByID(c.Request().Context(), student.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve organization details: "+err.Error())
	}

	if organization == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, organization)
}
