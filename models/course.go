package models

import (
	"time"
)

// Course represents a course within an organization
type Course struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	EnrollmentOpen bool      `json:"enrollment_open"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CourseWithTeachers represents a course with its assigned teachers
type CourseWithTeachers struct {
	Course   *Course `json:"course"`
	Teachers []*User `json:"teachers"`
}

// CourseWithStudentCount combines a course with student enrollment count
type CourseWithStudentCount struct {
	Course       *Course `json:"course"`
	StudentCount int     `json:"student_count"`
	TeacherCount int     `json:"teacher_count"`
}

// CourseWithDetails combines a course with organization and teacher details
type CourseWithDetails struct {
	Course           *Course `json:"course"`
	OrganizationID   string  `json:"organization_id"`
	OrganizationName string  `json:"organization_name"`
	Teachers         []*User `json:"teachers"`
	StudentCount     int     `json:"student_count"`
}

// CourseEnrollment represents a student's enrollment in a course
type CourseEnrollment struct {
	CourseID   string    `json:"course_id"`
	StudentID  string    `json:"student_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

// CourseTeacher represents a teacher assigned to a course
type CourseTeacher struct {
	CourseID  string    `json:"course_id"`
	TeacherID string    `json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCourseRequest represents the data needed to create a new course
type CreateCourseRequest struct {
	Name           string `json:"name" validate:"required,min=3,max=255"`
	Description    string `json:"description"`
	EnrollmentOpen bool   `json:"enrollment_open"`
	OrganizationID string `json:"organization_id" validate:"required"`
}

// UpdateCourseRequest represents the data needed to update a course
type UpdateCourseRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=3,max=255"`
	Description    *string `json:"description"`
	EnrollmentOpen *bool   `json:"enrollment_open"`
}

// AssignTeacherRequest represents the data needed to assign a teacher to a course
type AssignTeacherRequest struct {
	TeacherID string `json:"teacher_id" validate:"required"`
}

// EnrollStudentRequest represents the data needed to enroll a student in a course
type EnrollStudentRequest struct {
	StudentID string `json:"student_id" validate:"required"`
}

// BulkEnrollmentRequest represents the data needed for bulk enrollment
type BulkEnrollmentRequest struct {
	StudentIDs []string `json:"student_ids" validate:"required,min=1"`
}
