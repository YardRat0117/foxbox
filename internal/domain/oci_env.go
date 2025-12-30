package domain

import (
    "errors"
    "regexp"
)

// OCIEnvRef for OCI-based env (image)
type OCIEnvRef struct {
    Image string
}

// Kind returns the EnvKind (EnvOCI actually)
func (r OCIEnvRef) Kind() EnvKind {
    return EnvOCI
}

var imageRegex = regexp.MustCompile(
    `^[a-z0-9]+([._/-][a-z0-9]+)*(:[a-zA-Z0-9._-]+)?$`,
)

// NewOCIEnvRef initializes a new OCI-based env with image name verification
func NewOCIEnvRef(raw string) (OCIEnvRef, error) {
    if !imageRegex.MatchString(raw) {
        return OCIEnvRef{}, errors.New("invalid OCI image ref: " + raw)
    }
    return OCIEnvRef{Image: raw}, nil
}

