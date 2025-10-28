package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/version"
)

// func newVersionCommand creates a `cobra.Command` object `version` with given context
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show foxbox version",
		Run: func(_ *cobra.Command, _ []string) {
			// Both configured in package `version`
			fmt.Printf("FoxBox Commit: %s\nFoxbox Tag: %s\n", version.Commit, version.Tag)
		},
	}
}
