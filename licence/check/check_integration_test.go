//go:build integration

package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// @verifies SR_49@1.1.1
// @verifies SR_50@1.1.1
func TestCommandNamesWhatFailsAndExitsWithOne(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "licencecheck")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/licencecheck").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "sbom"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "sbom", "sbom.md"), []byte(sbomMD), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "licences_assessed.yaml"), []byte(
		"- purl: pkg:apk/alpine/nginx@1.30.5-r1\n  licence: BSD-2-Clause\n  reason: the text of nginx is the BSD-2-Clause\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command(bin, "-root", repo).CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Fatalf("exit: %v\n%s", err, out)
	}
	want := "pkg:golang/example.org/gpl@v1.0.0 GPL-3.0-only\npkg:npm/words@1.0.0 unknown\n6 packages checked, 2 fail\n"
	if string(out) != want {
		t.Errorf("got\n%s", out)
	}
	if strings.Contains(string(out), "nginx") {
		t.Error("the assessed licence is not taken")
	}
}
