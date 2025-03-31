// // // // backend/go/internal/api/routes.go
// // // package api

// // // import (
// // // 	"database/sql"
// // // 	"net/http"
// // // 	"sage-ai-v2/internal/api/handlers"
// // // 	"sage-ai-v2/internal/models"
// // // 	"sage-ai-v2/internal/services"
// // // 	"sage-ai-v2/pkg/logger"
// // // 	"time"
// // // )

// // // // SetupRoutes configures all API routes for the application
// // // func SetupRoutes(db *sql.DB, orch *orchestrator.Orchestrator) http.Handler {
// // // 	// Create a new ServeMux
// // // 	mux := http.NewServeMux()

// // // 	// Initialize authentication service
// // // 	authService := setupAuthService(db)

// // // 	// Add authentication routes
// // // 	AddAuthRoutes(mux, authService)

// // // 	// Register API routes
// // // 	mux.HandleFunc("/api/upload", handlers.UploadFileHandler)
// // // 	mux.HandleFunc("/api/query", handlers.QueryHandler)

// // // 	// Add health check endpoint
// // // 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// // // 		w.Header().Set("Content-Type", "text/plain")
// // // 		w.WriteHeader(http.StatusOK)
// // // 		w.Write([]byte("OK"))
// // // 	})

// // // 	// Apply middleware directly without using a separate middleware package
// // // 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // // 		// Set CORS headers for all responses
// // // 		w.Header().Set("Access-Control-Allow-Origin", "*")
// // // 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// // // 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// // // 		// Handle preflight requests
// // // 		if r.Method == "OPTIONS" {
// // // 			w.WriteHeader(http.StatusOK)
// // // 			return
// // // 		}

// // // 		// Log the request
// // // 		logger.InfoLogger.Printf("Request: %s %s", r.Method, r.URL.Path)

// // // 		// Continue to the next handler
// // // 		mux.ServeHTTP(w, r)
// // // 	})

// // // 	logger.InfoLogger.Printf("API routes configured")
// // // 	return handler
// // // }

// // // // setupAuthService creates and configures the authentication service
// // // func setupAuthService(db *sql.DB) *services.AuthService {
// // // 	// JWT configuration
// // // 	jwtSecret := "your-secret-key" // In production, load from environment variable
// // // 	jwtExpiry := 7 * 24 * time.Hour // 7 days

// // // 	// OAuth configurations
// // // 	oauthConfs := map[string]models.OAuthConfig{
// // // 		"google": {
// // // 			ClientID:     "64583008448-4aa9mivl1jurlp1bheabkc5m0irc6fsp.apps.googleusercontent.com",
// // // 			ClientSecret: "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",
// // // 			RedirectURI:  "http://localhost:3000/oauth-callback",
// // // 			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
// // // 			TokenURL:     "https://oauth2.googleapis.com/token",
// // // 			UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
// // // 			Scopes:       []string{"email", "profile"},
// // // 		},
// // // 		"github": {
// // // 			ClientID:     "Ov23liJMbcmt6eXGI7yN",
// // // 			ClientSecret: "04617e96169696a53048a2bdc886c5d9ae38268d",
// // // 			RedirectURI:  "http://localhost:3000/oauth-callback",
// // // 			AuthURL:      "https://github.com/login/oauth/authorize",
// // // 			TokenURL:     "https://github.com/login/oauth/access_token",
// // // 			UserInfoURL:  "https://api.github.com/user",
// // // 			Scopes:       []string{"read:user", "user:email"},
// // // 		},
// // // 		// "github": {
// // // 		// 	ClientID:     "Ov23liFLEJCnd0fpR0P0",
// // // 		// 	ClientSecret: "6707da4695ea33aea2b485c3ba4edb420ebe77f7",
// // // 		// 	RedirectURI:  "http://localhost:3000/oauth-callback",
// // // 		// 	AuthURL:      "https://github.com/login/oauth/authorize",
// // // 		// 	TokenURL:     "https://github.com/login/oauth/access_token",
// // // 		// 	UserInfoURL:  "https://api.github.com/user",
// // // 		// 	Scopes:       []string{"read:user", "user:email"},
// // // 		// },
// // // 	}

// // // 	return services.NewAuthService(db, jwtSecret, jwtExpiry, oauthConfs)
// // // }

// // package api

// // import (
// // 	"database/sql"
// // 	"net/http"
// // 	"sage-ai-v2/internal/api/handlers"
// // 	"sage-ai-v2/internal/models"
// // 	"sage-ai-v2/internal/orchestrator"
// // 	"sage-ai-v2/internal/services"
// // 	"sage-ai-v2/pkg/logger"
// // 	"time"
// // )

