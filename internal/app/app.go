// Package app implements the core application use cases of foxbox.
package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/YardRat0117/foxbox/internal/domain"
	"github.com/YardRat0117/foxbox/internal/runtime"
	"github.com/YardRat0117/foxbox/internal/types"
)

// App orchestrates tool execution inside containers.
type App struct {
	cfg *types.Config
	rt  runtime.Runtime
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

// ListTool lists all configured tool.
func (a *App) ListTool(ctx context.Context, args []string) error {
	images, err := a.rt.ListImage(ctx)
	if err != nil {
		return err
	}

	imageTags := make(map[string][]string)
	for _, img := range images {
		imageName := img.Ref.Raw
		imageTag := img.Tag
		// Slice for raw image name without tag
		if idx := strings.LastIndex(imageName, ":"); idx != -1 {
			imageName = imageName[:idx]
		}
		if idx := strings.LastIndex(imageTag, ":"); idx != -1 {
			imageTag = imageTag[idx+1:]
		}
		imageTags[imageName] = append(imageTags[imageName], imageTag)
	}

	const nameWidth, parenWidth = 10, 15
	fmt.Println("Configured tools:")

	for name, tool := range a.cfg.Tools {
		status := "[not installed]"
		tags := ""

		if tagList, exists := imageTags[tool.Image]; exists {
			status = "[installed]    "
			if len(tagList) > 0 {
				tags = "tags: " + strings.Join(tagList, ", ")
			}
		}

		fmt.Printf("- %-*s %-*s %s %s\n",
			nameWidth, name,
			parenWidth, fmt.Sprintf("(%s)", tool.Image),
			status, tags,
		)
	}

	return nil
}
