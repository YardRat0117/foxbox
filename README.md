# ratbox

A lightweight CLI tool to run common development tools in containers. No need to install compilers locally — you can just use `ratbox` and run commands in isolated containers.

---

## Features

- Run tools via containers with a single command:
- Auto-mount your current directory into the container
- Supports Podman (Docker support coming later)

---

## Installation

1. Build from source

```bash
git clone https://github.com/yardrat0117/ratbox.git
cd ratbox
go build -o ratbox ./rbox
```

2. Prepare Podman (e.g., via APT)

```bash
sudo apt update
sudo apt install podman
```

3. Ready to use!

---

## Usage

1. List all configured tools

```bash
ratbox list
```

2. Install a specified tool

```bash
ratbox install gcc
```

3. Run a specified tool (with arguments for the tool)

```bash
// If version not specified, `ratbox` pulls the `latest` by default
ratbox run gcc -- hello.c -o hello

// Version should be specified with `@`
ratbox run python@3.12 -- hello.py
```

---

## Configuration

Default tools are provided in `config/default.yml`. 
You can also DIY - if `~/.config/rbox.yml` is provided, this would override `config/default.yml`.
Note: `$(pwd)` is the supported representation for current working directory in `rbox.yml`. Please adhere to this represetation!

---

## Roadmap

- Functionality
    - [x] List available tools (`ratbox list`)
    - [x] Install/pull images (`ratbox install <tool>`)
    - [x] Support multiple versions (`ratbox python@3.9`)
    - [ ] Add Docker support
- Reliability
    - [ ] Boot speed optimization
    - [ ] Error hints
    - [ ] Logging
- Vim Integration
    - [ ] Vim plugin (planned for another project)
- Community
    - [ ] Across-platform support (Linux/macOS/WSL/Windows)
    - [ ] Detailed docs, tutorials and wikis

---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for details.

---

## About

I'm a 2nd-year CS undergrad. This is a personal project and I may not have much time to maintain it.

However, if you encounter any issues, feel free to contact me via GitHub or via email - I'll try to help if I can.

If you like it, star it please :)
