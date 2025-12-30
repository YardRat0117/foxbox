// Package foxtypes define internal types used, "fox" from "foxbox"
package foxtypes

import (
	"github.com/YardRat0117/foxbox/internal/domain"
)

// Tool contains basic info for something to run in the container
type Tool struct {
	Entry   string
	Workdir string
	Mounts  []domain.MountSpec
	Env     domain.EnvRef
}
