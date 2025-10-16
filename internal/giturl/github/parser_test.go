package github

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGithubURLParser(t *testing.T) {
	t.Run("valid github urls", func(t *testing.T) {
		for _, tc := range []testCase{
			{
				url:     "https://github.com/fredbi/go-vcsfetch/blob/master/README.md",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
		} {
			u, err := url.Parse(tc.url)
			require.NoErrorf(t, err, "test is wrongly configured: expected a valid URL string, but got: %q: %v", tc.url, err)

			res, err := Parse(u)
			require.NoError(t, err)

			require.Equal(t, tc.repo, res.RepoURL().String())
			require.Equal(t, tc.version, res.Version())
			require.Equal(t, tc.path, res.Path())
		}
	})

	t.Run("valid github raw content urls", func(t *testing.T) {
	})

	t.Run("invalid github urls", func(t *testing.T) {
	})

	t.Run("invalid github raw content urls", func(t *testing.T) {
	})
}

type testCase struct {
	url     string
	repo    string
	version string
	path    string
}
