// backend/go/internal/api/chat_routes.go
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sage-ai-v2/pkg/logger"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// ChatMessage represents a single message in a chat session
type ChatMessage struct {
	Type            string                 `json:"type"`
	Text            string                 `json:"text"`
	Sender          string                 `json:"sender,omitempty"`
	Timestamp       time.Time              `json:"timestamp"`
	Results         []map[string]interface{} `json:"results,omitempty"`
	GeneratedQuery  string                 `json:"generatedQuery,omitempty"`
	KnowledgeContext []map[string]interface{} `json:"knowledgeContext,omitempty"`
	File            string                 `json:"file,omitempty"`
}

// Chat represents a chat session
type Chat struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	File        string        `json:"file,omitempty"`
	FilePath    string        `json:"filePath,omitempty"`
	Headers     []string      `json:"headers,omitempty"`
	Messages    []ChatMessage `json:"messages,omitempty"`
	CreatedAt   time.Time     `json:"timestamp"`
	LastUpdated time.Time     `json:"lastUpdated,omitempty"`
	TrainingDataIDs  []string      `json:"trainingDataIds,omitempty"`
}

// ChatStore is a simple in-memory storage for chats
type ChatStore struct {
	Chats    map[string]Chat
	DataPath string
}

// NewChatStore creates a new chat store
func NewChatStore(dataPath string) *ChatStore {
	store := &ChatStore{
		Chats:    make(map[string]Chat),
		DataPath: dataPath,
	}
	
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		logger.ErrorLogger.Printf("Failed to create chat data directory: %v", err)
	}
	
	// Load existing chats
	store.loadChats()
	
	return store
}

// loadChats loads chats from disk
func (cs *ChatStore) loadChats() {
	// Read all files in the data directory
	files, err := os.ReadDir(cs.DataPath)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to read chat data directory: %v", err)
		return
	}
	
	// Load each chat file
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			chatPath := filepath.Join(cs.DataPath, file.Name())
			chatData, err := os.ReadFile(chatPath)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to read chat file %s: %v", file.Name(), err)
				continue
			}
			
			var chat Chat
			if err := json.Unmarshal(chatData, &chat); err != nil {
				logger.ErrorLogger.Printf("Failed to unmarshal chat data from %s: %v", file.Name(), err)
				continue
			}
			
			// Add to map
			cs.Chats[chat.ID] = chat
		}
	}
	
	logger.InfoLogger.Printf("Loaded %d chats from disk", len(cs.Chats))
}

// saveChat saves a chat to disk
func (cs *ChatStore) saveChat(chat Chat) error {
	// Marshal to JSON
	chatData, err := json.MarshalIndent(chat, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chat data: %w", err)
	}
	
	// Write to file
	chatPath := filepath.Join(cs.DataPath, chat.ID+".json")
	if err := os.WriteFile(chatPath, chatData, 0644); err != nil {
		return fmt.Errorf("failed to write chat file: %w", err)
	}
	
	return nil
}

// AddChat adds a new chat
func (cs *ChatStore) AddChat(chat Chat) error {
	// Add to map
	cs.Chats[chat.ID] = chat
	
	// Save to disk
	return cs.saveChat(chat)
}

// GetChat retrieves a chat by ID
func (cs *ChatStore) GetChat(id string) (Chat, bool) {
	chat, exists := cs.Chats[id]
	return chat, exists
}

// UpdateChat updates an existing chat
func (cs *ChatStore) UpdateChat(id string, updates map[string]interface{}) (Chat, error) {
	chat, exists := cs.Chats[id]
	if !exists {
		return Chat{}, fmt.Errorf("chat not found: %s", id)
	}
	
	// Apply updates
	if title, ok := updates["title"].(string); ok {
		chat.Title = title
	}
	
	if file, ok := updates["file"].(string); ok {
		chat.File = file
	}
	
	if filePath, ok := updates["filePath"].(string); ok {
		chat.FilePath = filePath
	}
	
	if headers, ok := updates["headers"].([]string); ok {
		chat.Headers = headers
	}
	
	if messages, ok := updates["messages"]; ok {
		// Convert messages to proper type
		messagesJSON, err := json.Marshal(messages)
		if err != nil {
			return Chat{}, fmt.Errorf("failed to marshal messages: %w", err)
		}
		
		var typedMessages []ChatMessage
		if err := json.Unmarshal(messagesJSON, &typedMessages); err != nil {
			return Chat{}, fmt.Errorf("failed to unmarshal messages: %w", err)
		}
		
		chat.Messages = typedMessages
	}
	
	// Update last updated timestamp
	if lastUpdated, ok := updates["lastUpdated"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, lastUpdated)
		if err == nil {
			chat.LastUpdated = parsedTime
		}
	} else {
		chat.LastUpdated = time.Now()
	}
	
	// Update in map
	cs.Chats[id] = chat
	
	// Save to disk
	if err := cs.saveChat(chat); err != nil {
		return Chat{}, err
	}
	
	return chat, nil
}

