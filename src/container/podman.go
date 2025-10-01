package container

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/YardRat0117/foxbox/src/types"
)

// PodmanRuntime implements the Runtime interface using Podman.
type PodmanRuntime struct{}

// NewPodmanRuntime creates a new PodmanRuntime.
func NewPodmanRuntime() Runtime {
	return &PodmanRuntime{}
}

// --- ImageManager ---

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

// --- ToolManager ---

// InstallTool pulls the corresponding image for required tool.
func (p *PodmanRuntime) InstallTool(toolName string, version string) error {
	image := fmt.Sprintf("%s:%s", toolName, version)

	// Hint
	fmt.Printf("Installing tool %s by pulling image %s...\n", toolName, image)

	return p.pullImage(image)
}

// RemoveTool removes the corresponding image for specified tool.
func (p *PodmanRuntime) RemoveTool(toolName string, imgName string, version string) error {
	image := fmt.Sprintf("%s:%s", imgName, version)

	fmt.Printf("Sure to remove tool %s by removing image %s? [y/N]", toolName, image)

	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}

	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Skipped: ", image)
		return nil
	}

	// Hint
	fmt.Printf("Removing tool %s by removing image %s...\n", toolName, image)

	return p.removeImage(image)
}

// RunTool runs the given tool inside a Podman container.
func (p *PodmanRuntime) RunTool(tool types.Tool, version string, args []string) error {
	// Get current working directory (cwd)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Can't get current working directory: %v", err)
	}

	// Build command `run --rm -i`
	podmanArgs := []string{"run", "--rm", "-i"}

	// Attached `-v <hostVol>` for directory mounting
	for _, vol := range tool.Volumes {
		hostVol := strings.ReplaceAll(vol, "$(pwd)", cwd)
		podmanArgs = append(podmanArgs, "-v", hostVol)
	}

	// Attached `<image>:<version>` for image and tag, `latest` by default
	image := tool.Image
	if version != "" {
		image = fmt.Sprintf("%s:%s", tool.Image, version)
	}

	// Attached `-w <workdir> <image> <entry>` for working directory etc.
	podmanArgs = append(podmanArgs, "-w", tool.Workdir, image, tool.Entry)

	// Attached other arguments
	podmanArgs = append(podmanArgs, args...)

	// Build the cmd
	// `podman run --rm -i -v <hostVol> <image>:<version> -w <workdir> <image> <entry>`
	// --rm : deprecate the container after termination
	// -i : interactive, which enables stdin
	// -v <hostVol> : mounting host directory for the container
	// <image>:<version> : specified image and tag
	// -w <workdir> : working directory in the container
	// <image> <entry> : run the entry command
	cmd := exec.Command("podman", podmanArgs...)
	// Stdout & Stderr redirected to OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Execute the command and return error info
	return cmd.Run()
}

// CheckTools inspects the given tools and returns their status using Podman.
func (p *PodmanRuntime) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	// Declare `ToolStatus` for return value
	status := make(map[string]types.ToolStatus)

	// Build command `podman images --format {{.Repository}}:{{.Tag}}`
	// This retrieves all local images in the format "repo:tag"
	cmd := exec.Command("podman", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output into localImages[repo][tag]
	localImages := map[string]map[string]struct{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		repo, tag := splitImage(line)
		if _, ok := localImages[repo]; !ok {
			localImages[repo] = map[string]struct{}{}
		}
		localImages[repo][tag] = struct{}{}
	}

	// Compare each tool with local images
	for name, tool := range tools {
		repo, tag := splitImage(tool.Image)
		installed := false
		localTags := []string{}

		// Find whether repo exists locally
		for localRepo, tags := range localImages {
			if strings.HasSuffix(localRepo, repo) {
				// Collect available tags
				for t := range tags {
					localTags = append(localTags, t)
				}

				// If tag matches (or defaults to "latest"), mark as installed
				if _, ok := tags[tag]; ok {
					installed = true
				} else if tag == "" || tag == "latest" {
					installed = true
				}
				break
			}
		}

		// Record result
		status[name] = types.ToolStatus{
			Installed: installed,
			LocalTags: localTags,
		}
	}

	return status, nil
}

// CleanTools removes all installed images for configurated tools
func (p *PodmanRuntime) CleanTools(tools map[string]types.Tool) error {
	var errs []error

	statuses, err := p.CheckTools(tools)
	if err != nil {
		return fmt.Errorf("check tools failed: %w", err)
	}

	for name, st := range statuses {
		tool, ok := tools[name]
		if !ok {
			continue
		}
		imgName, _ := splitImage(tool.Image)

		for _, tag := range st.LocalTags {
			if err := p.RemoveTool(name, imgName, tag); err != nil {
				errs = append(errs, fmt.Errorf("failed to remove %s:%s: %w", name, tag, err))
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// --- helpers ---
func splitImage(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
	}
	return image, "latest"
}
