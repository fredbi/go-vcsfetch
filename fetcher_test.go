package vcsfetch

import (
	"bytes"
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetcher(t *testing.T) {
	t.Run("with defaults", func(t *testing.T) {
		fetcher := NewFetcher()
		ctx := context.Background()
		w := new(bytes.Buffer)

		t.Run("with invalid URLs", func(t *testing.T) {
			t.Run("should NOT fetch from invalid URL string", func(t *testing.T) {
				const invalidLocator = "" // invalid URL
				err := fetcher.Fetch(ctx, w, invalidLocator)
				require.ErrorIs(t, err, Error)
			})
			t.Run("should NOT fetch from invalid locator URL", func(t *testing.T) {
				invalidLocator := &url.URL{} // invalid URL
				err := fetcher.FetchURL(ctx, w, invalidLocator)
				require.ErrorIs(t, err, Error)
			})
			t.Run("should NOT fetch from invalid Locator", func(t *testing.T) {
				invalidLocator := &MockLocator{} // invalid Locator
				err := fetcher.FetchLocator(ctx, w, invalidLocator)
				require.ErrorIs(t, err, Error)
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
	})

	t.Run("with ssh authentication", func(t *testing.T) {
	})

	t.Run("with recurse submodules", func(t *testing.T) {
	})

	t.Run("with slug shorthand", func(t *testing.T) {
	})
}
