package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"assessment-management-system/config"
	"assessment-management-system/db"
	customMiddleware "assessment-management-system/middleware"
	"assessment-management-system/migrations"
	"assessment-management-system/repositories"
	"assessment-management-system/routes"
	"assessment-management-system/services"
)

func init() {
	// Find the project root directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Load .env file if it exists
	if err := godotenv.Load(filepath.Join(basepath, ".env")); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}
}

func main() {
	// Initialize configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConn, err := db.Connect(appConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Run database migrations
	if err := migrations.RunMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed initial data
	orgRepo := repositories.NewOrganizationRepository(dbConn)
	userRepo := repositories.NewUserRepository(dbConn)
	seedService := services.NewSeedService(orgRepo, userRepo)

	if err := seedService.SeedInitialData(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(customMiddleware.LoggingMiddleware())

	// Initialize routes
	routes.SetupRoutes(e, dbConn, appConfig)

	// // Print all registered routes for debugging
	// for _, route := range e.Routes() {
	//         log.Printf("Route: %s %s -> %s\n", route.Method, route.Path, route.Name)
	// }

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}
