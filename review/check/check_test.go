package check

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/review"
)

// fakeGit is a history of one commit and its parent, kept in memory.
type fakeGit struct {
	files     map[string]map[string]string // commit -> path -> content
	parent    map[string]string
	changed   map[string][]string // path -> the commits that changed it, the newest first
	merges    map[string]bool
	authors   map[string]string // commit -> its author, where it is not author
	author    string
	coAuthors []string
}

func (f fakeGit) CommitExists(rev string) bool { _, ok := f.files[rev]; return ok }

func (f fakeGit) Show(rev, path string) ([]byte, error) {
	if c, ok := f.files[rev][path]; ok {
		return []byte(c), nil
	}
	return nil, os.ErrNotExist
}

func (f fakeGit) FileCommits(rev, path string) ([]string, error) {
	if commits, ok := f.changed[path]; ok {
		return commits, nil
	}
	return []string{rev}, nil
}

func (f fakeGit) Parent(commit string) (string, bool) {
	p, ok := f.parent[commit]
	return p, ok
}

func (f fakeGit) IsMerge(commit string) bool { return f.merges[commit] }

func (f fakeGit) Authors(commit string) (string, []string, error) {
	if a, ok := f.authors[commit]; ok {
		return a, nil, nil
	}
	return f.author, f.coAuthors, nil
}

const criteria = `level: tool
criteria:
- object_key: RC_01
- object_key: RC_02
`

const requirements = `type: requirements
level: tool
requirements:
- object_key: TOOL_03
  object_text: check every review result
`

func history(author string, coAuthors ...string) fakeGit {
	return fakeGit{
		files: map[string]map[string]string{
			"c1": {"reviews/criteria-tool.yaml": criteria, "requirements/tool-requirements.yaml": requirements},
		},
		parent:    map[string]string{},
		author:    author,
		coAuthors: coAuthors,
	}
}

// indexOf makes the index of the files the history holds in its newest commit.
func indexOf(git fakeGit) *review.Index {
	newest := map[string]string{}
	for _, files := range git.files {
		for file, content := range files {
			newest[file] = content
		}
	}
	var names []string
	for file := range newest {
		names = append(names, file)
	}
	sort.Strings(names)
	return review.NewIndex(func(file string) ([]byte, error) { return []byte(newest[file]), nil }, names)
}

// check writes a review result into a temporary root and checks it.
func check(t *testing.T, result string, git fakeGit) []string {
	t.Helper()
	return checkRecord(t, "status: open\n", result, git)
}

// checkRecord checks a review result with the fields of its record given.
func checkRecord(t *testing.T, record, result string, git fakeGit) []string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "reviews"), 0o750); err != nil {
		t.Fatal(err)
	}
	review := `level: review
created_by: Mathias Fischer
created_at: '2026-10-25T13:30:00+01:00'
` + record + `criteria: reviews/criteria-tool.yaml
criteria_revision: 'c1'
results:
` + result
	if err := os.WriteFile(filepath.Join(root, "reviews", "r.yaml"), []byte(review), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, err := Check(root, "reviews/r.yaml", git, indexOf(git))
	if err != nil {
		t.Fatal(err)
	}
	var messages []string
	for _, f := range findings {
		messages = append(messages, f.Message)
	}
	return messages
}

func expect(t *testing.T, messages []string, want string) {
	t.Helper()
	for _, m := range messages {
		if strings.Contains(m, want) {
			return
		}
	}
	t.Errorf("no finding %q in %q", want, messages)
}

// @implements TC_01@1.1.1, TC_06@1.1.1, TC_11@1.1.1, TC_18@1.1.1, TC_20@1.1.1, TC_37@1.1.1
func TestAGoodResultHasNoFinding(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	if len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_02@1.1.1, TC_03@1.1.1
func TestEveryCriterionOnce(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_09: passed}
  verdict: passed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "criterion RC_02 is missing")
	expect(t, got, "criterion RC_09 is not in")
}

