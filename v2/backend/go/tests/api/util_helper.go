	// backend/go/tests/api/util_helper.go
	package api

	import (
		"encoding/json"
		"errors"
		"io"
		"net/http"
	)

	// parsePenNameJSONRequest parses a JSON request body into a struct
	func parsePenNameJSONRequest(r *http.Request, v interface{}) error {
		if r.Body == nil {
			return errors.New("request body is empty")
		}
		
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		
		// Close the body
		defer r.Body.Close()
		
		// Parse the JSON
		return json.Unmarshal(body, v)
	}

	// writeJSON writes a JSON response
	func writeJSON(w http.ResponseWriter, v interface{}) error {
		// Marshal the response to JSON
		jsonData, err := json.Marshal(v)
		if err != nil {
			return err
		}
		
		// Set the content type
		w.Header().Set("Content-Type", "application/json")
		
		// Write the response
		_, err = w.Write(jsonData)
		return err
	}

	// ErrorResponse represents an error response
	type ErrorResponse struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
	}

	// WriteErrorResponse writes an error response in JSON format
	func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
		// Set the status code
		w.WriteHeader(statusCode)
		
		// Create the error response
		response := ErrorResponse{
			Error:   message,
			Success: false,
		}
		
		// Write the JSON response
		if err := writeJSON(w, response); err != nil {
			// If we can't write the JSON response, fallback to plain text
			http.Error(w, message, statusCode)
		}
	}