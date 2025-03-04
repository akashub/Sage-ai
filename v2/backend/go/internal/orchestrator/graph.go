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