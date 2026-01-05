## Runtime

### What it is

`Runtime` is the **execution backend abstraction** of foxbox.

It defines the minimal set of capabilities required to:

- manage container images,
- create and control containers,
- execute commands inside containers,
- and observe execution results.

`Runtime` does not describe *how* these capabilities are implemented, only *what* must be provided by a concrete backend.

---

### Where it lives

An `Runtime` would be injected as `runtime.Runtime` from the `app.App`.

### What it contains

The runtime layer consists of two core interfaces:

- `Runtime`

```go
type Runtime interface {
    // image lifecycle
    EnsureImage(ctx context.Context, ref domain.ImageRef) error
    RemoveImage(ctx context.Context, ref domain.ImageRef) error
    ListImage(ctx context.Context) ([]domain.ImageInfo, error)

    // container lifecycle
    Create(ctx context.Context, spec domain.ContainerSpec) (domain.ContainerID, error)
    Start(ctx context.Context, id domain.ContainerID) error
    Stop(ctx context.Context, id domain.ContainerID) error
    Remove(ctx context.Context, id domain.ContainerID) error

    // execution
    Exec(id domain.ContainerID) (Execution, error)
}
```

- `Execution`

```go
type Execution interface {
    Attach(ctx context.Context) error
    Wait(ctx context.Context) (int, error)
    Close(ctx context.Context) error
}
```

### Responsibilities

The runtime layer is responsible for:

- Ensuring container images are available locally

- Listing and removing images

- Creating and destroying container instances

- Starting and stopping containers

- Attaching to execution I/O streams

- Waiting for execution completion and returning exit codes

- Cleaning up execution-related resources
