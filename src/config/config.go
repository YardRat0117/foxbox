package config

import (
	"os"
	"path/filepath"

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

func LoadConfig() (*Config, error) {
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
