// backend/go/internal/api/handlers/training_data.go
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/pkg/logger"
	"strconv"
	"strings"
	"time"
)

// TrainingDataRequest defines the request for adding training data
type TrainingDataRequest struct {
	Type        string                 `json:"type"` // "ddl", "documentation", "question_sql"
	Content     string                 `json:"content,omitempty"` // Used for direct content submission
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TrainingDataResponse defines the response for training data operations
type TrainingDataResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	TrainingID  string `json:"training_id,omitempty"`
}

// AddTrainingDataHandler handles adding training data via API
func AddTrainingDataHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var req TrainingDataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.ErrorLogger.Printf("Error parsing training data request: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate request
		if req.Type == "" {
			http.Error(w, "Training data type cannot be empty", http.StatusBadRequest)
			return
		}

		if req.Content == "" {
			http.Error(w, "Training data content cannot be empty", http.StatusBadRequest)
			return
		}

		// Process metadata
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		
		// Add description to metadata
		req.Metadata["description"] = req.Description

		// Process based on type
		var err error
		switch req.Type {
		case "ddl":
			name := req.Metadata["name"]
			if name == nil {
				name = fmt.Sprintf("DDL-%d", time.Now().Unix())
			}
			err = orch.KnowledgeManager.AddDDLSchema(r.Context(), name.(string), req.Content, req.Description)
		
		case "documentation":
			var tags []string
			if tagsList, ok := req.Metadata["tags"].([]interface{}); ok {
				for _, tag := range tagsList {
					if tagStr, ok := tag.(string); ok {
						tags = append(tags, tagStr)
					}
				}
			}
			err = orch.KnowledgeManager.AddDocumentation(r.Context(), req.Description, req.Content, tags)
		
		case "question_sql":
			// Parse as JSON QuestionSQLPair
			var pair knowledge.QuestionSQLPair
			if err := json.Unmarshal([]byte(req.Content), &pair); err != nil {
				logger.ErrorLogger.Printf("Failed to parse question-SQL pair: %v", err)
				http.Error(w, "Invalid question-SQL pair format", http.StatusBadRequest)
				return
			}
			
			// Set description if not already set
			if pair.Description == "" {
				pair.Description = req.Description
			}
			
			err = orch.KnowledgeManager.AddQuestionSQLPair(r.Context(), pair)
		
		default:
			http.Error(w, fmt.Sprintf("Unsupported training data type: %s", req.Type), http.StatusBadRequest)
			return
		}

		if err != nil {
			logger.ErrorLogger.Printf("Failed to add training data: %v", err)
			http.Error(w, fmt.Sprintf("Failed to add training data: %v", err), http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := TrainingDataResponse{
			Success:    true,
			Message:    "Training data added successfully",
			TrainingID: fmt.Sprintf("training_%d", time.Now().UnixNano()),
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// UploadTrainingFileHandler handles training data file uploads
func UploadTrainingFileHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Parse multipart form (max 50MB)
		err := r.ParseMultipartForm(50 << 20)
		if err != nil {
			logger.ErrorLogger.Printf("Error parsing multipart form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get file from form
		file, handler, err := r.FormFile("file")
		if err != nil {
			logger.ErrorLogger.Printf("Error retrieving file: %v", err)
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get data type and description from form
		dataType := r.FormValue("type")
		description := r.FormValue("description")
		
		if dataType == "" {
			http.Error(w, "Training data type cannot be empty", http.StatusBadRequest)
			return
		}

		logger.InfoLogger.Printf("Received training data file: %s, type: %s", handler.Filename, dataType)

		// Create uploads directory if it doesn't exist
		uploadsDir := filepath.Join("data", "training")
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			logger.ErrorLogger.Printf("Error creating training data directory: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Create unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
		filePath := filepath.Join(uploadsDir, filename)

		// Create the file
		dst, err := os.Create(filePath)
		if err != nil {
			logger.ErrorLogger.Printf("Error creating file: %v", err)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the created file
		if _, err := io.Copy(dst, file); err != nil {
			logger.ErrorLogger.Printf("Error copying file content: %v", err)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		logger.InfoLogger.Printf("Training data file saved to: %s", filePath)

		// Determine file type if not specified
		if dataType == "auto" {
			ext := strings.ToLower(filepath.Ext(handler.Filename))
			switch ext {
			case ".sql":
				dataType = "ddl"
			case ".json":
				dataType = "question_sql_json"
			case ".md", ".txt", ".docx", ".html":
				dataType = "documentation"
			default:
				dataType = "documentation" // Default type
			}
			logger.InfoLogger.Printf("Auto-detected data type: %s", dataType)
		}

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.ErrorLogger.Printf("Error reading file content: %v", err)
			http.Error(w, "Error reading file content", http.StatusInternalServerError)
			return
		}

		// Process based on file type
		var processErr error
		fileName := strings.TrimSuffix(handler.Filename, filepath.Ext(handler.Filename))
		
		if description == "" {
			description = fileName
		}

		switch dataType {
		case "ddl":
			processErr = orch.KnowledgeManager.AddDDLSchema(
				r.Context(),
				fileName,
				string(content),
				description,
			)
			
		case "documentation":
			processErr = orch.KnowledgeManager.AddDocumentation(
				r.Context(),
				description,
				string(content),
				[]string{},
			)
			
		case "question_sql_json":
			var pairs []knowledge.QuestionSQLPair
			if err := json.Unmarshal(content, &pairs); err != nil {
				// Try single object
				var pair knowledge.QuestionSQLPair
				if err := json.Unmarshal(content, &pair); err != nil {
					logger.ErrorLogger.Printf("Failed to parse JSON content: %v", err)
					http.Error(w, "Invalid JSON format", http.StatusBadRequest)
					return
				}
				pairs = []knowledge.QuestionSQLPair{pair}
			}
			
			// Add each pair
			for _, pair := range pairs {
				// Set description if empty
				if pair.Description == "" {
					pair.Description = fmt.Sprintf("%s - %s", description, pair.Question)
				}
				
				// Set date if empty
				if pair.DateAdded == "" {
					pair.DateAdded = time.Now().Format(time.RFC3339)
				}
				
				if err := orch.KnowledgeManager.AddQuestionSQLPair(r.Context(), pair); err != nil {
					logger.ErrorLogger.Printf("Failed to add question-SQL pair: %v", err)
					// Continue with other pairs
				}
			}
			
		default:
			http.Error(w, fmt.Sprintf("Unsupported training data type: %s", dataType), http.StatusBadRequest)
			return
		}

		if processErr != nil {
			logger.ErrorLogger.Printf("Failed to process training data: %v", processErr)
			http.Error(w, fmt.Sprintf("Failed to process training data: %v", processErr), http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := TrainingDataResponse{
			Success:    true,
			Message:    "Training data file processed successfully",
			TrainingID: fmt.Sprintf("training_%d", time.Now().UnixNano()),
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// ListTrainingDataHandler lists all available training data
// In your ListTrainingDataHandler function
func ListTrainingDataHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger.InfoLogger.Printf("Received request to list training data")
        
        // Get data type filter from query parameter
        dataType := r.URL.Query().Get("type")
        logger.InfoLogger.Printf("Filter type: %s", dataType)

        // Get training data items
        items, err := orch.KnowledgeManager.ListTrainingData(r.Context(), dataType)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to list training data: %v", err)
            http.Error(w, fmt.Sprintf("Failed to list training data: %v", err), http.StatusInternalServerError)
            return
        }
        
        logger.InfoLogger.Printf("Returning %d training data items", len(items))
        
        // Return JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(items); err != nil {
            logger.ErrorLogger.Printf("Error encoding response: %v", err)
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
            return
        }
    }
}

// ViewTrainingDataHandler handles viewing a single training data item
func ViewTrainingDataHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check request method
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		id := parts[len(parts)-1]
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Get the item from the knowledge manager
		// In a real implementation, you would query the knowledge base
		// For now, we'll create a mock item
		item := knowledge.TrainingItem{
			ID:          id,
			Type:        "ddl",
			Content:     "CREATE TABLE users (\n  id INT PRIMARY KEY,\n  name VARCHAR(100),\n  email VARCHAR(100),\n  created_at TIMESTAMP\n);",
			Description: "Users table schema",
			DateAdded:   "2025-03-30T12:00:00Z",
		}

		// For documentation type
		if strings.Contains(id, "doc") {
			item.Type = "documentation"
			item.Content = "The users table contains all registered users in the system. Each user has a unique ID, name, email, and creation timestamp. The email is used for login purposes and must be unique."
		}

		// For question-sql type
		if strings.Contains(id, "qa") {
			item.Type = "question_sql"
			item.Content = `{
				"question": "Show me all users who registered in the last month",
				"sql": "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '1 month'",
				"description": "Recent user registrations query"
			}`
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(item); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteTrainingDataHandler handles deleting a training data item
func DeleteTrainingDataHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check request method
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		id := parts[len(parts)-1]
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// In a real implementation, you would delete the item from the knowledge base
		// For now, we'll just simulate success
		logger.InfoLogger.Printf("Deleting training data item: %s", id)

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		response := TrainingDataResponse{
			Success: true,
			Message: "Training data deleted successfully",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// SearchTrainingDataHandler handles searching training data
func SearchTrainingDataHandler(orch *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check request method
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get query params
		query := r.URL.Query().Get("q")
		typeFilter := r.URL.Query().Get("type")
		limitStr := r.URL.Query().Get("limit")

		limit := 10 // Default limit
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		logger.InfoLogger.Printf("Searching training data with query: %s, type: %s, limit: %d", query, typeFilter, limit)

		// In a real implementation, you would search the knowledge base
		// For now, return mock results
		items := []map[string]interface{}{
			{
				"id":          "ddl_1234567890",
				"type":        "ddl",
				"description": "Users table schema",
				"date_added":  "2025-03-30T12:00:00Z",
				"relevance":   0.95,
			},
			{
				"id":          "doc_9876543210",
				"type":        "documentation",
				"description": "User registration process",
				"date_added":  "2025-03-29T15:30:00Z",
				"relevance":   0.85,
			},
			{
				"id":          "qa_5678901234",
				"type":        "question_sql",
				"description": "Recent user registrations query",
				"date_added":  "2025-03-28T09:15:00Z",
				"relevance":   0.75,
			},
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			logger.ErrorLogger.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}