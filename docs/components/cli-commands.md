# CLI Commands

This document details all CLI commands available in CMDR.

## Root Command

```shell
cmdr [flags] [command]
```

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Path to config file (default: `~/.cmdr/config.yaml`) |
| `--help` | `-h` | Help for cmdr |

## Command Management

### `cmdr install`

Install a new command version into CMDR.

```shell
cmdr install -n <name> -v <version> -l <location> [-a]
```

**Note:** The old format `cmdr command install` is deprecated. Use `cmdr install` instead.

**Flags:**

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--name` | `-n` | Yes | Command name |
| `--version` | `-v` | Yes | Version string |
| `--location` | `-l` | Yes | URL or file path to the binary |
| `--activate` | `-a` | No | Activate immediately after install |

**Source:** [`cmd/command/install.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/install.go)[^1]

**Example:**

```shell
# Install kubectl 1.28.0 from URL
cmdr install -n kubectl -v 1.28.0 -l https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl

# Install and activate immediately
cmdr install -n kubectl -v 1.28.0 -l /path/to/kubectl -a
```

### `cmdr use`

Activate a specific version of a command.

```shell
cmdr use -n <name> -v <version>
```

**Note:** The old format `cmdr command use` is deprecated. Use `cmdr use` instead.

**Flags:**

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--name` | `-n` | Yes | Command name |
| `--version` | `-v` | Yes | Version to activate |

**Source:** [`cmd/command/use.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/use.go)[^2]

**Example:**

```shell
cmdr use -n kubectl -v 1.28.0
```

### `cmdr list`

List installed command versions.

```shell
cmdr list [-n <name>] [-v <version>] [-a]
```

**Note:** The old format `cmdr command list` is deprecated. Use `cmdr list` instead.

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--name` | `-n` | Filter by command name |
| `--version` | `-v` | Filter by version |
| `--activate` | `-a` | Show only activated commands |

**Source:** [`cmd/command/list.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/list.go)

**Example:**

```shell
# List all commands
cmdr list

# List all versions of kubectl
cmdr list -n kubectl

# List only activated commands
cmdr list -a
```

### `cmdr remove`

Remove a command version.

```shell
cmdr remove -n <name> -v <version>
```

**Note:** The old format `cmdr command remove` is deprecated. Use `cmdr remove` instead.

**Flags:**

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--name` | `-n` | Yes | Command name |
| `--version` | `-v` | Yes | Version to remove |

**Source:** [`cmd/command/remove.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/remove.go)

### `cmdr unset`

Deactivate a command (remove from shims).

```shell
cmdr unset -n <name>
```

**Note:** The old format `cmdr command unset` is deprecated. Use `cmdr unset` instead.

**Source:** [`cmd/command/unset.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/unset.go)

### `cmdr define`

Define a command without downloading (for local binaries).

```shell
cmdr define -n <name> -v <version> -l <path> [-a]
```

**Note:** The old format `cmdr command define` is deprecated. Use `cmdr define` instead.

**Source:** [`cmd/command/define.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/define.go)

## Configuration Management

### `cmdr config list`

List all configuration values.

```shell
cmdr config list
```

**Source:** [`cmd/config/list.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/config/list.go)

### `cmdr config get`

Get a specific configuration value.

```shell
cmdr config get -k <key>
```

**Source:** [`cmd/config/get.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/config/get.go)

### `cmdr config set`

Set a configuration value.

```shell
cmdr config set -k <key> -v <value>
```

**Source:** [`cmd/config/set.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/config/set.go)

**Example:**

```shell
# Set log level
cmdr config set -k log.level -v debug

# Set URL replacement for proxy
cmdr config set -k download.replace -v '{"match": "...", "template": "..."}'
```

## System Commands

### `cmdr clean`

Clean old **inactive** command versions.

This command:
- Groups versions by command name
- Sorts **inactive** versions by added time (based on shim file mtime)
- Keeps the newest inactive versions (default: 3)
- Moves versions older than a threshold (default: 100 days) into a trash directory

```shell
cmdr clean [-n <name> ...] [--age <days>] [--keep <count>]
```

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--name` / `-n` | No | Only clean specified command name(s). Can be repeated. If any specified name does not exist, the command fails. |
| `--age` | No | Only clean versions older than this many days (default: 100) |
| `--keep` | No | Keep this many newest inactive versions per command (default: 3) |

**Trash directory:**
- macOS: `~/.Trash/cmdr-cleaned`
- Linux: `/tmp/cmdr-cleaned`

**Source:** [`cmd/clean.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/clean.go)

### `cmdr init`

Initialize CMDR environment.

```shell
cmdr init
```

This command:
1. Creates directory structure (`bin/`, `shims/`, `profile/`)
2. Generates shell profile scripts
3. Registers CMDR itself as a managed command

**Source:** [`cmd/init.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/init.go)

### `cmdr upgrade`

Upgrade CMDR to the latest version.

```shell
cmdr upgrade [-r <release>] [-a <asset>]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--release` | `-r` | Specific release name (default: latest) |
| `--asset` | `-a` | Specific asset name |

**Source:** [`cmd/upgrade.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/upgrade.go)

### `cmdr doctor`

Diagnose CMDR installation issues.

```shell
cmdr doctor
```

Checks for:
- Directory structure integrity
- Database accessibility
- Shell profile configuration
- PATH configuration

**Source:** [`cmd/doctor.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/doctor.go)

### `cmdr version`

Display CMDR version information.

```shell
cmdr version
```

**Source:** [`cmd/version.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/version.go)

## Command Structure

```
cmdr
├── clean         # Clean old inactive versions
├── command
│   ├── define    # Define command from local path
│   ├── install   # Install command from URL/path
│   ├── list      # List installed commands
│   ├── remove    # Remove a command version
│   ├── unset     # Deactivate a command
│   └── use       # Activate a command version
├── config
│   ├── get       # Get config value
│   ├── list      # List all config
│   └── set       # Set config value
├── doctor        # Diagnose issues
├── init          # Initialize CMDR
├── upgrade       # Upgrade CMDR
└── version       # Show version
```

---

[^1]: Install command implementation in [`cmd/command/install.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/install.go) L12-L37
[^2]: Use command implementation in [`cmd/command/use.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/use.go) L12-L32