// // // SetupRoutes configures all API routes for the application
// // func SetupRoutes(db *sql.DB, orch *orchestrator.Orchestrator) http.Handler {
// // 	// Create a new ServeMux
// // 	mux := http.NewServeMux()

// // 	// Initialize authentication service
// // 	authService := setupAuthService(db)

// // 	// Add authentication routes
// // 	AddAuthRoutes(mux, authService)

// // 	// Register API routes
// // 	mux.HandleFunc("/api/upload", handlers.UploadFileHandler)

// // 	// Use the orchestrator for query handling
// // 	if orch != nil {
// // 		queryHandler := handlers.QueryHandler(orch)
// // 		mux.HandleFunc("/api/query", queryHandler)

// // 		// Add training data handlers
// // 		mux.HandleFunc("/api/training/upload", handlers.UploadTrainingFileHandler(orch))
// // 		mux.HandleFunc("/api/training/add", handlers.AddTrainingDataHandler(orch))
// // 		mux.HandleFunc("/api/training/list", handlers.ListTrainingDataHandler(orch))
// // 	} else {
// // 		// Fallback for when orchestrator is not provided
// // 		mux.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
// // 			http.Error(w, "Orchestrator not initialized", http.StatusServiceUnavailable)
// // 		})
// // 	}

// // 	// Add health check endpoint
// // 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// // 		w.Header().Set("Content-Type", "text/plain")
// // 		w.WriteHeader(http.StatusOK)
// // 		w.Write([]byte("OK"))
// // 	})

// // 	// Apply middleware directly without using a separate middleware package
// // 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		// Set CORS headers for all responses
// // 		w.Header().Set("Access-Control-Allow-Origin", "*")
// // 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// // 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// // 		// Handle preflight requests
// // 		if r.Method == "OPTIONS" {
// // 			w.WriteHeader(http.StatusOK)
// // 			return
// // 		}

// // 		// Log the request
// // 		logger.InfoLogger.Printf("Request: %s %s", r.Method, r.URL.Path)

// // 		// Continue to the next handler
// // 		mux.ServeHTTP(w, r)
// // 	})

// // 	logger.InfoLogger.Printf("API routes configured")
// // 	return handler
// // }

// // // setupAuthService creates and configures the authentication service
// // func setupAuthService(db *sql.DB) *services.AuthService {
// // 	// JWT configuration
// // 	jwtSecret := "your-secret-key" // In production, load from environment variable
// // 	jwtExpiry := 7 * 24 * time.Hour // 7 days

// // 	// OAuth configurations
// // 	oauthConfs := map[string]models.OAuthConfig{
// // 		"google": {
// // 			ClientID:     "64583008448-4aa9mivl1jurlp1bheabkc5m0irc6fsp.apps.googleusercontent.com",
// // 			ClientSecret: "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",
// // 			RedirectURI:  "http://localhost:3000/oauth-callback",
// // 			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
// // 			TokenURL:     "https://oauth2.googleapis.com/token",
// // 			UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
// // 			Scopes:       []string{"email", "profile"},
// // 		},
// // 		"github": {
// // 			ClientID:     "Ov23liJMbcmt6eXGI7yN",
// // 			ClientSecret: "04617e96169696a53048a2bdc886c5d9ae38268d",
// // 			RedirectURI:  "http://localhost:3000/oauth-callback",
// // 			AuthURL:      "https://github.com/login/oauth/authorize",
// // 			TokenURL:     "https://github.com/login/oauth/access_token",
// // 			UserInfoURL:  "https://api.github.com/user",
// // 			Scopes:       []string{"read:user", "user:email"},
// // 		},
// // 	}

// // 	return services.NewAuthService(db, jwtSecret, jwtExpiry, oauthConfs)
// // }

// // backend/go/internal/api/routes.go
// // backend/go/internal/api/routes.go
// package api

// import (
// 	"database/sql"
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"sage-ai-v2/internal/api/middleware"
// 	"sage-ai-v2/internal/knowledge"
// 	"sage-ai-v2/internal/orchestrator"
// 	"sage-ai-v2/pkg/logger"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// // SetupRoutes configures the API routes
// func SetupRoutes(db *sql.DB, orch *orchestrator.Orchestrator) *mux.Router {
//     // Create router
//     router := mux.NewRouter()
    
//     // Apply middleware - make sure this is before any routes are defined
//     router.Use(middleware.LoggingMiddleware)
//     router.Use(middleware.CORSMiddleware)
    