// DeleteChat deletes a chat
func (cs *ChatStore) DeleteChat(id string) error {
	// Remove from map
	delete(cs.Chats, id)
	
	// Delete file
	chatPath := filepath.Join(cs.DataPath, id+".json")
	if err := os.Remove(chatPath); err != nil {
		return fmt.Errorf("failed to delete chat file: %w", err)
	}
	
	return nil
}

// ListChats returns all chats, sorted by lastUpdated (most recent first)
func (cs *ChatStore) ListChats() []Chat {
	chats := make([]Chat, 0, len(cs.Chats))
	for _, chat := range cs.Chats {
		chats = append(chats, chat)
	}
	
	// Sort by lastUpdated (most recent first)
	sort.Slice(chats, func(i, j int) bool {
		// If LastUpdated is zero, use CreatedAt
		timeI := chats[i].LastUpdated
		if timeI.IsZero() {
			timeI = chats[i].CreatedAt
		}
		
		timeJ := chats[j].LastUpdated
		if timeJ.IsZero() {
			timeJ = chats[j].CreatedAt
		}
		
		return timeI.After(timeJ)
	})
	
	return chats
}

// Initialize the chat store
var chatStore *ChatStore

// SetupChatRoutes registers chat-related routes
// func SetupChatRoutes(router *mux.Router, km *knowledge.KnowledgeManager) {
// 	// Initialize chat store
// 	chatStore = NewChatStore("./data/chats")
	
// 	// Register routes
// 	router.HandleFunc("/api/chats", GetChatsHandler).Methods("GET", "OPTIONS")
//     router.HandleFunc("/api/chats", CreateChatHandler).Methods("POST", "OPTIONS")
//     router.HandleFunc("/api/chats/{id}", GetChatHandler).Methods("GET", "OPTIONS")
//     router.HandleFunc("/api/chats/{id}", UpdateChatHandler).Methods("PUT", "OPTIONS")
//     router.HandleFunc("/api/chats/{id}", DeleteChatHandler).Methods("DELETE", "OPTIONS")

// 	// Add routes for chat-specific training data
//     if km != nil {
//         router.HandleFunc("/api/chats/{id}/training", func(w http.ResponseWriter, r *http.Request) {
//             getChatTrainingDataHandler(w, r, km)
//         }).Methods("GET", "OPTIONS")
        
//         router.HandleFunc("/api/chats/{id}/training", func(w http.ResponseWriter, r *http.Request) {
//             updateChatTrainingDataHandler(w, r, km)
//         }).Methods("POST", "OPTIONS")
//     }
// }
func SetupChatRoutes(router *mux.Router) {
    // Initialize chat store
    chatStore = NewChatStore("./data/chats")
    
    // Register routes
    router.HandleFunc("/api/chats", GetChatsHandler).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/chats", CreateChatHandler).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", GetChatHandler).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", UpdateChatHandler).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", DeleteChatHandler).Methods("DELETE", "OPTIONS")
    
    // Add simplified training data routes - without using KM
    router.HandleFunc("/api/chats/{id}/training", getChatTrainingDataSimple).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/chats/{id}/training", updateChatTrainingDataSimple).Methods("POST", "OPTIONS")
}

// Add these simplified handlers that don't use the KM:

