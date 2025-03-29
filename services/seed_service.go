package services

import (
	"context"
	"log"

	"assessment-management-system/models"
	"assessment-management-system/repositories"
	"assessment-management-system/utils"
)

// SeedService handles seeding initial data in the database
type SeedService struct {
	orgRepo  *repositories.OrganizationRepository
	userRepo *repositories.UserRepository
}

// NewSeedService creates a new SeedService
func NewSeedService(
	orgRepo *repositories.OrganizationRepository,
	userRepo *repositories.UserRepository,
) *SeedService {
	return &SeedService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// SeedInitialData creates the default organization and admin user
func (s *SeedService) SeedInitialData(ctx context.Context) error {
	// Check if any organizations exist
	orgs, err := s.orgRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	// If no organizations exist, create a default one
	var defaultOrg *models.Organization
	if len(orgs) == 0 {
		defaultOrg, err = s.orgRepo.Create(ctx, "Default Organization", "A place for learning and growth")
		if err != nil {
			return err
		}
		log.Println("Created default organization")
	} else {
		defaultOrg = orgs[0]
	}

	// Check if any users exist for this organization
	count, err := s.userRepo.CountByOrganization(ctx, defaultOrg.ID)
	if err != nil {
		return err
	}

	// If no users exist, create a default admin
	if count == 0 {
		// Generate password hash
		passwordHash, err := utils.HashPassword("admin123")
		if err != nil {
			return err
		}

		_, err = s.userRepo.Create(
			ctx,
			defaultOrg.ID,
			"admin@example.com",
			passwordHash,
			"System",
			"Administrator",
			string(models.RoleAdmin),
		)
		if err != nil {
			return err
		}
		log.Println("Created default admin user (admin@example.com / admin123)")
	}

	return nil
}
