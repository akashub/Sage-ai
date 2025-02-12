// backend/go/internal/orchestrator/nodes/executor.go
package nodes

import (
	"context"
	"fmt"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/csv"
	"sage-ai-v2/pkg/logger"
)

type Executor struct {
    parser *csv.Parser
}

func CreateExecutor() NodeExecutor {
    return &Executor{
        parser: csv.CreateParser(),
    }
}

// func (e *Executor) Execute(ctx context.Context, state *types.State) error {
//     if !state.ValidationResult["isValid"].(bool) {
//         return fmt.Errorf("cannot execute invalid query")
//     }

//     results, err := e.parser.ExecuteQuery(
//         state.CSVPath,
//         state.GeneratedQuery,
//     )
//     if err != nil {
//         return err
//     }

//     state.ExecutionResult = results
//     return nil
// }

func (e *Executor) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Executor Node: Starting query execution")

    if !state.ValidationResult["isValid"].(bool) {
        logger.ErrorLogger.Printf("Executor Node: Cannot execute invalid query")
        return fmt.Errorf("cannot execute invalid query")
    }

    logger.DebugLogger.Printf("Executor Node: Executing query: %s", state.GeneratedQuery)
    logger.DebugLogger.Printf("Executor Node: CSV Path: %s", state.CSVPath)

    results, err := e.parser.ExecuteQuery(
        state.CSVPath,
        state.GeneratedQuery,
    )
    if err != nil {
        logger.ErrorLogger.Printf("Executor Node: Query execution failed: %v", err)
        return err
    }

    logger.InfoLogger.Printf("Executor Node: Query executed successfully")
    if resultsArray, ok := results.([]map[string]interface{}); ok {
        logger.DebugLogger.Printf("Executor Node: Found %d results", len(resultsArray))
    }

    state.ExecutionResult = results
    return nil
}