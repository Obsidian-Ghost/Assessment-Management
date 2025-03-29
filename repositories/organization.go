package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"assessment-management-system/db"
	"assessment-management-system/models"
)

// OrganizationRepository handles database operations for organizations
type OrganizationRepository struct {
	db *db.DB
}

// NewOrganizationRepository creates a new OrganizationRepository
func NewOrganizationRepository(db *db.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db: db,
	}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, name, slogan string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO organizations (name, slogan) 
                VALUES ($1, $2) 
                RETURNING id, name, slogan, created_at, updated_at`,
		name, slogan).Scan(&org.ID, &org.Name, &org.Slogan, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAll retrieves all organizations
func (r *OrganizationRepository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, name, slogan, created_at, updated_at 
                FROM organizations 
                ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Slogan, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgs, nil
}

// FindByID retrieves an organization by ID
func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, slogan, created_at, updated_at 
                FROM organizations 
                WHERE id = $1`,
		id).Scan(&org.ID, &org.Name, &org.Slogan, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &org, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, id string, name string, slogan string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.Pool.QueryRow(ctx,
		`UPDATE organizations 
                SET name = $2, slogan = $3, updated_at = $4
                WHERE id = $1 
                RETURNING id, name, slogan, created_at, updated_at`,
		id, name, slogan, time.Now()).Scan(&org.ID, &org.Name, &org.Slogan, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &org, nil
}

// Delete deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM organizations WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("organization not found")
	}
	return nil
}

// ExecuteInTransaction executes a function within a transaction
func (r *OrganizationRepository) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return r.db.ExecuteTransaction(ctx, fn)
}
