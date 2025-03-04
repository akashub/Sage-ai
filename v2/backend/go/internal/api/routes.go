// backend/go/internal/api/routes.go
// backend/go/internal/api/routes.go
package api

import (
	"net/http"
	"sage-ai-v2/internal/api/handlers"
	"sage-ai-v2/internal/api/middleware"
	"sage-ai-v2/pkg/logger"
)

// SetupRoutes configures all API routes for the application
func SetupRoutes() http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register API routes
	mux.HandleFunc("/api/upload", handlers.UploadFileHandler)
	mux.HandleFunc("/api/query", handlers.QueryHandler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware to all routes
	handler := middleware.ApplyAuth(mux)

	logger.InfoLogger.Printf("API routes configured")
	return handler
}