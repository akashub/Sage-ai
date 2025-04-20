// // backend/go/internal/orchestrator/nodes/executor.go
// package nodes

// import (
// 	"context"
// 	"fmt"
// 	"sage-ai-v2/internal/types"
// 	"sage-ai-v2/pkg/csv"
// 	"sage-ai-v2/pkg/logger"
// )

// type Executor struct {
//     parser *csv.Parser
// }

// func CreateExecutor() NodeExecutor {
//     return &Executor{
//         parser: csv.CreateParser(),
//     }
// }

// func (e *Executor) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Executor Node: Starting query execution")

//     if !state.ValidationResult["isValid"].(bool) {
//         logger.ErrorLogger.Printf("Executor Node: Cannot execute invalid query")
//         return fmt.Errorf("cannot execute invalid query")
//     }

//     logger.DebugLogger.Printf("Executor Node: Executing query: %s", state.GeneratedQuery)
//     logger.DebugLogger.Printf("Executor Node: CSV Path: %s", state.CSVPath)

//     results, err := e.parser.ExecuteQuery(
//         state.CSVPath,
//         state.GeneratedQuery,
//     )
//     if err != nil {
//         logger.ErrorLogger.Printf("Executor Node: Query execution failed: %v", err)
//         return err
//     }

//     logger.InfoLogger.Printf("Executor Node: Query executed successfully")
//     if resultsArray, ok := results.([]map[string]interface{}); ok {
//         logger.DebugLogger.Printf("Executor Node: Found %d results", len(resultsArray))
//     }

//     state.ExecutionResult = results
//     return nil
// }

// // NodeExecutor defines the interface for all processing nodes
// type NodeExecutor interface {
// 	Execute(ctx context.Context, state *types.State) error
// }

// backend/go/internal/orchestrator/nodes/executor.go
package nodes

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
	"sort"
	"strconv"
	"strings"
)

// Executor executes SQL queries against CSV data
type Executor struct{}

// CreateExecutor initializes a new executor node
func CreateExecutor() NodeExecutor {
    return &Executor{}
}

func (e *Executor) parseSQL(sql string) (map[string]interface{}, error) {
    // Simplified parser for SELECT statements
    sql = strings.TrimSpace(sql)
    
    // Basic structure check - ensure it's a SELECT statement
    if !strings.HasPrefix(strings.ToUpper(sql), "SELECT") {
        return nil, fmt.Errorf("only SELECT statements are supported")
    }
    
    query := map[string]interface{}{
        "type": "select",
    }
    
    // Extract columns - handle different SQL dialects and formats
    fromIdx := strings.Index(strings.ToUpper(sql), " FROM ")
    if fromIdx <= 0 {
        // Try alternate format without space
        fromIdx = strings.Index(strings.ToUpper(sql), "FROM")
        if fromIdx <= 0 {
            return nil, fmt.Errorf("invalid SELECT statement format, cannot find FROM clause")
        }
    }
    
    columnsStr := sql[6:fromIdx]
    columnsStr = strings.TrimSpace(columnsStr)
    
    if columnsStr == "*" {
        query["columns"] = []string{"*"}
    } else {
        // Split by commas, handle functions and aliases
        columns := []string{}
        for _, col := range strings.Split(columnsStr, ",") {
            col = strings.TrimSpace(col)
            columns = append(columns, col)
        }
        query["columns"] = columns
    }
    
    // Extract table
    remainingSQL := sql[fromIdx+5:] // +5 for "FROM "
    remainingSQL = strings.TrimSpace(remainingSQL)
    
    // If there's a WHERE, LIMIT, etc.
    endPos := len(remainingSQL)
    for _, keyword := range []string{" WHERE ", " ORDER BY ", " LIMIT ", " GROUP BY "} {
        if pos := strings.Index(strings.ToUpper(remainingSQL), keyword); pos > 0 && pos < endPos {
            endPos = pos
        }
    }
    
    tableName := remainingSQL[:endPos]
    tableName = strings.TrimSpace(tableName)
    query["from"] = tableName
    
    // Extract LIMIT if exists
    if limitPos := strings.Index(strings.ToUpper(sql), " LIMIT "); limitPos > 0 {
        limitStr := sql[limitPos+7:]
        limitStr = strings.TrimSpace(limitStr)
        
        // Remove trailing semicolon if any
        limitStr = strings.TrimSuffix(limitStr, ";")
        
        // Handle potential trailing clauses
        if spacePos := strings.Index(limitStr, " "); spacePos > 0 {
            limitStr = limitStr[:spacePos]
        }
        
        if limit, err := strconv.Atoi(limitStr); err == nil {
            query["limit"] = limit
        }
    }
    
    return query, nil
}

