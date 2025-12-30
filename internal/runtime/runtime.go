// Package runtime defines the execution capabilities of a container runtime.
package runtime

import (
	"context"
	"github.com/YardRat0117/foxbox/internal/domain"
)

// Runtime defines the abstract capabilities of a container runtime.
type Runtime interface {
	// execution environment lifecycle
	EnsureEnv(ctx context.Context, ref domain.EnvRef) error
	RemoveEnv(ctx context.Context, ref domain.EnvRef) error
	HasEnv(ctx context.Context, ref domain.EnvRef) (bool, error)

	// sandbox lifecycle
	CreateSandbox(ctx context.Context, spec domain.SandboxSpec) (domain.SandboxID, error)
	StartSandbox(ctx context.Context, id domain.SandboxID) error
	WaitSandbox(ctx context.Context, id domain.SandboxID) (domain.SandboxResult, error)
	RemoveSandbox(ctx context.Context, id domain.SandboxID) error
}
