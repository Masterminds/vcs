package vcs

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var svnDetectUrl = regexp.MustCompile("URL: (?P<foo>.+)\n")

// NewSvnRepo creates a new instance of SvnRepo. The remote and local directories
// need to be passed in. The remote location should include the branch for SVN.
// For example, if the package is https://github.com/Masterminds/cookoo/ the remote
// should be https://github.com/Masterminds/cookoo/trunk for the trunk branch.
func NewSvnRepo(remote, local string) (*SvnRepo, error) {
	ltype, err := detectVcsFromFS(local)

	// Found a VCS other than Svn. Need to report an error.
	if err == nil && ltype != SvnType {
		return nil, ErrWrongVCS
	} else if err == nil {
		// An SVN repo was found so test that the URL there matches
		// the repo passed in here.
		out, err := exec.Command("svn", "info", local).CombinedOutput()
		if err != nil {
			return nil, err
		}

		m := svnDetectUrl.FindStringSubmatch(string(out))
		if m[1] != "" && m[1] != remote {
			return nil, ErrWrongRemote
		}
	}

	r := &SvnRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)

	return r, nil
}

// SvnRepo implements the Repo interface for the Svn source control.
type SvnRepo struct {
	base
}

// Vcs retrieves the underlying VCS being implemented.
func (s SvnRepo) Vcs() VcsType {
	return SvnType
}

// Get is used to perform an initial checkout of a repository.
// Note, because SVN isn't distributed this is a checkout without
// a clone.
func (s *SvnRepo) Get() error {
	return s.run("svn", "checkout", s.Remote(), s.LocalPath())
}

// Update performs an SVN update to an existing checkout.
func (s *SvnRepo) Update() error {
	return s.runFromDir("svn", "update")
}

// UpdateVersion sets the version of a package currently checked out via SVN.
func (s *SvnRepo) UpdateVersion(version string) error {
	return s.runFromDir("svn", "update", "-r", version)
}

// Version retrieves the current version.
func (s *SvnRepo) Version() (string, error) {

	oldDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(s.LocalPath())
	defer os.Chdir(oldDir)

	out, err := exec.Command("svnversion", ".").CombinedOutput()
	l(out)
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
