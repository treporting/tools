//go:build integration

package skeleton

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/treporting/specs-and-deploy/tools/review"
	"github.com/treporting/specs-and-deploy/tools/review/check"
)

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Claude Opus 5.5", "GIT_AUTHOR_EMAIL=noreply@anthropic.com",
		"GIT_COMMITTER_NAME=Mathias Fischer", "GIT_COMMITTER_EMAIL=m@example.org")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func write(t *testing.T, root, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, filepath.Dir(path)), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// A skeleton on a real history: it is made, and reviewcheck then reports every todo in it - an
// unfilled skeleton cannot pass.
//
// @verifies SR_11@2.1.1, SR_12@1.1.1, SR_13@1.1.1, SR_14@1.1.1
func TestASkeletonOnAHistoryOfGit(t *testing.T) {

	root := t.TempDir()
	gitIn(t, root, "init", "-q")
	write(t, root, "reviews/criteria-tool.yaml", "level: tool\ncriteria:\n- object_key: RC_01\n- object_key: RC_02\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "criteria")
	write(t, root, "requirements/tool-requirements.yaml", "type: requirements\nlevel: tool\nrequirements:\n- object_key: TOOL_03\n  object_text: check\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "requirement")

	git := review.ExecGit{Root: root}
	index, err := review.IndexAt(root)
	if err != nil {
		t.Fatal(err)
	}
	path, err := Make(root, Request{Reviewer: "Mathias Fischer", Scope: "tool-03", Targets: []string{"TOOL_03"}, Now: time.Now()}, git, index)
	if err != nil {
		t.Fatal(err)
	}

	findings, err := check.Check(root, path, git, index)
	if err != nil {
		t.Fatal(err)
	}
	todo := 0
	for _, f := range findings {
		if strings.Contains(f.Message, "neither passed nor failed") {
			todo++
		} else {
			t.Errorf("want only the todos reported, got %v", f)
		}
	}
	if todo != 3 {
		t.Errorf("want 3 todos reported, two criteria and the verdict; got %d: %v", todo, findings)
	}
}
