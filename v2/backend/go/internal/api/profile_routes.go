// backend/go/internal/api/profile_routes.go
package api

import (
	"encoding/json"
	"net/http"
	"sage-ai-v2/pkg/logger"
	"time"

	"github.com/gorilla/mux"
)

// UserProfile represents a user's profile data
type UserProfile struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"createdAt"`
	LastLoginAt   time.Time `json:"lastLoginAt"`
	Plan          string    `json:"plan"`
	ProfilePicURL string    `json:"profilePicUrl,omitempty"`
	Stats         UserStats `json:"stats"`
}

// UserStats represents usage statistics for a user
type UserStats struct {
	TotalChats     int `json:"totalChats"`
	TotalQueries   int `json:"totalQueries"`
	APICredentials int `json:"apiCredentials"`
}

// SetupProfileRoutes registers profile-related routes
func SetupProfileRoutes(router *mux.Router) {
	router.HandleFunc("/api/profile", GetProfileHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/profile", UpdateProfileHandler).Methods("PUT", "OPTIONS")
}

// GetProfileHandler returns the user's profile
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get user ID from token/auth
	// In a real implementation, this would come from the JWT token
	// For now, we'll use a mock ID
	userID := "user123" // Would come from authentication middleware

	// In a real implementation, this would be fetched from the database
	// For demo purposes, we're creating mock data
	profile := UserProfile{
		ID:            userID,
		Name:          "Aakash Singh",
		Email:         "aakashsinghas03@gmail.com",
		CreatedAt:     time.Now().AddDate(0, -3, 0), // 3 months ago
		LastLoginAt:   time.Now(),
		Plan:          "Free",
		ProfilePicURL: "https://ui-avatars.com/api/?name=Test+User&background=5865F2&color=fff",
		Stats: UserStats{
			TotalChats:     len(chatStore.Chats),
			TotalQueries:   calculateTotalQueries(userID),
			APICredentials: 2, // Mock data
		},
	}

	// Count user's chats
	userChats := getUserChats(userID)
	profile.Stats.TotalChats = len(userChats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfileHandler updates the user's profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get user ID from token/auth
	userID := "user123" // Would come from authentication middleware

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would update the user in the database
	// For demo purposes, we'll just log and return success
	logger.InfoLogger.Printf("Profile update for user %s: %+v", userID, updates)

	// Return updated profile (mock data for now)
	profile := UserProfile{
		ID:            userID,
		Name:          updates["name"].(string),
		Email:         "user@example.com", // Email typically can't be changed easily
		CreatedAt:     time.Now().AddDate(0, -3, 0),
		LastLoginAt:   time.Now(),
		Plan:          "Free",
		ProfilePicURL: updates["profilePicUrl"].(string),
		Stats: UserStats{
			TotalChats:     len(chatStore.Chats),
			TotalQueries:   calculateTotalQueries(userID),
			APICredentials: 2,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// Helper functions for profile data

// getUserChats returns all chats for a specific user
func getUserChats(userID string) []Chat {
	// Filter chats by user ID
	var userChats []Chat
	// In a real implementation, you would have a UserID field in the Chat struct
	// For now, we'll return all chats as a demo
	userChats = append(userChats, chatStore.ListChats()...)
	return userChats
}

// calculateTotalQueries counts the total number of queries made by a user
func calculateTotalQueries(userID string) int {
	totalQueries := 0
	userChats := getUserChats(userID)

	for _, chat := range userChats {
		// Count user messages as queries
		for _, msg := range chat.Messages {
			if msg.Sender == "user" {
				totalQueries++
			}
		}
	}

	return totalQueries
}
