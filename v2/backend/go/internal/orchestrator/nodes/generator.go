// // backend/go/internal/orchestrator/nodes/generator.go
// package nodes

// import (
// 	"context"
// 	"fmt"
// 	"sage-ai-v2/internal/llm"
// 	"sage-ai-v2/internal/types"
// 	"sage-ai-v2/pkg/logger"
// )

// type Generator struct {
//     bridge *llm.Bridge
// }

// func CreateGenerator(bridge *llm.Bridge) NodeExecutor {
//     return &Generator{bridge: bridge}
// }

// func (g *Generator) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Generator Node: Starting SQL generation")

//     if state.Analysis == nil {
//         logger.ErrorLogger.Printf("Generator Node: No analysis available for query generation")
//         return fmt.Errorf("no analysis available for query generation")
//     }

//     logger.DebugLogger.Printf("Generator Node: Using analysis: %+v", state.Analysis)
//     logger.DebugLogger.Printf("Generator Node: Schema context: %+v", state.Schema)

//     query, err := g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
//     if err != nil {
//         logger.ErrorLogger.Printf("Generator Node: Query generation failed: %v", err)
//         return err
//     }

//     logger.InfoLogger.Printf("Generator Node: Successfully generated SQL query")
//     logger.DebugLogger.Printf("Generated Query: %s", query)
//     state.GeneratedQuery = query

//     return nil
// }

// backend/go/internal/orchestrator/nodes/generator.go
package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
	"strconv"
	"strings"
)

// Generator is responsible for SQL query generation
type Generator struct {
    bridge *llm.Bridge
}

// CreateGenerator initializes a new generator node
func CreateGenerator(bridge *llm.Bridge) NodeExecutor {
    return &Generator{bridge: bridge}
}

// Execute runs the generator node
// func (g *Generator) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Generator Node: Starting SQL generation")
    
//     if state.Analysis == nil {
//         logger.ErrorLogger.Printf("Generator Node: No analysis available for query generation")
//         return fmt.Errorf("no analysis available for query generation")
//     }

//     logger.DebugLogger.Printf("Generator Node: Using analysis: %+v", state.Analysis)
//     logger.DebugLogger.Printf("Generator Node: Schema context: %+v", state.Schema)

//     // Check if we should use knowledge base
//     useKnowledgeBase := true
//     if state.Options != nil {
//         if val, ok := state.Options["useKnowledgeBase"]; ok {
//             if boolVal, ok := val.(bool); ok {
//                 useKnowledgeBase = boolVal
//             }
//         }
//     }

//     var query string
//     var err error

//     if useKnowledgeBase && state.KnowledgeContext != nil {
//         // Use knowledge-enhanced generation
//         knowledgeContext := map[string]interface{}{
//             "csv_schema": state.Schema,
//         }
        
//         // Add DDL schemas if available
//         if len(state.KnowledgeContext.DDLSchemas) > 0 {
//             var ddlSchemas []string
//             for _, ddl := range state.KnowledgeContext.DDLSchemas {
//                 ddlSchemas = append(ddlSchemas, ddl.Content)
//             }
//             knowledgeContext["ddl_schemas"] = ddlSchemas
//         }
        
//         // Add question-SQL pairs if available
//         if len(state.KnowledgeContext.QuestionSQLPairs) > 0 {
//             var examples []map[string]string
//             for _, pair := range state.KnowledgeContext.QuestionSQLPairs {
//                 examples = append(examples, map[string]string{
//                     "question": pair.Question,
//                     "sql":      pair.SQL,
//                 })
//             }
//             knowledgeContext["examples"] = examples
//         }
        
//         // Convert knowledge context to JSON
//         knowledgeContextJSON, err := json.Marshal(knowledgeContext)
//         if err == nil {
//             logger.InfoLogger.Printf("Generator Node: Using knowledge context for generation")
            
//             request := map[string]interface{}{
//                 "analysis":         state.Analysis,
//                 "schema":           state.Schema,
//                 "knowledge_context": json.RawMessage(knowledgeContextJSON),
//             }
            
//             query, err = g.bridge.GenerateQueryWithKnowledge(ctx, request)
//         } else {
//             logger.ErrorLogger.Printf("Failed to marshal knowledge context: %v", err)
//             query, err = g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
//         }

        
//     } else {
//         // Use standard generation
//         query, err = g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
//     }

//     if err != nil {
//         logger.ErrorLogger.Printf("Generator Node: Query generation failed: %v", err)
//         return err
//     }

