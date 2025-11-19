package github

import (
	"fmt"
	"net/url"
	"path"
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
// For Github Enterprise, there is no way to guess the host.
//
// Examples:
//
//   - https://raw.githubusercontent.com/fredbi/go-vcsfetch/refs/heads/master/README.md
//   - https://raw.githubusercontent.com/fredbi/go-vcsfetch/master/README.md
func Raw(locator Locator) (*url.URL, error) {
	// TODO:
	// - check if already raw.githubusercontent.com host
	// - check that the scheme is http or https
	repo := locator.RepoURL()
	pth := locator.Path()
	if pth == "" {
		return nil, fmt.Errorf("returning a raw content url requires a non empty path to a file: %w", ErrGithub)
	}

	version := locator.Version()
	if version == "" {
		version = "HEAD"
	}

	host := repo.Host
	if host == "github.com" {
		u := repo
		u.Host = "raw.githubusercontent.com"
		u.Path = path.Join(u.Path, version, locator.Path())
		u.Fragment = ""
		u.RawFragment = ""

		return u, nil
	}

	return nil, fmt.Errorf("no way to guess the raw content host for github not hosted by github.com: %q: %w", host, ErrGithub)
}
