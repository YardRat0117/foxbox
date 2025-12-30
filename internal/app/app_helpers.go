package app

import (
	"strings"
)

// ParseToolArg converts CLI tool config arguments into name and version pair.
func parseToolArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)

	name = parts[0]

	version = "latest"
	if len(parts) == 2 && parts[1] != "" {
		version = parts[1]
	}

	return
}

// NewRunToolRequest builds a new RunToolRequest
func (a *App) NewRunToolRequest(toolName string, args []string) RunToolRequest {
	tool := a.cfg.Tools[toolName]

	return RunToolRequest{
		ToolArgs: args,
		Entry:    tool.Entry,
		Env:      tool.Env,
		Mounts:   tool.Mounts,
		WorkDir:  tool.Workdir,
	}
}
