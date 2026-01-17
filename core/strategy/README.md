# Download Strategy Feature

## Overview

The download strategy feature provides flexible download strategies for the `install` command. It supports multiple strategies that can be selectively enabled based on URI scheme and hostname patterns.

## Supported Strategies

### 1. Direct Strategy (Default)
- **Name**: `direct`
- **Description**: Downloads directly without any modifications
- **Configuration**:
  ```yaml
  download:
    direct:
      timeout: 30       # timeout in seconds (default: 30)
      max_retries: 3     # number of retries (default: 3)
      enabled: true       # enable strategy (default: true)
      condition:
        schemes: [http, https, git]   # only enable for these schemes
        hosts: [github.com, gitlab.com]  # only enable for these hosts
        patterns: ["*.github.com", "git*.com"]  # glob patterns for hosts
  ```
- **Conditional Usage**: Can be enabled/disabled for specific schemes, hosts, or patterns

### 2. Proxy Strategy
- **Name**: `proxy`
- **Description**: Downloads through an HTTP or SOCKS5 proxy
- **Configuration**:
  ```yaml
  download:
    proxy:
      enabled: true                         # enable strategy (default: false)
      type: http                            # proxy type: "http" or "socks5"
      address: http://proxy.example.com:8080  # proxy URL
      timeout: 60                           # timeout in seconds (default: 30)
      max_retries: 5                         # number of retries (default: 3)
      condition:
        schemes: [http, https]              # only use proxy for these schemes
        hosts: [github.com]                  # only use proxy for these hosts
        patterns: ["*.github.com"]             # only use proxy for matching hosts
  ```
- **Conditional Usage**: Can be configured to use proxy only for specific domains

### 3. Rewrite Strategy
- **Name**: `rewrite`
- **Description**: Rewrites download URLs using a template
- **Configuration**:
  ```yaml
  download:
    rewrite:
      rule: "https://mirror.example.com/{{.Path}}"  # URL template
      enabled: true                               # enable strategy (default: false)
      condition:
        schemes: [https]                         # only rewrite https URLs
        hosts: [github.com, gitlab.com]         # only rewrite these domains
        patterns: ["*.github.com"]                # only rewrite matching hosts
  ```
- **Template Variables**:
  - `{{.URI}}`: Full original URI
  - `{{.Scheme}}`: URI scheme (http, https, etc.)
  - `{{.Host}}`: URI host
  - `{{.Path}}`: URI path
  - `{{.Query}}`: URI query string
  - `{{.Fragment}}`: URI fragment
- **Conditional Usage**: Can be configured to rewrite only specific URLs

## Conditional Strategy Selection

Strategies can be conditionally enabled based on:

### 1. URI Schemes
```yaml
download:
  proxy:
    enabled: true
    condition:
      schemes: [http, https]  # only use proxy for http/https
```

### 2. Host Names
```yaml
download:
  proxy:
    enabled: true
    condition:
      hosts: [github.com, gitlab.com]  # only use proxy for these domains
```

### 3. Glob Patterns
```yaml
download:
  proxy:
    enabled: true
    condition:
      patterns: ["*.github.com", "*.gitlab.com"]  # match with glob patterns
```

### 4. Combined Conditions
```yaml
download:
  proxy:
    enabled: true
    condition:
      schemes: [https]
      hosts: [github.com]
```

This will use proxy **only** for:
- HTTPS URLs
- From github.com domain

## Strategy Chain Execution

Strategies are evaluated in order:

1. Each strategy checks if it's enabled for the current URI
2. Only enabled strategies are tried
3. Strategies are tried in the order they were added to the chain
4. If all strategies are disabled, all strategies are tried (backward compatibility)

### Example 1: GitHub with Proxy
```yaml
download:
  proxy:
    enabled: true
    condition:
      hosts: [github.com]
  rewrite:
    enabled: true
    rule: "https://mirror.example.com/{{.Path}}"
    condition:
      hosts: [github.com]
```

For `https://github.com/nodejs/node/archive/v18.0.0.tar.gz`:
1. Rewrite strategy: matches github.com → Rewrites to mirror
2. Proxy strategy: matches github.com → Uses proxy
3. Direct strategy: Always enabled → Skipped (previous succeeded)

### Example 2: Different Proxies for Different Domains
```yaml
download:
  proxy1:
    type: http
    address: http://proxy1.example.com:8080
    condition:
      hosts: [github.com]
  proxy2:
    type: http
    address: http://proxy2.example.com:8080
    condition:
      hosts: [gitlab.com]
```

For `https://github.com/repo/file` → Uses proxy1
For `https://gitlab.com/repo/file` → Uses proxy2
For `https://other.com/file` → Direct download (no proxy)

### Example 3: Rewrite Only for GitHub
```yaml
download:
  rewrite:
    rule: "https://mirror.example.com{{.Path}}"
    condition:
      hosts: [github.com]
```

