// SPDX-FileCopyrightText: Copyright 2025 Frédéric BIDON
// SPDX-License-Identifier: Apache-2.0

package vcsfetch

import (
	"fmt"
	"net/url"
	"os"

	"github.com/fredbi/go-vcsfetch/internal/download"
	"github.com/fredbi/go-vcsfetch/internal/git"
)

func optionsWithDefaults[O any, T ~func(*O)](opts []T) O {
	var o O
	var ptr any = &o
	if defaulter, ok := ptr.(interface{ applyDefaults() O }); ok {
		defaulter.applyDefaults()
	}

	for _, apply := range opts {
		apply(&o)
	}

	return o
}

// FetchOption configures a [Fetcher] with optional behavior.
type FetchOption func(*fetchOptions)

// FetchWithBackingDir tells the [Fetcher] to back the fetched resources
// on disk. By default, fetched resources are mapped in memory.
//
// If dir is empty, the default is given by [os.MkDirTemp] using "vcsclone" as the pattern.
// In this case, [FetchWithBackingDir] panics if it can't create a temporary directory.
//
// When using [FetchWithBackingDir] with a non-empty directory, the fetched content
// will not be removed after usage and left up to the caller to leave it or clean it if needed.
func FetchWithBackingDir(enabled bool, dir string) FetchOption {
	return func(o *fetchOptions) {
		withGitBackingDir(enabled, dir)(&o.gitOptions)
	}
}

// FetchWithGitSkipAutoDetect skips the auto-detection of a local git binary.
//
// Whenever enabled, git binary autodetection allows for some operations to be performed
// faster using git native implementation rather than the pure go implementation.
func FetchWithGitSkipAutoDetect(skipped bool) FetchOption {
	return func(o *fetchOptions) {
		withGitSkipAutodetect(skipped)(&o.gitOptions)
	}
}

// FetchWithGitDebug enables debug logging of the underlying git operations.
func FetchWithGitDebug(enabled bool) FetchOption {
	return func(o *fetchOptions) {
		withGitDebug(enabled)(&o.gitOptions)
	}
}

// FetchWithExactTag indicates that tag references are matched exactly.
//
// By default tags are resolved to match the latest semver tag, when a version
// tag is not fully specified, e.g. "v2" would look for the latest "v2.x.y" tag,
// and "v2.1" for the latest "v2.1.y" tag. "v2.3.4" would always resolve to "v2.3.4".
//
// When specifying an exact tag, there is no semver implied or filtering of prereleases.
func FetchWithExactTag(exact bool) FetchOption {
	return func(o *fetchOptions) {
		withGitResolveExactTag(exact)(&o.gitOptions)
	}
}

// FetchWithRequireVersion tells the [Fetcher] to check that the fetched location
// comes with an explicit version. No default to HEAD is applied.
func FetchWithRequireVersion(required bool) FetchOption {
	return func(o *fetchOptions) {
		withRequiredLocVersion(required)(&o.locOptions)
	}
}

// FetchWithSPDXOptions appends SPDX-specific options to apply to any SPDX locator to be fetched.
func FetchWithSPDXOptions(opts ...SPDXOption) FetchOption {
	return func(o *fetchOptions) {
		withSPDXOptions(opts...)(&o.locOptions)
	}
}

// FetchWithGitLocatorOptions appends giturl-specific options to apply to any git-url locator to be fetched.
func FetchWithGitLocatorOptions(opts ...GitLocatorOption) FetchOption {
	return func(o *fetchOptions) {
		withGitLocatorOptions(opts...)(&o.locOptions)
	}
}

// FetchWithSkipRawURL disables the attempt to short-circuit git if a SCM raw-content URL is available
// for the remote resource.
func FetchWithSkipRawURL(skipped bool) FetchOption {
	return func(o *fetchOptions) {
		withSkipRawURL(skipped)(&o.locOptions)
	}
}

// FetchWithAllowPrereleases includes pre-releases in semver tag resolution.
//
// By default pre-releases are ignored.
//
// This option is disabled when using [FetchWithExactTag].
//
// Example:
// for tag "v2", with pre-releases allowed, "v1.3.0-rc1" is a valid candidate.
func FetchWithAllowPrereleases(allowed bool) FetchOption {
	return func(o *fetchOptions) {
		withGitAllowPrereleases(allowed)(&o.gitOptions)
	}
}

// FetchWithRecurseSubmodules resolves submodules when fetching.
//
// By default, git submodules are not updated.
func FetchWithRecurseSubmodules(enabled bool) FetchOption {
	return func(o *fetchOptions) {
		withGitRecurseSubModules(enabled)(&o.gitOptions)
	}
}

type fetchOptions struct {
	gitOptions
	locOptions
}

// CloneOption configures a [Cloner] with optional behavior.
type CloneOption func(*cloneOptions)

