package git

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	/*
		"github.com/go-git/go-billy/v5/helper/iofs"
	*/)

type Options struct {
	IsFSBacked        bool
	Dir               string
	ResolveExactTag   bool
	RecurseSubModules bool // TODO
	AllowPreReleases  bool
	Debug             bool
	GitSkipAutoDetect bool
	// Auth
	// TLS
	// Proxy
}

type CloneOptions struct {
	SparseFilter []string
}

type Ref struct {
	*plumbing.Reference
	ShortName string
	IsTag     bool
	IsSemver  bool
	Version   semver.Version
}

type Repository struct {
	*Options

	repoURL  *url.URL
	repo     *gogit.Repository
	store    func() storage.Storer
	worktree func() billy.Filesystem
}

func NewRepo(repoURL *url.URL, opts *Options) *Repository {
	var (
		storeFunc    func() storage.Storer
		worktreeFunc func() billy.Filesystem
	)

	if opts != nil && opts.IsFSBacked && opts.Dir != "" {
		// optional osFS-backend
		fs := osfs.New(opts.Dir, osfs.WithBoundOS())
		lru := cache.NewObjectLRUDefault()

		storeFunc = func() storage.Storer {
			lru.Clear()

			return filesystem.NewStorage(fs, lru)
		}
		worktreeFunc = func() billy.Filesystem {
			fs.(*osfs.BoundOS).RemoveAll(fs.Root())

			return fs
		}
	} else {
		// default memFS-backend
		storeFunc = func() storage.Storer { return memory.NewStorage() }
		worktreeFunc = memfs.New
	}

	return &Repository{
		Options:  opts,
		repoURL:  repoURL,
		store:    storeFunc,
		worktree: worktreeFunc,
	}
}

func (r *Repository) Fetch(ctx context.Context, w io.Writer, file, ref string) error {
	t0 := time.Now()
	repo, remote, err := r.init()
	if err != nil {
		return fmt.Errorf("could not initialize git repo: %w", err)
	}
	t1 := time.Now()
	log.Printf("init: elapsed: %v", t1.Sub(t0))

	selectedRef, err := r.selectRef(ctx, remote, ref)
	if err != nil {
		return fmt.Errorf("could not resolve remote ref: %w", err)
	}
	t2 := time.Now()
	log.Printf("select: elapsed: %v", t2.Sub(t1))

	hash := selectedRef.Hash()
	err = r.fetch(ctx, remote, hash, file)
	if err != nil {
		return fmt.Errorf("could not fetch remote ref: %w", err)
	}
	t3 := time.Now()
	log.Printf("fetch: elapsed: %v", t3.Sub(t2))

	local, err := repo.Worktree()
	if err != nil {
		return err
	}

	// at this point we should only have a hash
	var filter []string
	dir := filepath.Dir(file)
	if dir != "." && dir != "/" {
		filter = []string{file}
	}
	err = local.Checkout(&gogit.CheckoutOptions{
		Hash:                      hash,
		Branch:                    selectedRef.Name(),
		Create:                    true,
		Force:                     true,
		SparseCheckoutDirectories: filter,
	})
	if err != nil {
		return err
	}
	t4 := time.Now()
	log.Printf("checkout: elapsed: %v", t4.Sub(t3))
	/*
		if err = fs.WalkDir(&fsWrapper{Filesystem: local.Filesystem}, "/", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			log.Printf("DEBUG: %s", path)
			return nil
		}); err != nil {
			return fmt.Errorf("DEBUG: walkdir error: %w", err)
		}
	*/

	path := filepath.Join(local.Filesystem.Root(), file)
	fd, err := local.Filesystem.Open(path)
	if err != nil {
		return fmt.Errorf("did not find %q on checkout: %w", path, err)
	}

	_, err = io.Copy(w, fd)
	t5 := time.Now()
	log.Printf("copy: elapsed: %v", t5.Sub(t4))

	return err

	// TODO: auto detect if git is installed, and if so, switch to the git archive command
	// Can we do git archive ? => not supported by go git
	/*
		git archive --remote=$REPO_URL HEAD:path/to -- file.xz |
		tar xO > /where/you/want/to/have.it
	*/
}

