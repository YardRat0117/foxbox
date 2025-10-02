package container

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/YardRat0117/foxbox/internal/types"
)

var _ toolManager = (*cliToolManager)(nil)

type cliToolManager struct {
	runtimeName string
	imageMgr    imageManager
}

// installTool pulls the corresponding image for required tool.
func (t *cliToolManager) installTool(toolName, version string) error {
	image := fmt.Sprintf("%s:%s", toolName, version)
	fmt.Printf("Installing tool %s by pulling image %s...\n", toolName, image)
	return t.imageMgr.pullImage(image)
}

// removeTool removes the corresponding image for specified tool.
func (t *cliToolManager) removeTool(toolName string, imgName string, version string) error {
	image := fmt.Sprintf("%s:%s", imgName, version)

	if !Confirm(fmt.Sprintf("Sure to remove tool %s by removing image %s@%s?", toolName, imgName, version)) {
		fmt.Println("Skipped: ", image)
		return nil
	}
	fmt.Printf("Removing tool %s by removing image %s@%s...\n", toolName, image, version)
	return t.imageMgr.removeImage(image)
}

// runTool runs the given tool.
func (t *cliToolManager) runTool(tool types.Tool, version string, args []string) error {
	cmdArgs := []string{"run", "--rm", "-i"}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Can't get current working directory: %v", err)
	}
	for _, vol := range tool.Volumes {
		hostVol := strings.ReplaceAll(vol, "$(pwd)", cwd)
		cmdArgs = append(cmdArgs, "-v", hostVol)
	}

	image := tool.Image
	if version != "" {
		image = fmt.Sprintf("%s:%s", tool.Image, version)
	}

	cmdArgs = append(cmdArgs, "-w", tool.Workdir, image, tool.Entry)
	cmdArgs = append(cmdArgs, args...)

	// `<RuntimeName> run --rm -i -v <hostVol> <image>:<version> -w <workdir> <image> <entry>`
	// #nosec G204: parameters are split, and RUntimeName is controlled
	cmd := exec.Command(t.runtimeName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// checkTools inspects the given tools and returns their status.
func (t *cliToolManager) checkTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	localImages, err := t.imageMgr.getLocalImages()
	if err != nil {
		return nil, err
	}

	status := make(map[string]types.ToolStatus)
	for name, tool := range tools {
		repo, tag := SplitImage(tool.Image)
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

// cleanTools removes all installed images for configurated tools
func (t *cliToolManager) cleanTools(tools map[string]types.Tool) error {
	var errs []error

	statuses, err := t.checkTools(tools)
	if err != nil {
		return fmt.Errorf("check tools failed: %w", err)
	}

	for name, st := range statuses {
		tool, ok := tools[name]
		if !ok {
			continue
		}
		imgName, _ := SplitImage(tool.Image)

		for _, tag := range st.LocalTags {
			if err := t.removeTool(name, imgName, tag); err != nil {
				errs = append(errs, fmt.Errorf("failed to remove %s:%s: %w", name, tag, err))
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
