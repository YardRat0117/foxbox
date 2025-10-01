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

// --- imageManager ---

// checkImage checks if the given image exists locally using Podman.
func (p *PodmanRuntime) checkImage(image string) (bool, error) {
	// Check Podman
	if _, err := exec.LookPath("podman"); err != nil {
		return false, fmt.Errorf("podman not found in PATH: %w", err)
	}

	// Quit if timed out
	const delay = 30 // in seconds
	ctx, cancel := context.WithTimeout(context.Background(), delay*time.Second)
	defer cancel()

	// Build command `podman image inspect`
	cmd := exec.CommandContext(ctx, "podman", "image", "inspect", image)
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
	const NF = 125 // Podman return value for not found image
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == NF {
		return false, nil
	}

	return false, fmt.Errorf("image inspect failed: %s", strings.TrimSpace(stderr.String()))
}

// pullImage pulls the specified image using Podman.
func (p *PodmanRuntime) pullImage(image string) error {
	// Build command `podman pull <image>`
	cmd := exec.Command("podman", "pull", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}

// removeImage removes the specified image using Podman.
func (p *PodmanRuntime) removeImage(image string) error {
	// Hint
	fmt.Printf("Removing image %s using podman...\n", image)

	// Build command `podman rmi -f <image>`
	cmd := exec.Command("podman", "rmi", "-f", image)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return error info
	return cmd.Run()
}
