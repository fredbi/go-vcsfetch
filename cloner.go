// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"

	"github.com/fredbi/go-vcsfetch/internal/git"
)

// Cloner allows for working with vcs repositories to perform cloning or sparse cloning.
//
// The [Cloner] is intended for read-only capture of remote resources. If you need to mutate
// the cloned resources, please consider using another tool.
//
// See [Fetcher] for available options.
//
// # Fetching multiple resources
//
// The [Cloner] may be used to fetch against the cloned resources using a similar syntax as
// with a [Fetcher], using the [Fetch.FetchFromClone] methods. All fetched locators must then match with the cloned base URL or will
// return an error.
//
// # Concurrency
//
// The [Cloner] is not intended for concurrent usage: it is a stateful object.
// Once a repository has been cloned, it becomes accessible via [Cloner.FS].
//
// You may use [Cloner.Close] to relinquish memory or temporary disk resources and reuse the [Cloner].
//
// Exception: when using [WithBackingDir] with a non-empty directory, the cloned content
// it not removed after usage and left up to the caller to leave it or clean it if needed.
type Cloner struct {
	cloneOptions

	clonedURL *url.URL
	clonedFS  fs.FS
}

// NewCloner builds a [Cloner] to retrieve an entire vcs repository.
func NewCloner(opts ...CloneOption) *Cloner {
	return &Cloner{
		cloneOptions: optionsWithDefaults(opts),
	}
}

// Clone a vcs repository.
//
// The repoURL string must be a valid URL.
//
// The URL is detected to be either a valid SPDX locator or a well-known giturl.
//
// The clone is accessible as a read-only [fs.FS] using [Cloner.FS].
func (f *Cloner) Clone(ctx context.Context, repoURL string) error {
	u, err := url.Parse(repoURL)
	if err != nil {
		return fmt.Errorf("expected a valid URL: %w: %w", err, ErrVCS)
	}

	return f.CloneURL(ctx, u)
}

// CloneLocator clones a vcs repository from a [Locator].
//
// The clone is accessible as a read-only [fs.FS] using [Cloner.FS].
func (f *Cloner) CloneLocator(ctx context.Context, locator Locator, opts ...CloneOption) error {
	repo := git.NewRepo(locator.RepoURL(), f.toInternalGitOptions())

	fs, err := repo.Clone(ctx, locator.Version(), f.toInternalGitCloneOptions())
	if err != nil {
		return err
	}

	f.clonedURL = locator.RepoURL()
	f.clonedFS = fs

	return nil
}

// CloneURL clones a vcs repository from a [url.URL].
//
// The clone is accessible as a read-only [fs.FS] using [Cloner.FS].
func (f *Cloner) CloneURL(ctx context.Context, u *url.URL) error {
	var locator Locator
	spdxLocator, err := SPDXLocatorFromURL(u, f.spdxOpts...)
	if err == nil {
		locator = spdxLocator
	} else {
		gitLocator, err := GitLocatorFromURL(u, f.gitLocOpts...)
		if err != nil {
			return fmt.Errorf("the provided URL is not a SPDX locator or a recognized git URL: %w: %w", err, ErrVCS)
		}
		locator = gitLocator
	}

	return f.CloneLocator(ctx, locator)
}

func (f *Cloner) FS() fs.FS {
	return f.clonedFS
}

// FetchFromClone fetches a single file from the cloned repository.
func (f *Cloner) FetchFromClone(ctx context.Context, w io.Writer, location string) error {
	u, err := url.Parse(location)
	if err != nil {
		return fmt.Errorf("expected a valid URL: %w: %w", err, ErrVCS)
	}

	return f.FetchURLFromClone(ctx, w, u)
}

// FetchLocatorFromClone fetches a single file from the cloned repository, using a [Locator].
func (f *Cloner) FetchLocatorFromClone(ctx context.Context, w io.Writer, locator Locator) error {
	if f.clonedURL == nil || f.clonedFS == nil {
		return fmt.Errorf("cannot fetch from clone: no clone available yet: %w", ErrVCS)
	}

	if locator.RepoURL().String() != f.clonedURL.String() {
		return fmt.Errorf("cannot fetch from clone not matching the cloned repo URL: %w", ErrVCS)
	}

	file, err := f.clonedFS.Open(locator.Path())
	if err != nil {
		return fmt.Errorf("cannot fetch from clone: %w: %w", err, ErrVCS)
	}

	_, err = io.Copy(w, file)

	return err
}

// FetchURLFromClone fetches a single file from the cloned repository, using a [url.URL].
func (f *Cloner) FetchURLFromClone(ctx context.Context, w io.Writer, u *url.URL) error {
	var locator Locator
	spdxLocator, err := SPDXLocatorFromURL(u, f.spdxOpts...)
	if err == nil {
		locator = spdxLocator
	} else {
		gitLocator, err := GitLocatorFromURL(u, f.gitLocOpts...)
		if err != nil {
			return fmt.Errorf("the provided URL is not a SPDX locator or a recognized git URL: %w: %w", err, ErrVCS)
		}
		locator = gitLocator
	}

	return f.FetchLocatorFromClone(ctx, w, locator)
}

// Close resets the state of the cloner.
func (f *Cloner) Close() error {
	if f.clonedFS == nil {
		return nil
	}

	return nil // TODO
}
