package vcs

import (
	"testing"
)

func TestVCSLookup(t *testing.T) {
	// TODO: Expand to make sure it detected the right vcs.
	urlList := map[string]bool{
		"https://github.com/masterminds":              false,
		"https://github.com/Masterminds/VCSTestRepo":  true,
		"https://bitbucket.org/mattfarina/testhgrepo": true,
	}

	for u, c := range urlList {
		_, err := detectVcsFromUrl(u)
		if err == nil && c == false {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err == ErrCannotDetectVCS && c == true {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err != nil && c == true {
			t.Errorf("Error detecting VCS from URL(%s): %s", u, err)
		}
	}
}
