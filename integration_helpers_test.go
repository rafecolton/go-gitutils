package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func beforeEach(t *testing.T) error {
	if !integration {
		t.Skip()
	}
	runner = nil
	cmd := exec.Command("git", "checkout", "-qf", testSha)
	cmd.Dir = testDirPath
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func setupTestDir() error {
	exec.Command("git", "clone", testRepo, testDirPath)
	_, err := os.Stat(testDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := exec.Command("git", "clone", testRepo, testDirPath).Run(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func resetToMaster() error {
	cmd := exec.Command("git", "checkout", "master")
	cmd.Dir = testDirPath
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func commit() error {
	if err := ioutil.WriteFile(testDirPath+"/README.md", []byte("foo"), 0644); err != nil {
		return err
	}
	cmd := exec.Command("git", "add", "README.md")
	cmd.Dir = testDirPath
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", "foo")
	cmd.Dir = testDirPath
	return cmd.Run()
}

func resetCommit() error {
	cmd := exec.Command("git", "reset", "--hard", "HEAD^")
	cmd.Dir = testDirPath
	return cmd.Run()
}
