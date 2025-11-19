package git

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
)

func getRemoteCapabilities(ctx context.Context, o *gogit.FetchOptions) (*capability.List, error) {
	s, err := newUploadPackSession(o.RemoteURL, o.Auth, o.InsecureSkipTLS, o.ClientCert, o.ClientKey, o.CABundle, o.ProxyOptions)
	if err != nil {
		return nil, err
	}

	ar, err := s.AdvertisedReferencesContext(ctx)
	if err != nil {
		return nil, err
	}

	return ar.Capabilities, nil
}

func newUploadPackSession(url string, auth transport.AuthMethod, insecure bool, clientCert, clientKey, caBundle []byte, proxyOpts transport.ProxyOptions) (transport.UploadPackSession, error) {
	c, ep, err := newClient(url, insecure, clientCert, clientKey, caBundle, proxyOpts)
	if err != nil {
		return nil, err
	}

	return c.NewUploadPackSession(ep, auth)
}

func newClient(url string, insecure bool, clientCert, clientKey, caBundle []byte, proxyOpts transport.ProxyOptions) (transport.Transport, *transport.Endpoint, error) {
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, nil, err
	}
	ep.InsecureSkipTLS = insecure
	ep.ClientCert = clientCert
	ep.ClientKey = clientKey
	ep.CaBundle = caBundle
	ep.Proxy = proxyOpts

	c, err := client.NewClient(ep)
	if err != nil {
		return nil, nil, err
	}

	return c, ep, err
}
