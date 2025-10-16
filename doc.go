// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

// Package vcsfetch provides a vcs fetcher and cloner for go.
//
// # URL formats for vcs locations
//
// We recommend the SPDX format, which is standardized and unambiguous.
// SPDX URLs must contain an URL fragment.
//
// [Fetcher] and [Cloner] also support well-known git-url schemes exposed by git platforms such as
// github, gitlab and gitea.
//
// URL shorthands using repo slugs: TODO
//
// # Supported vcs protocols
//
// Both the [Fetcher] and the [Cloner] come with native support for git, with no runtime dependencies.
// Supported transports for git  are: file, https, ssh and git over TCP.
//
// NOTES:
//
//   - http is also supported (e.g. for testing).
//   - git over TCP is not supported as a SPDX locator (TODO: check this)
//
// # Limitations
//
// At this moment, this package does not support mercurial ("hg"). We may add this feature later on,
// as mercurial is supported by go.
//
// [Fetcher] and [Cloner] do not support bazar ("bzr") or subversion ("svn"), and we currently have no plan
// to add support for those.
//
// # Versions
//
// Versions may specify a given commit sha (full or short sha), a branch (resolves as the HEAD of that
// branch) or a tag.
//
// The symbolic ref "HEAD" resolves as the HEAD of the default branch.
//
// Semver tags may be incomplete to refer to the latest semver with a given major or minor version:
//
// - v2 resolves as the latest v2.x.y tag (i.e. <v3)
// - v2.1 resolves as the latest v2.1.y tag (i.e. <v2.2)
//
// This behavior may be disabled with [WithExactTag].
//
// If no version information is provided, the default reference is the HEAD commit of the default branch
// (e.g. master or main).
//
// This behavior may be disabled with [WithRequireVersion].
//
// # Authentication
//
//   - TLS
//   - proxy
//
// TODO
package vcsfetch
