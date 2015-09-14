package vcs

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var svnDetectURL = regexp.MustCompile("URL: (?P<foo>.+)\n")

// NewSvnRepo creates a new instance of SvnRepo. The remote and local directories
// need to be passed in. The remote location should include the branch for SVN.
// For example, if the package is https://github.com/Masterminds/cookoo/ the remote
// should be https://github.com/Masterminds/cookoo/trunk for the trunk branch.
func NewSvnRepo(remote, local string) (*SvnRepo, error) {
	ltype, err := DetectVcsFromFS(local)

	// Found a VCS other than Svn. Need to report an error.
	if err == nil && ltype != Svn {
		return nil, ErrWrongVCS
	}

	r := &SvnRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)
	r.Logger = Logger

	// Make sure the local SVN repo is configured the same as the remote when
	// A remote value was passed in.
	if err == nil && r.CheckLocal() == true {
		// An SVN repo was found so test that the URL there matches
		// the repo passed in here.
		out, err := exec.Command("svn", "info", local).CombinedOutput()
		if err != nil {
			return nil, err
		}

		m := svnDetectURL.FindStringSubmatch(string(out))
		if m[1] != "" && m[1] != remote {
			return nil, ErrWrongRemote
		}

		// If no remote was passed in but one is configured for the locally
		// checked out Svn repo use that one.
		if remote == "" && m[1] != "" {
			r.setRemote(m[1])
		}
	}

	return r, nil
}

// SvnRepo implements the Repo interface for the Svn source control.
type SvnRepo struct {
	base
}

// Vcs retrieves the underlying VCS being implemented.
func (s SvnRepo) Vcs() Type {
	return Svn
}

// Get is used to perform an initial checkout of a repository.
// Note, because SVN isn't distributed this is a checkout without
// a clone.
func (s *SvnRepo) Get() error {
	_, err := s.run("svn", "checkout", s.Remote(), s.LocalPath())
	return err
}

// Update performs an SVN update to an existing checkout.
func (s *SvnRepo) Update() error {
	_, err := s.runFromDir("svn", "update")
	return err
}

// UpdateVersion sets the version of a package currently checked out via SVN.
func (s *SvnRepo) UpdateVersion(version string) error {
	_, err := s.runFromDir("svn", "update", "-r", version)
	return err
}

// Version retrieves the current version.
func (s *SvnRepo) Version() (string, error) {
	out, err := s.runFromDir("svnversion", ".")
	s.log(out)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CheckLocal verifies the local location is an SVN repo.
func (s *SvnRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.svn"); err == nil {
		return true
	}

	return false

}
