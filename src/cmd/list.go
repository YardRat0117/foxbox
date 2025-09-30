package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/config"
	"github.com/YardRat0117/foxbox/src/container"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured tools status",
	Run: func(cmd *cobra.Command, args []string) {
		// `runtime` defined in `rootCmd`
		listConfig(runtime)
	},
}

func listConfig(runtime container.Runtime) {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	// Load installed tools
	installed, err := runtime.CheckTools(cfg.Tools)
	if err != nil {
		fmt.Println("Error listing installed tools:", err)
		os.Exit(1)
	}

	// Constants for formatting output
	const nameWidth, parenWidth, insWidth = 15, 15, 5

	fmt.Println("Configured tools:")
	for name, tool := range cfg.Tools {
		s, exists := installed[name]

		// Installation status
		status := "[not installed]"
		if exists && s.Installed {
			status = "[Installed]" + strings.Repeat(" ", insWidth)
		}

		// Tags
		tags := ""
		if exists && len(s.LocalTags) > 0 {
			tags = "tags " + strings.Join(s.LocalTags, ", ")
		}

		// Output
		fmt.Printf("- %-*s %-*s %s %s\n",
			nameWidth, name,
			parenWidth, fmt.Sprintf("(%s)", tool.Image),
			status, tags)
	}
}
