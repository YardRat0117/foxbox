package container

import (
	"github.com/YardRat0117/foxbox/internal/types"
)

// imageManager manages images, NOTHING to do with `tool`
type imageManager interface {
	checkImage(image string) (bool, error)
	pullImage(image string) error
	removeImage(image string) error
	getLocalImages() (map[string]map[string]struct{}, error)
}

// toolManager manages images, and calls ImageManagers to work
type toolManager interface {
	installTool(toolName string, version string) error
	removeTool(toolName string, imgName string, version string) error
	runTool(tool types.Tool, version string, args []string) error
	checkTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error)
	cleanTools(tools map[string]types.Tool) error
}

// Runtime represents the interface where commands interact with container
type Runtime interface {
	InstallTool(toolName, version string) error
	RemoveTool(toolName, imgName, version string) error
	RunTool(tool types.Tool, version string, args []string) error
	CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error)
	CleanTools(tools map[string]types.Tool) error
}
