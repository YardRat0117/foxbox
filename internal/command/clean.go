package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// func newCleanCommand creates a `cobra.Command` object `clean` with given context
func newCleanCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean all installed container image",
		RunE: func(_ *cobra.Command, _ []string) error {
			// Assign panel and config
			panel, cfg := rootCtx.panel, rootCtx.cfg

			// Try call panel to clean tools
			if err := panel.CleanTools(cfg.Tools); err != nil {
				return fmt.Errorf("Error cleaning installed tools: %w", err)
			}

			// Hint
			fmt.Println("Tools cleaned successfully!")

			return nil
		},
	}
}
