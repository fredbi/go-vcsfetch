package git

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
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
