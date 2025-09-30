package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/config"
	"github.com/YardRat0117/foxbox/src/container"
)

var runCmd = &cobra.Command{
	Use:   "run <tool> [args..] -- [toolArgs...]",
	Short: "Run a tool inside its container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Split original parameters
		toolInfo := strings.SplitN(args[0], "@", 2)

		// Parse tool info
		toolName := toolInfo[0]
		toolVer := ""
		if len(toolInfo) == 2 {
			toolVer = toolInfo[1]
		}
		toolArgs := args[1:]

		// `runtime` defined in `rootCmd`
		runTool(runtime, toolName, toolVer, toolArgs)
	},
}

func runTool(runtime container.Runtime, toolName string, toolVer string, toolArgs []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	tool, ok := cfg.Tools[toolName]
	if !ok {
		fmt.Printf("Tool `%s` not found\n", toolName)
		os.Exit(1)
	}

	if err := runtime.RunTool(tool, toolVer, toolArgs); err != nil {
		fmt.Println("Error running tool:", err)
		os.Exit(1)
	}
}
