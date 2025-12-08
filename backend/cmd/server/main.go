package main

import (
	"log"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/router"
)

func main() {
	// Initialize configuration
	config.Init()
	log.Println("Configuration loaded")

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup router
	r := router.SetupRouter()

	// Start server
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
