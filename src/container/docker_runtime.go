package container

// DockerRuntime implements the Runtime interface using Docker.
type DockerRuntime struct{}

// NewDockerRuntime creates a new DockerRuntime.
func NewDockerRuntime() Runtime {
	return &DockerRuntime{}
}
