package command

import (
	"github.com/spf13/cobra"
)

func newListCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured tools status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return rootCtx.app.ListTool(cmd.Context())
		},
	}
}
