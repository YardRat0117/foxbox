// Package runner implements struct Runner
package runner

import (
	"context"

	"github.com/YardRat0117/foxbox/internal/domain"
	"github.com/YardRat0117/foxbox/internal/runtime"
)

// Runner is refactored and extracted for possible RPC interfaces
type Runner struct {
	rt runtime.Runtime
}

// RunRequest is the request sent to runtime for running a sandbox
type RunRequest struct {
	Env     domain.EnvRef
	Cmd     []string
	Version string
	WorkDir string
	Mounts  []domain.MountSpec
}

// Run sends request to and returns response from the runtime
func (r *Runner) Run(ctx context.Context, req RunRequest) (int, error) {
	if err := r.rt.EnsureEnv(ctx, req.Env); err != nil {
		return 0, err
	}

	spec := domain.SandboxSpec{
		Env:     req.Env,
		Cmd:     req.Cmd,
		WorkDir: req.WorkDir,
		Mounts:  req.Mounts,
	}

	id, err := r.rt.CreateSandbox(ctx, spec)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = r.rt.RemoveSandbox(context.Background(), id)
	}()

	if err := r.rt.StartSandbox(ctx, id); err != nil {
		return 0, err
	}

	// runtime.WaitSandbox returns a struct SandboxResult
	// More info than mere exit codeis fine to add
	// but should be handled by runner
	result, err := r.rt.WaitSandbox(ctx, id)
	if err != nil {
		return 0, err
	}

	return result.ExitCode, nil
}

// HasEnv passes by to runtime.HasEnv
func (r *Runner) HasEnv(ctx context.Context, ref domain.EnvRef) (bool, error) {
	return r.rt.HasEnv(ctx, ref)
}

// EnsureEnv passes by to runtime.EnsureEnv
func (r *Runner) EnsureEnv(ctx context.Context, ref domain.EnvRef) error {
	return r.rt.EnsureEnv(ctx, ref)
}

// RemoveEnv passes by to runtime.RemoveEnv
func (r *Runner) RemoveEnv(ctx context.Context, ref domain.EnvRef) error {
	return r.rt.RemoveEnv(ctx, ref)
}
