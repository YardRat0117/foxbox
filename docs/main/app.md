## App

### What it is

`App` is the *application layer* of foxbox.

It implementes the **core user cases exposed to the CLI**, and acts as the **orchestrator** between:

- user input (CLI arguments),

- static configuration (`Config` / `Tool`),

- and the runtime backend (`runtime.Runtime`).

`App` contains **no runtime-specific logic**, but is responsible for **assembling domain objects** and invoking runtime opoerations in the correct order.

### Where it lives

An `App` would be injected as `*app.App` from the `command.rootContext`.

### What it contains

```Go
type App struct {
    cfg *types.Config
    rt runtime.Runtime
}
```

- `cfg` - immutable configuration loaded at startup

- `rt` - an abstract runtime backend implementing `runtime.Runtime`

### Responsibilities

`App` is responsible for:

- Parsing CLI arguments into semantic actions

- Resolving tools from configuration

- Constructing runtime execution specifications

- Coordinating runtime lifecycle calls

- Translating runtime errors into user-facing failures

Detailed, tool-related use cases are documented in [this independent doc](tool-related.md).

> Note: App methods return:

  - either an exit code (for `RunTool`)

  - or an error

- Internal failures are mapped to a fixed exit code (`70`)

- Runtime errors are surfaced directly to the user
