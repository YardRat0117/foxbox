package command

import (
	"github.com/spf13/cobra"
)

func newRemoveCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <tool>",
		Short: "Remove a tool's corresponding container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCtx.app.RemoveTool(cmd.Context(), args)
		},
	}
}
