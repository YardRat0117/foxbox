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

// (true,nil)  -> Image exists
// (false,nil) -> Image doesn't exist
// (false,err) -> Error occurred
func (p *PodmanRuntime) ImageExists(image string) (bool, error) {
	if _, err := exec.LookPath("podman"); err != nil {
		return false, fmt.Errorf("podman not found in PATH: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "podman", "image", "inspect", image)

	// Discard stdout, capture stderr
	var stderr bytes.Buffer
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Image found
	if err == nil {
		return true, nil // Image exists
	}

	// Timeout
	if ctx.Err() == context.DeadlineExceeded {
		return false, fmt.Errorf("timed out checking image: %w", err)
	}

	// Not found
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		if code == 125 {
			// Not found, no error
			return false, nil
		}
	}

	// Unoccasional cases
	msg := strings.TrimSpace(stderr.String())
	return false, fmt.Errorf("image inspect failed: %s", msg)
}

func (p *PodmanRuntime) PullImage(tool config.Tool) error {
	fmt.Printf("Pulling image %s using podman...\n", tool.Image)
	cmd := exec.Command("podman", "pull", tool.Image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PodmanRuntime) BuildRunCmd(tool config.Tool, ver string, args []string) *exec.Cmd {
	podmanArgs := []string{"run", "--rm", "-i"}
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

// helper: split "repo:tag" -> repo, tag
func splitImage(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
	}
	return image, "latest"
}

// helper
func contains(slice []string, s string) bool {
    for _, v := range slice {
        if v == s {
            return true
        }
    }
    return false
}

func (p *PodmanRuntime) ListInstalled(tools map[string]config.Tool) (map[string]ToolStatus, error) {
	status := make(map[string]ToolStatus)

	// Retrieve all images and tags
	cmd := exec.Command("podman", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	localImages := make(map[string][]string)
	for _, line := range lines {
		repo, tag := splitImage(line)
		localImages[repo] = append(localImages[repo], tag)
	}

	for name, tool := range tools {
		installed := false
		var tags []string
		repo, tag := splitImage(tool.Image)
		for localRepo, localTags := range localImages {
			if strings.HasSuffix(localRepo, repo) {
				tags = localTags
				if tag == "" || tag == "latest" || contains(localTags, tag) {
					installed = true
				}
				break
			}
		}
		status[name] = ToolStatus{
			Installed: installed,
			LocalTags: tags,
		}
	}

	return status, nil
}
