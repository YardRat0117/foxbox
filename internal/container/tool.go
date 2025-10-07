package container

import (
	"slices"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/YardRat0117/foxbox/internal/types"
)

// toolManager manages the service logic, and calls ImageManagers to execute on actuall CRUD stuff
type toolManager struct {
	rt runtimeManager
}

// installTool pulls the corresponding image for required tool.
func (t *toolManager) installTool(toolName, version string) error {
	image := fmt.Sprintf("%s:%s", toolName, version)
	fmt.Printf("Installing tool %s by pulling image %s...\n", toolName, image)
	return t.rt.pullImage(image)
}

// removeTool removes the corresponding image for specified tool.
func (t *toolManager) removeTool(toolName string, imgName string, version string) error {
	image := fmt.Sprintf("%s:%s", imgName, version)

	if !Confirm(fmt.Sprintf("Sure to remove tool %s by removing image %s@%s?", toolName, imgName, version)) {
		fmt.Println("Skipped: ", image)
		return nil
	}
	fmt.Printf("Removing tool %s by removing image %s@%s...\n", toolName, image, version)
	return t.rt.removeImage(image)
}

// runTool runs the given tool.
func (t *toolManager) runTool(tool types.Tool, version string, args []string) error {
	image := tool.Image
	if version != "" {
		image = fmt.Sprintf("%s:%s", tool.Image, version)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Can't get current working directory: %v", err)
	}
	volumes := make([]string, len(tool.Volumes))
	for i, vol := range tool.Volumes {
		volumes[i] = strings.ReplaceAll(vol, "$(pwd)", cwd)
	}

	return t.rt.runImage(
		image,
		tool.Entry,
		tool.Workdir,
		volumes,
		args,
	)
}

// checkTools inspects the given tools and returns their status.
func (t *toolManager) checkTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	localImages, err := t.rt.localImages()
	if err != nil {
		return nil, err
	}

	status := make(map[string]types.ToolStatus)
	for name, tool := range tools {
		repo, tag := SplitImage(tool.Image)
		repoBase := path.Base(repo)
		installed := false
		localTags := []string{}

		for localRepo, toolStatus := range localImages {
			localBase := path.Base(localRepo)
			if localRepo == tool.Image || localBase == repoBase {
				localTags = append(localTags, toolStatus.LocalTags...)
				checkTags := []string{}
				if tag != "" {
					checkTags = append(checkTags, tag)
				}
				checkTags = append(checkTags, "latest")

				for _, t := range checkTags {
					if slices.Contains(toolStatus.LocalTags, t) {
							installed = true
						}
					if installed {
						break
					}
				}
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
func (t *toolManager) cleanTools(tools map[string]types.Tool) error {
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
