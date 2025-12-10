# Bitbucket URL Parser

Implementation of Bitbucket URL parsing for the `go-vcsfetch` library.

## Supported URL Formats

### Repository URL
```
https://bitbucket.org/{workspace}/{repo}
https://bitbucket.org/{workspace}/{repo}.git
```

### Browse URLs
```
https://bitbucket.org/{workspace}/{repo}/src/{ref}/{path}
```

### Raw Content URLs
```
https://bitbucket.org/{workspace}/{repo}/raw/{ref}/{path}
```

## Key Differences from Other Platforms

### Terminology
Bitbucket uses **"workspace"** instead of "owner" or "organization".

### Simpler URL Structure
Unlike GitHub (which uses `/blob/{ref}`) or Gitea (which uses `/src/branch/{ref}`),
Bitbucket has a simpler format where the ref comes directly after `src` or `raw`:

| Platform | Browse URL Pattern |
|----------|-------------------|
| GitHub   | `/blob/{ref}/{path}` |
| Gitea    | `/src/branch/{ref}/{path}` |
| **Bitbucket** | **`/src/{ref}/{path}`** |

### No Ref Type Discriminator
Bitbucket doesn't require specifying whether the ref is a branch, tag, or commit.
The same URL structure works for all:

```
/src/master/file.txt          (branch)
/src/v1.0.0/file.txt          (tag)
/src/abc123def456/file.txt    (commit sha)
```

## Examples

### Parse a Bitbucket browse URL
```go
u, _ := url.Parse("https://bitbucket.org/workspace/repo/src/master/README.md")
loc, err := bitbucket.Parse(u)
// loc.RepoURL() => https://bitbucket.org/workspace/repo
// loc.Version() => master
// loc.Path()    => README.md
```

### Parse a raw content URL
```go
u, _ := url.Parse("https://bitbucket.org/atlassian/python-bitbucket/raw/main/setup.py")
loc, err := bitbucket.Parse(u)
// Automatically recognizes raw URLs
```

### Generate a raw content URL
```go
rawURL, err := bitbucket.Raw(loc)
// rawURL => https://bitbucket.org/workspace/repo/raw/master/README.md
```

## Bitbucket Server (Self-Hosted)

This parser works with Bitbucket Server (self-hosted instances):

```go
u, _ := url.Parse("https://bitbucket.example.com/workspace/project/src/develop/code.js")
loc, err := bitbucket.Parse(u)
// Works seamlessly with custom domains
```

**Note:** Bitbucket Server may have different URL patterns than Bitbucket Cloud.
This implementation follows Bitbucket Cloud conventions.

## Real-World Examples

### Atlassian's Python Bitbucket Library
```
Browse: https://bitbucket.org/atlassian/python-bitbucket/src/main/pybitbucket/auth.py
Raw:    https://bitbucket.org/atlassian/python-bitbucket/raw/main/pybitbucket/auth.py
```

### With Commit SHA
```
Browse: https://bitbucket.org/workspace/repo/src/abc123def456789/src/main.go
Raw:    https://bitbucket.org/workspace/repo/raw/abc123def456789/src/main.go
```

### With Version Tag
```
Browse: https://bitbucket.org/workspace/repo/src/v2.1.3/CHANGELOG.md
Raw:    https://bitbucket.org/workspace/repo/raw/v2.1.3/CHANGELOG.md
```

## API Access

Bitbucket also provides an API for file access:
```
https://api.bitbucket.org/2.0/repositories/{workspace}/{repo}/src/{ref}/{path}
```

This implementation focuses on web URLs, not API endpoints.

## Testing

Run tests:
```bash
go test ./internal/giturl/bitbucket/...
```

Test coverage includes:
- Repository-only URLs
- Browse and raw URLs
- Tags, branches, and commit SHAs
- Nested file paths
- Self-hosted instances
- Error cases (invalid URLs)
