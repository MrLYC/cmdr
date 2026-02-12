# Architecture Overview

This document provides a high-level overview of CMDR's architecture and design principles.

## System Overview

CMDR is a command-line tool built in Go that manages multiple versions of CLI commands. It uses a layered architecture with clear separation between:

1. **CLI Layer** (`cmd/`) - User-facing commands built with Cobra
2. **Core Layer** (`core/`) - Business logic, interfaces, and abstractions
3. **Storage Layer** - BoltDB-based persistence via Storm ORM

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Layer (cmd/)                        │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌────────┐ │
│  │  init   │ │ command │ │ config  │ │ upgrade │ │ doctor │ │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └───┬────┘ │
└───────┼──────────┼─────────┼─────────┼──────────┼───────────┘
        │          │         │         │          │
        ▼          ▼         ▼         ▼          ▼
┌─────────────────────────────────────────────────────────────┐
│                     Core Layer (core/)                       │
│  ┌───────────────┐  ┌─────────────┐  ┌──────────────────┐   │
│  │ CommandManager│  │ Initializer │  │    Strategies    │   │
│  │   Interface   │  │  Interface  │  │ (shim creation)  │   │
│  └───────┬───────┘  └──────┬──────┘  └────────┬─────────┘   │
│          │                 │                   │             │
│  ┌───────┴───────┐  ┌──────┴──────┐  ┌────────┴─────────┐   │
│  │   Managers    │  │ Initializers│  │    Fetchers      │   │
│  │ (binary,db,..)│  │(profile,fs) │  │  (go-getter,go)  │   │
│  └───────────────┘  └─────────────┘  └──────────────────┘   │
└───────────────────────────────────────────────────────────────┘
        │                    │
        ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                    Storage Layer                             │
│  ┌──────────────────┐     ┌─────────────────────────────┐   │
│  │   Storm/BoltDB   │     │     Filesystem (shims,      │   │
│  │   (cmdr.db)      │     │      binaries, profile)     │   │
│  └──────────────────┘     └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Key Components

### Entry Point

The application starts in `main.go`, which:

1. Sets up panic recovery with custom exit codes[^1]
2. Publishes an exit event for cleanup
3. Delegates to the CLI layer

```go
// main.go L14-L33
func main() {
    defer func() {
        recovered := recover()
        // Handle ExitError for clean exit codes
    }()
    defer core.PublishEvent(core.EventExit)
    ctx := context.Background()
    cmd.ExecuteContext(ctx)
}
```

### CLI Layer (`cmd/`)

The CLI is built with [Cobra](https://github.com/spf13/cobra) and organized into subcommands:

| Command | Description | Source |
|---------|-------------|--------|
| `cmdr init` | Initialize CMDR environment | [`cmd/init.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/init.go) |
| `cmdr install/use/list/...` | Manage command versions (deprecated: `cmdr command xxx`) | [`cmd/command/`](https://github.com/mrlyc/cmdr/blob/master/cmd/command/) |
| `cmdr config` | Manage configuration | [`cmd/config/`](https://github.com/mrlyc/cmdr/blob/master/cmd/config/) |
| `cmdr upgrade` | Upgrade CMDR itself | [`cmd/upgrade.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/upgrade.go) |
| `cmdr doctor` | Diagnose issues | [`cmd/doctor.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/doctor.go) |
| `cmdr version` | Show version info | [`cmd/version.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/version.go) |

### Core Layer (`core/`)

The core layer defines interfaces and provides implementations:

#### Interfaces

- **`CommandManager`** - CRUD operations for commands[^2]
- **`Initializer`** - System initialization tasks[^3]
- **`CmdrSearcher`** - Finding CMDR releases[^4]
- **`Database`** - Data persistence abstraction[^5]

#### Implementations

Located in subdirectories:

- **`core/manager/`** - Command manager implementations (binary, database, download, doctor)
- **`core/initializer/`** - Initialization implementations (filesystem, profile, command)
- **`core/fetcher/`** - File download implementations (go-getter, go)
- **`core/strategy/`** - Shim creation strategies (direct, proxy, chain, rewrite)

### Configuration

Configuration is managed by Viper with support for:

- YAML configuration files
- Environment variables (prefixed with `CMDR_`)
- Command-line flags

Configuration initialization happens in `cmd/root.go`:

1. `preInitConfig()` - Set defaults and env binding[^6]
2. `initConfig()` - Load config file[^7]
3. `postInitConfig()` - Resolve relative paths[^8]

## Factory Pattern

CMDR uses the factory pattern extensively for extensibility:

```go
// Example from core/command.go
var factoriesCommandManager map[CommandProvider]func(cfg Configuration) (CommandManager, error)

func RegisterCommandManagerFactory(key CommandProvider, fn func(...) (CommandManager, error)) {
    factoriesCommandManager[key] = fn
}

func NewCommandManager(key CommandProvider, cfg Configuration) (CommandManager, error) {
    fn, ok := factoriesCommandManager[key]
    if !ok {
        return nil, ErrCommandManagerFactoryeNotFound
    }
    return fn(cfg)
}
```

This pattern allows:

- Decoupled implementation registration
- Easy testing with mock implementations
- Runtime provider selection

## Event System

CMDR uses an event bus for decoupled communication:

```go
// core/bus.go
func PublishEvent(event string)
func SubscribeEvent(event string, handler func())
```

Events like `EventExit` allow components to clean up resources without tight coupling.

---

[^1]: Exit handling in [`main.go`](https://github.com/mrlyc/cmdr/blob/master/main.go) L15-L28
[^2]: CommandManager interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L36-L47
[^3]: Initializer interface in [`core/initializer.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer.go) L5-L7
[^4]: CmdrSearcher interface in [`core/cmdr.go`](https://github.com/mrlyc/cmdr/blob/master/core/cmdr.go) L18-L20
[^5]: Database interface in [`core/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/database.go)
[^6]: Pre-initialization in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L67-L89
[^7]: Config loading in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L92-L103
[^8]: Path resolution in [`cmd/root.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/root.go) L105-L128
