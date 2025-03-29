package services

import (
	"assessment-management-system/models"
	"assessment-management-system/repositories"
	"context"
	"errors"
	"fmt"
)

// CourseService handles course-related business logic
type CourseService struct {
	courseRepo *repositories.CourseRepository
	userRepo   *repositories.UserRepository
	orgRepo    *repositories.OrganizationRepository
}

// NewCourseService creates a new CourseService
func NewCourseService(
	courseRepo *repositories.CourseRepository,
	userRepo *repositories.UserRepository,
	orgRepo *repositories.OrganizationRepository,
) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		userRepo:   userRepo,
		orgRepo:    orgRepo,
	}
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(ctx context.Context, organizationID string, req models.CreateCourseRequest) (*models.Course, error) {
	// Validate organization
	org, err := s.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Check if course name already exists in the organization
	existingCourse, err := s.courseRepo.FindByNameAndOrganization(ctx, req.Name, organizationID)
	if err != nil {
		return nil, err
	}

	if existingCourse != nil {
		return nil, errors.New("course name already exists in this organization")
	}

	if req.EnrollmentOpen {
		req.EnrollmentOpen = false
	}

	// Create course
	course, err := s.courseRepo.Create(ctx, organizationID, req.Name, req.Description, req.EnrollmentOpen)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// GetCoursesByOrganization retrieves all courses in an organization
func (s *CourseService) GetCoursesByOrganization(ctx context.Context, organizationID string) ([]*models.CourseWithStudentCount, error) {
	// Validate organization
	org, err := s.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Get courses
	courses, err := s.courseRepo.FindByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	// For each course, get the student and teacher count
	var coursesWithCount []*models.CourseWithStudentCount
	for _, course := range courses {
		studentCount, err := s.courseRepo.CountStudentsByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		teacherCount, err := s.courseRepo.CountTeachersByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		coursesWithCount = append(coursesWithCount, &models.CourseWithStudentCount{
			Course:       course,
			StudentCount: studentCount,
			TeacherCount: teacherCount,
		})
	}

	return coursesWithCount, nil
}

// GetCourseByID retrieves a course by ID
func (s *CourseService) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	return s.courseRepo.FindByID(ctx, id)
}

// UpdateCourse updates a course
func (s *CourseService) UpdateCourse(ctx context.Context, id string, req models.UpdateCourseRequest) (*models.Course, error) {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Check if name is being changed and is already taken
	if req.Name != nil && *req.Name != course.Name {
		existingCourse, err := s.courseRepo.FindByNameAndOrganization(ctx, *req.Name, course.OrganizationID)
		if err != nil {
			return nil, err
		}

		if existingCourse != nil && existingCourse.ID != id {
			return nil, errors.New("course name already exists in this organization")
		}
	}

	// Update course
	updatedCourse, err := s.courseRepo.Update(ctx, id, req.Name, req.Description, req.EnrollmentOpen)
	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}

// DeleteCourse deletes a course
func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// Check if students are enrolled
	studentCount, err := s.courseRepo.CountStudentsByCourse(ctx, id)
	if err != nil {
		return err
	}

	if studentCount > 0 {
		return errors.New("cannot delete course with enrolled students")
	}

	// Delete course
	return s.courseRepo.Delete(ctx, id)
}

// AssignTeacherToCourse assigns a teacher to a course
func (s *CourseService) AssignTeacherToCourse(ctx context.Context, courseID, teacherID string) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// Check if teacher exists and has role 'teacher'
	teacher, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return err
	}

	if teacher == nil {
		return errors.New("teacher not found")
	}

	if teacher.Role != models.RoleTeacher {
		return errors.New("user is not a teacher")
	}

	// Check if teacher is in the same organization as the course
	if teacher.OrganizationID != course.OrganizationID {
		return errors.New("teacher and course must be in the same organization")
	}

	// Check if teacher is already assigned to the course
	isAssigned, err := s.courseRepo.IsTeacherAssigned(ctx, courseID, teacherID)
	if err != nil {
		return err
	}

	if isAssigned {
		return errors.New("teacher is already assigned to this course")
	}

	// Assign teacher to course
	return s.courseRepo.AssignTeacher(ctx, courseID, teacherID)
}

