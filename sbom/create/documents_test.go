package create

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const documentsYAML = `# Every Markdown file.
schema: schemas/documents.schema.json
type: documents
documents:
- name: README.md
  type: document
  object_version: '3.1.1'
  object_hash: 'f2685eb'
  status: approved
  verification_method: [review]
`

const withSbom = documentsYAML + `- name: sbom/sbom.md
  type: document
  object_version: '2.1.1'
  object_hash: 'aaaaaaa'
  status: completed
  verification_method: [generated]
- name: devops.md
  type: document
  object_version: '1.1.1'
  object_hash: 'bbbbbbb'
  status: approved
  verification_method: [review]
`

// @implements TC_109@1.1.1
func TestVersionOfSbomMD(t *testing.T) {
	if got := Version([]byte(withSbom)); got != "2.1.1" {
		t.Errorf("got %q", got)
	}
	if got := Version([]byte(documentsYAML)); got != "" {
		t.Errorf("no entry: got %q", got)
	}
}

// @implements TC_109@1.1.1, TC_110@1.1.1
func TestNextRaisesX(t *testing.T) {
	for version, want := range map[string]string{"": "1.1.1", "1.1.1": "2.1.1", "9.2.3": "10.2.3"} {
		if got, err := Next(version); err != nil || got != want {
			t.Errorf("Next(%q) = %q, %v; want %q", version, got, err, want)
		}
	}
	for _, version := range []string{"1.1", "a.1.1"} {
		if _, err := Next(version); err == nil {
			t.Errorf("Next(%q): no error", version)
		}
	}
}

// @implements TC_120@1.1.1
func TestHashIsThatOfGitHashObject(t *testing.T) {
	content := []byte("# SBOM\n")
	file := filepath.Join(t.TempDir(), "sbom.md")
	if err := os.WriteFile(file, content, 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "hash-object", file).Output()
	if err != nil {
		t.Skip("git not there:", err)
	}
	if got := Hash(content); got != string(out[:7]) {
		t.Errorf("got %q, git %q", got, out[:7])
	}
}

// @implements TC_121@1.1.1
func TestEntryAddedWhereThereIsNone(t *testing.T) {
	got := string(Entry([]byte(documentsYAML), "1.1.1", "1234567"))
	want := documentsYAML + `- name: sbom/sbom.md
  type: document
  object_version: '1.1.1'
  object_hash: '1234567'
  status: completed
  verification_method: [generated]
`
	if got != want {
		t.Errorf("got\n%s", got)
	}
}

// @implements TC_120@1.1.1
func TestEntrySetInPlace(t *testing.T) {
	got := string(Entry([]byte(withSbom), "3.1.1", "1234567"))
	want := strings.Replace(strings.Replace(withSbom, "'2.1.1'", "'3.1.1'", 1), "'aaaaaaa'", "'1234567'", 1)
	if got != want {
		t.Errorf("got\n%s", got)
	}
}

// @implements TC_120@1.1.1
func TestEntryCompletedAndGenerated(t *testing.T) {
	changed := strings.Replace(withSbom, "status: completed\n  verification_method: [generated]", "status: draft\n  verification_method: [review]", 1)
	got := string(Entry([]byte(changed), "3.1.1", "1234567"))
	if !strings.Contains(got, "- name: sbom/sbom.md\n  type: document\n  object_version: '3.1.1'\n  object_hash: '1234567'\n  status: completed\n  verification_method: [generated]\n- name: devops.md") {
		t.Errorf("got\n%s", got)
	}
}
