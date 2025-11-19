package github

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGithubURLParser(t *testing.T) {
	t.Parallel()

	t.Run("with valid github urls", func(t *testing.T) {
		for tc := range parserTestCasesValid(t) {
			t.Run(fmt.Sprintf("should parse %v", tc.url), testShouldParseURL(tc))
		}
	})

	t.Run("with invalid github urls", func(t *testing.T) {
		for tc := range parserTestCasesInvalid(t) {
			t.Run(fmt.Sprintf("should NOT parse %v", tc.url), testShouldNotParseURL(tc))
		}
	})
}

func testShouldParseURL(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
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
}

func testShouldNotParseURL(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		u, err := url.Parse(tc.url)
		require.NoErrorf(t, err, "test is wrongly configured: expected a valid URL string, but got: %q: %v", tc.url, err)

		_, err = Parse(u)
		require.Errorf(t, err, "expected url %q to parse with an error", tc.url)
	}
}

type testCase struct {
	url     string
	repo    string
	version string
	path    string
}

func parserTestCasesValid(_ *testing.T) iter.Seq[testCase] {
	return slices.Values(
		[]testCase{
			{
				url:     "https://github.com/fredbi/go-vcsfetch/blob/master/README.md",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/master/",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "/",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/blob/HEAD/pkg/doc.go",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "HEAD",
				path:    "pkg/doc.go",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://raw.githubusercontent.com/fredbi/go-vcsfetch/refs/heads/master/README.md",
				repo:    "https://raw.githubusercontent.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "ssh://git@github.com/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "ssh://git@github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "ssh://git@github.com/fredbi/go-vcsfetch.git/tree/v2.1/pkg/doc",
				repo:    "ssh://git@github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://raw.githubusercontent.com/fredbi/go-vcsfetch/refs/remotes/release/README.md",
				repo:    "https://raw.githubusercontent.com/fredbi/go-vcsfetch",
				version: "release",
				path:    "README.md",
			},
			{
				url:     "fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "ssh://:443/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "ssh://github.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://github.com:443/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://github.com:443/fredbi/go-vcsfetch",
				repo:    "https://github.com:443/fredbi/go-vcsfetch",
				version: "",
				path:    "/",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/v2.1",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "/",
			},
			{
				url:     "https://raw.githubusercontent.com/fredbi/go-vcsfetch/v2.1/LICENSE",
				repo:    "https://raw.githubusercontent.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "LICENSE",
			},
			// TODO: escaped paths
		},
	)
}

func parserTestCasesInvalid(_ *testing.T) iter.Seq[testCase] {
	return slices.Values(
		[]testCase{
			{
				url: "https://github.com/fredbi/go-vcsfetch/blob",
			},
			{
				url: "https://github.com/fredbi/go-vcsfetch/refs/HEAD/pkg/doc.go",
			},
			{
				url: "https://raw.githubusercontent.com/fredbi/go-vcsfetch/blob/heads/master/README.md",
			},
			{
				url: "https://github.com/fredbi/",
			},
			{
				url: "https://raw.githubusercontent.com/fredbi/go-vcsfetch",
			},
			{
				url: "https://raw.githubusercontent.com/fredbi/go-vcsfetch/blob/README.md",
			},
			{
				url: "https://github.com/fredbi/go-vcsfetch/blob/master/",
			},
		},
	)
}
