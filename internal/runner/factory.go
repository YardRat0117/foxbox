package runner

import (
	"github.com/YardRat0117/foxbox/internal/runtime"
)

// New initializez a new runner (dependency injection)
func New(rt runtime.Runtime) *Runner {
	return &Runner{rt: rt}
}
