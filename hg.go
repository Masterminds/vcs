package vcs

import (
	"bytes"
	"encoding/xml"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var hgDetectURL = regexp.MustCompile("default = (?P<foo>.+)\n")

// NewHgRepo creates a new instance of HgRepo. The remote and local directories
// need to be passed in.
func NewHgRepo(remote, local string) (*HgRepo, error) {
	ltype, err := DetectVcsFromFS(local)

	// Found a VCS other than Hg. Need to report an error.
	if err == nil && ltype != Hg {
		return nil, ErrWrongVCS
	}

	r := &HgRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)
	r.Logger = Logger

	// Make sure the local Hg repo is configured the same as the remote when
	// A remote value was passed in.
	if err == nil && r.CheckLocal() == true {
		// An Hg repo was found so test that the URL there matches
		// the repo passed in here.
		c := exec.Command("hg", "paths")
		c.Dir = local
		c.Env = envForDir(c.Dir)
		out, err := c.CombinedOutput()
		if err != nil {
			return nil, err
		}

		m := hgDetectURL.FindStringSubmatch(string(out))
		if m[1] != "" && m[1] != remote {
			return nil, ErrWrongRemote
		}

		// If no remote was passed in but one is configured for the locally
		// checked out Hg repo use that one.
		if remote == "" && m[1] != "" {
			r.setRemote(m[1])
		}
	}

	return r, nil
}

// HgRepo implements the Repo interface for the Mercurial source control.
type HgRepo struct {
	base
}

// Vcs retrieves the underlying VCS being implemented.
func (s HgRepo) Vcs() Type {
	return Hg
}

// Get is used to perform an initial clone of a repository.
func (s *HgRepo) Get() error {
	_, err := s.run("hg", "clone", s.Remote(), s.LocalPath())
	return err
}

// Update performs a Mercurial pull to an existing checkout.
func (s *HgRepo) Update() error {
	_, err := s.runFromDir("hg", "update")
	return err
}

// UpdateVersion sets the version of a package currently checked out via Hg.
func (s *HgRepo) UpdateVersion(version string) error {
	_, err := s.runFromDir("hg", "pull")
	if err != nil {
		return err
	}
	_, err = s.runFromDir("hg", "update", version)
	return err
}

// Version retrieves the current version.
func (s *HgRepo) Version() (string, error) {
	out, err := s.runFromDir("hg", "identify")
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(out), " ", 2)
	sha := parts[0]
	return strings.TrimSpace(sha), nil
}

// Date retrieves the date on the latest commit.
func (s *HgRepo) Date() (time.Time, error) {
	version, err := s.Version()
	if err != nil {
		return time.Time{}, err
	}
	out, err := s.runFromDir("hg", "log", "-r", version, "--template", "{date|isodatesec}")
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(longForm, string(out))
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// CheckLocal verifies the local location is a Git repo.
func (s *HgRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.hg"); err == nil {
		return true
	}

	return false
}

// Branches returns a list of available branches
func (s *HgRepo) Branches() ([]string, error) {
	out, err := s.runFromDir("hg", "branches")
	if err != nil {
		return []string{}, err
	}
	branches := s.referenceList(string(out), `(?m-s)^(\S+)`)
	return branches, nil
}

// Tags returns a list of available tags
func (s *HgRepo) Tags() ([]string, error) {
	out, err := s.runFromDir("hg", "tags")
	if err != nil {
		return []string{}, err
	}
	tags := s.referenceList(string(out), `(?m-s)^(\S+)`)
	return tags, nil
}

// CurrentVersionsWithRevs returns a list of available branches and tags, and
// their underlying revision identifiers. The list is guaranteed to be current.
//
// TODO currently, this only retrieves the tags known to whatever the checked
// out branch happens to be. This is...kinda broken?
func (s *HgRepo) CurrentVersionsWithRevs() (versions []VersionInfo, localSynced bool, err error) {
	var out []byte
	err = s.Update()
	if err != nil {
		return
	}

	localSynced = true
	out, err = s.runFromDir("hg", "tags", "--debug", "--verbose")
	if err != nil {
		return
	}

	all := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	lbyt := []byte("local")
	nulrev := []byte("0000000000000000000000000000000000000000")
	for _, line := range all {
		if bytes.Equal(lbyt, line[len(line)-len(lbyt):]) {
			// Skip local tags
			continue
		}

		// tip is magic, don't include it
		if bytes.HasPrefix(line, []byte("tip")) {
			continue
		}

		// Split on colon; this gets us the rev and the tag plus local revno
		pair := bytes.Split(line, []byte(":"))
		if bytes.Equal(nulrev, pair[1]) {
			// null rev indicates this tag is marked for deletion
			continue
		}

		idx := bytes.IndexByte(pair[0], 32) // space
		versions = append(versions, VersionInfo{
			Name:     string(pair[0][:idx]),
			Revision: string(pair[1]),
		})
	}

	out, err = s.runFromDir("hg", "branches", "--debug", "--verbose")
	if err != nil {
		// better nothing than incomplete
		versions = nil
		return
	}

	all = bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	lbyt = []byte("(inactive)")
	for _, line := range all {
		if bytes.Equal(lbyt, line[len(line)-len(lbyt):]) {
			// Skip inactive branches
			continue
		}

		// Split on colon; this gets us the rev and the branch plus local revno
		pair := bytes.Split(line, []byte(":"))
		idx := bytes.IndexByte(pair[0], 32) // space
		versions = append(versions, VersionInfo{
			IsBranch: true,
			Name:     string(pair[0][:idx]),
			Revision: string(pair[1]),
		})
	}

	return
}

// IsReference returns if a string is a reference. A reference can be a
// commit id, branch, or tag.
func (s *HgRepo) IsReference(r string) bool {
	_, err := s.runFromDir("hg", "log", "-r", r)
	if err == nil {
		return true
	}

	return false
}

// IsDirty returns if the checkout has been modified from the checked
// out reference.
func (s *HgRepo) IsDirty() bool {
	out, err := s.runFromDir("hg", "diff")
	return err != nil || len(out) != 0
}

// CommitInfo retrieves metadata about a commit.
func (s *HgRepo) CommitInfo(id string) (*CommitInfo, error) {
	out, err := s.runFromDir("hg", "log", "-r", id, "--style=xml")
	if err != nil {
		return nil, ErrRevisionUnavailable
	}

	type Author struct {
		Name  string `xml:",chardata"`
		Email string `xml:"email,attr"`
	}
	type Logentry struct {
		Node   string `xml:"node,attr"`
		Author Author `xml:"author"`
		Date   string `xml:"date"`
		Msg    string `xml:"msg"`
	}
	type Log struct {
		XMLName xml.Name   `xml:"log"`
		Logs    []Logentry `xml:"logentry"`
	}

	logs := &Log{}
	err = xml.Unmarshal(out, &logs)
	if err != nil {
		return nil, err
	}
	if len(logs.Logs) == 0 {
		return nil, ErrRevisionUnavailable
	}

	ci := &CommitInfo{
		Commit:  logs.Logs[0].Node,
		Author:  logs.Logs[0].Author.Name + " <" + logs.Logs[0].Author.Email + ">",
		Message: logs.Logs[0].Msg,
	}

	if logs.Logs[0].Date != "" {
		ci.Date, err = time.Parse(time.RFC3339, logs.Logs[0].Date)
		if err != nil {
			return nil, err
		}
	}

	return ci, nil
}
