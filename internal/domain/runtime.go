package domain

// ContainerID identifies a created container instance.
type ContainerID string

// ImageInfo describes a locally available container image.
type ImageInfo struct {
    Ref   ImageRef
    Tags  []string
    Size  int64
}

// VolumeSpec defines a host-to-container volume mapping.
type VolumeSpec struct {
    Host  string
    Guest string
}

// ContainerSpec describes how a container should be created.
type ContainerSpec struct {
    Image   ImageRef
    Cmd     []string
    WorkDir string
    Volumes []VolumeSpec
}

