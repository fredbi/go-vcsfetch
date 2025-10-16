// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"fmt"
	"net/url"
	"strings"
)

var _ Locator = &SPDXLocator{}

// SPDXLocator describes a SPDX VCS locator, with all its components detailed.
//
// It implements the [Locator] interface.
//
// The SPDX (Software Package Data Exchange) specification provides a detailed framework for referencing software
// components, including through Version Control System (VCS) locators.
//
// Normative references
//
//   - https://spdx.github.io/spdx-spec/v2.3/package-information/#77-package-download-location-field
//   - https://spdx.org/rdf/spdx-terms-v2.0/dataproperties/downloadLocation___-576346518.html>
//
// TL;DR: the SPDX locator comes with a "@" in the URL path for the version and a "#" URL fragment for the
// target file (or directory).
//
// # SPDX VCS Locator format
//
// The VCS location syntax, as described in the latest SPDX version,
// resembles a URL with specific structure to accommodate different version control systems.
//
// # VCS Location structure
//
// Format:
//
//	<vcs_tool>+<transport>://<host_name>[/<path_to_repository>][@<revision_tag_or_branch>][#<sub_path>]
//
// Where:
//
//	<vcs_tool>: Specifies the type of version control system (e.g., git, hg, svn, bzr).
//	<transport>: Indicates the transport mechanism (e.g., ssh, https).
//	<host_name>: The server or host where the repository resides.
//	<path_to_repository>: The path to the repository if applicable.
//	<revision_tag_or_branch>: Identifies a specific commit, branch, or tag in the repository.
//	<sub_path>: Optional, specifies a sub-directory or file path within the repository.
//
// # Examples
//
//   - git:
//     git+https://github.com/user/repo.git@main#file
//
//   - mercurial (not supported by [Fetcher] and [Cloner] yet):
//     hg+https://www.mercurial-scm.org/repo/myrepo@branchname#file
//
//   - subversion (won't be supported):
//     svn+https://svn.example.com/repo/trunk@1234#subdir
//
// # Implementation tolerances and limitations
//
// Our use-case for SPDX locators is limited to single file retrieval:
//   - an URL fragment is required
//
// Our implementation supports a full URL with the following:
//
//   - an empty "vcs-tool" part is tolerated in the scheme and defaults to "git".
//     Therefore schemes such as "git+https" and "https" are equivalent.
//   - "username:password" credentials
//   - hostname port
//   - query parameters in URL are ignored but tolerated
//   - the absence of an explicit reference provided with "@" will be resolved as the head of the default branch
//
// Optionally, the [SPDXLocator] may support SCM-specific shorthands using "git repo slugs":
//
//   - github, gitlab:
//     fredbi/go-vcsfetcher@HEAD#.github/dependabot.yaml (implied: "https://github.com" or "https://gitlab.com")
//
// The implied vcs base URL is customizable with [WithRootURL].
type SPDXLocator struct {
	url.Userinfo

	Tool      string
	Transport string
	Host      string
	RepoPath  string
	Ref       string
	SubPath   string
}

// ParseSPDXLocator parses a VCS locator string and returns its components as a [SPDXLocator].
func ParseSPDXLocator(location string, opts ...SPDXOption) (*SPDXLocator, error) {
	if location == "" {
		return nil, fmt.Errorf("empty locator is invalid: %w", Error)
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("a SPDX locator should be a valid URL: %w: %w", err, Error)
	}

	return SPDXLocatorFromURL(u, opts...)
}

// SPDXLocatorFromURL parses an URL into a [SPDXLocator].
func SPDXLocatorFromURL(u *url.URL, opts ...SPDXOption) (*SPDXLocator, error) {
	o := optionsWithDefaults(opts)

	if u.Path == "" {
		return nil, fmt.Errorf("SPDX locator requires an URL path: %w", Error)
	}
	if u.Fragment == "" {
		return nil, fmt.Errorf("SPDX locator requires an URL fragment to specify a single file: %w", Error)
	}

	// scheme analyis
	var tool, transport string
	parts := strings.SplitN(u.Scheme, "+", 2)
	if len(parts) > 0 {
		tool = parts[0]
		transport = parts[1]
	} else {
		tool = "git"
		transport = u.Scheme
	}

	var repoPath, ref string
	parts = strings.SplitN(u.Path, "@", 2)
	if len(parts) > 0 {
		repoPath = parts[0]
		ref = parts[1]
		// 	} else {
		repoPath = u.Path
	}
	if o.requireVersion && ref == "" {
		return nil, fmt.Errorf("a non-empty version is required: %w", Error)
	}

	var userinfo url.Userinfo
	if u.User != nil {
		userinfo = *(u.User)
	}

	return &SPDXLocator{
		Userinfo:  userinfo,
		Tool:      tool,
		Transport: transport,
		Host:      u.Host,
		RepoPath:  repoPath,
		Ref:       ref,
		SubPath:   u.Fragment,
	}, nil
}

func (l *SPDXLocator) RepoURL() *url.URL {
	u := &url.URL{
		Scheme: l.Transport,
		User:   &l.Userinfo,
		Host:   l.Host,
		Path:   l.RepoPath,
	}

	return u
}

func (l *SPDXLocator) Version() string {
	return l.Ref
}

func (l *SPDXLocator) Path() string {
	return l.SubPath
}

func (l *SPDXLocator) IsLocal() bool {
	return l.Transport == "file"
}

func (l *SPDXLocator) HasAuth() bool {
	_, isSet := l.Password()
	return isSet
}

func (l *SPDXLocator) String() string {
	u := l.RepoURL()
	if l.Tool != "" {
		u.Scheme = l.Tool + "+" + u.Scheme
	}
	u.Path += "@" + l.Version()
	u.Fragment = l.Path()

	return u.String()
}
