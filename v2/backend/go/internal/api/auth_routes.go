// backend/go/internal/api/auth_routes.go
package api

import (
	"net/http"
	"sage-ai-v2/internal/api/handlers"
	"sage-ai-v2/internal/services"
)

func AddAuthRoutes(mux *http.ServeMux, authService *services.AuthService) {
	// Create auth handler
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	mux.HandleFunc("/api/auth/signin", authHandler.SignInHandler)
	mux.HandleFunc("/api/auth/signup", authHandler.SignUpHandler)
	mux.HandleFunc("/api/auth/signout", authHandler.SignOutHandler)

	// OAuth routes
	mux.HandleFunc("/api/auth/oauth/google", authHandler.OAuthSignInHandler)
	mux.HandleFunc("/api/auth/oauth/github", authHandler.OAuthSignInHandler)
	mux.HandleFunc("/api/auth/oauth/url/google", authHandler.OAuthURLHandler)
	mux.HandleFunc("/api/auth/oauth/url/github", authHandler.OAuthURLHandler)

	// Protected routes
	mux.Handle("/api/auth/user", authHandler.AuthMiddleware(http.HandlerFunc(authHandler.GetUserHandler)))
}