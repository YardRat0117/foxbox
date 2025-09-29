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
		toolInfo := strings.SplitN(args[0], "@", 2)
		toolName := toolInfo[0]
		version := ""
		if len(toolInfo) == 2 {
			version = toolInfo[1]
		}
		toolArgs := args[1:]
		RunTool(toolName, version, toolArgs)
	},
}

func RunTool(toolName string, version string, toolArgs []string) {
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

	runtime := container.NewRuntime()

	cmd := runtime.BuildRunCmd(tool, version, toolArgs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running container:", err)
		os.Exit(1)
	}
}
