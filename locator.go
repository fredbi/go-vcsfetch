// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

//go: generate mockery

package vcsfetch

import "net/url"

// Locator is the interface for types that know how to resolve a vcs URL.
//
// This package currently exposes two implementations: [SPDXLocator] and [GitLocator].
//
// Users of the [Fetcher] and the [Cloner] may implement a custom [Locator] to meet special requirements.
type Locator interface {
	// RepoURL yields the base URL of the vcs repository,
	// e.g. https://github.com/fredbi/go-vcsfetcher
	RepoURL() *url.URL

	// Version yields the ref identifying the desired version of a file, e.g. v0.0.1
	Version() string

	// Path yields the file path relative to the repository,
	// e.g. internal/git/api.go
	Path() string

	// IsLocal indicates if the repository is local,
	// e.g. the URL looks like file://src/fred/github.com/fredbi/go-vcsfetcher
	IsLocal() bool

	// HasAuth indicates if the [Locator] embeds some credentials,
	// e.g. the URL looks like https://fredbi:token@github.com/fredbi/go-vcsfetcher
	HasAuth() bool

	String() string
}
