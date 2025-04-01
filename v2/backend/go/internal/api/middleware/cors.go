// // backend/go/internal/api/middleware/cors.go
// package middleware

// import (
// 	"net/http"
// 	"sage-ai-v2/pkg/logger"
// )

// // CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers
// func CORSMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Set CORS headers for all responses
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		// Handle OPTIONS requests (preflight)
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		// Call the next handler
// 		next.ServeHTTP(w, r)
// 	})
// }

// // LoggingMiddleware logs all incoming requests
// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.InfoLogger.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
// 		next.ServeHTTP(w, r)
// 	})
// }

// backend/go/internal/api/middleware/cors.go
package middleware

import (
	"net/http"
	"sage-ai-v2/pkg/logger"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers
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

// LoggingMiddleware logs all incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLogger.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}