// Package main
package main

import (
	"fmt"
	"os"

	"github.com/YardRat0117/foxbox/internal/command"
)

func main() {
	rootCmd := command.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
