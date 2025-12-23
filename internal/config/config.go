package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/YardRat0117/foxbox/internal/types"
)

// Load loads application configuration from predefined locations.
func Load() (*types.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Note: user config file (located at `~/.config/foxbox.yml`) would be preffered
	userConfigPath := filepath.Join(home, ".config", "foxbox.yml")
	configPath := filepath.Join("config", "default.yml")
	if _, err := os.Stat(userConfigPath); err == nil {
		configPath = userConfigPath
	}

	// Open the config file
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decode the config file, and return as the type `Config`
	var cfg types.Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
