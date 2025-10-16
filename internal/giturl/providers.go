package giturl

import (
	"net/url"
)

type Provider string

const (
	ProviderGitHub    Provider = "github"
	ProviderGitlab    Provider = "gitlab"
	ProviderAzure     Provider = "azure"
	ProviderBitBucket Provider = "bitbucket"
	ProviderGitea     Provider = "gitea"
)

func (p Provider) String() string {
	return string(p)
}

type Parser interface {
	Parse(*url.URL) (any, error)
}
