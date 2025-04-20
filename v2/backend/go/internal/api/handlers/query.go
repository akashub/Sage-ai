// // backend/go/internal/api/handlers/query.go
// package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"sage-ai-v2/internal/llm"
// 	"sage-ai-v2/internal/orchestrator"
// 	"sage-ai-v2/pkg/logger"
// )

// // QueryRequest struct for query processing
// type QueryRequest struct {
// 	Query   string `json:"query"`
// 	CSVPath string `json:"csvPath"`
// }

// // QueryResponse struct for query results
// type QueryResponse struct {
// 	Success        bool                      `json:"success"`
// 	GeneratedQuery string                    `json:"generatedQuery"`
// 	Results        []map[string]interface{}  `json:"results"`
// 	Message        string                    `json:"message,omitempty"`
// }

// // findLatestMatchingFile finds the most recent file in the uploads directory that contains the base filename
// func findLatestMatchingFile(requestPath string) (string, error) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get current directory: %w", err)
// 	}
	
// 	// Normalize the requested path
// 	baseFilename := filepath.Base(requestPath)
	
// 	// Define search paths
// 	searchPaths := []string{
// 		filepath.Join(cwd, "data", "uploads"),
// 		"./data/uploads",
// 		"/data/uploads",
// 		requestPath,
// 	}
	
// 	logger.InfoLogger.Printf("Looking for file containing '%s'", baseFilename)
// 	logger.InfoLogger.Printf("Current working directory: %s", cwd)
	
// 	// Try each search path
// 	for _, searchDir := range searchPaths {
// 		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
// 			logger.InfoLogger.Printf("Directory does not exist: %s", searchDir)
// 			continue
// 		}
		
// 		logger.InfoLogger.Printf("Searching in directory: %s", searchDir)
		
// 		// Read all files in the directory
// 		files, err := ioutil.ReadDir(searchDir)
// 		if err != nil {
// 			logger.InfoLogger.Printf("Failed to read directory %s: %v", searchDir, err)
// 			continue
// 		}
		
// 		// Find files that match the pattern (contain the base filename)
// 		var matchingFiles []os.FileInfo
// 		for _, file := range files {
// 			if !file.IsDir() && strings.Contains(file.Name(), baseFilename) {
// 				matchingFiles = append(matchingFiles, file)
// 				logger.InfoLogger.Printf("Found matching file: %s", file.Name())
// 			}
// 		}
		
// 		if len(matchingFiles) > 0 {
// 			// Sort by modification time (most recent first)
// 			latestFile := matchingFiles[0]
// 			latestTime := latestFile.ModTime()
			
// 			for _, file := range matchingFiles[1:] {
// 				if file.ModTime().After(latestTime) {
// 					latestFile = file
// 					latestTime = file.ModTime()
// 				}
// 			}
			
// 			// Return the full path to the latest matching file
// 			fullPath := filepath.Join(searchDir, latestFile.Name())
// 			logger.InfoLogger.Printf("Using latest matching file: %s", fullPath)
// 			return fullPath, nil
// 		}
// 	}
	
// 	// If we got here, we didn't find any matching files
// 	return "", fmt.Errorf("no matching files found for %s in any search directory", baseFilename)
// }

// // QueryHandler handles the processing of user queries against CSV data
// func QueryHandler(w http.ResponseWriter, r *http.Request) {
// 	// Set CORS headers
// 	w.Header().Set("Access-Control-Allow-Origin", "*") 
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	// Handle preflight request
// 	if r.Method == "OPTIONS" {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	// Check request method
// 	if r.Method != "POST" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse JSON request
// 	var req QueryRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		logger.ErrorLogger.Printf("Error parsing request: %v", err)
// 		http.Error(w, "Invalid request format", http.StatusBadRequest)
// 		return
// 	}

// 	// Validate request
// 	if req.Query == "" {
// 		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
// 		return
// 	}

// 	if req.CSVPath == "" {
// 		http.Error(w, "CSV path cannot be empty", http.StatusBadRequest)
// 		return
// 	}

// 	// Find the actual CSV file that matches the request
// 	csvFilename, err := findLatestMatchingFile(req.CSVPath)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Failed to find CSV file: %v", err)
// 		http.Error(w, fmt.Sprintf("Failed to find CSV file: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	// Initialize components
// 	bridge := llm.CreateBridge("http://localhost:8000") // Connect to Python microservice
// 	orch := orchestrator.CreateOrchestrator(bridge)
// 	orch.NewSession()

// 	// Create context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()

// 	// Process the query
// 	logger.InfoLogger.Printf("Processing query: %s with CSV: %s", req.Query, csvFilename)
// 	result, err := orch.ProcessQuery(ctx, req.Query, csvFilename)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error processing query: %v", err)
// 		http.Error(w, fmt.Sprintf("Error processing query: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Extract results from state
// 	var results []map[string]interface{}
// 	if resultsData, ok := result.ExecutionResult.([]map[string]interface{}); ok {
// 		results = resultsData
// 	} else {
// 		logger.InfoLogger.Printf("Results not in expected format: %T", result.ExecutionResult)
// 		results = []map[string]interface{}{}
// 	}

// 	// Prepare response
// 	response := QueryResponse{
// 		Success:        true,
// 		GeneratedQuery: result.GeneratedQuery,
// 		Results:        results,
// 		Message:        "Query processed successfully",
// 	}

// 	// Return JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		logger.ErrorLogger.Printf("Error encoding response: %v", err)
// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
// 		return
// 	}
// }

// backend/go/internal/api/handlers/query.go
// backend/go/internal/api/handlers/query.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"
)

