package github

type githubError string

func (e githubError) Error() string {
	return string(e)
}

// ErrGithub is a sentinel error for all errors that originate from this package.
const ErrGithub githubError = "github provider error"
