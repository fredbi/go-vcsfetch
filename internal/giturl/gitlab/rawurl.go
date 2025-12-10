package gitlab

import (
	"fmt"
	"net/url"
	"path"
)

type Locator interface {
	RepoURL() *url.URL
	Path() string
	Version() string
}

// Raw returns the raw URL for a [Locator] hosted on any gitlab SCM instance.
//
// Example:
//
//   - https://gitlab.com/fredbi/go-vcsfetch/-/raw/release/README.md
func Raw(locator Locator) (*url.URL, error) {
	pth := locator.Path()
	if pth == "" {
		return nil, fmt.Errorf("returning a raw content url requires a non empty path to a file: %w", ErrGitlab)
	}

	version := locator.Version()
	if version == "" {
		version = "HEAD"
	}

	u := locator.RepoURL()
	u.Path = path.Join(u.Path, "-", "raw", version, locator.Path())
	u.Fragment = ""
	u.RawFragment = ""

	return u, nil
}
