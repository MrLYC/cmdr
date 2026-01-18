# Data Flow

This document describes how data flows through CMDR during common operations.

## Command Installation Flow

When a user runs `cmdr command install -n kubectl -v 1.28.0 -l <url>`:

```
┌──────────────┐
│    User      │
│  (CLI input) │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────┐
│                     cmd/command/install.go                    │
│  1. Parse flags (name, version, location)                     │
│  2. Create DownloadManager from factory                       │
│  3. Call manager.Define(name, version, location)              │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                  core/manager/download.go                     │
│  1. Fetch file from URL using Fetcher                         │
│  2. Validate downloaded binary                                │
│  3. Delegate to BinaryManager                                 │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                   core/manager/binary.go                      │
│  1. Copy binary to bin directory                              │
│  2. Set executable permissions                                │
│  3. Delegate to DatabaseManager                               │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                  core/manager/database.go                     │
│  1. Create Command record                                     │
│  2. Store in Storm/BoltDB                                     │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                      ~/.cmdr/cmdr.db                          │
│  Command: {name: "kubectl", version: "1.28.0", ...}           │
└──────────────────────────────────────────────────────────────┘
```

## Command Activation Flow

When a user runs `cmdr command use -n kubectl -v 1.28.0`:

```
┌──────────────┐
│    User      │
│  (CLI input) │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────┐
│                     cmd/command/use.go                        │
│  1. Parse flags (name, version)                               │
│  2. Query database for command                                │
│  3. Activate the command version                              │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                  core/manager/command.go                      │
│  1. Deactivate previous version (if any)                      │
│  2. Update database: set activated=true                       │
│  3. Create/update shim in shims directory                     │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                    core/strategy/                             │
│  Strategy pattern for shim creation:                          │
│  - DirectStrategy: Symlink to binary                          │
│  - ProxyStrategy: Script that execs binary                    │
│  - ChainStrategy: Multiple strategies in sequence             │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                 ~/.cmdr/shims/kubectl                         │
│  (Shim script or symlink pointing to active version)          │
└──────────────────────────────────────────────────────────────┘
```

## Command Execution Flow

When a user runs a managed command (e.g., `kubectl get pods`):

```
┌──────────────┐
│    User      │
│  $ kubectl   │
└──────┬───────┘
       │
       │ (PATH includes ~/.cmdr/shims)
       ▼
┌──────────────────────────────────────────────────────────────┐
│              ~/.cmdr/shims/kubectl                            │
│  (Shim created by CMDR)                                       │
│  Redirects to actual binary location                          │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│    ~/.cmdr/bin/kubectl-1.28.0                                 │
│    (Actual binary for the activated version)                  │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                     Command Output                            │
│    (Normal command execution)                                 │
└──────────────────────────────────────────────────────────────┘
```

## Initialization Flow

When `cmdr init` is run:

```
┌──────────────┐
│    User      │
│  cmdr init   │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────┐
│                      cmd/init.go                              │
│  1. Iterate through registered initializers                   │
│  2. Call Init(isUpgrade) on each                              │
└───────────────────────────┬──────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  Filesystem   │   │    Profile    │   │   Command     │
│  Initializer  │   │  Initializer  │   │  Initializer  │
│               │   │               │   │               │
│ Creates dirs: │   │ Creates shell │   │ Registers     │
│ - bin/        │   │ profile for   │   │ cmdr as a     │
│ - shims/      │   │ PATH setup    │   │ managed cmd   │
│ - profile/    │   │               │   │               │
└───────────────┘   └───────────────┘   └───────────────┘
```

## Database Schema

Commands are stored in BoltDB via Storm ORM:

```go
type CommandModel struct {
    ID        int    `storm:"id,increment"`
    Name      string `storm:"index"`
    Version   string
    Location  string
    Activated bool   `storm:"index"`
}
```

The database file is located at `~/.cmdr/cmdr.db` by default.

## Configuration Flow

```
┌─────────────────────┐
│  Environment Vars   │
│  (CMDR_*)           │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Config File        │◄──── ~/.cmdr/config.yaml
│  (Viper)            │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  CLI Flags          │◄──── --config, etc.
│  (Cobra)            │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Final Config       │
│  (merged values)    │
└─────────────────────┘
```

Priority (highest to lowest):
1. CLI flags
2. Environment variables
3. Config file
4. Default values
