package git

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5/plumbing"
)

const HEAD = "HEAD"

func pickRef(allRefs []*plumbing.Reference, ref string, opts *Options) (*Ref, error) {
	desiredVersion, err := semver.ParseTolerant(ref) // incomplete version specification is completed, e.g. "v2" becomes "2.0.0"
	isDesiredSemver := err == nil
	var versionUpperBound semver.Version
	allowPrereleases := opts != nil && opts.AllowPreReleases
	resolveExactTag := opts != nil && opts.ResolveExactTag

	if isDesiredSemver {
		var allow bool
		desiredSemverLevel := min(strings.Count(ref, "."), 2) + 1
		versionUpperBound, allow = getVersionUpperBound(desiredVersion, desiredSemverLevel)
		allowPrereleases = allowPrereleases || allow
	}

	ctx := &refFilterContext{
		ref:               ref,
		resolveExactTag:   resolveExactTag,
		isDesiredSemver:   isDesiredSemver,
		allowPrereleases:  allowPrereleases,
		versionUpperBound: versionUpperBound,
	}

	refs := make([]Ref, 0, len(allRefs))
	var selectedRef *Ref
	for _, rf := range allRefs {
		localRef, ok := filterRef(ctx, rf)
		if !ok {
			continue
		}
		refs = append(refs, localRef)

		if ref == "" || ref == HEAD || resolveExactTag {
			selectedRef = &localRef
			break
		}
	}

	if len(refs) == 0 {
		return nil, fmt.Errorf("could not resolve any remote reference for ref spec: %q", ref)
	}

	if selectedRef != nil {
		// exact tag match
		return selectedRef, nil
	}

	if len(refs) == 1 {
		selectedRef = &refs[0]
		return selectedRef, nil
	}

	if !isDesiredSemver {
		// this is possible because of semver tolerance, e.g. we may have both tags "v0.2.0" and "0.2.0"
		return nil, fmt.Errorf("ref spec resolved ambiguously to multiple refs: %q", ref)
	}

	// now for selecting among semver candidates
	return latestSemver(refs)
}

func latestSemver(refs []Ref) (*Ref, error) {
	eligibleTags := make([]Ref, 0, len(refs))
	for _, rf := range refs {
		if !rf.IsSemver {
			continue
		}
		eligibleTags = append(eligibleTags, rf)
	}

	if len(eligibleTags) == 0 {
		return nil, fmt.Errorf("no tag did match the version constraint")
	}

	sort.Slice(eligibleTags, func(i, j int) bool {
		return eligibleTags[i].Version.GT(eligibleTags[j].Version) // latest comes first
	})

	tag := eligibleTags[0]
	return &tag, nil
}

type refFilterContext struct {
	ref               string
	resolveExactTag   bool
	isDesiredSemver   bool
	allowPrereleases  bool
	versionUpperBound semver.Version
}

func filterRef(filter *refFilterContext, rf *plumbing.Reference) (localRef Ref, retained bool) {
	if rf.Type() != plumbing.HashReference && rf.Type() != plumbing.SymbolicReference {
		// only consider hash and symbolic references (ignore invalid)
		return localRef, false
	}

	name := rf.Name()
	isTag := name.IsTag()
	if !name.IsBranch() && !isTag && name != plumbing.HEAD {
		// only consider branch, tag and HEAD refs
		return localRef, false
	}

	if (filter.ref == "" || filter.ref == HEAD) && name != plumbing.HEAD {
		// if the desired ref is empty of HEAD, only pick the HEAD branch
		return localRef, false
	}

	if filter.isDesiredSemver && !isTag {
		// if the desired ref is a semver, only consider tag. Ignore branch names that _could_ have a semver name
		return localRef, false
	}

	short := name.Short() // removes the "refs/xxxx/" prefix
	if (filter.resolveExactTag || !filter.isDesiredSemver) && short != filter.ref {
		// if tags must be resolved exactly only consider an exact match
		return localRef, false
	}

	localRef = Ref{
		Reference: rf,
		ShortName: short,
		IsTag:     isTag,
	}

	if isTag {
		version, isVersionErr := semver.ParseTolerant(short)
		if isVersionErr == nil {
			localRef.IsSemver = true
			localRef.Version = version
		}
	}

	if !filter.resolveExactTag && filter.isDesiredSemver {
		if !localRef.IsSemver {
			// reject non-semver tags
			return localRef, false
		}

		// if we disallow pre-releases reject such versions
		if !filter.allowPrereleases && len(localRef.Version.Pre) > 0 {
			return localRef, false
		}

		// if we allow to resolve compatible version tags, reject versions higher than the upper bound
		if localRef.Version.GE(filter.versionUpperBound) {
			return localRef, false
		}
	}

	return localRef, true
}

func getVersionUpperBound(desiredVersion semver.Version, desiredSemverLevel int) (semver.Version, bool) {
	var allowPrereleases bool
	versionUpperBound := desiredVersion // shallow clone: upper bound (excluded) for select tagged version
	versionUpperBound.Pre = nil
	versionUpperBound.Build = nil

	finalized := desiredVersion
	finalized.Pre = nil
	finalized.Build = nil
	if desiredVersion.GE(finalized) {
		allowPrereleases = true // the ref spec containes a pre-release: imply that we accept those
	}

	switch desiredSemverLevel {
	case 3: // fully specified
		_ = versionUpperBound.IncrementPatch()
	case 2: // allow for patches
		_ = versionUpperBound.IncrementMinor()
	case 1: // allow for minor versions
		_ = versionUpperBound.IncrementMajor()
	}

	return versionUpperBound, allowPrereleases
}
