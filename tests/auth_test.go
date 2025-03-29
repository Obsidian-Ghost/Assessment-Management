package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"assessment-management-system/handlers"
	"assessment-management-system/middleware"
	"assessment-management-system/models"
	"assessment-management-system/utils"
)

// Mock services
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, userID string, jwtSecret string, expiration time.Duration) (string, error) {
	args := m.Called(ctx, userID, jwtSecret, expiration)
	return args.String(0), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserOrganization(ctx context.Context, userID string) (*models.Organization, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

// Test login handler
func TestHandleLogin(t *testing.T) {
	// Setup
	e := echo.New()
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	jwtSecret := "test-secret"
	jwtExpiration := time.Hour * 24

	authHandler := handlers.NewAuthHandler(mockAuthService, mockUserService, jwtSecret, jwtExpiration)

	// Test valid login
	t.Run("Valid Login", func(t *testing.T) {
		// Mock data
		validUser := &models.User{
			ID:             "user-id",
			OrganizationID: "org-id",
			Email:          "test@example.com",
			FirstName:      "Test",
			LastName:       "User",
			Role:           models.RoleTeacher,
		}

		// Setup expectations
		mockAuthService.On("Authenticate", mock.Anything, "test@example.com", "password123").Return(validUser, nil).Once()

		// Request body
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(loginReq)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assert
		if assert.NoError(t, authHandler.HandleLogin(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check for token and user in response
			assert.Contains(t, response, "token")
			assert.Contains(t, response, "user")
		}

		// Verify mocks
		mockAuthService.AssertExpectations(t)
	})

	// Test invalid credentials
	t.Run("Invalid Credentials", func(t *testing.T) {
		// Setup expectations
		mockAuthService.On("Authenticate", mock.Anything, "wrong@example.com", "wrongpass").Return(nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")).Once()

		// Request body
		loginReq := models.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpass",
		}
		jsonBody, _ := json.Marshal(loginReq)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assert
		err := authHandler.HandleLogin(c)
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		}

		// Verify mocks
		mockAuthService.AssertExpectations(t)
	})
}

// Test JWT validation
func TestJWTValidation(t *testing.T) {
	// Create a user for token generation
	user := &models.User{
		ID:             "user-id",
		OrganizationID: "org-id",
		Email:          "test@example.com",
		FirstName:      "Test",
		LastName:       "User",
		Role:           models.RoleTeacher,
	}

	jwtSecret := "test-secret"
	jwtExpiration := time.Hour * 24

	// Generate a valid token
	token, err := utils.GenerateToken(user, jwtSecret, jwtExpiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := utils.ValidateToken(token, jwtSecret)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)

	// Test with incorrect secret
	_, err = utils.ValidateToken(token, "wrong-secret")
	assert.Error(t, err)

	// Test middleware
	e := echo.New()
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// Define a simple handler
	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Valid token test
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := authMiddleware(testHandler)
	assert.NoError(t, h(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// Invalid token test
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	h = authMiddleware(testHandler)
	err = h(c)
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
