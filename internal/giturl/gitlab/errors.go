package gitlab

type gitlabError string

func (e gitlabError) Error() string {
	return string(e)
}

// ErrGitlab is a sentinel error for all errors that originate from this package.
const ErrGitlab gitlabError = "gitlab provider error"
