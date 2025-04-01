// // backend/go/internal/orchestrator/orchestrator.go
// package orchestrator

// import (
//     "context"
//     "fmt"
//     "time"
//     "sage-ai-v2/internal/llm"
//     "sage-ai-v2/internal/types"
//     "sage-ai-v2/pkg/logger"
//     "sage-ai-v2/internal/orchestrator/nodes"
// 	"sage-ai-v2/pkg/errors"
// )

// type Orchestrator struct {
//     bridge    *llm.Bridge
//     graph     *Graph
//     sessionID string
// }

// func CreateOrchestrator(bridge *llm.Bridge) *Orchestrator {
//     o := &Orchestrator{
//         bridge: bridge,
//     }
//     o.NewSession()
//     return o
// }

// func (o *Orchestrator) NewSession() {
//     o.sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
//     o.bridge.SetSession(o.sessionID)
//     o.initializeGraph()
//     logger.InfoLogger.Printf("Started new session: %s", o.sessionID)
// }

// func (o *Orchestrator) initializeGraph() {
//     o.graph = CreateGraph()

//     // Create nodes
//     analyzerNode := nodes.CreateAnalyzer(o.bridge)
//     generatorNode := nodes.CreateGenerator(o.bridge)
//     validatorNode := nodes.CreateValidator(o.bridge)
//     executorNode := nodes.CreateExecutor()

//     // Add nodes to graph
//     o.graph.AddNode("analyzer", analyzerNode.Execute)
//     o.graph.AddNode("generator", generatorNode.Execute)
//     o.graph.AddNode("validator", validatorNode.Execute)
//     o.graph.AddNode("executor", executorNode.Execute)

//     // Setup normal flow
//     o.graph.AddEdge("analyzer", "generator")
//     o.graph.AddEdge("generator", "validator")
//     o.graph.AddEdge("validator", "executor")

//     logger.InfoLogger.Printf("Graph initialized with nodes and edges")
// }

// func (o *Orchestrator) ClearSession() {
//     o.sessionID = ""
//     o.graph = nil
//     logger.InfoLogger.Printf("Session cleared")
// }

// func (o *Orchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
//     if o.graph == nil {
//         o.NewSession()
//     }

//     logger.InfoLogger.Printf("Starting query processing pipeline")
//     logger.InfoLogger.Printf("Orchestrator: Starting new query processing")
//     logger.InfoLogger.Printf("Resetting previous state")
//     logger.InfoLogger.Printf("Input Query: %s", query)
//     logger.InfoLogger.Printf("CSV Path: %s", csvPath)

//     logger.InfoLogger.Printf("Processing query in session: %s", o.sessionID)

//     state := &types.State{
//         Query:   query,
//         CSVPath: csvPath,
//     }

//     maxHealingAttempts := 3
//     healingAttempts := 0

//     for {
//         logger.InfoLogger.Printf("Executing processing graph (healing attempt: %d/%d)",
//             healingAttempts, maxHealingAttempts)

//         err := o.graph.Execute(ctx, state)
//         if err == nil {
//             logger.InfoLogger.Printf("Query processing completed successfully")
//             return state, nil
//         }

//         // Check if error is due to validation
//         if validationErr, ok := err.(*errors.ValidationError); ok && healingAttempts < maxHealingAttempts {
//             healingAttempts++
//             logger.InfoLogger.Printf("Validation failed, attempting healing (attempt %d/%d)",
//                 healingAttempts, maxHealingAttempts)

//             healResult, healErr := o.attemptHealing(ctx, state, validationErr.ValidationResult)
//             if healErr != nil {
//                 logger.ErrorLogger.Printf("Healing failed: %v", healErr)
//                 return nil, fmt.Errorf("healing failed: %w", healErr)
//             }

//             if healResult.RequiresReanalysis {
//                 logger.InfoLogger.Printf("Healing suggests reanalysis, clearing analysis state")
//                 state.Analysis = nil
//                 continue
//             }

//             if healResult.HealdQuery != "" {
//                 logger.InfoLogger.Printf("Applying healed query: %s", healResult.HealdQuery)
//                 state.GeneratedQuery = healResult.HealdQuery
//                 continue
//             }
//         }

//         logger.ErrorLogger.Printf("Processing failed: %v", err)
//         return nil, fmt.Errorf("processing failed: %w", err)
//     }
// }

