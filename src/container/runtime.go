package container

import (
	"os/exec"

	"github.com/YardRat0117/ratbox/src/config"
)

type ContainerRuntime interface {
	ImageExists(image string) (bool, error)

	PullImage(tool config.Tool) error

	BuildRunCmd(tool config.Tool, ver string, args []string) *exec.Cmd
}
