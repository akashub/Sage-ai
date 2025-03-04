// backend/go/cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sage-ai-v2/internal/api"
	"sage-ai-v2/internal/config"
	"sage-ai-v2/pkg/logger"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting server initialization...")
	ensureDirectoriesExist()

	// Load configuration
	cfg := loadConfig()

	fmt.Printf("Configuration loaded. Port: %d\n", cfg.Server.Port)

	// Set up router with all routes
	router := api.SetupRoutes()

	// Create the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
	}

	// Channel to listen for server errors
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		logger.InfoLogger.Printf("Server starting on port %d", cfg.Server.Port)
		logger.InfoLogger.Printf("LLM service URL: %s", cfg.LLM.ServiceURL)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for OS signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		logger.ErrorLogger.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		logger.InfoLogger.Println("Server is shutting down...")

		// Create a context with a timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeoutSeconds)*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := server.Shutdown(ctx); err != nil {
			logger.ErrorLogger.Printf("Server shutdown error: %v", err)
			os.Exit(1)
		}
	}
}

func ensureDirectoriesExist() {
    dirs := []string{
        "./data",
        "./data/uploads",
        "./logs",
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            logger.ErrorLogger.Fatalf("Failed to create directory %s: %v", dir, err)
        }
        
        // Log the absolute path of the directory
        absPath, err := filepath.Abs(dir)
        if err == nil {
            logger.InfoLogger.Printf("Directory ensured: %s (absolute: %s)", dir, absPath)
        }
    }
}

// loadConfig loads configuration with sensible defaults
func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		// If config can't be loaded, use defaults
		logger.InfoLogger.Printf("Using default configuration: %v", err)
		// No need to set defaults - Load() already returns default config with error
	}

	// Override with environment variables
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			cfg.Server.Port = p
		}
	}

	if llmURL := os.Getenv("LLM_SERVICE_URL"); llmURL != "" {
		cfg.LLM.ServiceURL = llmURL
	}

	return cfg
}