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
	"assessment-management-system/middleware"
	"assessment-management-system/models"
)

// MockAssessmentService is a mock implementation of AssessmentService
type MockAssessmentService struct {
	mock.Mock
}

func (m *MockAssessmentService) CreateAssessment(ctx context.Context, teacherID string, req models.CreateAssessmentRequest) (*models.Assessment, error) {
	args := m.Called(ctx, teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Assessment), args.Error(1)
}

func (m *MockAssessmentService) GetAssessmentsByOrganization(ctx context.Context, organizationID, courseID string) ([]*models.AssessmentWithDetails, error) {
	args := m.Called(ctx, organizationID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AssessmentWithDetails), args.Error(1)
}

func (m *MockAssessmentService) GetAssessmentsByCourse(ctx context.Context, courseID string) ([]*models.Assessment, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Assessment), args.Error(1)
}

func (m *MockAssessmentService) GetAssessmentByID(ctx context.Context, id string) (*models.Assessment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Assessment), args.Error(1)
}

func (m *MockAssessmentService) UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Assessment), args.Error(1)
}

func (m *MockAssessmentService) DeleteAssessment(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAssessmentService) GetAssessmentSubmissions(ctx context.Context, assessmentID string) ([]*models.AssessmentSubmission, error) {
	args := m.Called(ctx, assessmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AssessmentSubmission), args.Error(1)
}

func (m *MockAssessmentService) GetSubmissionByID(ctx context.Context, id string) (*models.AssessmentSubmission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssessmentSubmission), args.Error(1)
}

func (m *MockAssessmentService) GetStudentSubmission(ctx context.Context, assessmentID, studentID string) (*models.AssessmentSubmission, error) {
	args := m.Called(ctx, assessmentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssessmentSubmission), args.Error(1)
}

func (m *MockAssessmentService) HasStudentSubmitted(ctx context.Context, assessmentID, studentID string) (bool, error) {
	args := m.Called(ctx, assessmentID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAssessmentService) SubmitAssessment(ctx context.Context, assessmentID, studentID, content string) (*models.AssessmentSubmission, error) {
	args := m.Called(ctx, assessmentID, studentID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssessmentSubmission), args.Error(1)
}

func (m *MockAssessmentService) GetSubmissionGrade(ctx context.Context, submissionID string) (*models.Grade, error) {
	args := m.Called(ctx, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Grade), args.Error(1)
}

func (m *MockAssessmentService) GradeSubmission(ctx context.Context, submissionID, teacherID string, score float64, feedback string) (*models.Grade, error) {
	args := m.Called(ctx, submissionID, teacherID, score, feedback)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Grade), args.Error(1)
}

func (m *MockAssessmentService) GetStudentAssessmentStatus(ctx context.Context, assessmentID, studentID string) (*models.StudentAssessmentStatus, error) {
	args := m.Called(ctx, assessmentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentAssessmentStatus), args.Error(1)
}

// Helper functions to setup the test environment
func setupAdminAssessmentTest(t *testing.T) (*echo.Echo, *MockAssessmentService, *MockCourseService, *admin.AssessmentHandler) {
	e := echo.New()
	mockAssessmentService := new(MockAssessmentService)
	mockCourseService := new(MockCourseService)
	handler := admin.NewAssessmentHandler(mockAssessmentService, mockCourseService)
	return e, mockAssessmentService, mockCourseService, handler
}

func setupTeacherAssessmentTest(t *testing.T) (*echo.Echo, *MockAssessmentService, *MockCourseService, *teacher.AssessmentHandler) {
	e := echo.New()
	mockAssessmentService := new(MockAssessmentService)
	mockCourseService := new(MockCourseService)
	handler := teacher.NewAssessmentHandler(mockAssessmentService, mockCourseService)
	return e, mockAssessmentService, mockCourseService, handler
}

