# foxbox

A lightweight CLI tool to run common development tools in containers. No need to install compilers locally — you can just use `foxbox` and run commands in isolated containers.

---

## Features

- Run tools in containers with a single command
- Auto-mount directory into the container
- Supports Podman (Docker support coming later)

---

## Installation

1. Build from source

```bash
git clone https://github.com/yardrat0117/foxbox.git
cd foxbox
./build.sh
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

2. Install a configured tool

```bash
foxbox install gcc
foxbox install python@3.12
```

3. Run a specified tool

```bash
foxbox run gcc -- hello.c -o hello
foxbox run python@3.12 -- hello.py
```

4. Remove a specified tool
```bash
foxbox remove gcc
foxbox remove python@3.12
```


5. Check version

```bash
foxbox version
```

---

## Configuration

Default tools are provided in `config/default.yml`. 
You can also DIY - if `~/.config/foxbox.yml` is provided, this would override `config/default.yml`.
Note: `$(pwd)` is the supported representation for current working directory in `foxbox.yml`. Please adhere to this represetation!

---

## Roadmap

- Functionality
    - [x] List available tools in config file(`foxbox list`)
    - [x] Install tools/pull images (`foxbox install <tool>`)
    - [x] Support multiple versions (`foxbox python@3.9`)
    - [ ] Manage tools 
        - [x] Check installed tools (merged into `foxbox list`)
        - [x] Check tool tags (merged into `foxbox list`)
        - [x] Remove installed tools (`foxbox remove python@3.9`)
        - [ ] Prune cached data (`foxbox prune`)
    - [x] Check foxbox version (`foxbox version`)
    - [ ] Add Docker support
- Reliability
    - [ ] Code quality
        - [x] Comments
        - [x] Refactor
    - [ ] Boot speed optimization
        - [ ] Integrate with Podman REST API via Go bindings
    - [ ] Error hints
    - [ ] Logging
- Vim Integration
    - [ ] Vim plugin (planned for another project)
- Community
    - [ ] Across-platform support (Linux/macOS/WSL/Windows)
    - [ ] Detailed docs, tutorials and wikis

Current planned goals:

1. Implement the `prune` command
2. Turn to Podman REST API
3. Add Docker Support

---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for details.

---

## About

I'm a 2nd-year CS undergrad. This is a personal project and I may not have much time to maintain it.

However, if you encounter any issues, feel free to contact me via GitHub or via email - I'll try to help if I can.

If you like it, star it please :)
