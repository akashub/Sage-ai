// backend/go/tests/api/chat_test.go
package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sage-ai-v2/internal/api"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestChatStore creates a test chat store with a temporary directory
func setupTestChatStore(t *testing.T) (*api.ChatStore, string) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "chat-test")
	require.NoError(t, err)
	
	// Create chat store using the temporary directory
	store := api.NewChatStore(tempDir)
	
	return store, tempDir
}

// TestChatStore tests the basic functionality of the chat store
func TestChatStore(t *testing.T) {
	// Set up test environment
	store, tempDir := setupTestChatStore(t)
	defer os.RemoveAll(tempDir)
	
	// Create a test chat
	testChat := api.Chat{
		ID:        "chat_1234567890",
		Title:     "Test Chat",
		File:      "test.csv",
		FilePath:  "data/uploads/test.csv",
		Headers:   []string{"id", "name", "email"},
		CreatedAt: time.Now(),
		Messages: []api.ChatMessage{
			{
				Type:      "system",
				Text:      "Chat started",
				Timestamp: time.Now(),
			},
		},
	}
	
	// Test AddChat
	err := store.AddChat(testChat)
	assert.NoError(t, err)
	
	// Test GetChat
	chat, exists := store.GetChat("chat_1234567890")
	assert.True(t, exists)
	assert.Equal(t, testChat.ID, chat.ID)
	assert.Equal(t, testChat.Title, chat.Title)
	assert.Equal(t, testChat.File, chat.File)
	assert.Equal(t, len(testChat.Messages), len(chat.Messages))
	
	// Test UpdateChat
	updates := map[string]interface{}{
		"title": "Updated Chat Title",
		"messages": []api.ChatMessage{
			{
				Type:      "system",
				Text:      "Chat started",
				Timestamp: time.Now(),
			},
			{
				Type:      "user",
				Text:      "Hello",
				Sender:    "user",
				Timestamp: time.Now(),
			},
		},
	}
	
	updatedChat, err := store.UpdateChat("chat_1234567890", updates)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Chat Title", updatedChat.Title)
	assert.Len(t, updatedChat.Messages, 2)
	
	// Test ListChats
	chats := store.ListChats()
	assert.Len(t, chats, 1)
	assert.Equal(t, "chat_1234567890", chats[0].ID)
	
	// Test DeleteChat
	err = store.DeleteChat("chat_1234567890")
	assert.NoError(t, err)
	
	// Verify chat was deleted
	_, exists = store.GetChat("chat_1234567890")
	assert.False(t, exists)
	
	chats = store.ListChats()
	assert.Len(t, chats, 0)
}

// TestChatRoutes tests the chat API routes
func TestChatRoutes(t *testing.T) {
	// Create router
	router := mux.NewRouter()
	
	// Set up test environment
	chatStore, tempDir := setupTestChatStore(t)
	defer os.RemoveAll(tempDir)
	
	// Replace ChatStore with test store
	api.SetupChatRoutesForTesting(router, chatStore)
	
	// Test creating a chat
	t.Run("CreateChat", func(t *testing.T) {
		chatData := map[string]interface{}{
			"title": "New Test Chat",
			"file":  "test.csv",
		}
		
		jsonBody, err := json.Marshal(chatData)
		require.NoError(t, err)
		
		req, err := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var response api.Chat
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "New Test Chat", response.Title)
		
		// Verify chat was added to store
		chat, exists := chatStore.GetChat(response.ID)
		assert.True(t, exists)
		assert.Equal(t, response.ID, chat.ID)
	})
	
	// Test getting chat list
	t.Run("GetChats", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/chats", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var chats []api.Chat
		err = json.Unmarshal(rr.Body.Bytes(), &chats)
		require.NoError(t, err)
		
		assert.Len(t, chats, 1)
		assert.Equal(t, "New Test Chat", chats[0].Title)
	})
	
	// Get the chat ID from previous test
	var chatID string
	chats := chatStore.ListChats()
	if len(chats) > 0 {
		chatID = chats[0].ID
	}
	
	// Test getting a single chat
	t.Run("GetChat", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/chats/%s", chatID), nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var chat api.Chat
		err = json.Unmarshal(rr.Body.Bytes(), &chat)
		require.NoError(t, err)
		
		assert.Equal(t, chatID, chat.ID)
		assert.Equal(t, "New Test Chat", chat.Title)
	})
	
	// Test updating a chat
	t.Run("UpdateChat", func(t *testing.T) {
		updateData := map[string]interface{}{
			"title": "Updated Chat Title",
			"messages": []map[string]interface{}{
				{
					"type":      "user",
					"text":      "Test message",
					"sender":    "user",
					"timestamp": time.Now(),
				},
			},
		}
		
		jsonBody, err := json.Marshal(updateData)
		require.NoError(t, err)
		
		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/chats/%s", chatID), bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var chat api.Chat
		err = json.Unmarshal(rr.Body.Bytes(), &chat)
		require.NoError(t, err)
		
		assert.Equal(t, "Updated Chat Title", chat.Title)
		assert.Len(t, chat.Messages, 1)
		
		// Verify chat was updated in store
		updatedChat, exists := chatStore.GetChat(chatID)
		assert.True(t, exists)
		assert.Equal(t, "Updated Chat Title", updatedChat.Title)
	})
	
	// Test deleting a chat
	t.Run("DeleteChat", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/chats/%s", chatID), nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusNoContent, rr.Code)
		
		// Verify chat was deleted from store
		_, exists := chatStore.GetChat(chatID)
		assert.False(t, exists)
		
		chats := chatStore.ListChats()
		assert.Len(t, chats, 0)
	})
}

// Helper function to add to api/chat_routes.go for testing
// Add this to your api/chat_routes.go file
