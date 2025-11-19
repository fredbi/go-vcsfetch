package giturl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fredbi/go-vcsfetch/internal/giturl/github"
	"github.com/fredbi/go-vcsfetch/internal/giturl/gitlab"
)

// Provider represents a SCM platform with a proprietary git-url format.
type Provider string

const (
	ProviderUnknown   Provider = "unknown"
	ProviderGithub    Provider = "github"
	ProviderGitlab    Provider = "gitlab"
	ProviderAzure     Provider = "azure"
	ProviderBitBucket Provider = "bitbucket"
	ProviderGitea     Provider = "gitea"
)

func (p Provider) String() string {
	return string(p)
}

// Locator is the minimal interface returned by a parsed URL.
type Locator interface {
	RepoURL() *url.URL
	Path() string
	Version() string
}

// AutoDetect tries to determine the [Provider] that corresponds to a given [url.URL].
//
// Detection is rather crude and based on the host in the URL.
//
// It may not work for SCMs deployed on-premises.
func AutoDetect(u *url.URL) (Provider, Locator, error) {
	host := strings.ToLower(u.Host)

	switch {
	case strings.Contains(host, ProviderGithub.String()):
		locator, err := github.Parse(u)

		return ProviderGithub, locator, err
	case strings.Contains(host, ProviderGitlab.String()):
		locator, err := github.Parse(u)
		return ProviderGitlab, locator, err
	case strings.Contains(host, ProviderAzure.String()):
		panic("not implemented") // TODO
	case strings.Contains(host, ProviderBitBucket.String()):
		panic("not implemented") // TODO
	case strings.Contains(host, ProviderGitea.String()):
		panic("not implemented") // TODO
	default:
		return ProviderUnknown, nil, fmt.Errorf("url=%q: %w: %w", u.String(), ErrUnknownProvider, ErrProvider)
	}
}

// Raw transforms a [Locator] into a raw-content URL to retrieve a vcs resource from well-known SCM providers.
//
// This allows to bypass the use of git and is usually faster.
func Raw(locator Locator) (*url.URL, error) {
	provider, _, err := AutoDetect(locator.RepoURL())
	if err != nil {
		return nil, err
	}

	switch provider {
	case ProviderGithub:
		return github.Raw(locator)
	case ProviderGitlab:
		return gitlab.Raw(locator)
	default:
		panic("not implemented") // TODO
	}
}