//     // Extract knowledge manager from orchestrator
//     km := orch.KnowledgeManager

//     // Create uploads directory if it doesn't exist
//     os.MkdirAll("data/uploads", 0755)

//     // Define all your routes below...
    
//     // Query API
//     router.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
//         handleQueryRequest(w, r, orch)
//     }).Methods("GET", "POST", "OPTIONS") // Make sure OPTIONS is included

//     // File Upload API
//     router.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
//         handleFileUpload(w, r)
//     }).Methods("POST", "OPTIONS") // Make sure OPTIONS is included

//     // Training Data API
//     setupTrainingDataRoutes(router, km)
    
//     // Chat History API
//     SetupChatRoutes(router)

//     // Auth routes (if implemented)
//     setupAuthRoutes(router, db)
    
//     return router
// }

// // handleQueryRequest processes the query against a CSV file
// func handleQueryRequest(w http.ResponseWriter, r *http.Request, orch *orchestrator.Orchestrator) {
// 	// Parse request
// 	var req struct {
// 		Query           string                 `json:"query"`
// 		CSVPath         string                 `json:"csvPath"`
// 		UseKnowledgeBase bool                  `json:"useKnowledgeBase"`
// 		Options         map[string]interface{} `json:"options"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	logger.InfoLogger.Printf("Query received: %s", req.Query)
// 	logger.InfoLogger.Printf("CSV Path: %s", req.CSVPath)

// 	// Set default options if not provided
// 	if req.Options == nil {
// 		req.Options = make(map[string]interface{})
// 	}
	
// 	// Add the useKnowledgeBase flag to options
// 	req.Options["useKnowledgeBase"] = req.UseKnowledgeBase
	
// 	// Process the query
// 	ctx := r.Context()
// 	result, err := orch.ProcessQueryWithOptions(ctx, req.Query, req.CSVPath, req.Options)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error processing query: %v", err)
// 		http.Error(w, fmt.Sprintf("Error processing query: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Prepare response
// 	response := map[string]interface{}{
// 		"query":           req.Query,
// 		"sql":             result.GeneratedQuery,
// 		"results":         result.ExecutionResult,
// 		"knowledgeContext": nil,
// 	}
	
// 	// Add natural language response if available
// 	if result.Analysis != nil {
// 		if resp, ok := result.Analysis["response"].(string); ok {
// 			response["response"] = resp
// 		}
// 	}
	
// 	// Add knowledge context if available
// 	if result.KnowledgeContext != nil {
// 		// Prepare knowledge context for response
// 		knowledgeContext := make([]map[string]interface{}, 0)
		
// 		// Add DDL schemas
// 		for _, ddl := range result.KnowledgeContext.DDLSchemas {
// 			knowledgeContext = append(knowledgeContext, map[string]interface{}{
// 				"type":        "ddl",
// 				"id":          ddl.ID,
// 				"description": ddl.Description,
// 				"content":     "",  // Don't include full content in response
// 			})
// 		}
		
// 		// Add documentation
// 		for _, doc := range result.KnowledgeContext.Documentation {
// 			knowledgeContext = append(knowledgeContext, map[string]interface{}{
// 				"type":        "documentation",
// 				"id":          doc.ID,
// 				"description": doc.Description,
// 				"content":     "",  // Don't include full content in response
// 			})
// 		}
		
// 		// Add question-SQL pairs
// 		for _, pair := range result.KnowledgeContext.QuestionSQLPairs {
// 			knowledgeContext = append(knowledgeContext, map[string]interface{}{
// 				"type":        "question_sql",
// 				"id":          pair.Description, // Using description as ID
// 				"description": pair.Description,
// 				"question":    pair.Question,
// 				"sql":         pair.SQL,
// 			})
// 		}
		
// 		response["knowledgeContext"] = knowledgeContext
// 	}

// 	// Send response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // handleFileUpload processes file uploads for CSV data
// func handleFileUpload(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
//     // Handle OPTIONS request
//     if r.Method == "OPTIONS" {
//         w.WriteHeader(http.StatusOK)
//         return
//     }
// 	// Parse multipart form file
// 	err := r.ParseMultipartForm(32 << 20) // 32MB max size
// 	if err != nil {
// 		http.Error(w, "Could not parse form", http.StatusBadRequest)
// 		return
// 	}

// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "No file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a unique filename
// 	timestamp := time.Now().UnixNano()
// 	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
// 	filepath := filepath.Join("data", "uploads", filename)

// 	// Create the file
// 	dst, err := os.Create(filepath)
// 	if err != nil {
// 		http.Error(w, "Failed to create file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer dst.Close()

