// backend/go/internal/models/auth.go
package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID            string    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at" db:"last_login_at"`
	ProviderType  string    `json:"provider_type" db:"provider_type"`
	ProviderID    string    `json:"provider_id" db:"provider_id"`
	RefreshToken  string    `json:"-" db:"refresh_token"`
	ProfilePicURL string    `json:"profile_pic_url" db:"profile_pic_url"`
}

// SignInRequest represents the request body for sign-in
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpRequest represents the request body for sign-up
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// OAuthRequest represents the request body for OAuth sign-in/sign-up
type OAuthRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	Provider    string `json:"provider"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

// OAuthConfig contains the configuration for OAuth providers
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       []string
}