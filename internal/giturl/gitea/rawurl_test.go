// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package gitea

import (
	"iter"
	"net/url"
	"slices"
	"testing"

	"github.com/go-openapi/testify/v2/require"
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
		const emptyPath = "https://gitea.com/owner/repo/"

		u, err := url.Parse(emptyPath)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyPath, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid gitea URL string, but got: %q: %v",
			emptyPath, err,
		)

		_, err = Raw(raw)
		require.Errorf(t, err, "expected an empty path to return an error")
	})

	t.Run("should convert URL with empty version to raw", func(t *testing.T) {
		const emptyVersion = "https://gitea.com/owner/repo/src/branch/main/file"

		u, err := url.Parse(emptyVersion)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyVersion, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid gitea URL string, but got: %q: %v",
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
			"test is wrongly configured: expected a valid gitea locator string, but got: %q: %v",
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
			"test is wrongly configured: expected a valid gitea locator string, but got: %q: %v",
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
				url:     "https://gitea.com/fredbi/go-vcsfetch/src/branch/master/README.md",
				repo:    "https://gitea.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://gitea.com/fredbi/go-vcsfetch/src/branch/HEAD/pkg/doc.go",
				repo:    "https://gitea.com/fredbi/go-vcsfetch",
				version: "HEAD",
				path:    "pkg/doc.go",
			},
			{
				url:     "https://gitea.com/fredbi/go-vcsfetch/raw/branch/master/README.md",
				repo:    "https://gitea.com/fredbi/go-vcsfetch",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://gitea.com/fredbi/go-vcsfetch/src/tag/v1.0.0/LICENSE",
				repo:    "https://gitea.com/fredbi/go-vcsfetch",
				version: "v1.0.0",
				path:    "LICENSE",
			},
			{
				url:     "https://gitea.com/fredbi/go-vcsfetch/src/commit/abc123def/file.txt",
				repo:    "https://gitea.com/fredbi/go-vcsfetch",
				version: "abc123def",
				path:    "file.txt",
			},
			{
				url:     "https://gitea.com/owner/repo/src/branch/develop/internal/util.go",
				repo:    "https://gitea.com/owner/repo",
				version: "develop",
				path:    "internal/util.go",
			},
			{
				url:     "https://try.gitea.io/owner/project/src/branch/main/docs/api.md",
				repo:    "https://try.gitea.io/owner/project",
				version: "main",
				path:    "docs/api.md",
			},
			{
				url:     "https://gitea.com/owner/repo.git/src/branch/main/file.go",
				repo:    "https://gitea.com/owner/repo",
				version: "main",
				path:    "file.go",
			},
		},
	)
}

func rawTestCasesInvalid(_ *testing.T) iter.Seq[testCase] {
	return slices.Values(
		[]testCase{
			{
				url:     "https://gitea.com/owner/repo",
				repo:    "https://gitea.com/owner/repo",
				version: "",
				path:    "/",
			},
			{
				url:     "https://gitea.com/owner/repo/src/branch/main",
				repo:    "https://gitea.com/owner/repo",
				version: "main",
				path:    "/",
			},
			{
				url:     "ssh://git@gitea.com/owner/repo/src/branch/main/file.go",
				repo:    "ssh://git@gitea.com/owner/repo",
				version: "main",
				path:    "file.go",
			},
			{
				url:     "https://gitea.com:8080/owner/repo/src/branch/main/file.go",
				repo:    "https://gitea.com:8080/owner/repo",
				version: "main",
				path:    "file.go",
			},
			{
				url:     "https://git.example.com:8443/org/repo/src/branch/release/v2/config.yaml",
				repo:    "https://git.example.com:8443/org/repo",
				version: "release",
				path:    "v2/config.yaml",
			},
		},
	)
}
