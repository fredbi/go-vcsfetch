package git

/*
const (
	shaPattern = `^[a-f0-9]{40}$|^[a-f0-9]{7}$`
)

var shaRegex = regexp.MustCompile(shaPattern)
*/

/*
	// we don't really need this
	var commitSha, tag, branch string
	const (
		gitTagsPrefix    = "refs/tags"
		gitHeadsPrefix   = "refs/heads"
		gitRemotesPrefix = "refs/remotes"
	)

	if ref != "" {
		if shaRegex.MatchString(ref) {
			commitSha = ref
		} else if trimmed, isTag := strings.CutPrefix(ref, gitTagsPrefix); isTag {
			tag = trimmed
		} else if trimmed, isHead := strings.CutPrefix(ref, gitHeadsPrefix); isHead {
			branch = trimmed
		} else if trimmed, isRemote := strings.CutPrefix(ref, gitRemotesPrefix); isRemote {
			branch = trimmed
		} else {
			// don't know for sure
		}
	}
*/
