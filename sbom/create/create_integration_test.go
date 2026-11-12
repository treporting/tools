//go:build integration

package create

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The module of the tools is the subsystem here: its commands link yaml.v3, jsonschema and x/text.
const toolsRoot = "../.."

// @verifies SR_36@2.1.1
func TestGoListOfTheTools(t *testing.T) {
	list, err := GoList(toolsRoot)
	if err != nil {
		t.Fatal(err)
	}
	modules := map[string]Package{}
	for _, p := range GoModules(list, LicenceIn) {
		modules[p.Name] = p
	}
	if p := modules["gopkg.in/yaml.v3"]; p.PURL != "pkg:golang/gopkg.in/yaml.v3@v3.0.1" || p.Licence != "Apache-2.0 AND MIT" {
		t.Errorf("yaml.v3: %+v", p)
	}
	if p := modules["golang.org/x/text"]; p.Licence != "BSD-3-Clause" {
		t.Errorf("x/text: %+v", p)
	}
	if _, ok := modules["github.com/treporting/specs-and-deploy/tools"]; ok {
		t.Error("the module of the subsystem is listed")
	}
}

// @verifies SR_44@1.1.1
func TestGoEnvNamesTheVersionOfGo(t *testing.T) {
	version, err := GoEnv(toolsRoot, "GOVERSION")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(Stdlib(version, "").PURL, "pkg:golang/stdlib@go1.") {
		t.Errorf("version %q", version)
	}
}

// @implements TC_108@2.1.1
func TestCommandWritesSbomMDInTheRootAndKeepsTheHistory(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sbom")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/sbom").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	repo := filepath.Join(t.TempDir(), "viewer-trace")
	if err := os.MkdirAll(repo, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LICENSE"), []byte(apacheText), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "package-lock.json"), []byte(lock), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "documents.yaml"), []byte(documentsYAML), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, label := range []string{"1.0.0", "1.1.0"} {
		if out, err := exec.Command(bin, "-root", repo, "-label", label, "-ticket", "343").CombinedOutput(); err != nil {
			t.Fatalf("sbom: %v\n%s", err, out)
		}
	}
	md, err := os.ReadFile(filepath.Join(repo, Name))
	if err != nil {
		t.Fatal(err)
	}
	rows := History(md)
	if len(rows) != 2 || !strings.HasPrefix(rows[0], "| 1.1.1 | ") || !strings.Contains(rows[0], " | 1.0.0 | #343 | sbom ") ||
		!strings.HasPrefix(rows[1], "| 2.1.1 | ") || !strings.Contains(rows[1], " | 1.1.0 | #343 | sbom ") {
		t.Errorf("rows %q", rows)
	}
	documents, err := os.ReadFile(filepath.Join(repo, "documents.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(documents), "- name: "+Name+"\n  type: document\n  object_version: '2.1.1'\n  object_hash: '"+Hash(md)+"'\n  status: completed\n  verification_method: [generated]\n") {
		t.Errorf("documents.yaml:\n%s", documents)
	}
	if !strings.Contains(string(md), "| pkg:github/treporting/viewer-trace@1.1.0 | viewer-trace | 1.1.0 | Apache-2.0 |") {
		t.Errorf("own package missing:\n%s", md)
	}
}

// @implements TC_122@1.1.1, TC_123@1.1.1
func TestCheckFailsWhereSbomMDDoesNotListWhatIsDelivered(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sbom")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/sbom").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	repo := filepath.Join(t.TempDir(), "viewer-trace")
	if err := os.MkdirAll(repo, 0o750); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{"LICENSE": apacheText, "package-lock.json": lock, "documents.yaml": documentsYAML} {
		if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if out, err := exec.Command(bin, "-root", repo, "-label", "1.0.0", "-ticket", "343").CombinedOutput(); err != nil {
		t.Fatalf("sbom: %v\n%s", err, out)
	}
	written, err := os.ReadFile(filepath.Join(repo, Name))
	if err != nil {
		t.Fatal(err)
	}

	if out, err := exec.Command(bin, "-root", repo, "-check").CombinedOutput(); err != nil {
		t.Errorf("check of the sbom.md just written: %v\n%s", err, out)
	}
	after, err := os.ReadFile(filepath.Join(repo, Name))
	if err != nil || string(after) != string(written) {
		t.Errorf("check wrote sbom.md: %v", err)
	}

	changed := strings.Replace(lock, `"license": "MIT"`, `"license": "ISC"`, 1)
	if err := os.WriteFile(filepath.Join(repo, "package-lock.json"), []byte(changed), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(bin, "-root", repo, "-check").CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 || !strings.Contains(string(out), "delivered is") {
		t.Errorf("check of a changed lock: %v\n%s", err, out)
	}
}

// sbomIn builds sbom and runs it with args in a repo of git with a package-lock.json in its root,
// and answers the names in its folder sbom/.
func sbomIn(t *testing.T, args ...string) (string, []string) {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "sbom")
	if out, err := exec.Command("go", "build", "-o", bin, "../cmd/sbom").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	repo := filepath.Join(t.TempDir(), "viewer-trace")
	if err := os.MkdirAll(repo, 0o750); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{"LICENSE": apacheText, "package-lock.json": lock, "documents.yaml": documentsYAML} {
		if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, git := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "a"}} {
		if out, err := exec.Command("git", append([]string{"-C", repo}, git...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", git, err, out)
		}
	}
	if out, err := exec.Command(bin, append([]string{"-root", repo, "-label", "1.0.0", "-ticket", "343"}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("sbom: %v\n%s", err, out)
	}
	entries, err := os.ReadDir(filepath.Join(repo, Folder))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return repo, names
}

// @implements TC_149@1.2.1
func TestCommandWithoutStoreWritesOnlySbomMD(t *testing.T) {
	repo, names := sbomIn(t)
	if len(names) != 1 || names[0] != "sbom.md" {
		t.Errorf("sbom/ holds %v", names)
	}
	md, err := os.ReadFile(filepath.Join(repo, Name))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(md), "| root-package-lock.json | `npm install` | - |") {
		t.Errorf("Original data:\n%s", md)
	}
}

// @implements TC_150@1.2.1, TC_151@1.1.1
func TestCommandWithStoreWritesTheOriginalsAndTheirRecord(t *testing.T) {
	repo, names := sbomIn(t, "-store")
	if len(names) != 3 {
		t.Fatalf("sbom/ holds %v", names)
	}
	record, lockName := names[0], names[1] // <ts>.yaml sorts before <ts>_root-package-lock.json
	ts, ok := strings.CutSuffix(lockName, "_root-package-lock.json")
	if !ok || len(ts) != len("20261105T100000") || record != ts+".yaml" || names[2] != "sbom.md" {
		t.Fatalf("sbom/ holds %v", names)
	}
	md, err := os.ReadFile(filepath.Join(repo, Name))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(md), "| root-package-lock.json | `npm install` | ./sbom/"+lockName+" |") {
		t.Errorf("Original data:\n%s", md)
	}
	data, err := os.ReadFile(filepath.Join(repo, Folder, record))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "- file: "+lockName+"\n    made_by: npm install\n") || !strings.Contains(string(data), "document: sbom/sbom.md\n") {
		t.Errorf("%s:\n%s", record, data)
	}
}
