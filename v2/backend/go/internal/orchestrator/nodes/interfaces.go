// backend/go/internal/orchestrator/nodes/interfaces.go
package nodes

import (
    "context"
    "sage-ai-v2/internal/types"
)

type NodeExecutor interface {
    Execute(ctx context.Context, state *types.State) error
}