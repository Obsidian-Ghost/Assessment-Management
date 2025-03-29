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

// MockOrganizationService is a mock implementation of OrganizationService
type MockOrganizationService struct {
	mock.Mock
}

func (m *MockOrganizationService) CreateOrganization(ctx context.Context, req models.CreateOrganizationRequest) (*models.Organization, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) GetAllOrganizations(ctx context.Context) ([]*models.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) UpdateOrganization(ctx context.Context, id string, req models.UpdateOrganizationRequest) (*models.Organization, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationService) GetOrganizationStats(ctx context.Context, id string) (*models.OrganizationStats, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OrganizationStats), args.Error(1)
}

// Helper function to setup the test
func setupOrganizationTest(t *testing.T) (*echo.Echo, *MockOrganizationService, *admin.OrganizationHandler) {
	e := echo.New()
	mockService := new(MockOrganizationService)
	handler := admin.NewOrganizationHandler(mockService)
	return e, mockService, handler
}

// Test creating an organization
func TestHandleCreateOrganization(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Valid Create Organization", func(t *testing.T) {
		// Mock data
		req := models.CreateOrganizationRequest{
			Name:   "Test Organization",
			Slogan: "Testing is Important",
		}

		org := &models.Organization{
			ID:        "org-id",
			Name:      "Test Organization",
			Slogan:    "Testing is Important",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set up expectations
		mockService.On("CreateOrganization", mock.Anything, req).Return(org, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/organizations", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Perform the test
		if assert.NoError(t, handler.HandleCreateOrganization(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.Organization
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, org.ID, response.ID)
			assert.Equal(t, org.Name, response.Name)
			assert.Equal(t, org.Slogan, response.Slogan)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Create Organization", func(t *testing.T) {
		// Empty name (invalid)
		req := models.CreateOrganizationRequest{
			Name:   "",
			Slogan: "Testing is Important",
		}

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/organizations", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Perform the test - should error due to validation
		err := handler.HandleCreateOrganization(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := models.CreateOrganizationRequest{
			Name:   "Test Organization",
			Slogan: "Testing is Important",
		}

		// Service fails
		mockService.On("CreateOrganization", mock.Anything, req).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/admin/organizations", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Perform the test
		err := handler.HandleCreateOrganization(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting all organizations
func TestHandleGetAllOrganizations(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Get All Organizations Success", func(t *testing.T) {
		// Mock data
		orgs := []*models.Organization{
			{
				ID:        "org-id-1",
				Name:      "Organization 1",
				Slogan:    "Slogan 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "org-id-2",
				Name:      "Organization 2",
				Slogan:    "Slogan 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Set up expectations
		mockService.On("GetAllOrganizations", mock.Anything).Return(orgs, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAllOrganizations(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.Organization
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, orgs[0].ID, response[0].ID)
			assert.Equal(t, orgs[1].ID, response[1].ID)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get All Organizations Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetAllOrganizations", mock.Anything).Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Perform the test
		err := handler.HandleGetAllOrganizations(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting an organization by ID
func TestHandleGetOrganizationByID(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Get Organization by ID Success", func(t *testing.T) {
		// Mock data
		org := &models.Organization{
			ID:        "org-id",
			Name:      "Test Organization",
			Slogan:    "Testing is Important",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set up expectations
		mockService.On("GetOrganizationByID", mock.Anything, "org-id").Return(org, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/org-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		if assert.NoError(t, handler.HandleGetOrganizationByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Organization
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, org.ID, response.ID)
			assert.Equal(t, org.Name, response.Name)
			assert.Equal(t, org.Slogan, response.Slogan)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Organization by ID Not Found", func(t *testing.T) {
		// Set up expectations - organization not found
		mockService.On("GetOrganizationByID", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/nonexistent-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Perform the test
		err := handler.HandleGetOrganizationByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Organization by ID Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetOrganizationByID", mock.Anything, "org-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/org-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		err := handler.HandleGetOrganizationByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test updating an organization
func TestHandleUpdateOrganization(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Update Organization Success", func(t *testing.T) {
		// Mock data
		name := "Updated Organization"
		slogan := "Updated Slogan"
		req := models.UpdateOrganizationRequest{
			Name:   &name,
			Slogan: &slogan,
		}

		org := &models.Organization{
			ID:        "org-id",
			Name:      name,
			Slogan:    slogan,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set up expectations
		mockService.On("UpdateOrganization", mock.Anything, "org-id", req).Return(org, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/organizations/org-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		if assert.NoError(t, handler.HandleUpdateOrganization(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Organization
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, org.ID, response.ID)
			assert.Equal(t, org.Name, response.Name)
			assert.Equal(t, org.Slogan, response.Slogan)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Update Organization Not Found", func(t *testing.T) {
		// Mock data
		name := "Updated Organization"
		slogan := "Updated Slogan"
		req := models.UpdateOrganizationRequest{
			Name:   &name,
			Slogan: &slogan,
		}

		// Set up expectations - organization not found
		mockService.On("UpdateOrganization", mock.Anything, "nonexistent-id", req).Return(nil, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/organizations/nonexistent-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Perform the test
		err := handler.HandleUpdateOrganization(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Update Organization Error", func(t *testing.T) {
		// Mock data
		name := "Updated Organization"
		slogan := "Updated Slogan"
		req := models.UpdateOrganizationRequest{
			Name:   &name,
			Slogan: &slogan,
		}

		// Service fails
		mockService.On("UpdateOrganization", mock.Anything, "org-id", req).Return(nil, errors.New("database error")).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/admin/organizations/org-id", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		err := handler.HandleUpdateOrganization(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test deleting an organization
func TestHandleDeleteOrganization(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Delete Organization Success", func(t *testing.T) {
		// Set up expectations
		mockService.On("DeleteOrganization", mock.Anything, "org-id").Return(nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/organizations/org-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		if assert.NoError(t, handler.HandleDeleteOrganization(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Delete Organization Error", func(t *testing.T) {
		// Service fails
		mockService.On("DeleteOrganization", mock.Anything, "org-id").Return(errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/admin/organizations/org-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		err := handler.HandleDeleteOrganization(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}

// Test getting organization statistics
func TestHandleGetOrganizationStats(t *testing.T) {
	e, mockService, handler := setupOrganizationTest(t)

	t.Run("Get Organization Stats Success", func(t *testing.T) {
		// Mock data
		org := &models.Organization{
			ID:        "org-id",
			Name:      "Test Organization",
			Slogan:    "Testing is Important",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		stats := &models.OrganizationStats{
			Organization: org,
			UserCount:    10,
			AdminCount:   2,
			TeacherCount: 3,
			StudentCount: 5,
			CourseCount:  5,
		}

		// Set up expectations
		mockService.On("GetOrganizationStats", mock.Anything, "org-id").Return(stats, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/org-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		if assert.NoError(t, handler.HandleGetOrganizationStats(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.OrganizationStats
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, stats.UserCount, response.UserCount)
			assert.Equal(t, stats.AdminCount, response.AdminCount)
			assert.Equal(t, stats.TeacherCount, response.TeacherCount)
			assert.Equal(t, stats.StudentCount, response.StudentCount)
			assert.Equal(t, stats.CourseCount, response.CourseCount)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Get Organization Stats Not Found", func(t *testing.T) {
		// Set up expectations - organization not found
		mockService.On("GetOrganizationStats", mock.Anything, "nonexistent-id").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/nonexistent-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		// Perform the test
		err := handler.HandleGetOrganizationStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Get Organization Stats Error", func(t *testing.T) {
		// Service fails
		mockService.On("GetOrganizationStats", mock.Anything, "org-id").Return(nil, errors.New("database error")).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/organizations/org-id/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("org-id")

		// Perform the test
		err := handler.HandleGetOrganizationStats(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockService.AssertExpectations(t)
	})
}
