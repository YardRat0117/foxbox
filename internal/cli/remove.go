package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
)

var removeCmd = &cobra.Command{
	Use:   "remove <tool>",
	Short: "Remove a tool's container image",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		// Split original parameters
		toolInfo := strings.SplitN(args[0], "@", 2)

		// Parse tool info
		toolName := toolInfo[0]
		toolVer := "latest" // `latest` by default
		if (len(toolInfo)) == 2 {
			toolVer = toolInfo[1]
		}

		// `runtime` defined in `rootCmd`
		removeTool(runtime, toolName, toolVer)
	},
}

func removeTool(runtime container.Runtime, toolName string, toolVer string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	tool, ok := cfg.Tools[toolName]
	if !ok {
		fmt.Printf("Tool '%s' not found in config\n", toolName)
		os.Exit(1)
	}

	if err := runtime.RemoveTool(toolName, tool.Image, toolVer); err != nil {
		fmt.Printf("Error removing tool: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Image %s removed successfully!\n", tool.Image)
}
