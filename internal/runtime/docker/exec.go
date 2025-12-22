package docker

import (
	"context"
	"io"
	"os"

	mobyContainer "github.com/moby/moby/api/types/container"
	mobyClient "github.com/moby/moby/client"
)

type dockerExecution struct {
	cli *mobyClient.Client
	id  string
}

func (e *dockerExecution) Attach(ctx context.Context) error {
	attach, err := e.cli.ContainerAttach(ctx, e.id, mobyClient.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}

	go func() {
		// Just ignore
		_, _ = io.Copy(os.Stdout, attach.Reader)
	}()
	return nil
}

func (e *dockerExecution) Wait(ctx context.Context) error {
	statusCh, errCh := e.cli.ContainerWait(ctx, e.id, mobyContainer.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return err
	case <-statusCh:
		return nil
	}
}

func (e *dockerExecution) Close(ctx context.Context) error {
	return e.cli.ContainerRemove(ctx, e.id, mobyClient.ContainerRemoveOptions{
		Force: true,
	})
}
