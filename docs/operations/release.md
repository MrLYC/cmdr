# Release

This document describes the release process for CMDR.

## Release Strategy

CMDR follows semantic versioning (SemVer):

- **Major** (vX.0.0): Breaking changes
- **Minor** (vX.Y.0): New features, backward compatible
- **Patch** (vX.Y.Z): Bug fixes, backward compatible

## Release Process

### 1. Prepare Release

```bash
# Ensure you're on master branch
git checkout master
git pull origin master

# Run all tests
go test ./...
cd integration-tests && ./start.sh && cd ..

# Update version references if needed
# (Currently version is determined by git tags)
```

### 2. Create Release Tag

```bash
# Create annotated tag
git tag -a v1.2.3 -m "Release v1.2.3: Description of changes"

# Push tag to trigger release workflow
git push origin v1.2.3
```

### 3. Automated Build

The GitHub Actions `release.yml` workflow automatically:

1. **Builds** binaries for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)

2. **Creates** GitHub release with:
   - Release notes
   - Pre-built binaries
   - Checksums

3. **Publishes** release assets

**Workflow:** `.github/workflows/release.yml`

## GoReleaser Configuration

CMDR uses [GoReleaser](https://goreleaser.com/) for building and publishing releases.

**Configuration:** `.goreleaser.yml`

### Build Configuration

```yaml
# .goreleaser.yml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}
```

### Archive Configuration

```yaml
archives:
  - format: tar.gz
    name_template: "cmdr_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
```

## Release Checklist

Before creating a release:

- [ ] All tests pass (`go test ./...`)
- [ ] Integration tests pass (`integration-tests/start.sh`)
- [ ] Code is linted (`golangci-lint run`)
- [ ] CHANGELOG updated (if maintaining one)
- [ ] Documentation updated
- [ ] No breaking changes (or major version bump)

## Manual Release (If Needed)

If automated release fails, build manually:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Build release (dry run)
goreleaser release --snapshot --rm-dist

# Build actual release
export GITHUB_TOKEN=your_token
goreleaser release --rm-dist
```

## Release Assets

Each release includes:

### Binaries

| Asset | Platform | Architecture |
|-------|----------|--------------|
| `cmdr_vX.Y.Z_linux_amd64.tar.gz` | Linux | x86_64 |
| `cmdr_vX.Y.Z_linux_arm64.tar.gz` | Linux | ARM64 |
| `cmdr_vX.Y.Z_darwin_amd64.tar.gz` | macOS | Intel |
| `cmdr_vX.Y.Z_darwin_arm64.tar.gz` | macOS | Apple Silicon |
| `cmdr_vX.Y.Z_windows_amd64.zip` | Windows | x86_64 |

### Checksums

- `checksums.txt` - SHA256 checksums of all assets

## Upgrading CMDR

Users can upgrade using the built-in command:

```bash
cmdr upgrade

# Or specific version
cmdr upgrade --release v1.2.3
```

The upgrade process:

1. Queries GitHub API for latest release (or specified version)
2. Downloads appropriate binary for current platform
3. Replaces current binary
4. Runs `cmdr init --upgrade` to update profile

**Implementation:** [`cmd/upgrade.go`](https://github.com/mrlyc/cmdr/blob/master/cmd/upgrade.go)

## Version Detection

CMDR detects its version from:

1. Build-time ldflags (`-X main.version=...`)
2. Git tags (if built from source)

**Display version:**

```bash
cmdr version
```

## Pre-Release Versions

For testing before official release:

```bash
# Create pre-release tag
git tag -a v1.2.3-rc.1 -m "Release candidate 1"
git push origin v1.2.3-rc.1
```

GoReleaser will mark it as "pre-release" on GitHub.

## Hotfix Process

For critical bugs in released versions:

1. Create hotfix branch from tag:
   ```bash
   git checkout -b hotfix/v1.2.4 v1.2.3
   ```

2. Fix bug and commit:
   ```bash
   git commit -m "fix: critical bug description"
   ```

3. Create patch release:
   ```bash
   git tag -a v1.2.4 -m "Hotfix v1.2.4: Fix critical bug"
   git push origin v1.2.4
   ```

4. Merge back to master:
   ```bash
   git checkout master
   git merge hotfix/v1.2.4
   git push origin master
   ```

## Release Announcement

After release is published:

1. **GitHub Release Notes**: Auto-generated from commits or manually edited
2. **README Badge**: Automatically updated to show latest version
3. **Documentation**: Update if needed

## Troubleshooting Releases

### Build Fails

Check `.github/workflows/release.yml` logs for errors:

```bash
# Common issues:
# - Go version mismatch
# - Missing dependencies
# - Linting failures
```

### GoReleaser Fails

```bash
# Test locally
goreleaser release --snapshot --rm-dist

# Check configuration
goreleaser check
```

### Binary Not Working

```bash
# Verify binary is executable
chmod +x cmdr

# Check if it's for correct platform
file cmdr

# Test locally before release
./cmdr version
```

## Security Considerations

### Signing Releases

Currently, releases are not signed. To add signing:

1. Generate GPG key
2. Add key to GoReleaser config
3. Upload public key to GitHub

### Checksum Verification

Users should verify checksums:

```bash
# Download checksum file
curl -LO https://github.com/mrlyc/cmdr/releases/download/v1.2.3/checksums.txt

# Verify binary
sha256sum -c checksums.txt --ignore-missing
```

## Rollback

If a release has critical issues:

1. **Mark as Pre-release**: Edit GitHub release, check "This is a pre-release"
2. **Delete Tag** (if not yet widely distributed):
   ```bash
   git tag -d v1.2.3
   git push origin :refs/tags/v1.2.3
   ```
3. **Create Patch Release**: Follow hotfix process

## Future Improvements

Potential enhancements to release process:

- [ ] Automated changelog generation from commits
- [ ] GPG signing of release binaries
- [ ] Docker image publishing
- [ ] Homebrew formula auto-update
- [ ] Release announcement automation
