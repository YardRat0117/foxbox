package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newRemoveCommand creates a `cobra.Command` object `remove` with given context
func newRemoveCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <tool>",
		Short: "Remove a tool's container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Assign panel and config
			panel, cfg := rootCtx.panel, rootCtx.cfg

			// Split tool info
			toolName, toolVer := parseToolArg(args[0])

			// Check if tool configured
			tool, ok := cfg.Tools[toolName]
			if !ok {
				return fmt.Errorf("Tool `%s` not configured\n", toolName)
			}

			// Call panel to remove tool
			if err := panel.RemoveTool(toolName, tool.Image, toolVer); err != nil {
				return fmt.Errorf("Error removing tool: %e", err)
			}

			// Hint
			fmt.Printf("Image %s removed successfully!\n", tool.Image)

			return nil
		},
	}
}
