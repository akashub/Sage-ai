package types

type State struct {
    Query            string                 `json:"query"`
    CSVPath          string                 `json:"csv_path"`
    Schema           map[string]interface{} `json:"schema"`
    Analysis         map[string]interface{} `json:"analysis"`
    GeneratedQuery   string                 `json:"generated_query"`
    ValidationResult map[string]interface{} `json:"validation_result"`
    ExecutionResult  interface{}            `json:"execution_result"`
    Error            string                 `json:"error,omitempty"`
}