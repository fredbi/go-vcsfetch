package giturl

type providerError string

func (e providerError) Error() string {
	return string(e)
}

const (
	// ErrProvider is a sentinel error for all errors that originate from this package.
	ErrProvider providerError = "git-url provider detection error"

	// ErrUnknownProvider is raised whenever a URL cannot be associated with a well-known SCM provider.
	ErrUnknownProvider providerError = "unrecognized git-url provider in URL"

	// ErrNotImplementedProvider is currently raised for the azure provider.
	ErrNotImplementedProvider providerError = "provider is detected but not implemented yet"
)
