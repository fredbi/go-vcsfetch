package github

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
// This package exposes [URL] as an implementation for github.
type Locator interface {
	RepoURL() *url.URL
	Path() string
	Version() string
}

// Raw returns the raw.githubusercontent URL for a [Locator] hosted on github.com.
//
// Only https url's are supported.
//
// For Github Enterprise, there is no way to guess the host: this only works on github.com
//
// Examples:
//
//   - https://raw.githubusercontent.com/fredbi/go-vcsfetch/refs/heads/master/README.md
//   - https://raw.githubusercontent.com/fredbi/go-vcsfetch/master/README.md
func Raw(locator Locator) (*url.URL, error) {
	repo := locator.RepoURL()
	pth := strings.Trim(locator.Path(), "/")
	if pth == "" {
		return nil, fmt.Errorf("returning a raw content url requires a non empty path to a file: %w", ErrGithub)
	}

	version := locator.Version()
	if version == "" {
		version = "HEAD"
	}

	scheme, _ := strings.CutSuffix(repo.Scheme, "+git")

	if scheme != "https" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL scheme: %w", ErrGithub)
	}

	if port := repo.Port(); port != "" && port != "443" {
		return nil, fmt.Errorf("returning a raw content url requires a https URL with standard port (443 or unspecified): %w", ErrGithub)
	}

	host := repo.Hostname()
	if host == defaultHost || host == rawHost {
		u := repo
		u.Host = "raw.githubusercontent.com"
		u.Path = path.Join(u.Path, version, pth)
		u.Fragment = ""
		u.RawFragment = ""

		return u, nil
	}

	return nil, fmt.Errorf("no way to guess the raw content host for github not hosted by github.com: %q: %w", host, ErrGithub)
}
