package git

// Options for a git [Repository]
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

// / CloneOptions to tune the behavior of git clone.
type CloneOptions struct {
	SparseFilter []string
}
