package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all installed container image",
	Run: func(_ *cobra.Command, _ []string) {
		if err := panel.CleanTools(cfg.Tools); err != nil {
			fatal("Error cleaning installed tools:", err)
		}
		fmt.Println("Tools cleaned successfully!")
	},
}