// @implements TC_04@1.1.1, TC_05@1.1.1
func TestTheVerdictFollowsTheCriteria(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: {result: failed, note: unclear}}
  verdict: passed
- object_key: RR_02
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: failed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "the verdict is passed, but not every criterion passed")
	expect(t, got, "the verdict is failed, but every criterion passed")
}

// @implements TC_07@1.1.1, TC_09@1.1.1
func TestAFailedVerdictWithANote(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: {result: failed, note: unclear}}
  verdict: failed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	if len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_08@1.1.1
func TestAFailedCriterionHasANote(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: failed}
  verdict: failed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "criterion RC_02 failed with no note")
}

// @implements TC_10@1.1.1
func TestTheRevisionExists(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'nope'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "the revision nope is no commit")
}

// @implements TC_12@1.1.1, TC_13@1.1.1
func TestTheReviewerIsNotTheAuthor(t *testing.T) {
	result := `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`
	//	written by the reviewer alone
	expect(t, check(t, result, history("Mathias Fischer")), "the reviewer Mathias Fischer is an author of TOOL_03")

	//	written by another, with the reviewer as co-author
	expect(t, check(t, result, history("Copilot", "Mathias Fischer")), "the reviewer Mathias Fischer is an author of TOOL_03")

	//	written by another, with no co-author: whatever name stands there
	if got := check(t, result, history("Copilot")); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_12@1.1.1
func TestASoftwareRequirementOfTheToolsHasAnAuthor(t *testing.T) {
	git := history("Mathias Fischer")
	git.files["c1"]["tools/software-requirements.yaml"] = `type: requirements
level: r5
requirements:
- object_key: SR_01
  object_text: report a criterion that is missing
`
	got := check(t, `- object_key: RR_01
  target_ids: [SR_01]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`, git)
	expect(t, got, "the reviewer Mathias Fischer is an author of SR_01")
}

// @implements TC_16@1.1.1
func TestAValueOtherThanPassedOrFailed(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: todo}
  verdict: passed
- object_key: RR_02
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: pased
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "criterion RC_02 is neither passed nor failed")
	//	counted as not passed
	expect(t, got, "the verdict is passed, but not every criterion passed")
	expect(t, got, "the verdict is neither passed nor failed")
}

// @implements TC_19@1.1.1
func TestTheCriteriaRevisionChangedTheCriteria(t *testing.T) {
	git := history("Claude Opus 5.5", "Claude Opus 5.5")
	git.changed = map[string][]string{"reviews/criteria-tool.yaml": {"c0"}}
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
- object_key: RR_02
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`, git)
	if n := strings.Count(strings.Join(got, "\n"), "the criteria revision c1 did not change reviews/criteria-tool.yaml"); n != 1 {
		t.Errorf("want the criteria revision reported once, got %d in %q", n, got)
	}
}

// @implements TC_21@1.1.1
func TestTheReviewerIsNotTheAuthorOfADocument(t *testing.T) {
	result := `- object_key: RR_01
  target_ids: [formats.md]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`
	git := history("Mathias Fischer")
	git.files["c1"]["formats.md"] = "# formats\n"
	expect(t, check(t, result, git), "the reviewer Mathias Fischer is an author of formats.md")

	git = history("Claude Opus 5.5")
	git.files["c1"]["formats.md"] = "# formats\n"
	if got := check(t, result, git); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_14@1.1.1
func TestAMergeIsNoAuthorship(t *testing.T) {
	//	c1 on a branch, written by Claude Code; m1 merges it into a main that changed the file
	//	as well, so m1 equals neither parent - merged by the reviewer
	git := history("Claude Opus 5.5")
	git.files["c0"] = map[string]string{"formats.md": "a\n"}
	git.files["c1"]["formats.md"] = "a\nb\n"
	git.files["m1"] = map[string]string{"reviews/criteria-tool.yaml": criteria, "formats.md": "z\na\nb\n"}
	git.parent = map[string]string{"c1": "c0", "m1": "c0"}
	git.changed = map[string][]string{"formats.md": {"m1", "c1"}}
	git.merges = map[string]bool{"m1": true}
	git.authors = map[string]string{"m1": "Mathias Fischer"}
	got := check(t, `- object_key: RR_01
  target_ids: [formats.md]
  revision: 'm1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`, git)
	if len(got) != 0 {
		t.Errorf("want no finding, the merge is no authorship; got %q", got)
	}
}

