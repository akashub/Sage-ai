// middleware/cors.go
package middleware

import (
	"net/http"
	"sage-ai-v2/pkg/logger"
)

// CORSMiddleware applies CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLogger.Printf("CORS middleware processing request: %s %s", r.Method, r.URL.Path)
		
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600") // Cache preflight for 1 hour

		// Handle OPTIONS requests (preflight)
		if r.Method == "OPTIONS" {
			logger.InfoLogger.Printf("Handling OPTIONS request for: %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}