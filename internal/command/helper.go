package command

import (
	"strings"
)

func parseToolArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)
	name = parts[0]
	version = "latest"
	if len(parts) == 2 && parts[1] != "" {
		version = parts[1]
	}
	return
}