// Clone the repository defined by an URL
func (r *Repository) Clone(ctx context.Context, ref string, opts *CloneOptions) (fs.FS, error) {
	return nil, nil
	/*
		// Branches and tags are safe to fetch when cloning. This is not the case
		// of notes, for example so we only pass a reference to clone if we're
		// dealing with a brach or tag.
		var reference plumbing.ReferenceName
		switch {
		case components.Branch != "":
			reference = plumbing.NewBranchReferenceName(components.Branch)
		case components.Tag != "":
			reference = plumbing.NewTagReferenceName(components.Tag)
		}

		var fsobj billy.Filesystem
		if opts.ClonePath == "" {
			fsobj = memfs.New()
		} else {
			fsobj = osfs.New(opts.ClonePath)
		}

		// Handle cloning from repos with file: transport
		repourl := components.RepoURL()
		if components.Transport == "file" {
			repourl = components.RepoPath
		}

		// Make a shallow clone of the repo to memory
		if len(opts.Filter) > 0 {

		}
		repo, err := git.Clone(memory.NewStorage(), fsobj, &git.CloneOptions{
			URL: repourl,
			// Progress:      os.Stdout,
			ReferenceName: reference,
			SingleBranch:  true,
			// Depth:         1,
			// RecurseSubmodules: 0,
			// ShallowSubmodules: false,
			// TODO(fred): depth
			// TODO(fred): how to achieve sparse checkout?
		})
		if err != nil {
			return nil, fmt.Errorf("cloning repo: %w", err)
		}

		commitHash := components.Commit
		// Here we handle commits and other references (not tags or branches)
		if reference == "" && components.Commit == "" {
			// But also ensuring we are note refetching a previous commit
			if components.RefString != "" && components.RefString != components.Commit {
				// Since this ref was not fetched at clone time, we do a fetch here
				// to make sure it is available. This is especially important for
				// git notes that are never transferred by default and cannot be
				// fetched at clone time, I thing because of a bug that somewhere
				// changes the ref string from refs/notes/commits to refs/heads/notes/commits
				//
				if err := repo.Fetch(&git.FetchOptions{
					RefSpecs: []config.RefSpec{
						config.RefSpec(fmt.Sprintf("%s:%s", components.RefString, components.RefString)),
					},
				}); err != nil {
					return nil, fmt.Errorf("late fetching ref %q: %w", components.RefString, err)
				}

				// Resolve the reference, it should not fail as we fetched it already
				ref, err := repo.Reference(plumbing.ReferenceName(components.RefString), true)
				if err != nil {
					return nil, fmt.Errorf("resolving reference %q: %w", components.RefString, err)
				}

				// Resolve the reference to a commit hash
				hach, err := repo.ResolveRevision(plumbing.Revision(ref.Name().String()))
				if err != nil {
					return nil, fmt.Errorf("resolving latest revision on %q to commit: %w", ref.Name().String(), err)
				}
				commitHash = hach.String()
			}
		}

		// If a revision was specified, check it out
		if commitHash != "" {
			wt, err := repo.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting repository worktree: %w", err)
			}

			if err = wt.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(commitHash),
			}); err != nil {
				return nil, fmt.Errorf("checking out commit %s: %w", commitHash, err)
			}
		}

		return iofs.New(fsobj), nil
	*/
}

func (r *Repository) init() (*gogit.Repository, *gogit.Remote, error) {
	repo, err := gogit.Init(r.store(), r.worktree())
	if err != nil {
		return nil, nil, err
	}

	// TODO: config (auth, ...)

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{r.repoURL.String()},
	})
	if err != nil {
		return nil, nil, err
	}

	return repo, remote, nil
}

