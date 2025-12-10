// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fredbi/go-vcsfetch/internal/giturl"
)

var _ Locator = &GitLocator{}

// GitLocator describes an URL used to access a vcs resource over git
// using common URL formats (github, gitlab, ...).
//
// The URL may use schemes git, http, https or ssh.
//
// See https://git-scm.com/docs/git-fetch#_git_urls for reference.
type GitLocator struct {
	repo *url.URL
	url.Userinfo

	Provider  string
	Transport string
	Host      string
	RepoPath  string
	Ref       string
	SubPath   string
}

// ParseGitLocator builds a [GitLocator] from an URL string.
func ParseGitLocator(location string, opts ...GitLocatorOption) (*GitLocator, error) {
	if location == "" {
		return nil, fmt.Errorf("empty locator is invalid: %w", ErrVCS)
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("a git locator should be a valid URL: %w: %w", err, ErrVCS)
	}

	return GitLocatorFromURL(u, opts...)
}

// GitLocatorFromURL builds a [GitLocator] from an [url.URL].
func GitLocatorFromURL(u *url.URL, opts ...GitLocatorOption) (*GitLocator, error) {
	ref := ""
	o := optionsWithDefaults(opts)
	if o.requireVersion && ref == "" {
		return nil, fmt.Errorf("a non-empty version is required: %w", ErrVCS)
	}

	provider, loc, err := giturl.AutoDetect(u)
	if err != nil {
		return nil, fmt.Errorf("invalid git locator: %w: %w", err, ErrVCS)
	}

	var userinfo url.Userinfo
	if u.User != nil {
		userinfo = *(u.User)
	}

	gl := &GitLocator{
		repo:      loc.RepoURL(),
		Provider:  string(provider),
		Userinfo:  userinfo,
		Transport: u.Scheme, // TODO: factorize with spdx
		Host:      u.Host,
		Ref:       loc.Version(),
		SubPath:   loc.Path(),
	}

	return gl, nil // TODO
}

func (l *GitLocator) RepoURL() *url.URL {
	return l.repo
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
