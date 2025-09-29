package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

	// Prefer `~/.config/rbox.yml`
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

// (true,nil) -> Image exists
// (false,nil) -> Image doesn't exist, no error
// (false,err) -> Error occurred
func imageExists(image string) (bool, error) {
	if _, err := exec.LookPath("podman"); err != nil {
		return false, fmt.Errorf("podman not found in PATH: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "podman", "image", "inspect", image)

	// Keep stderr for debugging, discard stdout
	var stderr bytes.Buffer
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		return true, nil // Image exists
	} else {
		// Timeout
		if ctx.Err() == context.DeadlineExceeded {
			return false, fmt.Errorf("timed out checking image: %w", err)
		}

		var msg = strings.TrimSpace(stderr.String())
		// Occasional cases for container not exist
		if strings.Contains(msg, "No such object") ||
			strings.Contains(msg, "No such image") ||
			strings.Contains(strings.ToLower(msg), "not found") {
			return false, nil
		}

		// Unoccasional cases
		return false, fmt.Errorf("image inspect failed: %s", msg)
	}
}

func pullImage(image string) error {
	fmt.Printf("Pulling image %s using podman...\n", image)
	cmd := exec.Command("podman", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rbox <tool> [args...]",
		Short: "Run development tools in containers",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			toolName := args[0]

			if toolName == "list" {
				fmt.Println("Use `rbox list` instead of `rbox list <args>`")
				os.Exit(0)
			}

			toolArgs := args[1:]

			cfg, err := loadConfig()
			if err != nil {
				fmt.Println("Failed to load config:", err)
				os.Exit(1)
			}

			tool, ok := cfg.Tools[toolName]
			if !ok {
				fmt.Printf("Tool '%s' not found in config\n", toolName)
				fmt.Println("Available tools:")
				for name := range cfg.Tools {
					fmt.Printf("  - %s\n", name)
				}
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

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available tools",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := loadConfig()
			if err != nil {
				fmt.Println("Failed to load config:", err)
				os.Exit(1)
			}

			fmt.Println("Available tools:")
			for name, tool := range cfg.Tools {
				fmt.Printf("  - %-10s (%s)\n", name, tool.Image)
			}
		},
	}

	var installCmd = &cobra.Command{
		Use:   "install <tool>",
		Short: "Install (pull) a tool's container image",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			toolName := args[0]

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

			if err := pullImage(tool.Image); err != nil {
				fmt.Println("Error pulling image:", err)
				os.Exit(1)
			}

			fmt.Printf("Image %s installed successfully!\n", tool.Image)
		},
	}

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
