// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

type errVcsFetch string

func (e errVcsFetch) Error() string {
	return string(e)
}

// Error is a sentinel error for all errors that originate from this package.
const Error errVcsFetch = "vcsfetch error"
