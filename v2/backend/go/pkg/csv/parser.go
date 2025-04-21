// // backend/go/pkg/csv/parser.go
// package csv

// import (
//     "encoding/csv"
//     "os"
//     "strings"
// )

// type Parser struct {
//     cache map[string][][]string
// }

// func CreateParser() *Parser {
//     return &Parser{
//         cache: make(map[string][][]string),
//     }
// }

// type Operation struct {
//     Type   string
//     Params map[string]interface{}
// }

// func (p *Parser) ExecuteQuery(csvPath string, query string) (interface{}, error) {
//     // Read CSV if not in cache
//     data, err := p.readCSV(csvPath)
//     if err != nil {
//         return nil, err
//     }

//     // Execute the query on the data
//     result, err := p.processQuery(data, query)
//     if err != nil {
//         return nil, err
//     }

//     return result, nil
// }

// func (p *Parser) readCSV(path string) ([][]string, error) {
//     // Check cache first
//     if data, exists := p.cache[path]; exists {
//         return data, nil
//     }

//     // Open file
//     file, err := os.Open(path)
//     if err != nil {
//         return nil, err
//     }
//     defer file.Close()

//     // Read CSV
//     reader := csv.NewReader(file)
//     data, err := reader.ReadAll()
//     if err != nil {
//         return nil, err
//     }

//     // Cache the data
//     p.cache[path] = data
//     return data, nil
// }

// func (p *Parser) processQuery(data [][]string, query string) (interface{}, error) {
//     // Convert the query into operations
//     ops := p.parseQuery(query)

//     // Process operations
//     result := p.executeOperations(data, ops)

//     return result, nil
// }

// func (p *Parser) parseQuery(query string) []Operation {
//     // Parse the query string into operations
//     var ops []Operation

//     // Example parsing
//     if strings.Contains(query, "SELECT") {
//         ops = append(ops, Operation{
//             Type: "select",
//             Params: map[string]interface{}{
//                 "columns": []string{"*"},
//             },
//         })
//     }

//     return ops
// }

// func (p *Parser) executeOperations(data [][]string, ops []Operation) interface{} {
//     result := data

//     for _, op := range ops {
//         switch op.Type {
//         case "select":
//             result = p.executeSelect(result, op.Params)
//         case "filter":
//             result = p.executeFilter(result, op.Params)
//         case "aggregate":
//             result = p.executeAggregate(result, op.Params)
//         }
//     }

//     return result
// }

// func (p *Parser) executeSelect(data [][]string, _ map[string]interface{}) [][]string {
//     // Implementation for SELECT operations
//     return data
// }

// func (p *Parser) executeFilter(data [][]string, _ map[string]interface{}) [][]string {
//     // Implementation for WHERE conditions
//     return data
// }

// func (p *Parser) executeAggregate(data [][]string, _ map[string]interface{}) [][]string {
//     // Implementation for GROUP BY and aggregate functions
//     return data
// }

// backend/go/pkg/csv/parser.go
// backend/go/pkg/csv/parser.go
// package csv

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"sage-ai-v2/pkg/logger"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// type Parser struct {
// 	sessionID string
// 	cache     map[string]map[string][][]string // sessionID -> filePath -> data
// }

// func CreateParser() *Parser {
// 	return &Parser{
// 		sessionID: generateSessionID(),
// 		cache:     make(map[string]map[string][][]string),
// 	}
// }

// func (p *Parser) NewSession() {
// 	p.sessionID = generateSessionID()
// 	// Clear cache for previous session
// 	delete(p.cache, p.sessionID)
// }

// type QueryInfo struct {
// 	SelectColumns []string
// 	OrderBy       string
// 	OrderDesc     bool
// 	Limit         int
// 	Conditions    []string
// }


// func generateSessionID() string {
// 	return fmt.Sprintf("session_%d", time.Now().UnixNano())
// }

// func (p *Parser) readCSV(path string) ([][]string, error) {
// 	// Check session cache first
// 	if sessionCache, ok := p.cache[p.sessionID]; ok {
// 		if data, ok := sessionCache[path]; ok {
// 			return data, nil
// 		}
// 	}

// 	// Initialize session cache if needed
// 	if _, ok := p.cache[p.sessionID]; !ok {
// 		p.cache[p.sessionID] = make(map[string][][]string)
// 	}

