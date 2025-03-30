package models

import (
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

// User represents a user of the system
type User struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Role           UserRole  `json:"role"`
	PasswordHash   string    `json:"-"` // Omitted from JSON responses
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UserWithOrgName includes organization name with user details
type UserWithOrgName struct {
	User             *User  `json:"user"`
	OrganizationName string `json:"organization_name"`
}

// TeacherStats represents statistics for a teacher
type TeacherStats struct {
	User               *User `json:"user"`
	AssignedCourses    int   `json:"assigned_courses"`
	CreatedAssessments int   `json:"created_assessments"`
	PendingGrading     int   `json:"pending_grading"`
}

// StudentStats represents statistics for a student
type StudentStats struct {
	User                 *User   `json:"user"`
	EnrolledCourses      int     `json:"enrolled_courses"`
	CompletedAssessments int     `json:"completed_assessments"`
	PendingAssessments   int     `json:"pending_assessments"`
	AverageGrade         float64 `json:"average_grade"`
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Email          string   `json:"email" validate:"required,email"`
	Password       string   `json:"password" validate:"required,min=8"`
	FirstName      string   `json:"first_name" validate:"required"`
	LastName       string   `json:"last_name" validate:"required"`
	Role           UserRole `json:"role" validate:"required,oneof=admin teacher student"`
	OrganizationID string   `json:"organization_id"` // Optional, will use admin's organization if not provided
}

// UpdateUserRequest represents the data needed to update a user
type UpdateUserRequest struct {
	Email     *string   `json:"email" validate:"omitempty,email"`
	Password  *string   `json:"password" validate:"omitempty,min=8"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Role      *UserRole `json:"role" validate:"omitempty,oneof=admin teacher student"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// BulkUserUploadRequest represents the data format for bulk user upload
type BulkUserUploadRequest struct {
	OrganizationID string              `json:"organization_id" validate:"required"`
	Users          []CreateUserRequest `json:"users" validate:"required,min=1,dive"`
}
