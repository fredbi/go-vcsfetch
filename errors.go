// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

type vcsFetchError string

func (e vcsFetchError) Error() string {
	return string(e)
}

// ErrVCS is a sentinel error for all errors that originate from this package.
const ErrVCS vcsFetchError = "vcsfetch error"
