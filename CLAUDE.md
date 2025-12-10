# CLAUDE.md - AI Assistant Guide for go-vcsfetch

## Project Overview

**go-vcsfetch** is a Go library that provides VCS (Version Control System) fetching and cloning capabilities with no required runtime dependencies. It enables retrieval of individual files or entire repositories from Git-based version control systems.

**Primary Use Cases:**
- Fetch single files from remote repositories (e.g., configuration files)
- Clone entire repositories or specific folders
- Support SPDX downloadLocation attributes
- Work with common git-url schemes (GitHub, GitLab, etc.)

**Key Features:**
- Pure Go implementation (no git binary required, though used when available for performance)
- Supports http, https, ssh, and git TCP protocols
- Memory-backed or filesystem-backed operations
- Sparse cloning support
- Semver-aware tag resolution
- Optimized fetching for GitHub and GitLab via raw content URLs

## Repository Structure

```
go-vcsfetch/
├── .github/workflows/       # CI/CD workflows
│   ├── 01-golang-lint.yaml  # golangci-lint checks
│   ├── 02-test.yaml         # Tests on Ubuntu/Windows with coverage
│   └── 03-govulncheck.yaml  # Vulnerability scanning
├── internal/                # Internal packages (not exported)
│   ├── download/            # Raw content download from SCMs
│   ├── git/                 # Git operations wrapper
│   │   ├── capabilities.go  # Git binary detection
│   │   ├── fs.go            # Filesystem operations
│   │   ├── git.go           # Core git repository logic
│   │   ├── native.go        # Native git binary operations
│   │   └── ref.go           # Reference resolution (branches, tags, semver)
│   └── giturl/              # URL parsing for git platforms
│       ├── github/          # GitHub-specific URL handling
│       ├── gitlab/          # GitLab-specific URL handling
│       └── providers.go     # Provider detection and routing
├── cloner.go                # Cloner API (repository cloning)
├── fetcher.go               # Fetcher API (single file fetching)
├── giturl.go                # GitLocator implementation
├── spdx.go                  # SPDXLocator implementation
├── locator.go               # Locator interface definition
├── options.go               # Options for Fetcher/Cloner configuration
├── errors.go                # Error definitions
├── doc.go                   # Package documentation
└── *_test.go                # Test files

```

## Core Architecture

### Public API Components

#### 1. Fetcher (`fetcher.go`)
- **Purpose**: Retrieve single files from VCS repositories
- **Thread Safety**: Stateless, safe for concurrent use
- **Key Methods**:
  - `Fetch(ctx, writer, location)` - Fetch from URL string
  - `FetchURL(ctx, writer, url)` - Fetch from parsed URL
  - `FetchLocator(ctx, writer, locator)` - Fetch using Locator interface
- **Optimization**: Automatically uses raw content URLs for GitHub/GitLab when possible

#### 2. Cloner (`cloner.go`)
- **Purpose**: Clone entire repositories or folders
- **Thread Safety**: Stateful, NOT safe for concurrent use
- **Key Methods**:
  - `Clone(ctx, repoURL)` - Clone repository
  - `CloneURL(ctx, url)` - Clone from parsed URL
  - `CloneLocator(ctx, locator)` - Clone using Locator interface
  - `FS()` - Access cloned content as `fs.FS`
  - `FetchFromClone(ctx, writer, location)` - Fetch files from cloned repo
  - `Close()` - Clean up resources
- **Use Case**: Multiple file fetches from same repository

#### 3. Locators
Two implementations of the `Locator` interface:

**SPDXLocator (`spdx.go`)**:
- Standard SPDX format: `<vcs>+<transport>://<host>/<path>@<version>#<file>`
- Example: `git+https://github.com/user/repo@v1.0.0#README.md`
- Requires URL fragment (file path)
- Used for standardized, unambiguous references

**GitLocator (`giturl.go`)**:
- Platform-specific URL formats (GitHub, GitLab, etc.)
- Auto-detects provider and transforms to standard git URL
- Fallback when SPDX parsing fails

### Internal Packages

#### internal/git
- **Purpose**: Wrapper around go-git library with enhancements
- **Key Files**:
  - `git.go`: Core `Repository` type, fetch/clone operations
  - `capabilities.go`: Detect and use native git binary when available
  - `native.go`: Execute native git commands for performance
  - `ref.go`: Reference resolution including semver-aware tag matching
  - `fs.go`: Filesystem abstraction (memory or disk-backed)

