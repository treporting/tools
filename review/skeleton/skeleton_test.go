package skeleton

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/treporting/tools/review"
)

// fakeGit is a history of two commits kept in memory: c1 changed the criteria, c2 is HEAD.
type fakeGit struct {
	files map[string]string // path -> content, the same in both commits
}

func (f fakeGit) Show(rev, path string) ([]byte, error) {
	if c, ok := f.files[path]; ok {
		return []byte(c), nil
	}
	return nil, os.ErrNotExist
}

func (f fakeGit) FileCommits(rev, path string) ([]string, error) {
	if strings.HasPrefix(path, "reviews/criteria-") {
		return []string{"c1full"}, nil
	}
	return []string{"c2full"}, nil
}

func (f fakeGit) Short(rev string) (string, error) {
	if rev == "HEAD" {
		return "c2", nil
	}
	return strings.TrimSuffix(rev, "full"), nil
}

func history() fakeGit {
	return fakeGit{files: map[string]string{
		"reviews/criteria-tool.yaml":      "level: tool\ncriteria:\n- object_key: RC_1\n- object_key: RC_2\n",
		"reviews/criteria-documents.yaml": "level: document\ncriteria:\n- object_key: RC_1\n",
		"requirements/tool-requirements.yaml": `type: requirements
level: tool
requirements:
- object_key: TOOL_3
  object_text: check every review result
- object_key: TOOL_4
  object_text: create the skeleton
`,
		"harness.md": "# Harness\n",
	}}
}

// indexOf makes the index of the files of the history.
func indexOf(git fakeGit) *review.Index {
	var names []string
	for file := range git.files {
		names = append(names, file)
	}
	sort.Strings(names)
	return review.NewIndex(func(file string) ([]byte, error) { return git.Show("HEAD", file) }, names)
}

var now = time.Date(2026, 10, 26, 14, 30, 5, 0, time.FixedZone("CET", 3600))

// makeSkeleton writes a skeleton into a temporary root and answers its path and content.
func makeSkeleton(t *testing.T, targets ...string) (string, string, error) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "reviews"), 0o750); err != nil {
		t.Fatal(err)
	}
	path, err := Make(root, Request{Reviewer: "Mathias Fischer", Scope: "tool-03-04", Targets: targets, Now: now}, history(), indexOf(history()))
	if err != nil {
		entries, _ := os.ReadDir(filepath.Join(root, "reviews"))
		if len(entries) != 0 {
			t.Errorf("want nothing written, got %d files", len(entries))
		}
		return "", "", err
	}
	data, _ := os.ReadFile(filepath.Join(root, path)) // #nosec G304 -- the path comes from the user of the tool, not from outside
	return path, string(data), nil
}

func contains(t *testing.T, content string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(content, w) {
			t.Errorf("no %q in\n%s", w, content)
		}
	}
}

// @implements TC_23@1.1.1
func TestTheNameTheReviewerAndTheTime(t *testing.T) {
	path, content, err := makeSkeleton(t, "TOOL_3", "TOOL_4")
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join("reviews", "results", "20261026T143005_tool-03-04.yaml") {
		t.Errorf("path %s", path)
	}
	contains(t, content, "created_by: Mathias Fischer\n", "created_at: '2026-10-26T14:30:05+01:00'\n")
}

// @implements TC_24@1.1.1, TC_25@1.1.1
func TestTheCriteriaOfTheLevel(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "criteria: reviews/criteria-tool.yaml\n", "criteria_revision: 'c1'\n")

	_, content, err = makeSkeleton(t, "harness.md")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "criteria: reviews/criteria-documents.yaml\n", "  target_ids: [harness.md]\n")
}

// @implements TC_26@1.1.1
func TestTheRevisionIsHead(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "  revision: 'c2'\n")
}

