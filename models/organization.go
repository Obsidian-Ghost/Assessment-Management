package models

import (
	"time"
)

// Organization represents an educational organization in the system
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slogan    string    `json:"slogan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrganizationStats represents statistics for an organization
type OrganizationStats struct {
	Organization *Organization `json:"organization"`
	UserCount    int           `json:"user_count"`
	AdminCount   int           `json:"admin_count"`
	TeacherCount int           `json:"teacher_count"`
	StudentCount int           `json:"student_count"`
	CourseCount  int           `json:"course_count"`
}

// CreateOrganizationRequest represents the data needed to create a new organization
type CreateOrganizationRequest struct {
	Name   string `json:"name" validate:"required,min=3,max=255"`
	Slogan string `json:"slogan" validate:"max=1000"`
}

// UpdateOrganizationRequest represents the data needed to update an organization
type UpdateOrganizationRequest struct {
	Name   *string `json:"name" validate:"omitempty,min=3,max=255"`
	Slogan *string `json:"slogan" validate:"omitempty,max=1000"`
}
