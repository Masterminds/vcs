package vcs

import (
	"os"
	"os/exec"
	"strings"
)

// NewHgRepo creates a new instance of HgRepo. The remote and local directories
// need to be passed in.
func NewHgRepo(remote, local string) *HgRepo {
	r := &HgRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)

	return r
}

// HgRepo implements the Repo interface for the Mercurial source control.
type HgRepo struct {
	base
}

// Get is used to perform an initial clone of a repository.
func (s *HgRepo) Get() error {
	return s.run("hg", "clone", "-U", s.Remote(), s.LocalPath())
}

// Update performs a Mercurial pull to an existing checkout.
func (s *HgRepo) Update() error {
	return s.runFromDir("hg", "update")
}

// UpdateVersion sets the version of a package currently checked out via Hg.
func (s *HgRepo) UpdateVersion(version string) error {
	return s.runFromDir("hg", "update", version)
}

// Version retrieves the current version.
func (s *HgRepo) Version() (string, error) {

	oldDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(s.LocalPath())
	defer os.Chdir(oldDir)

	out, err := exec.Command("hg", "identify").CombinedOutput()
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(out), " ", 2)
	sha := parts[0]
	return strings.TrimSpace(sha), nil
}

// CheckLocal verifies the local location is a Git repo.
func (s *HgRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.hg"); err == nil {
		return true
	}

	return false

}
