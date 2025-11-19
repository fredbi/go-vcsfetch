// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/fredbi/go-vcsfetch/internal/download"
	"github.com/fredbi/go-vcsfetch/internal/git"
	"github.com/fredbi/go-vcsfetch/internal/giturl"
)

// Fetcher allows for working with vcs repositories to perform cloning, sparse cloning
// and single file fetching.
//
// The [Fetcher] is intended for read-only capture of remote resources. If you need to mutate
// the cloned resources, please consider using another tool.
//
// # Concurrency
//
// The [Fetcher] is stateles and may be called concurrently.
//
// All fetches are carried out independently. If you plan to fetch multiple resources against a single
// repository, consider using a [Cloner] for improved performances.
type Fetcher struct {
	fetchOptions
}

// NewFetcher builds a [Fetcher] to retrieve single files from a vcs repository.
func NewFetcher(opts ...FetchOption) *Fetcher {
	return &Fetcher{
		fetchOptions: optionsWithDefaults(opts),
	}
}

// Fetch a single file from a vcs location string.
//
// The content of the fetched file is copied to the passed [io.Writer].
//
// The string argument must be a valid URL.
func (f *Fetcher) Fetch(ctx context.Context, w io.Writer, location string) error {
	u, err := url.Parse(location)
	if err != nil {
		return fmt.Errorf("expected a valid URL: %w: %w", err, ErrVCS)
	}

	return f.FetchURL(ctx, w, u)
}

// FetchLocator fetches a single file specified by a [Locator] from a vcs location.
//
// The content of the fetched file is copied to the passed [io.Writer].
//
// If you want to retrieve a locator representing a folder, use [Cloner.CloneLocator] with sparse option.
//
// NOTE: this package provides 2 implementations of the [Locator].
// You may pass your own implementation of this interface to this method.
func (f *Fetcher) FetchLocator(ctx context.Context, w io.Writer, locator Locator) error {
	if f.requireVersion && locator.Version() == "" {
		return fmt.Errorf("an explicit version is required, but %v does not specify a version: %w", locator, ErrVCS)
	}

	// short-circuit that avoids the use of git thanks to a direct raw-content download URL from the SCM.
	//
	// This works fine on github.com and all gitlab instances.
	if download.Supported(locator.RepoURL()) {
		rawURL, err := giturl.Raw(locator)
		if err == nil {
			if e := download.Content(ctx, rawURL, w, f.toInternalDownloadOptions()); e != nil {
				return fmt.Errorf("could not fetch raw content from %q: %w: %w", rawURL, e, ErrVCS)
			}
		}
	}

	// general-purpose git retrieval
	repo := git.NewRepo(locator.RepoURL(), f.toInternalGitOptions())
	if err := repo.Fetch(ctx, w, locator.Path(), locator.Version()); err != nil {
		return errors.Join(err, ErrVCS)
	}

	return nil
}

// FetchURL fetches a single file from a vcs location as an URL.
//
// The content of the fetched file is copied to the passed [io.Writer].
//
// If the URL is detected to be a valid SPDX locator, it is equivalent to [Fetcher.FetchLocator] with a [SPDXLocator].
// Otherwise, it falls back to git-url parsing and is equivalent to [Fetcher.FetchLocator] with a [GitLocator].
//
// If you want to retrieve an URL representing a folder, use [Cloner.CloneURL] with sparse option instead.
func (f *Fetcher) FetchURL(ctx context.Context, w io.Writer, u *url.URL) error {
	var locator Locator
	spdxLocator, err := SPDXLocatorFromURL(u, f.spdxOpts...)
	if err == nil {
		// prioritize spdx locator
		locator = spdxLocator
	} else {
		// fallback on a giturl
		gitLocator, err := GitLocatorFromURL(u, f.gitLocOpts...)
		if err != nil {
			return fmt.Errorf("the provided URL is not a SPDX locator or a recognized git URL: %w: %w", err, ErrVCS)
		}
		locator = gitLocator
	}

	if err := f.FetchLocator(ctx, w, locator); err != nil {
		return errors.Join(err, ErrVCS)
	}

	return nil
}
