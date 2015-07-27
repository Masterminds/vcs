package vcs

import (
	"log"
)

// Logger is where you can provide a logger, implementing the log.Logger interface,
// where verbose output from each VCS will be written.
var Logger log.Logger

// Repo provides an interface to work with repositories using different source
// control systems such as Git, Bzr, Mercurial, and SVN. For implementations
// of this interface see BzrRepo, GitRepo, HgRepo, and SvnRepo.
type Repo interface {

	// Remote retrieves the remote location for a repo.
	Remote() string

	// LocalPath retrieves the local file system location for a repo.
	LocalPath() string

	// Get is used to perform an initial clone/checkout of a repository.
	Get() error

	// Update performs an update to an existing checkout of a repository.
	Update() error

	// UpdateVersion sets the version of a package of a repository.
	UpdateVersion(string) error

	// Version retrieves the current version.
	Version() (string, error)
}

type helpers struct {
	logger log.Logger
}
