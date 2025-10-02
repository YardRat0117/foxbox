package main

import (
	"fmt"
	"os"

	"github.com/YardRat0117/foxbox/internal/cli"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
