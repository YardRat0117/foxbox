package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
)

var installCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install (pull) a tool's container image",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		// Split original parameters
		toolInfo := strings.SplitN(args[0], "@", 2)

		// Parse tool info
		toolName := toolInfo[0]
		toolVer := "latest" // `latest` by dafault
		if len(toolInfo) == 2 {
			toolVer = toolInfo[1]
		}

		// `runtime` defined in `rootCmd`
		installTool(runtime, toolName, toolVer)
	},
}

func installTool(runtime container.Runtime, toolName string, toolVer string) {
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

	if err := runtime.InstallTool(tool.Image, toolVer); err != nil {
		fmt.Println("Error installing tool:", err)
		os.Exit(1)
	}

	fmt.Printf("Image %s installed successfully!\n", tool.Image)
}
