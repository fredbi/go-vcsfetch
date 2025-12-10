// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package bitbucket

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Locator redefines locally the common minimal locator interface.
//
// This avoids cross-dependencies between repos.
//
// This package exposes [URL] as an implementation for bitbucket.
type Locator interface {
	RepoURL() *url.URL
	Path() string
	Version() string
}

// Raw returns the raw content URL for a [Locator] hosted on Bitbucket.
//
// Only https URL's are supported.
//
// For self-hosted Bitbucket Server instances, this only works for instances
// accessible via standard https (port 443 or unspecified).
//
// Examples:
//
//   - https://bitbucket.org/workspace/repo/raw/master/README.md
//   - https://bitbucket.org/atlassian/python-bitbucket/raw/main/setup.py
func Raw(locator Locator) (*url.URL, error) {
	repo := locator.RepoURL()
	pth := strings.Trim(locator.Path(), "/")
	if pth == "" {
		return nil, fmt.Errorf("returning a raw content url requires a non empty path to a file: %w", ErrBitbucket)
	}

	version := locator.Version()
	if version == "" {
		version = "HEAD"
	}

	scheme, _ := strings.CutSuffix(repo.Scheme, "+git")

	if scheme != "https" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL scheme: %w", ErrBitbucket)
	}

	if port := repo.Port(); port != "" && port != "443" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL with standard port (443 or unspecified): %w", ErrBitbucket)
	}

	u := &url.URL{}
	*u = *repo // shallow clone

	// Bitbucket raw URL format: /{workspace}/{repo}/raw/{ref}/{path}
	u.Path = path.Join(u.Path, "raw", version, pth)
	u.Fragment = ""
	u.RawFragment = ""

	return u, nil
}
