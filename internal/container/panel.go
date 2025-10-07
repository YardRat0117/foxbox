package container

import (
	"github.com/YardRat0117/foxbox/internal/types"
)

var _ PanelInterface = (*Panel)(nil)

type Panel struct {
	toolManager *toolManager
}

func (c *Panel) InstallTool(toolName, version string) error {
	return c.toolManager.installTool(toolName, version)
}

func (c *Panel) RemoveTool(toolName, imgName, version string) error {
	return c.toolManager.removeTool(toolName, imgName, version)
}

func (c *Panel) RunTool(tool types.Tool, version string, args []string) error {
	return c.toolManager.runTool(tool, version, args)
}

func (c *Panel) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	return c.toolManager.checkTools(tools)
}

func (c *Panel) CleanTools(tools map[string]types.Tool) error {
	return c.toolManager.cleanTools(tools)
}
