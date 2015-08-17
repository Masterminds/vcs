package vcs

import (
	"io/ioutil"
)

func ExampleNewRepo() {
	remote := "https://github.com/Masterminds/vcs"
	local, _ := ioutil.TempDir("", "go-vcs")
	repo, _ := NewRepo(remote, local)
	// Returns: instance of GitRepo

	repo.Vcs()
	// Returns Git as this is a Git repo
}
