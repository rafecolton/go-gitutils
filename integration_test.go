package git

import (
	"fmt"
	"os"
	"os/exec"
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

func setupTestDir() error {
	exec.Command("git", "clone", testRepo, testDirPath)
	_, err := os.Stat(testDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := exec.Command("git", "clone", testRepo, testDirPath).Run(); err != nil {
				return err
			}
			cmd := exec.Command("git", "checkout", "-qf", testSha)
			cmd.Dir = testDirPath
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func TestIntegration(t *testing.T) {
	if !integration {
		t.Skip()
	}
	runner = nil

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
