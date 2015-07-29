package vcs

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Logger is where you can provide a logger, implementing the log.Logger interface,
// where verbose output from each VCS will be written. The default logger does
// not log data. To log data supply your own logger or change the output location
// of the provided logger.
var Logger *log.Logger = log.New(ioutil.Discard, "go-vcs", log.LstdFlags)

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

	// CheckLocal verifies the local location is of the correct VCS type
	CheckLocal() bool
}

func l(v interface{}) {
	Logger.Printf("%s", v)
}

type base struct {
	remote, local string
}

// Remote retrieves the remote location for a repo.
func (b *base) Remote() string {
	return b.remote
}

// LocalPath retrieves the local file system location for a repo.
func (b *base) LocalPath() string {
	return b.local
}

func (b *base) setRemote(remote string) {
	b.remote = remote
}

func (b *base) setLocalPath(local string) {
	b.local = local
}

func (b base) run(cmd string, args ...string) error {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	l(out)
	return err
}

func (b *base) runFromDir(cmd string, args ...string) error {
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(b.local)
	defer os.Chdir(oldDir)

	err = b.run(cmd, args...)

	return err
}
