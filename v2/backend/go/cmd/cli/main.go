// // backend/go/cmd/cli/main.go
// package main

// import (
// 	"bufio"
// 	"context"
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"sage-ai-v2/internal/llm"
// 	"sage-ai-v2/internal/orchestrator"
// 	"sage-ai-v2/internal/types"
// )

// // func main() {
// //     // Initialize components
// //     bridge := llm.CreateBridge("http://localhost:8000")
// //     orch := orchestrator.CreateOrchestrator(bridge)

// //     // Get CSV file path
// //     fmt.Print("Enter the path to your CSV file: ")
// //     scanner := bufio.NewScanner(os.Stdin)
// //     scanner.Scan()
// //     csvPath := scanner.Text()

// //     // Validate CSV file
// //     headers, err := getCSVHeaders(csvPath)
// //     if err != nil {
// //         fmt.Printf("Error reading CSV: %v\n", err)
// //         return
// //     }

// //     fmt.Printf("\nFound columns: %v\n", strings.Join(headers, ", "))

// //     // Interactive question loop
// //     for {
// //         fmt.Print("\nEnter your question (or 'quit' to exit): ")
// //         scanner.Scan()
// //         question := scanner.Text()

// //         if strings.ToLower(question) == "quit" {
// //             break
// //         }

// //         // Process query
// //         result, err := orch.ProcessQuery(context.Background(), question, csvPath)
// //         if err != nil {
// //             fmt.Printf("Error processing query: %v\n", err)
// //             continue
// //         }

// //         // Print results
// //         printResults(result)
// //     }
// // }
// func main() {
//     // Initialize components
//     bridge := llm.CreateBridge("http://localhost:8000")
//     orch := orchestrator.CreateOrchestrator(bridge)

//     // Get CSV file path
//     fmt.Print("Enter the path to your CSV file: ")
//     scanner := bufio.NewScanner(os.Stdin)
//     scanner.Scan()
//     csvPath := scanner.Text()

//     // Validate CSV file
//     headers, err := getCSVHeaders(csvPath)
//     if err != nil {
//         fmt.Printf("Error reading CSV: %v\n", err)
//         return
//     }

//     fmt.Printf("\nFound columns: %v\n", strings.Join(headers, ", "))

//     // Interactive question loop
//     for {
//         fmt.Print("\nEnter your question (or 'quit' to exit): ")
//         scanner.Scan()
//         question := scanner.Text()

//         if strings.ToLower(question) == "quit" {
//             break
//         }

//         // Process query with new context
//         ctx := context.Background()
//         result, err := orch.ProcessQuery(ctx, question, csvPath)
//         if err != nil {
//             fmt.Printf("Error processing query: %v\n", err)
//             continue
//         }

//         // Print results
//         printResults(result)
//     }
// }

// func getCSVHeaders(filepath string) ([]string, error) {
//     file, err := os.Open(filepath)
//     if err != nil {
//         return nil, fmt.Errorf("error opening file: %w", err)
//     }
//     defer file.Close()

//     reader := csv.NewReader(file)
//     headers, err := reader.Read()
//     if err != nil {
//         return nil, fmt.Errorf("error reading headers: %w", err)
//     }

//     return headers, nil
// }

// // func printResults(state *types.State) {
// //     fmt.Println("\nResults:")
// //     fmt.Println("Generated SQL:", state.GeneratedQuery)
    
// //     if state.Error != "" {
// //         fmt.Println("Error:", state.Error)
// //         return
// //     }

// //     fmt.Println("\nExecution Results:")
// //     switch results := state.ExecutionResult.(type) {
// //     case []map[string]interface{}:
// //         printTableFormat(results)
// //     default:
// //         fmt.Printf("%+v\n", state.ExecutionResult)
// //     }
// // }
// // func printResults(state *types.State) {
// //     fmt.Println("\nResults:")
// //     fmt.Println("Generated SQL:", state.GeneratedQuery)
    
// //     if state.Error != "" {
// //         fmt.Println("Error:", state.Error)
// //         return
// //     }

