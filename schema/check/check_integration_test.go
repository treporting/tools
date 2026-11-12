//go:build integration

package check

import (
	"os/exec"
	"testing"
)

// The files of a real repo: a YAML of our data is checked; one in .github/ and one not under
// version control are not.
//
// @implements TC_94@1.1.1
func TestTheFilesOfAHistoryOfGit(t *testing.T) {

	dir := root(t, map[string]string{
		"ours.yaml":                   "schema: schemas/t.schema.json\nlevel: r1\n",
		".github/workflows/tools.yml": "name: tools\non: pull_request\n",
		"loose.yaml":                  "level: r1\n",
	})
	for _, args := range [][]string{{"init", "-q"}, {"add", "ours.yaml", ".github", "schemas"}} {
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	files, err := Files(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != "ours.yaml" {
		t.Fatalf("want only ours.yaml, got %q", files)
	}
	findings, err := Check(dir, files)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("want no finding, got %q", messages(findings))
	}
}
