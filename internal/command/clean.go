package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCleanCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean all installed container image",
		RunE: func(_ *cobra.Command, _ []string) error {
			panel, cfg := rootCtx.panel, rootCtx.cfg

			if err := panel.CleanTools(cfg.Tools); err != nil {
				return fmt.Errorf("Error cleaning installed tools: %w", err)
			}
			fmt.Println("Tools cleaned successfully!")
			return nil
		},
	}
}
