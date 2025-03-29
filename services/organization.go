package services

import (
	"context"
	"errors"

	"assessment-management-system/models"
	"assessment-management-system/repositories"
)

// OrganizationService handles organization-related business logic
type OrganizationService struct {
	orgRepo    *repositories.OrganizationRepository
	userRepo   *repositories.UserRepository
	courseRepo *repositories.CourseRepository
}

// NewOrganizationService creates a new OrganizationService
func NewOrganizationService(
	orgRepo *repositories.OrganizationRepository,
	userRepo *repositories.UserRepository,
	courseRepo *repositories.CourseRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:    orgRepo,
		userRepo:   userRepo,
		courseRepo: courseRepo,
	}
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, req models.CreateOrganizationRequest) (*models.Organization, error) {
	if req.Name == "" {
		return nil, errors.New("organization name is required")
	}

	// Create organization
	org, err := s.orgRepo.Create(ctx, req.Name, req.Slogan)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetAllOrganizations retrieves all organizations
func (s *OrganizationService) GetAllOrganizations(ctx context.Context) ([]*models.Organization, error) {
	return s.orgRepo.FindAll(ctx)
}

// GetOrganizationByID retrieves an organization by ID
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	return s.orgRepo.FindByID(ctx, id)
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id string, req models.UpdateOrganizationRequest) (*models.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Update the organization fields if provided
	if req.Name != nil {
		org.Name = *req.Name
	}

	if req.Slogan != nil {
		org.Slogan = *req.Slogan
	}

	// Save the updated organization
	updatedOrg, err := s.orgRepo.Update(ctx, id, org.Name, org.Slogan)
	if err != nil {
		return nil, err
	}

	return updatedOrg, nil
}

// DeleteOrganization deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	// Check if organization exists
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if org == nil {
		return errors.New("organization not found")
	}

	// Delete the organization (cascade will handle related records)
	return s.orgRepo.Delete(ctx, id)
}

// GetOrganizationStats retrieves statistics for an organization
func (s *OrganizationService) GetOrganizationStats(ctx context.Context, id string) (*models.OrganizationStats, error) {
	// Check if organization exists
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Get user counts
	userCount, err := s.userRepo.CountByOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	adminCount, err := s.userRepo.CountByOrganizationAndRole(ctx, id, string(models.RoleAdmin))
	if err != nil {
		return nil, err
	}

	teacherCount, err := s.userRepo.CountByOrganizationAndRole(ctx, id, string(models.RoleTeacher))
	if err != nil {
		return nil, err
	}

	studentCount, err := s.userRepo.CountByOrganizationAndRole(ctx, id, string(models.RoleStudent))
	if err != nil {
		return nil, err
	}

	// Get course count
	courseCount, err := s.courseRepo.CountByOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create the stats object
	stats := &models.OrganizationStats{
		Organization: org,
		UserCount:    userCount,
		AdminCount:   adminCount,
		TeacherCount: teacherCount,
		StudentCount: studentCount,
		CourseCount:  courseCount,
	}

	return stats, nil
}
