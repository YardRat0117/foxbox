package command

import (
	"github.com/spf13/cobra"
)

func newRunCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "run <tool> [args]",
		Short: "Run a tool inside its container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCtx.app.RunTool(cmd.Context(), args)
		},
	}
}