// RemoveTeacherFromCourse removes a teacher from a course
func (s *CourseService) RemoveTeacherFromCourse(ctx context.Context, courseID, teacherID string) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// Check if teacher is assigned to the course
	isAssigned, err := s.courseRepo.IsTeacherAssigned(ctx, courseID, teacherID)
	if err != nil {
		return err
	}

	if !isAssigned {
		return errors.New("teacher is not assigned to this course")
	}

	// Check if this is the last teacher (can't remove the last teacher)
	teacherCount, err := s.courseRepo.CountTeachersByCourse(ctx, courseID)
	if err != nil {
		return err
	}

	if teacherCount <= 1 {
		return errors.New("cannot remove the last teacher from a course")
	}

	// Remove teacher from course
	return s.courseRepo.RemoveTeacher(ctx, courseID, teacherID)
}

// GetCourseTeachers retrieves all teachers assigned to a course
func (s *CourseService) GetCourseTeachers(ctx context.Context, courseID string) ([]*models.User, error) {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Get teachers
	return s.courseRepo.FindTeachersByCourse(ctx, courseID)
}

// ToggleCourseEnrollment toggles a course's enrollment status
func (s *CourseService) ToggleCourseEnrollment(ctx context.Context, courseID string, enrollmentOpen bool) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// If trying to open enrollment, ensure at least one teacher is assigned
	if enrollmentOpen {
		hasTeachers, err := s.CourseHasTeachers(ctx, courseID)
		if err != nil {
			return err
		}

		if !hasTeachers {
			return errors.New("cannot open enrollment for a course without teachers")
		}
	}

	// Update enrollment status
	return s.courseRepo.UpdateEnrollmentStatus(ctx, courseID, enrollmentOpen)
}

// EnrollStudentInCourse enrolls a student in a course
func (s *CourseService) EnrollStudentInCourse(ctx context.Context, courseID, studentID string) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// Check if student exists and has role 'student'
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return errors.New("student not found")
	}

	if student.Role != models.RoleStudent {
		return errors.New("user is not a student")
	}

	// Check if student is in the same organization as the course
	if student.OrganizationID != course.OrganizationID {
		return errors.New("student and course must be in the same organization")
	}

	// Check if course is open for enrollment
	if !course.EnrollmentOpen {
		return errors.New("course is not open for enrollment")
	}

	// Check if student is already enrolled
	isEnrolled, err := s.courseRepo.IsStudentEnrolled(ctx, courseID, studentID)
	if err != nil {
		return err
	}

	if isEnrolled {
		return errors.New("student is already enrolled in this course")
	}

	// Check if course has teachers
	hasTeachers, err := s.CourseHasTeachers(ctx, courseID)
	if err != nil {
		return err
	}

	if !hasTeachers {
		return errors.New("cannot enroll in a course without teachers")
	}

	// Enroll student
	return s.courseRepo.EnrollStudent(ctx, courseID, studentID)
}

// UnenrollStudentFromCourse removes a student from a course
func (s *CourseService) UnenrollStudentFromCourse(ctx context.Context, courseID, studentID string) error {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("course not found")
	}

	// Check if student is enrolled
	isEnrolled, err := s.courseRepo.IsStudentEnrolled(ctx, courseID, studentID)
	if err != nil {
		return err
	}

	if !isEnrolled {
		return errors.New("student is not enrolled in this course")
	}

	// Unenroll student
	return s.courseRepo.UnenrollStudent(ctx, courseID, studentID)
}

// GetCourseStudents retrieves all students enrolled in a course
func (s *CourseService) GetCourseStudents(ctx context.Context, courseID string) ([]*models.User, error) {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Get students
	return s.courseRepo.FindStudentsByCourse(ctx, courseID)
}

