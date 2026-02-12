# Quick Start

This guide will help you get started with CMDR in just a few minutes.

## Installing Your First Command

Let's install a command version. CMDR can download from URLs or use local files:

```shell
# Install from URL
cmdr install -n <command-name> -v <version> -l <url-or-path>

# Example: Install a specific script
cmdr install -n hello -v 1.0.0 -l /usr/local/bin/hello
```

**Note:** The `cmdr command xxx` format is deprecated. Use `cmdr xxx` directly.

### Flags Explained

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-n` | `--name` | Name of the command |
| `-v` | `--version` | Version string (semantic versioning recommended) |
| `-l` | `--location` | URL or file path to command binary |

## Listing Installed Commands

View all versions of a specific command:

```shell
cmdr list -n <command-name>
```

View all managed commands:

```shell
cmdr list
```

## Switching Versions

Activate a specific version of a command:

```shell
cmdr use -n <command-name> -v <version>
```

After running this, the command will point to the specified version.

## Removing Commands

Remove a specific version:

```shell
cmdr remove -n <command-name> -v <version>
```

## Practical Example

Let's walk through managing multiple versions of a tool:

```shell
# 1. Install version 1.0.0
cmdr install -n mytool -v 1.0.0 -l https://example.com/mytool-1.0.0

# 2. Install version 2.0.0
cmdr install -n mytool -v 2.0.0 -l https://example.com/mytool-2.0.0

# 3. List all installed versions
cmdr list -n mytool

# 4. Use version 1.0.0
cmdr use -n mytool -v 1.0.0

# 5. Verify
mytool --version  # Should show 1.0.0

# 6. Switch to version 2.0.0
cmdr use -n mytool -v 2.0.0

# 7. Verify again
mytool --version  # Should show 2.0.0
```

## How Version Switching Works

When you run `cmdr use`, CMDR:

1. Looks up the command in its database[^1]
2. Creates or updates a shim script in `~/.cmdr/shims/`[^2]
3. The shim redirects to the actual binary location

Since `~/.cmdr/shims` is in your PATH, running the command name executes the shim, which runs the correct version.

## Next Steps

- [Configure CMDR](configuration.md) - Customize behavior with configuration
- [Architecture Overview](../architecture/overview.md) - Understand how CMDR works internally

---

[^1]: Command query in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L25-L34
[^2]: Shim management in [`core/strategy/`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/)
