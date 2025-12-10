package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

// downloadError is a sentinel error type to report errors from this package.
type downloadError string

func (e downloadError) Error() string {
	return string(e)
}

// ErrDownload is a sentinel error to report errors from the download content package.
const ErrDownload downloadError = "error downloading file"

// Supported indicates if the provided URL can be downloaded.
//
// This works for http and https URL schemes, but not ssh or git.
func Supported(u *url.URL) bool {
	scheme, _ := strings.CutPrefix(u.Scheme, "git+")

	switch scheme {
	case schemeHTTP, schemeHTTPS:
		return true
	default:
		return false
	}
}

// Content downloads a file from a remote URL and copies the fetched content to an [io.Writer].
//
// [Content] currently supports only the http and https URL schemes (no support for local files).
func Content(ctx context.Context, u *url.URL, w io.Writer, opts *Options) error {
	scheme, _ := strings.CutPrefix(u.Scheme, "git+")
	v := *u
	v.Scheme, _ = strings.CutPrefix(u.Scheme, "git+")

	switch scheme {
	case schemeHTTP, schemeHTTPS:
		return httpContent(ctx, &v, w, opts)
	default:
		return fmt.Errorf("unsupported URL scheme: %s: %w", u.Scheme, ErrDownload)
	}
}

func httpContent(ctx context.Context, u *url.URL, w io.Writer, opts *Options) error {
	if opts == nil {
		opts = &defaultOptions
	}

	var client *http.Client
	if opts.Client != nil {
		client = opts.Client
	} else {
		client = http.DefaultClient
	}

	if opts.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
		ctx = timeoutCtx
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return errors.Join(err, ErrDownload)
	}

	for key, val := range opts.CustomHeaders {
		req.Header.Set(key, val)
	}

	if opts.BasicAuthUsername != "" && opts.BasicAuthPassword != "" {
		req.SetBasicAuth(opts.BasicAuthUsername, opts.BasicAuthPassword)
	}

	resp, err := client.Do(req)
	defer func() {
		if resp == nil || resp.Body == nil {
			return
		}

		_ = resp.Body.Close()
	}()

	if err != nil {
		return errors.Join(err, ErrDownload)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not fetch resource at %q [%s]: %w", u.String(), resp.Status, ErrDownload)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return errors.Join(err, ErrDownload)
	}

	return nil
}
