package vcs

import (
	"os"
	"os/exec"
	"strings"
)

// SvnRepo implements the Repo interface for the Svn source control.
type SvnRepo struct {
	helpers
	remote, local string
}

// Remote retrieves the remote location for a repo.
func (s *SvnRepo) Remote() string {
	return s.remote
}

// LocalPath retrieves the local file system location for a repo.
func (s *SvnRepo) LocalPath() string {
	return s.local
}

// Get is used to perform an initial checkout of a repository.
// Note, because SVN isn't distributed this is a checkout without
// a clone.
func (s *SvnRepo) Get() error {
	out, err := exec.Command("svn", "checkout", s.remote, s.local).CombinedOutput()
	s.logger.Print(out)
	return err
}

// Update performs an SVN update to an existing checkout.
func (s *SvnRepo) Update() error {

	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(s.local)
	defer os.Chdir(oldDir)

	out, err := exec.Command("svn", "update").CombinedOutput()
	s.logger.Print(out)
	return err
}

// UpdateVersion sets the version of a package currently checked out via SVN.
func (s *SvnRepo) UpdateVersion(version string) error {

	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(s.local)
	defer os.Chdir(oldDir)

	out, err := exec.Command("svn", "update", "-r", version).CombinedOutput()
	s.logger.Print(out)
	return err
}

// Version retrieves the current version.
func (s *SvnRepo) Version() (string, error) {

	oldDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(s.local)
	defer os.Chdir(oldDir)

	out, err := exec.Command("svnversion", ".").CombinedOutput()
	s.logger.Print(out)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
