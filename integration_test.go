package git

import (
	"fmt"
	"os"
	"testing"
)

var (
	integration bool
	testDirPath = os.Getenv("PWD") + "/_testing/kamino-test"
	testSha     = "97d2258b4a58d9bf07636d76c97a3eb09490cf70"
	testRepo    = "https://github.com/rafecolton/kamino-test.git"
)

func init() {
	if os.Getenv("INTEGRATION") != "" {
		integration = true
		if err := setupTestDir(); err != nil {
			fmt.Printf("error setting up the test dir: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func TestIntegration(t *testing.T) {
	if err := beforeEach(t); err != nil {
		t.Fatal(err)
	}

	branch := Branch(testDirPath)
	clean := IsClean(testDirPath)
	remote := RemoteAccount(testDirPath)
	sha := Sha(testDirPath)
	tag := Tag(testDirPath)
	upToDate := UpToDate(testDirPath)

	if branch != "removing-submodule-DO-NOT-DELETE-THIS-BRANCH" {
		t.Errorf("expected removing-submodule-DO-NOT-DELETE-THIS-BRANCH, got %s", branch)
	}

	if !clean {
		t.Error("expected repo to be clean")
	}

	if remote != "rafecolton" {
		t.Errorf("expected rafecolton, got %s", remote)
	}

	if sha != testSha {
		t.Errorf("expected %s, got %s", testSha, sha)
	}

	if tag != sha[0:7] {
		t.Errorf("expected %s, got %s", sha[0:7], tag)
	}

	if upToDate != StatusDiverged {
		t.Errorf("expected status %s, got %s", StatusDiverged, upToDate)
	}
}

func TestBranchUpToDate(t *testing.T) {
	if err := beforeEach(t); err != nil {
		t.Fatal(err)
	}
	if err := resetToMaster(); err != nil {
		t.Fatal(err)
	}

	expected := "master"
	actual := Branch(testDirPath)
	if actual != expected {
		t.Errorf("expected branch %q, got %q", expected, actual)
	}
}

func TestNeedToPushStatus(t *testing.T) {
	if err := beforeEach(t); err != nil {
		t.Fatal(err)
	}
	if err := resetToMaster(); err != nil {
		t.Fatal(err)
	}
	if err := commit(); err != nil {
		t.Fatal(err)
	}
	status := UpToDate(testDirPath)
	if err := resetCommit(); err != nil {
		t.Fatal(err)
	}
	if status != StatusNeedToPush {
		t.Errorf("expected StatusNeedToPush, got " + status.String())
	}
}
