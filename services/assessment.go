package services

import (
	"context"
	"errors"
	"time"

	"assessment-management-system/models"
	"assessment-management-system/repositories"
)

// AssessmentService handles assessment-related business logic
type AssessmentService struct {
	assessmentRepo *repositories.AssessmentRepository
	courseRepo     *repositories.CourseRepository
	userRepo       *repositories.UserRepository
}

// NewAssessmentService creates a new AssessmentService
func NewAssessmentService(
	assessmentRepo *repositories.AssessmentRepository,
	courseRepo *repositories.CourseRepository,
	userRepo *repositories.UserRepository,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: assessmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
	}
}

// CreateAssessment creates a new assessment
func (s *AssessmentService) CreateAssessment(ctx context.Context, teacherID string, req models.CreateAssessmentRequest) (*models.Assessment, error) {
	// Validate teacher
	teacher, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	if teacher == nil || teacher.Role != models.RoleTeacher {
		return nil, errors.New("invalid teacher")
	}

	// Validate course
	course, err := s.courseRepo.FindByID(ctx, req.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Check if teacher is assigned to the course
	isAssigned, err := s.courseRepo.IsTeacherAssigned(ctx, req.CourseID, teacherID)
	if err != nil {
		return nil, err
	}

	if !isAssigned {
		return nil, errors.New("teacher is not assigned to this course")
	}

	// Create assessment
	assessment, err := s.assessmentRepo.Create(ctx, req.CourseID, teacherID, req.Title, req.Description, string(req.Type), req.MaxScore, req.DueDate)
	if err != nil {
		return nil, err
	}

	return assessment, nil
}

// GetAssessmentsByOrganization retrieves all assessments for an organization
func (s *AssessmentService) GetAssessmentsByOrganization(ctx context.Context, organizationID, courseID string) ([]*models.AssessmentWithDetails, error) {
	assessments, err := s.assessmentRepo.FindByOrganization(ctx, organizationID, courseID)
	if err != nil {
		return nil, err
	}

	// Enrich with additional details
	var assessmentsWithDetails []*models.AssessmentWithDetails
	for _, assessment := range assessments {
		teacher, err := s.userRepo.FindByID(ctx, assessment.TeacherID)
		if err != nil {
			return nil, err
		}

		teacherName := "Unknown"
		if teacher != nil {
			teacherName = teacher.FirstName + " " + teacher.LastName
		}

		course, err := s.courseRepo.FindByID(ctx, assessment.CourseID)
		if err != nil {
			return nil, err
		}

		courseName := "Unknown"
		if course != nil {
			courseName = course.Name
		}

		assessmentsWithDetails = append(assessmentsWithDetails, &models.AssessmentWithDetails{
			Assessment:  assessment,
			TeacherName: teacherName,
			CourseName:  courseName,
		})
	}

	return assessmentsWithDetails, nil
}

// GetAssessmentsByCourse retrieves all assessments for a course
func (s *AssessmentService) GetAssessmentsByCourse(ctx context.Context, courseID string) ([]*models.Assessment, error) {
	// Validate course
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Get assessments
	return s.assessmentRepo.FindByCourse(ctx, courseID)
}

// GetAssessmentByID retrieves an assessment by ID
func (s *AssessmentService) GetAssessmentByID(ctx context.Context, id string) (*models.Assessment, error) {
	return s.assessmentRepo.FindByID(ctx, id)
}

// UpdateAssessment updates an assessment
func (s *AssessmentService) UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error) {
	// Check if assessment exists
	assessment, err := s.assessmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	// Update assessment
	updatedAssessment, err := s.assessmentRepo.Update(ctx, id, req.Title, req.Description, req.Type, req.MaxScore, req.DueDate)
	if err != nil {
		return nil, err
	}

	return updatedAssessment, nil
}

// DeleteAssessment deletes an assessment
func (s *AssessmentService) DeleteAssessment(ctx context.Context, id string) error {
	// Check if assessment exists
	assessment, err := s.assessmentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if assessment == nil {
		return errors.New("assessment not found")
	}

	// Delete assessment (cascade will handle submissions and grades)
	return s.assessmentRepo.Delete(ctx, id)
}

// GetAssessmentSubmissions retrieves all submissions for an assessment
func (s *AssessmentService) GetAssessmentSubmissions(ctx context.Context, assessmentID string) ([]*models.AssessmentSubmission, error) {
	// Check if assessment exists
	assessment, err := s.assessmentRepo.FindByID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	// Get submissions
	return s.assessmentRepo.FindSubmissionsByAssessment(ctx, assessmentID)
}

// GetSubmissionByID retrieves a submission by ID
func (s *AssessmentService) GetSubmissionByID(ctx context.Context, id string) (*models.AssessmentSubmission, error) {
	return s.assessmentRepo.FindSubmissionByID(ctx, id)
}