#### internal/giturl
- **Purpose**: Parse and transform platform-specific git URLs
- **Providers**: GitHub, GitLab (extensible for Gitea, Bitbucket, Azure)
- **Key Functions**:
  - `AutoDetect()`: Identify provider from URL
  - `Raw()`: Transform to raw content URL for direct download

#### internal/download
- **Purpose**: HTTP-based raw content download
- **Use**: Short-circuit git operations when possible (GitHub/GitLab)

## Development Workflows

### Building and Testing

```bash
# Run all tests with race detection and coverage
go test -v -race -cover -coverprofile=coverage.out -covermode=atomic ./...

# Run specific test
go test -v -run TestFetcher ./...

# Run linting
golangci-lint run

# Check for vulnerabilities
govulncheck ./...
```

### Code Quality Standards

**Linting Configuration** (`.golangci.yml`):
- Uses golangci-lint v2 with "all" linters enabled
- Disabled linters: See `.golangci.yml` lines 4-38 for rationale
- Key enabled checks: errcheck, govet, staticcheck, gosimple, ineffassign
- Cyclomatic complexity threshold: 45
- No limit on issues per linter (thorough checking)

**Testing Strategy**:
- Parallel tests when possible (`t.Parallel()`)
- Use testify/require for assertions
- Mock interfaces defined in `mocks_test.go`
- Test matrix: Ubuntu/Windows × Go oldstable/stable
- Race detection enabled in CI
- Coverage reporting via Coveralls

### Git Workflow

**Branch Strategy**:
- Main branch: `master` (default for PRs)
- Feature branches: `claude/claude-md-*` format for AI-assisted development
- Always develop on designated feature branch, never push directly to master

**Commit Guidelines**:
- Use conventional commit style (implied from project)
- Reference issue numbers where applicable
- Keep commits focused and atomic
- SPDX headers required in all source files

### CI/CD Pipeline

**GitHub Actions Workflows**:
1. **Linting** (`01-golang-lint.yaml`): golangci-lint checks
2. **Tests** (`02-test.yaml`): Cross-platform testing with coverage
3. **Security** (`03-govulncheck.yaml`): Vulnerability scanning

All workflows trigger on push and pull_request events.

## Key Conventions

### 1. File Headers
All source files must include SPDX headers:
```go
// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0
```

### 2. Error Handling
- All VCS-related errors should wrap `ErrVCS` using `errors.Join()` or `fmt.Errorf("%w", ErrVCS)`
- Use `fmt.Errorf()` for error wrapping with context
- Example: `fmt.Errorf("could not fetch: %w: %w", err, ErrVCS)`

### 3. Options Pattern
- Functional options pattern used throughout (see `options.go`)
- Options types: `FetchOption`, `CloneOption`, `SPDXOption`, `GitLocatorOption`
- Generic `optionsWithDefaults()` function applies defaults and user options
- Example: `NewFetcher(FetchWithBackingDir(true, "/tmp"), FetchWithExactTag(true))`

### 4. Context Usage
- All I/O operations accept `context.Context` as first parameter
- Respect context cancellation in long-running operations
- Pass context through call chain consistently

### 5. Concurrency
- `Fetcher`: Safe for concurrent use (stateless)
- `Cloner`: NOT safe for concurrent use (stateful)
- Document thread safety in godoc comments

### 6. Internal vs Public API
- `internal/` packages are implementation details, never exported
- Public API is minimal: `Fetcher`, `Cloner`, `Locator` interface, and option types
- Keep internal complexity hidden from users

### 7. Version Resolution
- **Default**: HEAD of default branch if no version specified
- **Semver**: Incomplete semver (v2, v2.1) resolves to latest compatible version
- **Exact Tag**: Can be enforced via `FetchWithExactTag(true)` option
- **Pre-releases**: Excluded by default, enable with `FetchWithAllowPrereleases(true)`

### 8. Testing Patterns
```go
// Parallel tests
func TestFoo(t *testing.T) {
    t.Parallel()
    // ...
}

// Table-driven tests (use when appropriate)
tests := []struct{
    name string
    input string
    want string
}{ /* ... */ }

// Use testify/require for assertions
require.NoError(t, err)
require.Equal(t, expected, actual)
```

