package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/container"
)

var panelOpt string
var panel *container.Panel

var rootCmd = &cobra.Command{
	Use:   "foxbox <cmd> -- [toolArgs...]",
	Short: "Foxbox - lightweight tool panel",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",

	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		switch panelOpt {
		case "podman":
			panel = container.NewPodmanPanel()

		case "docker":
			panel = container.NewDockerPanel()

		case "podman-api":
			fmt.Fprintln(os.Stderr, "[foxbox] panel=podman-api not implemented yet, falling back to podman")
			panel = container.NewPodmanPanel()

		case "docker-api":
			panel = container.NewDockerAPIPanel()

		case "podman-native":
			fmt.Fprintln(os.Stderr, "[foxbox] panel=podman-native not implemented yet, falling back to podman")
			panel = container.NewPodmanPanel()

		default:
			fmt.Fprintf(os.Stderr, "[foxbox] unknown panel=%s, falling back to podman\n", panelOpt)
			panel = container.NewPodmanPanel()
		}

	},
}

// Execute provides external interface for `main.go` to launch the program
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&panelOpt,
		"panel",
		"podman",
		"container panel backend (podman | docker | api | podapi)",
	)

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}
