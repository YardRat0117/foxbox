package container

// PodmanRuntime implements the Runtime interface using Podman.
type PodmanRuntime struct {
	CLIToolManager
}

// NewPodmanRuntime creates a new PodmanRuntime.
func NewPodmanRuntime() Runtime {
	imgMgr := &CLIImageManager{
		RuntimeName: "podman",
	}

	return &PodmanRuntime{
		CLIToolManager: CLIToolManager{
			RuntimeName: "podman",
			ImageMgr:    imgMgr,
		},
	}
}
