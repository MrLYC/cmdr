# CMDR
[![unittest](https://github.com/mrlyc/cmdr/actions/workflows/unittest.yml/badge.svg)](https://github.com/mrlyc/cmdr/actions/workflows/unittest.yml) [![integration-test](https://github.com/mrlyc/cmdr/actions/workflows/integration-test.yml/badge.svg)](https://github.com/mrlyc/cmdr/actions/workflows/integration-test.yml) [![codecov](https://codecov.io/gh/MrLYC/cmdr/branch/master/graph/badge.svg?token=mo4TJP4mQt)](https://codecov.io/gh/MrLYC/cmdr) ![Go version](https://img.shields.io/github/go-mod/go-version/mrlyc/cmdr) ![release](https://img.shields.io/github/v/release/mrlyc/cmdr?label=version)

CMDR is a simple command version management tool that helps you quickly switch from multiple command versions.

## Documentation

Full documentation is available at [https://mrlyc.github.io/cmdr/](https://mrlyc.github.io/cmdr/)

**中文文档 (Chinese Documentation):** [https://www.zdoc.app/zh/MrLYC/cmdr?refresh=1770878793863](https://www.zdoc.app/zh/MrLYC/cmdr?refresh=1770878793863)

- [Installation Guide](https://mrlyc.github.io/cmdr/getting-started/installation/) - Get started with CMDR
- [Quick Start](https://mrlyc.github.io/cmdr/getting-started/quick-start/) - First steps tutorial
- [Configuration](https://mrlyc.github.io/cmdr/getting-started/configuration/) - Configuration reference
- [Architecture](https://mrlyc.github.io/cmdr/architecture/overview/) - System design and internals
- [API Reference](https://mrlyc.github.io/cmdr/api/core-interfaces/) - Interface documentation
- [Development Guide](https://mrlyc.github.io/cmdr/operations/build-and-test/) - Build, test, and contribute

## Installation

### Script
Run one for the following command to install the latest version of CMDR:

```shell
curl -o- https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh | ${SHELL:-bash}
```

For Chinese users, you can install the CMDR via a proxy:

```shell
curl -o- https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh | bash -s -p
```

### Manual
1. Download the latest release from [GitHub](https://github.com/mrlyc/cmdr/releases/latest);
2. Make sure the download asset is executable;
3. Run the following command to install the binary:
    ```shell
    /path/to/cmdr init
    ```
4. Restart your shell and run the following command to verify the installation:
    ```shell
    cmdr version
    ```

## Get Started

**Note:** The `cmdr command xxx` format is deprecated. Use `cmdr xxx` directly (e.g., `cmdr install` instead of `cmdr command install`).

To install a new command, run the following command:
```shell
cmdr install -n <command-name> -v <version> -l <path-or-url>
```

Then you can list the installed commands by running the following command:
```shell
cmdr list -n <command-name>
```

Use a specified command version:
```shell
cmdr use -n <command-name> -v <version>
```

## Features

### Upgrade
To upgrade the CMDR, just run:
```shell
cmdr upgrade
```

### Url replacement
Speed up the download process by replacing the `url` to github proxy:
```shell
cmdr config set -k download.replace -v '{"match": "^https://raw.githubusercontent.com/.*$", "template": "https://ghproxy.com/{{ .input | urlquery }}"}'
cmdr install -n install.sh -v 0.0.0 -l https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh
```

## Supported Download Protocols

CMDR uses [go-getter](https://github.com/hashicorp/go-getter) for downloading commands, which supports multiple protocols:

### HTTP/HTTPS
Direct downloads from HTTP/HTTPS URLs:
```shell
cmdr install -n kubectl -v 1.28.0 -l https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl
```

### Git Repositories
Clone from Git repositories with support for branches, tags, and commits:
```shell
# Clone repository
cmdr install -n script -v 1.0.0 -l git::https://github.com/user/repo.git

# Specific branch
cmdr install -n script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=main

# Specific tag
cmdr install -n script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=v1.0.0

# Specific commit
cmdr install -n script -v 1.0.0 -l git::https://github.com/user/repo.git?ref=abc123

# Sub-path in repository
cmdr install -n script -v 1.0.0 -l git::https://github.com/user/repo.git//scripts/install.sh
```

### GitHub Shortcuts
Convenient shortcuts for GitHub repositories:
```shell
# Short form for GitHub repositories (defaults to HTTPS)
cmdr install -n tool -v 1.0.0 -l github.com/user/repo

# With sub-path
cmdr install -n tool -v 1.0.0 -l github.com/user/repo//path/to/binary
```

### Archive Support
Automatically extracts various archive formats:
- `.tar.gz`, `.tgz` - Gzip compressed tar archives
- `.tar.bz2`, `.tbz2` - Bzip2 compressed tar archives
- `.tar.xz`, `.txz` - XZ compressed tar archives
- `.tar.zst`, `.tzst` - Zstandard compressed tar archives
- `.zip` - ZIP archives
- `.gz` - Gzip compressed files
- `.bz2` - Bzip2 compressed files
- `.xz` - XZ compressed files
- `.zst` - Zstandard compressed files

```shell
# Automatically extracts tar.gz archive
cmdr install -n kubectl -v 1.28.0 -l https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl.tar.gz

# Select specific file from archive using // separator
cmdr install -n node -v 18.0.0 -l https://nodejs.org/dist/v18.0.0/node-v18.0.0-linux-x64.tar.gz//node-v18.0.0-linux-x64/bin/node

# Disable automatic extraction
cmdr install -n archive -v 1.0.0 -l https://example.com/file.tar.gz?archive=false
```

### Force Protocol Override
Use `::` to force a specific protocol detector:
```shell
# Force Git protocol even if URL looks like HTTPS
cmdr install -n tool -v 1.0.0 -l git::https://example.com/tool.git

# Force HTTP protocol
cmdr install -n tool -v 1.0.0 -l http::example.com/tool
```

### S3 Support
Download from Amazon S3 buckets:
```shell
# AWS S3 (various addressing schemes)
cmdr install -n tool -v 1.0.0 -l s3::https://s3.amazonaws.com/bucket/tool
cmdr install -n tool -v 1.0.0 -l s3::https://s3-eu-west-1.amazonaws.com/bucket/tool

# With credentials (prefer using AWS environment variables)
cmdr install -n tool -v 1.0.0 -l "s3::https://s3.amazonaws.com/bucket/tool?aws_access_key_id=KEYID&aws_access_key_secret=SECRET"
```

### GCS Support
Download from Google Cloud Storage:
```shell
# Google Cloud Storage
cmdr install -n tool -v 1.0.0 -l gcs::https://www.googleapis.com/storage/v1/bucket/foo.zip
```

### Local Files
Reference local files or directories:
```shell
# Relative path
cmdr install -n tool -v 1.0.0 -l ./path/to/binary

# Absolute path
cmdr install -n tool -v 1.0.0 -l /usr/local/bin/binary

# File URL
cmdr install -n tool -v 1.0.0 -l file:///path/to/binary
```

### Checksums
Verify downloads with checksums:
```shell
# Add checksum verification
cmdr install -n kubectl -v 1.28.0 -l "https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl?checksum=sha256:abc123..."

# Checksum file (automatically detected)
cmdr install -n tool -v 1.0.0 -l https://example.com/tool  # Looks for tool.sha256, tool.sha256sum, etc.
```

For more details and advanced options, see the [go-getter documentation](https://github.com/hashicorp/go-getter).