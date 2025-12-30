package app

import (
	"github.com/YardRat0117/foxbox/internal/domain"
)

// RunToolRequest is for CLI-to-APP request
type RunToolRequest struct {
	ToolArgs []string
	Entry    string
	Env      domain.EnvRef
	Mounts   []domain.MountSpec
	WorkDir  string
}
