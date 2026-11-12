//go:build integration

package check

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/review"
)

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// A requirement approved in the same pull request does not count: the base commit has it as draft.
//
// @implements TC_73@1.1.1
func TestApprovedInTheSamePullRequestDoesNotCount(t *testing.T) {

	dir := write(t, map[string]string{"tools/software-requirements.yaml": strings.Replace(requirements, "SR_01\n  status: approved", "SR_01\n  status: draft", 1)})
	gitIn(t, dir, "init", "-q")
	gitIn(t, dir, "-c", "user.name=x", "-c", "user.email=x@x", "add", "-A")
	gitIn(t, dir, "-c", "user.name=x", "-c", "user.email=x@x", "commit", "-q", "-m", "base")
	base := gitIn(t, dir, "rev-parse", "--short", "HEAD")

	write2 := map[string]string{
		"tools/software-requirements.yaml": requirements,
		"a.go":                             "// @implements SR_01\nfunc a() {}\n",
	}
	for path, content := range write2 {
		if err := writeFile(dir, path, content); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, dir, "-c", "user.name=x", "-c", "user.email=x@x", "add", "-A")
	gitIn(t, dir, "-c", "user.name=x", "-c", "user.email=x@x", "commit", "-q", "-m", "SR_01 approved, and its code")

	found, _, err := Find(dir, strings.Fields(gitIn(t, dir, "ls-files")), Kinds{Whitelist: []string{"*.go", "*.yaml"}})
	if err != nil {
		t.Fatal(err)
	}
	index, err := review.IndexAt(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := messages(Check(review.ExecGit{Root: dir}, base, index, found))
	if len(got) != 1 || !strings.Contains(got[0], "@implements SR_01: draft in the base commit") {
		t.Errorf("want SR_01 refused as draft in the base commit, got %q", got)
	}
}
