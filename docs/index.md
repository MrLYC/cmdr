# CMDR Documentation

[![unittest](https://github.com/mrlyc/cmdr/actions/workflows/unittest.yml/badge.svg)](https://github.com/mrlyc/cmdr/actions/workflows/unittest.yml)
[![integration-test](https://github.com/mrlyc/cmdr/actions/workflows/integration-test.yml/badge.svg)](https://github.com/mrlyc/cmdr/actions/workflows/integration-test.yml)
[![codecov](https://codecov.io/gh/MrLYC/cmdr/branch/master/graph/badge.svg?token=mo4TJP4mQt)](https://codecov.io/gh/MrLYC/cmdr)
![Go version](https://img.shields.io/github/go-mod/go-version/mrlyc/cmdr)
![release](https://img.shields.io/github/v/release/mrlyc/cmdr?label=version)

**CMDR** is a simple command version management tool that helps you quickly switch between multiple command versions.

## What is CMDR?

CMDR (Command Manager) solves a common developer problem: managing multiple versions of CLI tools. Whether you need different versions of `node`, `python`, `go`, or any other command-line tool, CMDR provides a unified interface to:

- **Install** commands from URLs or local paths
- **Switch** between different versions seamlessly
- **List** all installed versions of a command
- **Activate/Deactivate** specific versions

## Quick Example

```shell
# Install a command version
cmdr command install -n kubectl -v 1.28.0 -l https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl

# List installed versions
cmdr command list -n kubectl

# Switch to a specific version
cmdr command use -n kubectl -v 1.28.0
```

## Key Features

- **Version Management**: Install and switch between multiple versions of any command
- **URL Replacement**: Speed up downloads by configuring proxy URLs[^1]
- **Self-Upgrade**: Update CMDR itself with a single command
- **Cross-Platform**: Works on Linux and macOS

## Documentation Structure

| Section | Description |
|---------|-------------|
| [Getting Started](getting-started/installation.md) | Installation and first steps |
| [Architecture](architecture/overview.md) | System design and internals |
| [Components](components/cli-commands.md) | Detailed component documentation |
| [API](api/core-interfaces.md) | Interface and configuration reference |
| [Operations](operations/build-and-test.md) | Build, test, and release processes |

## Technology Stack

CMDR is built with Go 1.25 and uses:

- **[Cobra](https://github.com/spf13/cobra)** - CLI framework[^2]
- **[Viper](https://github.com/spf13/viper)** - Configuration management[^2]
- **[Storm](https://github.com/asdine/storm)** - BoltDB wrapper for data persistence[^3]
- **[go-getter](https://github.com/hashicorp/go-getter)** - URL-based file fetching[^4]

---

[^1]: See URL replacement feature in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L166-L179
[^2]: CLI setup in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L20-L56
[^3]: Database initialization in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L147-L164
[^4]: Fetcher implementation in [`core/fetcher/go_getter.go`](https://github.com/mrlyc/cmdr/blob/master/core/fetcher/go_getter.go)
