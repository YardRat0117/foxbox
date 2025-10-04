package container

func NewDockerPanel() *Panel {
	rtMgr := &cliRuntimeManager{
		runtimeName: "docker",
	}
	return &Panel{
		toolManager: &toolManager{
			rt: rtMgr,
		},
	}
}

func NewPodmanPanel() *Panel {
	rtMgr := &cliRuntimeManager{
		runtimeName: "podman",
	}
	return &Panel{
		toolManager: &toolManager{
			rt: rtMgr,
		},
	}
}

func NewDockerAPIPanel() *Panel {
	rtMgr := &dockerAPIManager{
		runtimeURL: "unix:///var/run/docker.sock",
	}
	return &Panel {
		toolManager: &toolManager{
			rt: rtMgr,
		},
	}

}