// 	// Read and cache the data for this session
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	data, err := reader.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	p.cache[p.sessionID][path] = data
// 	return data, nil
// }



// func (p *Parser) ExecuteQuery(csvPath string, query string) (interface{}, error) {
// 	logger.InfoLogger.Printf("CSV Parser: Starting query execution")
// 	logger.InfoLogger.Printf("Query: %s", query)

// 	data, err := p.readCSV(csvPath)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("CSV Parser: Failed to read CSV: %v", err)
// 		return nil, fmt.Errorf("error reading CSV: %w", err)
// 	}
// 	logger.InfoLogger.Printf("CSV Parser: Successfully read %d rows", len(data))

// 	if len(data) < 2 {
// 		logger.ErrorLogger.Printf("CSV Parser: Insufficient data in CSV")
// 		return nil, fmt.Errorf("insufficient data in CSV")
// 	}

// 	headers := data[0]
// 	logger.DebugLogger.Printf("CSV Headers: %v", headers)

// 	queryInfo := p.parseQuery(query, headers)
// 	logger.DebugLogger.Printf("Parsed Query Info: %+v", queryInfo)

// 	results := p.processQuery(headers, data[1:], queryInfo)
// 	logger.InfoLogger.Printf("CSV Parser: Query execution completed, found %d results", len(results))

// 	return results, nil
// }

// func (p *Parser) parseQuery(query string, headers []string) QueryInfo {
// 	logger.InfoLogger.Printf("CSV Parser: Parsing query: %s", query)

// 	selectStr := strings.ToLower(query)
// 	selectColumns := []string{}
// 	orderBy := "revenue" // Default
// 	var conditions []string

// 	limit := -1

// 	// Extract WHERE conditions
// 	if strings.Contains(selectStr, "where") {
// 		parts := strings.Split(selectStr, "where")
// 		if len(parts) > 1 {
// 			wherePart := strings.Split(parts[1], "order by")[0]
// 			wherePart = strings.Split(wherePart, "limit")[0]
// 			wherePart = strings.TrimSpace(wherePart)
// 			conditions = append(conditions, wherePart)
// 			logger.DebugLogger.Printf("CSV Parser: Found WHERE condition: %s", wherePart)
// 		}
// 	}

// 	// Extract LIMIT
// 	if strings.Contains(selectStr, "limit") {
// 		parts := strings.Split(selectStr, "limit")
// 		if len(parts) > 1 {
// 			limitStr := strings.TrimSpace(strings.Split(parts[1], ";")[0]) // Handle semicolon at end
// 			if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
// 				limit = parsedLimit
// 				logger.DebugLogger.Printf("CSV Parser: Found LIMIT: %d", limit)
// 			}
// 		}
// 	}

// 	// If no limit was found in query, set default
// 	if limit == -1 {
// 		limit = 10 // Only set default if no LIMIT clause was found
// 	}

// 	// Extract columns
// 	if strings.Contains(selectStr, "select") {
// 		selectPart := strings.Split(strings.Split(selectStr, "from")[0], "select")[1]
// 		columns := strings.Split(selectPart, ",")
// 		for _, col := range columns {
// 			col = strings.TrimSpace(col)
// 			// Match column names case-insensitively
// 			for _, header := range headers {
// 				if strings.ToLower(header) == col {
// 					selectColumns = append(selectColumns, header)
// 				}
// 			}
// 		}
// 	}

// 	logger.DebugLogger.Printf("CSV Parser: Selected columns: %v", selectColumns)

// 	return QueryInfo{
// 		SelectColumns: selectColumns,
// 		OrderBy:       orderBy,
// 		OrderDesc:     true,
// 		Limit:         limit,
// 		Conditions:    conditions,
// 	}
// }

// func (p *Parser) processQuery(headers []string, rows [][]string, query QueryInfo) []map[string]interface{} {
// 	logger.InfoLogger.Printf("CSV Parser: Processing query with limit %d", query.Limit)

// 	// Get column indexes
// 	colIndexes := make(map[string]int)
// 	for i, header := range headers {
// 		colIndexes[header] = i
// 	}

// 	// Convert rows to records
// 	var records []map[string]interface{}
// 	for _, row := range rows {
// 		if len(row) != len(headers) {
// 			continue // Skip malformed rows
// 		}

// 		record := make(map[string]interface{})

