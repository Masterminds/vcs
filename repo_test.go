package vcs

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func ExampleNewRepo() {
	remote := "https://github.com/Masterminds/vcs"
	local, _ := ioutil.TempDir("", "go-vcs")
	repo, _ := NewRepo(remote, local)
	// Returns: instance of GitRepo

	repo.Vcs()
	// Returns Git as this is a Git repo

	err := repo.Get()
	// Pulls down a repo, or a checkout in the case of SVN, and returns an
	// error if that didn't happen successfully.
	if err != nil {
		fmt.Println(err)
	}

	err = repo.UpdateVersion("master")
	// Checkouts out a specific version. In most cases this can be a commit id,
	// branch, or tag.
	if err != nil {
		fmt.Println(err)
	}
}

func TestDepInstalled(t *testing.T) {
	i := depInstalled("git")
	if !i {
		t.Error("depInstalled not finding installed dep.")
	}

	i = depInstalled("thisreallyisntinstalled")
	if i {
		t.Error("depInstalled finding not installed dep.")
	}
}

func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t}, "test", log.LstdFlags)
}

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}