For `https://github.com/nodejs/node/archive/v18.tar.gz`:
→ Rewritten to: `https://mirror.example.com/archive/v18.tar.gz`

For `https://nodejs.org/dist/node.tar.gz`:
→ No rewrite (not github.com)

### Example 4: Proxy Only for HTTPS
```yaml
download:
  proxy:
    enabled: true
    type: http
    address: http://proxy.example.com:8080
    condition:
      schemes: [https]
```

For `https://github.com/file` → Uses proxy
For `http://github.com/file` → Direct download

## Retry and Fallback Logic

1. **Retry**: If a strategy returns a retriable error (timeout, connection error), it will retry up to `max_retries` times
2. **Fallback**: If all retries fail with a network error, the next enabled strategy is tried
3. **Failure**: If all enabled strategies fail, the download is considered failed

### Retriable Errors
- Timeout errors
- Connection errors (refused, reset, unreachable)
- Temporary network errors

### Fallback Errors
- Any network error

## Complete Configuration Example

```yaml
download:
  # Direct strategy (always enabled)
  direct:
    timeout: 30
    max_retries: 3
    enabled: true

  # Proxy for GitHub
  proxy:
    enabled: true
    type: http
    address: http://proxy.example.com:8080
    timeout: 60
    max_retries: 5
    condition:
      hosts: [github.com]

  # Rewrite for GitHub to mirror
  rewrite:
    rule: "https://mirror.example.com{{.Path}}?token=abc123"
    enabled: true
    condition:
      hosts: [github.com]

  # Legacy replacement (still supported)
  replace:
    - pattern: "github.com"
      replacement: "github.com.cn"

# Proxy configuration (still supported for backward compatibility)
proxy:
  http: http://proxy.example.com:8080
  https: https://proxy.example.com:8443
```

## Usage Examples

### Install with direct download
```bash
cmdr install node 18.0.0 https://nodejs.org/dist/v18.0.0/node-v18.0.0-linux-x64.tar.xz
```

### Install with proxy for GitHub
```yaml
# config.yaml
download:
  proxy:
    enabled: true
    type: socks5
    address: socks5://proxy.example.com:1080
    condition:
      hosts: [github.com]
```

```bash
cmdr install node 18.0.0 https://github.com/nodejs/node/archive/v18.0.0.tar.gz
```

### Install with GitHub mirror
```yaml
# config.yaml
download:
  rewrite:
    rule: "https://mirror.example.com{{.Path}}"
    enabled: true
    condition:
      hosts: [github.com]
```

```bash
cmdr install node 18.0.0 https://github.com/nodejs/node/archive/v18.0.0.tar.gz
# Will download from: https://mirror.example.com/archive/v18.0.0.tar.gz
```

### Install from different sources with different strategies
```yaml
# config.yaml
download:
  # Proxy for GitHub
  proxy:
    enabled: true
    type: http
    address: http://github-proxy:8080
    condition:
      hosts: [github.com]

  # Rewrite for GitLab
  rewrite:
    rule: "https://gitlab-mirror{{.Path}}"
    enabled: true
    condition:
      hosts: [gitlab.com]

  # Direct for everything else
  direct:
    enabled: true
    timeout: 30
```

```bash
# GitHub → proxy
cmdr install node https://github.com/nodejs/node/archive/v18.tar.gz

# GitLab → rewrite
cmdr install gitlab https://gitlab.com/gitlab/gitlab/archive/v15.tar.gz

# Node.js → direct
cmdr install node https://nodejs.org/dist/node.tar.gz
```

## Advanced: Multiple Proxy Servers

You can define multiple proxy strategies for different domains by modifying the strategy chain in code:

```go
// In core/manager/download.go
proxy1 := strategy.NewProxyStrategy()
proxy2 := strategy.NewProxyStrategy()

// Configure proxy1 for GitHub
cfg.Set("download.proxy1.enabled", true)
cfg.Set("download.proxy1.address", "http://proxy1.example.com:8080")
cfg.Set("download.proxy1.condition.hosts", []string{"github.com"})

// Configure proxy2 for GitLab
cfg.Set("download.proxy2.enabled", true)
cfg.Set("download.proxy2.address", "http://proxy2.example.com:8080")
cfg.Set("download.proxy2.condition.hosts", []string{"gitlab.com"})

strategyChain := strategy.NewStrategyChain(
    strategy.NewDirectStrategy(),
    strategy.NewRewriteStrategy(),
    proxy1,
    proxy2,
)
```

## Backward Compatibility

The feature maintains full backward compatibility:

1. **Existing proxy configuration** (`proxy.http`, `proxy.https`) still works
2. **Existing download replacement** (`download.replace`) still works
3. **Default behavior** unchanged when no strategy configuration is provided
4. **No conditions configured** → all strategies are tried (as before)

