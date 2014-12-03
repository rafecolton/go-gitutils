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
	return runCommands(exec.Command("git", "checkout", "-qf", testSha))
}

func runCommands(commands ...*exec.Cmd) error {
	for _, cmd := range commands {
		cmd.Dir = testDirPath
		if err := cmd.Run(); err != nil {
			return err
		}
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
	return runCommands(exec.Command("git", "checkout", "master"))
}

func commit() error {
	if err := ioutil.WriteFile(testDirPath+"/README.md", []byte("foo"), 0644); err != nil {
		return err
	}
	commands := []*exec.Cmd{
		exec.Command("git", "config", "--global", "user.email", "foo@example.com"),
		exec.Command("git", "config", "--global", "user.name", "Foo Example"),
		exec.Command("git", "add", "README.md"),
		exec.Command("git", "commit", "-m", "foo"),
	}

	return runCommands(commands...)
}

func resetCommit() error {
	return runCommands(exec.Command("git", "reset", "--hard", "HEAD^"))
}
