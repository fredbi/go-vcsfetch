# go-vcsfetch

<!-- Badges: status  -->
[![Tests][test-badge]][test-url] [![Coverage][cov-badge]][cov-url] [![CI vuln scan][vuln-scan-badge]][vuln-scan-url] [![CodeQL][codeql-badge]][codeql-url]
<!-- Badges: release & docker images  -->
<!-- Badges: code quality  -->
<!-- Badges: license & compliance -->
[![Release][release-badge]][release-url] [![Go Report Card][gocard-badge]][gocard-url] [![CodeFactor Grade][codefactor-badge]][codefactor-url] [![License][license-badge]][license-url]
<!-- Badges: documentation & support -->
<!-- Badges: others & stats -->
[![GoDoc][godoc-badge]][godoc-url] [![go version][goversion-badge]][goversion-url] ![Top language][top-badge] ![Commits since latest release][commits-badge]

vcs fetcher and cloner for Go.

A Go library for fetching files from version control systems (vcs).

---

Easily retrieve individual files or repositories over a vcs location.

* [x] Support `git` repositories
* [x] Support SPDX Locators (spdx downloadLocation attribute)
* [x] Support common `git-url` schemes

All fetched resources are exposed for read-only operations only.

If you're looking for general purpose vcs support in Go for read/write or other git-heavy operations,
consider using `github.com/go-git/go-git` instead.

## Status

Work in progress. Unreleased.

## Use-cases

* retrieve a single file over a remote repo (e.g. config file)
* retrieve an entire folder at a specific version
* ...

**Not intended to work with local resources** (e.g. `file://...`).

## Features

**VCS (Version Control System)**

* [x] Works without git installed
* [x] Supported schemes: http, https, ssh, git TCP
* [x] Authentication (basic, ssh)
* [x] `Fetch` (single file) or `Clone` (folder or entire repo)
* [x] `Fetch` optimized for common SCMs (github.com, gitlab), with https raw content download to bypass pure-git operations
* [x] In memory or filesystem-backed
* [x] Supports sparse-cloning
* [x] Auto-detects the presence of the `git` binary for faster fetching using the `git` command line

**Resolving versions**

* [x] Ref as commit sha, branch or tag, with exact match
* [x] Semver tag resolution with incomplete semver: e.g. resolve `v2` as the latest tag `<v3`,
      and `2.1` as the latest tag `<v2.2`

**SCM-specific URLs**

* [x] `git-url` parses resource locators for well-known schemes
  * [ ] azure
  * [ ] bitbucket
  * [ ] gitea
  * [x] github
  * [x] gitlab
* [x] know how to transform a resource locator into a raw-content URL

## Quick start

```cmd
go get github.com/fredbi/go-vcsfetch
```

## Usage

### Basic usage

```go
import (
    "bytes"
    "context"
    "log"

    "github.com/fredbi/go-vcsfetch"
)

...

vf := vcsfetch.NewFetcher()
w := new(bytes.Buffer)
ctx := context.Background()

const spdxDownloadLocation = "https://github.com/fredbi/go-vcsfetch@HEAD#.golangci.yml"

if err := vf.Fetch(ctx, w, spdxDownloadLocation); err != nil {
    ...
}

log.Println(w.String())
```

### Advanced usage with options

Example use cases:

1. authentication
2. git over TCP
3. repo cloning
4. git-urls
5. folder retrieval and repeated fetches
6. exact vs semver tag resolution
7. using shorthand slugs
8. git-archive (with benchmark)
9. TLS settings

