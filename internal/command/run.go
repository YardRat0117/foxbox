package command

import (
	"os"

	"github.com/spf13/cobra"
)

func newRunCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "run <tool> [args]",
		Short: "Run a tool inside its sandbox",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName, toolTag := 
			code, err := rootCtx.app.RunTool(cmd.Context(), rootCtx.app.NewRunToolRequest(toolName, toolArgs))
			if err != nil {
				return err
			}
			os.Exit(code)
			return nil
		},
	}
}
