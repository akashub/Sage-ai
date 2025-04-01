// // backend/go/internal/orchestrator/nodes/analyzer.go
// package nodes

// import (
// 	"context"
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"sage-ai-v2/internal/llm"
// 	"sage-ai-v2/internal/types"
// 	"sage-ai-v2/pkg/logger"
// 	"strconv"
// 	"time"
//     "strings"
// )

// type Analyzer struct {
//     bridge *llm.Bridge
// }

// func CreateAnalyzer(bridge *llm.Bridge) NodeExecutor {
//     return &Analyzer{bridge: bridge}
// }

// func (a *Analyzer) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Analyzer Node: Starting analysis for query: %s", state.Query)

//     // Extract schema
//     schema, err := a.extractSchema(state.CSVPath)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to extract schema: %v", err)
//         return fmt.Errorf("schema extraction failed: %w", err)
//     }
//     logger.InfoLogger.Printf("Schema extracted successfully")
//     state.Schema = schema

//     // Call LLM for analysis
//     logger.InfoLogger.Printf("Sending query to LLM for analysis")
//     analysis, err := a.bridge.Analyze(ctx, state.Query, schema)
//     if err != nil {
//         logger.ErrorLogger.Printf("LLM analysis failed: %v", err)
//         return fmt.Errorf("LLM analysis failed: %w", err)
//     }

//     // Update state
//     if analysis == nil {
//         logger.ErrorLogger.Printf("Received nil analysis from LLM")
//         return fmt.Errorf("nil analysis from LLM")
//     }

//     state.Analysis = analysis
//     logger.InfoLogger.Printf("Analysis completed: %+v", analysis)

//     return nil
// }

// func (a *Analyzer) extractSchema(csvPath string) (map[string]interface{}, error) {
//     // Ensure the path is properly formed
//     if csvPath == "" {
//         return nil, fmt.Errorf("empty CSV path provided")
//     }

//     // Debug the current working directory
//     currentDir, _ := os.Getwd()
//     logger.InfoLogger.Printf("Current working directory: %s", currentDir)
//     logger.InfoLogger.Printf("Attempting to open CSV at path: %s", csvPath)

//     // Try to open the file
//     file, err := os.Open(csvPath)
//     if err != nil {
//         // If the path doesn't start with "./", try prepending it
//         if !strings.HasPrefix(csvPath, "./") && !strings.HasPrefix(csvPath, "/") {
//             altPath := "./" + csvPath
//             logger.InfoLogger.Printf("First attempt failed, trying alternative path: %s", altPath)
//             file, err = os.Open(altPath)
//             if err != nil {
//                 return nil, fmt.Errorf("failed to open CSV file at both %s and %s: %w", csvPath, altPath, err)
//             }
//         } else {
//             return nil, err
//         }
//     }
//     defer file.Close()

//     reader := csv.NewReader(file)

//     headers, err := reader.Read()
//     if err != nil {
//         return nil, err
//     }

//     firstRow, err := reader.Read()
//     if err != nil {
//         return nil, err
//     }

//     schema := make(map[string]interface{})
//     for i, header := range headers {
//         schema[header] = map[string]string{
//             "sample": firstRow[i],
//             "inferred_type": inferType(firstRow[i]),
//         }
//     }

//     return schema, nil
// }

// func inferType(value string) string {
//     // Try to parse as integer
//     if _, err := strconv.ParseInt(value, 10, 64); err == nil {
//         return "integer"
//     }

//     // Try to parse as float
//     if _, err := strconv.ParseFloat(value, 64); err == nil {
//         return "float"
//     }

//     // Try to parse as date
//     if _, err := time.Parse("2006-01-02", value); err == nil {
//         return "date"
//     }

//     // Try to parse as datetime
//     if _, err := time.Parse(time.RFC3339, value); err == nil {
//         return "datetime"
//     }

//     // Default to string
//     return "string"
// }

// backend/go/internal/orchestrator/nodes/analyzer.go
package nodes

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
	"strconv"
	"strings"
	"time"
)

