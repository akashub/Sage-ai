// backend/go/internal/llm/provider.go
package llm

type LLMProvider string

const (
    ProviderGemini   LLMProvider = "gemini"
    ProviderOpenAI   LLMProvider = "openai"
    ProviderAnthropic LLMProvider = "anthropic"
    ProviderMistral  LLMProvider = "mistral"
)

type LLMConfig struct {
    Provider LLMProvider `json:"provider"`
    APIKey   string     `json:"api_key"`
}