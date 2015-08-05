package vcs

import (
	"testing"
)

func TestVCSLookup(t *testing.T) {
	// TODO: Expand to make sure it detected the right vcs.
	urlList := map[string]struct {
		work bool
		t    VcsType
	}{
		"https://github.com/masterminds":                                   {work: false, t: GitType},
		"https://github.com/Masterminds/VCSTestRepo":                       {work: true, t: GitType},
		"https://bitbucket.org/mattfarina/testhgrepo":                      {work: true, t: HgType},
		"https://launchpad.net/govcstestbzrrepo/trunk":                     {work: true, t: BzrType},
		"https://launchpad.net/~mattfarina/+junk/mygovcstestbzrrepo":       {work: true, t: BzrType},
		"https://launchpad.net/~mattfarina/+junk/mygovcstestbzrrepo/trunk": {work: true, t: BzrType},
		"https://git.launchpad.net/govcstestgitrepo":                       {work: true, t: GitType},
		"https://git.launchpad.net/~mattfarina/+git/mygovcstestgitrepo":    {work: true, t: GitType},
		"http://farbtastic.googlecode.com/svn/":                            {work: true, t: SvnType},
		"http://farbtastic.googlecode.com/svn/trunk":                       {work: true, t: SvnType},
		"https://code.google.com/p/farbtastic":                             {work: false, t: SvnType},
		"https://code.google.com/p/plotinum":                               {work: true, t: HgType},
		"https://example.com/foo/bar.git":                                  {work: true, t: GitType},
		"https://example.com/foo/bar.svn":                                  {work: true, t: SvnType},
		"https://example.com/foo/bar/baz.bzr":                              {work: true, t: BzrType},
		"https://example.com/foo/bar/baz.hg":                               {work: true, t: HgType},
	}

	for u, c := range urlList {
		ty, err := detectVcsFromUrl(u)
		if err == nil && c.work == false {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err == ErrCannotDetectVCS && c.work == true {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err != nil && c.work == true {
			t.Errorf("Error detecting VCS from URL(%s): %s", u, err)
		}

		if c.work == true && ty != c.t {
			t.Errorf("Incorrect VCS type returned(%s)", u)
		}
	}
}
