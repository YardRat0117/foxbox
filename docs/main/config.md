## Config

### What it is

`Config` is the *only user-modifiable configuration input* of foxbox.

It is a **static YAML file** which declares which tools are available and how they should be executed.

At startup, foxbox loads the configuration once and deserializes it into an in-memory `types.Config` structure.

All later application behavior (listing tools, installing tools, running tools) is driven by this loaded configuration.

### Where it lives

Foxbox searches for the configuration file in the following order:

1. User configuration: `~/.config/foxbox.yml`

2. Default configuration: `config/default.yml`

If the user configuration existss, it takes precedence; Otherwise, the default configuration is used.

This behavior is implemented in:

```
internal/config/config.go
```

### What it contains

At the top level, the configuration consists primarily of a `tools` mapping.

Example:

```YAML
tools:
  gcc:
    image: gcc
    entry: gcc
    workdir: /work
    volumes:
      - $(pwd):/work
```

Conceputally, `tools` is a map; the key is the **tool name**, and the value is the **tool definition**

### Responsibilities

`Config` is responsible for:

- Declaring which tools exist

- Declaring how each tool shoulded be executed

- Providing static execution metadata to the application