// Update the Execute method to include better error handling
func (e *Executor) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Executor Node: Starting query execution")

    if state.GeneratedQuery == "" {
        logger.ErrorLogger.Printf("Executor Node: No query available to execute")
        return fmt.Errorf("no query available to execute")
    }

    // Preprocess the query to ensure it's a valid SELECT statement
    query := state.GeneratedQuery
    query = strings.TrimSpace(query)
    
    // Clean up the query - remove any leading or trailing quotes, backticks, etc.
    query = strings.Trim(query, "'\"`;")
    
    // Ensure the query is a SELECT statement
    if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
        logger.ErrorLogger.Printf("Executor Node: Query must be a SELECT statement: %s", query)
        return fmt.Errorf("only SELECT statements are supported: %s", query)
    }

    // Execute the query with the preprocessed statement
    results, err := e.executeQuery(state.CSVPath, query)
    if err != nil {
        logger.ErrorLogger.Printf("Executor Node: Query execution failed: %v", err)
        return fmt.Errorf("query execution failed: %w", err)
    }

    logger.InfoLogger.Printf("Executor Node: Query executed successfully")
    state.ExecutionResult = results
    return nil
}

// executeQuery executes a SQL-like query against CSV data
func (e *Executor) executeQuery(csvPath string, sqlQuery string) ([]map[string]interface{}, error) {
    // Parse SQL into structured query
    parsedQuery, err := e.parseSQL(sqlQuery)
    if err != nil {
        return nil, fmt.Errorf("error parsing SQL: %w", err)
    }
    
    // Open CSV file
    file, err := os.Open(csvPath)
    if err != nil {
        return nil, fmt.Errorf("error opening CSV file: %w", err)
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    
    // Read headers
    headers, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("error reading CSV headers: %w", err)
    }
    
    // Build column index map for quick lookup
    colMap := make(map[string]int)
    for i, h := range headers {
        colMap[h] = i
    }
    
    // Determine which columns to select
    selectedCols := []string{}
    columns, _ := parsedQuery["columns"].([]string)
    if len(columns) == 0 || (len(columns) == 1 && columns[0] == "*") {
        // SELECT * - select all columns
        selectedCols = headers
    } else {
        for _, col := range columns {
            // Handle potential "AS" aliases
            if strings.Contains(col, " AS ") {
                parts := strings.Split(col, " AS ")
                colName := strings.TrimSpace(parts[0])
                if _, ok := colMap[colName]; ok {
                    selectedCols = append(selectedCols, colName)
                }
            } else if _, ok := colMap[col]; ok {
                selectedCols = append(selectedCols, col)
            }
        }
    }
    
    // Read all rows
    var rows [][]string
    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("error reading CSV row: %w", err)
        }
        rows = append(rows, row)
    }
    
    // Apply filtering conditions if any
    var filteredRows [][]string
    conditions, hasConditions := parsedQuery["conditions"].([]string)
    
    if hasConditions {
        for _, row := range rows {
            matches := true
            for _, condition := range conditions {
                if !e.evaluateCondition(condition, row, headers, colMap) {
                    matches = false
                    break
                }
            }
            if matches {
                filteredRows = append(filteredRows, row)
            }
        }
    } else {
        filteredRows = rows
    }
    
    // Prepare result set
    var results []map[string]interface{}
    for _, row := range filteredRows {
        result := make(map[string]interface{})
        for _, col := range selectedCols {
            idx := colMap[col]
            if idx < len(row) {
                // Try to convert to appropriate type
                val := row[idx]
                if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
                    result[col] = intVal
                } else if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
                    result[col] = floatVal
                } else {
                    result[col] = val
                }
            } else {
                result[col] = nil
            }
        }
        results = append(results, result)
    }
    
    // Handle ORDER BY if specified
    orderBy, hasOrderBy := parsedQuery["orderBy"].(map[string]interface{})
    if hasOrderBy {
        column, hasColumn := orderBy["column"].(string)
        direction, hasDirection := orderBy["direction"].(string)
        
        if hasColumn {
            // Default to ascending if not specified
            ascending := true
            if hasDirection && strings.ToUpper(direction) == "DESC" {
                ascending = false
            }
            
            // Sort the results
            sort.Slice(results, func(i, j int) bool {
                // Get values to compare
                valI, hasI := results[i][column]
                valJ, hasJ := results[j][column]
                
                // Handle missing values
                if !hasI && !hasJ {
                    return false // Both missing, no change
                } else if !hasI {
                    return ascending // I missing, I comes first if ascending
                } else if !hasJ {
                    return !ascending // J missing, J comes first if ascending
                }
                
                // Compare based on types
                switch vi := valI.(type) {
                case int64:
                    if vj, ok := valJ.(int64); ok {
                        if ascending {
                            return vi < vj
                        } else {
                            return vi > vj
                        }
                    }
                case float64:
                    if vj, ok := valJ.(float64); ok {
                        if ascending {
                            return vi < vj
                        } else {
                            return vi > vj
                        }
                    }
                case string:
                    if vj, ok := valJ.(string); ok {
                        if ascending {
                            return vi < vj
                        } else {
                            return vi > vj
                        }
                    }
                }
                
                // If types don't match, convert to string and compare
                strI := fmt.Sprintf("%v", valI)
                strJ := fmt.Sprintf("%v", valJ)
                if ascending {
                    return strI < strJ
                } else {
                    return strI > strJ
                }
            })
        }
    } else {
        // If no explicit orderBy in parsed query, check for "sort" key
        // which might be used by the LLM in analysis
        if sortData, hasSort := parsedQuery["sort"].(map[string]interface{}); hasSort {
            column, hasColumn := sortData["column"].(string)
            order, hasOrder := sortData["order"].(string)
            
            if hasColumn {
                // Default to ascending if not specified
                ascending := true
                if hasOrder && strings.ToUpper(order) == "DESC" {
                    ascending = false
                }
                
                // Sort the results
                sort.Slice(results, func(i, j int) bool {
                    // Get values to compare
                    valI, hasI := results[i][column]
                    valJ, hasJ := results[j][column]
                    
                    // Handle missing values
                    if !hasI && !hasJ {
                        return false // Both missing, no change
                    } else if !hasI {
                        return ascending // I missing, I comes first if ascending
                    } else if !hasJ {
                        return !ascending // J missing, J comes first if ascending
                    }
                    
                    // Compare based on types
                    switch vi := valI.(type) {
                    case int64:
                        if vj, ok := valJ.(int64); ok {
                            if ascending {
                                return vi < vj
                            } else {
                                return vi > vj
                            }
                        }
                    case float64:
                        if vj, ok := valJ.(float64); ok {
                            if ascending {
                                return vi < vj
                            } else {
                                return vi > vj
                            }
                        }
                    case string:
                        if vj, ok := valJ.(string); ok {
                            if ascending {
                                return vi < vj
                            } else {
                                return vi > vj
                            }
                        }
                    }
                    
                    // If types don't match, convert to string and compare
                    strI := fmt.Sprintf("%v", valI)
                    strJ := fmt.Sprintf("%v", valJ)
                    if ascending {
                        return strI < strJ
                    } else {
                        return strI > strJ
                    }
                })
            }
        }
    }
    
    // Apply LIMIT if specified
    limit, hasLimit := parsedQuery["limit"].(int)
    if hasLimit && limit > 0 && limit < len(results) {
        results = results[:limit]
    }
    
    return results, nil
}

