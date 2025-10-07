package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install (pull) a tool's container image",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		toolName, toolVer := parseToolArg(args[0])
		tool, ok := cfg.Tools[toolName]
		if !ok {
			fatal(fmt.Sprintf("Tool `%s` not configured\n", toolName), nil)
		}
		if err := panel.InstallTool(tool.Image, toolVer); err != nil {
			fatal("Error installing tool: ", err)
		}
		fmt.Printf("Images %s installed successfully!\n", tool.Image)
	},
}
