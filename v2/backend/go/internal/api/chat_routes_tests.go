// backend/go/internal/api/chat_routes_test.go
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupChatTestEnvironment(t *testing.T) (*ChatStore, func()) {
	// Create a temporary directory for chat data
	tempDir, err := os.MkdirTemp("", "chat-test")
	require.NoError(t, err)

	// Create chat store with temporary directory
	chatStore := NewChatStore(tempDir)
	require.NotNil(t, chatStore)

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return chatStore, cleanup
}

func setupChatRouter(t *testing.T) (*mux.Router, *ChatStore, func()) {
	// Set up chat store
	store, cleanup := setupChatTestEnvironment(t)
	
	// Create router and register chat routes
	router := mux.NewRouter()
	
	// Store the global chatStore
	oldChatStore := chatStore
	chatStore = store
	
	// Register routes
	SetupChatRoutes(router)
	
	// Enhance cleanup to restore global state
	enhancedCleanup := func() {
		cleanup()
		chatStore = oldChatStore
	}
	
	return router, store, enhancedCleanup
}

func TestCreateChat(t *testing.T) {
	router, _, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create a test chat request
	reqBody := map[string]interface{}{
		"title": "Test Chat",
		"file": "test_file.csv",
	}
	
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response Chat
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Test Chat", response.Title)
	assert.Equal(t, "test_file.csv", response.File)
	assert.False(t, response.CreatedAt.IsZero())
}

