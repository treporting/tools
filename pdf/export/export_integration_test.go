//go:build integration

package export

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// repo makes a git repo with documents.yaml and two documents in two commits, and answers its
// directory and the revision of each document.
func repo(t *testing.T) (dir, readme, other string) {
	t.Helper()
	dir = t.TempDir()
	git := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t",
			"GIT_COMMITTER_EMAIL=t@t", "GIT_COMMITTER_DATE=2026-11-01T12:00:00Z", "GIT_AUTHOR_DATE=2026-11-01T12:00:00Z")
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
		return string(bytes.TrimSpace(out))
	}
	write := func(name, content string) {
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(name)), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-q")
	write("documents.yaml", documentsYAML)
	write("README.md", "# specs-and-deploy\n\nThe repo.\n")
	git("add", ".")
	git("commit", "-q", "-m", "one")
	readme = git("log", "-1", "--format=%h", "--abbrev=7")
	write("requirements/tool-requirements.md", "# Tool requirements\n\n| a | b |\n|---|---|\n| 1 | 2 |\n")
	git("add", ".")
	git("commit", "-q", "-m", "two")
	other = git("log", "-1", "--format=%h", "--abbrev=7")
	git("tag", "T1")
	write("README.md", "# changed after the tag\n")
	git("commit", "-q", "-am", "three")
	return dir, readme, other
}

// @implements TC_99@1.1.1, TC_133@1.1.1
func TestGitReadsTheRepoAtTheTag(t *testing.T) {
	dir, readme, other := repo(t)
	g := Git{Dir: dir, Ref: "T1"}
	if _, err := g.Commit(); err != nil {
		t.Fatal(err)
	}
	source, err := g.Show("README.md")
	if err != nil || string(source) != "# specs-and-deploy\n\nThe repo.\n" {
		t.Errorf("README.md at T1: %q %v", source, err)
	}
	for path, want := range map[string]string{"README.md": readme, "requirements/tool-requirements.md": other} {
		if got, err := g.Revision(path); err != nil || got != want {
			t.Errorf("revision of %s: %q %v, want %q", path, got, err, want)
		}
	}
	if date, err := g.Date(); err != nil || date != "2026-11-01" {
		t.Errorf("date %q %v", date, err)
	}
}

// @verifies SR_74@1.1.1
func TestGitFailsWhereThereIsNoSuchRef(t *testing.T) {
	dir, _, _ := repo(t)
	if _, err := (Git{Dir: dir, Ref: "T9"}).Commit(); err == nil {
		t.Error("no error for a tag that does not exist")
	}
}

// @verifies SR_68@1.1.1
func TestPrinterWritesAPDF(t *testing.T) {
	p, err := NewPrinter()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	page, err := Page("T", []byte("# T\n\n| a | b |\n|---|---|\n| 1 | 2 |\n"))
	if err != nil {
		t.Fatal(err)
	}
	pdf, err := p.Print(page, HeaderTemplate("T", "repo"), FooterTemplate("2026-11-01", "t.md · 1.1.1 · abcdef0 · draft"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Errorf("no PDF: %q", pdf[:min(len(pdf), 16)])
	}
}
