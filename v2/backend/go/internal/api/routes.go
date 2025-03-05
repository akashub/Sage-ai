// // backend/go/internal/api/routes.go
// // backend/go/internal/api/routes.go
// package api

// import (
// 	"net/http"
// 	"sage-ai-v2/internal/api/handlers"
// 	"sage-ai-v2/internal/api/middleware"
// 	"sage-ai-v2/pkg/logger"
// )

// // SetupRoutes configures all API routes for the application
// func SetupRoutes() http.Handler {
// 	// Create a new ServeMux
// 	mux := http.NewServeMux()

// 	// Register API routes
// 	mux.HandleFunc("/api/upload", handlers.UploadFileHandler)
// 	mux.HandleFunc("/api/query", handlers.QueryHandler)

// 	// Add health check endpoint
// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/plain")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	})

// 	// Apply middleware to all routes
// 	handler := middleware.ApplyAuth(mux)

// 	logger.InfoLogger.Printf("API routes configured")
// 	return handler
// }

// backend/go/internal/api/routes.go
// backend/go/internal/api/routes.go
package api

import (
	"database/sql"
	"net/http"
	"sage-ai-v2/internal/api/handlers"
	"sage-ai-v2/internal/models"
	"sage-ai-v2/internal/services"
	"sage-ai-v2/pkg/logger"
	"time"
)

// SetupRoutes configures all API routes for the application
func SetupRoutes(db *sql.DB) http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Initialize authentication service
	authService := setupAuthService(db)
	
	// Add authentication routes
	AddAuthRoutes(mux, authService)

	// Register API routes
	mux.HandleFunc("/api/upload", handlers.UploadFileHandler)
	mux.HandleFunc("/api/query", handlers.QueryHandler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware directly without using a separate middleware package
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		mux.ServeHTTP(w, r)
	})

	logger.InfoLogger.Printf("API routes configured")
	return handler
}

// setupAuthService creates and configures the authentication service
func setupAuthService(db *sql.DB) *services.AuthService {
	// JWT configuration
	jwtSecret := "your-secret-key" // In production, load from environment variable
	jwtExpiry := 7 * 24 * time.Hour // 7 days
	
	// OAuth configurations
	oauthConfs := map[string]models.OAuthConfig{
		"google": {
			ClientID:     "64583008448-4aa9mivl1jurlp1bheabkc5m0irc6fsp.apps.googleusercontent.com",
			ClientSecret: "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",
			RedirectURI:  "http://localhost:3000/oauth-callback",
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
			Scopes:       []string{"email", "profile"},
		},
		"github": {
			ClientID:     "Ov23liJMbcmt6eXGI7yN",
			ClientSecret: "ae1a717b3b1311ba7e3af4f356d37019fac61639",
			RedirectURI:  "http://localhost:3000/oauth-callback",
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
			Scopes:       []string{"read:user", "user:email"},
		},
	}
	
	return services.NewAuthService(db, jwtSecret, jwtExpiry, oauthConfs)
}