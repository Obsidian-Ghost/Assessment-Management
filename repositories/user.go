package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"assessment-management-system/db"
	"assessment-management-system/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *db.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, organizationID, email, passwordHash, firstName, lastName, role string) (*models.User, error) {
	var user models.User
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO users (organization_id, email, password_hash, first_name, last_name, role) 
                VALUES ($1, $2, $3, $4, $5, $6) 
                RETURNING id, organization_id, email, first_name, last_name, role, created_at, updated_at`,
		organizationID, email, passwordHash, firstName, lastName, role).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, organization_id, email, password_hash, first_name, last_name, role, created_at, updated_at 
                FROM users 
                WHERE id = $1`,
		id).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email (across all organizations)
// This is used to enforce global email uniqueness
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, organization_id, email, password_hash, first_name, last_name, role, created_at, updated_at 
                FROM users 
                WHERE email = $1
                LIMIT 1`,
		email).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAllByEmail retrieves all users with a given email (across all organizations)
func (r *UserRepository) FindAllByEmail(ctx context.Context, email string) ([]*models.User, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, organization_id, email, password_hash, first_name, last_name, role, created_at, updated_at 
                FROM users 
                WHERE email = $1`,
		email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindByEmailAndOrganization retrieves a user by email in a specific organization
func (r *UserRepository) FindByEmailAndOrganization(ctx context.Context, email, organizationID string) (*models.User, error) {
	var user models.User
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, organization_id, email, password_hash, first_name, last_name, role, created_at, updated_at 
                FROM users 
                WHERE email = $1 AND organization_id = $2`,
		email, organizationID).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByOrganization retrieves all users in an organization with optional filtering
func (r *UserRepository) FindByOrganization(ctx context.Context, organizationID, role, search string) ([]*models.User, error) {
	var query string
	var args []interface{}

	// Base query
	query = `SELECT id, organization_id, email, first_name, last_name, role, created_at, updated_at 
                        FROM users 
                        WHERE organization_id = $1`
	args = append(args, organizationID)

	// Add role filter if provided
	if role != "" {
		query += ` AND role = $2`
		args = append(args, role)
	}

	// Add search filter if provided
	if search != "" {
		if role != "" {
			query += ` AND (email ILIKE $3 OR first_name ILIKE $3 OR last_name ILIKE $3)`
			args = append(args, "%"+search+"%")
		} else {
			query += ` AND (email ILIKE $2 OR first_name ILIKE $2 OR last_name ILIKE $2)`
			args = append(args, "%"+search+"%")
		}
	}

	// Order by role and name
	query += ` ORDER BY role, first_name, last_name`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, id string, email *string, passwordHash string, firstName *string, lastName *string, role *models.UserRole, organizationID *string) (*models.User, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get current user
	var user models.User
	err = tx.QueryRow(ctx,
		`SELECT id, organization_id, email, password_hash, first_name, last_name, role, created_at, updated_at 
                FROM users 
                WHERE id = $1`,
		id).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields that are provided
	if email != nil {
		user.Email = *email
	}
	if passwordHash != "" {
		user.PasswordHash = passwordHash
	}
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if role != nil {
		user.Role = *role
	}
	if organizationID != nil {
		user.OrganizationID = *organizationID
	}

	// Update in database
	err = tx.QueryRow(ctx,
		`UPDATE users 
                SET email = $2, password_hash = $3, first_name = $4, last_name = $5, role = $6, organization_id = $7, updated_at = $8
                WHERE id = $1 
                RETURNING id, organization_id, email, first_name, last_name, role, created_at, updated_at`,
		id, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.OrganizationID, time.Now()).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`UPDATE users 
                SET password_hash = $2, updated_at = $3
                WHERE id = $1`,
		id, passwordHash, time.Now())

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM users WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// CountByOrganization counts users in an organization
func (r *UserRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE organization_id = $1`,
		organizationID).Scan(&count)
	return count, err
}

// CountByOrganizationAndRole counts users in an organization with a specific role
func (r *UserRepository) CountByOrganizationAndRole(ctx context.Context, organizationID, role string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE organization_id = $1 AND role = $2`,
		organizationID, role).Scan(&count)
	return count, err
}

// ExecuteInTransaction executes a function within a transaction
func (r *UserRepository) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return r.db.ExecuteTransaction(ctx, fn)
}
