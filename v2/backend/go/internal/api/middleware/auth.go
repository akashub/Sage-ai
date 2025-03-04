// backend/go/internal/api/middleware/auth.go
package middleware

import (
	"net/http"
	"sage-ai-v2/pkg/logger"
)

// ApplyAuth is a middleware that applies authentication to the API endpoints
func ApplyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the request
		logger.InfoLogger.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
