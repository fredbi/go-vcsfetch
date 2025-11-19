package giturl

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type testURL struct {
	u                *url.URL
	expectedProvider Provider
	expectError      bool
	expectedError    error
}

func TestAutoDetect(t *testing.T) {
	t.Parallel()

	for tc := range testURLs(t) {
		t.Run(fmt.Sprintf("with %v", tc.u), testSingleURL(tc))
	}
}

func testSingleURL(tc testURL) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		if tc.expectError {
			t.Run(fmt.Sprintf("should not auto-detect provider for %v", tc.u), func(t *testing.T) {
				provider, locator, err := AutoDetect(tc.u)

				require.Error(t, err)
				require.ErrorIs(t, err, ErrProvider)
				require.Nil(t, locator)
				require.Equal(t, tc.expectedProvider, provider)
				if tc.expectedError != nil {
					require.ErrorContains(t, err, tc.expectedError.Error())
				}
			})

			return
		}

		t.Run(fmt.Sprintf("should auto-detect provider %v for %v", tc.expectedProvider, tc.u), func(t *testing.T) {
			provider, locator, err := AutoDetect(tc.u)

			require.NoError(t, err)
			require.Equal(t, tc.expectedProvider, provider)
			require.NotNil(t, locator)
			t.Logf("%v", locator)
		})
	}
}

func testURLs(t *testing.T) iter.Seq[testURL] {
	t.Helper()

	return slices.Values(
		[]testURL{
			{
				u:                mustParseURL(t, "https://github.big-corporation.com/big-repo/blob/tree/master/README.md"),
				expectedProvider: ProviderGithub,
			},
			{
				u:                mustParseURL(t, "https://chez.com/big-repo/blob/tree/master/README.md"),
				expectedProvider: ProviderUnknown,
				expectError:      true,
				expectedError:    ErrUnknownProvider,
			},
		},
	)
}

func mustParseURL(t *testing.T, str string) *url.URL {
	t.Helper()

	u, err := url.Parse(str)
	require.NoError(t, err)

	return u
}
