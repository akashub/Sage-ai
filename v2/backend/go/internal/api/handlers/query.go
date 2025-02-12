// backend/go/internal/api/handlers/query.go
package handlers

import (
    "encoding/json"
    "net/http"
    "sage-ai-v2/internal/orchestrator"
)

type QueryRequest struct {
    Query    string `json:"query"`
    CSVPath  string `json:"csv_path"`
}

type QueryHandler struct {
    orchestrator *orchestrator.Orchestrator
}

func CreateQueryHandler(orch *orchestrator.Orchestrator) *QueryHandler {
    return &QueryHandler{orchestrator: orch}
}

func (h *QueryHandler) Handle(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req QueryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := h.orchestrator.ProcessQuery(r.Context(), req.Query, req.CSVPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}