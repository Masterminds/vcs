package vcs

import (
	"os"
	"os/exec"
	"strings"
)

// NewBzrRepo creates a new instance of BzrRepo. The remote and local directories
// need to be passed in.
func NewBzrRepo(remote, local string) (*BzrRepo, error) {
	ltype, err := detectVcsFromFS(local)

	// Found a VCS other than Bzr. Need to report an error.
	if err == nil && ltype != BzrType {
		return nil, ErrWrongVCS
	}

	r := &BzrRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)

	return r, nil
}

// BzrRepo implements the Repo interface for the Bzr source control.
type BzrRepo struct {
	base
}

// Get is used to perform an initial clone of a repository.
func (s *BzrRepo) Get() error {
	return s.run("bzr", "branch", s.Remote(), s.LocalPath())
}

// Update performs a Bzr pull and update to an existing checkout.
func (s *BzrRepo) Update() error {
	err := s.runFromDir("bzr", "pull")
	if err != nil {
		return err
	}
	return s.runFromDir("bzr", "update")
}

// UpdateVersion sets the version of a package currently checked out via Bzr.
func (s *BzrRepo) UpdateVersion(version string) error {
	return s.runFromDir("bzr", "update", "-r", version)
}

// Version retrieves the current version.
func (s *BzrRepo) Version() (string, error) {

	oldDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(s.LocalPath())
	defer os.Chdir(oldDir)

	out, err := exec.Command("bzr", "revno", "--tree").CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CheckLocal verifies the local location is a Bzr repo.
func (s *BzrRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.bzr"); err == nil {
		return true
	}

	return false
}
