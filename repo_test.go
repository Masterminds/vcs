package vcs

import (
	"io/ioutil"
)

func ExampleNewRepo() {
	remote := "https://github.com/Masterminds/go-vcs"
	local, err := ioutil.TempDir("", "go-vcs")
	repo, err := NewRepo(remote, local)
	// Returns: instance of GitRepo
}
