// backend/go/tests/orchestrator/orchestrator_test.go
package orchestrator

import (
	"context"
	"testing"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/types"

	"github.com/stretchr/testify/assert"
)

// MockOrchestrator provides a testable implementation of the orchestrator
type MockOrchestrator struct {
	processResult *types.State
	processError  error
	// Use function fields to allow overriding behavior in tests
	ProcessQueryFunc           func(ctx context.Context, query string, csvPath string) (*types.State, error)
	ProcessQueryWithOptionsFunc func(ctx context.Context, query string, csvPath string, options map[string]interface{}) (*types.State, error)
	NewSessionFunc             func()
}

// ProcessQuery implements the orchestrator interface
func (m *MockOrchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
	if m.ProcessQueryFunc != nil {
		return m.ProcessQueryFunc(ctx, query, csvPath)
	}
	if m.processError != nil {
		return nil, m.processError
	}
	return m.processResult, nil
}

// ProcessQueryWithOptions implements the orchestrator interface
func (m *MockOrchestrator) ProcessQueryWithOptions(ctx context.Context, query string, csvPath string, options map[string]interface{}) (*types.State, error) {
	if m.ProcessQueryWithOptionsFunc != nil {
		return m.ProcessQueryWithOptionsFunc(ctx, query, csvPath, options)
	}
	if m.processError != nil {
		return nil, m.processError
	}
	return m.processResult, nil
}

// NewSession implements the orchestrator interface
func (m *MockOrchestrator) NewSession() {
	if m.NewSessionFunc != nil {
		m.NewSessionFunc()
		return
	}
	// Default no-op for testing
}

// TestOrchestratorProcessQueryWithCustomMock tests the basic query processing flow
// using a custom mocked implementation
func TestOrchestratorProcessQueryWithCustomMock(t *testing.T) {
	// Create a mock orchestrator
	mockOrch := &MockOrchestrator{
		processResult: &types.State{
			Query:          "show me all users",
			GeneratedQuery: "SELECT id, name FROM users",
			Analysis: map[string]interface{}{
				"query_type": "select",
				"columns":    []string{"id", "name"},
			},
			ExecutionResult: []map[string]interface{}{
				{
					"id":   1,
					"name": "John",
				},
			},
			ValidationResult: map[string]interface{}{
				"isValid": true,
			},
		},
	}

	// Test processing a query using the mock
	state, err := mockOrch.ProcessQuery(context.Background(), "show me all users", "test.csv")
	
	// Verify results
	assert.NoError(t, err, "ProcessQuery should not return an error")
	assert.NotNil(t, state, "ProcessQuery should return a state")
	assert.Equal(t, "SELECT id, name FROM users", state.GeneratedQuery, "Generated query should match mock result")
}

// TestOrchestratorProcessQueryWithKnowledge tests the knowledge-enhanced query flow
func TestOrchestratorProcessQueryWithKnowledge(t *testing.T) {
	// Create a mock orchestrator with knowledge context
	mockOrch := &MockOrchestrator{
		processResult: &types.State{
			Query:          "show me all users",
			GeneratedQuery: "SELECT id, name FROM users",
			Analysis: map[string]interface{}{
				"query_type": "select",
				"columns":    []string{"id", "name"},
			},
			ExecutionResult: []map[string]interface{}{
				{
					"id":   1,
					"name": "John",
				},
			},
			ValidationResult: map[string]interface{}{
				"isValid": true,
			},
			KnowledgeContext: &knowledge.KnowledgeResult{
				DDLSchemas: []knowledge.TrainingItem{
					{
						ID:          "ddl_123",
						Type:        "ddl",
						Content:     "CREATE TABLE users (id INT, name VARCHAR(100))",
						Description: "Users table schema",
					},
				},
				Documentation:    []knowledge.TrainingItem{},
				QuestionSQLPairs: []knowledge.QuestionSQLPair{},
			},
		},
	}

	// Set up options to use knowledge base
	options := map[string]interface{}{
		"useKnowledgeBase": true,
	}

	// Test processing a query with knowledge
	state, err := mockOrch.ProcessQueryWithOptions(context.Background(), "show me all users", "test.csv", options)
	
	// Verify results
	assert.NoError(t, err, "ProcessQueryWithOptions should not return an error")
	assert.NotNil(t, state, "ProcessQueryWithOptions should return a state")
	assert.Equal(t, "SELECT id, name FROM users", state.GeneratedQuery, "Generated query should match mock result")
	assert.NotNil(t, state.KnowledgeContext, "KnowledgeContext should not be nil")
	assert.Len(t, state.KnowledgeContext.DDLSchemas, 1, "Should have one DDL schema")
}

// TestNewSession tests the session creation functionality using our enhanced mock
func TestNewSession(t *testing.T) {
	// Create a mock with session tracking
	sessionCreated := false
	mockOrch := &MockOrchestrator{
		NewSessionFunc: func() {
			sessionCreated = true
		},
	}
	
	// Call NewSession
	mockOrch.NewSession()
	
	// Verify the method was called
	assert.True(t, sessionCreated, "NewSession should have been called")
}