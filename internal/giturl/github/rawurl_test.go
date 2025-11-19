package github

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaw(t *testing.T) {
	t.Parallel()

	t.Run("with convertible URLs", func(t *testing.T) {
		for tc := range rawTestCasesValid(t) {
			t.Run("should convert to raw", testShouldRaw(tc))
		}
	})

	t.Run("with non-convertible URLs", func(t *testing.T) {
		for tc := range rawTestCasesInvalid(t) {
			t.Run("should NOT convert to raw", testShouldNotRaw(tc))
		}
	})
}

func testShouldRaw(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		u, err := url.Parse(tc.url)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			tc.url, err,
		)

		res, err := Raw(u)
		require.NoErrorf(t, err, "unexpected error: %v for %v", err, u)
	}
}

func testShouldNotRaw(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		u, err := url.Parse(tc.url)
		require.NoErrorf(t, err,
			"test is wrongly configured: expected a valid URL string, but got: %q: %v",
			tc.url, err,
		)

		res, err := Raw(u)
		require.Errorf(t, err, "expected error for %v", u)
	}
}
