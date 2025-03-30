package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"assessment-management-system/handlers/admin"
	"assessment-management-system/handlers/student"
	"assessment-management-system/handlers/teacher"
	"assessment-management-system/models"
)

// MockCourseService is a mock implementation of CourseService
type MockCourseService struct {
	mock.Mock
}

func (m *MockCourseService) CreateCourse(ctx context.Context, organizationID string, req models.CreateCourseRequest) (*models.Course, error) {
	args := m.Called(ctx, organizationID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseService) GetCoursesByOrganization(ctx context.Context, organizationID string) ([]*models.CourseWithStudentCount, error) {
	args := m.Called(ctx, organizationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CourseWithStudentCount), args.Error(1)
}

func (m *MockCourseService) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseService) UpdateCourse(ctx context.Context, id string, req models.UpdateCourseRequest) (*models.Course, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseService) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseService) AssignTeacherToCourse(ctx context.Context, courseID, teacherID string) error {
	args := m.Called(ctx, courseID, teacherID)
	return args.Error(0)
}

func (m *MockCourseService) RemoveTeacherFromCourse(ctx context.Context, courseID, teacherID string) error {
	args := m.Called(ctx, courseID, teacherID)
	return args.Error(0)
}

func (m *MockCourseService) GetCourseTeachers(ctx context.Context, courseID string) ([]*models.User, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockCourseService) ToggleCourseEnrollment(ctx context.Context, courseID string, enrollmentOpen bool) error {
	args := m.Called(ctx, courseID, enrollmentOpen)
	return args.Error(0)
}

func (m *MockCourseService) EnrollStudentInCourse(ctx context.Context, courseID, studentID string) error {
	args := m.Called(ctx, courseID, studentID)
	return args.Error(0)
}

func (m *MockCourseService) UnenrollStudentFromCourse(ctx context.Context, courseID, studentID string) error {
	args := m.Called(ctx, courseID, studentID)
	return args.Error(0)
}

func (m *MockCourseService) GetCourseStudents(ctx context.Context, courseID string) ([]*models.User, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockCourseService) BulkEnrollStudents(ctx context.Context, courseID string, studentIDs []string) (map[string]interface{}, error) {
	args := m.Called(ctx, courseID, studentIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCourseService) GetTeacherCourses(ctx context.Context, teacherID string) ([]*models.CourseWithStudentCount, error) {
	args := m.Called(ctx, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CourseWithStudentCount), args.Error(1)
}

func (m *MockCourseService) GetStudentCourses(ctx context.Context, studentID string) ([]*models.CourseWithDetails, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CourseWithDetails), args.Error(1)
}

func (m *MockCourseService) GetAvailableCourses(ctx context.Context, organizationID, studentID string) ([]*models.CourseWithDetails, error) {
	args := m.Called(ctx, organizationID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CourseWithDetails), args.Error(1)
}

