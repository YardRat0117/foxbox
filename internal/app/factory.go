package app

import (
	"github.com/YardRat0117/foxbox/internal/runtime"
	types "github.com/YardRat0117/foxbox/internal/foxtypes"
)

// New constructs a new App instance with configuration and runtime.
func New(cfg *types.Config, rt runtime.Runtime) *App {
	return &App{cfg: cfg, rt: rt}
}
