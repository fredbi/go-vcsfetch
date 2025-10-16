// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"fmt"
	"net/url"
	"strings"
)

var _ Locator = &GitLocator{}

// GitLocator describes an URL used to access a vcs resource over git
// using common URL formats (github, gitlab, ...).
//
// See also https://git-scm.com/docs/git-fetch#_git_urls
type GitLocator struct {
	url.Userinfo

	Transport string
	Host      string
	RepoPath  string
	Ref       string
	SubPath   string
}

func ParseGitLocator(location string, opts ...GitLocatorOption) (*GitLocator, error) {
	if location == "" {
		return nil, fmt.Errorf("empty locator is invalid: %w", Error)
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("a git locator should be a valid URL: %w: %w", err, Error)
	}

	return GitLocatorFromURL(u, opts...)
}

func GitLocatorFromURL(u *url.URL, opts ...GitLocatorOption) (*GitLocator, error) {
	ref := ""
	o := optionsWithDefaults(opts)
	if o.requireVersion && ref == "" {
		return nil, fmt.Errorf("a non-empty version is required: %w", Error)
	}
	return nil, nil // TODO
}

func (l *GitLocator) RepoURL() *url.URL {
	u := &url.URL{
		Scheme: l.Transport,
		User:   &l.Userinfo,
		Host:   l.Host,
		Path:   l.RepoPath,
	}

	return u
}

func (l *GitLocator) Version() string {
	return l.Ref
}

func (l *GitLocator) Path() string {
	return l.SubPath
}

func (l *GitLocator) IsLocal() bool {
	return l.Transport == "file"
}

func (l *GitLocator) HasAuth() bool {
	_, isSet := l.Password()
	return isSet
}

func (l *GitLocator) String() string {
	u := l.RepoURL()
	if !strings.HasPrefix(u.Scheme, "git+") {
		u.Scheme = "git+" + u.Scheme
	}
	u.Path += "@" + l.Version()
	u.Fragment = l.Path()

	return u.String()
}
