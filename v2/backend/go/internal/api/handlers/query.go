// // backend/go/internal/api/handlers/query.go
// package handlers

// import (
//     "encoding/json"
//     "net/http"
//     "sage-ai-v2/internal/orchestrator"
// )

// type QueryRequest struct {
//     Query    string `json:"query"`
//     CSVPath  string `json:"csv_path"`
// }

// type QueryHandler struct {
//     orchestrator *orchestrator.Orchestrator
// }

// func CreateQueryHandler(orch *orchestrator.Orchestrator) *QueryHandler {
//     return &QueryHandler{orchestrator: orch}
// }

// func (h *QueryHandler) Handle(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }

//     var req QueryRequest
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     result, err := h.orchestrator.ProcessQuery(r.Context(), req.Query, req.CSVPath)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(result)
// }
// backend/go/internal/api/handlers/query.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sage-ai-v2/internal/llm"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/pkg/logger"
)

// QueryRequest struct for query processing
type QueryRequest struct {
	Query   string `json:"query"`
	CSVPath string `json:"csvPath"`
}

// QueryResponse struct for query results
type QueryResponse struct {
	Success        bool                      `json:"success"`
	GeneratedQuery string                    `json:"generatedQuery"`
	Results        []map[string]interface{}  `json:"results"`
	Message        string                    `json:"message,omitempty"`
}

// findLatestMatchingFile finds the most recent file in the uploads directory that contains the base filename
func findLatestMatchingFile(requestPath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	
	// Normalize the requested path
	baseFilename := filepath.Base(requestPath)
	
	// Define search paths
	searchPaths := []string{
		filepath.Join(cwd, "data", "uploads"),
		"./data/uploads",
		"/data/uploads",
		requestPath,
	}
	
	logger.InfoLogger.Printf("Looking for file containing '%s'", baseFilename)
	logger.InfoLogger.Printf("Current working directory: %s", cwd)
	
	// Try each search path
	for _, searchDir := range searchPaths {
		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
			logger.InfoLogger.Printf("Directory does not exist: %s", searchDir)
			continue
		}
		
		logger.InfoLogger.Printf("Searching in directory: %s", searchDir)
		
		// Read all files in the directory
		files, err := ioutil.ReadDir(searchDir)
		if err != nil {
			logger.InfoLogger.Printf("Failed to read directory %s: %v", searchDir, err)
			continue
		}
		
		// Find files that match the pattern (contain the base filename)
		var matchingFiles []os.FileInfo
		for _, file := range files {
			if !file.IsDir() && strings.Contains(file.Name(), baseFilename) {
				matchingFiles = append(matchingFiles, file)
				logger.InfoLogger.Printf("Found matching file: %s", file.Name())
			}
		}
		
		if len(matchingFiles) > 0 {
			// Sort by modification time (most recent first)
			latestFile := matchingFiles[0]
			latestTime := latestFile.ModTime()
			
			for _, file := range matchingFiles[1:] {
				if file.ModTime().After(latestTime) {
					latestFile = file
					latestTime = file.ModTime()
				}
			}
			
			// Return the full path to the latest matching file
			fullPath := filepath.Join(searchDir, latestFile.Name())
			logger.InfoLogger.Printf("Using latest matching file: %s", fullPath)
			return fullPath, nil
		}
	}
	
	// If we got here, we didn't find any matching files
	return "", fmt.Errorf("no matching files found for %s in any search directory", baseFilename)
}

// QueryHandler handles the processing of user queries against CSV data
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check request method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorLogger.Printf("Error parsing request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	if req.CSVPath == "" {
		http.Error(w, "CSV path cannot be empty", http.StatusBadRequest)
		return
	}

	// Find the actual CSV file that matches the request
	csvFilename, err := findLatestMatchingFile(req.CSVPath)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to find CSV file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to find CSV file: %v", err), http.StatusBadRequest)
		return
	}

	// Initialize components
	bridge := llm.CreateBridge("http://localhost:8000") // Connect to Python microservice
	orch := orchestrator.CreateOrchestrator(bridge)
	orch.NewSession()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Process the query
	logger.InfoLogger.Printf("Processing query: %s with CSV: %s", req.Query, csvFilename)
	result, err := orch.ProcessQuery(ctx, req.Query, csvFilename)
	if err != nil {
		logger.ErrorLogger.Printf("Error processing query: %v", err)
		http.Error(w, fmt.Sprintf("Error processing query: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract results from state
	var results []map[string]interface{}
	if resultsData, ok := result.ExecutionResult.([]map[string]interface{}); ok {
		results = resultsData
	} else {
		logger.InfoLogger.Printf("Results not in expected format: %T", result.ExecutionResult)
		results = []map[string]interface{}{}
	}

	// Prepare response
	response := QueryResponse{
		Success:        true,
		GeneratedQuery: result.GeneratedQuery,
		Results:        results,
		Message:        "Query processed successfully",
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorLogger.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}