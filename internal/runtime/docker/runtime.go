// Package docker - docker impl package for runtime
package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/YardRat0117/foxbox/internal/domain"
	"github.com/YardRat0117/foxbox/internal/runtime"

	mobyContainer "github.com/moby/moby/api/types/container"
	mobyClient "github.com/moby/moby/client"
)

// Runtime - docker impl runtime
type Runtime struct {
	RuntimeURL string
}

var _ runtime.Runtime = (*Runtime)(nil)

// --- Runtime interface implementation ---

// EnsureEnv - docker impl
func (d *Runtime) EnsureEnv(
	ctx context.Context,
	ref domain.EnvRef,
) error {
	oci, ok := ref.(domain.OCIEnvRef)
	if !ok {
		return errors.New("docker runtime only supports OCI environment")
	}

	image := oci.Image

	// Use "latest" as default tag
	if !strings.Contains(image, ":") {
		image += ":latest"
	}

	cli, err := d.client()
	if err != nil {
		return err
	}

	images, err := cli.ImageList(ctx, mobyClient.ImageListOptions{})
	if err != nil {
		return err
	}

	for _, img := range images {
		if slices.Contains(img.RepoTags, image) {
			return nil
		}
	}

	// Pull image if missing
	fmt.Printf("Pulling image %s...\n", image)
	out, err := cli.ImagePull(ctx, image, mobyClient.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		// Simply ignore it
		_ = out.Close()
	}()

	decoder := json.NewDecoder(out)
	for {
		var msg map[string]any
		if err := decoder.Decode(&msg); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("decoding image pull response: %w", err)
		}
		if status, ok := msg["status"].(string); ok {
			if strings.Contains(status, "Pulling from") {
				if id, ok := msg["id"].(string); ok {
					fmt.Printf("%s:%s\n", status, id)
				} else {
					fmt.Println(status)
				}
			}
		}
	}

	return nil
}

// RemoveEnv - docker impl
func (d *Runtime) RemoveEnv(
	ctx context.Context,
	ref domain.EnvRef,
) error {
	oci, ok := ref.(domain.OCIEnvRef)
	if !ok {
		return errors.New("docker runtime only supports OCI environment")
	}

	cli, err := d.client()
	if err != nil {
		return err
	}

	_, err = cli.ImageRemove(
		ctx,
		oci.Image,
		mobyClient.ImageRemoveOptions{Force: true},
	)
	return err
}

// HasEnv - docker impl
func (d *Runtime) HasEnv(
	ctx context.Context,
	ref domain.EnvRef,
) (bool, error) {
	oci, ok := ref.(domain.OCIEnvRef)
	if !ok {
		return false, errors.New("docker runtime only supports OCI environment")
	}

	cli, err := d.client()
	if err != nil {
		return false, err
	}

	_, err = cli.ImageInspect(ctx, oci.Image)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateSandbox - docker impl
func (d *Runtime) CreateSandbox(
	ctx context.Context,
	spec domain.SandboxSpec,
) (domain.SandboxID, error) {
	cli, err := d.client()
	if err != nil {
		return "", err
	}

	// OCI Environment dispatch
	ociEnv, ok := spec.Env.(domain.OCIEnvRef)
	if !ok {
		return "", errors.New("docker runtime only supports OCI environment")
	}

	// Mounts
	binds := make([]string, 0, len(spec.Mounts))
	for _, m := range spec.Mounts {
		bind := m.Source + ":" + m.Target
		if m.ReadOnly {
			bind += ":ro"
		}
		binds = append(binds, bind)
	}

	// Env Vars
	var env []string
	for k, v := range spec.EnvVar {
		env = append(env, k+"="+v)
	}

	cfg := mobyContainer.Config{
		Image:      ociEnv.Image,
		Cmd:        spec.Cmd,
		WorkingDir: spec.WorkDir,
		Env:        env,
		User:       spec.User,
		Tty:        false,
	}

	hostCfg := mobyContainer.HostConfig{
		Binds: binds,
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&cfg,
		&hostCfg,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", err
	}

	return domain.SandboxID(resp.ID), nil
}

// StartSandbox - docker impl
func (d *Runtime) StartSandbox(
	ctx context.Context,
	id domain.SandboxID,
) error {
	cli, err := d.client()
	if err != nil {
		return err
	}
	return cli.ContainerStart(ctx, string(id), mobyClient.ContainerStartOptions{})
}

// WaitSandbox - docker impl
func (d *Runtime) WaitSandbox(
	ctx context.Context,
	id domain.SandboxID,
) (domain.SandboxResult, error) {
	cli, err := d.client()
	if err != nil {
		return domain.SandboxResult{}, err
	}

	statusCh, errCh := cli.ContainerWait(
		ctx,
		string(id),
		mobyContainer.WaitConditionNotRunning,
	)

	select {
	case err := <-errCh:
		if err != nil {
			return domain.SandboxResult{}, err
		}
	case status := <-statusCh:
		return domain.SandboxResult{
			ExitCode: int(status.StatusCode),
		}, nil
	}

	return domain.SandboxResult{}, ctx.Err()
}

// RemoveSandbox - docker impl
func (d *Runtime) RemoveSandbox(
	ctx context.Context,
	id domain.SandboxID,
) error {
	cli, err := d.client()
	if err != nil {
		return err
	}
	return cli.ContainerRemove(ctx, string(id), mobyClient.ContainerRemoveOptions{
		Force: true,
	})
}
