//go:build legacy

package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newInstallCommand creates a `cobra.Command` object `install` with given context
func newInstallCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "install <tool>",
		Short: "Install (pull) a tool's container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Assign panel and config
			panel, cfg := rootCtx.panel, rootCtx.cfg

			// Split tool info
			toolName, toolVer := parseToolArg(args[0])

			// Check if tool configured
			tool, ok := cfg.Tools[toolName]
			if !ok {
				return fmt.Errorf("Tool `%s` not configured\n", toolName)
			}

			// Call panel to install tool
			if err := panel.InstallTool(tool.Image, toolVer); err != nil {
				return fmt.Errorf("Error installing tool: %e", err)
			}

			// Hint
			fmt.Printf("Images %s installed successfully!\n", tool.Image)

			return nil
		},
	}
}
