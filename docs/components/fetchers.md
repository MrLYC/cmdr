# Fetchers

Fetchers are responsible for downloading files from various sources (URLs, local paths, Git repositories, etc.).

## Fetcher Interface

The core interface is minimal[^1]:

```go
type Fetcher interface {
    Fetch(source, destination string) error
}
```

## Fetcher Implementations

### GoGetterFetcher

**Source:** [`core/fetcher/go_getter.go`](https://github.com/mrlyc/cmdr/blob/master/core/fetcher/go_getter.go)

Uses HashiCorp's [go-getter](https://github.com/hashicorp/go-getter) library for flexible URL-based downloads.

**Supported Sources:**

| Protocol | Example | Description |
|----------|---------|-------------|
| HTTP/HTTPS | `https://example.com/file.tar.gz` | Direct downloads |
| Git | `git::https://github.com/user/repo.git` | Clone Git repositories |
| GitHub | `github.com/user/repo` | GitHub shortcuts |
| File | `file:///path/to/file` | Local files |
| S3 | `s3::https://s3.amazonaws.com/bucket/key` | AWS S3 objects |

**Features:**

- Automatic archive extraction (tar, zip, etc.)
- Checksum verification
- Sub-directory selection
- Git ref/branch/tag selection

**Configuration:**

Uses strategy chain for advanced download capabilities (proxy, rewrite, retry).

### GoFetcher

**Source:** [`core/fetcher/go.go`](https://github.com/mrlyc/cmdr/blob/master/core/fetcher/go.go)

Specialized fetcher for Go module downloads.

**Usage:**

```go
fetcher := NewGoFetcher()
err := fetcher.Fetch("golang.org/x/tools/cmd/goimports@latest", "/tmp/goimports")
```

**Features:**

- Go module resolution
- Version/tag specification
- Binary installation from Go packages

## Factory Pattern

Fetchers don't use a factory pattern in the current implementation. Instead, they're instantiated directly:

```go
// In DownloadManager
fetcher := fetcher.NewGoGetterFetcher(cfg, strategyChain)
```

## Usage in DownloadManager

The DownloadManager uses fetchers to download commands during installation:

```go
// core/manager/download.go
func (m *DownloadManager) Define(name, version, location string) (Command, error) {
    // 1. Create temp directory
    tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("cmdr-%s-%s", name, version))
    
    // 2. Fetch file to temp directory
    err := m.fetcher.Fetch(location, tempDir)
    
    // 3. Validate downloaded file
    // 4. Delegate to BinaryManager for storage
    return m.manager.Define(name, version, realLocation)
}
```

## Archive Handling

go-getter automatically handles archives:

| Format | Auto-Extract |
|--------|--------------|
| `.tar.gz` | Yes |
| `.tgz` | Yes |
| `.tar` | Yes |
| `.zip` | Yes |
| `.tar.bz2` | Yes |
| `.tar.xz` | Yes |

**Example:**

```bash
# Downloads and extracts tar.gz automatically
cmdr command install kubectl -v 1.28.0 -l https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl.tar.gz
```

## Sub-Directory Selection

You can select a specific file from an archive:

```bash
# Format: url//subpath
cmdr command install node -v 18.0.0 -l https://nodejs.org/dist/v18.0.0/node-v18.0.0-linux-x64.tar.gz//node-v18.0.0-linux-x64/bin/node
```

## Git Repository Cloning

Fetch from Git repositories:

```bash
# Clone repository
cmdr command install script -v 1.0.0 -l git::https://github.com/user/repo.git

# Specific branch
cmdr command install script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=main

# Specific tag
cmdr command install script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=v1.0.0

# Specific commit
cmdr command install script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=abc123

# Sub-path in repository
cmdr command install script -v 1.0.0 -l git::https://github.com/user/repo.git//scripts/install.sh
```

## Error Handling

Fetchers propagate errors from underlying libraries. The DownloadManager's strategy chain provides retry logic:

```
Fetch Error → DownloadManager → Strategy Chain
                                 ├─ Should Retry? → Retry
                                 └─ Should Fallback? → Next Strategy
```

## Checksum Verification

go-getter supports checksum verification:

```bash
# Format: url?checksum=type:value
cmdr command install kubectl -v 1.28.0 -l \
  "https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl?checksum=sha256:abc123..."
```

## Custom Fetcher Implementation

To implement a custom fetcher:

```go
type CustomFetcher struct {
    // Configuration
}

func (f *CustomFetcher) Fetch(source, destination string) error {
    // 1. Parse source URL
    // 2. Download file(s)
    // 3. Validate integrity
    // 4. Save to destination
    return nil
}

func NewCustomFetcher() *CustomFetcher {
    return &CustomFetcher{}
}
```

Then use in DownloadManager:

```go
customFetcher := fetcher.NewCustomFetcher()
downloadMgr := manager.NewDownloadManager(cfg, binaryMgr, customFetcher)
```

## Testing Fetchers

Fetchers can be tested with mock HTTP servers:

```go
func TestGoGetterFetcher(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("file content"))
    }))
    defer server.Close()
    
    // Test fetch
    fetcher := NewGoGetterFetcher(cfg, strategyChain)
    err := fetcher.Fetch(server.URL, "/tmp/test")
    // Assert no error, file exists
}
```

---

[^1]: Fetcher interface in [`core/fetcher.go`](https://github.com/mrlyc/cmdr/blob/master/core/fetcher.go)
