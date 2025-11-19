package git

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	u, err := url.Parse("https://github.com/go-swagger/go-swagger")
	require.NoError(t, err)

	r := NewRepo(u, &Options{GitSkipAutoDetect: true})
	require.NotNil(t, r)

	var w bytes.Buffer
	ctx := t.Context()

	require.NoError(t,
		r.Fetch(ctx, &w, ".github/CONTRIBUTING.md", "v0"),
	)
	t.Logf("%v", w.String())

	w.Reset()
	require.NoError(t,
		r.Fetch(ctx, &w, "notes/v0.33.0.md", "v0.33.0"),
	)

	t.Logf("%v", w.String())
}
