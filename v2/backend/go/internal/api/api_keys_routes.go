// backend/go/internal/api/api_keys_routes.go
package api

import (
	"encoding/json"
	"net/http"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// APIKeyEntry represents a saved API key
type APIKeyEntry struct {
	ID          string        `json:"id"`
	UserID      string        `json:"userId"`
	Provider    llm.LLMProvider `json:"provider"`
	Name        string        `json:"name"`
	LastUsed    time.Time     `json:"lastUsed"`
	CreatedAt   time.Time     `json:"createdAt"`
	IsDefault   bool          `json:"isDefault"`
	MaskedKey   string        `json:"maskedKey"`
}

// SetupAPIKeysRoutes registers API keys management routes
func SetupAPIKeysRoutes(router *mux.Router) {
	router.HandleFunc("/api/apikeys", GetAPIKeysHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/apikeys", SaveAPIKeyHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/apikeys/{id}", DeleteAPIKeyHandler).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/apikeys/{id}/default", SetDefaultAPIKeyHandler).Methods("PUT", "OPTIONS")
}

// GetAPIKeysHandler returns all saved API keys for a user
func GetAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get user ID from token/auth
	userID := "user123" // Would come from authentication middleware

	// Mock data for demo
	keys := []APIKeyEntry{
		{
			ID:        "key1",
			UserID:    userID,
			Provider:  llm.ProviderGemini,
			Name:      "Gemini API Key",
			LastUsed:  time.Now().AddDate(0, 0, -2), // 2 days ago
			CreatedAt: time.Now().AddDate(0, -1, 0), // 1 month ago
			IsDefault: true,
			MaskedKey: "***********************************ABC",
		},
		{
			ID:        "key2",
			UserID:    userID,
			Provider:  llm.ProviderOpenAI,
			Name:      "OpenAI GPT-4",
			LastUsed:  time.Now().AddDate(0, 0, -5), // 5 days ago
			CreatedAt: time.Now().AddDate(0, -2, 0), // 2 months ago
			IsDefault: false,
			MaskedKey: "sk-***********************************XYZ",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// SaveAPIKeyHandler saves a new API key for a user
func SaveAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get user ID from token/auth
	userID := "user123" // Would come from authentication middleware

	var keyData struct {
		Provider llm.LLMProvider `json:"provider"`
		APIKey   string          `json:"apiKey"`
		Name     string          `json:"name"`
		IsDefault bool           `json:"isDefault"`
	}

	if err := json.NewDecoder(r.Body).Decode(&keyData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the API key
	if keyData.APIKey == "" {
		http.Error(w, "API key is required", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would save to the database
	// For now, we'll just return success with a mock response
	newKey := APIKeyEntry{
		ID:        "key_" + time.Now().Format("20060102150405"),
		UserID:    userID,
		Provider:  keyData.Provider,
		Name:      keyData.Name,
		LastUsed:  time.Now(),
		CreatedAt: time.Now(),
		IsDefault: keyData.IsDefault,
		MaskedKey: maskAPIKey(keyData.APIKey),
	}

	logger.InfoLogger.Printf("Saved API key for user %s with provider %s", userID, keyData.Provider)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newKey)
}

// DeleteAPIKeyHandler deletes an API key
func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get user ID from token/auth
	userID := "user123" // Would come from authentication middleware

	// Get key ID from URL path
	vars := mux.Vars(r)
	keyID := vars["id"]

	// In a real implementation, this would delete from the database
	// For now, we'll just return success
	logger.InfoLogger.Printf("Deleted API key %s for user %s", keyID, userID)

	w.WriteHeader(http.StatusNoContent)
}

// SetDefaultAPIKeyHandler sets an API key as the default for a user
func SetDefaultAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get key ID from URL path
	vars := mux.Vars(r)
	keyID := vars["id"]

	// In a real implementation, this would update the database
	// For now, we'll just return success
	logger.InfoLogger.Printf("Set API key %s as default for user %s", keyID, userID)

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to mask API keys for display
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 3 {
		return "****"
	}
	
	// Keep the first few characters if it has a recognizable prefix
	prefixLength := 0
	if strings.HasPrefix(apiKey, "sk-") {
		prefixLength = 3
	}
	
	// Show only the last 3 characters
	masked := strings.Repeat("*", len(apiKey)-3-prefixLength)
	return apiKey[:prefixLength] + masked + apiKey[len(apiKey)-3:]
}