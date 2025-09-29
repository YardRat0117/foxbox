package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/ratbox/src/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tools in config",
	Run: func(cmd *cobra.Command, args []string) {
		ListTools()
	},
}

func ListTools() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	fmt.Println("Available tools:")
	for name, tool := range cfg.Tools {
		fmt.Printf("  - %-10s (%s)\n", name, tool.Image)
	}
}
