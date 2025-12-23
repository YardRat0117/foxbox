package command

import (
	"github.com/spf13/cobra"
)

func newInstallCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "install <tool>",
		Short: "Install a tool's corresponding container image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCtx.app.InstallTool(cmd.Context(), args)
		},
	}
}
