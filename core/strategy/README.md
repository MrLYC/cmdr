# Download Strategy Feature

## Overview

The download strategy feature provides flexible download strategies for the `install` command. It supports multiple strategies that can be tried in sequence when downloads fail.

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
