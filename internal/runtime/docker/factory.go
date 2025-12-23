package docker

import "github.com/YardRat0117/foxbox/internal/runtime"

// New constructs a Docker-based runtime.
func New(runtimeURL string) runtime.Runtime {
	return &Runtime{
		RuntimeURL: runtimeURL,
	}
}
