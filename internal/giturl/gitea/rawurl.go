// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package gitea

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
// This package exposes [URL] as an implementation for gitea.
type Locator interface {
	RepoURL() *url.URL
	Path() string
	Version() string
}

// Raw returns the raw content URL for a [Locator] hosted on a Gitea instance.
//
// Only https URL's are supported.
//
// For self-hosted instances, this only works for instances accessible via
// standard https (port 443 or unspecified).
//
// Examples:
//
//   - https://gitea.com/fredbi/go-vcsfetch/raw/branch/master/README.md
//   - https://try.gitea.io/owner/repo/raw/branch/main/file.txt
func Raw(locator Locator) (*url.URL, error) {
	repo := locator.RepoURL()
	pth := strings.Trim(locator.Path(), "/")
	if pth == "" {
		return nil, fmt.Errorf("returning a raw content url requires a non empty path to a file: %w", ErrGitea)
	}

	version := locator.Version()
	if version == "" {
		version = "HEAD"
	}

	scheme, _ := strings.CutSuffix(repo.Scheme, "+git")

	if scheme != "https" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL scheme: %w", ErrGitea)
	}

	if port := repo.Port(); port != "" && port != "443" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL with standard port (443 or unspecified): %w", ErrGitea)
	}

	u := &url.URL{}
	*u = *repo // shallow clone

	// Gitea raw URL format: /{owner}/{repo}/raw/branch/{ref}/{path}
	u.Path = path.Join(u.Path, "raw", "branch", version, pth)
	u.Fragment = ""
	u.RawFragment = ""

	return u, nil
}
