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

	"github.com/YardRat0117/foxbox/internal/types"
)

var _ runtimeManager = (*cliRuntimeManager)(nil)

type cliRuntimeManager struct {
	runtimeName string
}

// checkImage checks if the given image exists locally.
func (r *cliRuntimeManager) checkImage(image string) (bool, error) {
	if _, err := exec.LookPath(r.runtimeName); err != nil {
		return false, fmt.Errorf("%s not found in PATH: %w", r.runtimeName, err)
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
	cmd := exec.CommandContext(ctx, r.runtimeName, "image", "inspect", image)
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

// pullImage pulls the specified image.
func (r *cliRuntimeManager) pullImage(image string) error {
	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return err
	}

	// Build command `podman pull <image>`
	// #nosec G204: image name is validated by verifyImage
	cmd := exec.Command(r.runtimeName, "pull", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}

// removeImage removes the specified image.
func (r *cliRuntimeManager) removeImage(image string) error {
	// Hint
	fmt.Printf("Removing image %s using %s...\n", image, r.runtimeName)

	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return err
	}

	// Build command `podman rmi -f <image>`
	// #nosec G204: image name is validated by verifyImage
	cmd := exec.Command(r.runtimeName, "rmi", "-f", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}

// localImages retrieves pulled images
func (r *cliRuntimeManager) localImages() (map[string]*types.ToolStatus, error) {
	// #nosec G204: fixed args, safe to execute
	cmd := exec.Command(r.runtimeName, "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	localImages := map[string]*types.ToolStatus{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		repo, tag := SplitImage(line)

		if _, ok := localImages[repo]; !ok {
			localImages[repo] = &types.ToolStatus{
				Installed: true,
				LocalTags: []string{},
			}
		}

		localImages[repo].LocalTags = append(localImages[repo].LocalTags, tag)
	}
	return localImages, nil
}

// runImage runs a container from the image
func (r *cliRuntimeManager) runImage(
	image string,
	entry string,
	workdir string,
	volumes []string,
	args []string,
) error {
	// Build command arguments
	cmdArgs := []string{"run", "--rm", "-i"}

	// Mount volumes
	for _, vol := range volumes {
		cmdArgs = append(cmdArgs, "-v", vol)
	}

	// Set work directory
	cmdArgs = append(cmdArgs, "-w", workdir)

	// Image and entry command (with verification)
	if err := verifyImage(image); err != nil {
		return err
	}
	if err := verifyEntry(entry); err != nil {
		return err
	}
	cmdArgs = append(cmdArgs, image, entry)

	// Other arguments
	cmdArgs = append(cmdArgs, args...)

	// Build command
	// #nosec G204: r.runtimeName is fixed and cmdArgs only contains container image, entry, and tool flags
	cmd := exec.Command(r.runtimeName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running container: %s %s\n", r.runtimeName, strings.Join(cmdArgs, " "))

	return cmd.Run()
}
