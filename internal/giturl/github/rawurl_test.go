package github

import (
	"iter"
	"net/url"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaw(t *testing.T) {
	t.Parallel()

	t.Run("with valid raw URLs", func(t *testing.T) {
		for tc := range rawTestCasesValid(t) {
			t.Run("should convert to raw", testShouldRaw(tc))
		}
	})

	t.Run("with non-raw URLs", func(t *testing.T) {
		for tc := range rawTestCasesInvalid(t) {
			t.Run("should NOT convert to raw", testShouldNotRaw(tc))
		}
	})
}

func TestRawEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("should NOT convert URL with empty file path to raw", func(t *testing.T) {
		const emptyPath = "https://github.com/owner/repo/"

		u, err := url.Parse(emptyPath)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyPath, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid github URL string, but got: %q: %v",
			emptyPath, err,
		)

		_, err = Raw(raw)
		require.Errorf(t, err, "expected an empty path to return an error")
	})

	t.Run("should convert URL with empty version to raw", func(t *testing.T) {
		const emptyVersion = "https://github.com/owner/repo/tree/v2.1/file"

		u, err := url.Parse(emptyVersion)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyVersion, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid github URL string, but got: %q: %v",
			emptyVersion, err,
		)
		raw.version = "" // force empty version

		v, err := Raw(raw)
		require.NoErrorf(t, err, "expected an empty version to be supported")
		require.Contains(t, v.String(), "HEAD")
	})
}

func testShouldRaw(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		u, err := url.Parse(tc.url)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			tc.url, err,
		)

		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid github locator string, but got: %q: %v",
			tc.url, err,
		)

		res, err := Raw(raw)
		require.NoErrorf(t, err, "unexpected error: %v for %v", err, u)
		require.NotEmpty(t, res.String())
	}
}

func testShouldNotRaw(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		u, err := url.Parse(tc.url)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			tc.url, err,
		)

		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid github locator string, but got: %q: %v",
			tc.url, err,
		)

		res, err := Raw(raw)
		require.Errorf(t, err, "expected error for %v", u)
		require.Nil(t, res)
	}
}

func rawTestCasesValid(_ *testing.T) iter.Seq[testCase] {
	return slices.Values(
		[]testCase{
			{
				url:     "https://github.com/fredbi/go-vcsfetch/blob/master/README.md",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/blob/HEAD/pkg/doc.go",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "HEAD",
				path:    "pkg/doc.go",
			},
			{
				url:     "https://raw.githubusercontent.com/fredbi/go-vcsfetch/refs/heads/master/README.md",
				repo:    "https://raw.githubusercontent.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
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
				url:     "https://github.com:443/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://github.com:443/fredbi/go-vcsfetch/tree/v1.1.0/pkg/utils.go",
				repo:    "https://github.com:443/fredbi/go-vcsfetch",
				version: "",
				path:    "/",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/v2.1/README.md",
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
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			// TODO: escaped paths
		},
	)
}

func rawTestCasesInvalid(_ *testing.T) iter.Seq[testCase] {
	return slices.Values(
		[]testCase{
			{
				url:     "https://corporate.github.com/fredbi/go-vcsfetch/tree/v2.1/LICENSE",
				repo:    "https://corporate.github.com/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "LICENSE",
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
				url:     "ssh://:443/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "ssh://github.com:443/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
			},
			{
				url:     "https://github.com/fredbi/go-vcsfetch/tree/master/",
				repo:    "https://github.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "/",
			},
			{
				url:     "https://github.com:445/fredbi/go-vcsfetch/tree/v2.1/pkg/doc",
				repo:    "https://github.com:445/fredbi/go-vcsfetch",
				version: "v2.1",
				path:    "pkg/doc",
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
		},
	)
}
