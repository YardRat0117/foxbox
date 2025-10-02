package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/container"
)

var runtimeOpt string
var runtime container.Runtime

var rootCmd = &cobra.Command{
	Use:   "foxbox <cmd> -- [toolArgs...]",
	Short: "Foxbox - lightweight tool runtime",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",

	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		switch runtimeOpt {
		case "podman":
			runtime = container.NewPodmanRuntime()

		case "docker":
			runtime = container.NewDockerRuntime()

		case "podman-api":
			fmt.Fprintln(os.Stderr, "[foxbox] runtime=podman-api not implemented yet, falling back to podman")
			runtime = container.NewPodmanRuntime()

		case "docker-api":
			fmt.Fprintln(os.Stderr, "[foxbox] runtime=docker-api not implemented yet, falling back to podman")

		case "podman-native":
			fmt.Fprintln(os.Stderr, "[foxbox] runtime=podman-native not implemented yet, falling back to podman")
			runtime = container.NewPodmanRuntime()

		default:
			fmt.Fprintf(os.Stderr, "[foxbox] unknown runtime=%s, falling back to podman\n", runtimeOpt)
			runtime = container.NewPodmanRuntime()
		}

	},
}

// Execute provides external interface for `main.go` to launch the program
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&runtimeOpt,
		"runtime",
		"podman",
		"container runtime backend (podman | docker | api | podapi)",
	)

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}