// @implements TC_17@1.1.1
func TestAnUnfilledSkeletonIsReported(t *testing.T) {
	got := check(t, `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: todo, RC_02: todo}
  verdict: todo
`, history("Claude Opus 5.5", "Claude Opus 5.5"))
	expect(t, got, "criterion RC_01 is neither passed nor failed")
	expect(t, got, "criterion RC_02 is neither passed nor failed")
	expect(t, got, "the verdict is neither passed nor failed")
}

// order checks the result in twenty runs against the findings wanted, in their order: a map is
// walked in another order each time, so twenty runs show the order does not follow it.
func order(t *testing.T, revision string, git fakeGit, last string) {
	t.Helper()
	git.files["c1"]["reviews/criteria-tool.yaml"] = `level: tool
criteria:
- object_key: RC_01
- object_key: RC_02
- object_key: RC_03
- object_key: RC_04
- object_key: RC_05
`
	result := `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: '` + revision + `'
  criteria: {RC_05: todo, RC_04: todo, RC_01: todo, RC_03: failed, RC_09: passed, RC_07: passed}
  verdict: todo
`
	want := []string{
		"criterion RC_01 is neither passed nor failed",
		"criterion RC_02 is missing",
		"criterion RC_03 failed with no note",
		"criterion RC_04 is neither passed nor failed",
		"criterion RC_05 is neither passed nor failed",
		"criterion RC_07 is not in reviews/criteria-tool.yaml",
		"criterion RC_09 is not in reviews/criteria-tool.yaml",
		"the verdict is neither passed nor failed",
		last,
	}
	for run := 0; run < 20; run++ {
		got := check(t, result, git)
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Fatalf("run %d:\n got %q\nwant %q", run, got, want)
		}
	}
}

// @implements TC_31@2.1.1
func TestTheFindingsInTheSameOrderInEveryRun(t *testing.T) {
	order(t, "nope", history("Claude Opus 5.5", "Claude Opus 5.5"), "the revision nope is no commit")
}

// @implements TC_135@1.1.1
func TestTheFindingsUpToTheReviewerInTheSameOrder(t *testing.T) {
	order(t, "c1", history("Mathias Fischer"), "the reviewer Mathias Fischer is an author of TOOL_03")
}

// @implements TC_34@1.1.1, TC_35@1.1.1, TC_36@1.1.1
func TestACompletedReviewHasItsCompletedAt(t *testing.T) {
	result := `- object_key: RR_01
  target_ids: [TOOL_03]
  revision: 'c1'
  criteria: {RC_01: passed, RC_02: passed}
  verdict: passed
`
	git := history("Claude Opus 5.5", "Claude Opus 5.5")

	expect(t, checkRecord(t, "status: completed\n", result, git), "completed, and no completed_at")
	expect(t, checkRecord(t, "status: completed\ncompleted_at: '2026-10-25T13:29:59+01:00'\n", result, git),
		"the completed_at 2026-10-25T13:29:59+01:00 is before the created_at 2026-10-25T13:30:00+01:00")
	for _, at := range []string{"2026-10-25T13:30:00+01:00", "2026-10-25T14:00:00+01:00"} {
		if got := checkRecord(t, "status: completed\ncompleted_at: '"+at+"'\n", result, git); len(got) != 0 {
			t.Errorf("completed_at %s: want no finding, got %q", at, got)
		}
	}
}
