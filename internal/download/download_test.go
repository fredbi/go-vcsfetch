package download

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContent(t *testing.T) {
	var b bytes.Buffer
	const remote = "https://raw.githubusercontent.com/fredbi/go-vcsfetch/HEAD/LICENSE"
	u, err := url.Parse(remote)
	require.NoError(t, err)
	ctx := t.Context()

	require.NoError(t, Content(ctx, u, &b, nil))

	t.Log(b.String())
}
