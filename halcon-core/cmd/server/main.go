package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nietzshn/halcon-core/internal/config"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/handlers"
	custommw "github.com/nietzshn/halcon-core/internal/middleware"
	"github.com/nietzshn/halcon-core/internal/models"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data
	if err := database.Seed(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	// Create uploads directory
	if err := os.MkdirAll(config.AppConfig.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.AppConfig.CORSAllowedOrigins},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Static files for uploads
	e.Static("/uploads", config.AppConfig.UploadDir)

	// Public routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	e.POST("/api/auth/login", handlers.Login)
	e.GET("/api/track", handlers.TrackOrder)
	e.POST("/api/track", handlers.TrackOrder)

	// Protected routes
	api := e.Group("/api")
	api.Use(custommw.AuthMiddleware())

	// Auth routes
	api.GET("/auth/me", handlers.GetCurrentUser)

	// User management routes (Admin only)
	users := api.Group("/users")
	users.Use(custommw.RoleMiddleware(models.RoleAdmin))
	users.GET("", handlers.GetUsers)
	users.GET("/:id", handlers.GetUser)
	users.POST("", handlers.CreateUser)
	users.PUT("/:id", handlers.UpdateUser)
	users.DELETE("/:id", handlers.DeleteUser)

	// Order routes
	orders := api.Group("/orders")

	// All authenticated users can view orders (with role-based filtering in handler)
	orders.GET("", handlers.GetOrders)
	orders.GET("/:id", handlers.GetOrder)

	// Sales can create orders
	orders.POST("", handlers.CreateOrder, custommw.RoleMiddleware(models.RoleSales))

	// Warehouse and Route can update orders
	orders.PUT("/:id", handlers.UpdateOrder, custommw.RoleMiddleware(
		models.RoleWarehouse,
		models.RoleRoute,
		models.RoleSales,
	))

	// Soft delete and restore (Admin and Sales)
	orders.DELETE("/:id", handlers.SoftDeleteOrder, custommw.RoleMiddleware(
		models.RoleAdmin,
		models.RoleSales,
	))
	orders.POST("/:id/restore", handlers.RestoreOrder, custommw.RoleMiddleware(
		models.RoleAdmin,
		models.RoleSales,
	))

	// Upload evidence (Route only)
	orders.POST("/:id/evidence", handlers.UploadEvidence, custommw.RoleMiddleware(models.RoleRoute))

	// Start server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