// 	// Copy file contents
// 	if _, err := io.Copy(dst, file); err != nil {
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Extract CSV headers
// 	headers, err := extractCSVHeaders(filepath)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error extracting CSV headers: %v", err)
// 		// Don't fail the request, just log the error
// 		// The frontend will handle missing headers
// 	}

// 	// Return response with file info
// 	response := map[string]interface{}{
// 		"success":   true,
// 		"filename":  handler.Filename,
// 		"filePath":  filepath,
// 		"timestamp": timestamp,
// 		"headers":   headers,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // extractCSVHeaders reads the first row of a CSV file to get the column headers
// func extractCSVHeaders(filePath string) ([]string, error) {
//     // Open the file
//     file, err := os.Open(filePath)
//     if err != nil {
//         return nil, fmt.Errorf("failed to open file: %w", err)
//     }
//     defer file.Close()
    
//     // Create a new CSV reader
//     reader := csv.NewReader(file)
    
//     // Read the first row (headers)
//     headers, err := reader.Read()
//     if err != nil {
//         return nil, fmt.Errorf("failed to read CSV headers: %w", err)
//     }
    
//     return headers, nil
// }

// // setupTrainingDataRoutes configures training data API routes
// func setupTrainingDataRoutes(router *mux.Router, km *knowledge.KnowledgeManager) {
// 	// List training data
// 	router.HandleFunc("/api/training/list", func(w http.ResponseWriter, r *http.Request) {
// 		dataType := r.URL.Query().Get("type")
// 		items, err := km.ListTrainingData(r.Context(), dataType)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
		
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(items)
// 	}).Methods("GET")
	
// 	// Upload training file
// 	router.HandleFunc("/api/training/upload", func(w http.ResponseWriter, r *http.Request) {
// 		// Explicitly set CORS headers for every response
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
// 		// Handle preflight OPTIONS request
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
		
// 		// Continue with normal processing for POST
// 		err := r.ParseMultipartForm(32 << 20) // 32MB max
// 		if err != nil {
// 			http.Error(w, "Could not parse form", http.StatusBadRequest)
// 			return
// 		}
		
// 		file, handler, err := r.FormFile("file")
// 		if err != nil {
// 			http.Error(w, "No file provided", http.StatusBadRequest)
// 			return
// 		}
// 		defer file.Close()
		
// 		// Read file content
// 		content, err := io.ReadAll(file)
// 		if err != nil {
// 			http.Error(w, "Failed to read file", http.StatusInternalServerError)
// 			return
// 		}
		
// 		// Store file
// 		localPath, err := km.StoreFile(handler.Filename, content)
// 		if err != nil {
// 			http.Error(w, "Failed to store file", http.StatusInternalServerError)
// 			return
// 		}
		
// 		// Get training data type
// 		dataType := r.FormValue("type")
// 		if dataType == "" {
// 			dataType = "ddl" // Default type
// 		}
		
// 		// Get description
// 		description := r.FormValue("description")
// 		if description == "" {
// 			description = handler.Filename
// 		}
		
// 		// Process based on type
// 		ctx := r.Context()
// 		switch dataType {
// 		case "ddl":
// 			// Add DDL schema
// 			err = km.AddDDLSchema(ctx, handler.Filename, string(content), description)
// 		case "documentation":
// 			// Add documentation
// 			err = km.AddDocumentation(ctx, description, string(content), []string{})
// 		case "question_sql_json":
// 			// Load from JSON
// 			pairs, err := km.LoadQuestionSQLPairsFromJSON(ctx, localPath)
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Failed to load JSON: %v", err), http.StatusInternalServerError)
// 				return
// 			}
			
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"success": true,
// 				"count":   pairs,
// 				"path":    localPath,
// 			})
// 			return
// 		case "auto":
// 			// Auto-detect type based on file extension
// 			if strings.HasSuffix(handler.Filename, ".sql") {
// 				err = km.AddDDLSchema(ctx, handler.Filename, string(content), description)
// 			} else if strings.HasSuffix(handler.Filename, ".json") {
// 				pairs, err := km.LoadQuestionSQLPairsFromJSON(ctx, localPath)
// 				if err != nil {
// 					http.Error(w, fmt.Sprintf("Failed to load JSON: %v", err), http.StatusInternalServerError)
// 					return
// 				}
				
