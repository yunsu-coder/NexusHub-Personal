package main

import (
	"log"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/logger"
	"nexushub-personal/internal/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize configuration with validation
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	log.Println("Configuration loaded and validated successfully")

	// Initialize logger
	logDir := "./logs"
	if err := logger.Init(logDir, logger.INFO); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("NexusHub Personal starting...")
	logger.Info("Server mode: %s", config.AppConfig.Server.GinMode)

	// Initialize database
	if err := database.Init(); err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	logger.Info("Database initialized successfully")

	// Setup router
	r := router.SetupRouter()

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		logger.Close()
		os.Exit(0)
	}()

	// Start server
	addr := ":" + config.AppConfig.Server.Port
	logger.Info("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
