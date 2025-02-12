// // backend/go/internal/orchestrator/graph.go
// package orchestrator

// import (
//     "context"
//     "fmt"
// )

// type NodeFunc func(context.Context, *State) error

// type Graph struct {
//     nodes map[string]NodeFunc
//     edges map[string][]string
// }

// func NewGraph() *Graph {
//     return &Graph{
//         nodes: make(map[string]NodeFunc),
//         edges: make(map[string][]string),
//     }
// }

// func (g *Graph) AddNode(name string, fn NodeFunc) {
//     g.nodes[name] = fn
// }

// func (g *Graph) AddEdge(from, to string) {
//     g.edges[from] = append(g.edges[from], to)
// }

// func (g *Graph) Execute(ctx context.Context, state *State) error {
//     visited := make(map[string]bool)

//     var execute func(string) error
//     execute = func(node string) error {
//         if visited[node] {
//             return nil
//         }

//         visited[node] = true
//         if fn, ok := g.nodes[node]; ok {
//             if err := fn(ctx, state); err != nil {
//                 return fmt.Errorf("node %s failed: %w", node, err)
//             }

//             for _, next := range g.edges[node] {
//                 if err := execute(next); err != nil {
//                     return err
//                 }
//             }
//         }
//         return nil
//     }

//     return execute("analyzer")
// }

// backend/go/internal/orchestrator/graph.go
package orchestrator

import (
	"context"
	"fmt"
	"sage-ai-v2/internal/types"
	"sage-ai-v2/pkg/logger"
)

type NodeFunc func(context.Context, *types.State) error

type Graph struct {
    nodes map[string]NodeFunc
    edges map[string][]string
}

func CreateGraph() *Graph {
    return &Graph{
        nodes: make(map[string]NodeFunc),
        edges: make(map[string][]string),
    }
}

func (g *Graph) AddNode(name string, fn NodeFunc) {
    g.nodes[name] = fn
}

func (g *Graph) AddEdge(from, to string) {
    g.edges[from] = append(g.edges[from], to)
}

// func (g *Graph) Execute(ctx context.Context, state *types.State) error {
//     visited := make(map[string]bool)

//     var execute func(string) error
//     execute = func(node string) error {
//         if visited[node] {
//             return nil
//         }

//         visited[node] = true
//         if fn, ok := g.nodes[node]; ok {
//             if err := fn(ctx, state); err != nil {
//                 return fmt.Errorf("node %s failed: %w", node, err)
//             }

//             for _, next := range g.edges[node] {
//                 if err := execute(next); err != nil {
//                     return err
//                 }
//             }
//         }
//         return nil
//     }

//     return execute("analyzer")
// }

// func (g *Graph) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Graph: Starting execution pipeline")
//     visited := make(map[string]bool)

//     var execute func(string) error
//     execute = func(node string) error {
//         if visited[node] {
//             return nil
//         }

//         visited[node] = true
//         logger.InfoLogger.Printf("Graph: Executing node: %s", node)

//         if fn, ok := g.nodes[node]; ok {
//             if err := fn(ctx, state); err != nil {
//                 logger.ErrorLogger.Printf("Graph: Node %s failed: %v", node, err)
//                 return fmt.Errorf("node %s failed: %w", node, err)
//             }

//             logger.InfoLogger.Printf("Graph: Node %s completed successfully", node)

//             // Execute next nodes
//             if nextNodes, ok := g.edges[node]; ok {
//                 for _, next := range nextNodes {
//                     logger.DebugLogger.Printf("Graph: Moving to next node: %s -> %s", node, next)
//                     if err := execute(next); err != nil {
//                         return err
//                     }
//                 }
//             }
//         }
//         return nil
//     }

//     err := execute("analyzer")
//     if err != nil {
//         logger.ErrorLogger.Printf("Graph: Execution pipeline failed: %v", err)
//         return err
//     }

//     logger.InfoLogger.Printf("Graph: Execution pipeline completed successfully")
//     return nil
// }

// func (g *Graph) Execute(ctx context.Context, state *types.State) error {
//     logger.InfoLogger.Printf("Graph: Starting execution pipeline")
//     visited := make(map[string]bool)

//     var execute func(string) error
//     execute = func(node string) error {
//         if visited[node] {
//             return nil
//         }

//         visited[node] = true
//         logger.InfoLogger.Printf("Graph: Executing node: %s", node)

//         if fn, ok := g.nodes[node]; ok {
//             if err := fn(ctx, state); err != nil {
//                 logger.ErrorLogger.Printf("Node %s failed: %v", node, err)
//                 return fmt.Errorf("node %s failed: %w", node, err)
//             }

//             logger.InfoLogger.Printf("Node %s completed successfully", node)

//             for _, next := range g.edges[node] {
//                 if err := execute(next); err != nil {
//                     return err
//                 }
//             }
//         }
//         return nil
//     }

//     if err := execute("analyzer"); err != nil {
//         return err
//     }

//     logger.InfoLogger.Printf("Graph: Execution pipeline completed successfully")
//     return nil
// }

func (g *Graph) Execute(ctx context.Context, state *types.State) error {
    logger.InfoLogger.Printf("Graph: Starting execution pipeline")

    // Define execution order
    executionOrder := []string{"analyzer", "generator", "validator", "executor"}

    for _, nodeName := range executionOrder {
        if fn, ok := g.nodes[nodeName]; ok {
            logger.InfoLogger.Printf("Graph: Executing node: %s", nodeName)
            
            if err := fn(ctx, state); err != nil {
                logger.ErrorLogger.Printf("Graph: Node %s failed: %v", nodeName, err)
                return fmt.Errorf("node %s failed: %w", nodeName, err)
            }

            logger.InfoLogger.Printf("Graph: Node %s completed successfully", nodeName)

            // Check if state contains error after node execution
            if state.Error != "" {
                logger.ErrorLogger.Printf("Graph: State contains error after %s: %s", nodeName, state.Error)
                return fmt.Errorf("state error after %s: %s", nodeName, state.Error)
            }
        } else {
            logger.ErrorLogger.Printf("Graph: Node %s not found in graph", nodeName)
        }
    }

    logger.InfoLogger.Printf("Graph: Execution pipeline completed successfully")
    return nil
}