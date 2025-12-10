// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package bitbucket

type bitbucketError string

func (e bitbucketError) Error() string {
	return string(e)
}

// ErrBitbucket is a sentinel error for all errors that originate from this package.
const ErrBitbucket bitbucketError = "bitbucket error"
