// backend/go/internal/knowledge/types.go
package knowledge

// KnowledgeResult represents the result of a knowledge query
type KnowledgeResult struct {
	DDLSchemas       []TrainingItem    `json:"ddl_schemas"`
	Documentation    []TrainingItem    `json:"documentation"`
	QuestionSQLPairs []QuestionSQLPair `json:"question_sql_pairs"`
}

// TrainingItem represents a single item in the knowledge base
type TrainingItem struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "ddl", "documentation", "question_sql"
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Embedding   []float32              `json:"embedding,omitempty"`
	DateAdded   string                 `json:"date_added"` // ISO8601 string format
	Description string                 `json:"description,omitempty"`
}

// QuestionSQLPair represents a single QA item
type QuestionSQLPair struct {
	Question     string                 `json:"question"`
	SQL          string                 `json:"sql"`
	Schema       map[string]interface{} `json:"schema,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	DateAdded    string                 `json:"date_added"` // ISO8601 string format
	Verified     bool                   `json:"verified"`
}