type Analyzer struct {
	bridge          *llm.Bridge
	knowledgeManager *knowledge.KnowledgeManager
}

func CreateAnalyzer(bridge *llm.Bridge, km *knowledge.KnowledgeManager) NodeExecutor {
	return &Analyzer{
		bridge:          bridge,
		knowledgeManager: km,
	}
}

func (a *Analyzer) Execute(ctx context.Context, state *types.State) error {
	logger.InfoLogger.Printf("Analyzer Node: Starting analysis for query: %s", state.Query)

	// Extract schema from CSV
	schema, err := a.extractSchema(state.CSVPath)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to extract schema: %v", err)
		return fmt.Errorf("schema extraction failed: %w", err)
	}
	logger.InfoLogger.Printf("Schema extracted successfully")
	state.Schema = schema

	// Retrieve relevant knowledge
	knowledgeResult, err := a.retrieveKnowledge(ctx, state.Query)
	if err != nil {
		logger.ErrorLogger.Printf("Knowledge retrieval warning: %v", err)
		// Continue even if knowledge retrieval fails
	} else {
		logger.InfoLogger.Printf("Retrieved relevant knowledge: %d DDL schemas, %d documentation items, %d question-SQL pairs",
			len(knowledgeResult.DDLSchemas), len(knowledgeResult.Documentation), len(knowledgeResult.QuestionSQLPairs))
		
		// Add knowledge context to state
		state.KnowledgeContext = knowledgeResult
	}

	// Call LLM for analysis with enhanced context
	analysis, err := a.analyzeWithKnowledge(ctx, state)
	if err != nil {
		logger.ErrorLogger.Printf("LLM analysis failed: %v", err)
		return fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Update state
	if analysis == nil {
		logger.ErrorLogger.Printf("Received nil analysis from LLM")
		return fmt.Errorf("nil analysis from LLM")
	}

	state.Analysis = analysis
	logger.InfoLogger.Printf("Analysis completed: %+v", analysis)

	return nil
}

func (a *Analyzer) retrieveKnowledge(ctx context.Context, query string) (*knowledge.KnowledgeResult, error) {
    if a.knowledgeManager == nil {
        logger.InfoLogger.Printf("Knowledge manager not available, skipping knowledge retrieval")
        return &knowledge.KnowledgeResult{}, nil
    }
    
    logger.InfoLogger.Printf("Retrieving knowledge for query: %s", query)
    result, err := a.knowledgeManager.RetrieveRelevantKnowledge(ctx, query)
    
    if err != nil {
        logger.ErrorLogger.Printf("Error retrieving knowledge: %v", err)
        return &knowledge.KnowledgeResult{}, err
    }
    
    // Log retrieved knowledge
    logger.InfoLogger.Printf("Retrieved DDL schemas: %d", len(result.DDLSchemas))
    for i, schema := range result.DDLSchemas {
        logger.InfoLogger.Printf("  Schema %d: %s", i, schema.Description)
    }
    
    logger.InfoLogger.Printf("Retrieved Documentation: %d", len(result.Documentation))
    for i, doc := range result.Documentation {
        logger.InfoLogger.Printf("  Doc %d: %s", i, doc.Description)
    }
    
    logger.InfoLogger.Printf("Retrieved Question-SQL pairs: %d", len(result.QuestionSQLPairs))
    for i, pair := range result.QuestionSQLPairs {
        logger.InfoLogger.Printf("  Pair %d: %s", i, pair.Description)
    }
    
    return result, nil
}

