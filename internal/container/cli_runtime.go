package container

import (
	"github.com/YardRat0117/foxbox/internal/types"
)

func NewDockerRuntime() Runtime {
	return newCliRuntime("docker")
}

func NewPodmanRuntime() Runtime {
	return newCliRuntime("podman")
}

type cliRuntime struct {
	toolManager *cliToolManager
}

func newCliRuntime(runtimeName string) Runtime {
	imgMgr := &cliImageManager{
		runtimeName: runtimeName,
	}

	return &cliRuntime{
		toolManager: &cliToolManager{
			runtimeName: runtimeName,
			imageMgr:    imgMgr,
		},
	}
}

func (c *cliRuntime) InstallTool(toolName, version string) error {
	return c.toolManager.installTool(toolName, version)
}

func (c *cliRuntime) RemoveTool(toolName, imgName, version string) error {
	return c.toolManager.removeTool(toolName, imgName, version)
}

func (c *cliRuntime) RunTool(tool types.Tool, version string, args []string) error {
	return c.toolManager.runTool(tool, version, args)
}

func (c *cliRuntime) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	return c.toolManager.checkTools(tools)
}

func (c *cliRuntime) CleanTools(tools map[string]types.Tool) error {
	return c.toolManager.cleanTools(tools)
}
