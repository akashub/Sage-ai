// backend/go/internal/api/handlers/upload.go
package handlers

import (
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type UploadHandler struct {
    uploadDir string
}

func CreateUploadHandler(uploadDir string) *UploadHandler {
    return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Validate file type
    if !isCSV(header.Filename) {
        http.Error(w, "Only CSV files are allowed", http.StatusBadRequest)
        return
    }

    // Create upload directory if it doesn't exist
    if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Create file path
    filename := filepath.Join(h.uploadDir, header.Filename)
    
    // Create the file
    dst, err := os.Create(filename)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy the uploaded file
    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "filepath": filename,
        "message": "File uploaded successfully",
    })
}

func isCSV(filename string) bool {
    return filepath.Ext(filename) == ".csv"
}