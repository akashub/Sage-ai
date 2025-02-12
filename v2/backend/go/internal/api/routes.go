// backend/go/internal/api/routes.go
package api

import (
    "net/http"
    "sage-ai-v2/internal/api/handlers"
    "sage-ai-v2/internal/orchestrator"
)

type Router struct {
    *http.ServeMux
    queryHandler  *handlers.QueryHandler
    uploadHandler *handlers.UploadHandler
}

func CreateRouter(orch *orchestrator.Orchestrator) *Router {
    r := &Router{
        ServeMux:      http.NewServeMux(),
        queryHandler:  handlers.CreateQueryHandler(orch),
        uploadHandler: handlers.CreateUploadHandler("./uploads"),
    }

    r.setupRoutes()
    return r
}

func (r *Router) setupRoutes() {
    // API endpoints
    r.Handle("/api/query", http.HandlerFunc(r.queryHandler.Handle))
    r.Handle("/api/upload", http.HandlerFunc(r.uploadHandler.Handle))
}