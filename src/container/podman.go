package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/YardRat0117/ratbox/src/config"
)

type PodmanRuntime struct{}

func NewRuntime() ContainerRuntime {
	return &PodmanRuntime{}
}

func (p *PodmanRuntime) PullImage(tool config.Tool) error {
	fmt.Printf("Pulling image %s using podman...\n", tool.Image)
	cmd := exec.Command("podman", "pull", tool.Image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// (true,nil)  -> Image exists
// (false,nil) -> Image doesn't exist
// (false,err) -> Error occurred
func (p * PodmanRuntime) ImageExists(image string) (bool, error) {
	if _, err := exec.LookPath("podman"); err != nil {
		return false, fmt.Errorf("podman not found in PATH: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "podman", "image", "inspect", image)

	// Keep stderr for debugging, discard stdout
	var stderr bytes.Buffer
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		return true, nil // Image exists
	} else {
		// Timeout
		if ctx.Err() == context.DeadlineExceeded {
			return false, fmt.Errorf("timed out checking image: %w", err)
		}

		var msg = strings.TrimSpace(stderr.String())
		// Occasional cases for container not exist
		if strings.Contains(msg, "No such object") ||
			strings.Contains(msg, "No such image") ||
			strings.Contains(strings.ToLower(msg), "not found") {
			return false, nil
		}

		// Unoccasional cases
		return false, fmt.Errorf("image inspect failed: %s", msg)
	}
}

func (p * PodmanRuntime) BuildRunCmd(tool config.Tool, ver string, args []string) *exec.Cmd {
	podmanArgs := []string{"run", "--rm", "-it"}
	for _, vol := range tool.Volumes {
		hostVol := vol

		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		hostVol = strings.ReplaceAll(hostVol, "$(pwd)", cwd)
		podmanArgs = append(podmanArgs, "-v", hostVol)
	}

	image := tool.Image
	if ver != "" {
		image = fmt.Sprintf("%s:%s", tool.Image, ver)
	}

	podmanArgs = append(podmanArgs, "-w", tool.Workdir, image, tool.Entry)
	podmanArgs = append(podmanArgs, args...)

	return exec.Command("podman", podmanArgs...)
}