//     logger.InfoLogger.Printf("Generator Node: Successfully generated SQL query")
//     logger.DebugLogger.Printf("Generated Query: %s", query)
//     state.GeneratedQuery = query

//     return nil
// }

// Execute runs the generator node
// Execute runs the generator node
func (g *Generator) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Generator Node: Starting SQL generation")
    
    if state.Analysis == nil {
        logger.ErrorLogger.Printf("Generator Node: No analysis available for query generation")
        return fmt.Errorf("no analysis available for query generation")
    }

    logger.DebugLogger.Printf("Generator Node: Using analysis: %+v", state.Analysis)
    logger.DebugLogger.Printf("Generator Node: Schema context: %+v", state.Schema)

    // Check if we should use knowledge base
    useKnowledgeBase := true
    if state.Options != nil {
        if val, ok := state.Options["useKnowledgeBase"]; ok {
            if boolVal, ok := val.(bool); ok {
                useKnowledgeBase = boolVal
            }
        }
    }

    var query string
    var err error

    if useKnowledgeBase && state.KnowledgeContext != nil {
        // Use knowledge-enhanced generation
        knowledgeContext := map[string]interface{}{
            "csv_schema": state.Schema,
        }
        
        // Add DDL schemas if available
        if len(state.KnowledgeContext.DDLSchemas) > 0 {
            var ddlSchemas []string
            for _, ddl := range state.KnowledgeContext.DDLSchemas {
                ddlSchemas = append(ddlSchemas, ddl.Content)
            }
            knowledgeContext["ddl_schemas"] = ddlSchemas
        }
        
        // Add question-SQL pairs if available
        if len(state.KnowledgeContext.QuestionSQLPairs) > 0 {
            var examples []map[string]string
            for _, pair := range state.KnowledgeContext.QuestionSQLPairs {
                examples = append(examples, map[string]string{
                    "question": pair.Question,
                    "sql":      pair.SQL,
                })
            }
            knowledgeContext["examples"] = examples
        }
        
        // Convert knowledge context to JSON
        knowledgeContextJSON, err := json.Marshal(knowledgeContext)
        if err == nil {
            logger.InfoLogger.Printf("Generator Node: Using knowledge context for generation")
            
            request := map[string]interface{}{
                "analysis":         state.Analysis,
                "schema":           state.Schema,
                "knowledge_context": json.RawMessage(knowledgeContextJSON),
            }
            
            query, err = g.bridge.GenerateQueryWithKnowledge(ctx, request)
        } else {
            logger.ErrorLogger.Printf("Failed to marshal knowledge context: %v", err)
            query, err = g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
        }
    } else {
        // Use standard generation
        query, err = g.bridge.GenerateQuery(ctx, state.Analysis, state.Schema)
    }

    if err != nil {
        logger.ErrorLogger.Printf("Generator Node: Query generation failed: %v", err)
        return err
    }

    // Validate and clean up the generated query
    query = strings.TrimSpace(query)
    
    // Clean any prefix or suffix that might not be part of the SQL query
    // This fixes the issue where "sql" might be prefixed
    for _, prefix := range []string{"sql", "SQL"} {
        if strings.HasPrefix(query, prefix) {
            query = strings.TrimSpace(query[len(prefix):])
        }
    }
    
    // Clean any quotes or other extraneous characters
    query = strings.Trim(query, "`'\";")
    
    // Ensure it starts with SELECT
    if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
        logger.ErrorLogger.Printf("Generator Node: Generated non-SELECT query: %s", query)
        
        // Extract table name from schema if possible
        tableName := "data"
        if state.CSVPath != "" {
            // Use CSV filename as the table name
            base := filepath.Base(state.CSVPath)
            tableName = strings.TrimSuffix(base, filepath.Ext(base))
            // Remove timestamp prefix if present
            if parts := strings.SplitN(tableName, "_", 2); len(parts) > 1 {
                if _, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
                    tableName = parts[1]
                }
            }
        }
        
        // Create a fallback query
        query = fmt.Sprintf("SELECT * FROM %s LIMIT 10", tableName)
        logger.InfoLogger.Printf("Generator Node: Using fallback query: %s", query)
    }

    logger.InfoLogger.Printf("Generator Node: Successfully generated SQL query")
    logger.DebugLogger.Printf("Generated Query: %s", query)
    state.GeneratedQuery = query

    return nil
}