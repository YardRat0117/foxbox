package types

// ToolStatus contains detailed info (installation & tags) read from the container
type ToolStatus struct {
        Installed bool
        LocalTags []string
}

// Tool contains basic info for something to run in the container
type Tool struct {
	Image   string   `yaml:"image"`
	Entry   string   `yaml:"entry"`
	Workdir string   `yaml:"workdir"`
	Volumes []string `yaml:"volumes"`
}