func (m *MockCourseService) IsTeacherAssignedToCourse(ctx context.Context, courseID, teacherID string) (bool, error) {
	args := m.Called(ctx, courseID, teacherID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCourseService) IsStudentEnrolledInCourse(ctx context.Context, courseID, studentID string) (bool, error) {
	args := m.Called(ctx, courseID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCourseService) CourseHasTeachers(ctx context.Context, courseID string) (bool, error) {
	args := m.Called(ctx, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCourseService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

// Helper functions to setup the tests
func setupAdminCourseTest(t *testing.T) (*echo.Echo, *MockCourseService, *admin.CourseHandler) {
	e := echo.New()
	mockService := new(MockCourseService)
	handler := admin.NewCourseHandler(mockService)
	return e, mockService, handler
}

func setupTeacherCourseTest(t *testing.T) (*echo.Echo, *MockCourseService, *teacher.CourseHandler) {
	e := echo.New()
	mockService := new(MockCourseService)
	handler := teacher.NewCourseHandler(mockService)
	return e, mockService, handler
}

func setupStudentCourseTest(t *testing.T) (*echo.Echo, *MockCourseService, *student.CourseHandler) {
	e := echo.New()
	mockService := new(MockCourseService)
	handler := student.NewCourseHandler(mockService)
	return e, mockService, handler
}

// Mock user contexts
func mockAdminContext(c echo.Context) {
	admin := &models.User{
		ID:             "admin-id",
		OrganizationID: "org-id",
		Email:          "admin@example.com",
		FirstName:      "Admin",
		LastName:       "User",
		Role:           models.RoleAdmin,
	}
	c.Set("user", admin)
}

func mockTeacherContext(c echo.Context) {
	teacher := &models.User{
		ID:             "teacher-id",
		OrganizationID: "org-id",
		Email:          "teacher@example.com",
		FirstName:      "Teacher",
		LastName:       "User",
		Role:           models.RoleTeacher,
	}
	c.Set("user", teacher)
}

func mockStudentContext(c echo.Context) {
	student := &models.User{
		ID:             "student-id",
		OrganizationID: "org-id",
		Email:          "student@example.com",
		FirstName:      "Student",
		LastName:       "User",
		Role:           models.RoleStudent,
	}
	c.Set("user", student)
}

// Test admin course handler - creating a course
func TestHandleCreateCourse(t *testing.T) {
	e, mockService, handler := setupAdminCourseTest(t)

	t.Run("Valid Create Course", func(t *testing.T) {
		// Mock data
		req := models.CreateCourseRequest{
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false,
			OrganizationID: "org-id", // Add the organization ID
		}

		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations - the service now uses the organization ID from the request
		mockService.On("CreateCourse", mock.Anything, req.OrganizationID, req).Return(course, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleCreateCourse(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.Course
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, course.ID, response.ID)
			assert.Equal(t, course.Name, response.Name)
			assert.Equal(t, course.Description, response.Description)
			assert.Equal(t, course.EnrollmentOpen, response.EnrollmentOpen)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Create Course - No Name", func(t *testing.T) {
		// Invalid request (no name)
		req := models.CreateCourseRequest{
			Name:           "", // Name is required
			Description:    "Course description",
			EnrollmentOpen: false,
			OrganizationID: "org-id",
		}

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should error due to validation
		err := handler.HandleCreateCourse(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("Invalid Create Course - No Organization ID", func(t *testing.T) {
		// Invalid request (no organization ID)
		req := models.CreateCourseRequest{
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false,
			OrganizationID: "", // Organization ID is required
		}

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should error due to validation
		err := handler.HandleCreateCourse(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("Organization Mismatch", func(t *testing.T) {
		// Request with organization ID different from admin's org
		req := models.CreateCourseRequest{
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false,
			OrganizationID: "different-org-id", // Different from admin's org
		}

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should error due to org mismatch
		err := handler.HandleCreateCourse(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := models.CreateCourseRequest{
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false,
			OrganizationID: "org-id",
		}

		// Service fails
		mockService.On("CreateCourse", mock.Anything, req.OrganizationID, req).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleCreateCourse(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting all courses
func TestHandleGetAllCourses(t *testing.T) {
	e, mockService, handler := setupAdminCourseTest(t)

	t.Run("Get All Courses Success", func(t *testing.T) {
		// Mock data
		courses := []*models.CourseWithStudentCount{
			{
				Course: &models.Course{
					ID:             "course-id-1",
					OrganizationID: "org-id",
					Name:           "Course 1",
					Description:    "Description 1",
					EnrollmentOpen: true,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				StudentCount: 5,
				TeacherCount: 1,
			},
			{
				Course: &models.Course{
					ID:             "course-id-2",
					OrganizationID: "org-id",
					Name:           "Course 2",
					Description:    "Description 2",
					EnrollmentOpen: false,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				StudentCount: 0,
				TeacherCount: 2,
			},
		}

		// Set up expectations
		mockService.On("GetCoursesByOrganization", mock.Anything, "org-id").Return(courses, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAllCourses(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.CourseWithStudentCount
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, courses[0].Course.ID, response[0].Course.ID)
			assert.Equal(t, courses[1].Course.ID, response[1].Course.ID)
			assert.Equal(t, courses[0].StudentCount, response[0].StudentCount)
			assert.Equal(t, courses[1].StudentCount, response[1].StudentCount)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get All Courses Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetCoursesByOrganization", mock.Anything, "org-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetAllCourses(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting a course by ID
func TestHandleGetCourseByID(t *testing.T) {
	e, mockService, handler := setupAdminCourseTest(t)

	t.Run("Get Course by ID Success", func(t *testing.T) {
		// Mock data
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses/course-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetCourseByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Course
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, course.ID, response.ID)
			assert.Equal(t, course.Name, response.Name)
			assert.Equal(t, course.Description, response.Description)
			assert.Equal(t, course.EnrollmentOpen, response.EnrollmentOpen)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Course by ID Not Found", func(t *testing.T) {
		// Set up expectations - course not found
		mockService.On("GetCourseByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses/nonexistent-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetCourseByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Course by ID From Different Organization", func(t *testing.T) {
		// Mock data - course from different organization
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "other-org-id", // Different from admin's org
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses/course-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should fail with forbidden because course is from another org
		err := handler.HandleGetCourseByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Course by ID Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/courses/course-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetCourseByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test assigning a teacher to a course
func TestHandleAssignTeacher(t *testing.T) {
	e, mockService, handler := setupAdminCourseTest(t)

	t.Run("Assign Teacher Success", func(t *testing.T) {
		// Mock data
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		req := models.AssignTeacherRequest{
			TeacherID: "teacher-id",
		}

		// Set up expectations
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()
		mockService.On("AssignTeacherToCourse", mock.Anything, "course-id", "teacher-id").Return(nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses/course-id/teachers", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleAssignTeacher(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response, "message")
			assert.Equal(t, "Teacher assigned successfully", response["message"])
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Assign Teacher Course Not Found", func(t *testing.T) {
		// Set up expectations - course not found
		req := models.AssignTeacherRequest{
			TeacherID: "teacher-id",
		}

		mockService.On("GetCourseByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses/nonexistent-id/teachers", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleAssignTeacher(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Assign Teacher Error", func(t *testing.T) {
		// Mock data
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		req := models.AssignTeacherRequest{
			TeacherID: "teacher-id",
		}

		// Set up expectations - service fails
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()
		mockService.On("AssignTeacherToCourse", mock.Anything, "course-id", "teacher-id").Return(errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/courses/course-id/teachers", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleAssignTeacher(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test teacher-specific course handlers
func TestTeacherCourseHandlers(t *testing.T) {
	e, mockService, handler := setupTeacherCourseTest(t)

	t.Run("Get Assigned Courses", func(t *testing.T) {
		// Mock data
		courses := []*models.CourseWithStudentCount{
			{
				Course: &models.Course{
					ID:             "course-id-1",
					OrganizationID: "org-id",
					Name:           "Course 1",
					Description:    "Description 1",
					EnrollmentOpen: true,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				StudentCount: 5,
				TeacherCount: 2,
			},
		}

		// Set up expectations
		mockService.On("GetTeacherCourses", mock.Anything, "teacher-id").Return(courses, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/teacher/courses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAssignedCourses(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.CourseWithStudentCount
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 1)
			assert.Equal(t, courses[0].Course.ID, response[0].Course.ID)
			assert.Equal(t, courses[0].StudentCount, response[0].StudentCount)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Course by ID for Teacher", func(t *testing.T) {
		// Mock data
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("IsTeacherAssignedToCourse", mock.Anything, "course-id", "teacher-id").Return(true, nil).Once()
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/teacher/courses/course-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetCourseByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Course
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, course.ID, response.ID)
			assert.Equal(t, course.Name, response.Name)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Course Students for Teacher", func(t *testing.T) {
		// Mock data
		students := []*models.User{
			{
				ID:             "student-id-1",
				OrganizationID: "org-id",
				Email:          "student1@example.com",
				FirstName:      "Student",
				LastName:       "One",
				Role:           models.RoleStudent,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             "student-id-2",
				OrganizationID: "org-id",
				Email:          "student2@example.com",
				FirstName:      "Student",
				LastName:       "Two",
				Role:           models.RoleStudent,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		// Set up expectations
		mockService.On("IsTeacherAssignedToCourse", mock.Anything, "course-id", "teacher-id").Return(true, nil).Once()
		mockService.On("GetCourseStudents", mock.Anything, "course-id").Return(students, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/teacher/courses/course-id/students", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetCourseStudents(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, students[0].ID, response[0].ID)
			assert.Equal(t, students[1].ID, response[1].ID)
		}

		mockService.AssertExpectations(t)
	})
}

// Test student-specific course handlers
func TestStudentCourseHandlers(t *testing.T) {
	e, mockService, handler := setupStudentCourseTest(t)

	t.Run("Get Enrolled Courses", func(t *testing.T) {
		// Mock data
		courses := []*models.CourseWithDetails{
			{
				Course: &models.Course{
					ID:             "course-id-1",
					OrganizationID: "org-id",
					Name:           "Course 1",
					Description:    "Description 1",
					EnrollmentOpen: true,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				OrganizationID:   "org-id",
				OrganizationName: "Test Organization",
				Teachers: []*models.User{
					{
						ID:        "teacher-id",
						FirstName: "Teacher",
						LastName:  "User",
						Role:      models.RoleTeacher,
					},
				},
				StudentCount: 5,
			},
		}

		// Set up expectations
		mockService.On("GetStudentCourses", mock.Anything, "student-id").Return(courses, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/courses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetEnrolledCourses(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.CourseWithDetails
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 1)
			assert.Equal(t, courses[0].Course.ID, response[0].Course.ID)
			assert.Equal(t, courses[0].OrganizationName, response[0].OrganizationName)
			assert.Len(t, response[0].Teachers, 1)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Available Courses", func(t *testing.T) {
		// Mock data
		courses := []*models.CourseWithDetails{
			{
				Course: &models.Course{
					ID:             "course-id-2",
					OrganizationID: "org-id",
					Name:           "Course 2",
					Description:    "Description 2",
					EnrollmentOpen: true,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				OrganizationID:   "org-id",
				OrganizationName: "Test Organization",
				Teachers: []*models.User{
					{
						ID:        "teacher-id",
						FirstName: "Teacher",
						LastName:  "User",
						Role:      models.RoleTeacher,
					},
				},
				StudentCount: 3,
			},
		}

		// Set up expectations
		mockService.On("GetAvailableCourses", mock.Anything, "org-id", "student-id").Return(courses, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/courses/available", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAvailableCourses(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.CourseWithDetails
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 1)
			assert.Equal(t, courses[0].Course.ID, response[0].Course.ID)
			assert.Equal(t, courses[0].OrganizationName, response[0].OrganizationName)
			assert.Len(t, response[0].Teachers, 1)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Enroll in Course", func(t *testing.T) {
		// Mock data
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()
		mockService.On("CourseHasTeachers", mock.Anything, "course-id").Return(true, nil).Once()
		mockService.On("EnrollStudentInCourse", mock.Anything, "course-id", "student-id").Return(nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodPost, "/api/student/courses/course-id/enroll", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleEnrollInCourse(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response, "message")
			assert.Equal(t, "Successfully enrolled in course", response["message"])
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Enrollment Forbidden - Course Not Open", func(t *testing.T) {
		// Mock data - course with enrollment closed
		course := &models.Course{
			ID:             "course-id",
			OrganizationID: "org-id",
			Name:           "Test Course",
			Description:    "Course description",
			EnrollmentOpen: false, // Enrollment closed
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetCourseByID", mock.Anything, "course-id").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodPost, "/api/student/courses/course-id/enroll", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("course-id")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test - should fail because enrollment is closed
		err := handler.HandleEnrollInCourse(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockService.AssertExpectations(t)
	})
}