// 				w.Header().Set("Content-Type", "application/json")
// 				json.NewEncoder(w).Encode(map[string]interface{}{
// 					"success": true,
// 					"count":   pairs,
// 					"path":    localPath,
// 				})
// 				return
// 			} else {
// 				err = km.AddDocumentation(ctx, description, string(content), []string{})
// 			}
// 		default:
// 			http.Error(w, "Invalid data type", http.StatusBadRequest)
// 			return
// 		}
		
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to add training data: %v", err), http.StatusInternalServerError)
// 			return
// 		}
		
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"success": true,
// 			"type":    dataType,
// 			"path":    localPath,
// 		})
// 	}).Methods("POST", "OPTIONS")	
	
// 	// Add training data manually
// 	router.HandleFunc("/api/training/add", func(w http.ResponseWriter, r *http.Request) {
// 		var data struct {
// 			Type        string                 `json:"type"`
// 			Content     string                 `json:"content"`
// 			Description string                 `json:"description"`
// 			Metadata    map[string]interface{} `json:"metadata"`
// 		}
		
// 		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}
		
// 		// Process based on type
// 		ctx := r.Context()
// 		var err error
// 		switch data.Type {
// 		case "ddl":
// 			err = km.AddDDLSchema(ctx, data.Description, data.Content, data.Description)
// 		case "documentation":
// 			tags := []string{}
// 			if tagsData, ok := data.Metadata["tags"].([]interface{}); ok {
// 				for _, tag := range tagsData {
// 					if tagStr, ok := tag.(string); ok {
// 						tags = append(tags, tagStr)
// 					}
// 				}
// 			}
// 			err = km.AddDocumentation(ctx, data.Description, data.Content, tags)
// 		case "question_sql":
// 			pair := knowledge.QuestionSQLPair{
// 				Question:    data.Content,
// 				SQL:         "", // SQL should be in metadata
// 				Description: data.Description,
// 				DateAdded:   time.Now().Format(time.RFC3339),
// 				Verified:    true,
// 			}
			
// 			// Extract SQL from metadata
// 			if sqlData, ok := data.Metadata["sql"].(string); ok {
// 				pair.SQL = sqlData
// 			}
			
// 			// Extract tags from metadata
// 			if tagsData, ok := data.Metadata["tags"].([]interface{}); ok {
// 				for _, tag := range tagsData {
// 					if tagStr, ok := tag.(string); ok {
// 						pair.Tags = append(pair.Tags, tagStr)
// 					}
// 				}
// 			}
			
// 			err = km.AddQuestionSQLPair(ctx, pair)
// 		default:
// 			http.Error(w, "Invalid data type", http.StatusBadRequest)
// 			return
// 		}
		
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to add training data: %v", err), http.StatusInternalServerError)
// 			return
// 		}
		
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"success": true,
// 			"type":    data.Type,
// 		})
// 	}).Methods("POST")
	
// 	// View training data
// 	// router.HandleFunc("/api/training/view/{id}", func(w http.ResponseWriter, r *http.Request) {
// 	// 	// Extract ID from URL
// 	// 	vars := mux.Vars(r)
// 	// 	id := vars["id"]
		
// 	// 	// Get all training data items
// 	// 	items, err := km.ListTrainingData(r.Context(), "")
// 	// 	if err != nil {
// 	// 		http.Error(w, fmt.Sprintf("Failed to list training data: %v", err), http.StatusInternalServerError)
// 	// 		return
// 	// 	}
		
// 	// 	var item *knowledge.TrainingItem
// 	// 	for _, itemMap := range items {
// 	// 		// Check if this item matches the requested ID
// 	// 		if itemID, ok := itemMap["id"].(string); ok && itemID == id {
// 	// 			itemType, _ := itemMap["type"].(string)
// 	// 			description, _ := itemMap["description"].(string)
// 	// 			dateAdded, _ := itemMap["date_added"].(string)
				
// 	// 			// For a real system, we would retrieve the full content
// 	// 			// Here we'll just return a placeholder
// 	// 			item = &knowledge.TrainingItem{
// 	// 				ID:          id,
// 	// 				Type:        itemType,
// 	// 				Description: description,
// 	// 				DateAdded:   dateAdded,
// 	// 				Content:     "This is the full content of the training item " + id,
// 	// 			}
// 	// 			break
// 	// 		}
// 	// 	}
		
// 	// 	if item == nil {
// 	// 		http.Error(w, "Training item not found", http.StatusNotFound)
// 	// 		return
// 	// 	}
		
// 	// 	w.Header().Set("Content-Type", "application/json")
// 	// 	json.NewEncoder(w).Encode(item)
// 	// }).Methods("GET")
// 	// View training data
// 	router.HandleFunc("/api/training/view/{id}", func(w http.ResponseWriter, r *http.Request) {
// 		// Set CORS headers
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
// 		// Handle OPTIONS request
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
		
