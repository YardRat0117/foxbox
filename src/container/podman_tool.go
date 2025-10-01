package container

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/YardRat0117/foxbox/src/types"
)

// --- ToolManager ---

// InstallTool pulls the corresponding image for required tool.
func (p *PodmanRuntime) InstallTool(toolName string, version string) error {
	image := fmt.Sprintf("%s:%s", toolName, version)
	fmt.Printf("Installing tool %s by pulling image %s...\n", toolName, image)
	return p.pullImage(image)
}

// RemoveTool removes the corresponding image for specified tool.
func (p *PodmanRuntime) RemoveTool(toolName string, imgName string, version string) error {
	image := fmt.Sprintf("%s:%s", imgName, version)

	if !confirm(fmt.Sprintf("Sure to remove tool %s by removing image %s@%s?", toolName, imgName, version)) {
		fmt.Println("Skipped: ", image)
		return nil
	}
	fmt.Printf("Removing tool %s by removing image %s@%s...\n", toolName, image, version)
	return p.removeImage(image)
}

// RunTool runs the given tool inside a Podman container.
func (p *PodmanRuntime) RunTool(tool types.Tool, version string, args []string) error {
	podmanArgs := []string{"run", "--rm", "-i"}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Can't get current working directory: %v", err)
	}
	for _, vol := range tool.Volumes {
		hostVol := strings.ReplaceAll(vol, "$(pwd)", cwd)
		podmanArgs = append(podmanArgs, "-v", hostVol)
	}

	image := tool.Image
	if version != "" {
		image = fmt.Sprintf("%s:%s", tool.Image, version)
	}

	podmanArgs = append(podmanArgs, "-w", tool.Workdir, image, tool.Entry)
	podmanArgs = append(podmanArgs, args...)

	// `podman run --rm -i -v <hostVol> <image>:<version> -w <workdir> <image> <entry>`
	cmd := exec.Command("podman", podmanArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CheckTools inspects the given tools and returns their status using Podman.
func (p *PodmanRuntime) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	localImages, err := p.getLocalImages()
	if err != nil {
		return nil, err
	}

	status := make(map[string]types.ToolStatus)
	for name, tool := range tools {
		repo, tag := splitImage(tool.Image)
		installed := false
		localTags := []string{}

		for localRepo, tags := range localImages {
			if path.Base(localRepo) == repo {
				for t := range tags {
					localTags = append(localTags, t)
				}

				if _, ok := tags[tag]; ok || tag == "" || tag == "latest" {
					installed = true
				}
				break
			}
		}

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
