package domain

// MountSpec defines a filesystem mount from host to sandbox.
type MountSpec struct {
	Source   string
	Target   string
	ReadOnly bool
}
