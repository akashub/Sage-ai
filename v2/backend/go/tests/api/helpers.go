// backend/go/tests/api/handlers_test.go
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sage-ai-v2/internal/llm"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueryHandlerDirectly tests the query handler directly with HTTP requests
func TestQueryHandlerDirectly(t *testing.T) {
	// Create a test request with a simple query
	requestBody := map[string]interface{}{
		"query":            "show me all users",
		"csvPath":          "test.csv",
		"useKnowledgeBase": true,
	}
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Create a request
	req, err := http.NewRequest("POST", "/api/query", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create a test server with our mock handlers
	server := NewTestServer(t)
	defer server.Close()

	// Manually route the request to the test server
	server.Router.ServeHTTP(rr, req)

	// Check response status
	assert.Equal(t, http.StatusOK, rr.Code, "Handler should return 200 OK")

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Validate response content - it should contain fields like sql, results, etc.
	assert.Contains(t, response, "sql", "Response should contain SQL")
	assert.Contains(t, response, "results", "Response should contain results")

	// Additional checks based on the mock's expected behavior
	results, ok := response["results"].([]interface{})
	assert.True(t, ok, "Results should be an array")
	assert.NotEmpty(t, results, "Results should not be empty")
}

// TestUploadHandler tests the file upload endpoint
func TestUploadHandler(t *testing.T) {
	// Create a test CSV content
	csvContent := "id,name,email\n1,John,john@example.com\n2,Jane,jane@example.com"

	// Create a multipart form
	body := new(bytes.Buffer)
	writer := NewTestMultipartWriter(body)

	// Add CSV file to form
	err := writer.AddFile("file", "test.csv", csvContent)
	require.NoError(t, err)

	// Close the writer
	err = writer.Close()
	require.NoError(t, err)

	// Create request
	req, err := http.NewRequest("POST", "/api/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.GetContentType())

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create a test server
	server := NewTestServer(t)
	defer server.Close()

	// Route request to handler
	server.Router.ServeHTTP(rr, req)

	// Check for successful response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d: %s", rr.Code, rr.Body.String())
	}

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check if upload was successful (response format depends on your implementation)
	success, ok := response["success"].(bool)
	assert.True(t, ok, "Response should have a success field")
	assert.True(t, success, "Upload should be successful")
}

// TestValidateAPIKeyHandler tests the API key validation endpoint
func TestValidateAPIKeyHandler(t *testing.T) {
	// Test cases for different providers
	testCases := []struct {
		name           string
		provider       llm.LLMProvider
		apiKey         string
		expectedStatus int
		expectedValid  bool
	}{
		{
			name:           "Gemini Provider",
			provider:       llm.ProviderGemini,
			apiKey:         "test-gemini-key",
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "OpenAI Provider",
			provider:       llm.ProviderOpenAI,
			apiKey:         "test-openai-key",
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "Anthropic Provider",
			provider:       llm.ProviderAnthropic,
			apiKey:         "test-anthropic-key",
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "Mistral Provider",
			provider:       llm.ProviderMistral,
			apiKey:         "test-mistral-key",
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "Empty API Key",
			provider:       llm.ProviderGemini,
			apiKey:         "",
			expectedStatus: http.StatusBadRequest,
			expectedValid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server
			server := NewTestServer(t)
			defer server.Close()

			// Create a request
			reqBody := llm.LLMConfig{
				Provider: tc.provider,
				APIKey:   tc.apiKey,
			}

			reqBodyBytes, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/validate-api-key", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Send the request through the router
			server.Router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, rr.Code, "Handler returned wrong status code")

			// For successful responses, check the "valid" field
			if tc.expectedStatus == http.StatusOK {
				var response map[string]bool
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedValid, response["valid"], "Unexpected 'valid' value in response")
			}
		})
	}
}

// TestQueryWithLLMConfig tests query handling with different LLM configurations
func TestQueryWithLLMConfig(t *testing.T) {
	// Test different provider configs
	testCases := []struct {
		name     string
		provider llm.LLMProvider
		apiKey   string
	}{
		{"Gemini Provider", llm.ProviderGemini, "gemini-key"},
		{"OpenAI Provider", llm.ProviderOpenAI, "openai-key"},
		{"Anthropic Provider", llm.ProviderAnthropic, "anthropic-key"},
		{"Mistral Provider", llm.ProviderMistral, "mistral-key"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server
			server := NewTestServer(t)
			defer server.Close()

			// Create a query request with LLM config
			reqBody := map[string]interface{}{
				"query":            "What is the weather?",
				"csvPath":          "weather.csv",
				"useKnowledgeBase": true,
				"options": map[string]interface{}{
					"llmConfig": map[string]interface{}{
						"provider": string(tc.provider),
						"api_key":  tc.apiKey,
					},
				},
			}

			reqBodyBytes, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/query", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Send the request through the router
			server.Router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, http.StatusOK, rr.Code, "Handler should return 200 OK")

			// Parse response
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify SQL field exists
			assert.Contains(t, response, "sql", "Response should contain SQL")

			// Verify results field exists and is non-empty
			results, ok := response["results"].([]interface{})
			assert.True(t, ok, "Results should be an array")
			assert.NotEmpty(t, results, "Results should not be empty")

			// Verify the mock orchestrator received and used the LLM config
			// This is an indirect test - checking that the result contains expected content
			// that would only be present if the config was correctly passed

			// The actual check depends on how your mock orchestrator behaves with LLM configs
			// You could look for specific patterns in the SQL or results that indicate
			// the LLM config was used
		})
	}
}

// Helper for checking string map contains
func mapContainsSubstring(m map[string]interface{}, key, substr string) bool {
	if value, ok := m[key].(string); ok {
		return strings.Contains(value, substr)
	}
	return false
}
