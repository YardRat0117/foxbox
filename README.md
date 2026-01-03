# Foxbox

[English](./README.md) | [中文](./README_zh.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/YardRat0117/foxbox.svg)](https://pkg.go.dev/github.com/YardRat0117/foxbox)

## What is Foxbox

Foxbox is a Go-based, lightweight CLI tool for running development tools inside containers.

Though originally designed for compilers like `gcc` and `clang`, it now works great with almost any CLI tool.

As the name suggests, you can simply fetch a "box" and launch your tasks instantly - clean, isolated, and ready to **go**.

---

## Key Features

- **Reproducible and isolated environments**
  Run CLI-based tools in clean, isolated and reproducible containers - without any extra local installation.

- **Simplified container interaction**
  Enjoy containerized workflows without dealing with complex container arguments - even if you aren't a DevOps expert.

- **Automated environment management**
  Forget about chores like version control and directory mounting - Foxbox handles them for you.

---

## How Foxbox works

Foxbox originates from an alias:

```shell
alias c_gcc ='podman run -it --rm -v .:/workspace -w /workspace gcc:latest gcc'
```

And it builds upon this idea - **wrapping containerized tools behind concise commands**, with better usability and automation.

In essence, Foxbox does NOT manage the container by itself - it simply reads your command, translate it into the proper container arguments, and invokes the underlying tools to work - because we believes in the experts who built them.

---

## Installation

You can install Foxbox in two ways:

1. Build from source

You can simply clone the repo and build it locally.

```shell
git clone https://github.com/yardrat0117/foxbox.git
cd foxbox
make build
```
The `foxbox` binary will be placed under the `bin` directory in the repo.

Remember to move it somewhere in your `$PATH`.

> Note: if you care about binary size, try `make` for additional [UPX](https://upx.github.io/)-based compression. This step is optional, though.

2. Install using Go

If you already have Go installed, you can simply install it directly via Go.

```shell
go install github.com/YardRat0117/foxbox/cmd/foxbox@latest
```

The `foxbox` binary will be installed into your `$GOBIN` (typically `~/go/bin`, depending on your Go configuration).

Note: official Go module releases may lag slightly behind the latest source tags.

---

## Usage

Before using Foxbox, you need to "configure" tools to run later.

This is done through a YAML file named `foxbox.yml`, located under your config directory `~/.config` dir.

An example is provided in `configs/demo.yml`. And here's these fileds explanation:
| Field         | Description                                                                           |
| :-----------: | :-----------------------------------------------------------------------------------: |
| `tools`       | The root section containing all the tools for Foxbox.                                 |
| `<tool-name>` | A user-defined alias for the tool                                                     |
| `image`       | The container image name or tag to use.                                               |
| `entry`       | The entrypoint command to run. (most of cases the tool command)                       |
| `workdir`     | The working directory inside the container. (`/work` suggested, align with `volumes`) |
| `volumes`     | A list of host-directory mounts. (`$PWD:/work` suggested, align with `workdir`)       |

Besides, Foxbox provides user-friendly commands to use:

1. `list` - List all configured tools

2. `install` - Install a configured tool

3. `run` - Run a specified tool

4. `remove` - Remove a specified tool

5. `clean` - Clean all configured and installed tools

You can try `foxbox <command> --help` for detailed help.

---

## Development Roadmap

> Note: Foxbox is a solo project, launched and maintained independently.
> As a result, progress - both in adding new features and fixing bugs - may be slower than those community-driven projects.
> This roadmap may change or expand over time, so if you'd like to help out (I'd really appreciate it!), feel free to check back regularly.

- Core Functionality
    - [x] Basic commands (`run`, `install`, `remove` and etc.)
    - [x] Tool version control support (`<tool>@<tag>`)
    - [x] Docker Engine API integration
    - [ ] Podman Libpod API integration (Staged)
- Engineer Robustness
    - [x] Architecture Refactor
    - [ ] Sufficient test coverage
- Extended Functionality (Staged)
    - [ ] Unified API abstraction layer (`gRPC` / `HTTP RESTful`)
    - [ ] IaC based configuration compatibility (`Ansible Playbooks`)
    - [ ] UX enhancement (Vim plugin, GitHub CLI plugin)
- Community
    - [x] Updated `README.md`
    - [ ] Updated Chinese version `README_zh.md`
- Deprecated Features
    - Podman / Docker CLI support with `os/exec`
---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for details.

---

## About

Hi! I'm **YardRat**, the developer behind Foxbox - this solo project built and maintained in my spare time.

As a second-year CS undergrad, I launched Foxbox both as a personal trial to work on a serious OSS project, and as a practical solution to problems I've encountered during my everyday development on Linux.

While the project is still evolving, I'm also continuously refining it and exploring new directions to work on.

Progress may not be fast, since I've always got a bunch of chores, and lots of my time are spent on trials and refactoring. However, I'm happy that the progress is steady.

If you encounter isues or have suggestions, feel free to reach out, open an issue or open a PR - contributions are always welcome!

If you find my stuff useful, please consider giving me a star :)
