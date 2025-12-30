// Package app implements the core application use cases of foxbox.
package app

import (
	"context"
	"fmt"

	"github.com/YardRat0117/foxbox/internal/domain"
	types "github.com/YardRat0117/foxbox/internal/foxtypes"
	"github.com/YardRat0117/foxbox/internal/runner"
)

// App orchestrates tool execution inside containers.
type App struct {
	cfg    *types.Config
	runner *runner.Runner
}

// RunTool executes a configured tool inside a container.
func (a *App) RunTool(ctx context.Context, req RunToolRequest) (int, error) {
	// Foxbox internal error
	const exitCode = 70

	runReq := runner.RunRequest{
		Env:     req.Env,
		Cmd:     append([]string{req.Entry}, req.ToolArgs...),
		WorkDir: req.WorkDir,
		Mounts:  req.Mounts,
	}

	// Currently only exit code are directly returned
	// More detailed info would be handled by runner
	code, err := a.runner.Run(ctx, runReq)
	if err != nil {
		return exitCode, err
	}

	// Retain actual tool behavior
	return code, nil
}

// ListTool lists all configured tool.
func (a *App) ListTool(ctx context.Context) error {
	const nameWidth, parenWidth = 10, 15
	fmt.Println("Configured tools:")

	for name, tool := range a.cfg.Tools {
		installed, err := a.runner.HasEnv(ctx, tool.Env)
		var status string
		if err == nil {
			if installed {
				status = "[installed]"
			} else {
				status = "[not installed]"
			}
		} else {
			status = "[error]"
		}

		fmt.Printf(
			"- %-*s %-*s %s\n",
			nameWidth, name,
			parenWidth, fmt.Sprintf("(%s)", tool.Env),
			status,
		)
	}

	return nil
}

// InstallTool installs a configured tool.
func (a *App) InstallTool(ctx context.Context, ref domain.EnvRef) error {
	return a.runner.EnsureEnv(ctx, ref)
}

// RemoveTool removes a configured tool.
func (a *App) RemoveTool(ctx context.Context, ref domain.EnvRef) error {
	return a.runner.RemoveEnv(ctx, ref)
}

// CleanTool cleans all configured tool.
func (a *App) CleanTool(ctx context.Context) error {
	for _, tool := range a.cfg.Tools {
		installed, err := a.runner.HasEnv(ctx, tool.Env)
		if err != nil {
			return err
		}
		if installed {
			if err := a.runner.RemoveEnv(ctx, tool.Env); err != nil {
				return err
			}

		}
	}

	return nil
}
