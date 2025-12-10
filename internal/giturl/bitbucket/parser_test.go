// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package bitbucket

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	url     string
	repo    string
	version string
	path    string
}

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		wantRepo    string
		wantVersion string
		wantPath    string
		wantErr     bool
	}{
		{
			name:        "bitbucket.org repo only",
			input:       "https://bitbucket.org/workspace/repo",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "",
			wantPath:    "/",
			wantErr:     false,
		},
		{
			name:        "bitbucket.org src with ref and file",
			input:       "https://bitbucket.org/workspace/repo/src/master/README.md",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "master",
			wantPath:    "README.md",
			wantErr:     false,
		},
		{
			name:        "bitbucket.org raw with ref and file",
			input:       "https://bitbucket.org/workspace/repo/raw/main/path/to/file.go",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "main",
			wantPath:    "path/to/file.go",
			wantErr:     false,
		},
		{
			name:        "bitbucket.org with tag ref",
			input:       "https://bitbucket.org/workspace/repo/src/v1.0.0/LICENSE",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "v1.0.0",
			wantPath:    "LICENSE",
			wantErr:     false,
		},
		{
			name:        "bitbucket.org with commit sha",
			input:       "https://bitbucket.org/workspace/repo/src/abc123def456/file.txt",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "abc123def456",
			wantPath:    "file.txt",
			wantErr:     false,
		},
		{
			name:        "bitbucket.org with nested path",
			input:       "https://bitbucket.org/atlassian/python-bitbucket/src/main/pybitbucket/auth.py",
			wantRepo:    "https://bitbucket.org/atlassian/python-bitbucket",
			wantVersion: "main",
			wantPath:    "pybitbucket/auth.py",
			wantErr:     false,
		},
		{
			name:        "repo with .git suffix",
			input:       "https://bitbucket.org/workspace/repo.git/src/main/file",
			wantRepo:    "https://bitbucket.org/workspace/repo",
			wantVersion: "main",
			wantPath:    "file",
			wantErr:     false,
		},
		{
			name:        "bitbucket server (self-hosted)",
			input:       "https://bitbucket.example.com/workspace/project/src/develop/code.js",
			wantRepo:    "https://bitbucket.example.com/workspace/project",
			wantVersion: "develop",
			wantPath:    "code.js",
			wantErr:     false,
		},
		{
			name:    "invalid - missing workspace/repo",
			input:   "https://bitbucket.org/workspace",
			wantErr: true,
		},
		{
			name:    "invalid - wrong discriminator",
			input:   "https://bitbucket.org/workspace/repo/blob/main/file",
			wantErr: true,
		},
		{
			name:    "invalid - missing ref",
			input:   "https://bitbucket.org/workspace/repo/src",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u, err := url.Parse(tc.input)
			require.NoError(t, err)

			got, err := Parse(u)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tc.wantRepo, got.RepoURL().String())
			require.Equal(t, tc.wantVersion, got.Version())
			require.Equal(t, tc.wantPath, got.Path())
		})
	}
}
