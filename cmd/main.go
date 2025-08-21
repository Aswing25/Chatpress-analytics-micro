package main

import (
	v1 "chatpress-analytics/api/v1"
	"chatpress-analytics/config"
	"chatpress-analytics/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize repositories
	analyticsRepo := db.NewAnalyticsRepository(database)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	handler := v1.NewHandler(analyticsRepo, cfg.JWTSecret)

	// Setup analytics routes
	api := router.Group("/analytics")
	{
		api.GET("/api-usage-cost", handler.GetAPIUsageCost)
		api.GET("/monthly-usage", handler.GetMonthlyUsage)
		api.GET("/overall-status", handler.GetOverallStatus)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "chatpress-analytics",
			"status":  "running",
			"message": "Welcome to ChatPress Analytics Service 🚀",
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Prevent favicon.ico 404 spam
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	log.Printf("Analytics service starting on port %s", cfg.Port)
	log.Fatal(router.Run("127.0.0.1:" + cfg.Port))
}
