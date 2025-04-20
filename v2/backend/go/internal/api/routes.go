// backend/go/internal/api/routes.go
package api

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/api/middleware"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the API routes
func SetupRoutes(db *sql.DB, orch *orchestrator.Orchestrator) *mux.Router {
    // Create router
    router := mux.NewRouter()
    
    // Apply middleware
    router.Use(middleware.CORSMiddleware)
    
    // Extract knowledge manager from orchestrator
    km := orch.KnowledgeManager

    // Create uploads directory if it doesn't exist
    os.MkdirAll("data/uploads", 0755)

    // API routes
    router.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
        handleQueryRequest(w, r, orch)
    }).Methods("POST", "OPTIONS")

    router.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
        handleFileUpload(w, r)
    }).Methods("POST", "OPTIONS")

    // Training Data API endpoints - IMPORTANT FIX
    router.HandleFunc("/api/training/list", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        dataType := r.URL.Query().Get("type")
        items, err := km.ListTrainingData(r.Context(), dataType)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(items)
    }).Methods("GET", "OPTIONS")
    
    router.HandleFunc("/api/training/upload", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        err := r.ParseMultipartForm(32 << 20)
        if err != nil {
            http.Error(w, "Could not parse form", http.StatusBadRequest)
            return
        }
        
        file, handler, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "No file provided", http.StatusBadRequest)
            return
        }
        defer file.Close()
        
        content, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "Failed to read file", http.StatusInternalServerError)
            return
        }
        
        localPath, err := km.StoreFile(handler.Filename, content)
        if err != nil {
            http.Error(w, "Failed to store file", http.StatusInternalServerError)
            return
        }
        
        dataType := r.FormValue("type")
        if dataType == "" {
            dataType = "ddl"
        }
        
        description := r.FormValue("description")
        if description == "" {
            description = handler.Filename
        }
        
        id := fmt.Sprintf("%s_%d", dataType, time.Now().UnixNano())
        ctx := r.Context()
        
        switch dataType {
        case "ddl":
            err = km.AddDDLSchema(ctx, handler.Filename, string(content), description)
        case "documentation":
            err = km.AddDocumentation(ctx, description, string(content), []string{})
        case "question_sql_json":
            pairs, err := km.LoadQuestionSQLPairsFromJSON(ctx, localPath)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to load JSON: %v", err), http.StatusInternalServerError)
                return
            }
            
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": true,
                "count":   pairs,
                "path":    localPath,
                "id":      id,
            })
            return
        case "auto":
            if strings.HasSuffix(handler.Filename, ".sql") {
                err = km.AddDDLSchema(ctx, handler.Filename, string(content), description)
            } else if strings.HasSuffix(handler.Filename, ".json") {
                pairs, err := km.LoadQuestionSQLPairsFromJSON(ctx, localPath)
                if err != nil {
                    http.Error(w, fmt.Sprintf("Failed to load JSON: %v", err), http.StatusInternalServerError)
                    return
                }
                
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(map[string]interface{}{
                    "success": true,
                    "count":   pairs,
                    "path":    localPath,
                    "id":      id,
                })
                return
            } else {
                err = km.AddDocumentation(ctx, description, string(content), []string{})
            }
        default:
            http.Error(w, "Invalid data type", http.StatusBadRequest)
            return
        }
        
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to add training data: %v", err), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": true,
            "type":    dataType,
            "path":    localPath,
            "id":      id,
        })
    }).Methods("POST", "OPTIONS")
    
    router.HandleFunc("/api/training/add", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        var data struct {
            Type        string                 `json:"type"`
            Content     string                 `json:"content"`
            Description string                 `json:"description"`
            Metadata    map[string]interface{} `json:"metadata"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        id := fmt.Sprintf("%s_%d", data.Type, time.Now().UnixNano())
        ctx := r.Context()
        var err error
        
        switch data.Type {
        case "ddl":
            err = km.AddDDLSchema(ctx, data.Description, data.Content, data.Description)
        case "documentation":
            tags := []string{}
            if tagsData, ok := data.Metadata["tags"].([]interface{}); ok {
                for _, tag := range tagsData {
                    if tagStr, ok := tag.(string); ok {
                        tags = append(tags, tagStr)
                    }
                }
            }
            err = km.AddDocumentation(ctx, data.Description, data.Content, tags)
        case "question_sql":
            pair := knowledge.QuestionSQLPair{
                Question:    data.Content,
                SQL:         "",
                Description: data.Description,
                DateAdded:   time.Now().Format(time.RFC3339),
                Verified:    true,
            }
            
            if sqlData, ok := data.Metadata["sql"].(string); ok {
                pair.SQL = sqlData
            }
            
            if tagsData, ok := data.Metadata["tags"].([]interface{}); ok {
                for _, tag := range tagsData {
                    if tagStr, ok := tag.(string); ok {
                        pair.Tags = append(pair.Tags, tagStr)
                    }
                }
            }
            
            err = km.AddQuestionSQLPair(ctx, pair)
        default:
            http.Error(w, "Invalid data type", http.StatusBadRequest)
            return
        }
        
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to add training data: %v", err), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": true,
            "type":    data.Type,
            "id":      id,
        })
    }).Methods("POST", "OPTIONS")
    
    router.HandleFunc("/api/training/view/{id}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        vars := mux.Vars(r)
        id := vars["id"]
        
        if id == "" {
            http.Error(w, "ID is required", http.StatusBadRequest)
            return
        }
        
        item, err := km.GetTrainingItem(r.Context(), id)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to get training item: %v", err)
            http.Error(w, fmt.Sprintf("Failed to get training item: %v", err), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    }).Methods("GET", "OPTIONS")
    
    router.HandleFunc("/api/training/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        vars := mux.Vars(r)
        id := vars["id"]
        
        if id == "" {
            http.Error(w, "ID is required", http.StatusBadRequest)
            return
        }
        
        err := km.DeleteTrainingItem(r.Context(), id)
        if err != nil {
            logger.ErrorLogger.Printf("Failed to delete training item: %v", err)
            http.Error(w, fmt.Sprintf("Failed to delete training item: %v", err), http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusNoContent)
    }).Methods("DELETE", "OPTIONS")
    
    // Chat History API
    SetupChatRoutes(router)

    // Auth routes
    router.HandleFunc("/api/auth/signin", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "user": map[string]interface{}{
                "id": "user123",
                "email": req.Email,
                "name": "Test User",
                "createdAt": time.Now(),
                "lastLoginAt": time.Now(),
            },
            "accessToken": "mock_token_" + req.Email,
        })
    }).Methods("POST", "OPTIONS")
    
    router.HandleFunc("/api/auth/signup", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            Name     string `json:"name"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "user": map[string]interface{}{
                "id": "user123",
                "email": req.Email,
                "name": req.Name,
                "createdAt": time.Now(),
                "lastLoginAt": time.Now(),
            },
            "accessToken": "mock_token_" + req.Email,
        })
    }).Methods("POST", "OPTIONS")
    
    router.HandleFunc("/api/auth/signout", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{
            "success": true,
        })
    }).Methods("POST", "OPTIONS")
    
    // OAuth URL endpoint - CRITICAL FIX
    router.HandleFunc("/api/auth/oauth/url/{provider}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        vars := mux.Vars(r)
        provider := vars["provider"]
        
        redirectURI := r.URL.Query().Get("redirect_uri")
        if redirectURI == "" {
            http.Error(w, "Missing redirect_uri parameter", http.StatusBadRequest)
            return
        }
        
        logger.InfoLogger.Printf("OAuth URL request for provider: %s with redirect URI: %s", provider, redirectURI)
        
        var oauthURL string
        
        switch provider {
        case "github":
            oauthURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=Ov23liJMbcmt6eXGI7yN&redirect_uri=%s&scope=read:user,user:email", 
                url.QueryEscape(redirectURI))
        case "google":
            oauthURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=64583008448-4aa9mivl1jurlp1bheabkc5m0irc6fsp.apps.googleusercontent.com&redirect_uri=%s&response_type=code&scope=email+profile", 
                url.QueryEscape(redirectURI))
        default:
            http.Error(w, "Invalid provider", http.StatusBadRequest)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "url": oauthURL,
        })
    }).Methods("GET", "OPTIONS")
    
    // OAuth sign-in endpoint - CRITICAL FIX
    router.HandleFunc("/api/auth/oauth/{provider}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        vars := mux.Vars(r)
        provider := vars["provider"]
        
        var req struct {
            Code        string `json:"code"`
            RedirectURI string `json:"redirect_uri"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            logger.ErrorLogger.Printf("Failed to parse OAuth request: %v", err)
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }
        
        var name, email string
        switch provider {
        case "github":
            name = "GitHub User"
            email = "user@github.example.com"
        case "google":
            name = "Google User"
            email = "user@google.example.com"
        default:
            name = "OAuth User"
            email = "user@example.com"
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "user": map[string]interface{}{
                "id": fmt.Sprintf("%s_user_123", provider),
                "name": name,
                "email": email,
                "createdAt": time.Now(),
                "lastLoginAt": time.Now(),
            },
            "accessToken": "mock_oauth_token_" + provider,
        })
    }).Methods("POST", "OPTIONS")
    
    router.HandleFunc("/api/auth/user", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "id": "user123",
            "name": "Test User",
            "email": "user@example.com",
            "createdAt": time.Now(),
            "lastLoginAt": time.Now(),
        })
    }).Methods("GET", "OPTIONS")
    
    // Global handler for OPTIONS requests
    router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        w.Header().Set("Access-Control-Max-Age", "3600")
        w.WriteHeader(http.StatusOK)
    })
    
    // Health check endpoint
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")
    
    return router
}

// handleQueryRequest processes the query against a CSV file
func handleQueryRequest(w http.ResponseWriter, r *http.Request, orch *orchestrator.Orchestrator) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    var req struct {
        Query            string                 `json:"query"`
        CSVPath          string                 `json:"csvPath"`
        UseKnowledgeBase bool                   `json:"useKnowledgeBase"`
        TrainingDataIDs  []string               `json:"trainingDataIds"`
        Options          map[string]interface{} `json:"options"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Options == nil {
        req.Options = make(map[string]interface{})
    }
    
    req.Options["useKnowledgeBase"] = req.UseKnowledgeBase
    if len(req.TrainingDataIDs) > 0 {
        req.Options["trainingDataIds"] = req.TrainingDataIDs
    }
    
    result, err := orch.ProcessQueryWithOptions(r.Context(), req.Query, req.CSVPath, req.Options)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error processing query: %v", err), http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "query":  req.Query,
        "sql":    result.GeneratedQuery,
        "results": result.ExecutionResult,
        "knowledgeContext": nil,
    }
    
    if result.Analysis != nil {
        if resp, ok := result.Analysis["response"].(string); ok {
            response["response"] = resp
        }
    }
    
    if result.KnowledgeContext != nil {
        knowledgeContext := []map[string]interface{}{}
        
        for _, ddl := range result.KnowledgeContext.DDLSchemas {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "ddl",
                "id":          ddl.ID,
                "description": ddl.Description,
                "content":     ddl.Content,
            })
        }
        
        for _, doc := range result.KnowledgeContext.Documentation {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "documentation",
                "id":          doc.ID,
                "description": doc.Description,
                "content":     doc.Content,
            })
        }
        
        for _, pair := range result.KnowledgeContext.QuestionSQLPairs {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "question_sql",
                "id":          pair.Description,
                "description": pair.Description,
                "question":    pair.Question,
                "sql":         pair.SQL,
            })
        }
        
        response["knowledgeContext"] = knowledgeContext
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// handleFileUpload processes file uploads for CSV data
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, "Could not parse form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "No file provided", http.StatusBadRequest)
        return
    }
    defer file.Close()

    timestamp := time.Now().UnixNano()
    filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
    filepath := filepath.Join("data", "uploads", filename)

    dst, err := os.Create(filepath)
    if err != nil {
        http.Error(w, "Failed to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    headers, err := extractCSVHeaders(filepath)
    if err != nil {
        logger.ErrorLogger.Printf("Error extracting CSV headers: %v", err)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success":   true,
        "filename":  handler.Filename,
        "filePath":  filepath,
        "timestamp": timestamp,
        "headers":   headers,
    })
}

// extractCSVHeaders reads the first row of a CSV file to get the column headers
func extractCSVHeaders(filePath string) ([]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    headers, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV headers: %w", err)
    }
    
    return headers, nil
}