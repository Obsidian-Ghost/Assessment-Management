package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// DatabaseConfig holds database connection details
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	URL      string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	Database    DatabaseConfig
	JWT         JWTConfig
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*AppConfig, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// It's okay if .env is missing in production, so we log a warning instead of returning an error
		fmt.Println("Warning: No .env file found, using system environment variables")
	}

	// Environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Database configuration
	dbHost := os.Getenv("PGHOST")
	dbPortStr := os.Getenv("PGPORT")
	dbUser := os.Getenv("PGUSER")
	dbPassword := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")
	dbURL := os.Getenv("DATABASE_URL")

	// Convert dbPort to int
	dbPort := 5432 // Default PostgreSQL port
	if dbPortStr != "" {
		dbPort, err = strconv.Atoi(dbPortStr)
		if err != nil {
			return nil, errors.New("invalid database port: " + dbPortStr)
		}
	}

	// JWT configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // Default JWT secret (only for development)
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	jwtExpiration := 24 * time.Hour // Default expiration time
	if jwtExpirationStr != "" {
		jwtExpirationInt, err := strconv.Atoi(jwtExpirationStr)
		if err != nil {
			return nil, errors.New("invalid JWT expiration: " + jwtExpirationStr)
		}
		jwtExpiration = time.Duration(jwtExpirationInt) * time.Hour
	}

	return &AppConfig{
		Environment: env,
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
			URL:      dbURL,
		},
		JWT: JWTConfig{
			Secret:     jwtSecret,
			Expiration: jwtExpiration,
		},
	}, nil
}
