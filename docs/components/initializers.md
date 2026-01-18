# Initializers

Initializers are responsible for setting up the CMDR environment during the `cmdr init` command. They configure the filesystem, shell profiles, and other necessary components.

## Initializer Interface

The core interface is simple[^1]:

```go
type Initializer interface {
    Init(isUpgrade bool) error
}
```

The `isUpgrade` parameter indicates whether this is a fresh installation or an upgrade of an existing installation.

## Factory Registration

Initializers register themselves via factory pattern[^2]:

```go
core.RegisterInitializerFactory("filesystem", func(cfg Configuration) (Initializer, error) {
    return NewFilesystemInitializer(cfg), nil
})
```

## Initializer Implementations

### Filesystem Initializers

**Source:** [`core/initializer/filesystem.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/filesystem.go)

#### FSBackup

Backs up existing profile directory before modifications[^3]:

```go
type FSBackup struct {
    path   string  // Directory to backup
    target string  // Backup location
}
```

**Usage:** Registered as `"profile-dir-backup"`

#### EmbedFSExporter

Exports embedded filesystem to disk[^4]:

```go
type EmbedFSExporter struct {
    filesystem fs.FS       // Embedded filesystem
    srcPath    string      // Source path in embed
    dstPath    string      // Destination path on disk
    fileMode   os.FileMode // File permissions
}
```

**Usage:** Registered as `"profile-dir-export"` - exports shell profile scripts from embedded assets.

#### DirRender

Renders Go templates in a directory[^5]:

```go
type DirRender struct {
    data    interface{} // Template data
    srcPath string      // Directory to render
    ext     string      // Template extension (e.g., ".gotmpl")
}
```

**Usage:** Registered as `"profile-dir-render"` - renders profile templates with configuration.

### Profile Initializer

**Source:** [`core/initializer/profile.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/profile.go)

#### ProfileInjector

Injects CMDR initialization script into shell profile[^6]:

```go
type ProfileInjector struct {
    scriptPath  string // Path to cmdr_initializer.sh
    profilePath string // Path to shell profile (~/.bashrc, ~/.zshrc, etc.)
}
```

**Key Operations:**

1. **Detect Shell Profile**[^7]:
   - Bash → `~/.bashrc`
   - Zsh → `~/.zshrc`
   - Fish → `~/.config/fish/config.fish`
   - Ash/Sh → `~/.profile`

2. **Generate Profile Statement**:
   ```bash
   source '/path/to/.cmdr/profile/cmdr_initializer.sh'
   ```

3. **Update Profile**:
   - Reads existing profile
   - Replaces old CMDR source line (if exists)
   - Adds source line (if not exists)
   - Writes updated profile

**Shell Detection:** Uses parent process inspection to determine the current shell.

### Command Initializer

**Source:** [`core/initializer/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/command.go)

Registers CMDR itself as a managed command, enabling `cmdr upgrade`.

### Database Initializer

**Source:** [`core/initializer/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/database.go)

Initializes the BoltDB database schema.

## Initialization Flow

When `cmdr init` is run, initializers execute in this order:

```
1. profile-dir-backup
   └─> Backup existing profile directory

2. binary (BinaryManager.Init)
   └─> Create bin/ and shims/ directories

3. profile-dir-export
   └─> Export embedded profile scripts

4. profile-dir-render
   └─> Render profile templates with config

5. profile-injector
   └─> Add source line to shell profile

6. database
   └─> Initialize database schema

7. command
   └─> Register cmdr as managed command
```

## Embedded Assets

CMDR embeds shell profile scripts using Go 1.16+ embed:

```go
//go:embed embed/*
var EmbedFS embed.FS
```

The embedded profile directory contains:

```
embed/
└── profile/
    ├── cmdr_initializer.sh.gotmpl  # Shell initializer template
    └── ...
```

These are extracted to `~/.cmdr/profile/` and rendered with configuration values.

## Template Rendering

Templates use Go's `text/template` package with configuration data:

```go
// Example template
// cmdr_initializer.sh.gotmpl
export CMDR_ROOT="{{.Configuration.GetString "core.root_dir"}}"
export PATH="{{.Configuration.GetString "core.bin_dir"}}:$PATH"
```

Rendered output:

```bash
export CMDR_ROOT="/home/user/.cmdr"
export PATH="/home/user/.cmdr/bin:$PATH"
```

## Safety Features

1. **Backup Before Modify**: Profile directory is backed up before changes
2. **Idempotent Injection**: Multiple `cmdr init` calls don't duplicate source lines
3. **Version Migration**: Old format source lines are updated to new format
4. **Error Recovery**: Failures during init don't leave system in broken state

## Testing Initializers

Initializers can be tested independently:

```go
// Example test
cfg := core.NewConfiguration()
cfg.Set(core.CfgKeyCmdrProfileDir, "/tmp/test-profile")

injector := NewProfileInjector("/tmp/test-profile/init.sh", "/tmp/.bashrc")
err := injector.Init(false)
```

---

[^1]: Initializer interface in [`core/initializer.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer.go) L5-L7
[^2]: Factory registration pattern in [`core/initializer.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer.go) L18-L20
[^3]: FSBackup implementation in [`core/initializer/filesystem.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/filesystem.go) L21-L61
[^4]: EmbedFSExporter implementation in [`core/initializer/filesystem.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/filesystem.go) L63-L146
[^5]: DirRender implementation in [`core/initializer/filesystem.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/filesystem.go) L148-L264
[^6]: ProfileInjector implementation in [`core/initializer/profile.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/profile.go) L20-L118
[^7]: Shell detection in [`core/initializer/profile.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer/profile.go) L120-L151