// CloneWithBackingDir tells the [Cloner] to back the cloned resources
// on disk. By default, cloned resources are mapped in memory.
//
// If dir is empty, the default is given by [os.MkDirTemp] using "vcsclone" as the pattern.
// In this case, [CloneWithBackingDir] panics if it can't create a temporary directory.
//
// When using [CloneWithBackingDir] with a non-empty directory, the cloned content
// will not be removed after usage and left up to the caller to leave it or clean it if needed.
func CloneWithBackingDir(enabled bool, dir string) CloneOption {
	return func(o *cloneOptions) {
		withGitBackingDir(enabled, dir)(&o.gitOptions)
	}
}

// CloneWithGitSkipAutoDetect skips the auto-detection of a local git binary.
//
// Whenever enabled, git binary autodetection allows for some operations to be performed
// faster using git native implementation rather than the pure go implementation.
func CloneWithGitSkipAutoDetect(skipped bool) CloneOption {
	return func(o *cloneOptions) {
		withGitSkipAutodetect(skipped)(&o.gitOptions)
	}
}

// CloneWithGitDebug enables debug logging of the underlying git operations.
func CloneWithGitDebug(enabled bool) CloneOption {
	return func(o *cloneOptions) {
		withGitDebug(enabled)(&o.gitOptions)
	}
}

// CloneWithExactTag indicates that tag references are matched exactly.
//
// By default tags are resolved to match the latest semver tag, when a version
// tag is not fully specified, e.g. "v2" would look for the latest "v2.x.y" tag,
// and "v2.1" for the latest "v2.1.y" tag. "v2.3.4" would always resolve to "v2.3.4".
//
// When specifying an exact tag, there is no semver implied or filtering of prereleases.
func CloneWithExactTag(exact bool) CloneOption {
	return func(o *cloneOptions) {
		withGitResolveExactTag(exact)(&o.gitOptions)
	}
}

// CloneWithRequireVersion tells the [Cloner] to check that the cloned location
// comes with an explicit version. No default to HEAD is applied.
func CloneWithRequireVersion(required bool) CloneOption {
	return func(o *cloneOptions) {
		withRequiredLocVersion(required)(&o.locOptions)
	}
}

// CloneWithSPDXOptions appends SPDX-specific options to apply to any SPDX locator to be cloned.
func CloneWithSPDXOptions(opts ...SPDXOption) CloneOption {
	return func(o *cloneOptions) {
		withSPDXOptions(opts...)(&o.locOptions)
	}
}

// CloneWithGitLocatorOptions appends giturl-specific options to apply to any git-url locator to be cloned.
func CloneWithGitLocatorOptions(opts ...GitLocatorOption) CloneOption {
	return func(o *cloneOptions) {
		withGitLocatorOptions(opts...)(&o.locOptions)
	}
}

// CloneWithAllowPrereleases includes pre-releases in semver tag resolution.
//
// By default pre-releases are ignored.
//
// This option is disabled when using [CloneWithExactTag].
//
// Example:
// for tag "v2", with pre-releases allowed, "v1.3.0-rc1" is a valid candidate.
func CloneWithAllowPrereleases(allowed bool) CloneOption {
	return func(o *cloneOptions) {
		withGitAllowPrereleases(allowed)(&o.gitOptions)
	}
}

// CloneWithRecurseSubmodules resolves submodules when cloning.
//
// By default, git submodules are not updated.
func CloneWithRecurseSubmodules(enabled bool) CloneOption {
	return func(o *cloneOptions) {
		withGitRecurseSubModules(enabled)(&o.gitOptions)
	}
}

// CloneWithSparseFilter instructs the cloning to be performed only on the specified directories or files.
func CloneWithSparseFilter(filter ...string) CloneOption {
	return func(o *cloneOptions) {
		o.sparseFilter = append(o.sparseFilter, filter...)
	}
}

// SPDXOption is an option to parse a SPDX locator URL.
type SPDXOption func(*spdxOptions)

// GitLocatorOption is an option to parse a git locator (aka git-url).
type GitLocatorOption func(*gitLocatorOptions)

// SPDXWithRootURL declares an URL (as a [url.URL] or as a string) to prepend
// to "slug-like" abbreviated locators.
//
// Example to resolve github repo slugs: rootURL = https://github.com
//
//   - fredbi/go-vcsfetch#README.md resolves a https://github.com/fredbi/go-vcsfetch#README.md
//
// NOTE: [SPDXWithRootURL] panics if the argument passed is a string representing an invalid URL.
func SPDXWithRootURL[T string | *url.URL | url.URL](root T) SPDXOption {
	return func(o *spdxOptions) {
		withRootURL(root)(&o.commonLocOptions)
	}
}

