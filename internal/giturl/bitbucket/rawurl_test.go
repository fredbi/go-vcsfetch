// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package bitbucket

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
		const emptyPath = "https://bitbucket.org/workspace/repo/"

		u, err := url.Parse(emptyPath)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyPath, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid bitbucket URL string, but got: %q: %v",
			emptyPath, err,
		)

		_, err = Raw(raw)
		require.Errorf(t, err, "expected an empty path to return an error")
	})

	t.Run("should convert URL with empty version to raw", func(t *testing.T) {
		const emptyVersion = "https://bitbucket.org/workspace/repo/src/main/file"

		u, err := url.Parse(emptyVersion)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			emptyVersion, err,
		)
		raw, err := Parse(u)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid bitbucket URL string, but got: %q: %v",
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
			"test is wrongly configured: expected a valid bitbucket locator string, but got: %q: %v",
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
			"test is wrongly configured: expected a valid bitbucket locator string, but got: %q: %v",
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
				url:     "https://bitbucket.org/workspace/repo/src/master/README.md",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://bitbucket.org/workspace/repo/src/main/pkg/doc.go",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "main",
				path:    "pkg/doc.go",
			},
			{
				url:     "https://bitbucket.org/workspace/repo/raw/master/README.md",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "master",
				path:    "README.md",
			},
			{
				url:     "https://bitbucket.org/atlassian/python-bitbucket/src/main/pybitbucket/auth.py",
				repo:    "https://bitbucket.org/atlassian/python-bitbucket",
				version: "main",
				path:    "pybitbucket/auth.py",
			},
			{
				url:     "https://bitbucket.org/workspace/repo/src/v1.0.0/LICENSE",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "v1.0.0",
				path:    "LICENSE",
			},
			{
				url:     "https://bitbucket.org/workspace/repo/src/abc123def456/file.txt",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "abc123def456",
				path:    "file.txt",
			},
			{
				url:     "https://bitbucket.org/workspace/project/src/develop/internal/util.go",
				repo:    "https://bitbucket.org/workspace/project",
				version: "develop",
				path:    "internal/util.go",
			},
			{
				url:     "https://bitbucket.example.com/workspace/repo/src/release/docs/api.md",
				repo:    "https://bitbucket.example.com/workspace/repo",
				version: "release",
				path:    "docs/api.md",
			},
			{
				url:     "https://bitbucket.org:443/workspace/repo/src/main/config.yaml",
				repo:    "https://bitbucket.org:443/workspace/repo",
				version: "main",
				path:    "config.yaml",
			},
			{
				url:     "https://bitbucket.org/workspace/repo.git/src/main/file.go",
				repo:    "https://bitbucket.org/workspace/repo",
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
				url:     "https://bitbucket.org/workspace/repo",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "",
				path:    "/",
			},
			{
				url:     "https://bitbucket.org/workspace/repo/src/main",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "main",
				path:    "/",
			},
			{
				url:     "ssh://git@bitbucket.org/workspace/repo/src/main/file.go",
				repo:    "ssh://git@bitbucket.org/workspace/repo",
				version: "main",
				path:    "file.go",
			},
			{
				url:     "https://bitbucket.org:8080/workspace/repo/src/main/file.go",
				repo:    "https://bitbucket.org:8080/workspace/repo",
				version: "main",
				path:    "file.go",
			},
			{
				url:     "https://bitbucket.org/workspace/repo.git",
				repo:    "https://bitbucket.org/workspace/repo",
				version: "",
				path:    "/",
			},
		},
	)
}
