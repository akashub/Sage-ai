// backend/go/internal/api/handlers/training_data_test.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/orchestrator"
)

// MockKnowledgeManager implements the methods needed for testing
type MockKnowledgeManager struct {
	mock.Mock
	knowledge.KnowledgeManager
}

func (m *MockKnowledgeManager) AddDDLSchema(ctx context.Context, name string, content string, description string) error {
	args := m.Called(ctx, name, content, description)
	return args.Error(0)
}

func (m *MockKnowledgeManager) AddDocumentation(ctx context.Context, title string, content string, tags []string) error {
	args := m.Called(ctx, title, content, tags)
	return args.Error(0)
}

func (m *MockKnowledgeManager) AddQuestionSQLPair(ctx context.Context, pair knowledge.QuestionSQLPair) error {
	args := m.Called(ctx, pair)
	return args.Error(0)
}

func (m *MockKnowledgeManager) ListTrainingData(ctx context.Context, dataType string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, dataType)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockKnowledgeManager) GetTrainingItem(ctx context.Context, id string) (*knowledge.TrainingItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*knowledge.TrainingItem), args.Error(1)
}

func (m *MockKnowledgeManager) DeleteTrainingItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockKnowledgeManager) StoreFile(fileName string, content []byte) (string, error) {
	args := m.Called(fileName, content)
	return args.String(0), args.Error(1)
}

func setupMockOrchestrator() (*orchestrator.Orchestrator, *MockKnowledgeManager) {
	mockKM := new(MockKnowledgeManager)
	orch := &orchestrator.Orchestrator{}
	return orch, mockKM
}

