package container

import (
	"github.com/YardRat0117/foxbox/internal/types"
)

// runtimeManager manages images CRUD and running, NOTHING to do with `tool`
type runtimeManager interface {
	checkImage(image string) (bool, error)
	pullImage(image string) error
	removeImage(image string) error
	localImages() (map[string]*types.ToolStatus, error)
	runImage(image string, entry string, workdir string, volumes []string, args []string) error
}
