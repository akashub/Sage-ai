// backend/go/internal/orchestrator/orchestrator_test.go
package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/llm"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockLLMBridge is a mock implementation of llm.Bridge for testing
type MockLLMBridge struct {
	mock.Mock
	llm.Bridge
}

func (m *MockLLMBridge) SetSession(sessionID string) {
	m.Called(sessionID)
}

func (m *MockLLMBridge) Analyze(ctx context.Context, question string, schema map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, question, schema)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockLLMBridge) AnalyzeWithKnowledge(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockLLMBridge) GenerateQuery(ctx context.Context, analysis map[string]interface{}, schema map[string]interface{}) (string, error) {
	args := m.Called(ctx, analysis, schema)
	return args.String(0), args.Error(1)
}

func (m *MockLLMBridge) GenerateQueryWithKnowledge(ctx context.Context, req map[string]interface{}) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockLLMBridge) ValidateQuery(ctx context.Context, query string, schema map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, query, schema)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockLLMBridge) HealQuery(ctx context.Context, validationResult map[string]interface{}, originalQuery string, analysis map[string]interface{}, schema map[string]interface{}) (*llm.HealingResult, error) {
	args := m.Called(ctx, validationResult, originalQuery, analysis, schema)
	return args.Get(0).(*llm.HealingResult), args.Error(1)
}

func setupTestOrchestrator(t *testing.T) (*Orchestrator, *MockLLMBridge, *knowledge.KnowledgeManager, func()) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "orchestrator-test")
	require.NoError(t, err)

	// Create a mock LLM bridge
	mockBridge := new(MockLLMBridge)
	
	// Set up session ID expectations
	mockBridge.On("SetSession", mock.AnythingOfType("string")).Return()
	
	// Create a knowledge manager for testing
	vectorDB, err := knowledge.CreateMemoryVectorDB(filepath.Join(tempDir, "vector-store.json"))
	require.NoError(t, err)
	
	km := knowledge.CreateKnowledgeManager(vectorDB, tempDir)
	require.NotNil(t, km)
	
	// Create orchestrator with mock bridge
	orch := CreateOrchestrator(&mockBridge.Bridge, km)
	require.NotNil(t, orch)
	
	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	return orch, mockBridge, km, cleanup
}

func TestCreateOrchestrator(t *testing.T) {
	orch, _, _, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	assert.NotNil(t, orch)
	assert.NotNil(t, orch.bridge)
	assert.NotNil(t, orch.KnowledgeManager)
	assert.NotNil(t, orch.graph)
	assert.NotEmpty(t, orch.sessionID)
}

func TestNewSession(t *testing.T) {
	orch, mockBridge, _, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	// Initial session ID
	initialSessionID := orch.sessionID
	
	// Create a new session
	orch.NewSession()
	
	// Should have a different session ID
	assert.NotEqual(t, initialSessionID, orch.sessionID)
	
	// Mock should have been called again with the new session ID
	mockBridge.AssertCalled(t, "SetSession", orch.sessionID)
}

func TestClearSession(t *testing.T) {
	orch, _, _, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	// Clear the session
	orch.ClearSession()
	
	// Session ID should be empty
	assert.Empty(t, orch.sessionID)
	
	// Graph should be nil
	assert.Nil(t, orch.graph)
}

func TestProcessQueryWithOptions_WithoutKnowledgeBase(t *testing.T) {
	orch, mockBridge, _, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	ctx := context.Background()
	query := "Show me all users"
	csvPath := "users.csv"
	
	// Mock schema extraction (handled by nodes)
	// Mock analysis response
	mockAnalysis := map[string]interface{}{
		"query_type": "select",
		"tables":     []string{"users"},
		"columns":    []string{"name", "age"},
	}
	
	// Mock query generation
	generatedQuery := "SELECT name, age FROM users"
	
	// Mock validation response
	validationResult := map[string]interface{}{
		"isValid": true,
		"issues":  []string{},
	}
	
	// Set up expectations
	mockBridge.On("Analyze", mock.Anything, query, mock.AnythingOfType("map[string]interface {}")).Return(mockAnalysis, nil)
	mockBridge.On("GenerateQuery", mock.Anything, mockAnalysis, mock.AnythingOfType("map[string]interface {}")).Return(generatedQuery, nil)
	mockBridge.On("ValidateQuery", mock.Anything, generatedQuery, mock.AnythingOfType("map[string]interface {}")).Return(validationResult, nil)
	
	// Process query without knowledge base
	options := map[string]interface{}{
		"useKnowledgeBase": false,
	}
	
	result, err := orch.ProcessQueryWithOptions(ctx, query, csvPath, options)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Check result
	assert.Equal(t, generatedQuery, result.GeneratedQuery)
	
	// Verify expectations
	mockBridge.AssertCalled(t, "Analyze", mock.Anything, query, mock.AnythingOfType("map[string]interface {}"))
	mockBridge.AssertCalled(t, "GenerateQuery", mock.Anything, mockAnalysis, mock.AnythingOfType("map[string]interface {}"))
	mockBridge.AssertCalled(t, "ValidateQuery", mock.Anything, generatedQuery, mock.AnythingOfType("map[string]interface {}"))
	
	// Should NOT have called AnalyzeWithKnowledge or GenerateQueryWithKnowledge
	mockBridge.AssertNotCalled(t, "AnalyzeWithKnowledge", mock.Anything, mock.Anything)
	mockBridge.AssertNotCalled(t, "GenerateQueryWithKnowledge", mock.Anything, mock.Anything)
}