// 		// Extract ID from URL
// 		vars := mux.Vars(r)
// 		id := vars["id"]
		
// 		// Validate ID
// 		if id == "" {
// 			logger.ErrorLogger.Printf("Invalid or empty ID parameter")
// 			http.Error(w, "Invalid or empty ID parameter", http.StatusBadRequest)
// 			return
// 		}
		
// 		// Get the full training item
// 		item, err := km.GetTrainingItem(r.Context(), id)
// 		if err != nil {
// 			logger.ErrorLogger.Printf("Error retrieving training item: %v", err)
// 			http.Error(w, fmt.Sprintf("Error retrieving training item: %v", err), http.StatusInternalServerError)
// 			return
// 		}
		
// 		// If no content is available, add a placeholder
// 		if item.Content == "" {
// 			item.Content = "No content available for this item"
// 		}
		
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(item)
// 	}).Methods("GET", "OPTIONS")
	
// 	// Delete training data
// 	// router.HandleFunc("/api/training/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
// 	// 	// Set CORS headers immediately
// 	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
// 	// 	// Handle OPTIONS request
// 	// 	if r.Method == "OPTIONS" {
// 	// 		w.WriteHeader(http.StatusOK)
// 	// 		return
// 	// 	}
		
// 	// 	// Extract ID from URL
// 	// 	vars := mux.Vars(r)
// 	// 	id := vars["id"]
		
// 	// 	logger.InfoLogger.Printf("Deleting training item: %s", id)
		
// 	// 	// Important: Return success immediately without waiting for any background operations
// 	// 	w.WriteHeader(http.StatusNoContent)
		
// 	// 	// Run the actual deletion logic in a separate goroutine
// 	// 	go func() {
// 	// 		// This function will run in the background after response is sent
// 	// 		// Get all items - you can log errors but can't return them to client
// 	// 		items, err := km.ListTrainingData(context.Background(), "")
// 	// 		if err != nil {
// 	// 			logger.ErrorLogger.Printf("Error listing training data: %v", err)
// 	// 			return
// 	// 		}
			
// 	// 		// Check if item exists and log
// 	// 		var found bool = false
// 	// 		for _, item := range items {
// 	// 			if itemID, ok := item["id"].(string); ok && itemID == id {
// 	// 				found = true
// 	// 				logger.InfoLogger.Printf("Found item to delete: %s", itemID)
// 	// 				break
// 	// 			}
// 	// 		}
			
// 	// 		if !found {
// 	// 			logger.ErrorLogger.Printf("Training item not found: %s", id)
// 	// 			return
// 	// 		}
			
// 	// 		// Add actual deletion logic here when implemented
// 	// 		logger.InfoLogger.Printf("Successfully deleted training item: %s", id)
// 	// 	}()
// 	// }).Methods("DELETE", "OPTIONS")
// 	router.HandleFunc("/api/training/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
// 		// Set CORS headers
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
// 		// Handle OPTIONS request
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
		
// 		// Extract ID from URL
// 		vars := mux.Vars(r)
// 		id := vars["id"]
		
// 		logger.InfoLogger.Printf("Attempting to delete training item: %s", id)
		
// 		// Call the new deletion method
// 		err := km.DeleteTrainingItem(r.Context(), id)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "not found") {
// 				logger.ErrorLogger.Printf("Item not found: %s", id)
// 				http.Error(w, "Item not found", http.StatusNotFound)
// 				return
// 			}
			
// 			logger.ErrorLogger.Printf("Error deleting item: %v", err)
// 			http.Error(w, fmt.Sprintf("Error deleting item: %v", err), http.StatusInternalServerError)
// 			return
// 		}
		
// 		logger.InfoLogger.Printf("Successfully deleted training item: %s", id)
// 		w.WriteHeader(http.StatusNoContent)
// 	}).Methods("DELETE", "OPTIONS")
// }
// // setupAuthRoutes configures authentication API routes
// func setupAuthRoutes(router *mux.Router, db *sql.DB) {
// 	// These are placeholders for auth routes
// 	// Implement as needed
	
// 	// Login
// 	router.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
// 		// Mock implementation
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"token": "mock_token",
// 		})
// 	}).Methods("POST")
	
// 	// Register
// 	router.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
// 		// Mock implementation
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"token": "mock_token",
// 		})
// 	}).Methods("POST")

	
// }

// backend/go/internal/api/routes.go
package api

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sage-ai-v2/internal/knowledge"
	"sage-ai-v2/internal/orchestrator"
	"sage-ai-v2/internal/api/middleware"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the API routes