// evaluateCondition checks if a row matches a condition
func (e *Executor) evaluateCondition(condition string, row []string, headers []string, colMap map[string]int) bool {
    // Simple condition evaluator for LIKE and comparison operators
    condition = strings.TrimSpace(condition)
    
    // Handle LIKE operator
    if strings.Contains(condition, " LIKE ") {
        parts := strings.Split(condition, " LIKE ")
        if len(parts) != 2 {
            return true // Invalid condition, assume it passes
        }
        
        colName := strings.TrimSpace(parts[0])
        pattern := strings.TrimSpace(parts[1])
        pattern = strings.Trim(pattern, "'\"") // Remove quotes
        
        // Check if column exists
        colIdx, exists := colMap[colName]
        if !exists || colIdx >= len(row) {
            return true // Column doesn't exist, assume it passes
        }
        
        // Simple pattern matching (% is wildcard)
        value := row[colIdx]
        if pattern == "%" || pattern == "%%" {
            return true // Matches anything
        } else if strings.HasPrefix(pattern, "%") && strings.HasSuffix(pattern, "%") {
            substring := pattern[1 : len(pattern)-1]
            return strings.Contains(value, substring)
        } else if strings.HasPrefix(pattern, "%") {
            suffix := pattern[1:]
            return strings.HasSuffix(value, suffix)
        } else if strings.HasSuffix(pattern, "%") {
            prefix := pattern[:len(pattern)-1]
            return strings.HasPrefix(value, prefix)
        } else {
            return value == pattern
        }
    }
    
    // Handle comparison operators (=, >, <, >=, <=, !=)
    for _, op := range []string{">=", "<=", "!=", "=", ">", "<"} {
        if strings.Contains(condition, op) {
            parts := strings.Split(condition, op)
            if len(parts) != 2 {
                return true // Invalid condition, assume it passes
            }
            
            colName := strings.TrimSpace(parts[0])
            valStr := strings.TrimSpace(parts[1])
            valStr = strings.Trim(valStr, "'\"") // Remove quotes
            
            // Check if column exists
            colIdx, exists := colMap[colName]
            if !exists || colIdx >= len(row) {
                return true // Column doesn't exist, assume it passes
            }
            
            // Get the value from the row
            value := row[colIdx]
            
            // Try to parse numbers for numeric comparisons
            valFloat, valIsFloat := e.parseFloat(valStr)
            valueFloat, valueIsFloat := e.parseFloat(value)
            
            if valIsFloat && valueIsFloat {
                // Numeric comparison
                switch op {
                case "=":
                    return valueFloat == valFloat
                case "!=":
                    return valueFloat != valFloat
                case ">":
                    return valueFloat > valFloat
                case "<":
                    return valueFloat < valFloat
                case ">=":
                    return valueFloat >= valFloat
                case "<=":
                    return valueFloat <= valFloat
                }
            } else {
                // String comparison
                switch op {
                case "=":
                    return value == valStr
                case "!=":
                    return value != valStr
                case ">":
                    return value > valStr
                case "<":
                    return value < valStr
                case ">=":
                    return value >= valStr
                case "<=":
                    return value <= valStr
                }
            }
            
            break // Found and processed an operator
        }
    }
    
    // If we get here, condition format wasn't recognized
    return true // Assume it passes
}

// parseFloat attempts to parse a string as a float
func (e *Executor) parseFloat(s string) (float64, bool) {
    val, err := strconv.ParseFloat(s, 64)
    return val, err == nil
}