// GitWithRootURL declares an URL (as a [url.URL] or as a string) to prepend
// to "slug-like" abbreviated locators.
//
// Example to resolve github repo slugs: rootURL = https://github.com
//
//   - fredbi/go-vcsfetch#README.md resolves a https://github.com/fredbi/go-vcsfetch#README.md
//
// NOTE: [GitWithRootURL] panics if the argument passed is a string representing an invalid URL.
func GitWithRootURL[T string | *url.URL | url.URL](root T) GitLocatorOption {
	return func(o *gitLocatorOptions) {
		withRootURL(root)(&o.commonLocOptions)
	}
}

// SPDXWithRequiredVersion tells the [SPDXLocator] parser to check that the location
// comes with an explicit version.
func SPDXWithRequiredVersion(required bool) SPDXOption {
	return func(o *spdxOptions) {
		withRequiredVersion(required)(&o.commonLocOptions)
	}
}

// GitWithRequiredVersion tells the [GitLocator] parser to check that the location
// comes with an explicit version.
func GitWithRequiredVersion(required bool) GitLocatorOption {
	return func(o *gitLocatorOptions) {
		withRequiredVersion(required)(&o.commonLocOptions)
	}
}

type cloneOptions struct {
	gitOptions
	locOptions

	sparseFilter []string
}

type gitOption func(*gitOptions)

type gitOptions struct {
	isFSBacked        bool
	dir               string
	gitSkipAutodetect bool
	debug             bool
	resolveExactTag   bool
	allowPrereleases  bool
	recurseSubModules bool
	// auth TODO
}

type locOption func(*locOptions)

type locOptions struct {
	requireVersion bool
	skipRawURL     bool
	spdxOpts       []SPDXOption
	gitLocOpts     []GitLocatorOption
}

type spdxOptions struct {
	commonLocOptions
}

type gitLocatorOptions struct {
	commonLocOptions
}

type commonLocOption func(*commonLocOptions)

type commonLocOptions struct {
	requireVersion  bool
	useSCMshorthand string
	rootURL         *url.URL
}

func withGitBackingDir(enabled bool, dir string) gitOption {
	return func(o *gitOptions) {
		o.isFSBacked = enabled
		if !enabled {
			return
		}

		if dir == "" {
			tempDir, err := os.MkdirTemp("", "vcsclone")
			if err != nil {
				panic(fmt.Errorf("could not created temporary folder to clone: %w: %w", err, ErrVCS))
			}
			o.dir = tempDir
		} else {
			o.dir = dir
		}
	}
}

func withGitSkipAutodetect(skipped bool) gitOption {
	return func(o *gitOptions) {
		o.gitSkipAutodetect = skipped
	}
}

func withGitDebug(enabled bool) gitOption {
	return func(o *gitOptions) {
		o.debug = enabled
	}
}

func withGitResolveExactTag(exact bool) gitOption {
	return func(o *gitOptions) {
		o.resolveExactTag = exact
	}
}

func withGitAllowPrereleases(allowed bool) gitOption {
	return func(o *gitOptions) {
		o.allowPrereleases = allowed
	}
}

func withGitRecurseSubModules(enabled bool) gitOption {
	return func(o *gitOptions) {
		o.recurseSubModules = enabled
	}
}

func withSPDXOptions(opts ...SPDXOption) locOption {
	return func(o *locOptions) {
		o.spdxOpts = append(o.spdxOpts, opts...)
	}
}

func withGitLocatorOptions(opts ...GitLocatorOption) locOption {
	return func(o *locOptions) {
		o.gitLocOpts = append(o.gitLocOpts, opts...)
	}
}

func withRequiredLocVersion(required bool) locOption {
	return func(o *locOptions) {
		o.requireVersion = required
	}
}

func withSkipRawURL(skipped bool) locOption {
	return func(o *locOptions) {
		o.skipRawURL = skipped
	}
}

func withRootURL[T string | *url.URL | url.URL](root T) commonLocOption {
	return func(o *commonLocOptions) {
		var v any = root
		switch value := v.(type) {
		case string:
			u, err := url.Parse(value)
			if err != nil {
				panic(fmt.Errorf("invalid URL string passed as parameter: %q: %w: %w", value, err, ErrVCS))
			}
			o.rootURL = u
		case *url.URL:
			o.rootURL = value
		case url.URL:
			o.rootURL = &value
		}
	}
}

func withRequiredVersion(required bool) commonLocOption {
	return func(o *commonLocOptions) {
		o.requireVersion = required
	}
}

func (o locOptions) toInternalDownloadOptions() *download.Options {
	return &download.Options{}
}

func (o gitOptions) toInternalGitOptions() *git.Options {
	return &git.Options{
		IsFSBacked:        o.isFSBacked,
		Dir:               o.dir,
		GitSkipAutoDetect: o.gitSkipAutodetect,
		Debug:             o.debug,
		ResolveExactTag:   o.resolveExactTag,
	}
}

func (o cloneOptions) toInternalGitCloneOptions() *git.CloneOptions {
	return &git.CloneOptions{
		SparseFilter: o.sparseFilter,
	}
}
