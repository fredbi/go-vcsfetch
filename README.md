# go-vcsfetch
![Lint](https://github.com/fredbi/go-vcsfetcher/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/go-vcsfetcher/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/go-vcsfetcher/badge.svg?branch=master)](https://coveralls.io/github/fredbi/go-vcsfetcher?branch=master)
![Vulnerability Check](https://github.com/fredbi/go-vcsfetcher/actions/workflows/03-govulncheck.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/go-vcsfetcher)](https://goreportcard.com/report/github.com/fredbi/go-vcsfetcher)

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/go-vcsfetcher)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://raw.githubusercontent.com/fredbi/go-vcsfetcher/master/LICENSE.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-vcsfetcher.svg)](https://pkg.go.dev/github.com/fredbi/go-vcsfetcher)
![Go version](https://img.shields.io/github/go-mod/go-version/fredbi/go-vcsfetcher?color=violet)
![Top language](https://img.shields.io/github/languages/top/fredbi/uri?color=green)

A vcs fetcher and cloner for go.

Easily retrieve individual files or repositories over a vcs location.

* [x] Support `git` repositories
* [x] Support SPDX Locators (spdx downloadLocation attribute)
* [x] Support common `git-url` schemes

All fetched resources are exposed for read-only operations.

If you're looking for general purpose vcs support for go,
consider using `github.com/go-git/go-git` instead.

## Use-cases

* retrieve a single file over a remote repo (e.g. config file)
* retrieve an entire folder at a specific version
* ...

## Features

**VCS (Version Control System)**

* [x] No required runtime dependency (works without git installed)
* [x] Supported schemes: http, https, ssh, git TCP
* [x] Authentication (basic, ssh)
* [x] `Fetch` (single file) or `Clone` (folder or entire repo)
* [x] In memory or filesystem-backed
* [x] `git-url` for well-known schemes by gihub, gitlab, gitea.
* [x] Sparse-cloning
* [x] Auto-detects the presence of the `git` binary for faster fetching using `git-archive`

**Resolving versions**

* [x] Ref as commit sha, branch or tag, with exact match
* [x] Semver tag resolution with incomplete semver: e.g. resolve `v2` as the latest tag `<v3`,
      and `2.1` as the latest tag `<v2.2`

### Future developments

* [ ] Support for mercurial, with a runtime dependency on `hg` 
* [ ] native go git-archive support (or from go-git/v6?)
* [ ] mock git server

## Importing this library in your project

```cmd
go get github.com/fredbi/go-vcsfetcher
```

## Usage

### Basic usage

```go
import (
    "bytes"
    "context"
    "log"

    "github.com/fredbi/go-vcsfetcher"
)

...

vf := vcsfetcher.NewFetcher()
w := new(bytes.Buffer)
ctx := context.Background()

const spdxDownloadLocation = "https://github.com/fredbi/go-vcsfetcher@HEAD#.golangci.yml"

if err := vf.Fetch(ctx, w, spdxDownloadLocation) ; err != nil {
    ...
}

log.Println(w.String())
```

### Advanced usage with options

Examplified use cases:

1. authentication
2. git over TCP
3. repo cloning
4. git-urls
5. folder retrieval and repeated fetches
6. exact vs semver tag resolution
7. using shorthand slugs
8. git-archive (with benchmark)
9. TLS settings

Take a tour of the ![[Live examples](https://img.shields.io/badge/Live%20Examples-blue)](https://pkg.go.dev/github.com/fredbi/go-vcsfetcher#examples).

## Dependencies

This library is built on top `github.com/go-git/go-git`, a pure go git implementation.
It does not require runtime dependencies (e.g. like when using go-git bindings from `git2go`).

It does not require the `git` binary to be installed.

However, when the `git` binary is present and auto-detection is not disabled, the library may chose to perform
some operations using the native `git` implementation, which is usually faster than the native go port.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-vcsfetcher.svg)](https://pkg.go.dev/github.com/fredbi/go-vcsfetcher).

## Resource usage and performances

TODO

## License

This library is distributed under the Apache 2.0 license.

`SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON`
[`SPDX-License-Identifier: Apache-2.0`](./LICENSE.md)

## Credits and acknowledgments

Initially, my intent was to enable a shared golangci-lint config file
on a repository common to all repos within an organization.

Doing a little research on how to work with vcs resources in go, I stumbled over this little package
[github.com/carabiner-dev/vcslocator](https://pkg.go.dev/github.com/carabiner-dev/vcslocator).

I started to use it righ away, but quickly stumbled on a number of limitations.
My requirements departed quite a bit from that implementation, to the point that forking wasn't an option.
And so started this implementation.

Thank you to the guys at [`carabiner-dev`](https://carabiner.dev/), who provided me the inspiration
to use a SPDX locator on top of `go-git`.

Notice that this implementation is 100% original code and not a plagiarism of the above.
