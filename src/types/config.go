package types

// Config is a map of `Tool`, which corresponds to a config file
type Config struct {
	Tools map[string]Tool `yaml:"tools"`
}