func (r *Repository) selectRef(ctx context.Context, remote *gogit.Remote, ref string) (*Ref, error) {
	allRefs, err := remote.ListContext(ctx, &gogit.ListOptions{ // NOTE: unfortunately, there is no way to filter refs
		// Auth / TLS/ Proxy
	})
	if err != nil {
		return nil, err
	}

	desiredVersion, err := semver.ParseTolerant(ref) // incomplete version specification is completed, e.g. "v2" becomes "2.0.0"
	isDesiredSemver := err == nil
	isDesiredSemverLevel := 0

	versionUpperBound := desiredVersion // shallow clone: upper bound (excluded) for select tagged version
	versionUpperBound.Pre = nil
	versionUpperBound.Build = nil

	allowPrereleases := r.Options != nil && r.AllowPreReleases
	resolveExactTag := r.Options != nil && r.ResolveExactTag

	if isDesiredSemver {
		finalized := desiredVersion
		finalized.Pre = nil
		finalized.Build = nil
		if desiredVersion.GE(finalized) {
			allowPrereleases = true // the ref spec containes a pre-release: imply that we accept those
		}
		isDesiredSemverLevel = min(strings.Count(ref, "."), 2) + 1
		switch isDesiredSemverLevel {
		case 3: // fully specified
			_ = versionUpperBound.IncrementPatch()
		case 2: // allow for patches
			_ = versionUpperBound.IncrementMinor()
		case 1: // allow for minor versions
			_ = versionUpperBound.IncrementMajor()
		}
	}

	const HEAD = "HEAD"
	refs := make([]Ref, 0, len(allRefs))
	var selectedRef *Ref
	for _, rf := range allRefs {
		if rf.Type() != plumbing.HashReference && rf.Type() != plumbing.SymbolicReference {
			// only consider hash and symbolic references (ignore invalid)
			continue
		}

		name := rf.Name()
		isTag := name.IsTag()
		if !name.IsBranch() && !isTag && name != plumbing.HEAD {
			// only consider branch, tag and HEAD refs
			continue
		}

		if (ref == "" || ref == HEAD) && name != plumbing.HEAD {
			// if the desired ref is empty of HEAD, only pick the HEAD branch
			continue
		}

		if isDesiredSemver && !isTag {
			// if the desired ref is a semver, only consider tag. Ignore branch names that _could_ have a semver name
			continue
		}

		short := name.Short() // removes the "refs/xxxx/" prefix
		if (resolveExactTag || !isDesiredSemver) && short != ref {
			// if tags must be resolved exactly only consider an exact match
			continue
		}

		localRef := Ref{
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

		if !resolveExactTag && isDesiredSemver {
			if !localRef.IsSemver {
				// reject non-semver tags
				continue
			}

			// if we disallow pre-releases reject such versions
			if !allowPrereleases && len(localRef.Version.Pre) > 0 {
				continue
			}

			// if we allow to resolve compatible version tags, reject versions higher than the upper bound
			if localRef.Version.GE(versionUpperBound) {
				continue
			}
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

	if selectedRef == nil {
		// now for selecting among semver candidates
		if len(refs) > 1 && isDesiredSemver {
			eligibleTags := make([]Ref, 0, len(refs))
			for _, rf := range refs {
				if !rf.IsSemver {
					continue
				}
				eligibleTags = append(eligibleTags, rf)
			}
			if len(eligibleTags) == 0 {
				return nil, fmt.Errorf("no tag did match this constraint: %q", ref)
			}

			sort.Slice(eligibleTags, func(i, j int) bool {
				return eligibleTags[i].Version.GT(eligibleTags[j].Version) // latest comes first
			})

			tag := eligibleTags[0]
			selectedRef = &tag
		} else {
			// this is possible because of semver tolerance, e.g. we may have both tags "v0.2.0" and "0.2.0"
			return nil, fmt.Errorf("ref spec resolved ambiguously to multiple refs: %q", ref)
		}
	}

	return selectedRef, nil
}

func (r *Repository) fetch(ctx context.Context, remote *gogit.Remote, hash plumbing.Hash, file string) error {
	_ = file

	refSpec := config.RefSpec(fmt.Sprintf("+%[1]v:%[1]v", hash)) // build a hash ref
	err := remote.FetchContext(ctx, &gogit.FetchOptions{
		RefSpecs: []config.RefSpec{refSpec},
		Depth:    0,
		Tags:     gogit.NoTags,
		Force:    true,
		// Auth: / TLS / Proxy
	})
	if err != nil {
		return fmt.Errorf("fetch remote hash ref %v: %w", hash, err)
	}

	// TODO: if local fs, use Storer.AddAlternate?
	// RecurseSubModules???

	/*
		branch := "" // remote branch
			// required?
			err = repo.CreateBranch(&config.Branch{
				Name:   branch,
				Remote: remote.Config().Name,
			})
	*/
	return nil
}
