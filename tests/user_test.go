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
	"assessment-management-system/middleware"
	"assessment-management-system/models"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, organizationID string, req models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, organizationID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUsersByOrganization(ctx context.Context, organizationID, role, search string) ([]*models.User, error) {
	args := m.Called(ctx, organizationID, role, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) BulkCreateUsers(ctx context.Context, organizationID string, users []models.CreateUserRequest) (map[string]interface{}, error) {
	args := m.Called(ctx, organizationID, users)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockUserService) GetUserOrganization(ctx context.Context, userID string) (*models.Organization, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockUserService) GetTeacherStats(ctx context.Context, teacherID, organizationID string) (*models.TeacherStats, error) {
	args := m.Called(ctx, teacherID, organizationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TeacherStats), args.Error(1)
}

func (m *MockUserService) GetStudentStats(ctx context.Context, studentID, organizationID string) (*models.StudentStats, error) {
	args := m.Called(ctx, studentID, organizationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentStats), args.Error(1)
}

// Helper function to setup the test
func setupUserTest(t *testing.T) (*echo.Echo, *MockUserService, *admin.UserHandler) {
	e := echo.New()
	mockService := new(MockUserService)
	handler := admin.NewUserHandler(mockService)
	return e, mockService, handler
}

// Mock the auth middleware to provide the admin user context
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

// Test creating a user
func TestHandleCreateUser(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Valid Create User", func(t *testing.T) {
		// Mock data
		req := models.CreateUserRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
			Role:      models.RoleTeacher,
		}

		user := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleTeacher,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("CreateUser", mock.Anything, "org-id", req).Return(user, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleCreateUser(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, response.ID)
			assert.Equal(t, user.Email, response.Email)
			assert.Equal(t, user.FirstName, response.FirstName)
			assert.Equal(t, user.LastName, response.LastName)
			assert.Equal(t, user.Role, response.Role)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Create User", func(t *testing.T) {
		// Invalid email
		req := models.CreateUserRequest{
			Email:     "invalid-email",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
			Role:      models.RoleTeacher,
		}

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should error due to validation
		err := handler.HandleCreateUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := models.CreateUserRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
			Role:      models.RoleTeacher,
		}

		// Service fails
		mockService.On("CreateUser", mock.Anything, "org-id", req).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleCreateUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting all users
func TestHandleGetAllUsers(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Get All Users Success", func(t *testing.T) {
		// Mock data
		users := []*models.User{
			{
				ID:             "user-id-1",
				OrganizationID: "org-id",
				Email:          "user1@example.com",
				FirstName:      "User",
				LastName:       "One",
				Role:           models.RoleTeacher,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             "user-id-2",
				OrganizationID: "org-id",
				Email:          "user2@example.com",
				FirstName:      "User",
				LastName:       "Two",
				Role:           models.RoleStudent,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		// Set up expectations - no filters
		mockService.On("GetUsersByOrganization", mock.Anything, "org-id", "", "").Return(users, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAllUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, users[0].ID, response[0].ID)
			assert.Equal(t, users[1].ID, response[1].ID)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Users with Filters", func(t *testing.T) {
		// Mock data - only teachers
		users := []*models.User{
			{
				ID:             "user-id-1",
				OrganizationID: "org-id",
				Email:          "teacher1@example.com",
				FirstName:      "Teacher",
				LastName:       "One",
				Role:           models.RoleTeacher,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		// Set up expectations - filter by role=teacher
		mockService.On("GetUsersByOrganization", mock.Anything, "org-id", "teacher", "").Return(users, nil).Once()

		// Create request with query parameters
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users?role=teacher", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.QueryParams().Set("role", "teacher")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAllUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 1)
			assert.Equal(t, users[0].ID, response[0].ID)
			assert.Equal(t, models.RoleTeacher, response[0].Role)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Users Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetUsersByOrganization", mock.Anything, "org-id", "", "").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetAllUsers(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting a user by ID
func TestHandleGetUserByID(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Get User by ID Success", func(t *testing.T) {
		// Mock data
		user := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleTeacher,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(user, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetUserByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, response.ID)
			assert.Equal(t, user.Email, response.Email)
			assert.Equal(t, user.FirstName, response.FirstName)
			assert.Equal(t, user.LastName, response.LastName)
			assert.Equal(t, user.Role, response.Role)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get User by ID Not Found", func(t *testing.T) {
		// Set up expectations - user not found
		mockService.On("GetUserByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/nonexistent-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetUserByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get User by ID From Different Organization", func(t *testing.T) {
		// Mock data - user from different organization
		user := &models.User{
			ID:             "user-id",
			OrganizationID: "other-org-id", // Different from admin's org
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleTeacher,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(user, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test - should fail with forbidden because user is from another org
		err := handler.HandleGetUserByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get User by ID Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetUserByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test updating a user
func TestHandleUpdateUser(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Update User Success", func(t *testing.T) {
		// Mock data
		email := "updated@example.com"
		firstName := "Updated"
		lastName := "User"
		role := models.RoleTeacher
		req := models.UpdateUserRequest{
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
			Role:      &role,
		}

		// Original user (to check org ID)
		originalUser := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "original@example.com",
			FirstName:      "Original",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Updated user
		updatedUser := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          email,
			FirstName:      firstName,
			LastName:       lastName,
			Role:           role,
			CreatedAt:      originalUser.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(originalUser, nil).Once()
		mockService.On("UpdateUser", mock.Anything, "user-id", req).Return(updatedUser, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleUpdateUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.User
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, updatedUser.ID, response.ID)
			assert.Equal(t, updatedUser.Email, response.Email)
			assert.Equal(t, updatedUser.FirstName, response.FirstName)
			assert.Equal(t, updatedUser.LastName, response.LastName)
			assert.Equal(t, updatedUser.Role, response.Role)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Update User Different Organization", func(t *testing.T) {
		// Mock data
		email := "updated@example.com"
		req := models.UpdateUserRequest{
			Email: &email,
		}

		// Original user (from different organization)
		originalUser := &models.User{
			ID:             "user-id",
			OrganizationID: "other-org-id", // Different from admin's org
			Email:          "original@example.com",
			FirstName:      "Original",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(originalUser, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleUpdateUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Update User Not Found", func(t *testing.T) {
		// Mock data
		email := "updated@example.com"
		req := models.UpdateUserRequest{
			Email: &email,
		}

		// Set up expectations - user not found
		mockService.On("GetUserByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/users/nonexistent-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleUpdateUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Update User Error", func(t *testing.T) {
		// Mock data
		email := "updated@example.com"
		req := models.UpdateUserRequest{
			Email: &email,
		}

		// Original user
		originalUser := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "original@example.com",
			FirstName:      "Original",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(originalUser, nil).Once()
		mockService.On("UpdateUser", mock.Anything, "user-id", req).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleUpdateUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test deleting a user
func TestHandleDeleteUser(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Delete User Success", func(t *testing.T) {
		// Original user
		user := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(user, nil).Once()
		mockService.On("DeleteUser", mock.Anything, "user-id").Return(nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleDeleteUser(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Delete User Different Organization", func(t *testing.T) {
		// Original user (from different organization)
		user := &models.User{
			ID:             "user-id",
			OrganizationID: "other-org-id", // Different from admin's org
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(user, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleDeleteUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Delete User Not Found", func(t *testing.T) {
		// Set up expectations - user not found
		mockService.On("GetUserByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/users/nonexistent-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleDeleteUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Delete User Error", func(t *testing.T) {
		// Original user
		user := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockService.On("GetUserByID", mock.Anything, "user-id").Return(user, nil).Once()
		mockService.On("DeleteUser", mock.Anything, "user-id").Return(errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/users/user-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleDeleteUser(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test bulk upload users
func TestHandleBulkUploadUsers(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Bulk Upload Users Success", func(t *testing.T) {
		// Mock data
		req := models.BulkUserUploadRequest{
			Users: []models.CreateUserRequest{
				{
					Email:     "user1@example.com",
					Password:  "password123",
					FirstName: "User",
					LastName:  "One",
					Role:      models.RoleTeacher,
				},
				{
					Email:     "user2@example.com",
					Password:  "password123",
					FirstName: "User",
					LastName:  "Two",
					Role:      models.RoleStudent,
				},
			},
		}

		results := map[string]interface{}{
			"successful": []models.User{
				{
					ID:             "user-id-1",
					OrganizationID: "org-id",
					Email:          "user1@example.com",
					FirstName:      "User",
					LastName:       "One",
					Role:           models.RoleTeacher,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				{
					ID:             "user-id-2",
					OrganizationID: "org-id",
					Email:          "user2@example.com",
					FirstName:      "User",
					LastName:       "Two",
					Role:           models.RoleStudent,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
			},
			"failed": []map[string]string{},
		}

		// Set up expectations
		mockService.On("BulkCreateUsers", mock.Anything, "org-id", req.Users).Return(results, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/users/bulk", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleBulkUploadUsers(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response, "successful")
			assert.Contains(t, response, "failed")
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Bulk Upload Users Error", func(t *testing.T) {
		// Mock data
		req := models.BulkUserUploadRequest{
			Users: []models.CreateUserRequest{
				{
					Email:     "user1@example.com",
					Password:  "password123",
					FirstName: "User",
					LastName:  "One",
					Role:      models.RoleTeacher,
				},
			},
		}

		// Service fails
		mockService.On("BulkCreateUsers", mock.Anything, "org-id", req.Users).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/users/bulk", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleBulkUploadUsers(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting teacher stats
func TestHandleGetTeacherStats(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Get Teacher Stats Success", func(t *testing.T) {
		// Mock data
		teacher := &models.User{
			ID:             "teacher-id",
			OrganizationID: "org-id",
			Email:          "teacher@example.com",
			FirstName:      "Teacher",
			LastName:       "User",
			Role:           models.RoleTeacher,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		stats := &models.TeacherStats{
			User:               teacher,
			AssignedCourses:    2,
			CreatedAssessments: 5,
			PendingGrading:     3,
		}

		// Set up expectations
		mockService.On("GetTeacherStats", mock.Anything, "teacher-id", "org-id").Return(stats, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/teachers/teacher-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("teacher-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetTeacherStats(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.TeacherStats
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, stats.AssignedCourses, response.AssignedCourses)
			assert.Equal(t, stats.CreatedAssessments, response.CreatedAssessments)
			assert.Equal(t, stats.PendingGrading, response.PendingGrading)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Teacher Stats Not Found", func(t *testing.T) {
		// Set up expectations - teacher not found
		mockService.On("GetTeacherStats", mock.Anything, "nonexistent-id", "org-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/teachers/nonexistent-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetTeacherStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Teacher Stats Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetTeacherStats", mock.Anything, "teacher-id", "org-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/teachers/teacher-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("teacher-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetTeacherStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting student stats
func TestHandleGetStudentStats(t *testing.T) {
	e, mockService, handler := setupUserTest(t)

	t.Run("Get Student Stats Success", func(t *testing.T) {
		// Mock data
		student := &models.User{
			ID:             "student-id",
			OrganizationID: "org-id",
			Email:          "student@example.com",
			FirstName:      "Student",
			LastName:       "User",
			Role:           models.RoleStudent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		stats := &models.StudentStats{
			User:                 student,
			EnrolledCourses:      2,
			CompletedAssessments: 5,
			PendingAssessments:   3,
			AverageGrade:         85.5,
		}

		// Set up expectations
		mockService.On("GetStudentStats", mock.Anything, "student-id", "org-id").Return(stats, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/students/student-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("student-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetStudentStats(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.StudentStats
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, stats.EnrolledCourses, response.EnrolledCourses)
			assert.Equal(t, stats.CompletedAssessments, response.CompletedAssessments)
			assert.Equal(t, stats.PendingAssessments, response.PendingAssessments)
			assert.Equal(t, stats.AverageGrade, response.AverageGrade)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Student Stats Not Found", func(t *testing.T) {
		// Set up expectations - student not found
		mockService.On("GetStudentStats", mock.Anything, "nonexistent-id", "org-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/students/nonexistent-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetStudentStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Student Stats Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetStudentStats", mock.Anything, "student-id", "org-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/users/students/student-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("student-id")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetStudentStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}
