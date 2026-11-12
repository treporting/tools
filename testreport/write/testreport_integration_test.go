//go:build integration

package write

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// levelsYAML names r5tc as levels.yaml does.
const levelsYAML = "levels:\n- key: r5tc\n  name: subsystem test\n  abbreviation: SwVal\n"

// testreportRun builds testreport and runs it with args, for the test level r5tc of -dir tools, in a
// repo of git with levels.yaml as given, with junit-go.xml and junit-vitest.xml given from outside the
// repo; and answers the repo, its output and its error.
func testreportRun(t *testing.T, levels string, args ...string) (string, []byte, error) {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "testreport")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/testreport").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	repo := t.TempDir()
	for name, content := range map[string]string{
		"documents.yaml":                "documents:\n",
		"levels.yaml":                   levels,
		"tools/testcases/software.yaml": "testcases:\n",
	} {
		if err := os.MkdirAll(filepath.Join(repo, filepath.Dir(name)), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, git := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "a"}} {
		if out, err := exec.Command("git", append([]string{"-C", repo}, git...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", git, err, out)
		}
	}
	runs := t.TempDir()
	goXML, vitestXML := filepath.Join(runs, "junit-go.xml"), filepath.Join(runs, "junit-vitest.xml")
	if err := os.WriteFile(goXML, []byte(goJUnit), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(vitestXML, []byte(vitestJUnit), 0o600); err != nil {
		t.Fatal(err)
	}
	all := append([]string{"-root", repo, "-dir", "tools", "-level", "r5tc",
		"-junit", goXML, "-made-by", "gotestsum --junitfile junit-go.xml -- ./...",
		"-junit", vitestXML, "-made-by", "vitest run --reporter=junit"}, args...)
	out, err := exec.Command(bin, all...).CombinedOutput()
	return repo, out, err
}

// testreportIn runs testreport as testreportRun does, with levelsYAML, and answers the repo and the
// names in tools/testreports/.
func testreportIn(t *testing.T, args ...string) (string, []string) {
	t.Helper()
	repo, out, err := testreportRun(t, levelsYAML, args...)
	if err != nil {
		t.Fatalf("testreport: %v\n%s", err, out)
	}
	entries, err := os.ReadDir(filepath.Join(repo, "tools", "testreports"))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return repo, names
}

// head answers the test report of the repo up to its first section.
func head(t *testing.T, repo string) string {
	t.Helper()
	md, err := os.ReadFile(filepath.Join(repo, "tools", "testreports", "software.md"))
	if err != nil {
		t.Fatal(err)
	}
	text, _, _ := strings.Cut(string(md), "\n## ")
	return text
}

// @implements TC_161@1.2.1
func TestTitleNamesTheTestLevel(t *testing.T) {
	repo, _ := testreportIn(t)
	h := head(t, repo)
	if !strings.HasPrefix(h, "# Test report for subsystem test (SwVal)\n") {
		t.Errorf("title:\n%s", h)
	}
	if strings.Contains(h, "r5") {
		t.Errorf("the head names r5:\n%s", h)
	}
}

// @implements TC_162@1.1.1
func TestLevelsFromTheRootOrFromLevels(t *testing.T) {
	repo, _ := testreportIn(t)
	if !strings.Contains(head(t, repo), "subsystem test (SwVal)") {
		t.Errorf("not read from <root>/levels.yaml:\n%s", head(t, repo))
	}
	other := filepath.Join(t.TempDir(), "levels.yaml")
	if err := os.WriteFile(other, []byte("levels:\n- key: r5tc\n  name: other test\n  abbreviation: OT\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	repo, _ = testreportIn(t, "-levels", other)
	if !strings.HasPrefix(head(t, repo), "# Test report for other test (OT)\n") {
		t.Errorf("not read from -levels:\n%s", head(t, repo))
	}
}

// @implements TC_163@1.1.1
func TestLevelNotNamedFails(t *testing.T) {
	repo, out, err := testreportRun(t, "levels:\n- key: r6tc\n  name: integration test\n  abbreviation: SwIT\n")
	if err == nil || !strings.Contains(string(out), "r5tc") {
		t.Errorf("no error naming r5tc: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(repo, "tools", "testreports", "software.md")); !os.IsNotExist(err) {
		t.Errorf("a test report is written: %v", err)
	}
}

// @implements TC_139@2.2.1
func TestCommandWithoutStoreCopiesNothing(t *testing.T) {
	_, names := testreportIn(t)
	if len(names) != 1 || names[0] != "software.md" {
		t.Errorf("tools/testreports/ holds %v", names)
	}
}

// @implements TC_154@1.1.1, TC_155@1.1.1
func TestCommandWithStoreCopiesEachJUnit(t *testing.T) {
	repo, names := testreportIn(t, "-store")
	if len(names) != 4 {
		t.Fatalf("tools/testreports/ holds %v", names)
	}
	// <ts>_r5tc-junit-go.xml, <ts>_r5tc-junit-vitest.xml, <ts>_r5tc.yaml, software.md
	ts, ok := strings.CutSuffix(names[0], "_r5tc-junit-go.xml")
	if !ok || len(ts) != len("20261105T100000") || names[1] != ts+"_r5tc-junit-vitest.xml" || names[2] != ts+"_r5tc.yaml" || names[3] != "software.md" {
		t.Fatalf("tools/testreports/ holds %v", names)
	}
	copied, err := os.ReadFile(filepath.Join(repo, "tools", "testreports", names[0]))
	if err != nil || string(copied) != goJUnit {
		t.Errorf("the copy of junit-go.xml: %v", err)
	}
	md, err := os.ReadFile(filepath.Join(repo, "tools", "testreports", "software.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"| junit-go.xml | `gotestsum --junitfile junit-go.xml -- ./...` | ./tools/testreports/" + names[0] + " |",
		"| junit-vitest.xml | `vitest run --reporter=junit` | ./tools/testreports/" + names[1] + " |",
	} {
		if !strings.Contains(string(md), want) {
			t.Errorf("Original data has no %q:\n%s", want, md)
		}
	}
	record, err := os.ReadFile(filepath.Join(repo, "tools", "testreports", names[2]))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"document: tools/testreports/software.md\n",
		"- file: " + names[0] + "\n    made_by: gotestsum --junitfile junit-go.xml -- ./...\n",
		"- file: " + names[1] + "\n    made_by: vitest run --reporter=junit\n",
	} {
		if !strings.Contains(string(record), want) {
			t.Errorf("%s has no %q:\n%s", names[2], want, record)
		}
	}
}
