package github

import (
	"fmt"
	"net/url"
	"strings"
)

type GitHubURL struct {
	repoURL *url.URL
	path    string
	version string
}

const (
	defaultScheme = "https"
	defaultHost   = "github.com"
)

func Parse(githubURL *url.URL) (*GitHubURL, error) {
	u := &url.URL{}
	*u = *githubURL // shallow clone

	const (
		repoIndex = 2
		kindIndex = 3
		refIndex  = 4
	)

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

	// pth, _ := strings.CutPrefix(u.Path, "/")
	pth := strings.ToLower(u.Path)
	parts := strings.Split(pth, "/")

	if len(parts) < refIndex {
		return nil, fmt.Errorf("expected the URL path component to contain at least %d parts, but got %q", refIndex, pth)
	}

	repo := strings.Join(parts[:repoIndex], "/")
	repo = strings.TrimSuffix(repo, ".git")
	u.Path = repo

	var (
		// isFile, isDir bool
		ref   string
		index int
	)
	isRaw := strings.HasPrefix(u.Host, "raw")

	if isRaw {
		if parts[kindIndex] != "refs" {
			return nil, fmt.Errorf(`expected raw content URL path to contain "refs" but got %q in %q`, parts[kindIndex], pth)
		}
		ref = parts[refIndex+1]
		index = refIndex + 2
	} else {
		switch parts[kindIndex] {
		case "blob", "tree":
		default:
			return nil, fmt.Errorf(`expected URL path to contain "blob" or "tree" but got %q in %q`, parts[kindIndex], pth)
		}
		ref = parts[refIndex]
		index = refIndex + 1
	}

	if len(parts) < index {
		return nil, fmt.Errorf(`expected URL path to contain at least %d parts, but got %q`, index, pth)
	}

	path := strings.Join(parts[index:], "/")
	u.RawFragment = ""
	u.Fragment = ""
	u.RawQuery = ""

	gh := &GitHubURL{
		repoURL: u,
		path:    path,
		version: ref,
	}
	return gh, nil
}

// RepoURL yields the base URL of the vcs repository,
// e.g. https://github.com/fredbi/go-vcsfetcher
func (gh *GitHubURL) RepoURL() *url.URL {
	return gh.repoURL
}

// Version yields the ref identifying the desired version of a file,
// e.g. "master" in https://github.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *GitHubURL) Version() string {
	return gh.version
}

// Path yields the file path relative to the repository,
// e.g. "README.md" in https://github.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *GitHubURL) Path() string {
	return gh.path
}
