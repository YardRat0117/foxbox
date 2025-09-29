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
		ListConfig()
	},
}

func ListConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	runtime := container.NewRuntime()
	installed, err := runtime.ListInstalled(cfg.Tools)
	if err != nil {
		fmt.Println("Error listing installed tools:", err)
		os.Exit(1)
	}

	const parenWidth = 15
	const installPad = 5
	const tagsWidth = 30

	fmt.Println("Configured tools:")
	for name, tool := range cfg.Tools {
		status := "[not installed]"
		tags := ""
		if s, exists := installed[name]; exists {
			if s.Installed {
				status = "[installed]"
			}
			if len(s.LocalTags) > 0 {
				tags = "tags: " + strings.Join(s.LocalTags, ", ")
			}
		}
		paren := fmt.Sprintf("(%s)", tool.Image)

		if status == "[installed]" {
			status = status + strings.Repeat(" ", installPad)
		}

		fmt.Printf("  - %-8s %-*s %s %s\n", name, parenWidth, paren, status, tags)

	}
}
