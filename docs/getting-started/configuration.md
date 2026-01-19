# Configuration

CMDR uses a YAML configuration file and supports environment variables for customization.

## Configuration File

The default configuration file is located at `~/.cmdr/config.yaml`. You can specify a different location using the `-c` or `--config` flag.

```shell
cmdr --config /path/to/config.yaml <command>
```

## Managing Configuration

### View Configuration

List all configuration values:

```shell
cmdr config list
```

Get a specific value:

```shell
cmdr config get -k <key>
```

### Set Configuration

```shell
cmdr config set -k <key> -v <value>
```

## Configuration Keys

### Core Settings

| Key | Default | Description |
|-----|---------|-------------|
| `core.root_dir` | `~/.cmdr` | Root directory for CMDR data[^1] |
| `core.bin_dir` | `bin` | Directory for command binaries (relative to root) |
| `core.shims_dir` | `shims` | Directory for shim scripts (relative to root) |
| `core.profile_dir` | `profile` | Directory for shell profile scripts (relative to root) |
| `core.database_path` | `cmdr.db` | Path to command database (relative to root) |
| `core.link_mode` | - | How to link commands (symlink vs copy) |

### Logging

| Key | Default | Description |
|-----|---------|-------------|
| `log.level` | `info` | Log level: `debug`, `info`, `warn`, `error`[^2] |
| `log.output` | `stderr` | Log output: `stdout` or `stderr` |

### Proxy Settings

| Key | Description |
|-----|-------------|
| `proxy.go` | Go module proxy URL |
| `proxy.http` | HTTP proxy URL |
| `proxy.https` | HTTPS proxy URL |

### Download Settings

| Key | Description |
|-----|-------------|
| `download.replace` | URL replacement pattern for proxying downloads |
| `download.direct.timeout` | Timeout for direct downloads |
| `download.direct.max_retries` | Maximum retry attempts |
| `download.proxy.enabled` | Enable download proxy |
| `download.proxy.type` | Proxy type |
| `download.proxy.address` | Proxy address |

## Environment Variables

All configuration keys can be set via environment variables using `CMDR_` prefix and replacing `.` with `_`[^3]:

```shell
# Set log level via environment
export CMDR_LOG_LEVEL=debug

# Set root directory
export CMDR_CORE_ROOT_DIR=/custom/path
```

## URL Replacement

The URL replacement feature allows you to redirect downloads through a proxy server. This is useful for:

- Speeding up downloads in regions with slow GitHub access
- Using internal mirrors

### Configuration

Set a replacement pattern using JSON:

```shell
cmdr config set -k download.replace -v '{"match": "^https://raw.githubusercontent.com/.*$", "template": "https://ghproxy.com/{{ .input | urlquery }}"}'
```

### Example Usage

```shell
# After configuring URL replacement
cmdr command install -n install.sh -v 0.0.0 -l https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh
# The URL is automatically rewritten to use the proxy
```

## Sample Configuration File

```yaml
# ~/.cmdr/config.yaml

core:
  root_dir: ~/.cmdr
  bin_dir: bin
  shims_dir: shims
  profile_dir: profile
  database_path: cmdr.db

log:
  level: info
  output: stderr

proxy:
  go: https://goproxy.cn,direct
  http: ""
  https: ""

download:
  replace:
    match: "^https://raw.githubusercontent.com/.*$"
    template: "https://ghproxy.com/{{ .input | urlquery }}"
```

---

[^1]: Configuration keys defined in [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L23-L33
[^2]: Log initialization in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L130-L145
[^3]: Environment variable binding in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L67-L71
