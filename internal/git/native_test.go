package git

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNativeGithubRepository(t *testing.T) {
	u, err := url.Parse("ssh://git@github.com/go-swagger/go-swagger")
	require.NoError(t, err)

	r := NewRepo(u, &Options{
		GitSkipAutoDetect: false,
		Debug:             true,
	})
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

func TestNativeGitlabRepository(t *testing.T) {
	u, err := url.Parse("https://gitlab.com/gitlab-org/gitlab-runner")
	require.NoError(t, err)

	r := NewRepo(u, &Options{
		GitSkipAutoDetect: false,
		Debug:             true,
	})
	require.NotNil(t, r)

	var w bytes.Buffer
	ctx := t.Context()

	require.NoError(t,
		r.Fetch(ctx, &w, "LICENSE", "main"),
	)
	t.Logf("%v", w.String())
}
