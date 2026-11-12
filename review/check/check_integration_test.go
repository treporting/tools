//go:build integration

package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/review"
)

// gitIn runs git in dir and answers its output.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Mathias Fischer", "GIT_AUTHOR_EMAIL=m@example.org",
		"GIT_COMMITTER_NAME=Mathias Fischer", "GIT_COMMITTER_EMAIL=m@example.org")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
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

// The check on a real history: the requirements are written by Claude Code as author, then another
// requirement of the same file is changed by the reviewer alone - which must not make the reviewer
// an author of the first.
//
// @implements TC_15@1.1.1
func TestTheCheckOnAHistoryOfGit(t *testing.T) {

	root := t.TempDir()
	gitIn(t, root, "init", "-q")
	write(t, root, "reviews/criteria-tool.yaml", criteria)
	write(t, root, "requirements/tool-requirements.yaml", requirements+`- object_key: TOOL_04
  object_text: first
`)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "--author=Claude Opus 5.5 <noreply@anthropic.com>", "-m", "requirements")
	first := gitIn(t, root, "rev-parse", "--short", "HEAD")

	write(t, root, "requirements/tool-requirements.yaml", requirements+`- object_key: TOOL_04
  object_text: changed by the reviewer alone
`)
	gitIn(t, root, "commit", "-q", "-am", "TOOL_04 changed")
	second := gitIn(t, root, "rev-parse", "--short", "HEAD")

	write(t, root, "reviews/r.yaml", `level: review
created_by: Mathias Fischer
created_at: '2026-10-25T13:30:00+01:00'
criteria: reviews/criteria-tool.yaml
criteria_revision: '`+first+`'
results:
- object_key: RR_01
  target_ids: [TOOL_03]
  revision: '`+second+`'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
- object_key: RR_02
  target_ids: [TOOL_04]
  revision: '`+second+`'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`)

	findings, err := checkAt(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Result != "RR_02" ||
		!strings.Contains(findings[0].Message, "is an author of TOOL_04") {
		t.Errorf("want one finding, that the reviewer is an author of TOOL_04; got %v", findings)
	}
}

// The criteria revision on a real history: a short hash of the commit that changed the criteria is
// taken, a later commit that did not change them is not.
//
// @verifies SR_07@1.1.1
func TestTheCriteriaRevisionOnAHistoryOfGit(t *testing.T) {

	root := t.TempDir()
	gitIn(t, root, "init", "-q")
	write(t, root, "reviews/criteria-tool.yaml", criteria)
	write(t, root, "requirements/tool-requirements.yaml", requirements)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "--author=Claude Opus 5.5 <noreply@anthropic.com>", "-m", "criteria")
	first := gitIn(t, root, "rev-parse", "--short", "HEAD")

	write(t, root, "requirements/tool-requirements.yaml", requirements+"# changed\n")
	gitIn(t, root, "commit", "-q", "--author=Claude Opus 5.5 <noreply@anthropic.com>", "-am", "no criteria")
	second := gitIn(t, root, "rev-parse", "--short", "HEAD")

	for _, c := range []struct {
		revision string
		want     int
	}{{first, 0}, {second, 1}} {
		write(t, root, "reviews/r.yaml", `level: review
created_by: Mathias Fischer
created_at: '2026-10-25T13:30:00+01:00'
criteria: reviews/criteria-tool.yaml
criteria_revision: '`+c.revision+`'
results:
- object_key: RR_01
  target_ids: [TOOL_03]
  revision: '`+second+`'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`)
		findings, err := checkAt(t, root)
		if err != nil {
			t.Fatal(err)
		}
		if len(findings) != c.want {
			t.Errorf("criteria revision %s: want %d findings, got %v", c.revision, c.want, findings)
		}
	}
}

