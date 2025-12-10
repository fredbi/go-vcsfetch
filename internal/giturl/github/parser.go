package github

import (
	"fmt"
	"net/url"
	"strings"
)

// URL is a github-style URL to a vcs resource hosted by github SCM.
type URL struct {
	repoURL *url.URL
	path    string
	version string
}

const (
	defaultScheme = "https"
	defaultHost   = "github.com"
	rawHost       = "raw.githubusercontent.com"
)

// Parse a github URL.
func Parse(githubURL *url.URL) (*URL, error) {
	u := &url.URL{}
	*u = *githubURL // shallow clone

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
	isRaw := strings.HasPrefix(strings.ToLower(u.Host), "raw")
	pth := strings.Trim(u.Path, "/")

	const (
		repoIndex = 2
		refIndex  = 4
	)

	parts := strings.Split(pth, "/")
	if len(parts) < repoIndex {
		return nil, fmt.Errorf("expected the URL path component to contain at least %d parts, but got %q: %w", repoIndex, pth, ErrGithub)
	}

	repo := strings.Join(parts[:repoIndex], "/")
	repo = strings.TrimSuffix(repo, ".git")
	u.Path = repo

	if len(parts) == repoIndex {
		if isRaw {
			return nil, fmt.Errorf(`expected raw content URL path to contain "refs" but got empty path: %w`, ErrGithub)
		}

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

	var (
		ref    string
		isTree bool
	)

	if isRaw {
		discriminator := strings.ToLower(parts[0])
		switch discriminator {
		case "refs":
			const neededPartsForRaw = 3
			if len(parts) < neededPartsForRaw {
				return nil, fmt.Errorf(`expected raw content URL path to contain at least %d parts but got %q: %w`, neededPartsForRaw, pth, ErrGithub)
			}
			// skip parts[1] (e.g. "heads")
			ref = parts[2]
			parts = parts[3:]

		case "blob", "tree": // not sure how github behaves if there is a branch or a tag called "blob" or "tree"...
			return nil, fmt.Errorf(`expected raw content URL path to contain "refs" but got %q in %q: %w`, parts[0], pth, ErrGithub)
		default:
			// parts[0] is the ref
			const neededPartsForRaw = 2
			if len(parts) < neededPartsForRaw {
				return nil, fmt.Errorf(`expected raw content URL path to contain at least %d parts but got %q: %w`, neededPartsForRaw, pth, ErrGithub)
			}
			ref = parts[0]
			parts = parts[1:]
		}
	} else {
		const neededPartsForBlob = 2
		if len(parts) < neededPartsForBlob {
			return nil, fmt.Errorf(`expected URL path to contain at least %d parts but got %q: %w`, neededPartsForBlob, pth, ErrGithub)
		}

		switch strings.ToLower(parts[0]) {
		case "blob":
		case "tree":
			isTree = true
		default:
			return nil, fmt.Errorf(`expected URL path to contain "blob" or "tree" but got %q in %q: %w`, parts[0], pth, ErrGithub)
		}
		ref = parts[1]
		parts = parts[2:]
	}

	if len(parts) == 0 {
		if !isTree {
			return nil, fmt.Errorf(`expected URL path to contain at least %d parts in %q: %w`, refIndex, pth, ErrGithub)
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
// e.g. https://github.com/fredbi/go-vcsfetcher
func (gh *URL) RepoURL() *url.URL {
	return gh.repoURL
}

// Version yields the ref identifying the desired version of a file,
// e.g. "master" in https://github.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *URL) Version() string {
	return gh.version
}

// Path yields the file path relative to the repository,
// e.g. "README.md" in https://github.com/fredbi/go-vcsfetcher/blob/master/README.md
func (gh *URL) Path() string {
	return gh.path
}