// @implements TC_27@1.1.1
func TestEveryCriterionAndTheVerdictAreTodo(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3", "TOOL_4")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "- object_key: RR_1\n  target_ids: [TOOL_3]\n", "- object_key: RR_2\n  target_ids: [TOOL_4]\n")
	if n := strings.Count(content, ": todo\n"); n != 6 {
		t.Errorf("want 6 todo, two criteria and a verdict for each of two; got %d", n)
	}
	if strings.Contains(content, "passed") || strings.Contains(content, "failed") {
		t.Errorf("a skeleton holds no judgement:\n%s", content)
	}
}

// @implements TC_28@1.1.1, TC_29@1.1.1, TC_30@1.1.1
func TestNothingWhereATargetIsMissingOrTheLevelsDiffer(t *testing.T) {
	if _, _, err := makeSkeleton(t, "TOOL_99"); err == nil || !strings.Contains(err.Error(), "no file holds TOOL_99") {
		t.Errorf("want TOOL_99 not found, got %v", err)
	}
	if _, _, err := makeSkeleton(t, "nothing.md"); err == nil {
		t.Errorf("want nothing.md not found")
	}
	if _, _, err := makeSkeleton(t, "TOOL_3", "harness.md"); err == nil || !strings.Contains(err.Error(), "one review, one level") {
		t.Errorf("want two levels refused, got %v", err)
	}
}

// @implements TC_32@1.1.1
func TestTheSkeletonNamesItsSchemaTypeAndLevel(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "schema: schemas/review-result.schema.json\ntype: review\nlevel: tool\n")
}

// @implements TC_33@1.1.1
func TestTheSkeletonIsAnOpenRecord(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "level: tool\nobject_version: '1.1.1'\nstatus: open\n")
}

// twoSubsystems is the history of a code repo whose two subsystems both hold SR_1.
func twoSubsystems() fakeGit {
	sr := "type: requirements\nlevel: r5\nrequirements:\n- object_key: SR_1\n  object_text: a behaviour\n"
	return fakeGit{files: map[string]string{
		"reviews/criteria-r5.yaml": "level: r5\ncriteria:\n- object_key: RC_1\n",
		"a/requirements.yaml":      sr,
		"b/requirements.yaml":      sr,
	}}
}

// makeIn writes a skeleton of the history git into a temporary root, with file as -file.
func makeIn(t *testing.T, git fakeGit, file string, targets ...string) (string, string, error) {
	t.Helper()
	root := t.TempDir()
	path, err := Make(root, Request{Reviewer: "Mathias Fischer", Scope: "sr-1", Targets: targets, File: file, Now: now}, git, indexOf(git))
	if err != nil {
		if _, statErr := os.Stat(filepath.Join(root, results)); statErr == nil {
			t.Errorf("want nothing written, got %s", results)
		}
		return "", "", err
	}
	data, _ := os.ReadFile(filepath.Join(root, path)) // #nosec G304 -- the path comes from the user of the tool, not from outside
	return path, string(data), nil
}

// @implements TC_191@1.1.1
func TestTheFileGivenIsTheFileOfTheResult(t *testing.T) {
	_, content, err := makeIn(t, twoSubsystems(), "b/requirements.yaml", "SR_1")
	if err != nil {
		t.Fatal(err)
	}
	contains(t, content, "  target_ids: [SR_1]\n  file: b/requirements.yaml\n")
}

// @implements TC_192@1.1.1
func TestNothingWhereAnIdStandsInTwoFilesAndNoFileIsGiven(t *testing.T) {
	_, _, err := makeIn(t, twoSubsystems(), "", "SR_1")
	if err == nil || !strings.Contains(err.Error(), "a/requirements.yaml, b/requirements.yaml") {
		t.Errorf("want both files named, got %v", err)
	}
}

func TestNoFieldFileWhereTheIdIsUnique(t *testing.T) {
	_, content, err := makeSkeleton(t, "TOOL_3")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(content, "file:") {
		t.Errorf("an id in one file wants no field file:\n%s", content)
	}
}
