package vcs

import (
	"os"
)

// Detect the type from the local path. Is there a better way to do this?
func detectVcsFromFS(vcsPath string) (VcsType, error) {

	// When the local directory to the package doesn't exist
	// it's not yet downloaded so we can't detect the type
	// locally.
	if _, err := os.Stat(vcsPath); os.IsNotExist(err) {
		return "", ErrCannotDetectVCS
	}

	seperator := string(os.PathSeparator)

	// Walk through each of the different VCS types to see if
	// one can be detected. Do this is order of guessed popularity.
	if _, err := os.Stat(vcsPath + seperator + ".git"); err == nil {
		return GitType, nil
	}
	if _, err := os.Stat(vcsPath + seperator + ".svn"); err == nil {
		return SvnType, nil
	}
	if _, err := os.Stat(vcsPath + seperator + ".hg"); err == nil {
		return HgType, nil
	}
	if _, err := os.Stat(vcsPath + seperator + ".bzr"); err == nil {
		return BzrType, nil
	}

	// If one was not already detected than we default to not finding it.
	return "", ErrCannotDetectVCS

}