func setupStudentAssessmentTest(t *testing.T) (*echo.Echo, *MockAssessmentService, *MockCourseService, *student.AssessmentHandler) {
	e := echo.New()
	mockAssessmentService := new(MockAssessmentService)
	mockCourseService := new(MockCourseService)
	handler := student.NewAssessmentHandler(mockAssessmentService, mockCourseService)
	return e, mockAssessmentService, mockCourseService, handler
}

// Test admin assessment handler
func TestAdminAssessmentHandler(t *testing.T) {
	e, mockAssessmentService, mockCourseService, handler := setupAdminAssessmentTest(t)

	t.Run("GetAllAssessments", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessments := []*models.AssessmentWithDetails{
			{
				Assessment: &models.Assessment{
					ID:          "assessment1",
					CourseID:    "course1",
					TeacherID:   "teacher1",
					Title:       "Test Assessment 1",
					Description: "Description for test assessment 1",
					Type:        models.AssessmentTypeQuiz,
					MaxScore:    100,
					DueDate:     &dueDate,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				TeacherName: "John Doe",
				CourseName:  "Mathematics",
			},
			{
				Assessment: &models.Assessment{
					ID:          "assessment2",
					CourseID:    "course2",
					TeacherID:   "teacher2",
					Title:       "Test Assessment 2",
					Description: "Description for test assessment 2",
					Type:        models.AssessmentTypeExam,
					MaxScore:    50,
					DueDate:     &dueDate,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				TeacherName: "Jane Smith",
				CourseName:  "Science",
			},
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentsByOrganization", mock.Anything, "org-id", "").Return(assessments, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/assessments", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAllAssessments(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.AssessmentWithDetails
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, assessments[0].Assessment.ID, response[0].Assessment.ID)
			assert.Equal(t, assessments[1].Assessment.ID, response[1].Assessment.ID)
		}

		mockAssessmentService.AssertExpectations(t)
	})

	t.Run("GetAssessmentByID", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		course := &models.Course{
			ID:             "course1",
			OrganizationID: "org-id",
			Name:           "Math Course",
			Description:    "Mathematics course",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("GetCourseByID", mock.Anything, "course1").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/assessments/assessment1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAssessmentByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Assessment
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, assessment.ID, response.ID)
			assert.Equal(t, assessment.Title, response.Title)
			assert.Equal(t, assessment.Type, response.Type)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("GetAssessmentByID_FromDifferentOrganization", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		course := &models.Course{
			ID:             "course1",
			OrganizationID: "different-org-id", // Different from admin's org
			Name:           "Math Course",
			Description:    "Mathematics course",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("GetCourseByID", mock.Anything, "course1").Return(course, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/assessments/assessment1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		err := handler.HandleGetAssessmentByID(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("GetAssessmentSubmissions", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		course := &models.Course{
			ID:             "course1",
			OrganizationID: "org-id",
			Name:           "Math Course",
			Description:    "Mathematics course",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		submissions := []*models.AssessmentSubmission{
			{
				ID:           "submission1",
				AssessmentID: "assessment1",
				StudentID:    "student1",
				Content:      "This is student 1's submission",
				SubmittedAt:  time.Now(),
			},
			{
				ID:           "submission2",
				AssessmentID: "assessment1",
				StudentID:    "student2",
				Content:      "This is student 2's submission",
				SubmittedAt:  time.Now(),
			},
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("GetCourseByID", mock.Anything, "course1").Return(course, nil).Once()
		mockAssessmentService.On("GetAssessmentSubmissions", mock.Anything, "assessment1").Return(submissions, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/assessments/assessment1/submissions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAssessmentSubmissions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.AssessmentSubmission
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, submissions[0].ID, response[0].ID)
			assert.Equal(t, submissions[1].ID, response[1].ID)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("GetSubmissionGrades", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		course := &models.Course{
			ID:             "course1",
			OrganizationID: "org-id",
			Name:           "Math Course",
			Description:    "Mathematics course",
			EnrollmentOpen: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student1",
			Content:      "This is student 1's submission",
			SubmittedAt:  time.Now(),
		}

		grade := &models.Grade{
			SubmissionID: "submission1",
			Score:        85.5,
			Feedback:     "Good work, but could improve in some areas",
			GradedBy:     "teacher1",
			GradedAt:     time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetSubmissionByID", mock.Anything, "submission1").Return(submission, nil).Once()
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("GetCourseByID", mock.Anything, "course1").Return(course, nil).Once()
		mockAssessmentService.On("GetSubmissionGrade", mock.Anything, "submission1").Return(grade, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/admin/submissions/submission1/grade", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("submissionId")
		c.SetParamValues("submission1")

		// Mock the admin user context
		mockAdminContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetSubmissionGrades(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Grade
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, grade.SubmissionID, response.SubmissionID)
			assert.Equal(t, grade.Score, response.Score)
			assert.Equal(t, grade.Feedback, response.Feedback)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})
}

// Test teacher assessment handler
func TestTeacherAssessmentHandler(t *testing.T) {
	e, mockAssessmentService, mockCourseService, handler := setupTeacherAssessmentTest(t)

	t.Run("CreateAssessment", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		req := models.CreateAssessmentRequest{
			CourseID:    "course1",
			Title:       "New Assessment",
			Description: "Description for new assessment",
			Type:        models.AssessmentTypeAssignment,
			MaxScore:    100,
			DueDate:     &dueDate,
		}

		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       "New Assessment",
			Description: "Description for new assessment",
			Type:        models.AssessmentTypeAssignment,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Set up expectations
		mockCourseService.On("IsTeacherAssignedToCourse", mock.Anything, "course1", "teacher-id").Return(true, nil).Once()
		mockAssessmentService.On("CreateAssessment", mock.Anything, "teacher-id", req).Return(assessment, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/teacher/assessments", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleCreateAssessment(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.Assessment
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, assessment.ID, response.ID)
			assert.Equal(t, assessment.Title, response.Title)
			assert.Equal(t, assessment.Type, response.Type)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("CreateAssessment_NotAssignedToCourse", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		req := models.CreateAssessmentRequest{
			CourseID:    "course1",
			Title:       "New Assessment",
			Description: "Description for new assessment",
			Type:        models.AssessmentTypeAssignment,
			MaxScore:    100,
			DueDate:     &dueDate,
		}

		// Set up expectations
		mockCourseService.On("IsTeacherAssignedToCourse", mock.Anything, "course1", "teacher-id").Return(false, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/teacher/assessments", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		err := handler.HandleCreateAssessment(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockCourseService.AssertExpectations(t)
	})

	t.Run("GetAssessments", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessments := []*models.Assessment{
			{
				ID:          "assessment1",
				CourseID:    "course1",
				TeacherID:   "teacher-id",
				Title:       "Test Assessment 1",
				Description: "Description for test assessment 1",
				Type:        models.AssessmentTypeQuiz,
				MaxScore:    100,
				DueDate:     &dueDate,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "assessment2",
				CourseID:    "course1",
				TeacherID:   "teacher-id",
				Title:       "Test Assessment 2",
				Description: "Description for test assessment 2",
				Type:        models.AssessmentTypeExam,
				MaxScore:    50,
				DueDate:     &dueDate,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		// Set up expectations
		mockCourseService.On("IsTeacherAssignedToCourse", mock.Anything, "course1", "teacher-id").Return(true, nil).Once()
		mockAssessmentService.On("GetAssessmentsByCourse", mock.Anything, "course1").Return(assessments, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/teacher/courses/course1/assessments", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("courseId")
		c.SetParamValues("course1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAssessments(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []*models.Assessment
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, assessments[0].ID, response[0].ID)
			assert.Equal(t, assessments[1].ID, response[1].ID)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("UpdateAssessment", func(t *testing.T) {
		// Mock data
		title := "Updated Assessment"
		description := "Updated description"
		assessmentType := models.AssessmentTypeProject
		maxScore := 150
		dueDate := time.Now().Add(time.Hour * 24 * 14) // two weeks from now

		req := models.UpdateAssessmentRequest{
			Title:       &title,
			Description: &description,
			Type:        &assessmentType,
			MaxScore:    &maxScore,
			DueDate:     &dueDate,
		}

		originalAssessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       "Original Assessment",
			Description: "Original description",
			Type:        models.AssessmentTypeAssignment,
			MaxScore:    100,
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		updatedAssessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       title,
			Description: description,
			Type:        assessmentType,
			MaxScore:    maxScore,
			DueDate:     &dueDate,
			CreatedAt:   originalAssessment.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(originalAssessment, nil).Once()
		mockAssessmentService.On("UpdateAssessment", mock.Anything, "assessment1", req).Return(updatedAssessment, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/teacher/assessments/assessment1", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleUpdateAssessment(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Assessment
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, updatedAssessment.ID, response.ID)
			assert.Equal(t, updatedAssessment.Title, response.Title)
			assert.Equal(t, updatedAssessment.Description, response.Description)
			assert.Equal(t, updatedAssessment.Type, response.Type)
			assert.Equal(t, updatedAssessment.MaxScore, response.MaxScore)
			// Cannot directly compare time.Time objects due to possible precision issues during JSON marshalling
			assert.NotNil(t, response.DueDate)
		}

		mockAssessmentService.AssertExpectations(t)
	})

	t.Run("UpdateAssessment_NotOwner", func(t *testing.T) {
		// Mock data - assessment created by different teacher
		title := "Updated Assessment"
		req := models.UpdateAssessmentRequest{
			Title: &title,
		}

		originalAssessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "different-teacher-id", // Different from the requesting teacher
			Title:       "Original Assessment",
			Description: "Original description",
			Type:        models.AssessmentTypeAssignment,
			MaxScore:    100,
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(originalAssessment, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/api/teacher/assessments/assessment1", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		err := handler.HandleUpdateAssessment(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockAssessmentService.AssertExpectations(t)
	})

	t.Run("DeleteAssessment", func(t *testing.T) {
		// Mock data
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockAssessmentService.On("DeleteAssessment", mock.Anything, "assessment1").Return(nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodDelete, "/api/teacher/assessments/assessment1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleDeleteAssessment(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}

		mockAssessmentService.AssertExpectations(t)
	})

	t.Run("GradeSubmission", func(t *testing.T) {
		// Mock data
		req := models.GradeSubmissionRequest{
			Score:    85.5,
			Feedback: "Good work, could improve in some areas",
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student1",
			Content:      "This is student 1's submission",
			SubmittedAt:  time.Now(),
		}

		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		grade := &models.Grade{
			SubmissionID: "submission1",
			Score:        85.5,
			Feedback:     "Good work, could improve in some areas",
			GradedBy:     "teacher-id",
			GradedAt:     time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetSubmissionByID", mock.Anything, "submission1").Return(submission, nil).Once()
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsTeacherAssignedToCourse", mock.Anything, "course1", "teacher-id").Return(true, nil).Once()
		mockAssessmentService.On("GradeSubmission", mock.Anything, "submission1", "teacher-id", 85.5, "Good work, could improve in some areas").Return(grade, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/teacher/submissions/submission1/grade", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("submissionId")
		c.SetParamValues("submission1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGradeSubmission(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Grade
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, grade.SubmissionID, response.SubmissionID)
			assert.Equal(t, grade.Score, response.Score)
			assert.Equal(t, grade.Feedback, response.Feedback)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("GradeSubmission_ScoreExceedsMaximum", func(t *testing.T) {
		// Mock data - score exceeds max score
		req := models.GradeSubmissionRequest{
			Score:    150, // Max score is 100
			Feedback: "Good work",
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student1",
			Content:      "This is student 1's submission",
			SubmittedAt:  time.Now(),
		}

		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher-id",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100, // Max score is 100
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetSubmissionByID", mock.Anything, "submission1").Return(submission, nil).Once()
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsTeacherAssignedToCourse", mock.Anything, "course1", "teacher-id").Return(true, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/teacher/submissions/submission1/grade", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("submissionId")
		c.SetParamValues("submission1")

		// Mock the teacher user context
		mockTeacherContext(c)

		// Perform the test
		err := handler.HandleGradeSubmission(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})
}

// Test student assessment handlers
func TestStudentAssessmentHandler(t *testing.T) {
	e, mockAssessmentService, mockCourseService, handler := setupStudentAssessmentTest(t)

	t.Run("GetCourseAssessments", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessments := []*models.Assessment{
			{
				ID:          "assessment1",
				CourseID:    "course1",
				TeacherID:   "teacher1",
				Title:       "Test Assessment 1",
				Description: "Description for test assessment 1",
				Type:        models.AssessmentTypeQuiz,
				MaxScore:    100,
				DueDate:     &dueDate,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "assessment2",
				CourseID:    "course1",
				TeacherID:   "teacher1",
				Title:       "Test Assessment 2",
				Description: "Description for test assessment 2",
				Type:        models.AssessmentTypeExam,
				MaxScore:    50,
				DueDate:     &dueDate,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		statuses := []models.StudentAssessmentStatus{
			{
				Assessment:   assessments[0],
				Submission:   nil,
				Grade:        nil,
				HasSubmitted: false,
				IsGraded:     false,
				DaysUntilDue: 7,
				IsOverdue:    false,
			},
			{
				Assessment:   assessments[1],
				Submission:   nil,
				Grade:        nil,
				HasSubmitted: false,
				IsGraded:     false,
				DaysUntilDue: 7,
				IsOverdue:    false,
			},
		}

		// Set up expectations
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("GetAssessmentsByCourse", mock.Anything, "course1").Return(assessments, nil).Once()
		mockAssessmentService.On("GetStudentAssessmentStatus", mock.Anything, "assessment1", "student-id").Return(&statuses[0], nil).Once()
		mockAssessmentService.On("GetStudentAssessmentStatus", mock.Anything, "assessment2", "student-id").Return(&statuses[1], nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/courses/course1/assessments", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("courseId")
		c.SetParamValues("course1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetCourseAssessments(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []models.StudentAssessmentStatus
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, assessments[0].ID, response[0].Assessment.ID)
			assert.Equal(t, assessments[1].ID, response[1].Assessment.ID)
			assert.False(t, response[0].HasSubmitted)
			assert.False(t, response[0].IsGraded)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("GetAssessmentByID", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		status := &models.StudentAssessmentStatus{
			Assessment:   assessment,
			Submission:   nil,
			Grade:        nil,
			HasSubmitted: false,
			IsGraded:     false,
			DaysUntilDue: 7,
			IsOverdue:    false,
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("GetStudentAssessmentStatus", mock.Anything, "assessment1", "student-id").Return(status, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/assessments/assessment1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleGetAssessmentByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.StudentAssessmentStatus
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, assessment.ID, response.Assessment.ID)
			assert.Equal(t, assessment.Title, response.Assessment.Title)
			assert.False(t, response.HasSubmitted)
			assert.False(t, response.IsGraded)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("SubmitAssessment", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		req := models.CreateSubmissionRequest{
			Content: "This is my assessment submission",
		}

		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student-id",
			Content:      "This is my assessment submission",
			SubmittedAt:  time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("HasStudentSubmitted", mock.Anything, "assessment1", "student-id").Return(false, nil).Once()
		mockAssessmentService.On("SubmitAssessment", mock.Anything, "assessment1", "student-id", "This is my assessment submission").Return(submission, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/student/assessments/assessment1/submit", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleSubmitAssessment(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.AssessmentSubmission
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, submission.ID, response.ID)
			assert.Equal(t, submission.Content, response.Content)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("SubmitAssessment_AlreadySubmitted", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		req := models.CreateSubmissionRequest{
			Content: "This is my assessment submission",
		}

		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Set up expectations - student has already submitted
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("HasStudentSubmitted", mock.Anything, "assessment1", "student-id").Return(true, nil).Once()

		// Create request
		jsonBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/student/assessments/assessment1/submit", bytes.NewReader(jsonBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		err := handler.HandleSubmitAssessment(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpError.Code)

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("ViewSubmission", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student-id",
			Content:      "This is my assessment submission",
			SubmittedAt:  time.Now(),
		}

		grade := &models.Grade{
			SubmissionID: "submission1",
			Score:        85.5,
			Feedback:     "Good work, could improve in some areas",
			GradedBy:     "teacher1",
			GradedAt:     time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("GetStudentSubmission", mock.Anything, "assessment1", "student-id").Return(submission, nil).Once()
		mockAssessmentService.On("GetSubmissionGrade", mock.Anything, "submission1").Return(grade, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/assessments/assessment1/submission", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleViewSubmission(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response, "submission")
			assert.Contains(t, response, "grade")

			// Check some key details
			submissionData := response["submission"].(map[string]interface{})
			gradeData := response["grade"].(map[string]interface{})

			assert.Equal(t, "submission1", submissionData["id"])
			assert.Equal(t, 85.5, gradeData["score"])
			assert.Equal(t, "Good work, could improve in some areas", gradeData["feedback"])
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("ViewGrade", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student-id",
			Content:      "This is my assessment submission",
			SubmittedAt:  time.Now(),
		}

		grade := &models.Grade{
			SubmissionID: "submission1",
			Score:        85.5,
			Feedback:     "Good work, could improve in some areas",
			GradedBy:     "teacher1",
			GradedAt:     time.Now(),
		}

		// Set up expectations
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("GetStudentSubmission", mock.Anything, "assessment1", "student-id").Return(submission, nil).Once()
		mockAssessmentService.On("GetSubmissionGrade", mock.Anything, "submission1").Return(grade, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/assessments/assessment1/grade", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		if assert.NoError(t, handler.HandleViewGrade(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Grade
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, grade.SubmissionID, response.SubmissionID)
			assert.Equal(t, grade.Score, response.Score)
			assert.Equal(t, grade.Feedback, response.Feedback)
		}

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})

	t.Run("ViewGrade_NotGraded", func(t *testing.T) {
		// Mock data
		dueDate := time.Now().Add(time.Hour * 24 * 7) // one week from now
		assessment := &models.Assessment{
			ID:          "assessment1",
			CourseID:    "course1",
			TeacherID:   "teacher1",
			Title:       "Test Assessment",
			Description: "Description for test assessment",
			Type:        models.AssessmentTypeQuiz,
			MaxScore:    100,
			DueDate:     &dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		submission := &models.AssessmentSubmission{
			ID:           "submission1",
			AssessmentID: "assessment1",
			StudentID:    "student-id",
			Content:      "This is my assessment submission",
			SubmittedAt:  time.Now(),
		}

		// Set up expectations - submission not graded yet
		mockAssessmentService.On("GetAssessmentByID", mock.Anything, "assessment1").Return(assessment, nil).Once()
		mockCourseService.On("IsStudentEnrolledInCourse", mock.Anything, "course1", "student-id").Return(true, nil).Once()
		mockAssessmentService.On("GetStudentSubmission", mock.Anything, "assessment1", "student-id").Return(submission, nil).Once()
		mockAssessmentService.On("GetSubmissionGrade", mock.Anything, "submission1").Return(nil, nil).Once()

		// Create request
		httpReq := httptest.NewRequest(http.MethodGet, "/api/student/assessments/assessment1/grade", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetParamNames("id")
		c.SetParamValues("assessment1")

		// Mock the student user context
		mockStudentContext(c)

		// Perform the test
		err := handler.HandleViewGrade(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)

		mockAssessmentService.AssertExpectations(t)
		mockCourseService.AssertExpectations(t)
	})
}
