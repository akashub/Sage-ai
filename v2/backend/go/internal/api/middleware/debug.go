// middleware/debug.go
package middleware

import (
    "net/http"
    "strings"
    "sage-ai-v2/pkg/logger"
)

// DebugMiddleware logs detailed request and response information
func DebugMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only debug auth-related endpoints
        if !strings.Contains(r.URL.Path, "/auth/") {
            next.ServeHTTP(w, r)
            return
        }
        
        // Log the request method, URL, and headers
        logger.InfoLogger.Printf("Auth request: %s %s", r.Method, r.URL.Path)
        logger.InfoLogger.Printf("Auth request headers: %v", r.Header)
        
        // Create a response wrapper to capture the status
        rw := newResponseWriter(w)
        
        // Call the next handler
        next.ServeHTTP(rw, r)
        
        // Log the response status code
        logger.InfoLogger.Printf("Auth response status: %d for %s %s", rw.status, r.Method, r.URL.Path)
    })
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
    http.ResponseWriter
    status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    if rw.status == 0 {
        rw.status = http.StatusOK
    }
    return rw.ResponseWriter.Write(b)
}