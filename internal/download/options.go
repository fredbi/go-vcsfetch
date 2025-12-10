package download

import (
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// Options sets HTTP request options.
type Options struct {
	Timeout           time.Duration
	BasicAuthUsername string
	BasicAuthPassword string
	CustomHeaders     map[string]string
	Client            *http.Client
}

var defaultOptions = Options{
	Timeout: defaultTimeout,
	Client:  http.DefaultClient,
}