// QueryRequest defines the request structure for the query endpoint
type QueryRequest struct {
	Query           string                 `json:"query"`
	CSVPath         string                 `json:"csvPath"`
	UseKnowledgeBase bool                  `json:"useKnowledgeBase"`
	Options         map[string]interface{} `json:"options,omitempty"`
	Timestamp       int64                  `json:"timestamp,omitempty"`
}

// QueryResponse defines the response structure for the query endpoint
type QueryResponse struct {
	Success         bool                     `json:"success"`
	SQL             string                   `json:"sql,omitempty"`
	Results         []map[string]interface{} `json:"results,omitempty"`
	Response        string                   `json:"response,omitempty"`
	KnowledgeContext []map[string]interface{} `json:"knowledgeContext,omitempty"`
	Error           string                   `json:"error,omitempty"`
}

// QueryHandler processes natural language queries and returns SQL results
func QueryHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight CORS request
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set headers for actual response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Parse request body
		var req QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.ErrorLogger.Printf("Error parsing query request: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		logger.InfoLogger.Printf("Received query: %s", req.Query)
		
		// Check if file exists or is properly specified
		if req.CSVPath == "" {
			http.Error(w, "CSV path is required", http.StatusBadRequest)
			return
		}

		// Try to find the actual file if it's just a name without path
		csvPath := req.CSVPath
		if !filepath.IsAbs(csvPath) && !containsPath(csvPath, "uploads") {
			// Check if it's just a filename and add uploads prefix
			if !strings.ContainsAny(csvPath, "/\\") {
				csvPath = filepath.Join("data", "uploads", csvPath)
			}
		}

		// Prepare options for orchestrator
		options := make(map[string]interface{})
		if req.Options != nil {
			for k, v := range req.Options {
				options[k] = v
			}
		}
		
		// Add knowledge base flag to options
		options["useKnowledgeBase"] = req.UseKnowledgeBase

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Process the query through the orchestrator
		state, err := orch.ProcessQueryWithOptions(ctx, req.Query, csvPath, options)
		if err != nil {
			logger.ErrorLogger.Printf("Query processing failed: %v", err)
			sendErrorResponse(w, "Query processing failed", http.StatusInternalServerError)
			return
		}

		// Extract knowledge context if available
		var knowledgeContext []map[string]interface{}
		if state.KnowledgeContext != nil {
			// Extract DDL schemas
			for _, schema := range state.KnowledgeContext.DDLSchemas {
				knowledgeContext = append(knowledgeContext, map[string]interface{}{
					"type":        "ddl",
					"content":     schema.Content,
					"description": schema.Description,
				})
			}
			
			// Extract documentation
			for _, doc := range state.KnowledgeContext.Documentation {
				knowledgeContext = append(knowledgeContext, map[string]interface{}{
					"type":        "documentation",
					"content":     doc.Content,
					"description": doc.Description,
				})
			}
			
			// Extract question-SQL pairs
			for _, pair := range state.KnowledgeContext.QuestionSQLPairs {
				knowledgeContext = append(knowledgeContext, map[string]interface{}{
					"type":        "question_sql",
					"question":    pair.Question,
					"sql":         pair.SQL,
					"description": pair.Description,
				})
			}
		}

		// Convert the execution result to the expected format
		var results []map[string]interface{}
		if state.ExecutionResult != nil {
			if resultArray, ok := state.ExecutionResult.([]map[string]interface{}); ok {
				results = resultArray
			} else if resultMap, ok := state.ExecutionResult.(map[string]interface{}); ok {
				// Single result case
				results = []map[string]interface{}{resultMap}
			}
		}

		// Generate a natural language response
		response := generateNaturalLanguageResponse(state.Query, state.GeneratedQuery, results)

		// Create the response
		resp := QueryResponse{
			Success:         true,
			SQL:             state.GeneratedQuery,
			Results:         results,
			Response:        response,
			KnowledgeContext: knowledgeContext,
		}

		// Return the JSON response
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			sendErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	resp := QueryResponse{
		Success: false,
		Error:   message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// containsPath checks if the target path contains a specific subpath
func containsPath(path, subpath string) bool {
	return filepath.Dir(path) == subpath || filepath.Base(filepath.Dir(path)) == subpath
}

// generateNaturalLanguageResponse creates a human-readable response
func generateNaturalLanguageResponse(query, sql string, results []map[string]interface{}) string {
	// For a simple query with no results
	if len(results) == 0 {
		return "I couldn't find any results matching your query."
	}
	
	// For a query with results
	resultCount := len(results)
	
	// Get a small sample of results to mention
	maxSampleSize := 3
	if resultCount < maxSampleSize {
		maxSampleSize = resultCount
	}
	
	// Get the first few results
	sample := results[:maxSampleSize]
	
	// Generate a response based on result count
	var response string
	
	switch {
	case resultCount == 1:
		response = fmt.Sprintf("I found 1 result for your query. ")
	case resultCount <= 10:
		response = fmt.Sprintf("I found %d results for your query. ", resultCount)
	default:
		response = fmt.Sprintf("I found %d results for your query. Here are the first few: ", resultCount)
	}
	
	// Add sample data to the response
	for i, result := range sample {
		// Get the first couple of key-value pairs to mention
		var details []string
		count := 0
		for k, v := range result {
			if count >= 2 { // Limit to 2 fields per result
				break
			}
			details = append(details, fmt.Sprintf("%s: %v", k, v))
			count++
		}
		
		if i < len(sample)-1 {
			response += fmt.Sprintf("%s; ", strings.Join(details, ", "))
		} else {
			response += fmt.Sprintf("%s. ", strings.Join(details, ", "))
		}
	}
	
	return response
}