package container

// DockerRuntime implements the Runtime interface using Docker.
type DockerRuntime struct {
	CLIToolManager
}

// NewDockerRuntime creates a new DockerRuntime.
func NewDockerRuntime() Runtime {
	imgMgr := &CLIImageManager{
		RuntimeName: "docker",
	}

	return &DockerRuntime{
		CLIToolManager: CLIToolManager{
			RuntimeName: "docker",
			ImageMgr:    imgMgr,
		},
	}
}
