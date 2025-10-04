package container

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Confirm checks with user about msg
func Confirm(msg string) bool {
	fmt.Print(msg, " [y/N]")
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return false
	}
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// SplitImage splits given image:tag into separated segments
func SplitImage(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
	}
	return image, "latest"
}

var imageRegex = regexp.MustCompile(`^[a-z0-9]+([._/-][a-z0-9]+)*(:[a-zA-Z0-9._-]+)?$`)

// verifyImage checks whether the given image name is a valid container image name with an optional tag.
func verifyImage(image string) error {
	if !imageRegex.MatchString(image) {
		return errors.New("invalid image name: " + image)
	}
	return nil
}

var entryRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// verifyEntry checks whether the given entry command is valid
func verifyEntry(entry string) error {
    if !entryRegex.MatchString(entry) {
        return fmt.Errorf("invalid entry command: %s", entry)
    }
    return nil
}
