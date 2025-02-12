// // backend/go/internal/orchestrator/orchestrator.go
// package orchestrator

// import (
// 	"context"
// 	"fmt"
// 	"sage-ai-v2/internal/llm"
// 	"sage-ai-v2/internal/orchestrator/nodes"
// 	"sage-ai-v2/internal/types"
// 	"sage-ai-v2/pkg/errors"
// 	"sage-ai-v2/pkg/logger"
// 	"time"
// )

// // ValidationError represents a validation failure that might need healing
// type ValidationError struct {
//     ValidationResult map[string]interface{}
//     Message         string
// }

// func (e *ValidationError) Error() string {
//     return e.Message
// }

// // HealingError represents a failure in the healing process
// type HealingError struct {
//     Message string
//     Cause   error
// }

// func (e *HealingError) Error() string {
//     if e.Cause != nil {
//         return fmt.Sprintf("%s: %v", e.Message, e.Cause)
//     }
//     return e.Message
// }

// // Orchestrator manages the text-to-SQL pipeline
// type Orchestrator struct {
//     bridge *llm.Bridge
//     graph  *Graph
//     sessionID string
// }

// func (o *Orchestrator) NewSession() {
//     o.sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
//     o.graph = CreateGraph() // Reset graph for new session
//     logger.InfoLogger.Printf("Started new session: %s", o.sessionID)
// }

// func (o *Orchestrator) ClearSession() {
//     o.sessionID = ""
//     o.graph = nil
//     logger.InfoLogger.Printf("Cleared session state")
// }


// func CreateOrchestrator(bridge *llm.Bridge) *Orchestrator {
//     orch := &Orchestrator{
//         bridge: bridge,
//         graph:  CreateGraph(),
//     }
    
//     orch.setupGraph()
//     return orch
// }

// func (o *Orchestrator) setupGraph() {
//     // Initialize nodes
//     analyzer := nodes.CreateAnalyzer(o.bridge)
//     generator := nodes.CreateGenerator(o.bridge)
//     validator := nodes.CreateValidator(o.bridge)
//     executor := nodes.CreateExecutor()

//     // Add nodes to graph
//     o.graph.AddNode("analyzer", analyzer.Execute)
//     o.graph.AddNode("generator", generator.Execute)
//     o.graph.AddNode("validator", validator.Execute)
//     o.graph.AddNode("executor", executor.Execute)

//     // Setup normal flow
//     o.graph.AddEdge("analyzer", "generator")
//     o.graph.AddEdge("generator", "validator")
//     o.graph.AddEdge("validator", "executor")
// }

// // func (o *Orchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
// //     state := &types.State{
// //         Query:   query,
// //         CSVPath: csvPath,
// //     }

// //     maxHealingAttempts := 3
// //     healingAttempts := 0

// //     for {
// //         err := o.graph.Execute(ctx, state)
// //         if err == nil {
// //             return state, nil
// //         }

// //         // Check if error is due to validation
// //         if validationErr, ok := err.(*ValidationError); ok && healingAttempts < maxHealingAttempts {
// //             healingAttempts++
            
// //             // Log healing attempt
// //             fmt.Printf("Attempting healing (attempt %d/%d)\n", healingAttempts, maxHealingAttempts)
            
// //             // Attempt healing
// //             healResult, healErr := o.attemptHealing(ctx, state, validationErr.ValidationResult)
// //             if healErr != nil {
// //                 return nil, &HealingError{
// //                     Message: "healing failed",
// //                     Cause:   healErr,
// //                 }
// //             }

// //             // Apply healing results
// //             if healResult.RequiresReanalysis {
// //                 fmt.Println("Healing suggests reanalysis, clearing analysis state")
// //                 state.Analysis = nil
// //                 continue
// //             }

// //             if healResult.HealdQuery != "" {
// //                 fmt.Printf("Applying healed query: %s\n", healResult.HealdQuery)
// //                 state.GeneratedQuery = healResult.HealdQuery
// //                 continue
// //             }

// //             if healResult.RequiresHumanReview {
// //                 return nil, &HealingError{
// //                     Message: "Query requires human review: " + healResult.Notes,
// //                 }
// //             }
// //         }

// //         // If we reach here, either it's not a validation error or healing failed
// //         return nil, fmt.Errorf("processing failed: %w", err)
// //     }
// // }
// func (o *Orchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
//     logger.InfoLogger.Printf("Starting query processing pipeline")
// 	logger.InfoLogger.Printf("Orchestrator: Starting new query processing")
//     logger.InfoLogger.Printf("Resetting previous state")
//     logger.InfoLogger.Printf("Input Query: %s", query)
//     logger.InfoLogger.Printf("CSV Path: %s", csvPath)

