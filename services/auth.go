package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"assessment-management-system/models"
	"assessment-management-system/repositories"
	"assessment-management-system/utils"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repositories.UserRepository, refreshTokenRepo *repositories.RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// Authenticate authenticates a user using email and password
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	// First try to find user by email (might return multiple users across organizations)
	users, err := s.userRepo.FindAllByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("invalid credentials")
	}

	// If there's only one user with this email, verify password
	if len(users) == 1 {
		user := users[0]
		if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
			return nil, errors.New("invalid credentials")
		}
		return user, nil
	}

	// If multiple users with the same email exist (in different organizations),
	// try to verify password for each of them
	for _, user := range users {
		if err := utils.VerifyPassword(user.PasswordHash, password); err == nil {
			// If password matches, return this user
			return user, nil
		}
	}

	return nil, errors.New("invalid credentials")
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := utils.VerifyPassword(user.PasswordHash, currentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	err = s.userRepo.UpdatePassword(ctx, userID, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

// CreateTokenPair creates both access and refresh tokens for a user
func (s *AuthService) CreateTokenPair(ctx context.Context, user *models.User, jwtSecret string,
	jwtExpiration, refreshExpiration time.Duration) (*models.TokenResponse, error) {

	// Generate access token
	accessToken, err := utils.GenerateToken(user, jwtSecret, jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := utils.CreateRefreshTokenModel(user.ID, refreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token to database
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Return token response
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int64(jwtExpiration.Seconds()),
	}

	return tokenResponse, nil
}

// RefreshAccessToken validates a refresh token and generates a new access token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenStr string,
	jwtSecret string, jwtExpiration time.Duration) (*models.TokenResponse, error) {

	// Get refresh token from database
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is revoked
	if refreshToken.Revoked {
		return nil, errors.New("refresh token has been revoked")
	}

	// Check if token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	// Get user for this refresh token
	user, err := s.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(user, jwtSecret, jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Return token response
	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token, // Return the same refresh token
		ExpiresIn:    int64(jwtExpiration.Seconds()),
	}

	return tokenResponse, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshTokenStr string) error {
	return s.refreshTokenRepo.Revoke(ctx, refreshTokenStr)
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (s *AuthService) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
}
