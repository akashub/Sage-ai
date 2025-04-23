// backend/go/tests/llm/bridge_test.go
package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sage-ai-v2/internal/llm"
	"testing"
	"time"
)

func TestCreateBridge(t *testing.T) {
	baseURL := "http://example.com"
	bridge := llm.CreateBridge(baseURL)
	
	if bridge.GetSessionID() != "" {
		t.Errorf("Expected empty session ID initially, got %s", bridge.GetSessionID())
	}
}

func TestSetSession(t *testing.T) {
	bridge := llm.CreateBridge("http://example.com")
	sessionID := "test-session-123"
	
	bridge.SetSession(sessionID)
	
	if bridge.GetSessionID() != sessionID {
		t.Errorf("Expected session ID to be %s, got %s", sessionID, bridge.GetSessionID())
	}
}

func TestSetLLMConfig(t *testing.T) {
	bridge := llm.CreateBridge("http://example.com")
	config := &llm.LLMConfig{
		Provider: llm.ProviderOpenAI,
		APIKey:   "test-api-key",
	}
	
	bridge.SetLLMConfig(config)
	
	// Since llmConfig is a private field, we can only test indirectly
	// by checking that subsequent requests include the config
	// This would require integration with a test server
	// For now, we'll just check that the method doesn't panic
}

func TestMakeRequestWithLLMConfig(t *testing.T) {
	// Create a test server that captures the request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&requestBody)
		
		// Check if llm_config is present in the request
		if config, ok := requestBody["llm_config"].(map[string]interface{}); ok {
			// Respond with success and echo back the config
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"config_received": config,
			})
		} else {
			// Respond with error
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "No LLM config provided",
			})
		}
	}))
	defer server.Close()
	
	bridge := llm.CreateBridge(server.URL)
	bridge.SetSession("test-session")
	
	// Test different provider configurations
	testCases := []struct {
		name     string
		provider llm.LLMProvider
		apiKey   string
	}{
		{"Gemini Provider", llm.ProviderGemini, "gemini-api-key"},
		{"OpenAI Provider", llm.ProviderOpenAI, "openai-api-key"},
		{"Anthropic Provider", llm.ProviderAnthropic, "anthropic-api-key"},
		{"Mistral Provider", llm.ProviderMistral, "mistral-api-key"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &llm.LLMConfig{
				Provider: tc.provider,
				APIKey:   tc.apiKey,
			}
			
			bridge.SetLLMConfig(config)
			
			// Make a test request to check if config is passed correctly
			requestData := map[string]interface{}{
				"test": "data",
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			resp, err := bridge.MakeRequest(ctx, "/test", requestData, "test-session")
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			
			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(resp, &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}
			
			// Check if config was received correctly
			if _, ok := response["config_received"]; !ok {
				t.Errorf("Expected config_received in response, not found")
			}
		})
	}
}

func TestAnalyzeWithKnowledge(t *testing.T) {
	// Create a test server that responds with a successful analysis
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure we're hitting the analyze_with_knowledge endpoint
		if r.URL.Path != "/analyze_with_knowledge" {
			t.Errorf("Expected /analyze_with_knowledge endpoint, got %s", r.URL.Path)
		}
		
		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"analysis": map[string]interface{}{
				"intent": "query",
				"query_type": "select",
				"tables": []string{"users"},
				"knowledge_used": true,
			},
		})
	}))
	defer server.Close()
	
	bridge := llm.CreateBridge(server.URL)
	bridge.SetSession("test-session")
	
	// Configure with OpenAI
	bridge.SetLLMConfig(&llm.LLMConfig{
		Provider: llm.ProviderOpenAI,
		APIKey:   "test-key",
	})
	
	// Prepare a test request with knowledge context
	analysisRequest := map[string]interface{}{
		"query": "Find all users",
		"schema": map[string]interface{}{
			"users": map[string]interface{}{
				"columns": []string{"id", "name", "email"},
			},
		},
		"knowledge_context": map[string]interface{}{
			"ddl_schemas": []interface{}{
				map[string]interface{}{
					"content": "CREATE TABLE users (id INT, name TEXT, email TEXT);",
				},
			},
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Call the method under test
	result, err := bridge.AnalyzeWithKnowledge(ctx, analysisRequest)
	
	if err != nil {
		t.Fatalf("AnalyzeWithKnowledge failed: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	
	// Verify the response content
	if intent, ok := result["intent"].(string); !ok || intent != "query" {
		t.Errorf("Expected intent to be 'query', got %v", result["intent"])
	}
	
	if queryType, ok := result["query_type"].(string); !ok || queryType != "select" {
		t.Errorf("Expected query_type to be 'select', got %v", result["query_type"])
	}
	
	if knowledgeUsed, ok := result["knowledge_used"].(bool); !ok || !knowledgeUsed {
		t.Errorf("Expected knowledge_used to be true, got %v", result["knowledge_used"])
	}
}

func TestGenerateQueryWithKnowledge(t *testing.T) {
	// Create a test server that responds with a successful query generation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure we're hitting the generate_with_knowledge endpoint
		if r.URL.Path != "/generate_with_knowledge" {
			t.Errorf("Expected /generate_with_knowledge endpoint, got %s", r.URL.Path)
		}
		
		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"query": "SELECT id, name, email FROM users",
		})
	}))
	defer server.Close()
	
	bridge := llm.CreateBridge(server.URL)
	bridge.SetSession("test-session")
	
	// Configure with Anthropic
	bridge.SetLLMConfig(&llm.LLMConfig{
		Provider: llm.ProviderAnthropic,
		APIKey:   "test-key",
	})
	
	// Prepare a test request with knowledge context
	generateRequest := map[string]interface{}{
		"analysis": map[string]interface{}{
			"intent": "query",
			"query_type": "select",
			"tables": []string{"users"},
			"columns": []string{"id", "name", "email"},
		},
		"schema": map[string]interface{}{
			"users": map[string]interface{}{
				"columns": []string{"id", "name", "email"},
			},
		},
		"knowledge_context": map[string]interface{}{
			"examples": []interface{}{
				map[string]interface{}{
					"question": "Get all users",
					"sql": "SELECT * FROM users",
				},
			},
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Call the method under test
	query, err := bridge.GenerateQueryWithKnowledge(ctx, generateRequest)
	
	if err != nil {
		t.Fatalf("GenerateQueryWithKnowledge failed: %v", err)
	}
	
	expectedQuery := "SELECT id, name, email FROM users"
	if query != expectedQuery {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery, query)
	}
}