// Two pull requests change one document from the same main; the reviewer merges both, so the second
// merge equals neither parent. Claude Code wrote every line, and the reviewer is no author.
//
// @implements TC_22@1.1.1
func TestTwoMergesOfOneDocument(t *testing.T) {

	root := t.TempDir()
	claude := "--author=Claude Opus 5.5 <noreply@anthropic.com>"
	gitIn(t, root, "init", "-q", "-b", "main")
	write(t, root, "reviews/criteria-tool.yaml", criteria)
	write(t, root, "doc.md", "top\n\nmiddle\n\nbottom\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", claude, "-m", "doc")
	first := gitIn(t, root, "rev-parse", "--short", "HEAD")

	gitIn(t, root, "switch", "-q", "-c", "a")
	write(t, root, "doc.md", "top, changed\n\nmiddle\n\nbottom\n")
	gitIn(t, root, "commit", "-q", claude, "-am", "a")
	gitIn(t, root, "switch", "-q", "main")
	gitIn(t, root, "switch", "-q", "-c", "b")
	write(t, root, "doc.md", "top\n\nmiddle\n\nbottom, changed\n")
	gitIn(t, root, "commit", "-q", claude, "-am", "b")
	gitIn(t, root, "switch", "-q", "main")
	gitIn(t, root, "merge", "-q", "--no-ff", "-m", "merge a", "a")
	gitIn(t, root, "merge", "-q", "--no-ff", "-m", "merge b", "b")
	merged := gitIn(t, root, "rev-parse", "--short", "HEAD")

	write(t, root, "reviews/r.yaml", `level: review
created_by: Mathias Fischer
created_at: '2026-10-25T13:30:00+01:00'
criteria: reviews/criteria-tool.yaml
criteria_revision: '`+first+`'
results:
- object_key: RR_01
  target_ids: [doc.md]
  revision: '`+merged+`'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`)
	findings, err := checkAt(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("want no finding, the merges are no authorship; got %v", findings)
	}
}

// checkAt checks reviews/r.yaml on the history of git at root, with the index of its working tree.
func checkAt(t *testing.T, root string) ([]Finding, error) {
	t.Helper()
	index, err := review.IndexAt(root)
	if err != nil {
		t.Fatal(err)
	}
	return Check(root, "reviews/r.yaml", review.ExecGit{Root: root}, index)
}

// resultsRepo answers a repo of git with two review results: one committed at the base, which leaves
// a failed criterion without a note, and one added after it, which breaks no rule; and the base.
func resultsRepo(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	gitIn(t, root, "init", "-q")
	write(t, root, "reviews/criteria-tool.yaml", criteria)
	write(t, root, "requirements/tool-requirements.yaml", requirements)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "--author=Claude Opus 5.5 <noreply@anthropic.com>", "-m", "requirements")
	first := gitIn(t, root, "rev-parse", "--short", "HEAD")
	result := func(criterion, verdict string) string {
		return `level: review
created_by: Mathias Fischer
created_at: '2026-10-25T13:30:00+01:00'
criteria: reviews/criteria-tool.yaml
criteria_revision: '` + first + `'
results:
- object_key: RR_01
  target_ids: [TOOL_03]
  revision: '` + first + `'
  criteria: {RC_01: passed, RC_02: ` + criterion + `}
  verdict: ` + verdict + `
`
	}
	write(t, root, "reviews/results/1_old.yaml", result("failed", "failed"))
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "the old result")
	base := gitIn(t, root, "rev-parse", "HEAD")
	write(t, root, "reviews/results/2_new.yaml", result("passed", "passed"))
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "the new result")
	return root, base
}

// checkAll checks each of the files, and answers what it found.
func checkAll(t *testing.T, root string, files []string) []Finding {
	t.Helper()
	index, err := review.IndexAt(root)
	if err != nil {
		t.Fatal(err)
	}
	var all []Finding
	for _, f := range files {
		findings, err := Check(root, f, review.ExecGit{Root: root}, index)
		if err != nil {
			t.Fatal(err)
		}
		all = append(all, findings...)
	}
	return all
}

// @implements TC_152@1.1.1
func TestWithBaseOnlyTheResultsAdded(t *testing.T) {
	root, base := resultsRepo(t)
	files, err := Results(root, base, review.ExecGit{Root: root}.Changed)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != "reviews/results/2_new.yaml" {
		t.Fatalf("checked %v", files)
	}
	if findings := checkAll(t, root, files); len(findings) != 0 {
		t.Errorf("findings %v", findings)
	}
}

// @implements TC_153@1.1.1
func TestWithoutBaseEveryResult(t *testing.T) {
	root, _ := resultsRepo(t)
	files, err := Results(root, "", review.ExecGit{Root: root}.Changed)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("checked %v", files)
	}
	findings := checkAll(t, root, files)
	if len(findings) != 1 || findings[0].File != "reviews/results/1_old.yaml" || !strings.Contains(findings[0].Message, "note") {
		t.Errorf("findings %v", findings)
	}
}
