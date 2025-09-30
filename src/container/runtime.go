package container

import (
	"github.com/YardRat0117/foxbox/src/types"
)

// ImageManager manages images, NOTHING to do with `tool`
type ImageManager interface {
	ImageExists(image string) (bool, error)
	PullImage(image string) error
	RemoveImage(image string) error
}

// Runner builds and conducts commands
type Runner interface {
	RunTool(tool types.Tool, version string, args []string) error
}

// ToolInspector reflects tool info
type ToolInspector interface {
	CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error)
}

// Runtime represents the integrated interfaces
type Runtime interface {
	ImageManager
	Runner
	ToolInspector
}
