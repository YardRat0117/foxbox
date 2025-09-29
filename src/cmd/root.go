package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "foxbox <tool> [args...] -- [toolArgs...]",
	Short: "Foxbox - lightweight tool runtime",
	Long:  "Foxbox manages containerized developer tools with a simple interface.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
}
