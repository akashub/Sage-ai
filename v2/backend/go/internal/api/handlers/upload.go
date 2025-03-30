// // backend/go/internal/api/handlers/upload.go
// package handlers

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	"sage-ai-v2/pkg/logger"
// )

// // UploadResponse defines the response structure for upload endpoint
// type UploadResponse struct {
// 	Success  bool     `json:"success"`
// 	Filename string   `json:"filename"`
// 	Filepath string   `json:"filepath"`
// 	Headers  []string `json:"headers"`
// }

// // UploadFileHandler handles file uploads, specifically CSV files
// func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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

// 	// Parse multipart form (max 10MB)
// 	err := r.ParseMultipartForm(10 << 20)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error parsing multipart form: %v", err)
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// 	// Get file from form
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error retrieving file: %v", err)
// 		http.Error(w, "Error retrieving file", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Validate file type (must be CSV)
// 	if filepath.Ext(handler.Filename) != ".csv" {
// 		http.Error(w, "Only CSV files are allowed", http.StatusBadRequest)
// 		return
// 	}

// 	logger.InfoLogger.Printf("Received file: %s", handler.Filename)

// 	// Get current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error getting current directory: %v", err)
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.InfoLogger.Printf("Current working directory: %s", cwd)

// 	// Create uploads directory if it doesn't exist
// 	uploadsDir := filepath.Join(cwd, "data", "uploads")
// 	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
// 		logger.ErrorLogger.Printf("Error creating upload directory: %v", err)
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.InfoLogger.Printf("Upload directory: %s", uploadsDir)

// 	// Create unique filename to prevent overwrites
// 	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
// 	destPath := filepath.Join(uploadsDir, filename)
	
// 	// Create a simplified path to return to the frontend
// 	simplePath := filepath.Join("data", "uploads", filename)

// 	// Create the file
// 	dst, err := os.Create(destPath)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error creating file: %v", err)
// 		http.Error(w, "Error saving file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer dst.Close()

// 	// Copy the uploaded file to the created file
// 	if _, err := io.Copy(dst, file); err != nil {
// 		logger.ErrorLogger.Printf("Error copying file content: %v", err)
// 		http.Error(w, "Error saving file", http.StatusInternalServerError)
// 		return
// 	}

// 	logger.InfoLogger.Printf("File saved to: %s", destPath)

// 	// Read CSV headers
// 	headers, err := getCSVHeaders(destPath)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error reading CSV headers: %v", err)
// 		http.Error(w, "Error reading CSV headers", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if we can confirm the file exists where we think it is
// 	if _, err := os.Stat(destPath); os.IsNotExist(err) {
// 		logger.ErrorLogger.Printf("Warning: File not found at %s after saving", destPath)
// 	} else {
// 		logger.InfoLogger.Printf("Confirmed file exists at: %s", destPath)
// 	}

// 	// Create response object
// 	response := UploadResponse{
// 		Success:  true,
// 		Filename: filename,
// 		Filepath: simplePath, // Return a simplified path for the frontend
// 		Headers:  headers,
// 	}

// 	// Return success response with file information
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
	
// 	// Encode response to JSON
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		logger.ErrorLogger.Printf("Error encoding response: %v", err)
// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
// 		return
// 	}
// }

// // getCSVHeaders reads the first row of a CSV file and returns it as headers
// func getCSVHeaders(filepath string) ([]string, error) {
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening file: %w", err)
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	headers, err := reader.Read()
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading headers: %w", err)
// 	}

// 	return headers, nil
// }

// backend/go/internal/api/handlers/upload_fixed.go
package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"
)

// UploadFileHandler is the entry point for handling file uploads
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	HandleFileUpload(w, r)
}

// UploadResponse defines the response for file uploads
type UploadResponse struct {
	Success  bool     `json:"success"`
	FilePath string   `json:"filePath"`
	Headers  []string `json:"headers,omitempty"`
	Error    string   `json:"error,omitempty"`
}

// HandleFileUpload processes file uploads, particularly CSVs
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
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

	// Check request method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 50MB)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		logger.ErrorLogger.Printf("Error parsing multipart form: %v", err)
		sendUploadError(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.ErrorLogger.Printf("Error retrieving file: %v", err)
		sendUploadError(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".csv" {
		sendUploadError(w, "Only CSV files are supported", http.StatusBadRequest)
		return
	}

	logger.InfoLogger.Printf("Received file: %s, size: %d bytes", handler.Filename, handler.Size)

	// Create uploads directory if it doesn't exist
	uploadsDir := filepath.Join("data", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		logger.ErrorLogger.Printf("Error creating uploads directory: %v", err)
		sendUploadError(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create a unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		logger.ErrorLogger.Printf("Error creating file: %v", err)
		sendUploadError(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file
	if _, err := io.Copy(dst, file); err != nil {
		logger.ErrorLogger.Printf("Error copying file content: %v", err)
		sendUploadError(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	logger.InfoLogger.Printf("File saved to: %s", filePath)

	// Extract CSV headers
	headers, err := extractCSVHeaders(filePath)
	if err != nil {
		logger.ErrorLogger.Printf("Error extracting CSV headers: %v", err)
		// Continue despite header extraction failure
	}

	// Prepare response
	response := UploadResponse{
		Success:  true,
		FilePath: filename, // Return just the filename, not the full path
		Headers:  headers,
	}

	// Return JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorLogger.Printf("Error encoding response: %v", err)
		sendUploadError(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// extractCSVHeaders reads the first line of a CSV file to get column headers
func extractCSVHeaders(filePath string) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the first line which contains headers
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV headers: %w", err)
	}

	// Remove BOM if present in the first header
	if len(headers) > 0 && strings.HasPrefix(headers[0], "\ufeff") {
		headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
	}

	// Trim whitespace from headers
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	return headers, nil
}

// sendUploadError sends a JSON-formatted error response
func sendUploadError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := UploadResponse{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}