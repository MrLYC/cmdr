# Design Patterns

This document describes the key design patterns used in CMDR.

## Factory Pattern

CMDR extensively uses the factory pattern for creating implementations of interfaces. This allows for:

- **Loose coupling** between interface consumers and implementations
- **Easy testing** with mock implementations
- **Runtime selection** of implementations

### Implementation

Each interface has a corresponding factory registry:

```go
// core/command.go
type factoryCommandManager func(cfg Configuration) (CommandManager, error)

var factoriesCommandManager map[CommandProvider]factoryCommandManager

func RegisterCommandManagerFactory(key CommandProvider, fn factoryCommandManager) {
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

### Registration

Implementations register themselves in `init()` functions[^1]:

```go
// core/manager/binary.go
func init() {
    core.RegisterCommandManagerFactory(core.CommandProviderBinary, NewBinaryManager)
}
```

### Usage Examples

| Interface | Provider Key | Implementation |
|-----------|--------------|----------------|
| `CommandManager` | `CommandProviderDatabase` | `core/manager/database.go` |
| `CommandManager` | `CommandProviderBinary` | `core/manager/binary.go` |
| `CommandManager` | `CommandProviderDownload` | `core/manager/download.go` |
| `Initializer` | `"filesystem"` | `core/initializer/filesystem.go` |
| `Initializer` | `"profile"` | `core/initializer/profile.go` |

## Strategy Pattern

CMDR uses the strategy pattern for shim creation, allowing different approaches to linking commands[^2]:

```
┌─────────────────────────────────────────────────────────┐
│                   Strategy Interface                     │
│  type Strategy interface {                              │
│      Execute(ctx context.Context, cmd Command) error    │
│  }                                                      │
└───────────────────────────┬─────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│    Direct     │   │     Proxy     │   │    Rewrite    │
│   Strategy    │   │   Strategy    │   │   Strategy    │
│               │   │               │   │               │
│ Creates       │   │ Creates shell │   │ Rewrites URLs │
│ symlinks      │   │ script proxy  │   │ for download  │
└───────────────┘   └───────────────┘   └───────────────┘
```

### Chain of Responsibility

The `ChainStrategy` combines multiple strategies[^3]:

```go
// core/strategy/chain.go
type ChainStrategy struct {
    strategies []Strategy
}

func (s *ChainStrategy) Execute(ctx context.Context, cmd Command) error {
    for _, strategy := range s.strategies {
        if err := strategy.Execute(ctx, cmd); err != nil {
            return err
        }
    }
    return nil
}
```

## Observer Pattern (Event Bus)

CMDR uses an event bus for decoupled communication between components[^4]:

```go
// core/bus.go
var bus = EventBus.New()

func PublishEvent(event string) {
    bus.Publish(event)
}

func SubscribeEvent(event string, handler func()) {
    bus.Subscribe(event, handler)
}
```

### Usage

```go
// Subscribe to exit event for cleanup
core.SubscribeEvent(core.EventExit, func() {
    // Cleanup resources
})

// Publish event when application exits
defer core.PublishEvent(core.EventExit)
```

## Decorator Pattern (Manager Chain)

Command managers can wrap each other to add functionality:

```
┌─────────────────────────────────────────────────────────┐
│                  DownloadManager                         │
│  - Downloads file from URL                              │
│  - Delegates to BinaryManager                           │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                   BinaryManager                          │
│  - Copies binary to bin directory                       │
│  - Sets permissions                                     │
│  - Delegates to DatabaseManager                         │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                  DatabaseManager                         │
│  - Persists command record                              │
│  - Core CRUD operations                                 │
└─────────────────────────────────────────────────────────┘
```

Each layer adds specific functionality while delegating core operations to the next layer.

## Dependency Injection

Configuration is injected through function parameters rather than global state:

```go
// Good: Configuration passed as parameter
func NewCommandManager(key CommandProvider, cfg Configuration) (CommandManager, error) {
    fn := factoriesCommandManager[key]
    return fn(cfg)  // cfg passed to factory
}

// Factory receives configuration
func NewDatabaseManager(cfg Configuration) (CommandManager, error) {
    dbPath := cfg.GetString(core.CfgKeyCmdrDatabasePath)
    // ...
}
```

This pattern makes components:

- **Testable** - Configuration can be mocked
- **Flexible** - Different configurations for different contexts
- **Explicit** - Dependencies are clear from function signatures

## Interface Segregation

CMDR defines focused interfaces[^5]:

```go
// Command - Read-only view of a command
type Command interface {
    GetName() string
    GetVersion() string
    GetActivated() bool
    GetLocation() string
}

// CommandQuery - Query builder for commands
type CommandQuery interface {
    WithName(name string) CommandQuery
    WithVersion(version string) CommandQuery
    WithActivated(activated bool) CommandQuery
    All() ([]Command, error)
    One() (Command, error)
    Count() (int, error)
}

// CommandManager - Full management operations
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

This separation allows consumers to depend only on the interface they need.

---

[^1]: Factory registration via init in [`core/manager/`](https://github.com/mrlyc/cmdr/blob/master/core/manager/)
[^2]: Strategy implementations in [`core/strategy/`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/)
[^3]: Chain strategy in [`core/strategy/chain.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/chain.go)
[^4]: Event bus in [`core/bus.go`](https://github.com/mrlyc/cmdr/blob/master/core/bus.go)
[^5]: Interface definitions in [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go) L18-L47
