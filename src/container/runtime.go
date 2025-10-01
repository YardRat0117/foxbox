package container

import (
	"github.com/YardRat0117/foxbox/src/types"
)

// imageManager manages images, NOTHING to do with `tool`
type imageManager interface {
	checkImage(image string) (bool, error)
	pullImage(image string) error
	removeImage(image string) error
	getLocalImages() (map[string]map[string]struct{}, error)
}

// ToolManager manages images, and calls ImageManagers to work
type ToolManager interface {
	InstallTool(toolName string, version string) error
	RemoveTool(toolName string, imgName string, version string) error
	RunTool(tool types.Tool, version string, args []string) error
	CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error)
	CleanTools(tools map[string]types.Tool) error
}

// Runtime represents the interface where commands interact with container
type Runtime interface {
	ToolManager
}
