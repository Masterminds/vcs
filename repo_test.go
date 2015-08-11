package vcs

import (
	"io/ioutil"
)

func ExampleNewRepo() {
	remote := "https://github.com/Masterminds/go-vcs"
	local, _ := ioutil.TempDir("", "go-vcs")
	repo, _ := NewRepo(remote, local)
	// Returns: instance of GitRepo

	repo.Vcs()
	// Returns GitType as this is a Git repo
}
