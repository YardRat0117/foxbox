package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRunCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "run <tool> [args..] -- [toolArgs...]",
		Short: "Run a tool inside its container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			panel, cfg := rootCtx.panel, rootCtx.cfg

			toolName, toolVer := parseToolArg(args[0])
			toolArgs := args[1:]
			tool, ok := cfg.Tools[toolName]
			if !ok {
				return fmt.Errorf("Tool `%s` not found\n", toolName)
			}
			if err := panel.RunTool(tool, toolVer, toolArgs); err != nil {
				return fmt.Errorf("Error running tool: %e", err)
			}
			return nil
		},
	}
}
