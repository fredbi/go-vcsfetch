package gitlab

import (
	"fmt"
	"net/url"
	"strings"
)

// URL is a gitlab-style URL to a vcs resource hosted by gitlab SCM.
type URL struct {
	repoURL *url.URL
	path    string
	version string
}

const (
	defaultScheme = "https"
	defaultHost   = "gitlab.com"
)

// Parse a gitlab URL.
func Parse(gitlabURL *url.URL) (*URL, error) {
	u := &url.URL{}
	*u = *gitlabURL // shallow clone

	if u.Scheme == "" {
		u.Scheme = defaultScheme
	} else {
		u.Scheme, _ = strings.CutPrefix(u.Scheme, "git+")
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
		refIndex  = 4
	)

	parts := strings.Split(pth, "/")
	if len(parts) < repoIndex {
		return nil, fmt.Errorf("expected the URL path component to contain at least %d parts, but got %q: %w", refIndex, pth, ErrGitlab)
	}

	repo := strings.Join(parts[:repoIndex], "/")
	repo = strings.TrimSuffix(repo, ".git")
	u.Path = repo

	if len(parts) == repoIndex || len(parts) == repoIndex+1 && parts[repoIndex] == "-" {
		// entire repo
		u.RawFragment = ""
		u.Fragment = ""
		u.RawQuery = ""

		gh := &URL{
			repoURL: u,
			path:    "/",
			version: "",
		}

		return gh, nil
	}

	parts = parts[repoIndex:]
	if parts[0] != "-" {
		return nil, fmt.Errorf(`expected URL path to contain a "-" separator: %w`, ErrGitlab)
	}

	parts = parts[1:]

	var (
		ref    string
		isTree bool
	)

	const neededPartsAfterDash = 2
	if len(parts) < neededPartsAfterDash {
		return nil, fmt.Errorf(`expected URL path to contain at least 2 parts but got %q: %w`, pth, ErrGitlab)
	}

	switch strings.ToLower(parts[0]) {
	case "blob", "raw":
	case "tree":
		isTree = true
	default:
		return nil, fmt.Errorf(`expected URL path to contain "blob" or "tree" but got %q in %q: %w`, parts[0], pth, ErrGitlab)
	}

	ref = parts[1]
	parts = parts[2:]

	if len(parts) == 0 {
		if !isTree {
			return nil, fmt.Errorf(`expected URL path to contain at least %d parts in %q: %w`, refIndex, pth, ErrGitlab)
		}

		parts = []string{"/"}
	}

	repoPath := strings.Join(parts, "/")
	u.RawFragment = ""
	u.Fragment = ""
	u.RawQuery = ""

	gh := &URL{
		repoURL: u,
		path:    repoPath,
		version: ref,
	}

	return gh, nil
}

// RepoURL yields the base URL of the vcs repository,
// e.g. https://gitlab.com/fredbi/go-vcsfetcher
func (gh *URL) RepoURL() *url.URL {
	return gh.repoURL
}

// Version yields the ref identifying the desired version of a file,
// e.g. "master" in https://gitlab.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *URL) Version() string {
	return gh.version
}

// Path yields the file path relative to the repository,
// e.g. "README.md" in https://gitlab.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *URL) Path() string {
	return gh.path
}
