package app

import (
	"github.com/YardRat0117/foxbox/internal/runtime"
	"github.com/YardRat0117/foxbox/internal/types"
)

// New constructs a new App instance with configuration and runtime.
func New(cfg *types.Config, rt runtime.Runtime) *App {
	return &App{cfg: cfg, rt: rt}
}