func TestGetChatById(t *testing.T) {
	router, store, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create a test chat
	testChat := Chat{
		ID:        "test-chat-1",
		Title:     "Test Chat",
		File:      "test_file.csv",
		CreatedAt: time.Now(),
		Messages: []ChatMessage{
			{
				Type:      "user",
				Text:      "Test message",
				Sender:    "user",
				Timestamp: time.Now(),
			},
		},
	}
	
	// Add to store
	err := store.AddChat(testChat)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("GET", "/api/chats/test-chat-1", nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response Chat
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response
	assert.Equal(t, testChat.ID, response.ID)
	assert.Equal(t, testChat.Title, response.Title)
	assert.Equal(t, testChat.File, response.File)
	assert.Len(t, response.Messages, 1)
	assert.Equal(t, testChat.Messages[0].Text, response.Messages[0].Text)
}

func TestGetAllChats(t *testing.T) {
	router, store, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create test chats
	chats := []Chat{
		{
			ID:        "test-chat-1",
			Title:     "Test Chat 1",
			File:      "file1.csv",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			LastUpdated: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "test-chat-2",
			Title:     "Test Chat 2",
			File:      "file2.csv",
			CreatedAt: time.Now().Add(-1 * time.Hour),
			LastUpdated: time.Now(),
		},
	}
	
	// Add to store
	for _, chat := range chats {
		err := store.AddChat(chat)
		require.NoError(t, err)
	}
	
	// Create request
	req, err := http.NewRequest("GET", "/api/chats", nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response []Chat
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response
	assert.Len(t, response, 2)
	
	// Chats should be sorted by last updated (most recent first)
	assert.Equal(t, "test-chat-2", response[0].ID)
	assert.Equal(t, "test-chat-1", response[1].ID)
}

func TestUpdateChat(t *testing.T) {
	router, store, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create a test chat
	testChat := Chat{
		ID:        "test-chat-1",
		Title:     "Test Chat",
		File:      "test_file.csv",
		CreatedAt: time.Now(),
	}
	
	// Add to store
	err := store.AddChat(testChat)
	require.NoError(t, err)
	
	// Create update request
	updateData := map[string]interface{}{
		"title": "Updated Title",
		"messages": []ChatMessage{
			{
				Type:      "user",
				Text:      "New message",
				Sender:    "user",
				Timestamp: time.Now(),
			},
		},
	}
	
	body, err := json.Marshal(updateData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("PUT", "/api/chats/test-chat-1", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var response Chat
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response
	assert.Equal(t, testChat.ID, response.ID)
	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, testChat.File, response.File)
	assert.Len(t, response.Messages, 1)
	assert.Equal(t, "New message", response.Messages[0].Text)
	assert.False(t, response.LastUpdated.IsZero())
}

func TestDeleteChat(t *testing.T) {
	router, store, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create a test chat
	testChat := Chat{
		ID:        "test-chat-1",
		Title:     "Test Chat",
		File:      "test_file.csv",
		CreatedAt: time.Now(),
	}
	
	// Add to store
	err := store.AddChat(testChat)
	require.NoError(t, err)
	
	// First, verify the chat exists
	_, exists := store.GetChat(testChat.ID)
	assert.True(t, exists)
	
	// Create request
	req, err := http.NewRequest("DELETE", "/api/chats/test-chat-1", nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusNoContent, rr.Code)
	
	// Verify the chat is deleted
	_, exists = store.GetChat(testChat.ID)
	assert.False(t, exists)
}

func TestChatStoreOperations(t *testing.T) {
	store, cleanup := setupChatTestEnvironment(t)
	defer cleanup()
	
	// Test AddChat and GetChat
	chat := Chat{
		ID:        "test-chat-ops",
		Title:     "Test Operations",
		File:      "test.csv",
		CreatedAt: time.Now(),
	}
	
	err := store.AddChat(chat)
	assert.NoError(t, err)
	
	retrievedChat, exists := store.GetChat(chat.ID)
	assert.True(t, exists)
	assert.Equal(t, chat.ID, retrievedChat.ID)
	assert.Equal(t, chat.Title, retrievedChat.Title)
	
	// Test UpdateChat
	updates := map[string]interface{}{
		"title": "Updated Title",
		"file":  "updated.csv",
	}
	
	updatedChat, err := store.UpdateChat(chat.ID, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedChat.Title)
	assert.Equal(t, "updated.csv", updatedChat.File)
	
	// Test ListChats
	chats := store.ListChats()
	assert.Len(t, chats, 1)
	assert.Equal(t, chat.ID, chats[0].ID)
	
	// Test DeleteChat
	err = store.DeleteChat(chat.ID)
	assert.NoError(t, err)
	
	_, exists = store.GetChat(chat.ID)
	assert.False(t, exists)
	
	chats = store.ListChats()
	assert.Len(t, chats, 0)
}

func TestChatPersistence(t *testing.T) {
	// Create a temporary directory for chat data
	tempDir, err := os.MkdirTemp("", "chat-persist-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// First store instance
	store1 := NewChatStore(tempDir)
	
	// Create a test chat
	chat := Chat{
		ID:        "chat-persist",
		Title:     "Persistence Test",
		File:      "persist.csv",
		CreatedAt: time.Now(),
		Messages: []ChatMessage{
			{
				Type:      "user",
				Text:      "Test message",
				Sender:    "user",
				Timestamp: time.Now(),
			},
		},
	}
	
	// Add to store
	err = store1.AddChat(chat)
	require.NoError(t, err)
	
	// Create a new store instance with the same directory
	store2 := NewChatStore(tempDir)
	
	// Verify the chat is loaded
	retrievedChat, exists := store2.GetChat(chat.ID)
	assert.True(t, exists)
	assert.Equal(t, chat.ID, retrievedChat.ID)
	assert.Equal(t, chat.Title, retrievedChat.Title)
	assert.Len(t, retrievedChat.Messages, 1)
	assert.Equal(t, chat.Messages[0].Text, retrievedChat.Messages[0].Text)
}

func TestChatTrainingData(t *testing.T) {
	router, store, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create a test chat
	testChat := Chat{
		ID:             "test-chat-training",
		Title:          "Training Data Test",
		File:           "test.csv",
		CreatedAt:      time.Now(),
		TrainingDataIDs: []string{"training-1", "training-2"},
	}
	
	// Add to store
	err := store.AddChat(testChat)
	require.NoError(t, err)
	
	// Test GET training data
	req, err := http.NewRequest("GET", "/api/chats/test-chat-training/training", nil)
	require.NoError(t, err)
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var getResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &getResponse)
	require.NoError(t, err)
	
	trainingDataIDs, ok := getResponse["trainingDataIds"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, trainingDataIDs, 2)
	
	// Test POST to update training data
	updateData := map[string]interface{}{
		"trainingDataIds": []string{"training-2", "training-3", "training-4"},
	}
	
	body, err := json.Marshal(updateData)
	require.NoError(t, err)
	
	req, err = http.NewRequest("POST", "/api/chats/test-chat-training/training", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify the update
	updatedChat, exists := store.GetChat(testChat.ID)
	assert.True(t, exists)
	assert.Len(t, updatedChat.TrainingDataIDs, 3)
	assert.Contains(t, updatedChat.TrainingDataIDs, "training-2")
	assert.Contains(t, updatedChat.TrainingDataIDs, "training-3")
	assert.Contains(t, updatedChat.TrainingDataIDs, "training-4")
	assert.NotContains(t, updatedChat.TrainingDataIDs, "training-1")
}

func TestChatNotFound(t *testing.T) {
	router, _, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Test GET non-existent chat
	req, err := http.NewRequest("GET", "/api/chats/non-existent", nil)
	require.NoError(t, err)
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusNotFound, rr.Code)
	
	// Test PUT non-existent chat
	updateData := map[string]interface{}{
		"title": "Updated Title",
	}
	
	body, err := json.Marshal(updateData)
	require.NoError(t, err)
	
	req, err = http.NewRequest("PUT", "/api/chats/non-existent", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusNotFound, rr.Code)
	
	// Test DELETE non-existent chat - should not error
	req, err = http.NewRequest("DELETE", "/api/chats/non-existent", nil)
	require.NoError(t, err)
	
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Some implementations return 204 even if the chat doesn't exist
	assert.True(t, rr.Code == http.StatusNoContent || rr.Code == http.StatusNotFound)
}

func TestCreateChatWithMessages(t *testing.T) {
	router, _, cleanup := setupChatRouter(t)
	defer cleanup()
	
	// Create chat with messages
	reqBody := map[string]interface{}{
		"title": "Chat With Messages",
		"file": "test.csv",
		"messages": []map[string]interface{}{
			{
				"type": "system",
				"text": "Chat initialized",
				"timestamp": time.Now(),
			},
			{
				"type": "user",
				"text": "Hello",
				"sender": "user",
				"timestamp": time.Now(),
			},
		},
	}
	
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response Chat
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response.Messages, 2)
	assert.Equal(t, "system", response.Messages[0].Type)
	assert.Equal(t, "user", response.Messages[1].Type)
}