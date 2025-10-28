package command

import (
	"strings"
)

// func parseToolArg is a utility func for slicing tool name and version
func parseToolArg(arg string) (name, version string) {
	// split by `@`
	parts := strings.SplitN(arg, "@", 2)

	// Assign name
	name = parts[0]

	// Assign version, with "latest" as default
	version = "latest"
	if len(parts) == 2 && parts[1] != "" {
		version = parts[1]
	}

	return
}
