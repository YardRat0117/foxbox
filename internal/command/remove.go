package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <tool>",
	Short: "Remove a tool's container image",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		toolName, toolVer := parseToolArg(args[0])
		tool, ok := cfg.Tools[toolName]

		if !ok {
			fatal(fmt.Sprintf("Tool `%s` not configured\n", toolName), nil)
		}

		if err := panel.RemoveTool(toolName, tool.Image, toolVer); err != nil {
			fatal("Error removing tool: ", err)
		}
		fmt.Printf("Image %s removed successfully!\n", tool.Image)
	},
}
