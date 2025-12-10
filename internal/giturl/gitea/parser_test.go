// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package gitea

import (
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
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
			name:        "gitea.com repo only",
			input:       "https://gitea.com/owner/repo",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "",
			wantPath:    "/",
			wantErr:     false,
		},
		{
			name:        "gitea.com src with branch and file",
			input:       "https://gitea.com/owner/repo/src/branch/master/README.md",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "master",
			wantPath:    "README.md",
			wantErr:     false,
		},
		{
			name:        "gitea.com raw with branch and file",
			input:       "https://gitea.com/owner/repo/raw/branch/main/path/to/file.go",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "main",
			wantPath:    "path/to/file.go",
			wantErr:     false,
		},
		{
			name:        "gitea.com with tag",
			input:       "https://gitea.com/owner/repo/src/tag/v1.0.0/LICENSE",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "v1.0.0",
			wantPath:    "LICENSE",
			wantErr:     false,
		},
		{
			name:        "gitea.com with commit",
			input:       "https://gitea.com/owner/repo/src/commit/abc123/file.txt",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "abc123",
			wantPath:    "file.txt",
			wantErr:     false,
		},
		{
			name:        "custom gitea instance",
			input:       "https://git.example.com/owner/repo/src/branch/develop/code.js",
			wantRepo:    "https://git.example.com/owner/repo",
			wantVersion: "develop",
			wantPath:    "code.js",
			wantErr:     false,
		},
		{
			name:        "repo with .git suffix",
			input:       "https://gitea.com/owner/repo.git/src/branch/main/file",
			wantRepo:    "https://gitea.com/owner/repo",
			wantVersion: "main",
			wantPath:    "file",
			wantErr:     false,
		},
		{
			name:    "invalid - missing ref type",
			input:   "https://gitea.com/owner/repo/src/master/file",
			wantErr: true,
		},
		{
			name:    "invalid - missing owner/repo",
			input:   "https://gitea.com/owner",
			wantErr: true,
		},
		{
			name:    "invalid - wrong discriminator",
			input:   "https://gitea.com/owner/repo/blob/branch/main/file",
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
