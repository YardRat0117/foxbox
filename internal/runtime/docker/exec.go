package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/moby/api/pkg/stdcopy"
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
		Logs:   true,
	})
	if err != nil {
		return err
	}

	go func() {
		defer attach.Close()
		// Just ignore
		_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader)
	}()
	return nil
}

func (e *dockerExecution) Wait(ctx context.Context) (int, error) {
	statusCh, errCh := e.cli.ContainerWait(
		ctx,
		e.id,
		mobyContainer.WaitConditionNotRunning,
	)

	select {
	case err := <-errCh:
		// Docker wait internal failure
		return 0, err
	case status := <-statusCh:
		// Container exited
		if status.Error != nil {
			return int(status.StatusCode), fmt.Errorf("%s", status.Error.Message)
		}
		return int(status.StatusCode), nil
	}
}

func (e *dockerExecution) Close(ctx context.Context) error {
	return e.cli.ContainerRemove(ctx, e.id, mobyClient.ContainerRemoveOptions{
		Force: true,
	})
}
