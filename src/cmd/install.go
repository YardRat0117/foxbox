package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/config"
	"github.com/YardRat0117/foxbox/src/container"
)

var installCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install (pull) a tool's container image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// `runtime` defined in `rootCmd`
		installTool(runtime, args[0])
	},
}

func installTool(runtime container.Runtime, toolName string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	tool, ok := cfg.Tools[toolName]
	if !ok {
		fmt.Printf("Tool '%s' not found in config\n", toolName)
		os.Exit(1)
	}

	if err := runtime.PullImage(tool.Image); err != nil {
		fmt.Println("Error pulling image:", err)
		os.Exit(1)
	}

	fmt.Printf("Image %s installed successfully!\n", tool.Image)
}