// 		// Check if this row matches the WHERE conditions
// 		matches := true
// 		for _, condition := range query.Conditions {
// 			if strings.Contains(strings.ToLower(condition), "like '%action%'") {
// 				if idx, ok := colIndexes["genres"]; ok {
// 					genres := strings.ToLower(row[idx])
// 					if !strings.Contains(genres, "action") {
// 						matches = false
// 						break
// 					}
// 				}
// 			} else if strings.Contains(strings.ToLower(condition), "like '%horror%'") {
// 				if idx, ok := colIndexes["genres"]; ok {
// 					genres := strings.ToLower(row[idx])
// 					if !strings.Contains(genres, "horror") {
// 						matches = false
// 						break
// 					}
// 				}
// 			}
// 			// Add more condition types as needed
// 		}

// 		if !matches {
// 			continue
// 		}

// 		// Add all requested columns
// 		for _, col := range query.SelectColumns {
// 			if idx, ok := colIndexes[col]; ok && idx < len(row) {
// 				value := row[idx]
// 				switch col {
// 				case "revenue":
// 					if rev, err := strconv.ParseFloat(value, 64); err == nil {
// 						record[col] = rev
// 					}
// 				case "vote_average":
// 					if rating, err := strconv.ParseFloat(value, 64); err == nil {
// 						record[col] = rating
// 					}
// 				default:
// 					record[col] = value
// 				}
// 			}
// 		}

// 		// Only add records with valid values
// 		if _, hasRequired := record[query.OrderBy]; hasRequired {
// 			records = append(records, record)
// 		}
// 	}

// 	logger.InfoLogger.Printf("CSV Parser: Found %d matching records before sorting", len(records))

// 	// Sort by specified column
// 	sort.Slice(records, func(i, j int) bool {
// 		valI, okI := records[i][query.OrderBy].(float64)
// 		valJ, okJ := records[j][query.OrderBy].(float64)

// 		if !okI || !okJ {
// 			return false
// 		}

// 		if query.OrderDesc {
// 			return valI > valJ
// 		}
// 		return valI < valJ
// 	})

// 	// Apply limit
// 	if query.Limit > 0 && query.Limit < len(records) {
// 		records = records[:query.Limit]
// 	}

// 	// Format numerical values for display
// 	for _, record := range records {
// 		if revenue, ok := record["revenue"].(float64); ok {
// 			record["revenue"] = fmt.Sprintf("$%.2f Million", revenue/1000000)
// 		}
// 		if rating, ok := record["vote_average"].(float64); ok {
// 			record["vote_average"] = fmt.Sprintf("%.1f", rating)
// 		}
// 	}

// 	logger.InfoLogger.Printf("CSV Parser: Returning %d records after processing", len(records))
// 	return records
// }

// v2/backend/go/pkg/csv/parser.go
package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sage-ai-v2/pkg/logger"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	sessionID string
	cache     map[string]map[string][][]string // sessionID -> filePath -> data
}

func CreateParser() *Parser {
	return &Parser{
		sessionID: generateSessionID(),
		cache:     make(map[string]map[string][][]string),
	}
}

func (p *Parser) NewSession() {
	p.sessionID = generateSessionID()
	// Clear cache for previous session
	delete(p.cache, p.sessionID)
}

