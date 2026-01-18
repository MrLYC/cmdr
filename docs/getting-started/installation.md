# Installation

This guide covers how to install CMDR on your system.

## Prerequisites

- A Unix-like operating system (Linux, macOS)
- Bash, Zsh, or compatible shell
- `curl` for script installation

## Quick Install (Recommended)

Run the following command to install the latest version of CMDR:

```shell
curl -o- https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh | ${SHELL:-bash}
```

### For Users in China

Speed up the installation using the proxy option:

```shell
curl -o- https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh | bash -s -p
```

## Manual Installation

If you prefer manual installation or need more control:

1. **Download the latest release** from [GitHub Releases](https://github.com/mrlyc/cmdr/releases/latest)

2. **Make the binary executable**:
   ```shell
   chmod +x /path/to/cmdr
   ```

3. **Initialize CMDR**:
   ```shell
   /path/to/cmdr init
   ```

4. **Restart your shell** and verify:
   ```shell
   cmdr version
   ```

## What the Installer Does

The installation script performs these steps[^1]:

1. Detects your OS and architecture
2. Downloads the appropriate binary from GitHub releases
3. Places the binary in a configured location
4. Runs `cmdr init` to set up the environment
5. Configures shell profile for PATH integration

## Directory Structure

After installation, CMDR creates the following directory structure[^2]:

```
~/.cmdr/
├── bin/           # Managed command binaries
├── shims/         # Shim scripts for version switching
├── profile/       # Shell initialization scripts
├── config.yaml    # Configuration file
└── cmdr.db        # Command database
```

## Verifying Installation

After installation, verify everything is working:

```shell
# Check CMDR version
cmdr version

# Check system status
cmdr doctor
```

## Uninstalling

To uninstall CMDR:

1. Remove the CMDR directory:
   ```shell
   rm -rf ~/.cmdr
   ```

2. Remove the CMDR source line from your shell profile (`.bashrc`, `.zshrc`, etc.)

---

[^1]: Installation script: [`install.sh`](https://github.com/mrlyc/cmdr/blob/master/install.sh)
[^2]: Default paths configured in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L73-L77
