package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Supported indicates if the provided URL can be downloaded.
//
// This works for http and https URL schemes, but not ssh or git.
func Supported(u *url.URL) bool {
	scheme, _ := strings.CutPrefix(u.Scheme, "git+")

	switch scheme {
	case "http", "https":
		return true
	default:
		return false
	}
}

func Content(ctx context.Context, u *url.URL, w io.Writer, opts *Options) error {
	scheme, _ := strings.CutPrefix(u.Scheme, "git+")
	v := *u
	v.Scheme, _ = strings.CutPrefix(u.Scheme, "git+")

	if scheme == "http" || scheme == "https" {
		return httpContent(ctx, &v, w, opts)
	}

	return fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
}

func httpContent(ctx context.Context, u *url.URL, w io.Writer, opts *Options) error {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	// TODO: other options for auth, headers etc
	resp, err := client.Do(req)
	defer func() {
		if resp == nil || resp.Body == nil {
			return
		}

		_ = resp.Body.Close()
	}()

	if err != nil {
		return err // TODO: wrap errors
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not fetch resource at %q [%s]", u.String(), resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err // TODO wrap error
	}

	return nil
}
