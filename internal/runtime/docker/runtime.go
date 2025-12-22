// Package docker - docker impl package for runtime
package docker

import (
	"context"
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

// EnsureImage - docker impl
func (d *Runtime) EnsureImage(
	ctx context.Context,
	ref domain.ImageRef,
) error {
	image := ref.Raw

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

	out, err := cli.ImagePull(ctx, image, mobyClient.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		// Simply ignore it
		_ = out.Close()
	}()

	_, _ = io.Copy(io.Discard, out)

	return nil
}

// Create - docker impl
func (d *Runtime) Create(
	ctx context.Context,
	spec domain.ContainerSpec,
) (domain.ContainerID, error) {
	cli, err := d.client()
	if err != nil {
		return "", err
	}

	cmd := append([]string{spec.Entry}, spec.Args...)

	binds := make([]string, 0, len(spec.Volumes))
	for _, v := range spec.Volumes {
		binds = append(binds, v.Host+":"+v.Guest)
	}

	cfg := mobyContainer.Config{
		Image:      spec.Image.Raw,
		Cmd:        cmd,
		WorkingDir: spec.Workdir,
		Tty:        true,
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

	return domain.ContainerID(resp.ID), nil
}

//Start - docker impl
func (d *Runtime) Start(
	ctx context.Context,
	id domain.ContainerID,
) error {
	cli, err := d.client()
	if err != nil {
		return err
	}
	return cli.ContainerStart(ctx, string(id), mobyClient.ContainerStartOptions{})
}

//Stop - docker impl
func (d *Runtime) Stop(
	ctx context.Context,
	id domain.ContainerID,
) error {
	cli, err := d.client()
	if err != nil {
		return err
	}
	return cli.ContainerStop(ctx, string(id), mobyClient.ContainerStopOptions{})
}

//RemoveImage - docker impl
func (d *Runtime) RemoveImage(
	ctx context.Context,
	ref domain.ImageRef,
) error {
	cli, err := d.client()
	if err != nil {
		return err
	}

	_, err = cli.ImageRemove(
		ctx,
		ref.Raw,
		mobyClient.ImageRemoveOptions{Force: true},
	)
	return err
}

//ListImage - docker impl
func (d *Runtime) ListImage(
	ctx context.Context,
) ([]domain.ImageInfo, error) {
	cli, err := d.client()
	if err != nil {
		return nil, err
	}

	images, err := cli.ImageList(ctx, mobyClient.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	var result []domain.ImageInfo
	for _, img := range images {
		for _, tag := range img.RepoTags {
			imageRef, _ := domain.NewImageRef(tag)
			result = append(result, domain.ImageInfo{
				Ref: imageRef,
			})
		}
	}

	return result, nil
}

//Exec - docker impl
func (d *Runtime) Exec(
	id domain.ContainerID,
) (runtime.Execution, error) {
	cli, err := d.client()
	if err != nil {
		return nil, err
	}
	return &dockerExecution{
		cli: cli,
		id:  string(id),
	}, nil
}
