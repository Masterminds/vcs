package vcs

import (
	"io/ioutil"
	//"log"
	"os"
	"testing"
)

// Canary test to ensure HgRepo implements the Repo interface.
var _ Repo = &HgRepo{}

// To verify hg is working we perform intergration testing
// with a known hg service.

func TestHg(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "go-vcs-hg-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewHgRepo("https://bitbucket.org/mattfarina/testhgrepo", tempDir+"/testhgrepo")

	// Check the basic getters.
	if repo.Remote() != "https://bitbucket.org/mattfarina/testhgrepo" {
		t.Error("Remote not set properly")
	}
	if repo.LocalPath() != tempDir+"/testhgrepo" {
		t.Error("Local disk location not set properly")
	}

	//Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Do an initial clone.
	err = repo.Get()
	if err != nil {
		t.Errorf("Unable to clone Hg repo. Err was %s", err)
	}

	// Verify Hg repo is a Hg repo
	if repo.CheckLocal() == false {
		t.Error("Problem checking out repo or Hg CheckLocal is not working")
	}

	// Test internal lookup mechanism used outside of Hg specific functionality.
	ltype, err := detectVcsFromFS(tempDir + "/testhgrepo")
	if err != nil {
		t.Error("detectVcsFromFS unable to Hg repo")
	}
	if ltype != HgType {
		t.Errorf("detectVcsFromFS detected %s instead of Hg type", ltype)
	}

	// Set the version using the short hash.
	err = repo.UpdateVersion("a5494ba2177f")
	if err != nil {
		t.Errorf("Unable to update Hg repo version. Err was %s", err)
	}

	// Use Version to verify we are on the right version.
	v, err := repo.Version()
	if v != "a5494ba2177f" {
		t.Error("Error checking checked out Hg version")
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
	if v != "d680e82228d2" {
		t.Error("Error checking checked out Hg version")
	}
	if err != nil {
		t.Error(err)
	}

}

func TestHgCheckLocal(t *testing.T) {
	// Verify repo.CheckLocal fails for non-Hg directories.
	// TestHg is already checking on a valid repo
	tempDir, err := ioutil.TempDir("", "go-vcs-hg-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo := NewHgRepo("", tempDir)
	if repo.CheckLocal() == true {
		t.Error("Hg CheckLocal does not identify non-Hg location")
	}
}
