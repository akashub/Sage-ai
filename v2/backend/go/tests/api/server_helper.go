// backend/go/tests/api/server_helper.go
package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sage-ai-v2/internal/api/handlers"
	"sage-ai-v2/internal/types"
	"testing"

	"github.com/gorilla/mux"
)

// NewTestServer creates a test server with mocked handlers
func NewTestServer(t *testing.T) *TestServer {
	// Create router
	router := mux.NewRouter()

	// Create predefined result for the mock query handler
	result := &types.State{
		Query:          "test query",
		GeneratedQuery: "SELECT * FROM test",
		Analysis: map[string]interface{}{
			"response": "Test response",
		},
		ExecutionResult: []map[string]interface{}{
			{
				"id":   1,
				"name": "Test",
			},
		},
	}

	// Register a mock query handler that always returns the same result
	router.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var req struct {
			Query           string                 `json:"query"`
			CSVPath         string                 `json:"csvPath"`
			UseKnowledgeBase bool                   `json:"useKnowledgeBase"`
			Options         map[string]interface{} `json:"options,omitempty"`
		}
		
		// Just parse the request for validation, but ignore the values
		if err := parsePenNameJSONRequest(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Return a predefined successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Convert the execution result to a response format
		response := map[string]interface{}{
			"success":  true,
			"sql":      result.GeneratedQuery,
			"results":  result.ExecutionResult,
			"response": result.Analysis["response"],
		}
		
		// Write the response
		if err := writeJSON(w, response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("POST", "OPTIONS")

	// Register the real upload handler
	router.HandleFunc("/api/upload", handlers.HandleFileUpload).Methods("POST", "OPTIONS")

	// Create a test HTTP server
	server := httptest.NewServer(router)

	return &TestServer{
		Router: router,
		Server: server,
	}
}

// Close closes the test server
func (s *TestServer) Close() {
	if s.Server != nil {
		s.Server.Close()
	}
}

// TestServer holds a test HTTP server and router
type TestServer struct {
	Router *mux.Router
	Server *httptest.Server
}

// MockQueryHandler returns a mock implementation of the QueryHandler
func MockQueryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return a predefined successful response
		w.Header().Set("Content-Type", "application/json")
		
		response := map[string]interface{}{
			"success": true,
			"sql":     "SELECT * FROM test",
			"results": []map[string]interface{}{
				{
					"id":   1,
					"name": "Test",
				},
			},
			"response": "This is a test response",
		}
		
		writeJSON(w, response)
	}
}

// MockProcessQuery mocks the ProcessQuery method of the Orchestrator
func MockProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
	// Return a predefined result
	return &types.State{
		Query:          query,
		CSVPath:        csvPath,
		GeneratedQuery: "SELECT * FROM test",
		ExecutionResult: []map[string]interface{}{
			{
				"id":   1,
				"name": "Test",
			},
		},
		Analysis: map[string]interface{}{
			"response": "This is a test response",
		},
	}, nil
}