package container

// PodmanRuntime implements the Runtime interface using Podman.
type PodmanRuntime struct{}

// NewPodmanRuntime creates a new PodmanRuntime.
func NewPodmanRuntime() Runtime {
	return &PodmanRuntime{}
}