// GetStudentSubmission retrieves a student's submission for an assessment
func (s *AssessmentService) GetStudentSubmission(ctx context.Context, assessmentID, studentID string) (*models.AssessmentSubmission, error) {
	return s.assessmentRepo.FindSubmissionByStudentAndAssessment(ctx, assessmentID, studentID)
}

// HasStudentSubmitted checks if a student has submitted an assessment
func (s *AssessmentService) HasStudentSubmitted(ctx context.Context, assessmentID, studentID string) (bool, error) {
	submission, err := s.assessmentRepo.FindSubmissionByStudentAndAssessment(ctx, assessmentID, studentID)
	if err != nil {
		return false, err
	}

	return submission != nil, nil
}

// SubmitAssessment submits an assessment
func (s *AssessmentService) SubmitAssessment(ctx context.Context, assessmentID, studentID, content string) (*models.AssessmentSubmission, error) {
	// Check if assessment exists
	assessment, err := s.assessmentRepo.FindByID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	// Check if student exists
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if student == nil || student.Role != models.RoleStudent {
		return nil, errors.New("invalid student")
	}

	// Check if student is enrolled in the course
	isEnrolled, err := s.courseRepo.IsStudentEnrolled(ctx, assessment.CourseID, studentID)
	if err != nil {
		return nil, err
	}

	if !isEnrolled {
		return nil, errors.New("student is not enrolled in this course")
	}

	// Check if student has already submitted
	hasSubmitted, err := s.HasStudentSubmitted(ctx, assessmentID, studentID)
	if err != nil {
		return nil, err
	}

	if hasSubmitted {
		return nil, errors.New("student has already submitted this assessment")
	}

	// Check if assessment is past due
	if assessment.DueDate != nil && assessment.DueDate.Before(time.Now()) {
		return nil, errors.New("assessment is past due")
	}

	// Create submission
	submission, err := s.assessmentRepo.CreateSubmission(ctx, assessmentID, studentID, content)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

// GetSubmissionGrade retrieves the grade for a submission
func (s *AssessmentService) GetSubmissionGrade(ctx context.Context, submissionID string) (*models.Grade, error) {
	return s.assessmentRepo.FindGradeBySubmission(ctx, submissionID)
}

// GradeSubmission grades a submission
func (s *AssessmentService) GradeSubmission(ctx context.Context, submissionID, teacherID string, score float64, feedback string) (*models.Grade, error) {
	// Check if submission exists
	submission, err := s.assessmentRepo.FindSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, errors.New("submission not found")
	}

	// Check if teacher exists
	teacher, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	if teacher == nil || teacher.Role != models.RoleTeacher {
		return nil, errors.New("invalid teacher")
	}

	// Get assessment to check max score
	assessment, err := s.assessmentRepo.FindByID(ctx, submission.AssessmentID)
	if err != nil {
		return nil, err
	}

	// Validate score
	if score < 0 || score > float64(assessment.MaxScore) {
		return nil, errors.New("score must be between 0 and the maximum score")
	}

	// Check if submission is already graded
	existingGrade, err := s.assessmentRepo.FindGradeBySubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	var grade *models.Grade
	if existingGrade == nil {
		// Create grade
		grade, err = s.assessmentRepo.CreateGrade(ctx, submissionID, score, feedback, teacherID)
	} else {
		// Update grade
		grade, err = s.assessmentRepo.UpdateGrade(ctx, submissionID, score, feedback, teacherID)
	}

	if err != nil {
		return nil, err
	}

	return grade, nil
}

// GetStudentAssessmentStatus retrieves a student's status for an assessment
func (s *AssessmentService) GetStudentAssessmentStatus(ctx context.Context, assessmentID, studentID string) (*models.StudentAssessmentStatus, error) {
	// Check if assessment exists
	assessment, err := s.assessmentRepo.FindByID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	// Get student's submission
	submission, err := s.assessmentRepo.FindSubmissionByStudentAndAssessment(ctx, assessmentID, studentID)
	if err != nil {
		return nil, err
	}

	// Check submission status
	hasSubmitted := submission != nil

	// Get grade if submitted
	var grade *models.Grade
	var isGraded bool
	if hasSubmitted {
		grade, err = s.assessmentRepo.FindGradeBySubmission(ctx, submission.ID)
		if err != nil {
			return nil, err
		}
		isGraded = grade != nil
	}

	// Calculate days until due
	var daysUntilDue int
	var isOverdue bool
	if assessment.DueDate != nil {
		if assessment.DueDate.After(time.Now()) {
			daysUntilDue = int(assessment.DueDate.Sub(time.Now()).Hours() / 24)
		} else {
			isOverdue = true
		}
	}

	// Create status object
	status := &models.StudentAssessmentStatus{
		Assessment:   assessment,
		Submission:   submission,
		Grade:        grade,
		HasSubmitted: hasSubmitted,
		IsGraded:     isGraded,
		DaysUntilDue: daysUntilDue,
		IsOverdue:    isOverdue,
	}

	return status, nil
}