// BulkEnrollStudents enrolls multiple students in a course
func (s *CourseService) BulkEnrollStudents(ctx context.Context, courseID string, studentIDs []string) (map[string]interface{}, error) {
	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Check if course is open for enrollment
	if !course.EnrollmentOpen {
		return nil, errors.New("course is not open for enrollment")
	}

	// Check if course has teachers
	hasTeachers, err := s.CourseHasTeachers(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if !hasTeachers {
		return nil, errors.New("cannot enroll in a course without teachers")
	}

	results := map[string]interface{}{
		"successful": make([]string, 0),
		"failed":     make([]map[string]string, 0),
	}

	for _, studentID := range studentIDs {
		err := s.EnrollStudentInCourse(ctx, courseID, studentID)
		if err != nil {
			results["failed"] = append(results["failed"].([]map[string]string), map[string]string{
				"student_id": studentID,
				"error":      err.Error(),
			})
		} else {
			results["successful"] = append(results["successful"].([]string), studentID)
		}
	}

	return results, nil
}

// GetTeacherCourses retrieves all courses assigned to a teacher
func (s *CourseService) GetTeacherCourses(ctx context.Context, teacherID string) ([]*models.CourseWithStudentCount, error) {
	// Check if teacher exists
	teacher, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	if teacher == nil {
		return nil, errors.New("teacher not found")
	}

	if teacher.Role != models.RoleTeacher {
		return nil, errors.New("user is not a teacher")
	}

	// Get courses
	courses, err := s.courseRepo.FindByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// For each course, get the student count
	var coursesWithCount []*models.CourseWithStudentCount
	for _, course := range courses {
		studentCount, err := s.courseRepo.CountStudentsByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		teacherCount, err := s.courseRepo.CountTeachersByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		coursesWithCount = append(coursesWithCount, &models.CourseWithStudentCount{
			Course:       course,
			StudentCount: studentCount,
			TeacherCount: teacherCount,
		})
	}

	return coursesWithCount, nil
}

// GetStudentCourses retrieves all courses a student is enrolled in
func (s *CourseService) GetStudentCourses(ctx context.Context, studentID string) ([]*models.CourseWithDetails, error) {
	// Check if student exists
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, errors.New("student not found")
	}

	if student.Role != models.RoleStudent {
		return nil, errors.New("user is not a student")
	}

	// Get courses
	courses, err := s.courseRepo.FindByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// For each course, get the teacher details
	var coursesWithDetails []*models.CourseWithDetails
	for _, course := range courses {
		teachers, err := s.courseRepo.FindTeachersByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		org, err := s.orgRepo.FindByID(ctx, course.OrganizationID)
		if err != nil {
			return nil, err
		}

		if org == nil {
			return nil, fmt.Errorf("organization not found for course: %s", course.ID)
		}

		studentCount, err := s.courseRepo.CountStudentsByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		coursesWithDetails = append(coursesWithDetails, &models.CourseWithDetails{
			Course:           course,
			OrganizationID:   course.OrganizationID,
			OrganizationName: org.Name,
			Teachers:         teachers,
			StudentCount:     studentCount,
		})
	}

	return coursesWithDetails, nil
}

// GetAvailableCourses retrieves all courses available for enrollment
func (s *CourseService) GetAvailableCourses(ctx context.Context, organizationID, studentID string) ([]*models.CourseWithDetails, error) {
	// Get all courses in the organization that are open for enrollment
	courses, err := s.courseRepo.FindAvailableCourses(ctx, organizationID, studentID)
	if err != nil {
		return nil, err
	}

	// For each course, get the teacher details
	var coursesWithDetails []*models.CourseWithDetails
	for _, course := range courses {
		teachers, err := s.courseRepo.FindTeachersByCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		// Only include courses with teachers
		if len(teachers) > 0 {
			org, err := s.orgRepo.FindByID(ctx, course.OrganizationID)
			if err != nil {
				return nil, err
			}

			if org == nil {
				return nil, fmt.Errorf("organization not found for course: %s", course.ID)
			}

			studentCount, err := s.courseRepo.CountStudentsByCourse(ctx, course.ID)
			if err != nil {
				return nil, err
			}

			coursesWithDetails = append(coursesWithDetails, &models.CourseWithDetails{
				Course:           course,
				OrganizationID:   course.OrganizationID,
				OrganizationName: org.Name,
				Teachers:         teachers,
				StudentCount:     studentCount,
			})
		}
	}

	return coursesWithDetails, nil
}

// IsTeacherAssignedToCourse checks if a teacher is assigned to a course
func (s *CourseService) IsTeacherAssignedToCourse(ctx context.Context, courseID, teacherID string) (bool, error) {
	return s.courseRepo.IsTeacherAssigned(ctx, courseID, teacherID)
}

// IsStudentEnrolledInCourse checks if a student is enrolled in a course
func (s *CourseService) IsStudentEnrolledInCourse(ctx context.Context, courseID, studentID string) (bool, error) {
	return s.courseRepo.IsStudentEnrolled(ctx, courseID, studentID)
}

// CourseHasTeachers checks if a course has at least one teacher assigned
func (s *CourseService) CourseHasTeachers(ctx context.Context, courseID string) (bool, error) {
	teacherCount, err := s.courseRepo.CountTeachersByCourse(ctx, courseID)
	if err != nil {
		return false, err
	}

	return teacherCount > 0, nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *CourseService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	return s.orgRepo.FindByID(ctx, id)
}
