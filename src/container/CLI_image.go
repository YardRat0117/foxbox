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
)

var _ imageManager = (*CLIImageManager)(nil)

type CLIImageManager struct {
	RuntimeName string
}

// checkImage checks if the given image exists locally.
func (i *CLIImageManager) checkImage(image string) (bool, error) {
	if _, err := exec.LookPath(i.RuntimeName); err != nil {
		return false, fmt.Errorf("%s not found in PATH: %w", i.RuntimeName, err)
	}

	// Quit if timed out
	const delay = 30 // in seconds
	ctx, cancel := context.WithTimeout(context.Background(), delay*time.Second)
	defer cancel()

	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return false, err
	}

	// Build command `image inspect`
	// #nosec G204: image name is validated by verifyImage
	cmd := exec.CommandContext(ctx, i.RuntimeName, "image", "inspect", image)
	// Discard stdout, reserve stderr
	cmd.Stdout = io.Discard
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	// Found
	if err == nil {
		return true, nil
	}
	// Timed out
	if ctx.Err() == context.DeadlineExceeded {
		return false, fmt.Errorf("timed out checking image: %w", err)
	}
	// Not found
	const NFPodman, NFDocker = 125, 1
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		if code == NFPodman || code == NFDocker {
			return false, nil
		}
		return false, fmt.Errorf("image inspect failed: %s", strings.TrimSpace(stderr.String()))
	}

	return false, fmt.Errorf("image inspect failed: %s", strings.TrimSpace(stderr.String()))
}

// pullImage pulls the specified image using Podman.
func (i *CLIImageManager) pullImage(image string) error {
	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return err
	}

	// Build command `podman pull <image>`
	// #nosec G204: image name is validated by verifyImage
	cmd := exec.Command(i.RuntimeName, "pull", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}

// removeImage removes the specified image using Podman.
func (i *CLIImageManager) removeImage(image string) error {
	// Hint
	fmt.Printf("Removing image %s using %s...\n", image, i.RuntimeName)

	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return err
	}

	// Build command `podman rmi -f <image>`
	// #nosec G204: image name is validated by verifyImage
	cmd := exec.Command(i.RuntimeName, "rmi", "-f", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}

// GetLocalImages retrieves pulled images
func (i *CLIImageManager) getLocalImages() (map[string]map[string]struct{}, error) {
	// #nosec G204: fixed args, safe to execute
	cmd := exec.Command(i.RuntimeName, "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	localImages := map[string]map[string]struct{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		repo, tag := SplitImage(line)
		if _, ok := localImages[repo]; !ok {
			localImages[repo] = map[string]struct{}{}
		}
		localImages[repo][tag] = struct{}{}
	}
	return localImages, nil
}
