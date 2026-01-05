## Tool-related Use Cases

### RunTool

```go
func (a *App) RunTool(ctx context.Context, args []string) (int, error)
```

#### Purpose

Executes a configured tool inside a container and returns the tool’s exit code.

This is the **primary execution path** of foxbox.

#### Behavior

- Parse CLI arguments

- Resolve Tool from Config

- Resolve image reference

- Ensure image exists

- Create container

- Execute container

- Cleanup

- Return exit code

### ListTool

```go
func (a *App) ListTool(ctx context.Context) error
```

#### Purpose

Lists all configured tools and indicates whether their images are installed locally.


#### Behavior

- Queries the runtime for locally available images

- Groups images by image name and tag

- For each configured tool, checks if matching local image, and print results

### InstallTool

```go
func (a *App) InstallTool(ctx context.Context, args []string) error
```

#### Purpose

Ensures that a tool’s container image is available locally.

#### Behavior

- Parses `<tool>[@version]`

- Resolves tool from Config

- Constructs a concrete image reference

- Calls `rt.EnsureImage`

### RemoveTool

```go
func (a *App) RemoveTool(ctx context.Context, args []string) error
```

#### Purpose

Removes a specific version of a tool’s image from the local runtime.

---

#### Behavior

- Parses `<tool>[@version]`

- Resolves tool from Config

- Constructs image reference

- Calls `rt.RemoveImage`

### CleanTool

```go
func (a *App) CleanTool(ctx context.Context) error
```

#### Purpose

Removes **all local images** corresponding to configured tools.

#### Behavior

- Lists all local images

- Matches images against configured tool image names

- Removes every matching tag