// //     fmt.Println("\nTop Movies by Revenue:")
// //     fmt.Println("----------------------------------------")
// //     if results, ok := state.ExecutionResult.([]map[string]interface{}); ok {
// //         for i, movie := range results {
// //             fmt.Printf("%d. %s - %s\n", 
// //                 i+1, 
// //                 movie["title"], 
// //                 movie["revenue"])
// //         }
// //     }
// // }

// // func printResults(state *types.State) {
// //     fmt.Println("\nProcessing Steps:")
// //     fmt.Println("----------------------------------------")
// //     fmt.Printf("1. Analysis:\n   - Query understood and analyzed\n   - Schema extracted from CSV\n")
// //     fmt.Printf("2. Query Generation:\n   - SQL Generated: %s\n", state.GeneratedQuery)
// //     fmt.Printf("3. Validation:\n   - Query validated for syntax and schema\n")
// //     fmt.Printf("4. Execution:\n   - Query executed against CSV data\n")
    
// //     fmt.Println("\nResults:")
// //     fmt.Println("----------------------------------------")
// //     if results, ok := state.ExecutionResult.([]map[string]interface{}); ok && len(results) > 0 {
// //         for _, movie := range results {
// //             fmt.Printf("\nTitle: %s\n", movie["title"])
// //             fmt.Printf("Revenue: %s\n", movie["revenue"])
// //             if overview, ok := movie["overview"].(string); ok {
// //                 fmt.Printf("Overview: %s\n", overview)
// //             }
// //             fmt.Println("----------------------------------------")
// //         }
// //     } else {
// //         fmt.Println("No results found")
// //     }
// // }
// func printResults(state *types.State) {
//     fmt.Println("\nProcessing Steps:")
//     fmt.Println("----------------------------------------")
//     fmt.Printf("1. Analysis:\n   - Query understood and analyzed\n   - Schema extracted from CSV\n")
//     fmt.Printf("2. Query Generation:\n   - SQL Generated: %s\n", state.GeneratedQuery)
//     fmt.Printf("3. Validation:\n   - Query validated for syntax and schema\n")
//     fmt.Printf("4. Execution:\n   - Query executed against CSV data\n")
    
//     fmt.Println("\nResults:")
//     fmt.Println("----------------------------------------")
    
//     if results, ok := state.ExecutionResult.([]map[string]interface{}); ok && len(results) > 0 {
//         for _, result := range results {
//             fmt.Printf("\nTitle: %s\n", result["title"])
//             fmt.Printf("Revenue: %s\n", result["revenue"])
//             if genres, ok := result["genres"].(string); ok {
//                 fmt.Printf("Genres: %s\n", genres)
//             }
//             fmt.Println("----------------------------------------")
//         }
//     } else {
//         fmt.Println("No results found")
//         fmt.Println("Debug Info:")
//         fmt.Printf("Execution Result Type: %T\n", state.ExecutionResult)
//         fmt.Printf("State: %+v\n", state)
//     }
// }

// func printTableFormat(results []map[string]interface{}) {
//     if len(results) == 0 {
//         fmt.Println("No results found")
//         return
//     }

//     // Get headers from first result
//     var headers []string
//     for k := range results[0] {
//         headers = append(headers, k)
//     }

//     // Print headers
//     for _, h := range headers {
//         fmt.Printf("%-15s", h)
//     }
//     fmt.Println()

//     // Print separator
//     fmt.Println(strings.Repeat("-", len(headers)*15))

//     // Print rows
//     for _, row := range results {
//         for _, h := range headers {
//             fmt.Printf("%-15v", row[h])
//         }
//         fmt.Println()
//     }
// }

// backend/go/cmd/cli/main.go
package main

import (
    "bufio"
    "context"
    "encoding/csv"
    "fmt"
    "os"
    "strings"
    "time"

    "sage-ai-v2/internal/llm"
    "sage-ai-v2/internal/orchestrator"
    "sage-ai-v2/internal/types"
    "sage-ai-v2/pkg/logger"
)