type QueryInfo struct {
	SelectColumns []string
	OrderBy       string
	OrderDesc     bool
	Limit         int
	Conditions    []string
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// findCSVFile tries multiple path variants to find the CSV file
func (p *Parser) findCSVFile(csvPath string) (string, error) {
	// Variant paths to try
	variants := []string{
		csvPath,
		filepath.Clean(csvPath),
	}

	// If the path doesn't have a leading ".", add variants with "./"
	if !strings.HasPrefix(csvPath, "./") && !strings.HasPrefix(csvPath, "/") {
		variants = append(variants, 
			"./"+csvPath,
			filepath.Join("./", csvPath),
		)
	}

	// Try paths relative to working directory
	cwd, _ := os.Getwd()
	variants = append(variants,
		filepath.Join(cwd, csvPath),
		filepath.Join(cwd, "data", "uploads", filepath.Base(csvPath)),
	)

	// Log the variants we're trying
	logger.InfoLogger.Printf("CSV Parser: Trying to find CSV file with these paths: %v", variants)

	// Try each variant
	for _, path := range variants {
		logger.DebugLogger.Printf("CSV Parser: Checking if file exists at: %s", path)
		if _, err := os.Stat(path); err == nil {
			logger.InfoLogger.Printf("CSV Parser: Found CSV file at: %s", path)
			return path, nil
		}
	}

	return "", fmt.Errorf("CSV file not found in any of these locations: %v", variants)
}

func (p *Parser) readCSV(path string) ([][]string, error) {
	// Check session cache first
	if sessionCache, ok := p.cache[p.sessionID]; ok {
		if data, ok := sessionCache[path]; ok {
			logger.InfoLogger.Printf("CSV Parser: Using cached data for %s", path)
			return data, nil
		}
	}

	// Initialize session cache if needed
	if _, ok := p.cache[p.sessionID]; !ok {
		p.cache[p.sessionID] = make(map[string][][]string)
	}

	// Find the actual file path
	resolvedPath, err := p.findCSVFile(path)
	if err != nil {
		logger.ErrorLogger.Printf("CSV Parser: File not found: %v", err)
		return nil, err
	}

	// Read and cache the data for this session
	file, err := os.Open(resolvedPath)
	if err != nil {
		logger.ErrorLogger.Printf("CSV Parser: Failed to open file: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Printf("CSV Parser: Failed to read CSV data: %v", err)
		return nil, err
	}

	p.cache[p.sessionID][path] = data
	logger.InfoLogger.Printf("CSV Parser: Successfully read and cached %d rows from %s", len(data), resolvedPath)
	return data, nil
}

func (p *Parser) ExecuteQuery(csvPath string, query string) (interface{}, error) {
	logger.InfoLogger.Printf("CSV Parser: Starting query execution")
	logger.InfoLogger.Printf("CSV Parser: Query: %s", query)
	logger.InfoLogger.Printf("CSV Parser: CSV Path: %s", csvPath)

	data, err := p.readCSV(csvPath)
	if err != nil {
		logger.ErrorLogger.Printf("CSV Parser: Failed to read CSV: %v", err)
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}
	logger.InfoLogger.Printf("CSV Parser: Successfully read %d rows", len(data))

	if len(data) < 2 {
		logger.ErrorLogger.Printf("CSV Parser: Insufficient data in CSV")
		return nil, fmt.Errorf("insufficient data in CSV")
	}

	headers := data[0]
	logger.DebugLogger.Printf("CSV Parser: Headers: %v", headers)

	queryInfo := p.parseQuery(query, headers)
	logger.DebugLogger.Printf("CSV Parser: Parsed Query Info: %+v", queryInfo)

	results := p.processQuery(headers, data[1:], queryInfo)
	logger.InfoLogger.Printf("CSV Parser: Query execution completed, found %d results", len(results))

	return results, nil
}

func (p *Parser) parseQuery(query string, headers []string) QueryInfo {
	logger.InfoLogger.Printf("CSV Parser: Parsing query: %s", query)

	selectStr := strings.ToLower(query)
	selectColumns := []string{}
	orderBy := "revenue" // Default
	var conditions []string

	limit := -1

	// Extract WHERE conditions
	if strings.Contains(selectStr, "where") {
		parts := strings.Split(selectStr, "where")
		if len(parts) > 1 {
			wherePart := strings.Split(parts[1], "order by")[0]
			wherePart = strings.Split(wherePart, "limit")[0]
			wherePart = strings.TrimSpace(wherePart)
			conditions = append(conditions, wherePart)
			logger.DebugLogger.Printf("CSV Parser: Found WHERE condition: %s", wherePart)
		}
	}

	// Extract LIMIT
	if strings.Contains(selectStr, "limit") {
		parts := strings.Split(selectStr, "limit")
		if len(parts) > 1 {
			limitStr := strings.TrimSpace(strings.Split(parts[1], ";")[0]) // Handle semicolon at end
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
				limit = parsedLimit
				logger.DebugLogger.Printf("CSV Parser: Found LIMIT: %d", limit)
			}
		}
	}

	// If no limit was found in query, set default
	if limit == -1 {
		limit = 10 // Only set default if no LIMIT clause was found
	}

	// Extract columns
	if strings.Contains(selectStr, "select") {
		selectPart := strings.Split(strings.Split(selectStr, "from")[0], "select")[1]
		columns := strings.Split(selectPart, ",")
		for _, col := range columns {
			col = strings.TrimSpace(col)
			// Match column names case-insensitively
			for _, header := range headers {
				if strings.ToLower(header) == col {
					selectColumns = append(selectColumns, header)
				}
			}
		}
	}
	
	// If no columns were selected, select all columns
	if len(selectColumns) == 0 {
		selectColumns = headers
		logger.DebugLogger.Printf("CSV Parser: No specific columns selected, using all headers")
	}

	logger.DebugLogger.Printf("CSV Parser: Selected columns: %v", selectColumns)

	return QueryInfo{
		SelectColumns: selectColumns,
		OrderBy:       orderBy,
		OrderDesc:     true,
		Limit:         limit,
		Conditions:    conditions,
	}
}

func (p *Parser) processQuery(headers []string, rows [][]string, query QueryInfo) []map[string]interface{} {
	logger.InfoLogger.Printf("CSV Parser: Processing query with limit %d", query.Limit)

	// Get column indexes
	colIndexes := make(map[string]int)
	for i, header := range headers {
		colIndexes[header] = i
	}

	// Convert rows to records
	var records []map[string]interface{}
	for _, row := range rows {
		if len(row) != len(headers) {
			continue // Skip malformed rows
		}

		record := make(map[string]interface{})

		// Check if this row matches the WHERE conditions
		matches := true
		for _, condition := range query.Conditions {
			if strings.Contains(strings.ToLower(condition), "like '%action%'") {
				if idx, ok := colIndexes["genres"]; ok {
					genres := strings.ToLower(row[idx])
					if !strings.Contains(genres, "action") {
						matches = false
						break
					}
				}
			} else if strings.Contains(strings.ToLower(condition), "like '%horror%'") {
				if idx, ok := colIndexes["genres"]; ok {
					genres := strings.ToLower(row[idx])
					if !strings.Contains(genres, "horror") {
						matches = false
						break
					}
				}
			}
			// Add more condition types as needed
		}

		if !matches {
			continue
		}

		// Add all requested columns
		for _, col := range query.SelectColumns {
			if idx, ok := colIndexes[col]; ok && idx < len(row) {
				value := row[idx]
				switch col {
				case "revenue":
					if rev, err := strconv.ParseFloat(value, 64); err == nil {
						record[col] = rev
					} else {
						record[col] = value // Keep as string if not a valid number
					}
				case "vote_average":
					if rating, err := strconv.ParseFloat(value, 64); err == nil {
						record[col] = rating
					} else {
						record[col] = value // Keep as string if not a valid number
					}
				default:
					record[col] = value
				}
			}
		}

		// Add record if it has any data
		if len(record) > 0 {
			records = append(records, record)
		}
	}

	logger.InfoLogger.Printf("CSV Parser: Found %d matching records before sorting", len(records))

	// Sort by specified column if possible
	orderByColumn := query.OrderBy
	if _, hasOrderColumn := colIndexes[orderByColumn]; hasOrderColumn {
		sort.Slice(records, func(i, j int) bool {
			// Try numeric comparison first
			valI, okI := records[i][orderByColumn].(float64)
			valJ, okJ := records[j][orderByColumn].(float64)

			if okI && okJ {
				// Numeric comparison
				if query.OrderDesc {
					return valI > valJ
				}
				return valI < valJ
			}

			// Fallback to string comparison
			strI, _ := records[i][orderByColumn].(string)
			strJ, _ := records[j][orderByColumn].(string)
			
			if query.OrderDesc {
				return strI > strJ
			}
			return strI < strJ
		})
	} else {
		logger.InfoLogger.Printf("CSV Parser: Order by column '%s' not found, skipping sort", orderByColumn)
	}

	// Apply limit
	if query.Limit > 0 && query.Limit < len(records) {
		records = records[:query.Limit]
	}

	// Format numerical values for display
	for _, record := range records {
		if revenue, ok := record["revenue"].(float64); ok {
			record["revenue"] = fmt.Sprintf("$%.2f Million", revenue/1000000)
		}
		if rating, ok := record["vote_average"].(float64); ok {
			record["vote_average"] = fmt.Sprintf("%.1f", rating)
		}
	}

	logger.InfoLogger.Printf("CSV Parser: Returning %d records after processing", len(records))
	return records
}