func TestAddTrainingDataHandler(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	// Setup expectations
	mockKM.On("AddDDLSchema", 
		mock.Anything, 
		"Test Schema",
		"CREATE TABLE test (id INT)",
		"Test schema description").Return(nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/add", AddTrainingDataHandler(orch)).Methods("POST", "OPTIONS")
	
	// Create request
	reqBody := map[string]interface{}{
		"type":        "ddl",
		"content":     "CREATE TABLE test (id INT)",
		"description": "Test schema description",
		"metadata": map[string]interface{}{
			"name": "Test Schema",
		},
	}
	
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/training/add", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, true, response["success"])
	assert.Contains(t, response, "trainingID")
	
	// Verify mock was called
	mockKM.AssertExpectations(t)
}

func TestListTrainingDataHandler(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	// Mock data
	trainingData := []map[string]interface{}{
		{
			"id":          "ddl-1",
			"type":        "ddl",
			"description": "Test DDL",
			"date_added":  time.Now().Format(time.RFC3339),
		},
		{
			"id":          "doc-1",
			"type":        "documentation",
			"description": "Test Doc",
			"date_added":  time.Now().Format(time.RFC3339),
		},
	}
	
	// Setup expectations
	mockKM.On("ListTrainingData", mock.Anything, "").Return(trainingData, nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/list", ListTrainingDataHandler(orch)).Methods("GET", "OPTIONS")
	
	// Create request
	req, err := http.NewRequest("GET", "/api/training/list", nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify response
	var response []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response, 2)
	assert.Equal(t, "ddl-1", response[0]["id"])
	assert.Equal(t, "doc-1", response[1]["id"])
	
	// Verify mock was called
	mockKM.AssertExpectations(t)
}

func TestListTrainingDataHandlerWithFilter(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	// Mock data
	trainingData := []map[string]interface{}{
		{
			"id":          "ddl-1",
			"type":        "ddl",
			"description": "Test DDL",
			"date_added":  time.Now().Format(time.RFC3339),
		},
	}
	
	// Setup expectations
	mockKM.On("ListTrainingData", mock.Anything, "ddl").Return(trainingData, nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/list", ListTrainingDataHandler(orch)).Methods("GET", "OPTIONS")
	
	// Create request
	req, err := http.NewRequest("GET", "/api/training/list?type=ddl", nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify response
	var response []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response, 1)
	assert.Equal(t, "ddl-1", response[0]["id"])
	assert.Equal(t, "ddl", response[0]["type"])
	
	// Verify mock was called
	mockKM.AssertExpectations(t)
}

func TestViewTrainingDataHandler(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	itemID := "ddl-1"
	
	// Mock data
	trainingItem := &knowledge.TrainingItem{
		ID:          itemID,
		Type:        "ddl",
		Content:     "CREATE TABLE test (id INT)",
		Description: "Test DDL",
		DateAdded:   time.Now().Format(time.RFC3339),
	}
	
	// Setup expectations
	mockKM.On("GetTrainingItem", mock.Anything, itemID).Return(trainingItem, nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/view/{id}", ViewTrainingDataHandler(orch)).Methods("GET", "OPTIONS")
	
	// Create request
	req, err := http.NewRequest("GET", "/api/training/view/"+itemID, nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify response
	var response knowledge.TrainingItem
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, itemID, response.ID)
	assert.Equal(t, "ddl", response.Type)
	assert.Equal(t, "CREATE TABLE test (id INT)", response.Content)
	assert.Equal(t, "Test DDL", response.Description)
	
	// Verify mock was called
	mockKM.AssertExpectations(t)
}

func TestDeleteTrainingDataHandler(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	itemID := "ddl-1"
	
	// Setup expectations
	mockKM.On("DeleteTrainingItem", mock.Anything, itemID).Return(nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/delete/{id}", DeleteTrainingDataHandler(orch)).Methods("DELETE", "OPTIONS")
	
	// Create request
	req, err := http.NewRequest("DELETE", "/api/training/delete/"+itemID, nil)
	require.NoError(t, err)
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response - should be 204 No Content
	assert.Equal(t, http.StatusNoContent, rr.Code)
	
	// Verify mock was called
	mockKM.AssertExpectations(t)
}

func TestUploadTrainingFileHandler(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	// Setup expectations for file storing
	mockKM.On("StoreFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return("/path/to/stored/file.sql", nil)
	
	// Setup expectations for adding DDL
	mockKM.On("AddDDLSchema", 
		mock.Anything, 
		"Test Schema",
		"CREATE TABLE test (id INT)",
		"Test schema").Return(nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/upload", UploadTrainingFileHandler(orch)).Methods("POST", "OPTIONS")
	
	// Create multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// Add file
	fileWriter, err := w.CreateFormFile("file", "test_schema.sql")
	require.NoError(t, err)
	_, err = io.Copy(fileWriter, strings.NewReader("CREATE TABLE test (id INT)"))
	require.NoError(t, err)
	
	// Add form fields
	err = w.WriteField("type", "ddl")
	require.NoError(t, err)
	err = w.WriteField("description", "Test schema")
	require.NoError(t, err)
	
	// Close writer
	w.Close()
	
	// Create request
	req, err := http.NewRequest("POST", "/api/training/upload", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, true, response["success"])
	assert.Contains(t, response, "trainingID")
	
	// Verify mocks were called
	mockKM.AssertExpectations(t)
}

func TestUploadTrainingFileHandlerWithAutoDetect(t *testing.T) {
	orch, mockKM := setupMockOrchestrator()
	
	// Setup expectations for file storing
	mockKM.On("StoreFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return("/path/to/stored/file.sql", nil)
	
	// Setup expectations for adding DDL
	mockKM.On("AddDDLSchema", 
		mock.Anything, 
		"auto_detect_test",
		"CREATE TABLE test (id INT)",
		"Auto-detect test").Return(nil)
	
	// Create router with handler
	router := mux.NewRouter()
	router.HandleFunc("/api/training/upload", UploadTrainingFileHandler(orch)).Methods("POST", "OPTIONS")
	
	// Create multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// Add file
	fileWriter, err := w.CreateFormFile("file", "auto_detect_test.sql")
	require.NoError(t, err)
	_, err = io.Copy(fileWriter, strings.NewReader("CREATE TABLE test (id INT)"))
	require.NoError(t, err)
	
	// Add form fields
	err = w.WriteField("type", "auto")
	require.NoError(t, err)
	err = w.WriteField("description", "Auto-detect test")
	require.NoError(t, err)
	
	// Close writer
	w.Close()
	
	// Create request
	req, err := http.NewRequest("POST", "/api/training/upload", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	
	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify mocks were called
	mockKM.AssertExpectations(t)
}
