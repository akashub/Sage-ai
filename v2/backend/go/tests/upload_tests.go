// backend/go/tests/api/upload_test.go
package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/api/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleFileUpload tests the file upload handler with a real file
func TestHandleFileUpload(t *testing.T) {
	// Create a temporary directory for uploads
	tempDir, err := os.MkdirTemp("", "uploads-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create data directory structure
	uploadsDir := filepath.Join(tempDir, "data", "uploads")
	err = os.MkdirAll(uploadsDir, 0755)
	require.NoError(t, err)
	
	// Create a test CSV file
	csvContent := "id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com"
	csvPath := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)
	
	// Create a multipart form with the CSV file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add the file to the form
	part, err := writer.CreateFormFile("file", "test.csv")
	require.NoError(t, err)
	
	// Open the file and copy its contents to the form
	file, err := os.Open(csvPath)
	require.NoError(t, err)
	defer file.Close()
	
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	
	// Close the writer before sending the request
	err = writer.Close()
	require.NoError(t, err)
	
	// Create the request
	req := httptest.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler with patched environment
	// We need to temporarily set the UPLOAD_DIR environment variable
	originalUploadsDir := os.Getenv("UPLOAD_DIR")
	os.Setenv("UPLOAD_DIR", uploadsDir)
	defer os.Setenv("UPLOAD_DIR", originalUploadsDir)
	
	// Call the handler with our request and response recorder
	handlers.HandleFileUpload(rr, req)
	
	// Check response status
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Check response content
	assert.Equal(t, true, response["success"])
	assert.NotEmpty(t, response["filePath"])
	
	// Check that headers were extracted
	headers, ok := response["headers"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, headers, 3)
	assert.Equal(t, "id", headers[0])
	assert.Equal(t, "name", headers[1])
	assert.Equal(t, "email", headers[2])
	
	// Verify file was saved to the uploads directory
	filePath, ok := response["filePath"].(string)
	assert.True(t, ok)
	
	// Check file exists in uploads directory (just the filename part, which should be stored in the database)
	hasFile := false
	err = filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Base(path) == filepath.Base(filePath) {
			hasFile = true
		}
		return nil
	})
	require.NoError(t, err)
	
	assert.True(t, hasFile, "Uploaded file should exist in uploads directory")
}

// TestHandleFileUploadRejection tests the file upload handler with invalid input
func TestHandleFileUploadRejection(t *testing.T) {
	t.Run("NonCSVFile", func(t *testing.T) {
		// Create a non-CSV file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		// Add a non-CSV file to the form
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		
		_, err = part.Write([]byte("This is not a CSV file"))
		require.NoError(t, err)
		
		// Close the writer before sending the request
		err = writer.Close()
		require.NoError(t, err)
		
		// Create the request
		req := httptest.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		// Create response recorder
		rr := httptest.NewRecorder()
		
		// Call the handler
		handlers.HandleFileUpload(rr, req)
		
		// Should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		
		// Parse response body
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Verify rejection reason
		assert.Equal(t, false, response["success"])
		assert.Contains(t, response["error"], "Only CSV files are supported")
	})
	
	t.Run("MissingFile", func(t *testing.T) {
		// Create a request with no file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		// Close the writer without adding a file
		err := writer.Close()
		require.NoError(t, err)
		
		// Create the request
		req := httptest.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		// Create response recorder
		rr := httptest.NewRecorder()
		
		// Call the handler
		handlers.HandleFileUpload(rr, req)
		
		// Should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		
		// Parse response body
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Verify rejection reason
		assert.Equal(t, false, response["success"])
		assert.Contains(t, response["error"], "file")
	})
	
	t.Run("NonMultipartForm", func(t *testing.T) {
		// Create a request with JSON body instead of multipart form
		body := bytes.NewBufferString(`{"file": "not-a-file"}`)
		
		// Create the request
		req := httptest.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Content-Type", "application/json")
		
		// Create response recorder
		rr := httptest.NewRecorder()
		
		// Call the handler
		handlers.HandleFileUpload(rr, req)
		
		// Should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		
		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Verify rejection reason
		assert.Equal(t, false, response["success"])
	})
}