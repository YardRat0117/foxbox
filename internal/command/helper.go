package command

import (
	"fmt"
	"strings"
	"os"
)

func fatal(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}

func parseToolArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)
	name = parts[0]
	version = "latest"
	if len(parts) == 2 && parts[1] != "" {
		version = parts[1]
	}
	return
}