func TestProcessQueryWithOptions_WithKnowledgeBase(t *testing.T) {
	orch, mockBridge, km, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	ctx := context.Background()
	query := "Show me all users"
	csvPath := "users.csv"
	
	// Add some test knowledge
	err := km.AddDDLSchema(ctx, "users_schema", 
		"CREATE TABLE users (id INT, name VARCHAR(100), age INT)", 
		"Users table schema")
	assert.NoError(t, err)
	
	// Mock schema extraction (handled by nodes)
	// Mock analysis response
	mockAnalysis := map[string]interface{}{
		"query_type": "select",
		"tables":     []string{"users"},
		"columns":    []string{"name", "age"},
	}
	
	// Mock query generation
	generatedQuery := "SELECT name, age FROM users"
	
	// Mock validation response
	validationResult := map[string]interface{}{
		"isValid": true,
		"issues":  []string{},
	}
	
	// Set up expectations for knowledge-based analysis
	mockBridge.On("AnalyzeWithKnowledge", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(mockAnalysis, nil)
	mockBridge.On("GenerateQueryWithKnowledge", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(generatedQuery, nil)
	mockBridge.On("ValidateQuery", mock.Anything, generatedQuery, mock.AnythingOfType("map[string]interface {}")).Return(validationResult, nil)
	
	// Process query with knowledge base
	options := map[string]interface{}{
		"useKnowledgeBase": true,
	}
	
	result, err := orch.ProcessQueryWithOptions(ctx, query, csvPath, options)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Check result
	assert.Equal(t, generatedQuery, result.GeneratedQuery)
	
	// Verify knowledge context
	assert.NotNil(t, result.KnowledgeContext)
	
	// Should have called knowledge-based methods
	mockBridge.AssertCalled(t, "AnalyzeWithKnowledge", mock.Anything, mock.Anything)
	mockBridge.AssertCalled(t, "GenerateQueryWithKnowledge", mock.Anything, mock.Anything)
	
	// Should NOT have called regular methods
	mockBridge.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything, mock.Anything)
	mockBridge.AssertNotCalled(t, "GenerateQuery", mock.Anything, mock.Anything, mock.Anything)
}

// TestProcessQueryWithOptions_ErrorHandling tests error handling in the process pipeline
func TestProcessQueryWithOptions_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockLLMBridge)
		expectedError bool
	}{
		{
			name: "Analysis Error",
			setupMocks: func(mockBridge *MockLLMBridge) {
				mockBridge.On("Analyze", mock.Anything, mock.Anything, mock.Anything).
					Return(map[string]interface{}{}, fmt.Errorf("analysis error"))
			},
			expectedError: true,
		},
		{
			name: "Query Generation Error",
			setupMocks: func(mockBridge *MockLLMBridge) {
				mockBridge.On("Analyze", mock.Anything, mock.Anything, mock.Anything).
					Return(map[string]interface{}{"query_type": "select"}, nil)
				mockBridge.On("GenerateQuery", mock.Anything, mock.Anything, mock.Anything).
					Return("", fmt.Errorf("generation error"))
			},
			expectedError: true,
		},
		{
			name: "Validation Error",
			setupMocks: func(mockBridge *MockLLMBridge) {
				mockBridge.On("Analyze", mock.Anything, mock.Anything, mock.Anything).
					Return(map[string]interface{}{"query_type": "select"}, nil)
				mockBridge.On("GenerateQuery", mock.Anything, mock.Anything, mock.Anything).
					Return("SELECT * FROM users", nil)
				mockBridge.On("ValidateQuery", mock.Anything, mock.Anything, mock.Anything).
					Return(map[string]interface{}{}, fmt.Errorf("validation error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orch, mockBridge, _, cleanup := setupTestOrchestrator(t)
			defer cleanup()
			
			// Setup mocks for this test case
			tt.setupMocks(mockBridge)
			
			ctx := context.Background()
			query := "Show me all users"
			csvPath := "users.csv"
			options := map[string]interface{}{"useKnowledgeBase": false}
			
			// Process query
			_, err := orch.ProcessQueryWithOptions(ctx, query, csvPath, options)
			
			// Check if error occurred as expected
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessQuery_DefaultOptions(t *testing.T) {
	orch, mockBridge, _, cleanup := setupTestOrchestrator(t)
	defer cleanup()
	
	ctx := context.Background()
	query := "Show me all users"
	csvPath := "users.csv"
	
	// Mock analysis response
	mockAnalysis := map[string]interface{}{
		"query_type": "select",
		"tables":     []string{"users"},
		"columns":    []string{"*"},
	}
	
	// Mock query generation
	generatedQuery := "SELECT * FROM users"
	
	// Mock validation response
	validationResult := map[string]interface{}{
		"isValid": true,
		"issues":  []string{},
	}
	
	// Set up expectations - should use non-knowledge methods by default
	mockBridge.On("Analyze", mock.Anything, query, mock.AnythingOfType("map[string]interface {}")).Return(mockAnalysis, nil)
	mockBridge.On("GenerateQuery", mock.Anything, mockAnalysis, mock.AnythingOfType("map[string]interface {}")).Return(generatedQuery, nil)
	mockBridge.On("ValidateQuery", mock.Anything, generatedQuery, mock.AnythingOfType("map[string]interface {}")).Return(validationResult, nil)
	
	// Process query using ProcessQuery which should use default options
	result, err := orch.ProcessQuery(ctx, query, csvPath)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Check result
	assert.Equal(t, generatedQuery, result.GeneratedQuery)
	
	// Verify standard methods were called
	mockBridge.AssertCalled(t, "Analyze", mock.Anything, query, mock.AnythingOfType("map[string]interface {}"))
	mockBridge.AssertCalled(t, "GenerateQuery", mock.Anything, mockAnalysis, mock.AnythingOfType("map[string]interface {}"))
	mockBridge.AssertCalled(t, "ValidateQuery", mock.Anything, generatedQuery, mock.AnythingOfType("map[string]interface {}"))
}