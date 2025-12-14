package git

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/blang/semver/v4"
	"github.com/davecgh/go-spew/spew"
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
)

// Ref wraps a git [plumbing.Reference].
type Ref struct {
	*plumbing.Reference

	ShortName string
	IsTag     bool
	IsSemver  bool
	Version   semver.Version
}

// Repository is a git repo.
type Repository struct {
	*Options

	repoURL  *url.URL
	repo     *gogit.Repository
	store    func() storage.Storer
	worktree func() billy.Filesystem
	debug    func(string, ...any)
}

// NewRepo initializes a new git repository for a given URL.
//
// No resources are actually fetched or stored yet.
func NewRepo(repoURL *url.URL, opts *Options) *Repository {
	var debug func(string, ...any)

	if opts != nil && opts.Debug {
		debug = log.Printf
	} else {
		debug = noDebug
	}

	if opts != nil && opts.IsFSBacked && opts.Dir != "" {
		// optional osFS-backend
		fs := osfs.New(opts.Dir, osfs.WithBoundOS())
		lru := cache.NewObjectLRUDefault()

		initStoreFunc := func() storage.Storer {
			lru.Clear()

			return filesystem.NewStorage(fs, lru)
		}
		initWorktreeFunc := func() billy.Filesystem {
			fs.(*osfs.BoundOS).RemoveAll(fs.Root())

			return fs
		}

		return &Repository{
			Options:  opts,
			repoURL:  repoURL,
			store:    initStoreFunc,
			worktree: initWorktreeFunc,
			debug:    debug,
		}
	}

	// default is MemFS backend
	initStoreFunc := func() storage.Storer { return memory.NewStorage() }
	initWorktreeFunc := memfs.New

	return &Repository{
		Options:  opts,
		repoURL:  repoURL,
		store:    initStoreFunc,
		worktree: initWorktreeFunc,
		debug:    debug,
	}
}

// Fetch a file at a given ref from the [Repository].
//
// The file is copied to the given [io.Writer].
func (r *Repository) Fetch(ctx context.Context, w io.Writer, file, ref string) error {
	// initialize git with proper remote
	repo, remote, err := r.init()
	if err != nil {
		return fmt.Errorf("could not initialize git repo: %w", err)
	}

	// figure out the hash for the desired ref
	selectedRef, err := r.selectRef(ctx, remote, ref)
	if err != nil {
		return fmt.Errorf("could not resolve remote ref: %w", err)
	}

	remoteCapabilities, err := getRemoteCapabilities(ctx, &gogit.FetchOptions{
		RemoteURL: r.repoURL.String(),
	})
	if err != nil {
		return fmt.Errorf("unable to retrieve the git protocol capabilities for the remote server: %w", err)
	}
	spew.Dump(remoteCapabilities)

	if r.Options == nil || !r.GitSkipAutoDetect {
		if r.supportArchive() && isGitInstalled() {
			r.debug("git is installed")
			// use installed git command
			return r.nativeExtractGitArchive(ctx, w, file, selectedRef)
		}
	}

	// use go-git implementation
	return r.fetchAndSparseCheckout(ctx, repo, remote, w, file, selectedRef)
}

func (r *Repository) supportArchive() bool {
	if r.repoURL.Scheme != "git" && r.repoURL.Scheme != "ssh" {
		return false
	}

	return true
}

func (r *Repository) fetchAndSparseCheckout(ctx context.Context, repo *gogit.Repository, remote *gogit.Remote, w io.Writer, file string, selectedRef *Ref) error {
	// fetch ref
	t2 := time.Now()
	hash := selectedRef.Hash()
	if err := r.fetch(ctx, remote, hash, file); err != nil {
		return fmt.Errorf("could not fetch remote ref: %w", err)
	}
	t3 := time.Now()
	r.debug("fetch: elapsed: %v", t3.Sub(t2))

	local, err := repo.Worktree()
	if err != nil {
		return err
	}

	// sparse checkout of the file.
	// At this point we should only have a hash
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
	r.debug("checkout: elapsed: %v", t4.Sub(t3))

	path := filepath.Join(local.Filesystem.Root(), file)
	fd, err := local.Filesystem.Open(path)
	if err != nil {
		return fmt.Errorf("did not find %q on checkout: %w", path, err)
	}

	_, err = io.Copy(w, fd)
	t5 := time.Now()
	r.debug("copy: elapsed: %v", t5.Sub(t4))

	return err
}

// Clone the repository defined by an URL.
func (r *Repository) Clone(ctx context.Context, ref string, opts *CloneOptions) (fs.FS, error) {
	// TODO: clone repo as fs.FS
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
	if r.repoURL == nil || r.repoURL.String() == "" {
		return nil, nil, fmt.Errorf("cannot init repo with empty URL")
	}

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

	// pick the best matching ref depending on chosen options
	return pickRef(allRefs, ref, r.Options)
}

func (r *Repository) fetch(ctx context.Context, remote *gogit.Remote, hash plumbing.Hash, file string) error {
	_ = file

	refSpec := config.RefSpec(fmt.Sprintf("+%[1]v:%[1]v", hash)) // build a hash ref
	err := remote.FetchContext(ctx, &gogit.FetchOptions{         // TODO: bug if repo maps HEAD to main (see gitlab test)
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

func noDebug(format string, args ...any) {
}

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
