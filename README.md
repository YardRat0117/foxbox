# ratbox

A lightweight CLI tool to run common development tools in containers. No need to install gcc, python, node, or rust locally — you can just use `ratbox` and run commands in isolated containers.

---

## Features

- Run tools via containers with a single command:

```bash
  # The `--` ensures arguments are passed to the contained tool.
  ratbox python -- main.py
  ratbox gcc -- hello.c -o hello
```

- Auto-mount your current directory into the container
- Configurable tool mappings in `~/.ratbox/config.yml`
- Supports Podman (Docker support coming later)

---

## Installation

### From source

```bash
git clone https://github.com/yardrat0117/ratbox.git
cd ratbox
go build -o ratbox ./rbox
```

### Add to PATH
```bash
export PATH=$PATH:/path/to/ratbox
```

---

## Configuration

Default tools are defined in `config/default.yml`:

```yml
tools:
  gcc:
    image: gcc:latest
    entry: gcc
    workdir: /work
    volumes:
      - $(pwd):/work

  python:
    image: python:latest
    entry: python
    workdir: /work
    volumes:
      - $(pwd):/work

  node:
    image: node:latest
    entry: node
    workdir: /work
    volumes:
      - $(pwd):/work

  rust:
    image: rust:latest
    entry: rustc
    workdir: /work
    volumes:
      - $(pwd):/work
```


---

## Roadmap

- [ ] List available tools (`ratbox list`)
- [ ] Install/pull images (`ratbox install <tool>`)
- [ ] Support multiple versions (`ratbox python@3.9`)
- [ ] Add Docker support

---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for details.

---

## About

I'm a 2nd-year CS undergrad. This is a personal project and I may not have much time to maintain it.

However, if you encounter any issues, feel free to contact me via GitHub or via email - I'll try to help if I can.

If you like it, star it please :)
