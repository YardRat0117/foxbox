package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all installed container image",
	Run: func(_ *cobra.Command, _ []string) {
		// `panel` defined in `rootCmd`
		cleanTools(panel)
	},
}

func cleanTools(panel *container.Panel) {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config: ", err)
		os.Exit(1)
	}

	// Delete all tools
	if err := panel.CleanTools(cfg.Tools); err != nil {
		fmt.Println("Error cleaning installed tools: ", err)
		os.Exit(1)
	}

	fmt.Println("Tools cleaned successfully!")
}
