package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
	"github.com/YardRat0117/foxbox/internal/types"
)

var (
	cfg      *types.Config
	panelOpt string
	panel    container.PanelInterface
)

var rootCmd = &cobra.Command{
	Use:   "foxbox <cmd> -- [toolArgs...]",
	Short: "Foxbox - lightweight tool panel",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",

	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		// Load config
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			fatal("Failed to load config:", err)
		}

		// Init panel
		switch panelOpt {
		case "podman":
			panel = container.NewPodmanPanel()
		case "docker":
			panel = container.NewDockerPanel()
		case "docker-api":
			panel = container.NewDockerAPIPanel()
		case "libpod-api":
			fatal("Panel `libpod-api` under dev", nil)
		default:
			fatal(fmt.Sprintf("Unknown panel %s", panelOpt), nil)
		}

	},
}

// Execute provides external interface for `main.go` to launch the program
func Execute() error {
	rootCmd.PersistentFlags().StringVar(&panelOpt, "panel", "podman", "container panel backend (podman | docker | api)")
	rootCmd.AddCommand(installCmd, removeCmd, runCmd, listCmd, cleanCmd, versionCmd)
	return rootCmd.Execute()
}
