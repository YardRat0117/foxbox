// Package domain defines core domain types and invariants.
package domain

import (
	"errors"
	"regexp"
)


// ImageRef represents a validated container image reference.
type ImageRef struct {
	Raw string
}

var imageRegex = regexp.MustCompile(
	`^[a-z0-9]+([._/-][a-z0-9]+)*(:[a-zA-Z0-9._-]+)?$`,
)

// NewImageRef constructs an ImageRef from a raw string after validation.
func NewImageRef(raw string) (ImageRef, error) {
	if !imageRegex.MatchString(raw) {
		return ImageRef{}, errors.New("invalid image ref: " + raw)
	}
	return ImageRef{Raw: raw}, nil
}
