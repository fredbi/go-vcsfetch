// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package gitea

import (
	"fmt"
	"net/url"
	"strings"
)

// URL is a gitea-style URL to a vcs resource hosted by gitea SCM.
type URL struct {
	repoURL *url.URL
	path    string
	version string
}

const (
	defaultScheme = "https"
	defaultHost   = "gitea.com"
)

// Parse a gitea URL.
//
// Gitea URL formats:
//   - Browse: https://gitea.com/{owner}/{repo}/src/branch/{ref}/{path}
//   - Raw: https://gitea.com/{owner}/{repo}/raw/branch/{ref}/{path}
//   - Repo: https://gitea.com/{owner}/{repo}
func Parse(giteaURL *url.URL) (*URL, error) {
	u := &url.URL{}
	*u = *giteaURL // shallow clone

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
		return nil, fmt.Errorf("expected the URL path component to contain at least %d parts, but got %q: %w", repoIndex, pth, ErrGitea)
	}

	repo := strings.Join(parts[:repoIndex], "/")
	repo = strings.TrimSuffix(repo, ".git")
	u.Path = repo

	if len(parts) == repoIndex {
		// entire repo
		u.RawFragment = ""
		u.Fragment = ""
		u.RawQuery = ""

		gt := &URL{
			repoURL: u,
			path:    "/",
			version: "",
		}

		return gt, nil
	}

	parts = parts[repoIndex:]

	var (
		ref    string
		isTree bool
	)

	// Gitea uses "src" or "raw" as first part
	const neededPartsAfterRepo = 2
	if len(parts) < neededPartsAfterRepo {
		return nil, fmt.Errorf(`expected URL path to contain at least %d parts after repo but got %q: %w`, neededPartsAfterRepo, pth, ErrGitea)
	}

	discriminator := strings.ToLower(parts[0])
	switch discriminator {
	case "src":
		// Browse URL: /src/branch/{ref}/{path}
	case "raw":
		// Raw URL: /raw/branch/{ref}/{path}
	default:
		return nil, fmt.Errorf(`expected URL path to contain "src" or "raw" but got %q in %q: %w`, parts[0], pth, ErrGitea)
	}

	parts = parts[1:]

	// Next part should be "branch", "tag", "commit"
	if len(parts) < 2 {
		return nil, fmt.Errorf(`expected URL path to contain ref type (branch/tag/commit) and ref name but got %q: %w`, pth, ErrGitea)
	}

	refType := strings.ToLower(parts[0])
	switch refType {
	case "branch", "tag", "commit":
		// Valid ref types
		ref = parts[1]
		parts = parts[2:]
	default:
		return nil, fmt.Errorf(`expected URL path to contain "branch", "tag", or "commit" but got %q in %q: %w`, parts[0], pth, ErrGitea)
	}

	if len(parts) == 0 {
		// No file path - this is a tree view
		isTree = true
		parts = []string{"/"}
	}

	repoPath := strings.Join(parts, "/")
	u.RawFragment = ""
	u.Fragment = ""
	u.RawQuery = ""

	gt := &URL{
		repoURL: u,
		path:    repoPath,
		version: ref,
	}

	_ = isTree // may be used for validation in the future

	return gt, nil
}

// RepoURL yields the base URL of the vcs repository,
// e.g. https://gitea.com/fredbi/go-vcsfetch
func (gt *URL) RepoURL() *url.URL {
	return gt.repoURL
}

// Version yields the ref identifying the desired version of a file,
// e.g. "master" in https://gitea.com/fredbi/go-vcsfetch/src/branch/master/README.md
func (gt *URL) Version() string {
	return gt.version
}

// Path yields the file path relative to the repository,
// e.g. "README.md" in https://gitea.com/fredbi/go-vcsfetch/src/branch/master/README.md
func (gt *URL) Path() string {
	return gt.path
}
