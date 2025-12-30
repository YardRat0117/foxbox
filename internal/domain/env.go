package domain

// EnvKind identifies the type of execution environment.
type EnvKind string

// EnvKind consts for runtimes
const (
	EnvOCI EnvKind = "oci"
)

// EnvRef represents a runtime-resolvable execution environment.
type EnvRef interface {
	Kind() EnvKind
}

// EnvInfo represents meta info about the env
type EnvInfo struct {
    Kind EnvKind
    Ref  string
    Size int64
}
