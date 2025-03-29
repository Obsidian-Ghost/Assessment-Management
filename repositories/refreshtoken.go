package repositories

import (
	"context"
	"fmt"
	"time"

	"assessment-management-system/db"
	"assessment-management-system/models"
)

// RefreshTokenRepository handles database operations for refresh tokens
type RefreshTokenRepository struct {
	db *db.DB
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository
func NewRefreshTokenRepository(db *db.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// Create creates a new refresh token in the database
func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
	).Scan(&refreshToken.ID, &refreshToken.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetByToken retrieves a refresh token by its token value
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked, revoked_at
		FROM refresh_tokens
		WHERE token = $1
	`

	refreshToken := &models.RefreshToken{}
	var revokedAt *time.Time

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		token,
	).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.Revoked,
		&revokedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if revokedAt != nil {
		refreshToken.RevokedAt = *revokedAt
	}

	return refreshToken, nil
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = NOW()
		WHERE token = $1
	`

	_, err := r.db.Pool.Exec(
		ctx,
		query,
		token,
	)

	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a specific user
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = NOW()
		WHERE user_id = $1 AND revoked = false
	`

	_, err := r.db.Pool.Exec(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to revoke all user refresh tokens: %w", err)
	}

	return nil
}

// CleanExpired deletes all expired refresh tokens
func (r *RefreshTokenRepository) CleanExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
	`

	_, err := r.db.Pool.Exec(ctx, query)

	if err != nil {
		return fmt.Errorf("failed to clean expired refresh tokens: %w", err)
	}

	return nil
}
