package services

import (
	"context"
	"errors"
	"strings"

	"assessment-management-system/models"
	"assessment-management-system/repositories"
	"assessment-management-system/utils"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo       *repositories.UserRepository
	orgRepo        *repositories.OrganizationRepository
	courseRepo     *repositories.CourseRepository
	assessmentRepo *repositories.AssessmentRepository
}

// NewUserService creates a new UserService
func NewUserService(
	userRepo *repositories.UserRepository,
	orgRepo *repositories.OrganizationRepository,
	courseRepo *repositories.CourseRepository,
	assessmentRepo *repositories.AssessmentRepository,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		courseRepo:     courseRepo,
		assessmentRepo: assessmentRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, organizationID string, req models.CreateUserRequest) (*models.User, error) {
	// Validate organization
	org, err := s.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Check if email already exists in any organization (global uniqueness)
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.userRepo.Create(ctx, organizationID, req.Email, passwordHash, req.FirstName, req.LastName, string(req.Role))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsersByOrganization retrieves all users in an organization with optional filtering
func (s *UserService) GetUsersByOrganization(ctx context.Context, organizationID, role, search string) ([]*models.User, error) {
	// Validate organization
	org, err := s.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Check if role is valid
	if role != "" && role != string(models.RoleAdmin) && role != string(models.RoleTeacher) && role != string(models.RoleStudent) {
		return nil, errors.New("invalid role")
	}

	return s.userRepo.FindByOrganization(ctx, organizationID, role, search)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and is already taken (globally)
	if req.Email != nil && *req.Email != user.Email {
		// Check if email is already in use by any user
		existingUser, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}

		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email is already in use")
		}
	}

	// Hash password if provided
	var passwordHash string
	if req.Password != nil {
		passwordHash, err = utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
	}

	// Update user (passing nil for organizationID to keep it unchanged)
	updatedUser, err := s.userRepo.Update(ctx, id, req.Email, passwordHash, req.FirstName, req.LastName, req.Role, nil)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// If the user is a teacher, check if they are assigned to any courses
	if user.Role == models.RoleTeacher {
		courses, err := s.courseRepo.FindByTeacher(ctx, id)
		if err != nil {
			return err
		}

		if len(courses) > 0 {
			// Get the list of course names
			var courseNames []string
			for _, course := range courses {
				courseNames = append(courseNames, course.Name)
			}

			return errors.New("teacher is assigned to courses: " + strings.Join(courseNames, ", ") + ". Please reassign courses before deleting")
		}
	}

	// Delete the user
	return s.userRepo.Delete(ctx, id)
}

// BulkCreateUsers creates multiple users at once
func (s *UserService) BulkCreateUsers(ctx context.Context, organizationID string, users []models.CreateUserRequest) (map[string]interface{}, error) {
	// Validate organization
	org, err := s.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	results := map[string]interface{}{
		"successful": make([]models.User, 0),
		"failed":     make([]map[string]string, 0),
	}

	// Pre-check all emails to avoid partial creation if there are duplicates
	emailSet := make(map[string]bool)
	for _, req := range users {
		// First check if the email is a duplicate within the batch
		if emailSet[req.Email] {
			results["failed"] = append(results["failed"].([]map[string]string), map[string]string{
				"email": req.Email,
				"error": "duplicate email within batch",
			})
			continue
		}
		emailSet[req.Email] = true

		// Then check if the email already exists in the system
		existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			results["failed"] = append(results["failed"].([]map[string]string), map[string]string{
				"email": req.Email,
				"error": "error checking email: " + err.Error(),
			})
			continue
		}
		if existingUser != nil {
			results["failed"] = append(results["failed"].([]map[string]string), map[string]string{
				"email": req.Email,
				"error": "email is already in use",
			})
			continue
		}
	}

	// Process the users that passed pre-check
	for _, req := range users {
		// Skip users that failed pre-check
		skip := false
		for _, failed := range results["failed"].([]map[string]string) {
			if failed["email"] == req.Email {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		user, err := s.CreateUser(ctx, organizationID, req)
		if err != nil {
			results["failed"] = append(results["failed"].([]map[string]string), map[string]string{
				"email": req.Email,
				"error": err.Error(),
			})
		} else {
			results["successful"] = append(results["successful"].([]models.User), *user)
		}
	}

	return results, nil
}

// GetUserOrganization retrieves a user's organization
func (s *UserService) GetUserOrganization(ctx context.Context, userID string) (*models.Organization, error) {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Get organization
	return s.orgRepo.FindByID(ctx, user.OrganizationID)
}

// GetTeacherStats retrieves statistics for a teacher
func (s *UserService) GetTeacherStats(ctx context.Context, teacherID, organizationID string) (*models.TeacherStats, error) {
	// Check if user exists and is a teacher in the given organization
	user, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	if user == nil || user.Role != models.RoleTeacher || user.OrganizationID != organizationID {
		return nil, errors.New("teacher not found in the organization")
	}

	// Get assigned courses count
	assignedCourses, err := s.courseRepo.CountByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// Get created assessments count
	createdAssessments, err := s.assessmentRepo.CountByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// Get pending grading count
	pendingGrading, err := s.assessmentRepo.CountUngradedSubmissionsByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// Create the stats object
	stats := &models.TeacherStats{
		User:               user,
		AssignedCourses:    assignedCourses,
		CreatedAssessments: createdAssessments,
		PendingGrading:     pendingGrading,
	}

	return stats, nil
}

// GetStudentStats retrieves statistics for a student
func (s *UserService) GetStudentStats(ctx context.Context, studentID, organizationID string) (*models.StudentStats, error) {
	// Check if user exists and is a student in the given organization
	user, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if user == nil || user.Role != models.RoleStudent || user.OrganizationID != organizationID {
		return nil, errors.New("student not found in the organization")
	}

	// Get enrolled courses count
	enrolledCourses, err := s.courseRepo.CountByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Get completed assessments count
	completedAssessments, err := s.assessmentRepo.CountSubmissionsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Get pending assessments count
	pendingAssessments, err := s.assessmentRepo.CountPendingAssessmentsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Get average grade
	averageGrade, err := s.assessmentRepo.GetAverageGradeForStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Create the stats object
	stats := &models.StudentStats{
		User:                 user,
		EnrolledCourses:      enrolledCourses,
		CompletedAssessments: completedAssessments,
		PendingAssessments:   pendingAssessments,
		AverageGrade:         averageGrade,
	}

	return stats, nil
}
