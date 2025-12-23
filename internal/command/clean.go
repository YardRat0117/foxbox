package command

import (
	"github.com/spf13/cobra"
)

func newCleanCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean all configured tools' corresponding images",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return rootCtx.app.CleanTool(cmd.Context())
		},
	}
}