// 	if o.sessionID == "" {
//         o.NewSession()
//     }
//     logger.InfoLogger.Printf("Processing query in session: %s", o.sessionID)
// 	// Creating context with timeout
//     // ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
//     // defer cancel()

//     state := &types.State{
//         Query:   query,
//         CSVPath: csvPath,
// 		Schema:           make(map[string]interface{}),
//         Analysis:         make(map[string]interface{}),
//         ValidationResult: make(map[string]interface{}),
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

// // Additional utility methods

// func (o *Orchestrator) GetExecutionStatus(state *types.State) string {
//     if state.Error != "" {
//         return "ERROR"
//     }
//     if state.ExecutionResult != nil {
//         return "COMPLETED"
//     }
//     if state.ValidationResult != nil {
//         return "VALIDATED"
//     }
//     if state.GeneratedQuery != "" {
//         return "GENERATED"
//     }
//     if state.Analysis != nil {
//         return "ANALYZED"
//     }
//     return "INITIALIZED"
// }

// func (o *Orchestrator) Reset() error {
//     o.graph = CreateGraph()
//     o.setupGraph()
//     return nil
// }

// backend/go/internal/orchestrator/orchestrator.go
package orchestrator

import (
    "context"
    "fmt"
    "time"
    "sage-ai-v2/internal/llm"
    "sage-ai-v2/internal/types"
    "sage-ai-v2/pkg/logger"
    "sage-ai-v2/internal/orchestrator/nodes"
	"sage-ai-v2/pkg/errors"
)

type Orchestrator struct {
    bridge    *llm.Bridge
    graph     *Graph
    sessionID string
}

func CreateOrchestrator(bridge *llm.Bridge) *Orchestrator {
    o := &Orchestrator{
        bridge: bridge,
    }
    o.NewSession()
    return o
}

func (o *Orchestrator) NewSession() {
    o.sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
    o.bridge.SetSession(o.sessionID)
    o.initializeGraph()
    logger.InfoLogger.Printf("Started new session: %s", o.sessionID)
}

func (o *Orchestrator) initializeGraph() {
    o.graph = CreateGraph()

    // Create nodes
    analyzerNode := nodes.CreateAnalyzer(o.bridge)
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

func (o *Orchestrator) ClearSession() {
    o.sessionID = ""
    o.graph = nil
    logger.InfoLogger.Printf("Session cleared")
}

func (o *Orchestrator) ProcessQuery(ctx context.Context, query string, csvPath string) (*types.State, error) {
    if o.graph == nil {
        o.NewSession()
    }

    logger.InfoLogger.Printf("Starting query processing pipeline")
    logger.InfoLogger.Printf("Orchestrator: Starting new query processing")
    logger.InfoLogger.Printf("Resetting previous state")
    logger.InfoLogger.Printf("Input Query: %s", query)
    logger.InfoLogger.Printf("CSV Path: %s", csvPath)
    
    logger.InfoLogger.Printf("Processing query in session: %s", o.sessionID)

    state := &types.State{
        Query:   query,
        CSVPath: csvPath,
    }

    maxHealingAttempts := 3
    healingAttempts := 0

    for {
        logger.InfoLogger.Printf("Executing processing graph (healing attempt: %d/%d)", 
            healingAttempts, maxHealingAttempts)

        err := o.graph.Execute(ctx, state)
        if err == nil {
            logger.InfoLogger.Printf("Query processing completed successfully")
            return state, nil
        }

        // Check if error is due to validation
        if validationErr, ok := err.(*errors.ValidationError); ok && healingAttempts < maxHealingAttempts {
            healingAttempts++
            logger.InfoLogger.Printf("Validation failed, attempting healing (attempt %d/%d)", 
                healingAttempts, maxHealingAttempts)
            
            healResult, healErr := o.attemptHealing(ctx, state, validationErr.ValidationResult)
            if healErr != nil {
                logger.ErrorLogger.Printf("Healing failed: %v", healErr)
                return nil, fmt.Errorf("healing failed: %w", healErr)
            }

            if healResult.RequiresReanalysis {
                logger.InfoLogger.Printf("Healing suggests reanalysis, clearing analysis state")
                state.Analysis = nil
                continue
            }

            if healResult.HealdQuery != "" {
                logger.InfoLogger.Printf("Applying healed query: %s", healResult.HealdQuery)
                state.GeneratedQuery = healResult.HealdQuery
                continue
            }
        }

        logger.ErrorLogger.Printf("Processing failed: %v", err)
        return nil, fmt.Errorf("processing failed: %w", err)
    }
}

func (o *Orchestrator) attemptHealing(ctx context.Context, state *types.State, validationResult map[string]interface{}) (*llm.HealingResult, error) {
    return o.bridge.HealQuery(
        ctx,
        validationResult,
        state.GeneratedQuery,
        state.Analysis,
        state.Schema,
    )
}