package main

import (
	"competition-app/config"
	"competition-app/models"
	"competition-app/routes"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the database connection
	if err := models.InitDB(cfg); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer models.DB.Close()

	// Initialize Redis connection
	if err := models.InitRedis(cfg); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		// Redis isn't critical, so we don't exit if it fails
	}
	defer models.CloseRedis()

	// Initialize router
	router := routes.SetupRouter()

	// Start the server
	port := cfg.ServerPort
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
