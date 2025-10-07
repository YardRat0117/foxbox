package container

import (
	"github.com/YardRat0117/foxbox/internal/types"
)

type PanelInterface interface {
	InstallTool(toolName, version string) error
	RemoveTool(toolName, imgName, version string) error
	RunTool(tool types.Tool, version string, args []string) error
	CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error)
	CleanTools(tools map[string]types.Tool) error
}