func main() {
    // Initialize components
    bridge := llm.CreateBridge("http://localhost:8000")
    orch := orchestrator.CreateOrchestrator(bridge)

    for {
        // Start new session
        sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
        logger.InfoLogger.Printf("Starting new session: %s", sessionID)
        orch.NewSession()

        // Get CSV file path
        fmt.Print("Enter the path to your CSV file (or 'quit' to exit): ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        csvPath := scanner.Text()

        if strings.ToLower(csvPath) == "quit" {
            break
        }

        // Validate CSV file
        headers, err := getCSVHeaders(csvPath)
        if err != nil {
            fmt.Printf("Error reading CSV: %v\n", err)
            continue
        }

        fmt.Printf("\nFound columns: %v\n", strings.Join(headers, ", "))

        // Interactive question loop for this session
        for {
            fmt.Print("\nEnter your question (or 'new' for new session, 'quit' to exit): ")
            scanner.Scan()
            question := scanner.Text()

            switch strings.ToLower(question) {
            case "quit":
                logger.InfoLogger.Printf("Ending session: %s", sessionID)
                return
            case "new":
                logger.InfoLogger.Printf("Ending current session: %s", sessionID)
                orch.ClearSession()
                fmt.Println("\nStarting new session...")
                break
            default:
                // Create context with timeout
                ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
                
                // Process query
                result, err := orch.ProcessQuery(ctx, question, csvPath)
                
                // Cancel context after processing
                cancel()
                
                if err != nil {
                    fmt.Printf("Error processing query: %v\n", err)
                    continue
                }

                // Print results
                printResults(result)
            }

            // Break out of question loop if starting new session
            if strings.ToLower(question) == "new" {
                break
            }
        }
    }
}

func getCSVHeaders(filepath string) ([]string, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, fmt.Errorf("error opening file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    headers, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("error reading headers: %w", err)
    }

    return headers, nil
}

func printResults(state *types.State) {
    fmt.Println("\nProcessing Steps:")
    fmt.Println("----------------------------------------")
    fmt.Printf("1. Analysis:\n   - Query understood and analyzed\n   - Schema extracted from CSV\n")
    fmt.Printf("2. Query Generation:\n   - SQL Generated: %s\n", state.GeneratedQuery)
    fmt.Printf("3. Validation:\n   - Query validated for syntax and schema\n")
    fmt.Printf("4. Execution:\n   - Query executed against CSV data\n")
    
    fmt.Println("\nResults:")
    fmt.Println("----------------------------------------")
    
    if results, ok := state.ExecutionResult.([]map[string]interface{}); ok && len(results) > 0 {
        for _, result := range results {
            fmt.Printf("\nTitle: %s\n", result["title"])
            if revenue, ok := result["revenue"].(string); ok {
                fmt.Printf("Revenue: %s\n", revenue)
            }
            if genres, ok := result["genres"].(string); ok {
                fmt.Printf("Genres: %s\n", genres)
            }
            if vote_average, ok := result["vote_average"].(string); ok {
                fmt.Printf("Rating: %s\n", vote_average)
            }
            fmt.Println("----------------------------------------")
        }
    } else {
        fmt.Println("No results found")
        fmt.Println("Debug Info:")
        fmt.Printf("Execution Result Type: %T\n", state.ExecutionResult)
        fmt.Printf("State: %+v\n", state)
    }
}

func printTableFormat(results []map[string]interface{}) {
    if len(results) == 0 {
        fmt.Println("No results found")
        return
    }

    // Get headers from first result
    var headers []string
    for k := range results[0] {
        headers = append(headers, k)
    }

    // Print headers
    for _, h := range headers {
        fmt.Printf("%-15s", h)
    }
    fmt.Println()

    // Print separator
    fmt.Println(strings.Repeat("-", len(headers)*15))

    // Print rows
    for _, row := range results {
        for _, h := range headers {
            fmt.Printf("%-15v", row[h])
        }
        fmt.Println()
    }
}