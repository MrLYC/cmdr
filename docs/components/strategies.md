# Strategies

Strategies implement different approaches for linking commands to the `bin/` directory. CMDR uses the Strategy pattern to allow flexible command activation methods.

## Strategy Interface

The core download strategy interface[^1]:

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

## Strategy Implementations

### DirectStrategy

**Source:** [`core/strategy/direct.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/direct.go)

Downloads files directly without modifications.

**Features:**

- **Retry Logic**[^2]: Retries on timeout and connection errors
- **Fallback Logic**: Falls back to next strategy on network errors
- **Configuration**:
  ```yaml
  download:
    direct:
      timeout: 30        # seconds
      max_retries: 3
  ```

**Error Handling:**

| Error Type | Action |
|------------|--------|
| Timeout | Retry |
| Connection refused/reset | Retry |
| Temporary network error | Retry |
| Other network errors | Fallback to next strategy |

### ProxyStrategy

**Source:** [`core/strategy/proxy.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/proxy.go)

Downloads through HTTP or SOCKS5 proxy.

**Configuration:**

```yaml
download:
  proxy:
    enabled: true
    type: http                           # or "socks5"
    address: http://proxy.example.com:8080
    timeout: 60
    max_retries: 5
```

**Proxy Types:**

- `http` - HTTP/HTTPS proxy
- `socks5` - SOCKS5 proxy

### RewriteStrategy

**Source:** [`core/strategy/rewrite.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/rewrite.go)

Rewrites download URLs using templates (for mirrors, CDNs).

**Configuration:**

```yaml
download:
  rewrite:
    rule: "https://mirror.example.com{{.Path}}"
```

**Template Variables:**

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.URI}}` | Full original URI | `https://github.com/user/repo/file.tar.gz` |
| `{{.Scheme}}` | URI scheme | `https` |
| `{{.Host}}` | URI host | `github.com` |
| `{{.Path}}` | URI path | `/user/repo/file.tar.gz` |
| `{{.Query}}` | URI query | `?token=abc` |
| `{{.Fragment}}` | URI fragment | `#section` |

**Example:**

```yaml
# Redirect GitHub downloads to mirror
download:
  rewrite:
    rule: "https://ghproxy.com/{{.URI}}"
```

Original: `https://github.com/nodejs/node/archive/v18.0.0.tar.gz`
Rewritten: `https://ghproxy.com/https://github.com/nodejs/node/archive/v18.0.0.tar.gz`

### ChainStrategy

**Source:** [`core/strategy/chain.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/chain.go)

Combines multiple strategies with retry and fallback logic.

**Execution Flow:**

```
┌─────────────────────────────────────────────────────────┐
│                    ChainStrategy                         │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ Rewrite  │→ │  Proxy   │→ │  Direct  │→ │  ...    │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬────┘ │
│       │ Retry       │ Retry       │ Retry       │       │
│       │ ↓           │ ↓           │ ↓           │       │
│       │ Fallback →  │ Fallback →  │ Fallback →  │       │
└───────┴─────────────┴─────────────┴─────────────┴───────┘
```

**Logic:**

1. Try first enabled strategy
2. On retriable error: retry up to `max_retries`
3. On fallback error: try next strategy
4. If all strategies fail: return error

## Conditional Strategy Execution

Strategies can be conditionally enabled based on URI patterns:

```yaml
download:
  proxy:
    enabled: true
    condition:
      schemes: [https]              # Only HTTPS
      hosts: [github.com]           # Only GitHub
      patterns: ["*.github.com"]    # Glob patterns
```

**Matching Rules:**

- **Schemes**: Match URI scheme (http, https, git, etc.)
- **Hosts**: Match exact hostname
- **Patterns**: Match hostname with glob patterns

**Example:**

```yaml
# Use proxy only for GitHub
download:
  proxy:
    enabled: true
    address: http://github-proxy:8080
    condition:
      hosts: [github.com]

# Rewrite only for GitLab
  rewrite:
    rule: "https://gitlab-mirror.com{{.Path}}"
    condition:
      hosts: [gitlab.com]
```

Result:
- `https://github.com/file` → Uses proxy
- `https://gitlab.com/file` → Uses rewrite
- `https://nodejs.org/file` → Uses direct

## Strategy Configuration

Complete strategy configuration:

```yaml
download:
  # Always try direct first
  direct:
    timeout: 30
    max_retries: 3
    enabled: true

  # Use proxy for GitHub
  proxy:
    enabled: true
    type: http
    address: http://proxy.example.com:8080
    timeout: 60
    max_retries: 5
    condition:
      hosts: [github.com]

  # Use mirror for GitLab
  rewrite:
    rule: "https://mirror.com{{.Path}}"
    enabled: true
    condition:
      hosts: [gitlab.com]
```

## Legacy URL Replacement

CMDR still supports the legacy `download.replace` configuration:

```yaml
download:
  replace:
    match: "^https://raw.githubusercontent.com/.*$"
    template: "https://ghproxy.com/{{ .input | urlquery }}"
```

This is internally converted to a RewriteStrategy.

## Adding Custom Strategies

To add a new strategy:

1. **Implement the interface:**

```go
type CustomStrategy struct {
    config *StrategyConfig
}

func (s *CustomStrategy) Name() string {
    return "custom"
}

func (s *CustomStrategy) Prepare(uri string) (string, error) {
    // Transform URI as needed
    return transformedURI, nil
}

func (s *CustomStrategy) ShouldRetry(err error) bool {
    // Define retry logic
    return isRetriableError(err)
}

func (s *CustomStrategy) ShouldFallback(err error) bool {
    // Define fallback logic
    return isFallbackError(err)
}

func (s *CustomStrategy) Configure(cfg Configuration) error {
    // Load configuration
    return nil
}

func (s *CustomStrategy) IsEnabled(uri string) bool {
    return s.config.Matches(uri)
}
```

2. **Register in DownloadManager:**

```go
customStrategy := NewCustomStrategy()
chain := NewStrategyChain(
    NewRewriteStrategy(),
    NewProxyStrategy(),
    customStrategy,
    NewDirectStrategy(),
)
```

## Testing Strategies

Strategies can be unit tested:

```go
func TestDirectStrategy(t *testing.T) {
    strategy := NewDirectStrategy()
    cfg := core.NewConfiguration()
    cfg.Set("download.direct.timeout", 30)
    
    err := strategy.Configure(cfg)
    // Assert no error
    
    uri, err := strategy.Prepare("https://example.com/file")
    // Assert URI unchanged
}
```

---

[^1]: Strategy interface in [`core/strategy/strategy.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/strategy.go)
[^2]: DirectStrategy retry logic in [`core/strategy/direct.go`](https://github.com/mrlyc/cmdr/blob/master/core/strategy/direct.go) L31-L50
