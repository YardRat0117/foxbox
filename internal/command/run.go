package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newRunCommand creates a `cobra.Command` object `run` with given context
func newRunCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "run <tool> [args..] -- [toolArgs...]",
		Short: "Run a tool inside its container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Assign panel and config
			panel, cfg := rootCtx.panel, rootCtx.cfg

			// Split tool info
			toolName, toolVer := parseToolArg(args[0])
			toolArgs := args[1:]

			// Check if tool configured
			tool, ok := cfg.Tools[toolName]
			if !ok {
				return fmt.Errorf("Tool `%s` not configured\n", toolName)
			}

			// Call panel to run tool
			if err := panel.RunTool(tool, toolVer, toolArgs); err != nil {
				return fmt.Errorf("Error running tool: %e", err)
			}

			return nil
		},
	}
}
