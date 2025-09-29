package cmd

import (

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ratbox <tool> [args...] -- [toolArgs...]",
	Short: "Ratbox - lightweight tool runtime",
	Long: "Ratbox manages containerized developer tools with a simple interface.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
}
