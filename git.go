package git

import (
	"os/exec"
	"regexp"
	"strings"
)

// GitRemoteRegex is a regex for pulling account owner from the output of `git remote -v`
var GitRemoteRegex = regexp.MustCompile(`^([^\t\n\f\r ]+)[\t\n\v\f\r ]+(git@github\.com:|(http[s]?|git):\/\/github\.com\/)([a-zA-Z0-9]{1}[a-zA-Z0-9-]*)\/([a-zA-Z0-9_.-]+)(\.git|[^\t\n\f\r ])+.*$`)

/*
Branch determines the git branch in the repo located at `top`.  Two attempts
are made to determine branch. First:

  git rev-parse -q --abbrev-ref HEAD

If the current working directory is not on a branch, the result will return
"HEAD". In this case, branch will be estimated by parsing the output of the
following:

  git branch --contains $(git rev-parse -q HEAD)
*/
func Branch(top string) string {
	var branchCmd = exec.Command("git", "rev-parse", "-q", "--abbrev-ref", "HEAD")
	branchCmd.Dir = top
	branchBytes, err := branchCmd.Output()
	if err != nil {
		return ""
	}
	branch := string(branchBytes)[:len(branchBytes)-1]
	if branch == "HEAD" {
		branchCmd = exec.Command("git", "branch", "--contains", Sha(top))
		branchCmd.Dir = top
		branchBytes, err := branchCmd.Output()
		if err != nil {
			return branch
		}
		branches := strings.Split(string(branchBytes)[:len(branchBytes)-1], "\n")
		for _, branchStr := range branches {
			if len(branchStr) < 1 || string(branchStr[0]) == "*" {
				continue
			}
			return strings.Trim(branchStr, " ")
		}
	}
	return branch
}

// Sha produces the git branch at `top` as determined by `git rev-parse -q HEAD`
func Sha(top string) string {
	var revCmd = exec.Command("git", "rev-parse", "-q", "HEAD")
	revCmd.Dir = top
	revBytes, err := revCmd.Output()
	if err != nil {
		return ""
	}
	rev := string(revBytes)[:len(revBytes)-1]
	return rev
}

// Tag produces the git tag at `top` as determined by `git describe --always --dirty --tags`
func Tag(top string) string {
	var shortCmd = exec.Command("git", "describe", "--always", "--dirty", "--tags")
	shortCmd.Dir = top
	shortBytes, err := shortCmd.Output()
	if err != nil {
		return ""
	}
	short := string(shortBytes)[:len(shortBytes)-1]
	return short
}

// IsClean returns true `git diff --shortstat` produces no output
func IsClean(top string) bool {
	var cmd = exec.Command("git", "diff", "--shortstat")
	cmd.Dir = top
	outBytes, err := cmd.Output()
	if err != nil {
		return false
	}
	if len(outBytes) > 0 {
		return false
	}
	return true
}

const (
	// StatusUpToDate means the local repo matches origin
	StatusUpToDate = iota

	// StatusNeedToPull means the local repo needs to pull from the remote
	StatusNeedToPull

	// StatusNeedToPush means the local repo needs to be pushed to the remote
	StatusNeedToPush

	// StatusDiverged means the local repo has diverged from the remote
	StatusDiverged
)

//UpToDate returns the status of the repo as determined by the above constants
func UpToDate(top string) int {
	var cmdLocal = exec.Command("git", "rev-parse", "@")
	cmdLocal.Dir = top
	local, err := runCmd(cmdLocal)
	if err != nil {
		return StatusDiverged
	}

	var cmdRemote = exec.Command("git", "rev-parse", "@{u}")
	cmdRemote.Dir = top
	remote, err := runCmd(cmdRemote)
	if err != nil {
		return StatusDiverged
	}

	if local == remote {
		return StatusUpToDate
	}

	var cmdBase = exec.Command("git", "merge-base", "@", "@{u}")
	cmdBase.Dir = top
	base, err := runCmd(cmdBase)
	if err != nil {
		return StatusDiverged
	}
	if local == base {
		return StatusNeedToPull
	} else if remote == base {
		return StatusNeedToPush
	}
	return StatusDiverged
}

/*
RemoteAccount returns the github account as determined by the output of `git
remote -v`
*/
func RemoteAccount(top string) string {
	return AccountFromRemotes(Remotes(top))
}

/*
AccountFromRemotes produces the (GitHub) account from the output of `git remote
-v` - this ouput is nominally provided by the Remotes function but may be
provided otherwise for testing purposes
*/
func AccountFromRemotes(remotes string) string {
	lines := strings.Split(remotes, "\n")

	var ret string

	for _, line := range lines {
		matches := GitRemoteRegex.FindStringSubmatch(line)
		if len(matches) == 7 {
			if matches[1] == "origin" {
				return matches[4]
			}
		}
	}

	return ret
}

// Remotes returns the output of `git remote -v
func Remotes(top string) string {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = top
	outBytes, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(outBytes)
}

/* HELPERS */

func runCmd(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(outBytes)[:len(outBytes)-1], nil
}
