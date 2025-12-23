// Package runtime defines the execution capabilities of a container runtime.
package runtime

import (
	"context"
	"github.com/YardRat0117/foxbox/internal/domain"
)

// Runtime defines the abstract capabilities of a container runtime.
type Runtime interface {
	// image lifecycle
	EnsureImage(ctx context.Context, ref domain.ImageRef) error
	RemoveImage(ctx context.Context, ref domain.ImageRef) error
	ListImage(ctx context.Context) ([]domain.ImageInfo, error)

	Create(ctx context.Context, spec domain.ContainerSpec) (domain.ContainerID, error)
	Start(ctx context.Context, id domain.ContainerID) error
	Stop(ctx context.Context, id domain.ContainerID) error
	Remove(ctx context.Context, id domain.ContainerID) error

	Exec(id domain.ContainerID) (Execution, error)
}
