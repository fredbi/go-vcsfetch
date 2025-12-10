package git

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestGitlabRepository(t *testing.T) {
	u, err := url.Parse("https://gitlab.com/gitlab-org/gitlab-runner")
	require.NoError(t, err)

	r := NewRepo(u, &Options{GitSkipAutoDetect: true})
	require.NotNil(t, r)

	var w bytes.Buffer
	ctx := t.Context()

	require.NoError(t,
		r.Fetch(ctx, &w, "LICENSE", "main"),
	)
	t.Logf("%v", w.String())
}
