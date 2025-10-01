package container

import (
	"fmt"
	"os/exec"
	"strings"
)

// --- helpers ---

// confirm checks with user about msg
func (d *DockerRuntime) confirm(msg string) bool {
	fmt.Print(msg, " [y/N]")
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return false
	}
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// splitImage splits given image:tag into separated segments
func (d *DockerRuntime) splitImage(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
	}
	return image, "latest"
}

// getLocalImages retrieves pulled images
func (d *DockerRuntime) getLocalImages() (map[string]map[string]struct{}, error) {
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	localImages := map[string]map[string]struct{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		repo, tag := d.splitImage(line)
		if _, ok := localImages[repo]; !ok {
			localImages[repo] = map[string]struct{}{}
		}
		localImages[repo][tag] = struct{}{}
	}
	return localImages, nil
}