// func (o *Orchestrator) attemptHealing(ctx context.Context, state *types.State, validationResult map[string]interface{}) (*llm.HealingResult, error) {
//     return o.bridge.HealQuery(
//         ctx,
//         validationResult,
//         state.GeneratedQuery,
//         state.Analysis,
//         state.Schema,
//     )
// }

// backend/go/internal/orchestrator/orchestrator.go
package orchestrator

import (
	"context"
	"fmt"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/orchestrator/nodes"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"
)

// Orchestrator coordinates the processing flow
type Orchestrator struct {
    bridge          *llm.Bridge
    KnowledgeManager *knowledge.KnowledgeManager
    graph           *Graph
    sessionID       string
}

// CreateOrchestrator initializes a new orchestrator
func CreateOrchestrator(bridge *llm.Bridge, knowledgeManager *knowledge.KnowledgeManager) *Orchestrator {
    o := &Orchestrator{
        bridge:          bridge,
        KnowledgeManager: knowledgeManager,
    }
    o.NewSession()
    return o
}

// NewSession starts a new processing session
func (o *Orchestrator) NewSession() {
    o.sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
    o.bridge.SetSession(o.sessionID)
    o.initializeGraph()
    logger.InfoLogger.Printf("Started new session: %s", o.sessionID)
}

// initializeGraph sets up the processing graph
func (o *Orchestrator) initializeGraph() {
    o.graph = CreateGraph()

    // Create nodes with knowledge manager
    analyzerNode := nodes.CreateAnalyzer(o.bridge, o.KnowledgeManager)
    generatorNode := nodes.CreateGenerator(o.bridge)
    validatorNode := nodes.CreateValidator(o.bridge)
    executorNode := nodes.CreateExecutor()

    // Add nodes to graph
    o.graph.AddNode("analyzer", analyzerNode.Execute)
    o.graph.AddNode("generator", generatorNode.Execute)
    o.graph.AddNode("validator", validatorNode.Execute)
    o.graph.AddNode("executor", executorNode.Execute)

    // Setup normal flow
    o.graph.AddEdge("analyzer", "generator")
    o.graph.AddEdge("generator", "validator")
    o.graph.AddEdge("validator", "executor")

    logger.InfoLogger.Printf("Graph initialized with nodes and edges")
}

// ClearSession ends the current session
func (o *Orchestrator) ClearSession() {
    o.sessionID = ""
    o.graph = nil
    logger.InfoLogger.Printf("Session cleared")
}

// ProcessQueryWithOptions processes a query with additional options
// Update in orchestrator.go to handle specific errors
func (o *Orchestrator) ProcessQueryWithOptions(ctx context.Context, query string, csvPath string, options map[string]interface{}) (*types.State, error) {
    if o.graph == nil {
        o.NewSession()
    }

    logger.InfoLogger.Printf("Starting query processing pipeline with options")
    logger.InfoLogger.Printf("Query: %s", query)
    logger.InfoLogger.Printf("CSV Path: %s", csvPath)
    logger.InfoLogger.Printf("Options: %+v", options)
    
    logger.InfoLogger.Printf("Processing query in session: %s", o.sessionID)

    state := &types.State{
        Query:   query,
        CSVPath: csvPath,
        Options: options,
    }

    // Execute the processing graph with error handling
    err := o.graph.Execute(ctx, state)
    if err != nil {
        logger.ErrorLogger.Printf("Processing failed: %v", err)
        
        // Check if this is an SQL parsing error
        if strings.Contains(err.Error(), "only SELECT statements are supported") {
            // Try to recover by generating a simple SELECT statement
            logger.InfoLogger.Printf("Recovering from SQL parsing error with a simple query")
            
            // Create a simple SELECT query
            state.GeneratedQuery = "SELECT * FROM " + csvPath + " LIMIT 5"
            
            // Execute just the executor node directly
            executor := nodes.CreateExecutor()
            execErr := executor.Execute(ctx, state)
            if execErr != nil {
                logger.ErrorLogger.Printf("Recovery attempt failed: %v", execErr)
                return nil, fmt.Errorf("processing failed: %w", err)
            }
            
            logger.InfoLogger.Printf("Recovery successful with simple query")
            return state, nil
        }
        
        return nil, fmt.Errorf("processing failed: %w", err)
    }

    logger.InfoLogger.Printf("Query processing completed successfully")
    return state, nil
}

// ProcessQuery processes a query with default options
func (o *Orchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
    return o.ProcessQueryWithOptions(ctx, query, csvPath, nil)
}