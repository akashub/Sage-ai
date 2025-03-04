// backend/go/internal/orchestrator/nodes/generator.go
package nodes

import (
	"context"
	"fmt"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
)

type Generator struct {
    bridge *llm.Bridge
}

func CreateGenerator(bridge *llm.Bridge) NodeExecutor {
    return &Generator{bridge: bridge}
}

func (g *Generator) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Generator Node: Starting SQL generation")
    
    if state.Analysis == nil {
        logger.ErrorLogger.Printf("Generator Node: No analysis available for query generation")
        return fmt.Errorf("no analysis available for query generation")
    }

    logger.DebugLogger.Printf("Generator Node: Using analysis: %+v", state.Analysis)
    logger.DebugLogger.Printf("Generator Node: Schema context: %+v", state.Schema)

    query, err := g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
    if err != nil {
        logger.ErrorLogger.Printf("Generator Node: Query generation failed: %v", err)
        return err
    }

    logger.InfoLogger.Printf("Generator Node: Successfully generated SQL query")
    logger.DebugLogger.Printf("Generated Query: %s", query)
    state.GeneratedQuery = query

    return nil
}

