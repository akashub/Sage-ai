// backend/go/internal/orchestrator/nodes/validator.go
package nodes

import (
	"context"
	"fmt"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/errors"
	"sage-ai-v2/pkg/logger"
)

type Validator struct {
    bridge *llm.Bridge
}

func CreateValidator(bridge *llm.Bridge) NodeExecutor {
    return &Validator{bridge: bridge}
}

func (v *Validator) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Validator Node: Starting query validation")

    if state.GeneratedQuery == "" {
        logger.ErrorLogger.Printf("Validator Node: No query available for validation")
        return fmt.Errorf("no query available for validation")
    }

    logger.DebugLogger.Printf("Validator Node: Validating query: %s", state.GeneratedQuery)
    logger.DebugLogger.Printf("Validator Node: Using schema: %+v", state.Schema)

    validation, err := v.bridge.ValidateQuery(ctx, state.GeneratedQuery, state.Schema)
    if err != nil {
        logger.ErrorLogger.Printf("Validator Node: Validation request failed: %v", err)
        return err
    }

    logger.InfoLogger.Printf("Validator Node: Validation completed")
    logger.DebugLogger.Printf("Validation Result: %+v", validation)

    if !validation["isValid"].(bool) {
        issues := validation["issues"].([]string)
        logger.ErrorLogger.Printf("Validator Node: Query validation failed with issues: %v", issues)
        return errors.NewValidationError(
            validation,
            fmt.Sprintf("Query validation failed: %v", issues),
        )
    }

    state.ValidationResult = validation
    logger.InfoLogger.Printf("Validator Node: Query validated successfully")
    return nil
}