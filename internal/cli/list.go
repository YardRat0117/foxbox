package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured tools status",
	Run: func(_ *cobra.Command, _ []string) {
		// `panel` defined in `rootCmd`
		listConfig(panel)
	},
}

func listConfig(panel *container.Panel) {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config: ", err)
		os.Exit(1)
	}

	// Load installed tools
	installed, err := panel.CheckTools(cfg.Tools)
	if err != nil {
		fmt.Println("Error listing installed tools: ", err)
		os.Exit(1)
	}

	// Constants for formatting output
	const nameWidth, parenWidth, insWidth = 15, 25, 5

	fmt.Println("Configured tools:")
	for name, tool := range cfg.Tools {
		s, exists := installed[name]

		// Installation status
		statusWidth := len("[not installed]")
		status := "[not installed]"
		if exists && s.Installed {
			status = "[Installed]"
		}
		status = fmt.Sprintf("%-*s", statusWidth, status)

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
