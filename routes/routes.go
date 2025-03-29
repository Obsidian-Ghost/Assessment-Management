package routes

import (
	"github.com/labstack/echo/v4"

	"assessment-management-system/config"
	"assessment-management-system/db"
	"assessment-management-system/handlers"
	"assessment-management-system/handlers/admin"
	"assessment-management-system/handlers/student"
	"assessment-management-system/handlers/teacher"
	customMiddleware "assessment-management-system/middleware"
	"assessment-management-system/repositories"
	"assessment-management-system/services"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo, db *db.DB, cfg *config.AppConfig) {
	// Create repositories
	orgRepo := repositories.NewOrganizationRepository(db)
	userRepo := repositories.NewUserRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	assessmentRepo := repositories.NewAssessmentRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// Create services
	authService := services.NewAuthService(userRepo, refreshTokenRepo)
	orgService := services.NewOrganizationService(orgRepo, userRepo, courseRepo)
	userService := services.NewUserService(userRepo, orgRepo, courseRepo, assessmentRepo)
	courseService := services.NewCourseService(courseRepo, userRepo, orgRepo)
	assessmentService := services.NewAssessmentService(assessmentRepo, courseRepo, userRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, userService, cfg.JWT.Secret, cfg.JWT.Expiration)

	// Admin handlers
	adminOrgHandler := admin.NewOrganizationHandler(orgService)
	adminUserHandler := admin.NewUserHandler(userService)
	adminCourseHandler := admin.NewCourseHandler(courseService)
	adminAssessmentHandler := admin.NewAssessmentHandler(assessmentService, courseService)

	// Teacher handlers
	teacherCourseHandler := teacher.NewCourseHandler(courseService)
	teacherAssessmentHandler := teacher.NewAssessmentHandler(assessmentService, courseService)

	// Student handlers
	studentCourseHandler := student.NewCourseHandler(courseService)
	studentAssessmentHandler := student.NewAssessmentHandler(assessmentService, courseService)

	// Auth middleware
	authMiddleware := customMiddleware.AuthMiddleware(cfg.JWT.Secret)

	// Role-based middleware
	adminOnly := customMiddleware.AdminOnly()
	teacherOnly := customMiddleware.TeacherOnly()
	studentOnly := customMiddleware.StudentOnly()

	// Public routes
	api := e.Group("/api")
	api.POST("/auth/login", authHandler.HandleLogin)
	api.POST("/auth/token/refresh", authHandler.HandleRefreshToken)

	// Protected routes
	apiAuth := api.Group("", authMiddleware)

	// User profile
	apiAuth.GET("/auth/me", authHandler.HandleGetMe)
	apiAuth.POST("/auth/change-password", authHandler.HandleChangePassword)
	apiAuth.POST("/auth/token/revoke", authHandler.HandleRevokeToken)
	apiAuth.POST("/auth/token/revoke-all", authHandler.HandleRevokeAllTokens)

	// Admin routes
	adminRoutes := apiAuth.Group("/admin", adminOnly)

	// Organization management
	adminRoutes.POST("/organizations", adminOrgHandler.HandleCreateOrganization)
	adminRoutes.GET("/organizations", adminOrgHandler.HandleGetAllOrganizations)
	adminRoutes.GET("/organizations/:id", adminOrgHandler.HandleGetOrganizationByID)
	adminRoutes.PUT("/organizations/:id", adminOrgHandler.HandleUpdateOrganization)
	adminRoutes.DELETE("/organizations/:id", adminOrgHandler.HandleDeleteOrganization)
	adminRoutes.GET("/organizations/:id/stats", adminOrgHandler.HandleGetOrganizationStats)

	// User management
	adminRoutes.POST("/users", adminUserHandler.HandleCreateUser)
	adminRoutes.GET("/users", adminUserHandler.HandleGetAllUsers)
	adminRoutes.GET("/users/:id", adminUserHandler.HandleGetUserByID)
	adminRoutes.PUT("/users/:id", adminUserHandler.HandleUpdateUser)
	adminRoutes.DELETE("/users/:id", adminUserHandler.HandleDeleteUser)
	adminRoutes.POST("/users/bulk", adminUserHandler.HandleBulkUploadUsers)
	adminRoutes.GET("/users/teachers/:id/stats", adminUserHandler.HandleGetTeacherStats)
	adminRoutes.GET("/users/students/:id/stats", adminUserHandler.HandleGetStudentStats)

	// Course management
	adminRoutes.POST("/courses", adminCourseHandler.HandleCreateCourse)
	adminRoutes.GET("/courses", adminCourseHandler.HandleGetAllCourses)
	adminRoutes.GET("/courses/:id", adminCourseHandler.HandleGetCourseByID)
	adminRoutes.PUT("/courses/:id", adminCourseHandler.HandleUpdateCourse)
	adminRoutes.DELETE("/courses/:id", adminCourseHandler.HandleDeleteCourse)
	adminRoutes.POST("/courses/:id/teachers", adminCourseHandler.HandleAssignTeacher)
	adminRoutes.DELETE("/courses/:id/teachers/:teacherId", adminCourseHandler.HandleRemoveTeacher)
	adminRoutes.GET("/courses/:id/teachers", adminCourseHandler.HandleGetCourseTeachers)
	adminRoutes.PUT("/courses/:id/enrollment", adminCourseHandler.HandleToggleEnrollment)
	adminRoutes.POST("/courses/:id/students", adminCourseHandler.HandleManageStudentEnrollment)
	adminRoutes.POST("/courses/:id/students/bulk", adminCourseHandler.HandleBulkEnrollStudents)
	adminRoutes.GET("/courses/:id/students", adminCourseHandler.HandleGetCourseStudents)

	// Assessment management (read-only for admin)
	adminRoutes.GET("/assessments", adminAssessmentHandler.HandleGetAllAssessments)
	adminRoutes.GET("/assessments/:id", adminAssessmentHandler.HandleGetAssessmentByID)
	adminRoutes.GET("/assessments/:id/submissions", adminAssessmentHandler.HandleGetAssessmentSubmissions)
	adminRoutes.GET("/submissions/:submissionId/grade", adminAssessmentHandler.HandleGetSubmissionGrades)

	// Teacher routes
	teacherRoutes := apiAuth.Group("/teacher", teacherOnly)

	// Course management for teachers
	teacherRoutes.GET("/courses", teacherCourseHandler.HandleGetAssignedCourses)
	teacherRoutes.GET("/courses/:id", teacherCourseHandler.HandleGetCourseByID)
	teacherRoutes.GET("/courses/:id/students", teacherCourseHandler.HandleGetCourseStudents)
	teacherRoutes.GET("/organization", teacherCourseHandler.HandleGetOrganizationDetails)

	// Assessment management for teachers
	teacherRoutes.POST("/assessments", teacherAssessmentHandler.HandleCreateAssessment)
	teacherRoutes.GET("/courses/:courseId/assessments", teacherAssessmentHandler.HandleGetAssessments)
	teacherRoutes.GET("/assessments/:id", teacherAssessmentHandler.HandleGetAssessmentByID)
	teacherRoutes.PUT("/assessments/:id", teacherAssessmentHandler.HandleUpdateAssessment)
	teacherRoutes.DELETE("/assessments/:id", teacherAssessmentHandler.HandleDeleteAssessment)
	teacherRoutes.GET("/assessments/:id/submissions", teacherAssessmentHandler.HandleGetSubmissions)
	teacherRoutes.POST("/submissions/:submissionId/grade", teacherAssessmentHandler.HandleGradeSubmission)

	// Student routes
	studentRoutes := apiAuth.Group("/student", studentOnly)

	// Course management for students
	studentRoutes.GET("/courses", studentCourseHandler.HandleGetEnrolledCourses)
	studentRoutes.GET("/courses/:id", studentCourseHandler.HandleGetCourseByID)
	studentRoutes.GET("/courses/available", studentCourseHandler.HandleGetAvailableCourses)
	studentRoutes.POST("/courses/:id/enroll", studentCourseHandler.HandleEnrollInCourse)
	studentRoutes.GET("/organization", studentCourseHandler.HandleGetOrganizationDetails)

	// Assessment management for students
	studentRoutes.GET("/courses/:courseId/assessments", studentAssessmentHandler.HandleGetCourseAssessments)
	studentRoutes.GET("/assessments/:id", studentAssessmentHandler.HandleGetAssessmentByID)
	studentRoutes.POST("/assessments/:id/submit", studentAssessmentHandler.HandleSubmitAssessment)
	studentRoutes.GET("/assessments/:id/submission", studentAssessmentHandler.HandleViewSubmission)
	studentRoutes.GET("/assessments/:id/grade", studentAssessmentHandler.HandleViewGrade)
}