func (a *Analyzer) analyzeWithKnowledge(ctx context.Context, state *types.State) (map[string]interface{}, error) {
	// Prepare knowledge context
	knowledgeContext := map[string]interface{}{
		"csv_schema": state.Schema,
	}
	
	// Add DDL schemas if available
	if state.KnowledgeContext != nil && len(state.KnowledgeContext.DDLSchemas) > 0 {
		var ddlSchemas []string
		for _, ddl := range state.KnowledgeContext.DDLSchemas {
			ddlSchemas = append(ddlSchemas, ddl.Content)
		}
		knowledgeContext["ddl_schemas"] = ddlSchemas
	}
	
	// Add documentation if available
	if state.KnowledgeContext != nil && len(state.KnowledgeContext.Documentation) > 0 {
		var docs []map[string]string
		for _, doc := range state.KnowledgeContext.Documentation {
			docs = append(docs, map[string]string{
				"title":   doc.Description,
				"content": doc.Content,
			})
		}
		knowledgeContext["documentation"] = docs
	}
	
	// Add question-SQL pairs if available
	if state.KnowledgeContext != nil && len(state.KnowledgeContext.QuestionSQLPairs) > 0 {
		var examples []map[string]string
		for _, pair := range state.KnowledgeContext.QuestionSQLPairs {
			examples = append(examples, map[string]string{
				"question": pair.Question,
				"sql":      pair.SQL,
			})
		}
		knowledgeContext["examples"] = examples
	}
	
	// Convert knowledge context to string representation for the LLM
	knowledgeContextJSON, err := json.Marshal(knowledgeContext)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to marshal knowledge context: %v", err)
		// Fall back to standard analysis without knowledge context
		return a.bridge.Analyze(ctx, state.Query, state.Schema)
	}
	
	// Add knowledge context to analysis request
	analysisRequest := map[string]interface{}{
		"query":            state.Query,
		"schema":           state.Schema,
		"knowledge_context": json.RawMessage(knowledgeContextJSON),
	}
	
	return a.bridge.AnalyzeWithKnowledge(ctx, analysisRequest)
}

// Improved extractSchema function for backend/go/internal/orchestrator/nodes/analyzer.go

func (a *Analyzer) extractSchema(csvPath string) (map[string]interface{}, error) {
    // Ensure the path is properly formed
    if csvPath == "" {
        return nil, fmt.Errorf("empty CSV path provided")
    }
    
    logger.InfoLogger.Printf("Received CSV path: %s", csvPath)
    
    // Current working directory for logging purposes
    currentDir, _ := os.Getwd()
    logger.InfoLogger.Printf("Current working directory: %s", currentDir)
    
    // Create a list of possible paths to try
    possiblePaths := []string{
        csvPath,                                // As provided
        filepath.Join("data", "uploads", csvPath), // In uploads folder
    }
    
    // If the path contains a filename without timestamp, try to find it
    if !strings.Contains(csvPath, "_") && strings.Contains(csvPath, ".") {
        // List files in uploads directory
        uploadsDir := filepath.Join("data", "uploads")
        files, err := os.ReadDir(uploadsDir)
        if err == nil {
            for _, file := range files {
                fileName := file.Name()
                // Check if the file ends with the provided filename (after the timestamp)
                if strings.HasSuffix(fileName, csvPath) {
                    possiblePaths = append(possiblePaths, filepath.Join(uploadsDir, fileName))
                }
            }
        }
    }
    
    // Try each possible path
    var lastErr error
    for _, path := range possiblePaths {
        logger.InfoLogger.Printf("Attempting to open CSV at path: %s", path)
        
        // Try to open the file
        file, err := os.Open(path)
        if err == nil {
            defer file.Close()
            logger.InfoLogger.Printf("Successfully opened file at: %s", path)
            
            reader := csv.NewReader(file)
            
            headers, err := reader.Read()
            if err != nil {
                return nil, err
            }

            firstRow, err := reader.Read()
            if err != nil {
                return nil, err
            }

            schema := make(map[string]interface{})
            for i, header := range headers {
                schema[header] = map[string]string{
                    "sample": firstRow[i],
                    "inferred_type": inferType(firstRow[i]),
                }
            }

            return schema, nil
        }
        
        lastErr = err
    }
    
    // If we're here, none of the paths worked
    return nil, fmt.Errorf("failed to open CSV file, tried multiple paths: %w", lastErr)
}

func inferType(value string) string {
	// Try to parse as integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "integer"
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float"
	}

	// Try to parse as date
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return "date"
	}

	// Try to parse as datetime
	if _, err := time.Parse(time.RFC3339, value); err == nil {
		return "datetime"
	}

	// Default to string
	return "string"
}