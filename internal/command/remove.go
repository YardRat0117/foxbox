package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newRemoveCmd is a factory func that provides dependency inection into the `remove` command
func newRemoveCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <tool>",
		Short: "Remove a tool's container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			panel, cfg := rootCtx.panel, rootCtx.cfg

			toolName, toolVer := parseToolArg(args[0])
			tool, ok := cfg.Tools[toolName]

			if !ok {
				return fmt.Errorf("Tool `%s` not configured\n", toolName)
			}

			if err := panel.RemoveTool(toolName, tool.Image, toolVer); err != nil {
				return fmt.Errorf("Error removing tool: %e", err)
			}
			fmt.Printf("Image %s removed successfully!\n", tool.Image)
			return nil
		},
	}
}
