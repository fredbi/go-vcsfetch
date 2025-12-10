# Gitea URL Parser

Implementation of Gitea URL parsing for the `go-vcsfetch` library.

## Supported URL Formats

### Repository URL
```
https://gitea.com/{owner}/{repo}
https://gitea.com/{owner}/{repo}.git
```

### Browse URLs
```
https://gitea.com/{owner}/{repo}/src/branch/{branch-name}/{path}
https://gitea.com/{owner}/{repo}/src/tag/{tag-name}/{path}
https://gitea.com/{owner}/{repo}/src/commit/{commit-sha}/{path}
```

### Raw Content URLs
```
https://gitea.com/{owner}/{repo}/raw/branch/{branch-name}/{path}
https://gitea.com/{owner}/{repo}/raw/tag/{tag-name}/{path}
https://gitea.com/{owner}/{repo}/raw/commit/{commit-sha}/{path}
```

## Examples

### Parse a Gitea browse URL
```go
u, _ := url.Parse("https://gitea.com/owner/repo/src/branch/master/README.md")
loc, err := gitea.Parse(u)
// loc.RepoURL() => https://gitea.com/owner/repo
// loc.Version() => master
// loc.Path()    => README.md
```

### Generate a raw content URL
```go
rawURL, err := gitea.Raw(loc)
// rawURL => https://gitea.com/owner/repo/raw/branch/master/README.md
```

## Self-Hosted Gitea Instances

This parser works with any Gitea instance, not just gitea.com:

```go
u, _ := url.Parse("https://git.example.com/org/project/src/branch/main/file.go")
loc, err := gitea.Parse(u)
// Works seamlessly with custom domains
```

## Differences from GitHub

While Gitea is based on GitHub's design, the URL structure differs:

| Platform | Browse URL Pattern |
|----------|-------------------|
| GitHub   | `/blob/{ref}/{path}` |
| Gitea    | `/src/branch/{ref}/{path}` |

| Platform | Raw URL Pattern |
|----------|-----------------|
| GitHub   | `raw.githubusercontent.com/{owner}/{repo}/{ref}/{path}` |
| Gitea    | `/{owner}/{repo}/raw/branch/{ref}/{path}` |

## Testing

Run tests:
```bash
go test ./internal/giturl/gitea/...
```

All tests include parallel execution for performance.
