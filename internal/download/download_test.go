package download

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestContent(t *testing.T) {
	t.Parallel()

	const (
		remote   = "https://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
		expected = "END OF TERMS AND CONDITIONS"
	)
	remoteURL := mustURL(t, remote)

	t.Run("with default options", func(t *testing.T) {
		var b bytes.Buffer
		ctx := t.Context()

		require.NoError(t, Content(ctx, remoteURL, &b, nil))
		require.Contains(t, b.String(), expected)
	})

	t.Run("with options", func(t *testing.T) {
		t.Run("with nil Client", func(t *testing.T) {
			var b bytes.Buffer
			ctx := t.Context()

			require.NoError(t, Content(ctx, remoteURL, &b, &Options{}))
			require.Contains(t, b.String(), expected)
		})
	})

	t.Run("with unsupported scheme", func(t *testing.T) {
		const remoteSSH = "ssh://git@raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"

		var b bytes.Buffer
		invalidURL := mustURL(t, remoteSSH)
		ctx := t.Context()
		require.Error(t, Content(ctx, invalidURL, &b, nil))
	})
}

func TestSupported(t *testing.T) {
	t.Parallel()

	t.Run("http[s] should be supported", func(t *testing.T) {
		const (
			remoteInsecure    = "http://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
			remoteSecure      = "https://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
			remoteInsecureGit = "git+http://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
		)

		require.True(t, Supported(mustURL(t, remoteInsecure)))
		require.True(t, Supported(mustURL(t, remoteSecure)))
		require.True(t, Supported(mustURL(t, remoteInsecureGit)))
	})

	t.Run("Other URL schemes should NOT be supported", func(t *testing.T) {
		const (
			remoteSSH         = "ssh://git@raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
			remoteTCP         = "git://git@raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
			localFile         = "file://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
			remoteInsecureGit = "http+git://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
		)

		require.False(t, Supported(mustURL(t, remoteSSH)))
		require.False(t, Supported(mustURL(t, remoteTCP)))
		require.False(t, Supported(mustURL(t, localFile)))
		require.False(t, Supported(mustURL(t, remoteInsecureGit)))
	})
}

func mustURL(t *testing.T, str string) *url.URL {
	t.Helper()

	u, err := url.Parse(str)
	require.NoError(t, err)

	return u
}
