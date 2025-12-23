package runtime

import (
	"context"
)

// Execution representes a running container execution instance.
type Execution interface {
	Attach(ctx context.Context) error
	Wait(ctx context.Context) (int, error)
	Close(ctx context.Context) error
}