## Extending Strategies

To add a new download strategy with conditional support:

1. Create a new strategy file in `core/strategy/`
2. Implement the `DownloadStrategy` interface:
   ```go
   type DownloadStrategy interface {
       Name() string
       Prepare(uri string) (string, error)
       ShouldRetry(err error) bool
       ShouldFallback(err error) bool
       Configure(cfg core.Configuration) error
       IsEnabled(uri string) bool
       SetEnabled(enabled bool)
   }
   ```
3. Parse condition configuration in `Configure()` method
4. Implement `IsEnabled(uri)` to check if strategy should be used for given URI
5. Add strategy to the chain in `core/manager/download.go`

## Testing

Run tests for the strategy package:

```bash
go test ./core/strategy/... -v
```


### 2. Proxy Strategy
- **Name**: `proxy`
- **Description**: Downloads through an HTTP or SOCKS5 proxy
- **Configuration**:
  ```yaml
  download:
    proxy:
      enabled: true                    # enable proxy (default: false)
      type: http                         # proxy type: "http" or "socks5"
      address: http://proxy.example.com:8080  # proxy URL
      timeout: 60                       # timeout in seconds (default: 30)
      max_retries: 5                     # number of retries (default: 3)
  ```

### 3. Rewrite Strategy
- **Name**: `rewrite`
- **Description**: Rewrites download URLs using a template
- **Configuration**:
  ```yaml
  download:
    rewrite:
      rule: "https://mirror.example.com/{{.Path}}"  # URL template
  ```
- **Template Variables**:
  - `{{.URI}}`: Full original URI
  - `{{.Scheme}}`: URI scheme (http, https, etc.)
  - `{{.Host}}`: URI host
  - `{{.Path}}`: URI path
  - `{{.Query}}`: URI query string
  - `{{.Fragment}}`: URI fragment

## Strategy Chain

Strategies are executed in the following order:
1. Direct strategy (always)
2. Rewrite strategy (if configured)
3. Proxy strategy (if configured)

### Retry and Fallback Logic

1. **Retry**: If a strategy returns a retriable error (timeout, connection error), it will retry up to `max_retries` times
2. **Fallback**: If all retries fail with a network error, the next strategy is tried
3. **Failure**: If all strategies fail, the download is considered failed

### Retriable Errors

- Timeout errors
- Connection errors (refused, reset, unreachable)
- Temporary network errors

### Fallback Errors

- Any network error

## Configuration Example

Complete configuration example:

```yaml
download:
  # Direct strategy
  direct:
    timeout: 30
    max_retries: 3

  # Proxy strategy
  proxy:
    enabled: true
    type: http
    address: http://proxy.example.com:8080
    timeout: 60
    max_retries: 5

  # Rewrite strategy
  rewrite:
    rule: "https://cdn.example.com{{.Path}}"

  # Legacy replacement (still supported)
  replace:
    - pattern: "github.com"
      replacement: "github.com.cn"

# Proxy configuration (still supported for backward compatibility)
proxy:
  http: http://proxy.example.com:8080
  https: https://proxy.example.com:8443
```

## Usage Examples

### Install with direct download
```bash
cmdr install node 18.0.0 https://nodejs.org/dist/v18.0.0/node-v18.0.0-linux-x64.tar.xz
```

### Install with proxy
```yaml
# config.yaml
download:
  proxy:
    enabled: true
    type: socks5
    address: socks5://proxy.example.com:1080
```

```bash
cmdr install node 18.0.0 https://nodejs.org/dist/v18.0.0/node-v18.0.0-linux-x64.tar.xz
```

### Install with URL rewrite
```yaml
# config.yaml
download:
  rewrite:
    rule: "https://mirror.example.com{{.Path}}?token=abc123"
```

```bash
cmdr install node 18.0.0 https://github.com/nodejs/node/archive/v18.0.0.tar.gz
# Will download from: https://mirror.example.com/archive/v18.0.0.tar.gz?token=abc123
```

## Backward Compatibility

The feature maintains full backward compatibility:

1. **Existing proxy configuration** (`proxy.http`, `proxy.https`) still works
2. **Existing download replacement** (`download.replace`) still works
3. **Default behavior** unchanged when no strategy configuration is provided

## Extending Strategies

To add a new download strategy:

1. Create a new strategy file in `core/strategy/`
2. Implement the `DownloadStrategy` interface:
   ```go
   type DownloadStrategy interface {
       Name() string
       Prepare(uri string) (string, error)
       ShouldRetry(err error) bool
       ShouldFallback(err error) bool
       Configure(cfg core.Configuration) error
   }
   ```

3. Add the strategy to the chain in `core/manager/download.go`

## Testing

Run tests for the strategy package:

```bash
go test ./core/strategy/... -v
```
