// internal/api/middleware/logging.go
package middleware

import (
    "net/http"
    "time"
    "sage-ai-v2/pkg/logger"
)

// LoggingMiddleware logs information about each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log request details
        logger.InfoLogger.Printf("Request started: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Log completion time
        duration := time.Since(start)
        logger.InfoLogger.Printf("Request completed: %s %s in %v", r.Method, r.URL.Path, duration)
    })
}