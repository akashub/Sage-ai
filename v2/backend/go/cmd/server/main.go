// backend/go/cmd/server/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sage-ai-v2/internal/api"
	"sage-ai-v2/internal/config"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/pkg/logger"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	fmt.Println("Starting server initialization...")

	// Setup directories
	ensureDirectoriesExist()

	// Load configuration
	cfg := loadConfig()

	// Initialize database
	db, err := setupDatabase()
	if err != nil {
		logger.ErrorLogger.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Initialize knowledge components
	knowledgeManager, err := setupKnowledge(cfg)
	if err != nil {
		logger.ErrorLogger.Printf("Knowledge system initialization failed: %v", err)
		logger.InfoLogger.Println("Continuing without knowledge graph support...")
	}

	// Initialize LLM bridge
	bridge := llm.CreateBridge(cfg.LLM.ServiceURL)

	// Initialize orchestrator with knowledge manager
	orch := orchestrator.CreateOrchestrator(bridge, knowledgeManager)

	// Set up router with all routes
	router := api.SetupRoutes(db, orch)

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

// ensureDirectoriesExist creates necessary directories if they don't exist
func ensureDirectoriesExist() {
	dirs := []string{
		"./data",
		"./data/uploads",
		"./data/training",
		"./data/knowledge",
		"./logs",
		"./data/db", // For SQLite database
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

// setupDatabase initializes the database connection and creates tables if needed
func setupDatabase() (*sql.DB, error) {
	dbPath := "./data/db/sage.db"

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	logger.InfoLogger.Printf("Database initialized successfully at %s", dbPath)
	return db, nil
}

// setupKnowledge initializes knowledge management components
func setupKnowledge(cfg *config.Config) (*knowledge.KnowledgeManager, error) {
	// Path for vector database persistence
	persistPath := "./data/knowledge/vector_store.json"
	trainingDir := "./data/training"

	// Create in-memory vector database
	vectorDB, err := knowledge.CreateMemoryVectorDB(persistPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector database: %w", err)
	}

	// Create knowledge manager
	knowledgeManager := knowledge.CreateKnowledgeManager(vectorDB, trainingDir)
	logger.InfoLogger.Printf("Knowledge management system initialized successfully")

	return knowledgeManager, nil
}

// createTables creates necessary database tables if they don't exist
func createTables(db *sql.DB) error {
	// Users table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT,
		name TEXT,
		created_at TIMESTAMP NOT NULL,
		last_login_at TIMESTAMP,
		provider_type TEXT,
		provider_id TEXT,
		refresh_token TEXT,
		profile_pic_url TEXT
	)`)

	if err != nil {
		return err
	}

	// Create index on email
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	if err != nil {
		return err
	}

	// Create index on provider_id
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider_type, provider_id)`)
	if err != nil {
		return err
	}

	return nil
}
