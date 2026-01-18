# Command Manager

The Command Manager is the core abstraction for managing command versions in CMDR. It provides a layered architecture with different providers for different concerns.

## Overview

Command managers follow a chain-of-responsibility pattern where each manager wraps another, adding specific functionality:

```
DownloadManager → BinaryManager → DatabaseManager
```

Each layer delegates to the next while adding its own behavior.

## CommandManager Interface

The core interface defines operations for managing commands[^1]:

```go
type CommandManager interface {
    Close() error
    Provider() CommandProvider
    Query() (CommandQuery, error)
    Define(name, version, location string) (Command, error)
    Undefine(name, version string) error
    Activate(name, version string) error
    Deactivate(name string) error
}
```

## Manager Implementations

### DatabaseManager

**Source:** [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go)

**Provider:** `CommandProviderDatabase`

**Responsibilities:**

- Persist command metadata in BoltDB via Storm ORM
- Query commands by name, version, location, activation status
- Wrap BinaryManager for actual file operations

**Key Operations:**

```go
// Define a command - saves to database and delegates to BinaryManager
func (m *DatabaseManager) Define(name, version, location string) (Command, error)

// Query commands with flexible filtering
func (m *DatabaseManager) Query() (CommandQuery, error)

// Activate - deactivates others, updates DB, delegates to BinaryManager
func (m *DatabaseManager) Activate(name, version string) error
```

**Implementation Details:**

- Uses version matching to support both raw and semantic versions[^2]
- Enforces single activation per command name
- Prevents deletion of activated commands[^3]

### BinaryManager

**Source:** [`core/manager/binary.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/binary.go)

**Provider:** `CommandProviderBinary`

**Responsibilities:**

- Manage physical binary files in `~/.cmdr/shims/<command>/`
- Create symlinks or copies of binaries
- Handle version normalization (e.g., `1.0` → `1.0.0`)
- Create activation symlinks in `bin/` directory

**Directory Structure:**

```
~/.cmdr/
├── bin/
│   └── kubectl -> ../shims/kubectl/kubectl_1.28.0
└── shims/
    └── kubectl/
        ├── kubectl_1.28.0
        └── kubectl_1.29.0
```

**Key Operations:**

```go
// Define - copy or link binary to shims directory
func (m *BinaryManager) Define(name, version, location string) (Command, error)

// Activate - create symlink in bin/ directory
func (m *BinaryManager) Activate(name, version string) error

// Deactivate - remove symlink from bin/ directory
func (m *BinaryManager) Deactivate(name string) error
```

**Link Modes:**

The manager supports two link modes[^4]:

| Mode | Behavior | Use Case |
|------|----------|----------|
| `copy` (default) | Copy binary to shims directory | Maximum compatibility |
| `link` | Create symlink to original location | Save disk space |

Set via configuration:

```shell
cmdr config set -k core.link_mode -v link
```

### DownloadManager

**Source:** [`core/manager/download.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/download.go)

**Provider:** `CommandProviderDownload`

**Responsibilities:**

- Download binaries from URLs
- Validate downloaded files
- Apply URL rewrite rules
- Delegate to BinaryManager for storage

**Key Operations:**

```go
// Define - download from URL, then delegate to BinaryManager
func (m *DownloadManager) Define(name, version, location string) (Command, error)
```

**URL Rewriting:**

Supports URL replacement patterns for proxying downloads:

```go
// Configuration example
{
    "match": "^https://raw.githubusercontent.com/.*$",
    "template": "https://ghproxy.com/{{ .input | urlquery }}"
}
```

### DoctorManager

**Source:** [`core/manager/doctor.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/doctor.go)

**Provider:** `CommandProviderDoctor`

**Responsibilities:**

- Diagnose installation issues
- Verify directory structure
- Check database integrity
- Validate PATH configuration

## Command Interface

Commands are read-only views of managed commands[^5]:

```go
type Command interface {
    GetName() string
    GetVersion() string
    GetActivated() bool
    GetLocation() string
}
```

## CommandQuery Interface

Provides fluent API for filtering commands[^6]:

```go
type CommandQuery interface {
    WithName(name string) CommandQuery
    WithVersion(version string) CommandQuery
    WithActivated(activated bool) CommandQuery
    WithLocation(location string) CommandQuery
    
    All() ([]Command, error)
    One() (Command, error)
    Count() (int, error)
}
```

**Example Usage:**

```go
// Find activated kubectl command
query, _ := manager.Query()
cmd, _ := query.WithName("kubectl").WithActivated(true).One()

// List all versions of kubectl
query, _ := manager.Query()
commands, _ := query.WithName("kubectl").All()
```

## Factory Registration

Managers register themselves via the factory pattern[^7]:

```go
// In init() function
core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(cfg Configuration) (CommandManager, error) {
    // Create BinaryManager first
    mgr, _ := core.NewCommandManager(core.CommandProviderBinary, cfg)
    db, _ := core.GetDatabase()
    return NewDatabaseManager(db, mgr), nil
})
```

## Version Handling

CMDR normalizes versions using semantic versioning:

- Input: `1.0`, `v1.0.0`, `1.0.0`
- Normalized: `1.0.0`

This ensures consistent storage and lookup across different version formats.

---

[^1]: CommandManager interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L36-L47
[^2]: Version matching in [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go) L12-L18
[^3]: Activated command protection in [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go) L104-L106
[^4]: Link mode configuration in [`core/manager/binary.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/binary.go) L415-L425
[^5]: Command interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L18-L23
[^6]: CommandQuery interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L25-L34
[^7]: Factory registration in [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go) L182-L196
