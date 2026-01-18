# Configuration Keys

Complete reference of all configuration keys in CMDR.

## Core Configuration

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `core.root_dir` | `~/.cmdr` | string | Root directory for CMDR data |
| `core.bin_dir` | `bin` | string | Directory for active command symlinks (relative to root) |
| `core.shims_dir` | `shims` | string | Directory for command version binaries (relative to root) |
| `core.profile_dir` | `profile` | string | Directory for shell initialization scripts (relative to root) |
| `core.database_path` | `cmdr.db` | string | Path to BoltDB database (relative to root) |
| `core.profile_path` | (auto-detected) | string | Path to shell profile file (~/.bashrc, ~/.zshrc, etc.) |
| `core.shell` | (auto-detected) | string | Current shell executable |
| `core.config_path` | `~/.cmdr/config.yaml` | string | Configuration file path |
| `core.link_mode` | `default` | string | How to link binaries: `copy` or `link` |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L23-L34

## Logging Configuration

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `log.level` | `info` | string | Log level: `trace`, `debug`, `info`, `warn`, `error` |
| `log.output` | `stderr` | string | Log output: `stdout` or `stderr` |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L40-L42

## Proxy Configuration

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `proxy.go` | - | string | Go module proxy URL (sets `GOPROXY`) |
| `proxy.http` | - | string | HTTP proxy URL (sets `HTTP_PROXY`) |
| `proxy.https` | - | string | HTTPS proxy URL (sets `HTTPS_PROXY`) |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L35-L38

## Download Configuration

### General

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `download.replace` | - | object | Legacy URL replacement pattern |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L44-L45

### Direct Strategy

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `download.direct.timeout` | 30 | int | Timeout in seconds |
| `download.direct.max_retries` | 3 | int | Maximum retry attempts |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L48-L49

### Proxy Strategy

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `download.proxy.enabled` | false | bool | Enable proxy strategy |
| `download.proxy.type` | `http` | string | Proxy type: `http` or `socks5` |
| `download.proxy.address` | - | string | Proxy server address |
| `download.proxy.timeout` | 30 | int | Timeout in seconds |
| `download.proxy.max_retries` | 3 | int | Maximum retry attempts |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L51-L55

### Rewrite Strategy

| Key | Default | Type | Description |
|-----|---------|------|-------------|
| `download.rewrite.rule` | - | string | URL rewrite template |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L57

## CLI Command Configuration

These keys are transient, used only during command execution:

### command install

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.install.name` | `-n, --name` | Command name |
| `_.command.install.version` | `-v, --version` | Version string |
| `_.command.install.location` | `-l, --location` | Download URL or file path |
| `_.command.install.activate` | `-a, --activate` | Activate after install |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L65-L68

### command define

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.define.name` | `-n, --name` | Command name |
| `_.command.define.version` | `-v, --version` | Version string |
| `_.command.define.location` | `-l, --location` | File path |
| `_.command.define.activate` | `-a, --activate` | Activate after define |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L60-L63

### command list

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.list.name` | `-n, --name` | Filter by name |
| `_.command.list.version` | `-v, --version` | Filter by version |
| `_.command.list.location` | `-l, --location` | Filter by location |
| `_.command.list.activate` | `-a, --activate` | Filter by activation |
| `_.command.list.fields` | `--fields` | Fields to display |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L70-L74

### command remove

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.remove.name` | `-n, --name` | Command name |
| `_.command.remove.version` | `-v, --version` | Version to remove |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L76-L77

### command unset

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.unset.name` | `-n, --name` | Command name |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L79

### command use

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.command.use.name` | `-n, --name` | Command name |
| `_.command.use.version` | `-v, --version` | Version to activate |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L81-L82

### config get

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.config.get.key` | `-k, --key` | Configuration key |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L85

### config set

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.config.set.key` | `-k, --key` | Configuration key |
| `_.config.set.value` | `-v, --value` | Configuration value |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L87-L88

### init

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.init.upgrade` | `--upgrade` | Whether this is an upgrade |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L91

### upgrade

| Key | CLI Flag | Description |
|-----|----------|-------------|
| `_.upgrade.release` | `-r, --release` | Release name to upgrade to |
| `_.upgrade.asset` | `-a, --asset` | Asset name to download |
| `_.upgrade.args` | `--args` | Additional arguments |

**Source:** [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L94-L96

## Environment Variable Mapping

All configuration keys can be set via environment variables:

**Format:** `CMDR_<KEY>` where `<KEY>` is the config key with `.` replaced by `_` and uppercased.

**Examples:**

| Config Key | Environment Variable |
|------------|---------------------|
| `core.root_dir` | `CMDR_CORE_ROOT_DIR` |
| `log.level` | `CMDR_LOG_LEVEL` |
| `download.direct.timeout` | `CMDR_DOWNLOAD_DIRECT_TIMEOUT` |
| `proxy.http` | `CMDR_PROXY_HTTP` |

## Configuration File Example

```yaml
# ~/.cmdr/config.yaml

# Core settings
core:
  root_dir: ~/.cmdr
  bin_dir: bin
  shims_dir: shims
  profile_dir: profile
  database_path: cmdr.db
  link_mode: copy

# Logging
log:
  level: info
  output: stderr

# Proxy settings
proxy:
  go: https://goproxy.cn,direct
  http: http://proxy.example.com:8080
  https: http://proxy.example.com:8080

# Download settings
download:
  # Direct strategy
  direct:
    timeout: 30
    max_retries: 3
  
  # Proxy strategy
  proxy:
    enabled: true
    type: http
    address: http://download-proxy.example.com:8080
    timeout: 60
    max_retries: 5
  
  # URL rewrite for GitHub
  rewrite:
    rule: "https://ghproxy.com/{{.URI}}"
  
  # Legacy replacement (still supported)
  replace:
    match: "^https://raw.githubusercontent.com/.*$"
    template: "https://ghproxy.com/{{ .input | urlquery }}"
```
