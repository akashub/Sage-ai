// backend/go/internal/types/state.go
package types

import (
	"sage-ai-v2/internal/knowledge"
)

// State represents the current state of query processing
type State struct {
	// Input
	Query   string                 `json:"query"`
	CSVPath string                 `json:"csv_path"`
	Options map[string]interface{} `json:"options,omitempty"`

	// Context
	Schema           map[string]interface{}     `json:"schema,omitempty"`
	KnowledgeContext *knowledge.KnowledgeResult `json:"knowledge_context,omitempty"`

	// Processing
	Analysis         map[string]interface{} `json:"analysis,omitempty"`
	GeneratedQuery   string                 `json:"generated_query,omitempty"`
	ValidationResult map[string]interface{} `json:"validation_result,omitempty"`
	ExecutionResult  interface{}            `json:"execution_result,omitempty"`

	// Errors
	Error string `json:"error,omitempty"`
}
