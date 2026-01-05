## Tool

### What it is

A `Tool` is a *user-facing execution unit* declared in `Config`.

It represents **"a predefined way to execute a command inside a virtual environment"** (for `docker` runtime, a container), and should be referenced directly by name from the CLI.

### Where it lives

A `Tool` is declared in `Config`.

### What it contains

A `Tool` contains following fields

- a name (used by CLI, usually the same as the entrypoint command)

- an entrypoint command

- a default working directory

- a set of associated mounting volumes

Example: suppose the following `Config` and CLI command

```YAML
myg++:
  image: gcc
  entry: g++
  workdir: /workspace
  volumes:
    - $PWD:/workspace
```

```shell
foxbox run myg++@15 -- --version
```

#### Image

The `image` field specifies the virtual environment template (for `docker` runtime, container image) to be used.

> Note: Versioning (tags) is **not enforced or interpreted at the `Tool` level**

#### Entry

The `entry` field defines the entrypoint command inside the virtual environment.

For example, the given CLI command would be actually executed as:

```shell
g++ --version
```

in the virtual environment with version tag `15` for tool `myg++`

#### Workdir and Volumes

The `workdir` and `volumes` defines how host directories would be mounted into the virtual environment.

For example, the given CLI command would mount `$PWD` (actually, the current working directory for the host) into the virtual environment as `/workspace`; and the command would be executed exactly at the `workdir`, namely, `/workspace`.

### Responsibilities

A `Tool` is a **static** and **declarative** description of:

- "which template to use"

- "which command to run"

- "where to run it"

- "which directories to mount"
