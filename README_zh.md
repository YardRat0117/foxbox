# Foxbox

[English](./README.md) | [中文](./README_zh.md)

## 什么是 Foxbox

Foxbox 是一个使用 Go 编写的轻量级命令行工具，用于在容器中运行开发工具。

它最初是为 `gcc`、`clang` 等编译器设计的，但现在几乎适用于任何命令行工具。

正如名字所示，你可以轻松获取一个 “box”，立刻启动任务 —— 干净、隔离、即开即用。

---

## 主要特性

- **可复现、隔离的环境**
  在干净、可复现的容器中运行命令行工具，无需本地安装。

- **简化容器交互**
  不需要复杂的容器参数，即可享受容器化带来的开发体验 —— 即使你不是 DevOps 专家。

- **自动化的环境管理**
  不再为版本控制、挂载目录这些繁琐的操作烦恼 —— Foxbox 会自动帮你处理好。

---

## 工作原理

Foxbox 最初的灵感来源于这样一个简单的别名：

```shell
alias c_gcc='podman run -it --rm -v .:/workspace -w /workspace gcc:latest gcc'
```

而 Foxbox 也在此基础上一脉相承：**用简洁的命令封装容器化的工具**，让它们更易用、更自动化。

从本质上讲，Foxbox 并不直接管理容器 —— 它只是读取你的命令，翻译成合适的容器参数，并调用底层工具执行 ——  我们相信容器开发者们，并让他们这些专业的人来做专业的事。

---

## 安装方法

你可以通过以下两种方式安装 Foxbox：

### 1. 从源码构建

直接克隆仓库并构建即可：

```shell
git clone https://github.com/yardrat0117/foxbox.git
cd foxbox
make build
```

生成的 `foxbox` 可执行文件会放在仓库的 `bin` 目录下。
记得将它移动到你的 `$PATH` 目录中。

> 提示：如果你在意二进制文件大小，可以使用 `make` 命令来进行 [UPX](https://upx.github.io/) 压缩（可选）。

---

### 2. 通过 Go 安装

如果你的系统中已经安装了 Go，可以直接使用以下命令安装：

```shell
go install github.com/YardRat0117/foxbox@latest
```

可执行文件会被安装到你的 `$GOBIN`（通常是 `~/go/bin`，取决于 Go 的配置）。

> 注意：Go 模块发布的版本可能略落后于源码的最新标签。

---

## 使用方法

在使用 Foxbox 之前，你需要先配置要运行的工具。

配置文件是一个名为 `foxbox.yml` 的 YAML 文件，位于 `~/.config` 目录下。
示例配置文件见 `configs/demo.yml`。以下是字段说明：

| 字段名        | 说明                                    |
| :------------ | :-------------------------------------- |
| `tools`       | 工具配置的根节点                        |
| `<tool-name>` | 用户自定义的工具别名                    |
| `image`       | 使用的容器镜像名称或标签                |
| `entry`       | 要执行的命令（通常是工具的主命令）      |
| `workdir`     | 容器内的工作目录（推荐 `/work`）        |
| `volumes`     | 主机目录挂载列表（推荐 `$(pwd):/work`） |

此外，Foxbox 提供了友好的命令接口：

1. `list` - 列出所有已配置的工具
2. `install` - 安装一个工具
3. `run` - 运行一个工具
4. `remove` - 移除一个工具
5. `clean` - 清除所有配置和已安装工具
6. `version` - 查看版本信息

可以通过 `foxbox <command> --help` 查看详细帮助。

---

## 开发路线图

> 说明：Foxbox 是一个个人独立项目，由我独自开发和维护。
> 因此功能更新与问题修复的进度可能会比社区驱动的项目慢一些。
> 该路线图可能会随时间调整或扩展，如果你愿意参与，非常欢迎关注与贡献！

* 核心功能

  * [x] 基本命令支持（`run`、`install`、`remove` 等）
  * [x] 工具版本控制支持（`<tool>@<tag>`）
  * [x] Podman CLI 支持
  * [x] Docker CLI 支持
  * [x] Docker Engine API 集成
  * [ ] Podman Libpod API 集成
* 工程实践

  * [ ] 注释覆盖完善
  * [ ] 单元测试覆盖完善
  * [ ] 结构化日志支持
  * [ ] 统一的 API 抽象层（`gRPC` / `HTTP RESTful`）
* 扩展功能

  * [ ] 基于 IaC 的配置兼容性（`Ansible Playbooks`）
  * [ ] 分离式运行后端（`Daemon backend` / `Remote K8s backend`）
  * [ ] 可观测性增强（`Prometheus`）
  * [ ] 用户体验增强（Vim 插件、GitHub CLI 插件）
* 社区

  * [x] 更新英文版 `README.md`
  * [ ] 更新中文版 `README_zh.md`

---

## 许可证

本项目基于 Apache License 2.0 开源。
详情请参见 [LICENSE](./LICENSE)。

---

## 关于我

你好！我是 **YardRat**，Foxbox 的开发者 —— 这个项目是我在业余时间独立完成并持续维护的。

我目前是一名计算机科学专业的大二学生。
Foxbox 最初是我尝试做一个“认真点的开源项目”的尝试，同时也是为了解决我在日常 Linux 开发中遇到的一些痛点。

虽然项目仍在不断演进中，但我也在持续改进它、探索新的方向。

进度可能不会太快——我平时总是有很多课内的杂事，而且大多数开发都是在试错或者重构已有代码——但我很高兴它在稳步前进。

如果你在使用中遇到问题或有建议，欢迎通过 GitHub 提 Issue 或直接提交 PR！
如果你觉得这个项目有点意思，请帮我点个 Star 吧 :)
