package gitlab

import (
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestGitlabURLParser(t *testing.T) {
	t.Run("valid gitlab urls", func(t *testing.T) {
		for _, tc := range []testCase{
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/blob/master/README.md",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/tree/master/",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "/",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/blob/HEAD/pkg/doc.go",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "HEAD",
				path:    "pkg/doc.go",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/tree/v2.1/pkg/doc",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "ssh://git@gitlab.com/fredbi/go-vcsfetch/-/tree/v2.1/pkg/doc",
				repo:    "ssh://git@gitlab.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "ssh://git@gitlab.com/fredbi/go-vcsfetch.git/-/tree/v2.1/pkg/doc",
				repo:    "ssh://git@gitlab.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/raw/release/README.md",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "release",
				path:    "README.md",
			},
			{
				url:     "fredbi/go-vcsfetch/-/tree/v2.1/pkg/doc",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "ssh://:443/fredbi/go-vcsfetch/-/tree/v2.1/pkg/doc",
				repo:    "ssh://gitlab.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://gitlab.com:443/fredbi/go-vcsfetch/-/tree/v2.1/pkg/doc",
				repo:    "https://gitlab.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://gitlab.com:443/fredbi/go-vcsfetch",
				repo:    "https://gitlab.com:443/fredbi/go-vcsfetch",
				version: "",
				path:    "/",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "",
				path:    "/",
			},
			{
				url:     "https://gitlab.com/fredbi/go-vcsfetch/-/",
				repo:    "https://gitlab.com/fredbi/go-vcsfetch",
				version: "",
				path:    "/",
			},
			// TODO: escaped paths
		} {
			u, err := url.Parse(tc.url)
			require.NoErrorf(t, err,
				"test is wrongly configured: expected a valid URL string, but got: %q: %v",
				tc.url, err,
			)

			res, err := Parse(u)
			require.NoErrorf(t, err, "unexpected error: %v for %v", err, u)

			require.Equal(t, tc.repo, res.RepoURL().String())
			require.Equal(t, tc.version, res.Version())
			require.Equal(t, tc.path, res.Path())
		}
	})

	t.Run("invalid gitlab urls", func(t *testing.T) {
		for _, tc := range []testCase{
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/blob",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/blob/main",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/refs/HEAD/pkg/doc.go",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/tree/v2.1",
			},
			{
				url: "https://raw.gitlabusercontent.com/fredbi/go-vcsfetch/blob/heads/master/README.md",
			},
			{
				url: "https://raw.gitlabusercontent.com/fredbi/go-vcsfetch/-/refs/heads/master/README.md",
			},
			{
				url: "https://gitlab.com/fredbi/",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/blob/README.md",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/blob/master/",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/-/refs/heads/master/README.md",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/-/blob",
			},
			{
				url: "https://gitlab.com/fredbi/go-vcsfetch/-/blob/master",
			},
		} {
			u, err := url.Parse(tc.url)
			require.NoErrorf(t, err, "test is wrongly configured: expected a valid URL string, but got: %q: %v", tc.url, err)

			_, err = Parse(u)
			require.Errorf(t, err, "expected url %q to parse with an error", tc.url)
		}
	})
}

type testCase struct {
	url     string
	repo    string
	version string
	path    string
}
