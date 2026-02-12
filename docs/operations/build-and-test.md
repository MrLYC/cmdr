# Build & Test

This document describes how to build, test, and develop CMDR.

## Prerequisites

- Go 1.25 or later
- Make (optional, for convenience commands)
- Git

## Building from Source

### Clone Repository

```bash
git clone https://github.com/mrlyc/cmdr.git
cd cmdr
```

### Build Binary

```bash
# Using go build
go build -o cmdr main.go

# Or using make
make build
```

The binary will be created in the current directory.

### Install Locally

```bash
# Build and initialize
go build -o cmdr main.go
./cmdr init
```

## Running Tests

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Integration Tests

CMDR has integration tests in the `integration-tests/` directory:

```bash
cd integration-tests
./start.sh
```

**What Integration Tests Do:**

1. Build CMDR binary
2. Run `cmdr init`
3. Install test commands
4. Verify command activation
5. Test upgrade functionality
6. Clean up

**Environment Variables:**

| Variable | Description |
|----------|-------------|
| `SKIP_INSTALL` | Skip installation step (for CI) |
| `CMDR_CORE_ROOT_DIR` | Override CMDR root directory |
| `CMDR_CORE_PROFILE_DIR` | Override profile directory |

## Development Workflow

### Project Structure

```
cmdr/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command & initialization
│   ├── command/           # command subcommands
│   ├── config/            # config subcommands
│   ├── init.go            # init command
│   ├── upgrade.go         # upgrade command
│   └── ...
├── core/                   # Core business logic
│   ├── command.go         # Command interface
│   ├── initializer.go     # Initializer interface
│   ├── config.go          # Configuration keys
│   ├── manager/           # CommandManager implementations
│   ├── initializer/       # Initializer implementations
│   ├── strategy/          # Download strategies
│   ├── fetcher/           # File fetchers
│   └── ...
├── integration-tests/      # Integration test scripts
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── Makefile               # Build automation
└── README.md
```

### Code Generation

CMDR uses go:generate for code generation:

```bash
# Generate string methods for enums
go generate ./...
```

Generated files:
- `*_string.go` - String methods for enum types (via stringer)
- `mock/*.go` - Mock implementations (via mockgen)

### Adding a New Command

1. Create command file in `cmd/` or `cmd/<subcommand>/`:

```go
// cmd/mycommand.go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/mrlyc/cmdr/core"
)

var myCmd = &cobra.Command{
    Use:   "my-command",
    Short: "Description",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    
    // Add flags
    myCmd.Flags().StringP("flag", "f", "", "flag description")
    
    // Bind to configuration
    cfg := core.GetConfiguration()
    cfg.BindPFlag("my.flag", myCmd.Flags().Lookup("flag"))
}
```

2. Register with root command in `init()`

3. Add tests in `cmd/*_test.go`

### Adding a New Manager

1. Create manager file in `core/manager/`:

```go
// core/manager/custom.go
package manager

import "github.com/mrlyc/cmdr/core"

type CustomManager struct {
    manager core.CommandManager // Wrapped manager
}

func (m *CustomManager) Define(name, version, location string) (core.Command, error) {
    // Custom logic
    return m.manager.Define(name, version, location)
}

// Implement other CommandManager methods...

func init() {
    core.RegisterCommandManagerFactory(core.CommandProviderCustom, func(cfg core.Configuration) (core.CommandManager, error) {
        // Get wrapped manager
        wrapped, _ := core.NewCommandManager(core.CommandProviderBinary, cfg)
        return &CustomManager{manager: wrapped}, nil
    })
}
```

2. Add provider constant in `core/command.go`

3. Add tests in `core/manager/custom_test.go`

## Testing Guidelines

### Unit Test Structure

```go
package manager_test

import (
    "testing"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
)

func TestManager(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Manager Suite")
}

var _ = ginkgo.Describe("CustomManager", func() {
    var manager *manager.CustomManager
    
    ginkgo.BeforeEach(func() {
        cfg := core.NewConfiguration()
        manager, _ = manager.NewCustomManager(cfg)
    })
    
    ginkgo.AfterEach(func() {
        manager.Close()
    })
    
    ginkgo.It("should define command", func() {
        cmd, err := manager.Define("test", "1.0.0", "/path/to/bin")
        gomega.Expect(err).ToNot(gomega.HaveOccurred())
        gomega.Expect(cmd.GetName()).To(gomega.Equal("test"))
    })
})
```

### Running Specific Tests

```bash
# Run specific package
go test ./core/manager/...

# Run specific test
go test ./core/manager/ -run TestDatabaseManager

# Run with verbose output
go test ./... -v

# Run with race detector
go test ./... -race
```

## Debugging

### Enable Debug Logging

```bash
# Via environment variable
export CMDR_LOG_LEVEL=debug
cmdr list

# Via config flag
cmdr --config /path/to/debug-config.yaml command list
```

### Inspect Database

```bash
# Install BoltDB browser
go install github.com/br0xen/boltbrowser@latest

# Open database
boltbrowser ~/.cmdr/cmdr.db
```

### Check Directory Structure

```bash
# List all managed files
tree ~/.cmdr

# Check symlinks
ls -la ~/.cmdr/bin/
ls -la ~/.cmdr/shims/
```

## CI/CD

### GitHub Actions Workflows

CMDR uses GitHub Actions for CI/CD:

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `unittest.yml` | Push, PR | Run unit tests on multiple Go versions |
| `integration-test.yml` | Push, PR | Run integration tests |
| `release.yml` | Tag push | Build and publish releases |
| `codeql-analysis.yml` | Schedule, PR | Security analysis |

**Workflow Files:** `.github/workflows/`

### Running CI Locally

```bash
# Unit tests (matches CI)
go test ./... -race -coverprofile=coverage.out

# Integration tests (matches CI)
export SKIP_INSTALL=true
export CMDR_CORE_ROOT_DIR=$PWD/.cmdr-test
export CMDR_CORE_PROFILE_DIR=$PWD/.cmdr-test/profile
cd integration-tests && ./start.sh
```

## Performance Profiling

### CPU Profiling

```bash
# Build with profiling
go test -cpuprofile cpu.prof -bench=. ./core/manager/

# Analyze profile
go tool pprof cpu.prof
```

### Memory Profiling

```bash
# Build with profiling
go test -memprofile mem.prof -bench=. ./core/manager/

# Analyze profile
go tool pprof mem.prof
```

## Code Quality

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run ./...
```

### Formatting

```bash
# Format code
go fmt ./...

# Or using gofmt
gofmt -w .
```

### Vet

```bash
# Static analysis
go vet ./...
```