### 9. Documentation
- Godoc comments on all exported types, functions, and methods
- Examples in godoc format (see pkg.go.dev examples)
- README.md kept up-to-date with feature status
- TODO comments for incomplete features

## Common Development Tasks

### Adding a New SCM Provider

1. Create new package under `internal/giturl/<provider>/`
2. Implement parser and raw URL generator
3. Register provider in `internal/giturl/providers.go`
4. Add tests following existing patterns (see `github/parser_test.go`)
5. Update README.md feature list

### Adding New Fetch/Clone Options

1. Define option function type in `options.go`
2. Add to appropriate options struct (`fetchOptions`, `cloneOptions`, etc.)
3. Implement internal option setter
4. Map to internal package options in `toInternal*Options()` methods
5. Document in godoc with examples

### Debugging Git Operations

- Enable debug logging: `FetchWithGitDebug(true)` or `CloneWithGitDebug(true)`
- Debug output goes to `log.Printf` (stdout)
- Shows git commands executed and go-git operations

### Performance Optimization

- Native git binary used automatically when detected (faster than go-git)
- Raw content URLs used for GitHub/GitLab (bypasses git operations)
- Disable autodetection: `FetchWithGitSkipAutoDetect(true)` for pure go-git
- Memory-backed by default, use `FetchWithBackingDir(true, dir)` for disk

## Important Notes for AI Assistants

### When Modifying Code

1. **Always preserve SPDX headers** in existing files
2. **Run tests** after changes: `go test ./...`
3. **Check linting**: `golangci-lint run`
4. **Update godoc** for API changes
5. **Add tests** for new functionality
6. **Consider backward compatibility** (public API is stable)

### When Adding Features

1. Check if it aligns with project scope (VCS fetching/cloning)
2. Implement in internal package first, expose via options
3. Follow existing patterns (options, error handling, context)
4. Add examples if adding public API
5. Update README.md feature checklist

### When Fixing Bugs

1. Add failing test that reproduces the issue
2. Fix the bug
3. Verify test passes
4. Check for similar issues in related code
5. Consider if fix affects public API behavior (document if yes)

### Code Style Notes

- Use `go fmt` and `goimports` (enforced by CI)
- Prefer table-driven tests for multiple scenarios
- Keep functions focused and reasonably sized
- Use meaningful variable names (no excessive abbreviation)
- Comment non-obvious logic, especially in internal packages

### Dependencies

- **Primary**: `github.com/go-git/go-git/v5` (pure Go git implementation)
- **Semver**: `github.com/blang/semver/v4` (semantic version parsing)
- **Testing**: `github.com/stretchr/testify` (assertions)
- **Go Version**: 1.24.0+ (see `go.mod`)

Keep dependencies minimal. Prefer standard library when possible.

## Troubleshooting Common Issues

### Tests Failing on Windows
- Check file path separators (use `filepath.Join`)
- Be aware of line ending differences (CRLF vs LF)
- Test matrix includes Windows, so CI will catch platform issues

### Git Operations Slow
- Native git binary may not be detected
- Check `capabilities.go` logic for your platform
- Try enabling debug mode to see if native git is used

### SPDX Locator Parsing Fails
- Ensure URL has fragment (`#file.txt`)
- Check URL format matches spec in `spdx.go` godoc
- Try `GitLocator` as fallback (auto-detected in `Fetcher.FetchURL`)

### Clone vs Fetch Confusion
- Use `Fetcher` for single file
- Use `Cloner` for multiple files from same repo or entire repo
- `Cloner` maintains state, `Fetcher` does not

## Related Resources

- [SPDX Package Download Location Spec](https://spdx.github.io/spdx-spec/v2.3/package-information/#77-package-download-location-field)
- [Git URL Formats](https://git-scm.com/docs/git-fetch#_git_urls)
- [go-git Documentation](https://pkg.go.dev/github.com/go-git/go-git/v5)
- [Project GoDoc](https://pkg.go.dev/github.com/fredbi/go-vcsfetch)

---

**Last Updated**: 2025-12-10
**Go Version**: 1.24.0
**License**: Apache-2.0
