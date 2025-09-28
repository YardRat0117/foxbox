package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Tool struct {
	Image   string   `yaml:"image"`
	Entry   string   `yaml:"entry"`
	Workdir string   `yaml:"workdir"`
	Volumes []string `yaml:"volumes"`
}

type Config struct {
	Tools map[string]Tool `yaml:"tools"`
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	userConfigPath := filepath.Join(home, ".config", "rbox.yml")
	configPath := filepath.Join("config", "default.yml")

	// Overwrite `config/default.yml` with `~/.config/rbox.yml`
	if _, err := os.Stat(userConfigPath); err == nil {
		configPath = userConfigPath
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rbox <tool> [args...]",
		Short: "Run development tools in containers",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			toolName := args[0]
			toolArgs := args[1:]

			cfg, err := loadConfig()
			if err != nil {
				fmt.Println("Failed to load config:", err)
				os.Exit(1)
			}

			tool, ok := cfg.Tools[toolName]
			if !ok {
				fmt.Printf("Tool '%s' not found in config\n", toolName)
				os.Exit(1)
			}

			// Build podman command
			podmanArgs := []string{"run", "--rm", "-it"}
			for _, vol := range tool.Volumes {
				hostVol := vol

				cwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}

				hostVol = strings.ReplaceAll(hostVol, "$(pwd)", cwd)
				podmanArgs = append(podmanArgs, "-v", hostVol)
			}
			podmanArgs = append(podmanArgs, "-w", tool.Workdir, tool.Image, tool.Entry)
			podmanArgs = append(podmanArgs, toolArgs...)

			// Execute podman
			execCmd := exec.Command("podman", podmanArgs...)
			execCmd.Stdin = os.Stdin
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr

			if err := execCmd.Run(); err != nil {
				fmt.Println("Error running container:", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
