# Core Interfaces

This document provides a reference for all core interfaces in CMDR.

## Command Interfaces

### Command

Read-only view of a managed command[^1]:

```go
type Command interface {
    GetName() string        // Command name (e.g., "kubectl")
    GetVersion() string     // Version string (e.g., "1.28.0")
    GetActivated() bool     // Whether this version is activated
    GetLocation() string    // File path to the binary
}
```

### CommandQuery

Fluent API for querying commands[^2]:

```go
type CommandQuery interface {
    WithName(name string) CommandQuery
    WithVersion(version string) CommandQuery
    WithActivated(activated bool) CommandQuery
    WithLocation(location string) CommandQuery
    
    All() ([]Command, error)   // Get all matching commands
    One() (Command, error)     // Get first matching command
    Count() (int, error)       // Count matching commands
}
```

**Usage Example:**

```go
query, _ := manager.Query()

// Find activated kubectl
cmd, _ := query.WithName("kubectl").WithActivated(true).One()

// List all versions of node
commands, _ := query.WithName("node").All()

// Count total managed commands
count, _ := query.Count()
```

### CommandManager

Full CRUD operations for commands[^3]:

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

**Providers:**

| Provider | Value | Implementation |
|----------|-------|----------------|
| Database | `CommandProviderDatabase` | [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go) |
| Binary | `CommandProviderBinary` | [`core/manager/binary.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/binary.go) |
| Download | `CommandProviderDownload` | [`core/manager/download.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/download.go) |
| Doctor | `CommandProviderDoctor` | [`core/manager/doctor.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/doctor.go) |

## Initialization Interfaces

### Initializer

System initialization tasks[^4]:

```go
type Initializer interface {
    Init(isUpgrade bool) error
}
```

**Registered Initializers:**

| Key | Implementation | Purpose |
|-----|----------------|---------|
| `"binary"` | `BinaryManager` | Create bin/ and shims/ directories |
| `"profile-dir-backup"` | `FSBackup` | Backup profile directory |
| `"profile-dir-export"` | `EmbedFSExporter` | Export embedded profile scripts |
| `"profile-dir-render"` | `DirRender` | Render profile templates |
| `"profile-injector"` | `ProfileInjector` | Inject source line into shell profile |
| `"database"` | `DatabaseInitializer` | Initialize database schema |
| `"command"` | `CommandInitializer` | Register cmdr as managed command |

## Download Interfaces

### Fetcher

File download abstraction[^5]:

```go
type Fetcher interface {
    Fetch(source, destination string) error
}
```

**Implementations:**

- `GoGetterFetcher` - Uses HashiCorp go-getter for flexible downloads
- `GoFetcher` - Specialized for Go module downloads

### DownloadStrategy

Strategy for URL transformationand retry logic:

```go
type DownloadStrategy interface {
    Name() string
    Prepare(uri string) (string, error)
    ShouldRetry(err error) bool
    ShouldFallback(err error) bool
    Configure(cfg Configuration) error
    IsEnabled(uri string) bool
    SetEnabled(enabled bool)
}
```

**Implementations:**

- `DirectStrategy` - Direct downloads
- `ProxyStrategy` - HTTP/SOCKS5 proxy
- `RewriteStrategy` - URL rewriting
- `ChainStrategy` - Chain of strategies with retry/fallback

## Configuration Interface

### Configuration

Type alias for Viper configuration[^6]:

```go
type Configuration = *viper.Viper
```

**Common Methods:**

```go
// Get values
GetString(key string) string
GetInt(key string) int
GetBool(key string) bool
GetStringSlice(key string) []string

// Set values
Set(key string, value interface{})
SetDefault(key string, value interface{})

// Check existence
IsSet(key string) bool

// Bind flags
BindPFlag(key string, flag *pflag.Flag) error
```

## Database Interface

### Database

Abstraction over Storm/BoltDB[^7]:

```go
type Database interface {
    Init() error
    Close() error
    
    Save(data interface{}) error
    Select(matchers ...q.Matcher) Query
    DeleteStruct(data interface{}) error
}
```

**Used By:** `DatabaseManager` for persisting command metadata.

## CMDR Searcher Interface

### CmdrSearcher

Search for CMDR releases (for self-upgrade)[^8]:

```go
type CmdrSearcher interface {
    GetReleaseAsset(ctx context.Context, releaseName, assetName string) (CmdrReleaseAsset, error)
}

type CmdrReleaseAsset struct {
    Name    string  // Release name
    Version string  // Version string
    Asset   string  // Asset name
    Url     string  // Download URL
}
```

**Providers:**

| Provider | Source |
|----------|--------|
| `CmdrSearcherProviderApi` | GitHub API |
| `CmdrSearcherProviderAtom` | GitHub Atom feed |

## Logger Interface

CMDR uses [logur](https://github.com/logur-dev/logur) for logging:

```go
type Logger interface {
    Trace(msg string, fields ...map[string]interface{})
    Debug(msg string, fields ...map[string]interface{})
    Info(msg string, fields ...map[string]interface{})
    Warn(msg string, fields ...map[string]interface{})
    Error(msg string, fields ...map[string]interface{})
}
```

**Usage:**

```go
logger := core.GetLogger()
logger.Info("command installed", map[string]interface{}{
    "name": "kubectl",
    "version": "1.28.0",
})
```

## Event Bus

Simple publish-subscribe for lifecycle events:

```go
func PublishEvent(event string)
func SubscribeEvent(event string, handler func())
```

**Events:**

- `EventExit` - Application exit

**Usage:**

```go
// Subscribe
core.SubscribeEvent(core.EventExit, func() {
    // Cleanup
})

// Publish
defer core.PublishEvent(core.EventExit)
```

---

[^1]: Command interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L18-L23
[^2]: CommandQuery interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L25-L34
[^3]: CommandManager interface in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L36-L47
[^4]: Initializer interface in [`core/initializer.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer.go) L5-L7
[^5]: Fetcher interface in [`core/fetcher.go`](https://github.com/mrlyc/cmdr/blob/master/core/fetcher.go)
[^6]: Configuration type in [`core/config.go`](https://github.com/mrlyc/cmdr/blob/master/core/config.go) L7
[^7]: Database interface in [`core/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/database.go)
[^8]: CmdrSearcher interface in [`core/cmdr.go`](https://github.com/mrlyc/cmdr/blob/master/core/cmdr.go) L18-L20
