package vcs

import (
	"io/ioutil"
	//"log"
	"os"
	"testing"
)

// Canary test to ensure BzrRepo implements the Repo interface.
var _ Repo = &BzrRepo{}

// To verify bzr is working we perform intergration testing
// with a known bzr service.

func TestBzr(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "go-vcs-bzr-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewBzrRepo("https://launchpad.net/govcstestbzrrepo", tempDir+"/govcstestbzrrepo")

	// Check the basic getters.
	if repo.Remote() != "https://launchpad.net/govcstestbzrrepo" {
		t.Error("Remote not set properly")
	}
	if repo.LocalPath() != tempDir+"/govcstestbzrrepo" {
		t.Error("Local disk location not set properly")
	}

	//Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Do an initial clone.
	err = repo.Get()
	if err != nil {
		t.Errorf("Unable to clone Bzr repo. Err was %s", err)
	}

	// Verify Bzr repo is a Bzr repo
	if repo.CheckLocal() == false {
		t.Error("Problem checking out repo or Bzr CheckLocal is not working")
	}

	// Test internal lookup mechanism used outside of Bzr specific functionality.
	ltype, err := detectVcsFromFS(tempDir + "/govcstestbzrrepo")
	if err != nil {
		t.Error("detectVcsFromFS unable to Bzr repo")
	}
	if ltype != BzrType {
		t.Errorf("detectVcsFromFS detected %s instead of Bzr type", ltype)
	}

	err = repo.UpdateVersion("2")
	if err != nil {
		t.Errorf("Unable to update Bzr repo version. Err was %s", err)
	}

	// Use Version to verify we are on the right version.
	v, err := repo.Version()
	if v != "2" {
		t.Error("Error checking checked out Bzr version")
	}
	if err != nil {
		t.Error(err)
	}

	// Perform an update.
	err = repo.Update()
	if err != nil {
		t.Error(err)
	}

	v, err = repo.Version()
	if v != "3" {
		t.Error("Error checking checked out Bzr version")
	}
	if err != nil {
		t.Error(err)
	}

}

func TestBzrCheckLocal(t *testing.T) {
	// Verify repo.CheckLocal fails for non-Bzr directories.
	// TestBzr is already checking on a valid repo
	tempDir, err := ioutil.TempDir("", "go-vcs-bzr-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewBzrRepo("", tempDir)
	if repo.CheckLocal() == true {
		t.Error("Bzr CheckLocal does not identify non-Bzr location")
	}
}
