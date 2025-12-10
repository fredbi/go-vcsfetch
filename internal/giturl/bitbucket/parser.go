// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package bitbucket

import (
	"fmt"
	"net/url"
	"strings"
)

// URL is a bitbucket-style URL to a vcs resource hosted by Bitbucket SCM.
type URL struct {
	repoURL *url.URL
	path    string
	version string
}

const (
	defaultScheme = "https"
	defaultHost   = "bitbucket.org"
)

// Parse a bitbucket URL.
//
// Bitbucket URL formats:
//   - Browse: https://bitbucket.org/{workspace}/{repo}/src/{ref}/{path}
//   - Raw: https://bitbucket.org/{workspace}/{repo}/raw/{ref}/{path}
//   - Repo: https://bitbucket.org/{workspace}/{repo}
//
// Note: Bitbucket uses "workspace" terminology instead of "owner".
func Parse(bitbucketURL *url.URL) (*URL, error) {
	u := &url.URL{}
	*u = *bitbucketURL // shallow clone

	if u.Scheme == "" {
		u.Scheme = defaultScheme
	}

	if u.Hostname() == "" {
		if u.Port() == "" {
			u.Host = defaultHost
		} else {
			u.Host = defaultHost + ":" + u.Port()
		}
	}

	u.Host = strings.ToLower(u.Host)
	pth := strings.Trim(u.Path, "/")

	const (
		repoIndex = 2
	)

	parts := strings.Split(pth, "/")
	if len(parts) < repoIndex {
		return nil, fmt.Errorf("expected the URL path component to contain at least %d parts, but got %q: %w", repoIndex, pth, ErrBitbucket)
	}

	repo := strings.Join(parts[:repoIndex], "/")
	repo = strings.TrimSuffix(repo, ".git")
	u.Path = repo

	if len(parts) == repoIndex {
		// entire repo
		u.RawFragment = ""
		u.Fragment = ""
		u.RawQuery = ""

		bb := &URL{
			repoURL: u,
			path:    "/",
			version: "",
		}

		return bb, nil
	}

	parts = parts[repoIndex:]

	// Bitbucket uses "src" or "raw" as first part (simpler than gitea - no branch/tag/commit discriminator)
	const neededPartsAfterRepo = 2
	if len(parts) < neededPartsAfterRepo {
		return nil, fmt.Errorf(`expected URL path to contain at least %d parts after repo but got %q: %w`, neededPartsAfterRepo, pth, ErrBitbucket)
	}

	discriminator := strings.ToLower(parts[0])
	switch discriminator {
	case "src":
		// Browse URL: /src/{ref}/{path}
	case "raw":
		// Raw URL: /raw/{ref}/{path}
	default:
		return nil, fmt.Errorf(`expected URL path to contain "src" or "raw" but got %q in %q: %w`, parts[0], pth, ErrBitbucket)
	}

	parts = parts[1:]

	// Next part is the ref (branch, tag, or commit - Bitbucket doesn't require specifying which)
	if len(parts) < 1 {
		return nil, fmt.Errorf(`expected URL path to contain ref name but got %q: %w`, pth, ErrBitbucket)
	}

	ref := parts[0]
	parts = parts[1:]

	if len(parts) == 0 {
		// No file path - this is a tree/directory view
		parts = []string{"/"}
	}

	repoPath := strings.Join(parts, "/")
	u.RawFragment = ""
	u.Fragment = ""
	u.RawQuery = ""

	bb := &URL{
		repoURL: u,
		path:    repoPath,
		version: ref,
	}

	return bb, nil
}

// RepoURL yields the base URL of the vcs repository,
// e.g. https://bitbucket.org/workspace/repo
func (bb *URL) RepoURL() *url.URL {
	return bb.repoURL
}

// Version yields the ref identifying the desired version of a file,
// e.g. "master" in https://bitbucket.org/workspace/repo/src/master/README.md
func (bb *URL) Version() string {
	return bb.version
}

// Path yields the file path relative to the repository,
// e.g. "README.md" in https://bitbucket.org/workspace/repo/src/master/README.md
func (bb *URL) Path() string {
	return bb.path
}