func SetupRoutes(db *sql.DB, orch *orchestrator.Orchestrator) *mux.Router {
    // Create router
    router := mux.NewRouter()
    
    // Global handler for OPTIONS requests - this must come before other routes
    router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.InfoLogger.Printf("Global OPTIONS handler called for: %s", r.URL.Path)
        
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        w.Header().Set("Access-Control-Max-Age", "3600") // Cache preflight for 1 hour
        
        w.WriteHeader(http.StatusOK)
    })
    
    // Apply middleware
    router.Use(middleware.LoggingMiddleware)
    router.Use(middleware.CORSMiddleware)
    
    // Extract knowledge manager from orchestrator
    km := orch.KnowledgeManager

    // Create uploads directory if it doesn't exist
    os.MkdirAll("data/uploads", 0755)

    // Query API
    router.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
        handleQueryRequest(w, r, orch)
    }).Methods("POST", "OPTIONS")

    // File Upload API
    router.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
        handleFileUpload(w, r)
    }).Methods("POST", "OPTIONS")

    // Training Data API
    setupTrainingDataRoutes(router, km)
    
    // Chat History API - using the simplified version
    SetupChatRoutes(router)

    // Auth routes (if implemented)
    setupAuthRoutes(router, db)
    
    return router
}

// handleQueryRequest processes the query against a CSV file
func handleQueryRequest(w http.ResponseWriter, r *http.Request, orch *orchestrator.Orchestrator) {
    // Set CORS headers explicitly
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Parse request
    var req struct {
        Query           string                 `json:"query"`
        CSVPath         string                 `json:"csvPath"`
        UseKnowledgeBase bool                  `json:"useKnowledgeBase"`
        Options         map[string]interface{} `json:"options"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    logger.InfoLogger.Printf("Query received: %s", req.Query)
    logger.InfoLogger.Printf("CSV Path: %s", req.CSVPath)

    // Set default options if not provided
    if req.Options == nil {
        req.Options = make(map[string]interface{})
    }
    
    // Add the useKnowledgeBase flag to options
    req.Options["useKnowledgeBase"] = req.UseKnowledgeBase
    
    // Process the query
    ctx := r.Context()
    result, err := orch.ProcessQueryWithOptions(ctx, req.Query, req.CSVPath, req.Options)
    if err != nil {
        logger.ErrorLogger.Printf("Error processing query: %v", err)
        http.Error(w, fmt.Sprintf("Error processing query: %v", err), http.StatusInternalServerError)
        return
    }

    // Prepare response
    response := map[string]interface{}{
        "query":           req.Query,
        "sql":             result.GeneratedQuery,
        "results":         result.ExecutionResult,
        "knowledgeContext": nil,
    }
    
    // Add natural language response if available
    if result.Analysis != nil {
        if resp, ok := result.Analysis["response"].(string); ok {
            response["response"] = resp
        }
    }
    
    // Add knowledge context if available
    if result.KnowledgeContext != nil {
        // Prepare knowledge context for response
        knowledgeContext := make([]map[string]interface{}, 0)
        
        // Add DDL schemas
        for _, ddl := range result.KnowledgeContext.DDLSchemas {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "ddl",
                "id":          ddl.ID,
                "description": ddl.Description,
                "content":     "",  // Don't include full content in response
            })
        }
        
        // Add documentation
        for _, doc := range result.KnowledgeContext.Documentation {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "documentation",
                "id":          doc.ID,
                "description": doc.Description,
                "content":     "",  // Don't include full content in response
            })
        }
        
        // Add question-SQL pairs
        for _, pair := range result.KnowledgeContext.QuestionSQLPairs {
            knowledgeContext = append(knowledgeContext, map[string]interface{}{
                "type":        "question_sql",
                "id":          pair.Description, // Using description as ID
                "description": pair.Description,
                "question":    pair.Question,
                "sql":         pair.SQL,
            })
        }
        
        response["knowledgeContext"] = knowledgeContext
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// handleFileUpload processes file uploads for CSV data
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers explicitly
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Parse multipart form file
    err := r.ParseMultipartForm(32 << 20) // 32MB max size
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

    // Create a unique filename
    timestamp := time.Now().UnixNano()
    filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
    filepath := filepath.Join("data", "uploads", filename)

    // Create the file
    dst, err := os.Create(filepath)
    if err != nil {
        http.Error(w, "Failed to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy file contents
    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    // Extract CSV headers
    headers, err := extractCSVHeaders(filepath)
    if err != nil {
        logger.ErrorLogger.Printf("Error extracting CSV headers: %v", err)
        // Don't fail the request, just log the error
        // The frontend will handle missing headers
    }

    // Return response with file info
    response := map[string]interface{}{
        "success":   true,
        "filename":  handler.Filename,
        "filePath":  filepath,
        "timestamp": timestamp,
        "headers":   headers,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// extractCSVHeaders reads the first row of a CSV file to get the column headers
func extractCSVHeaders(filePath string) ([]string, error) {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    // Create a new CSV reader
    reader := csv.NewReader(file)
    
    // Read the first row (headers)
    headers, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV headers: %w", err)
    }
    
    return headers, nil
}

// setupTrainingDataRoutes configures training data API routes
func setupTrainingDataRoutes(router *mux.Router, km *knowledge.KnowledgeManager) {
    // List training data
    router.HandleFunc("/api/training/list", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
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
    
    // Upload training file
    router.HandleFunc("/api/training/upload", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        err := r.ParseMultipartForm(32 << 20) // 32MB max
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
        
        // Read file content
        content, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "Failed to read file", http.StatusInternalServerError)
            return
        }
        
        // Store file
        localPath, err := km.StoreFile(handler.Filename, content)
        if err != nil {
            http.Error(w, "Failed to store file", http.StatusInternalServerError)
            return
        }
        
        // Get training data type
        dataType := r.FormValue("type")
        if dataType == "" {
            dataType = "ddl" // Default type
        }
        
        // Get description
        description := r.FormValue("description")
        if description == "" {
            description = handler.Filename
        }
        
        // Generate a unique ID for the new training item
        id := fmt.Sprintf("%s_%d", dataType, time.Now().UnixNano())
        
        // Process based on type
        ctx := r.Context()
        switch dataType {
        case "ddl":
            // Add DDL schema
            err = km.AddDDLSchema(ctx, handler.Filename, string(content), description)
        case "documentation":
            // Add documentation
            err = km.AddDocumentation(ctx, description, string(content), []string{})
        case "question_sql_json":
            // Load from JSON
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
            // Auto-detect type based on file extension
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
    
    // Add training data manually
    router.HandleFunc("/api/training/add", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
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
        
        // Generate a unique ID for the new training item
        id := fmt.Sprintf("%s_%d", data.Type, time.Now().UnixNano())
        
        // Process based on type
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
                SQL:         "", // SQL should be in metadata
                Description: data.Description,
                DateAdded:   time.Now().Format(time.RFC3339),
                Verified:    true,
            }
            
            // Extract SQL from metadata
            if sqlData, ok := data.Metadata["sql"].(string); ok {
                pair.SQL = sqlData
            }
            
            // Extract tags from metadata
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
    
    // View training data
    router.HandleFunc("/api/training/view/{id}", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Extract ID from URL
        vars := mux.Vars(r)
        id := vars["id"]
        
        // Get all items
        items, err := km.ListTrainingData(r.Context(), "")
        if err != nil {
            logger.ErrorLogger.Printf("Error listing training data: %v", err)
            http.Error(w, "Error listing training data", http.StatusInternalServerError)
            return
        }
        
        // Find the item with matching ID
        var found bool = false
        var item map[string]interface{}
        
        for _, itemData := range items {
            itemID, ok := itemData["id"].(string)
            if ok && itemID == id {
                found = true
                item = itemData
                
                // Add a placeholder content field
                item["content"] = fmt.Sprintf("This is the full content of training item %s", id)
                break
            }
        }
        
        if !found {
            http.Error(w, "Training item not found", http.StatusNotFound)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    }).Methods("GET", "OPTIONS")
    
    // Delete training data
    router.HandleFunc("/api/training/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Extract ID from URL
        vars := mux.Vars(r)
        id := vars["id"]
        
        logger.InfoLogger.Printf("Deleting training item: %s", id)
        
        // For now, just log the deletion and return success
        // In a real implementation, this would delete from the vector database
        logger.InfoLogger.Printf("Successfully deleted training item: %s (simulated)", id)
        
        // Return success immediately
        w.WriteHeader(http.StatusNoContent)
    }).Methods("DELETE", "OPTIONS")
}

// setupAuthRoutes configures authentication API routes
func setupAuthRoutes(router *mux.Router, db *sql.DB) {
    // These are placeholders for auth routes
    // Implement as needed
    
    // Login
    router.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Mock implementation
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "token": "mock_token",
        })
    }).Methods("POST", "OPTIONS")
    
    // Register
    router.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers explicitly
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Mock implementation
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "token": "mock_token",
        })
    }).Methods("POST", "OPTIONS")
}