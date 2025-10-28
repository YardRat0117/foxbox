package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newInstallCommand is a factory func that provides dependency injection into the `install` command
func newInstallCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "install <tool>",
		Short: "Install (pull) a tool's container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			panel, cfg := rootCtx.panel, rootCtx.cfg

			toolName, toolVer := parseToolArg(args[0])
			tool, ok := cfg.Tools[toolName]
			if !ok {
				return fmt.Errorf("Tool `%s` not configured\n", toolName)
			}
			if err := panel.InstallTool(tool.Image, toolVer); err != nil {
				return fmt.Errorf("Error installing tool: %e", err)
			}
			fmt.Printf("Images %s installed successfully!\n", tool.Image)
			return nil
		},
	}
}
