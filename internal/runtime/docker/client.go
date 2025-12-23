package docker

import (
	mobyClient "github.com/moby/moby/client"
)

func (d *Runtime) client() (*mobyClient.Client, error) {
	return mobyClient.NewClientWithOpts(
		mobyClient.WithHost(d.RuntimeURL),
		mobyClient.WithAPIVersionNegotiation(),
	)
}
