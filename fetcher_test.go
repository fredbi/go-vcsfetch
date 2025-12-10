package vcsfetch

import (
	"bytes"
	"context"
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestFetcher(t *testing.T) {
	t.Parallel()

	t.Run("with defaults", func(t *testing.T) {
		fetcher := NewFetcher()
		ctx := context.Background()
		w := new(bytes.Buffer)

		t.Run("with invalid URLs", func(t *testing.T) {
			t.Run("should NOT fetch from invalid URL string", func(t *testing.T) {
				const invalidLocator = "" // invalid URL
				err := fetcher.Fetch(ctx, w, invalidLocator)
				require.ErrorIs(t, err, ErrVCS)
			})
			t.Run("should NOT fetch from invalid locator URL", func(t *testing.T) {
				invalidLocator := &url.URL{} // invalid URL
				err := fetcher.FetchURL(ctx, w, invalidLocator)
				require.ErrorIs(t, err, ErrVCS)
			})
			t.Run("should NOT fetch from invalid Locator", func(t *testing.T) {
				invalidLocator := invalidLocator(t)
				err := fetcher.FetchLocator(ctx, w, invalidLocator)
				require.ErrorIs(t, err, ErrVCS)
			})
		})

		t.Run("with valid URLs", func(t *testing.T) {
			t.Run("should fetch HEAD from master", func(t *testing.T) {
				t.SkipNow()
			})
			t.Run("should fetch HEAD from branch", func(t *testing.T) {
				t.SkipNow()
			})
			t.Run("should fetch tag", func(t *testing.T) {
				t.SkipNow()
			})
			t.Run("should fetch latest compatible semver tag (minor release)", func(t *testing.T) {
				t.SkipNow()
			})
			t.Run("should fetch latest compatible semver tag (major release)", func(t *testing.T) {
				t.SkipNow()
			})
		})
	})

	t.Run("with options", func(t *testing.T) {
		t.Run("with version required", func(t *testing.T) {
			t.Run("should NOT fetch HEAD from branch by default", func(t *testing.T) {
				t.SkipNow()
			})
		})

		t.Run("with exact tag", func(t *testing.T) {
			t.Run("should NOT fetch latest compatible semver tag (minor release)", func(t *testing.T) {
				t.SkipNow()
			})
			t.Run("should fetch exact tag", func(t *testing.T) {
				t.SkipNow()
			})
		})

		t.Run("with pre-released allowed", func(t *testing.T) {
			t.Run("should fetch latest pre-release semver tag (major release)", func(t *testing.T) {
				t.SkipNow()
			})
		})

		t.Run("with https authentication", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("with ssh authentication", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("with recurse submodules", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("with slug shorthand", func(t *testing.T) {
			t.SkipNow()
		})

		t.Run("with backing directory", func(t *testing.T) {
			t.SkipNow()
		})
	})
}

func invalidLocator(t *testing.T) *MockLocator {
	return &MockLocator{
		RepoURLFunc: func() *url.URL {
			return &url.URL{}
		},
		PathFunc: func() string {
			return ""
		},
		VersionFunc: func() string {
			return ""
		},
	}
}
