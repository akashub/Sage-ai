// // backend/go/cmd/main.go
// package main

// import (
//     "context"
//     "log"
//     "net/http"
//     "os"
//     "os/signal"
//     "syscall"
//     "time"

//     "github.com/yourusername/sage-ai-v2/internal/api"
//     "github.com/yourusername/sage-ai-v2/internal/config"
//     "github.com/yourusername/sage-ai-v2/internal/llm"
//     "github.com/yourusername/sage-ai-v2/internal/orchestrator"
// )

// func main() {
//     // Load configuration
//     cfg := config.Load()

//     // Initialize LLM bridge
//     llmBridge := llm.PythonBridge(cfg.PythonServiceURL)

//     // Initialize orchestrator
//     orch := orchestrator.Orchestrator(llmBridge)

//     // Initialize API routes
//     router := api.Router(orch)

//     // Create server
//     srv := &http.Server{
//         Addr:    ":" + cfg.Port,
//         Handler: router,
//     }

//     // Handle graceful shutdown
//     done := make(chan bool)
//     quit := make(chan os.Signal, 1)
//     signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

//     go func() {
//         <-quit
//         log.Println("Server is shutting down...")

//         ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//         defer cancel()

//         if err := srv.Shutdown(ctx); err != nil {
//             log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
//         }
//         close(done)
//     }()

//     // Start server
//     log.Printf("Server is starting on port %s\n", cfg.Port)
//     if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//         log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
//     }

//     <-done
//     log.Println("Server stopped")
// }

// backend/go/cmd/main.go
// backend/go/cmd/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "sage-ai-v2/internal/api"
    "sage-ai-v2/internal/config"
    "sage-ai-v2/internal/llm"
    "sage-ai-v2/internal/orchestrator"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize components
    bridge := llm.CreateBridge(cfg.PythonServiceURL)
    orch := orchestrator.CreateOrchestrator(bridge)
    router := api.CreateRouter(orch)

    // Create server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }

    // Handle graceful shutdown
    done := make(chan bool)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-quit
        log.Println("Server is shutting down...")

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := srv.Shutdown(ctx); err != nil {
            log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
        }
        close(done)
    }()

    // Start server
    log.Printf("Server is starting on port %s\n", cfg.Port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
    }

    <-done
    log.Println("Server stopped")
}