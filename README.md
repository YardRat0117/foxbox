# foxbox

A lightweight CLI tool to run common development tools in containers. No need to install compilers locally — you can just use `foxbox` and run commands in isolated containers.

---

## Features

- Run tools in containers with a single command
- Auto-mount your current directory into the container
- Supports Podman (Docker support coming later)

---

## Installation

1. Build from source

```bash
git clone https://github.com/yardrat0117/foxbox.git
cd foxbox
go build -o foxbox ./src
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
foxbox list
```

2. Install a specified tool (`latest` only)

```bash
foxbox install gcc
```

3. Run a specified tool (with arguments for the tool)

```bash
# If version not specified, `foxbox` pulls the `latest` by default
foxbox run gcc -- hello.c -o hello

# Version should be specified with `@`
foxbox run python@3.12 -- hello.py
```

---

## Configuration

Default tools are provided in `config/default.yml`. 
You can also DIY - if `~/.config/rbox.yml` is provided, this would override `config/default.yml`.
Note: `$(pwd)` is the supported representation for current working directory in `rbox.yml`. Please adhere to this represetation!

---

## Roadmap

- Functionality
    - [x] List available tools in config file(`foxbox list`)
    - [x] Install tools/pull images (`foxbox install <tool>`)
    - [x] Support multiple versions (`foxbox python@3.9`)
    - [ ] Manage tools 
        - [x] Check installed tools (merged into `foxbox list`)
        - [x] Check tool tags (merged into `foxbox list`)
        - [ ] Remove installed tools (`foxbox rm python@3.9`)
        - [ ] Prune cached data (`foxbox prune`)
    - [ ] Check foxbox version (`foxbox version`)
    - [ ] Add Docker support
- Reliability
    - [ ] Code quality
        - [ ] Comments
        - [ ] Refactoring
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
