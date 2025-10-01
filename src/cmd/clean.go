package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/config"
	"github.com/YardRat0117/foxbox/src/container"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all installed container image",
	Run: func(_ *cobra.Command, _ []string) {
		// `runtime` defined in `rootCmd`
		cleanTools(runtime)
	},
}

func cleanTools(runtime container.Runtime) {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config: ", err)
		os.Exit(1)
	}

	// Delete all tools
	if err := runtime.CleanTools(cfg.Tools); err != nil {
		fmt.Println("Error cleaning installed tools: ", err)
		os.Exit(1)
	}

	fmt.Println("Tools cleaned successfully!")
}
