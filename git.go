package vcs

import (
	"os"
	"os/exec"
	"strings"
)

// NewGitRepo creates a new instance of GitRepo. The remote and local directories
// need to be passed in.
func NewGitRepo(remote, local string) (*GitRepo, error) {
	ltype, err := detectVcsFromFS(local)

	// Found a VCS other than Git. Need to report an error.
	if err == nil && ltype != GitType {
		return nil, ErrWrongVCS
	}

	r := &GitRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)
	r.RemoteLocation = "origin"

	return r, nil
}

// GitRepo implements the Repo interface for the Git source control.
type GitRepo struct {
	base
	RemoteLocation string
}

// Get is used to perform an initial clone of a repository.
func (s *GitRepo) Get() error {
	return s.run("git", "clone", s.Remote(), s.LocalPath())
}

// Update performs an Git fetch and pull to an existing checkout.
func (s *GitRepo) Update() error {
	// Perform a fetch to make sure everything is up to date.
	err := s.runFromDir("git", "fetch", s.RemoteLocation)
	if err != nil {
		return err
	}
	return s.runFromDir("git", "pull")
}

// UpdateVersion sets the version of a package currently checked out via Git.
func (s *GitRepo) UpdateVersion(version string) error {
	return s.runFromDir("git", "checkout", version)
}

// Version retrieves the current version.
func (s *GitRepo) Version() (string, error) {

	oldDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(s.LocalPath())
	defer os.Chdir(oldDir)

	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CheckLocal verifies the local location is a Git repo.
func (s *GitRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.git"); err == nil {
		return true
	}

	return false

}
