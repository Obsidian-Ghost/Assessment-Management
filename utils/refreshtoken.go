package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"assessment-management-system/models"
)

// GenerateRefreshToken generates a cryptographically secure refresh token
func GenerateRefreshToken() (string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to base64 string
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateRefreshTokenModel creates a refresh token model with the given user ID and expiration
func CreateRefreshTokenModel(userID string, expires time.Duration) (*models.RefreshToken, error) {
	token, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(expires),
		Revoked:   false,
	}

	return refreshToken, nil
}