// getChatTrainingDataSimple gets training data IDs for a chat
// func getChatTrainingDataSimple(w http.ResponseWriter, r *http.Request) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // Just return the training data IDs instead of trying to load the actual data
//     response := map[string]interface{}{
//         "trainingDataIds": chat.TrainingDataIDs,
//     }
    
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }
func getChatTrainingDataSimple(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Extract chat ID
    vars := mux.Vars(r)
    chatID := vars["id"]
    
    // Get the chat
    chat, exists := chatStore.GetChat(chatID)
    if !exists {
        http.Error(w, "Chat not found", http.StatusNotFound)
        return
    }
    
    // Just return the training data IDs instead of trying to load the actual data
    response := map[string]interface{}{
        "trainingDataIds": chat.TrainingDataIDs,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


// updateChatTrainingDataSimple updates training data IDs for a chat
// func updateChatTrainingDataSimple(w http.ResponseWriter, r *http.Request) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // Parse request
//     var req struct {
//         TrainingDataIDs []string `json:"trainingDataIds"`
//     }
    
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }
    
//     // Update the chat
//     chat.TrainingDataIDs = req.TrainingDataIDs
//     chat.LastUpdated = time.Now()
    
//     // Save the chat
//     if err := chatStore.AddChat(chat); err != nil {
//         logger.ErrorLogger.Printf("Failed to update chat: %v", err)
//         http.Error(w, "Failed to update chat", http.StatusInternalServerError)
//         return
//     }
    
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(chat)
// }
func updateChatTrainingDataSimple(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Extract chat ID
    vars := mux.Vars(r)
    chatID := vars["id"]
    
    // Get the chat
    chat, exists := chatStore.GetChat(chatID)
    if !exists {
        http.Error(w, "Chat not found", http.StatusNotFound)
        return
    }
    
    // Parse request
    var req struct {
        TrainingDataIDs []string `json:"trainingDataIds"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Update the chat
    chat.TrainingDataIDs = req.TrainingDataIDs
    chat.LastUpdated = time.Now()
    
    // Save the chat
    if err := chatStore.AddChat(chat); err != nil {
        logger.ErrorLogger.Printf("Failed to update chat: %v", err)
        http.Error(w, "Failed to update chat", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chat)
}

// GetChatsHandler returns all chats
func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Existing code continues...
    chats := chatStore.ListChats()
    
    // For list view, we don't need to include all messages
    for i := range chats {
        // Keep only the most recent message for preview
        if len(chats[i].Messages) > 0 {
            chats[i].Messages = []ChatMessage{chats[i].Messages[len(chats[i].Messages)-1]}
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chats)
}

// CreateChatHandler creates a new chat
func CreateChatHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var chatData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&chatData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Generate ID
	chatID := fmt.Sprintf("chat_%d", time.Now().UnixNano())
	
	// Create chat
	chat := Chat{
		ID:        chatID,
		Title:     fmt.Sprintf("New Chat %s", time.Now().Format("2006-01-02")),
		CreatedAt: time.Now(),
	}
	
	// Set title if provided
	if title, ok := chatData["title"].(string); ok && title != "" {
		chat.Title = title
	}
	
	// Add chat
	if err := chatStore.AddChat(chat); err != nil {
		logger.ErrorLogger.Printf("Failed to create chat: %v", err)
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}
	
	// Return the new chat
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// GetChatHandler returns a specific chat
func GetChatHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	vars := mux.Vars(r)
	chatID := vars["id"]
	
	chat, exists := chatStore.GetChat(chatID)
	if !exists {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// UpdateChatHandler updates a chat
func UpdateChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	vars := mux.Vars(r)
	chatID := vars["id"]
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	chat, err := chatStore.UpdateChat(chatID, updates)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to update chat: %v", err)
		http.Error(w, "Failed to update chat", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// DeleteChatHandler deletes a chat
func DeleteChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	vars := mux.Vars(r)
	chatID := vars["id"]
	
	if err := chatStore.DeleteChat(chatID); err != nil {
		logger.ErrorLogger.Printf("Failed to delete chat: %v", err)
		http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func SetupChatRoutesForTesting(router *mux.Router, store *ChatStore) {
	// Set the package-level chatStore to the test store
	chatStore = store
	
	// Register routes
	router.HandleFunc("/api/chats", GetChatsHandler).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/chats", CreateChatHandler).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", GetChatHandler).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", UpdateChatHandler).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/chats/{id}", DeleteChatHandler).Methods("DELETE", "OPTIONS")
}

// getChatTrainingDataHandler gets training data for a chat
// func getChatTrainingDataHandler(w http.ResponseWriter, r *http.Request, km *knowledge.KnowledgeManager) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // Get training data items
//     if len(chat.TrainingDataIDs) == 0 {
//         // Return empty array if no training data
//         w.Header().Set("Content-Type", "application/json")
//         json.NewEncoder(w).Encode([]interface{}{})
//         return
//     }
    
//     // Get all training data
//     allData, err := km.ListTrainingData(r.Context(), "")
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to list training data: %v", err)
//         http.Error(w, "Failed to list training data", http.StatusInternalServerError)
//         return
//     }
    
//     // Filter for items that belong to this chat
//     chatTrainingData := []map[string]interface{}{}
//     for _, item := range allData {
//         itemID, ok := item["id"].(string)
//         if !ok {
//             continue
//         }
        
//         for _, id := range chat.TrainingDataIDs {
//             if itemID == id {
//                 chatTrainingData = append(chatTrainingData, item)
//                 break
//             }
//         }
//     }
    
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(chatTrainingData)
// }

// In chat_routes.go - Get training data for a specific chat
// func getChatTrainingDataHandler(w http.ResponseWriter, r *http.Request, km *knowledge.KnowledgeManager) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID from URL
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // If knowledge manager is available, fetch actual training data items
//     if km != nil && len(chat.TrainingDataIDs) > 0 {
//         // Get all training data first
//         allItems, err := km.ListTrainingData(r.Context(), "")
//         if err != nil {
//             logger.ErrorLogger.Printf("Failed to list training data: %v", err)
//             http.Error(w, fmt.Sprintf("Failed to list training data: %v", err), http.StatusInternalServerError)
//             return
//         }
        
//         // Filter to include only items associated with this chat
//         chatTrainingData := []map[string]interface{}{}
//         for _, item := range allItems {
//             itemID, ok := item["id"].(string)
//             if !ok {
//                 continue
//             }
            
//             for _, id := range chat.TrainingDataIDs {
//                 if itemID == id {
//                     chatTrainingData = append(chatTrainingData, item)
//                     break
//                 }
//             }
//         }
        
//         // Return both IDs and actual data items
//         response := map[string]interface{}{
//             "trainingDataIds": chat.TrainingDataIDs,
//             "trainingData": chatTrainingData,
//         }
        
//         w.Header().Set("Content-Type", "application/json")
//         json.NewEncoder(w).Encode(response)
//         return
//     }
    
//     // If no knowledge manager or no training data IDs, just return the IDs
//     response := map[string]interface{}{
//         "trainingDataIds": chat.TrainingDataIDs,
//     }
    
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

// // updateChatTrainingDataHandler updates training data for a chat
// func updateChatTrainingDataHandler(w http.ResponseWriter, r *http.Request, km *knowledge.KnowledgeManager) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // Parse request
//     var req struct {
//         TrainingDataIDs []string `json:"trainingDataIds"`
//     }
    
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }
    
//     // Update the chat
//     chat.TrainingDataIDs = req.TrainingDataIDs
//     chat.LastUpdated = time.Now()
    
//     // Save the chat
//     if err := chatStore.AddChat(chat); err != nil {
//         logger.ErrorLogger.Printf("Failed to update chat: %v", err)
//         http.Error(w, "Failed to update chat", http.StatusInternalServerError)
//         return
//     }
    
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(chat)
// }

// // In chat_routes.go - Update training data for a specific chat
// func updateChatTrainingDataHandler(w http.ResponseWriter, r *http.Request, km *knowledge.KnowledgeManager) {
//     // Set CORS headers
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
    
//     // Extract chat ID from URL
//     vars := mux.Vars(r)
//     chatID := vars["id"]
    
//     // Get the chat
//     chat, exists := chatStore.GetChat(chatID)
//     if !exists {
//         http.Error(w, "Chat not found", http.StatusNotFound)
//         return
//     }
    
//     // Parse request body
//     var req struct {
//         TrainingDataIDs []string `json:"trainingDataIds"`
//     }
    
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }
    
//     // Update the chat with new training data IDs
//     chat.TrainingDataIDs = req.TrainingDataIDs
//     chat.LastUpdated = time.Now()
    
//     // Save the updated chat
//     if err := chatStore.AddChat(chat); err != nil {
//         logger.ErrorLogger.Printf("Failed to update chat with training data: %v", err)
//         http.Error(w, "Failed to update chat", http.StatusInternalServerError)
//         return
//     }
    
//     // Return the updated chat
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(chat)
// }