// package command parses cli arguments
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/container"
	"github.com/YardRat0117/foxbox/internal/types"
)

type rootContext struct {
	panel container.PanelInterface
	cfg   *types.Config
}

var (
	panelOpt string
	rootCtx  rootContext
)

var rootCmd = &cobra.Command{
	Use:   "foxbox <cmd> -- [toolArgs...]",
	Short: "Foxbox - lightweight tool panel",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",

	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Load config
		var errLoadCfg error
		rootCtx.cfg, errLoadCfg = config.LoadConfig()
		if errLoadCfg != nil {
			return fmt.Errorf("Failed to load config: %w", errLoadCfg)
		}

		// Init panel
		switch panelOpt {
		case "podman":
			rootCtx.panel = container.NewPodmanPanel()
		case "docker":
			rootCtx.panel = container.NewDockerPanel()
		case "docker-api":
			rootCtx.panel = container.NewDockerAPIPanel()
		case "libpod-api":
			return fmt.Errorf("Panel `libpod-api` under dev")
		default:
			return fmt.Errorf("Unknown panel %s", panelOpt)
		}

		return nil
	},
}

// Execute provides external interface for `main.go` to launch the program
func Execute() error {
	rootCmd.PersistentFlags().StringVar(&panelOpt, "panel", "podman", "container panel backend (podman | docker | api)")
	rootCmd.AddCommand(
		newInstallCommand(&rootCtx),
		newListCommand(&rootCtx),
		newRunCommand(&rootCtx),
		newRemoveCommand(&rootCtx),
		newCleanCommand(&rootCtx),
		newVersionCommand(),
	)
	return rootCmd.Execute()
}