Take a tour of the [![Live examples](https://img.shields.io/badge/Live%20Examples-blue)](https://pkg.go.dev/github.com/fredbi/go-vcsfetch#pkg-examples).

## Dependencies

This library is built on top of `github.com/go-git/go-git`, a pure Go git implementation.
It does not require runtime dependencies (e.g. like when using go-git bindings from `git2go`).

It does not require the `git` binary to be installed.

However, when the `git` binary is present and auto-detection is not disabled, the library may chose to perform
some operations using the native `git` implementation, which is usually faster than the native go port.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-vcsfetch.svg)](https://pkg.go.dev/github.com/fredbi/go-vcsfetch).

## Resource usage and performances

TODO

### Roadmap

* [ ] Support for `git-archive` download, when well-known SCM will start support this protocol
* [ ] Support for mercurial, with a runtime dependency on `hg`. 
* [ ] native go git-archive support (or from go-git/v6?)
* [ ] support semver version constraint such as `^v1.2.3` or `~v1.2.3`
* [ ] mock git server

## License

This library is distributed under the Apache 2.0 license.

`SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON`
[`SPDX-License-Identifier: Apache-2.0`](./LICENSE.md)

## Credits and acknowledgments

Initially, my intent was to enable a shared golangci-lint config file
on a repository common to all repos within an organization.

> Doing a little research on how to work with vcs resources in Go, I stumbled over this little package
> [github.com/carabiner-dev/vcslocator](https://pkg.go.dev/github.com/carabiner-dev/vcslocator).
>
> I started to use it right away, but was quickly hindered by a number of limitations.
> My requirements departed quite a bit from that implementation, to the point that forking wasn't an option.
> And so started this implementation.

Thank you to the guys at [`carabiner-dev`](https://carabiner.dev/), who provided me the inspiration
to use a SPDX locator on top of `go-git`.

Notice that this implementation is 100% original code and not a plagiarism of the above.

<!-- Badges: status  -->
[test-badge]: https://github.com/fredbi/go-vcsfetch/actions/workflows/go-test.yml/badge.svg
[test-url]: https://github.com/fredbi/go-vcsfetch/actions/workflows/go-test.yml
[cov-badge]: https://codecov.io/gh/fredbi/go-vcsfetch/branch/master/graph/badge.svg
[cov-url]: https://codecov.io/gh/fredbi/go-vcsfetch
[vuln-scan-badge]: https://github.com/fredbi/go-vcsfetch/actions/workflows/scanner.yml/badge.svg
[vuln-scan-url]: https://github.com/fredbi/go-vcsfetch/actions/workflows/scanner.yml
[codeql-badge]: https://github.com/fredbi/go-vcsfetch/actions/workflows/codeql.yml/badge.svg
[codeql-url]: https://github.com/fredbi/go-vcsfetch/actions/workflows/codeql.yml
<!-- Badges: release & docker images  -->
[release-badge]: https://badge.fury.io/gh/fredbi%2Fgo-vcsfetch.svg
[release-url]: https://badge.fury.io/gh/fredbi%2Fgo-vcsfetch
[gomod-badge]: https://badge.fury.io/go/github.com%2Ffredbi%2Fgo-vcsfetch.svg
[gomod-url]: https://badge.fury.io/go/github.com%2Ffredbi%2Fgo-vcsfetch
<!-- Badges: code quality  -->
[gocard-badge]: https://goreportcard.com/badge/github.com/fredbi/go-vcsfetch
[gocard-url]: https://goreportcard.com/report/github.com/fredbi/go-vcsfetch
[codefactor-badge]: https://img.shields.io/codefactor/grade/github/fredbi/go-vcsfetch
[codefactor-url]: https://www.codefactor.io/repository/github/fredbi/go-vcsfetch
<!-- Badges: documentation & support -->
[doc-badge]: https://img.shields.io/badge/doc-site-blue?link=https%3A%2F%2Fgoswagger.io%2Ffredbi%2F
[doc-url]: https://goswagger.io/fredbi
[godoc-badge]: https://pkg.go.dev/badge/github.com/fredbi/go-vcsfetch
[godoc-url]: http://pkg.go.dev/github.com/fredbi/go-vcsfetch
<!-- Badges: license & compliance -->
[license-badge]: http://img.shields.io/badge/license-Apache%20v2-orange.svg
[license-url]: https://github.com/fredbi/go-vcsfetch/?tab=Apache-2.0-1-ov-file#readme
<!-- Badges: others & stats -->
[goversion-badge]: https://img.shields.io/github/go-mod/go-version/fredbi/go-vcsfetch
[goversion-url]: https://github.com/fredbi/go-vcsfetch/blob/master/go.mod
[top-badge]: https://img.shields.io/github/languages/top/fredbi/go-vcsfetch
[commits-badge]: https://img.shields.io/github/commits-since/fredbi/go-vcsfetch/latest
