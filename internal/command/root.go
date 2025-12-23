// Package command parses cli arguments
package command

import (
	"github.com/spf13/cobra"

	"github.com/YardRat0117/foxbox/internal/app"
	"github.com/YardRat0117/foxbox/internal/config"
	"github.com/YardRat0117/foxbox/internal/runtime/docker"
)

type rootContext struct {
	app *app.App
}

// NewRootCommand creates noew root command, including dependency injection
func NewRootCommand() *cobra.Command {
	// TODO
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// TODO - more runtime support
	rt := docker.New("unix:///var/run/docker.sock")

	application := app.New(cfg, rt)

	rootCtx := &rootContext{app: application}

	rootCmd := &cobra.Command{
		Use:   "foxbox",
		Short: "Foxbox - run tools in containers",
	}

	// TODO - more function support
	rootCmd.AddCommand(newRunCommand(rootCtx))

	return rootCmd
}
