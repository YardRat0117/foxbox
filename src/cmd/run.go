package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/ratbox/src/config"
	"github.com/YardRat0117/ratbox/src/container"
)

var runCmd = &cobra.Command{
	Use: "run <tool> [args..] -- [toolArgs...]",
	Short: "Run a tool inside its container",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]
		toolArgs := args[1:]
		RunTool(toolName, toolArgs)
	},
}

func RunTool(toolName string, toolArgs []string) {
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

	cmd := runtime.BuildRunCmd(tool, toolArgs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running container:", err)
		os.Exit(1)
	}
}
