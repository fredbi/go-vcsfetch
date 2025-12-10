// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package gitea

type giteaError string

func (e giteaError) Error() string {
	return string(e)
}

// ErrGitea is a sentinel error for all errors that originate from this package.
const ErrGitea giteaError = "gitea error"
