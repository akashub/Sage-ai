// backend/go/internal/orchestrator/nodes/analyzer.go
package nodes

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
	"strconv"
	"time"
    "strings"
)

type Analyzer struct {
    bridge *llm.Bridge
}

func CreateAnalyzer(bridge *llm.Bridge) NodeExecutor {
    return &Analyzer{bridge: bridge}
}

func (a *Analyzer) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Analyzer Node: Starting analysis for query: %s", state.Query)

    // Extract schema
    schema, err := a.extractSchema(state.CSVPath)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to extract schema: %v", err)
        return fmt.Errorf("schema extraction failed: %w", err)
    }
    logger.InfoLogger.Printf("Schema extracted successfully")
    state.Schema = schema

    // Call LLM for analysis
    logger.InfoLogger.Printf("Sending query to LLM for analysis")
    analysis, err := a.bridge.Analyze(ctx, state.Query, schema)
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

func (a *Analyzer) extractSchema(csvPath string) (map[string]interface{}, error) {
    // Ensure the path is properly formed
    if csvPath == "" {
        return nil, fmt.Errorf("empty CSV path provided")
    }
    
    // Debug the current working directory
    currentDir, _ := os.Getwd()
    logger.InfoLogger.Printf("Current working directory: %s", currentDir)
    logger.InfoLogger.Printf("Attempting to open CSV at path: %s", csvPath)
    
    // Try to open the file
    file, err := os.Open(csvPath)
    if err != nil {
        // If the path doesn't start with "./", try prepending it
        if !strings.HasPrefix(csvPath, "./") && !strings.HasPrefix(csvPath, "/") {
            altPath := "./" + csvPath
            logger.InfoLogger.Printf("First attempt failed, trying alternative path: %s", altPath)
            file, err = os.Open(altPath)
            if err != nil {
                return nil, fmt.Errorf("failed to open CSV file at both %s and %s: %w", csvPath, altPath, err)
            }
        } else {
            return nil, err
        }
    }
    defer file.Close()

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