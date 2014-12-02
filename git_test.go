package git

import (
	"testing"
)

func init() {
	runner = &fakeRunner{}
}

var (
	testingPath = "_testing/kamino-test"
	testingSha  = "97d2258b4a58d9bf07636d76c97a3eb09490cf70"
)

func TestRemoteAccountHTTP(t *testing.T) {
	runner.(*fakeRunner).remoteV = `origin  http://github.com/rafecolton/docker-builder.git (fetch)`
	expected := "rafecolton"
	actual := RemoteAccount("")
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestRemoteAccountHTTPS(t *testing.T) {
	runner.(*fakeRunner).remoteV = `origin  https://github.com/rafecolton/docker-builder.git (fetch)`
	expected := "rafecolton"
	actual := RemoteAccount("")
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestRemoteAccountSSH(t *testing.T) {
	runner.(*fakeRunner).remoteV = `origin  git@github.com:rafecolton/docker-builder.git (fetch)`
	expected := "rafecolton"
	actual := RemoteAccount("")
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestRemoteAccountGit(t *testing.T) {
	runner.(*fakeRunner).remoteV = `origin  git://github.com/rafecolton/docker-builder.git (fetch)`
	expected := "rafecolton"
	actual := RemoteAccount("")
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestRemoteAccountWithoutSuffix(t *testing.T) {
	runner.(*fakeRunner).remoteV = `origin  git://github.com/rafecolton/docker-builder (fetch)`
	expected := "rafecolton"
	actual := RemoteAccount("")
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestSha(t *testing.T) {
	runner.(*fakeRunner).sha = "abc123\n"
	actual := Sha("")
	expected := "abc123"
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestBranch(t *testing.T) {
	runner.(*fakeRunner).branch = "asdf\n"
	actual := Branch("")
	expected := "asdf"
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestBranchAlt(t *testing.T) {
	runner.(*fakeRunner).branch = "HEAD\n"
	runner.(*fakeRunner).branch2 = `  master
  move-mithril-to-quay
* update-loadbalancer-role
  using-mighril-from-quay-instead-of-docker-hub`
	actual := Branch("")
	expected := "master"
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestCleanClean(t *testing.T) {
	if !IsClean("") {
		t.Errorf("expected cleanliness")
	}
}

func TestCleanDirty(t *testing.T) {
	runner.(*fakeRunner).clean = " 2 files changed, 139 insertions(+), 105 deletions(-)\n"
	if IsClean("") {
		t.Errorf("repo is dirty")
	}
}

func TestUpToDateDiverged(t *testing.T) {
	runner.(*fakeRunner).upToDateLocal = "abc"
	runner.(*fakeRunner).upToDateRemote = "123"
	runner.(*fakeRunner).upToDateBase = "def"
	if UpToDate("") != StatusDiverged {
		t.Errorf("status should be StatusDiverged")
	}
}

func TestUpToDateUpToDate(t *testing.T) {
	runner.(*fakeRunner).upToDateLocal = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	runner.(*fakeRunner).upToDateRemote = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	if UpToDate("") != StatusUpToDate {
		t.Errorf("status should be StatusUpToDate")
	}
}

func TestUpToDateNeedToPush(t *testing.T) {
	runner.(*fakeRunner).upToDateLocal = "f4c103a85141c59749ef24320a538ae7ed238909"
	runner.(*fakeRunner).upToDateRemote = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	runner.(*fakeRunner).upToDateBase = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	if UpToDate("") != StatusNeedToPush {
		t.Errorf("status should be StatusNeedToPush")
	}
}

func TestUpToDateNeedToPull(t *testing.T) {
	runner.(*fakeRunner).upToDateLocal = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	runner.(*fakeRunner).upToDateRemote = "f4c103a85141c59749ef24320a538ae7ed238909"
	runner.(*fakeRunner).upToDateBase = "50d1ab234ffa3df05162c8eae4dddef1d907faa8"
	if UpToDate("") != StatusNeedToPull {
		t.Errorf("status should be StatusNeedToPull")
	}
}

func TestTag(t *testing.T) {
	runner.(*fakeRunner).tag = "foo-tag"
	if Tag("") != "foo-tag" {
		t.Errorf("expected foo-tag, got %s", Tag(""))
	}
}

func TestStatusString(t *testing.T) {
	assertStatusString(StatusUpToDate, "StatusUpToDate", t)
	assertStatusString(StatusNeedToPull, "StatusNeedToPull", t)
	assertStatusString(StatusNeedToPush, "StatusNeedToPush", t)
	assertStatusString(StatusDiverged, "StatusDiverged", t)
}

func assertStatusString(status Status, str string, t *testing.T) {
	if status.String() != str {
		t.Errorf("expected " + str + ", got " + status.String())
	}
}

func TestPrintInvalidStatus(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			println("wheeeeee")
		}
	}()
	var invalid uint8 = 86
	assertStatusString(Status(invalid), "doesn't matter", t)
	t.Error("should not get here")
}
