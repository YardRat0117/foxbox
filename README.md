# foxbox

A lightweight CLI tool to run common development tools in containers. 

No need to install compilers locally — simply use `foxbox` to run commands in isolated containers.

---

## Features

- Run tools in containers with a single command
- Auto-mount directory into the container
- Supports Podman (Docker support coming soon)
- Manage multiple versions of tools effortlessly

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

```bash
foxbox version
```

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

5. Clean all installed tools
```bash
foxbox clean
```


6. Check version

```bash
foxbox version
```

---

## Configuration

- Default tools:`config/default.yml` 
- User override: `~/.config/foxbox.yml`
> Note: Use `$(pwd)` to represent the current working directory in your configuration.

---

## Roadmap

- Basic Functionality
	- [x] Simply run tools (`foxbox run`)
	- [x] Install tools (`foxbox install <tool>`)
	- [x] Support multiple versions (`<tool>@<tag>`)
    - [x] Check installed tools and tags (`foxbox list`)
    - [x] Remove installed tools (`foxbox remove <tool>`)
    - [x] Clean all installed tools (`foxbox clean`)
    - [x] Check foxbox version (`foxbox version`)
- Advanced Functionality
    - [x] Podman CLI support
    - [x] Docker CLI support
    - [x] Docker Engine API support
    - [ ] Podman Docker API support
    - [ ] Podman Libpod API support
	- [ ] Vim Plugin Integration (side-project)
- Community
    - [ ] Across-platform support (Linux/macOS/WSL)
    - [ ] Tutorials
    - [ ] Docs

---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for details.

---

## About

I'm a 2nd-year CS undergrad. This is a personal project and I may not have much time to maintain it.

However, if you encounter any issues, feel free to contact me via GitHub or via email - I'll try to help if I can.

If you like it, star it please :)
