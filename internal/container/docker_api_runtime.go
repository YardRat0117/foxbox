package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	foxTypes "github.com/YardRat0117/foxbox/internal/types"
	mobyContainer "github.com/moby/moby/api/types/container"
	mobyClient "github.com/moby/moby/client"
)

var _ runtimeManager = (*dockerAPIManager)(nil)

// dockerAPIManager implements imageManager with Docker Remote API
type dockerAPIManager struct {
	runtimeURL string
}

// checkImage checks if the given image exists locally.
func (d *dockerAPIManager) checkImage(image string) (bool, error) {
	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return false, err
	}

	// Attach `latest` by default
	if !strings.Contains(image, ":") {
		image = image + ":latest"
	}

	// Create a Docker client with the runtime URL (can be unix:///var/run/docker.sock)
	cli, err := mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.runtimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	// List all local images
	images, err := cli.ImageList(ctx, mobyClient.ImageListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list local images: %w", err)
	}

	// Iterate through all images and their tags to check for a match
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == image {
				// Image found locally
				return true, nil
			}
		}
	}

	// Image not found locally
	return false, nil
}

// pullImage pulls the specified image.
func (d *dockerAPIManager) pullImage(image string) error {
	// Verify image in case injection attack
	if err := verifyImage(image); err != nil {
		return err
	}

	cli, err := mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.runtimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	// Pull image
	out, err := cli.ImagePull(ctx, image, mobyClient.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, err)
	}
	defer out.Close()
	if _, err := io.Copy(io.Discard, out); err != nil {
		return fmt.Errorf("failed to discard image pull output: %w", err)
	}
	return nil
}

// removeImage removes the specified image.
func (d *dockerAPIManager) removeImage(image string) error {
	if err := verifyImage(image); err != nil {
		return err
	}

	cli, err := mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.runtimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()
	if _, err := cli.ImageRemove(ctx, image, mobyClient.ImageRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove image %s: %w", image, err)
	}

	return nil
}

// localImages retrieves pulled images
func (d *dockerAPIManager) localImages() (map[string]*foxTypes.ToolStatus, error) {
	cli, err := mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.runtimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()
	images, err := cli.ImageList(ctx, mobyClient.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list local images: %w", err)
	}

	localImages := map[string]*foxTypes.ToolStatus{}
	for _, img := range images {
		for _, tag := range img.RepoTags {
			repo, version := SplitImage(tag)
			if _, ok := localImages[repo]; !ok {
				localImages[repo] = &foxTypes.ToolStatus{
					Installed: true,
					LocalTags: []string{},
				}
			}
			localImages[repo].LocalTags = append(localImages[repo].LocalTags, version)
		}
	}

	return localImages, nil
}

// runImage runs a container from the image
func (d *dockerAPIManager) runImage(
	image string,
	entry string,
	workdir string,
	volumes []string,
	args []string,
) error {
	cli, err := mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.runtimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	// Ensure image is pulled
	exists, err := d.checkImage(image)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Printf("Image %s not found locally, pulling...\n", image)
		if err := d.pullImage(image); err != nil {
			return err
		}
	} else {
		fmt.Printf("Image %s found locally, skipping pull\n", image)
	}

	// Build container command
	cmd := append([]string{entry}, args...)

	resp, err := cli.ContainerCreate(ctx,
		&mobyContainer.Config{
			Image:      image,
			Cmd:        cmd,
			WorkingDir: workdir,
			Tty:        true,
		},
		&mobyContainer.HostConfig{
			Binds: volumes,
		},
		nil, nil, "",
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Attach container I/O
	attachResp, err := cli.ContainerAttach(ctx, resp.ID, mobyClient.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("failed to attach to container: %w", err)
	}
	defer attachResp.Close()

	// Copy container output to stdout asynchronously
	go func() {
		if _, err := io.Copy(os.Stdout, attachResp.Reader); err != nil {
			fmt.Fprintf(os.Stderr, "container output copy error: %v\n", err)
		}
	}()

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, mobyClient.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, mobyContainer.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error while waiting container: %w", err)
		}
	case <-statusCh:
	}

	// Remove container afterwards
	if err := cli.ContainerRemove(ctx, resp.ID, mobyClient.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	fmt.Printf("Container %s finished\n", resp.ID)
	return nil
}
