//go:build integration

package write

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// objectcode builds the command objectcode into a folder of the test.
func objectcode(t *testing.T) string {

	t.Helper()
	bin := filepath.Join(t.TempDir(), "objectcode")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/objectcode").CombinedOutput(); err != nil { // #nosec G204 -- the arguments come from the test
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

// repo is a git repo of one commit with a documents.yaml, and a section beside it.
func repo(t *testing.T) (root, section string) {

	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "documents.yaml"), []byte("- name: README.md\n  type: document\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, git := range [][]string{
		{"init", "-q"}, {"add", "documents.yaml"},
		{"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "one"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", root}, git...)...).CombinedOutput(); err != nil { // #nosec G204 -- the arguments come from the test
			t.Fatalf("git %v: %v\n%s", git, err, out)
		}
	}
	section = filepath.Join(t.TempDir(), "a.md")
	if err := os.WriteFile(section, []byte("## a\n\n- Binary: a (cmd/a)\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, section
}

func TestDocumentKeepsItsHistoryAndSetsItsEntry(t *testing.T) {

	bin := objectcode(t)
	root, section := repo(t)
	for _, args := range [][]string{
		{"-document", "-root", root, section},
		{"-document", "-root", root, "-label", "T047.0", "-ticket", "37", section},
	} {
		if out, err := exec.Command(bin, args...).CombinedOutput(); err != nil { // #nosec G204 -- the arguments come from the test
			t.Fatalf("objectcode %v: %v\n%s", args, err, out)
		}
	}
	md, err := os.ReadFile(filepath.Join(root, "testreports", "object-code.md")) // #nosec G304 -- a folder of the test
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(md), "- Label: T047.0\n") || !strings.Contains(string(md), "## a\n") {
		t.Errorf("the report is\n%s", md)
	}
	if rows := strings.Count(string(md), "| objectcode "); rows != 2 {
		t.Errorf("the Change history has %d rows, want 2:\n%s", rows, md)
	}
	if !strings.Contains(string(md), " | T047.0 | #37 | objectcode ") {
		t.Errorf("the second row names not the label and the ticket:\n%s", md)
	}
	documents, err := os.ReadFile(filepath.Join(root, "documents.yaml")) // #nosec G304 -- a folder of the test
	if err != nil {
		t.Fatal(err)
	}
	want := "- name: testreports/object-code.md\n  type: document\n  object_version: '2.1.1'\n  object_hash: '"
	if !strings.Contains(string(documents), want) || !strings.Contains(string(documents), "  status: completed\n  verification_method: [generated]\n") {
		t.Errorf("documents.yaml is\n%s", documents)
	}
}

func TestDocumentWithALabelWithoutATicketFails(t *testing.T) {

	bin := objectcode(t)
	root, section := repo(t)
	if out, err := exec.Command(bin, "-document", "-root", root, "-label", "T047.0", section).CombinedOutput(); err == nil { // #nosec G204 -- the arguments come from the test
		t.Errorf("objectcode wrote a report with a label without a ticket:\n%s", out)
	}
}
