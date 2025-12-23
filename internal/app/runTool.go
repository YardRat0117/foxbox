// Package app implements the core application use cases of foxbox.
package app

import (
	"context"
	"fmt"

	"github.com/YardRat0117/foxbox/internal/domain"
	"github.com/YardRat0117/foxbox/internal/runtime"
	"github.com/YardRat0117/foxbox/internal/types"
)

// App orchestrates tool execution inside containers.
type App struct {
	cfg *types.Config
	rt  runtime.Runtime
}

// New constructs a new App instance with configuration and runtime.
func New(cfg *types.Config, rt runtime.Runtime) *App {
	return &App{cfg: cfg, rt: rt}
}

// RunTool executes a configured tool inside a container.
func (a *App) RunTool(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no tool specified")
	}

	toolName, toolVer := parseToolArg(args[0])
	toolArgs := args[1:]

	tool, ok := a.cfg.Tools[toolName]
	if !ok {
		return fmt.Errorf("tool %q not configured", toolName)
	}

	imgRef, err := domain.NewImageRef(tool.Image + ":" + toolVer)
	if err != nil {
		return fmt.Errorf("invalid image reference for tool %q: %w", toolName, err)
	}

	if err := a.rt.EnsureImage(ctx, imgRef); err != nil {
		return fmt.Errorf("failed to ensure image %q: %w", imgRef.Raw, err)
	}

	spec := domain.ContainerSpec{
		Image:   imgRef,
		Cmd:     append([]string{tool.Entry}, toolArgs...),
		WorkDir: tool.Workdir,
		Volumes: parseVolumes(tool.Volumes),
	}

	id, err := a.rt.Create(ctx, spec)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	defer func() {
		_ = a.rt.Stop(ctx, id)
		_ = a.rt.Remove(ctx, id)
	}()

	exec, err := a.rt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to exec in container: %w", err)
	}
	defer func() {
		_ = exec.Close(ctx)
	}()

	if err := exec.Attach(ctx); err != nil {
		return fmt.Errorf("failed to attach to container: %w", err)
	}

	if err := a.rt.Start(ctx, id); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	if err := exec.Wait(ctx); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}
