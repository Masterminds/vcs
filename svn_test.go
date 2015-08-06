package vcs

import (
	"io/ioutil"
	//"log"
	"os"
	"testing"
)

// To verify svn is working we perform intergration testing
// with a known svn service.

// Canary test to ensure SvnRepo implements the Repo interface.
var _ Repo = &SvnRepo{}

func TestSvn(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "go-vcs-svn-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewSvnRepo("https://github.com/Masterminds/VCSTestRepo/trunk", tempDir+"/VCSTestRepo")

	// Check the basic getters.
	if repo.Remote() != "https://github.com/Masterminds/VCSTestRepo/trunk" {
		t.Error("Remote not set properly")
	}
	if repo.LocalPath() != tempDir+"/VCSTestRepo" {
		t.Error("Local disk location not set properly")
	}

	//Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Do an initial checkout.
	err = repo.Get()
	if err != nil {
		t.Errorf("Unable to checkout SVN repo. Err was %s", err)
	}

	// Verify SVN repo is a SVN repo
	if repo.CheckLocal() == false {
		t.Error("Problem checking out repo or SVN CheckLocal is not working")
	}

	// Test internal lookup mechanism used outside of Hg specific functionality.
	ltype, err := detectVcsFromFS(tempDir + "/VCSTestRepo")
	if err != nil {
		t.Error("detectVcsFromFS unable to Svn repo")
	}
	if ltype != SvnType {
		t.Errorf("detectVcsFromFS detected %s instead of Svn type", ltype)
	}

	// Update the version to a previous version.
	err = repo.UpdateVersion("r2")
	if err != nil {
		t.Errorf("Unable to update SVN repo version. Err was %s", err)
	}

	// Use Version to verify we are on the right version.
	v, err := repo.Version()
	if v != "2" {
		t.Error("Error checking checked SVN out version")
	}
	if err != nil {
		t.Error(err)
	}

	// Perform an update which should take up back to the latest version.
	err = repo.Update()
	if err != nil {
		t.Error(err)
	}

	// Make sure we are on a newer version because of the update.
	v, err = repo.Version()
	if v == "2" {
		t.Error("Error with version. Still on old version. Update failed")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestSvnCheckLocal(t *testing.T) {
	// Verify repo.CheckLocal fails for non-SVN directories.
	// TestSvn is already checking on a valid repo
	tempDir, err := ioutil.TempDir("", "go-vcs-svn-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewSvnRepo("", tempDir)
	if repo.CheckLocal() == true {
		t.Error("SVN CheckLocal does not identify non-SVN location")
	}
}
