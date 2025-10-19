package app

import (
	"backend/controllers"
	"time"

	"github.com/gin-contrib/cors"
)

func mapUrls() {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, //almacena la configuracion de CORS por 12 horas
	}))

	// Public endpoints (no authentication required)
	router.POST("/users/register", controllers.Register)                   // Register new user
	router.POST("/users/verify-email", controllers.VerifyEmail)            // Verify email with code
	router.POST("/users/resend-code", controllers.ResendVerificationCode) // Resend verification code
	router.POST("/users/login", controllers.Login)                         // Login with credentials
	router.POST("/users/refresh-token", controllers.RefreshToken)         // Refresh access token

	// Protected endpoints (authentication required)
	router.GET("/users/:id", controllers.VerifyToken, controllers.GetUserByID) // Get user by ID

	// Admin endpoints (admin authentication required)
	router.GET("/users/admin", controllers.VerifyAdminToken)                       // Verify admin token
	router.POST("/users/promote-admin", controllers.VerifyAdminToken, controllers.PromoteToAdmin) // Promote user to admin (admin only)
}
