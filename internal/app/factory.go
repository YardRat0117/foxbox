package app

import (
	types "github.com/YardRat0117/foxbox/internal/foxtypes"
	"github.com/YardRat0117/foxbox/internal/runner"
)

// New constructs a new App instance with configuration and runtime.
func New(cfg *types.Config, runner *runner.Runner) *App {
	return &App{cfg: cfg, runner: runner}
}
