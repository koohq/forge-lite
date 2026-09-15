# Contributing to Forge Lite

Thank you for your interest in contributing to Forge Lite! This document provides guidelines for development, testing, and submitting contributions.

## Prerequisites

- **Go**: 1.27.1 or higher
- **Task** ([go-task](https://taskfile.dev)): Recommended for task automation
- **golangci-lint**: 2.13.2 or higher
- **GoReleaser**: v2.18.1 or higher (for release snapshot builds)

## Development Workflow

1. Fork the repository and create a feature branch:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your code changes adhering to project standards.

## Code Quality & Testing

Before submitting a pull request, ensure all validation checks pass:

### Formatting
Format all Go source files with `gofmt`:
```bash
gofmt -l .
```
No output indicates that all files are properly formatted. To format files in place:
```bash
gofmt -w .
```

### Static Analysis & Vet
Run Go static analysis:
```bash
go vet ./...
```

### Unit Tests
Run the test suite:
```bash
# Using Task
task test

# Or using Go directly
go test -v ./...
```
In Linux / CI environments with CGo enabled, race conditions can also be checked with:
```bash
go test -race ./...
```

### Linter
Run `golangci-lint`:
```bash
golangci-lint run
```

### Build & Release Verification
Test local development compilation:
```bash
# Using Task
task build

# Or using Go directly
go build -trimpath -o bin/forge-lite .
```

Test multi-platform snapshot builds:
```bash
# Using Task
task release

# Or using GoReleaser directly
goreleaser build --snapshot --clean
```

## Commit Conventions

This project strictly adheres to [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` A new user-facing feature (triggers a Minor release)
- `fix:` A bug fix for user-facing behavior (triggers a Patch release)
- `docs:` Documentation changes only
- `style:` Code style changes (formatting, missing semicolons, etc.)
- `refactor:` Code changes that neither fix a bug nor add a feature
- `perf:` Performance improvements
- `test:` Adding or updating tests
- `chore:` Repository maintenance, tooling, dependencies
- `ci:` CI/CD workflow configuration

All commit messages must be written in English.

## Automated Releases & Tool Maintenance

- **Release Automation**: Releases are automated using Release Please and GoReleaser. Merging a release pull request automatically creates a tag and publishes cross-compiled binary assets to GitHub Releases.
- **Tool Version Maintenance**: Note that tool versions specified in CI workflows (such as `golangci-lint` in `ci.yml` and `goreleaser` in `ci.yml` / `release.yml`) are pinned by version strings and must be verified and updated manually when new tool versions are released.
