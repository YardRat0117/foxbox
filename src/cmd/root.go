package cmd

import (
	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/container"
)

var runtime container.Runtime

var rootCmd = &cobra.Command{
	Use:   "foxbox <cmd> -- [toolArgs...]",
	Short: "Foxbox - lightweight tool runtime",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Podman currently
		runtime = container.NewPodmanRuntime()
	},
}

// Execute provides external interface for `main.go` to launch the program
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}
