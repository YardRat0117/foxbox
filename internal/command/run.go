package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <tool> [args..] -- [toolArgs...]",
	Short: "Run a tool inside its container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		toolName, toolVer := parseToolArg(args[0])
		toolArgs := args[1:]
		tool, ok := cfg.Tools[toolName]
		if !ok {
			fatal(fmt.Sprintf("Tool `%s` not found\n", toolName), nil)
		}
		if err := panel.RunTool(tool, toolVer, toolArgs); err != nil {
			fatal("Error running tool:", err)
		}
	},
}
