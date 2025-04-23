// backend/go/tests/llm/provider_test.go
package llm

import (
	"encoding/json"
	"sage-ai-v2/internal/llm"
	"testing"
)

func TestLLMProviderConstants(t *testing.T) {
	// Test that provider constants have expected values
	tests := []struct {
		provider llm.LLMProvider
		expected string
	}{
		{llm.ProviderGemini, "gemini"},
		{llm.ProviderOpenAI, "openai"},
		{llm.ProviderAnthropic, "anthropic"},
		{llm.ProviderMistral, "mistral"},
	}

	for _, test := range tests {
		t.Run(string(test.provider), func(t *testing.T) {
			if string(test.provider) != test.expected {
				t.Errorf("Expected %s provider to have value %s, got %s",
					test.provider, test.expected, string(test.provider))
			}
		})
	}
}

func TestLLMConfig(t *testing.T) {
	// Test creating and accessing LLMConfig
	config := llm.LLMConfig{
		Provider: llm.ProviderOpenAI,
		APIKey:   "test-key-123",
	}

	if config.Provider != llm.ProviderOpenAI {
		t.Errorf("Expected Provider to be %s, got %s", llm.ProviderOpenAI, config.Provider)
	}

	if config.APIKey != "test-key-123" {
		t.Errorf("Expected APIKey to be %s, got %s", "test-key-123", config.APIKey)
	}
}

func TestLLMConfigJSONSerialization(t *testing.T) {
	// Test JSON serialization of LLMConfig
	config := llm.LLMConfig{
		Provider: llm.ProviderOpenAI,
		APIKey:   "test-key-123",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal LLMConfig to JSON: %v", err)
	}

	// Check if JSON has expected fields and values
	expectedJSON := `{"provider":"openai","api_key":"test-key-123"}`
	if string(jsonBytes) != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, string(jsonBytes))
	}

	// Test deserialization
	var deserializedConfig llm.LLMConfig
	if err := json.Unmarshal(jsonBytes, &deserializedConfig); err != nil {
		t.Fatalf("Failed to unmarshal JSON to LLMConfig: %v", err)
	}

	if deserializedConfig.Provider != config.Provider {
		t.Errorf("Expected Provider to be %s, got %s", config.Provider, deserializedConfig.Provider)
	}

	if deserializedConfig.APIKey != config.APIKey {
		t.Errorf("Expected APIKey to be %s, got %s", config.APIKey, deserializedConfig.APIKey)
	}
}
