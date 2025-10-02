package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show foxbox version",
	Run: func(_ *cobra.Command, _ []string) {
		showVersion()
	},
}

func showVersion() {
	// `Commit` is declared in package `main`
	fmt.Printf("FoxBox Commit: %s\n", version.Commit)
}
