package cmd

import (
	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/src/container"
)

var runtime container.ContainerRuntime

var rootCmd = &cobra.Command{
	Use:   "foxbox <tool> [args...] -- [toolArgs...]",
	Short: "Foxbox - lightweight tool runtime",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		runtime = container.NewRuntime()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
}
