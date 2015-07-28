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

	repo := NewSvnRepo("https://github.com/Masterminds/cookoo/trunk", tempDir+"/cookoo")

	// Check the basic getters.
	if repo.Remote() != "https://github.com/Masterminds/cookoo/trunk" {
		t.Error("Remote not set properly")
	}
	if repo.LocalPath() != tempDir+"/cookoo" {
		t.Error("Local disk location not set properly")
	}

	//Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Do an initial checkout.
	err = repo.Get()
	if err != nil {
		t.Errorf("Unable to checkout SVN repo. Err was %s", err)
	}

	// Update the version to a previous version.
	err = repo.UpdateVersion("r100")
	if err != nil {
		t.Errorf("Unable to update SVN repo version. Err was %s", err)
	}

	// Use Version to verify we are on the right version.
	v, err := repo.Version()
	if v != "100" {
		t.Error("Error checking checked out version")
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
	if v == "100" {
		t.Error("Error with version. Still on old version. Update failed")
	}
	if err != nil {
		t.Error(err)
	}

}
