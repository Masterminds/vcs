// Package vcs provides the ability to work with varying version control systems
// (VCS),  also known as source control systems (SCM) though the same interface.
//
// This package includes a function that attempts to detect the repo type from
// the remote URL and return the proper type. For example,
//
//     remote := "https://github.com/Masterminds/go-vcs"
//     local, _ := ioutil.TempDir("", "go-vcs")
//     repo, err := NewRepo(remote, local)
//
// In this case repo will be a GitRepo instance. NewRepo can detect the VCS for
// numerous popular VCS and from the URL. For example, a URL ending in .git
// that's not from one of the popular VCS will be detected as a Git repo and
// the correct type will be returned.
//
// If you know the repository type and would like to create an instance of a
// specific type you can use one of constructors for a type. They are NewGitRepo,
// NewSvnRepo, NewBzrRepo, and NewHgRepo. The definition and usage is the same
// as NewRepo.
//
// Once you have an object implementing the Repo interface the operations are
// the same no matter which VCS you're using. There are some caveats. For
// example, each VCS has its own version formats that need to be respected and
// checkout out branches, if a branch is being worked with, is different in
// each VCS.
package vcs

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	// ErrWrongVCS is returned when an action is tried on the wrong VCS.
	ErrWrongVCS = errors.New("Wrong VCS detected")

	// ErrCannotDetectVCS is returned when VCS cannot be detected from URI string.
	ErrCannotDetectVCS = errors.New("Cannot detect VCS")

	// ErrWrongRemote occurs when the passed in remote does not match the VCS
	// configured endpoint.
	ErrWrongRemote = errors.New("The Remote does not match the VCS endpoint")
)

// Logger is where you can provide a logger, implementing the log.Logger interface,
// where verbose output from each VCS will be written. The default logger does
// not log data. To log data supply your own logger or change the output location
// of the provided logger.
var Logger *log.Logger

func init() {
	// Initialize the logger to one that does not actually log anywhere. This is
	// to be overridden by the package user by setting vcs.Logger to a different
	// logger.
	Logger = log.New(ioutil.Discard, "go-vcs", log.LstdFlags)
}

// Type describes the type of VCS
type Type string

// VCS types
const (
	Git Type = "git"
	Svn Type = "svn"
	Bzr Type = "bzr"
	Hg  Type = "hg"
)

// Repo provides an interface to work with repositories using different source
// control systems such as Git, Bzr, Mercurial, and SVN. For implementations
// of this interface see BzrRepo, GitRepo, HgRepo, and SvnRepo.
type Repo interface {

	// Vcs retrieves the underlying VCS being implemented.
	Vcs() Type

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

// NewRepo returns a Repo based on trying to detect the source control from the
// remote and local locations. The appropriate implementation will be returned
// or an ErrCannotDetectVCS if the VCS type cannot be detected.
// Note, this function may make calls to the Internet to determind help determine
// the VCS.
func NewRepo(remote, local string) (Repo, error) {
	vtype, err := detectVcsFromFS(local)

	// When the VCS cannot be detected from the local checkout attempt to
	// determine the type from the remote url. Note, some remote urls such
	// as bitbucket require going out to the Internet to detect the type.
	if err == ErrCannotDetectVCS {
		vtype, remote, err = detectVcsFromRemote(remote)
	}

	if err != nil {
		return nil, err
	}

	switch vtype {
	case Git:
		return NewGitRepo(remote, local)
	case Svn:
		return NewSvnRepo(remote, local)
	case Hg:
		return NewHgRepo(remote, local)
	case Bzr:
		return NewBzrRepo(remote, local)
	}

	// Should never fall through to here but just in case.
	return nil, ErrCannotDetectVCS
}

type base struct {
	remote, local string
	Logger        *log.Logger
}

func (b *base) log(v interface{}) {
	b.Logger.Printf("%s", v)
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
	b.log(out)
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
