package create

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/schema/check"
)

// originals answers the original data of a repo in root with a go.mod in tools/, a
// package-lock.json in the root and the image docker.io/nginx:1.30.
func originals(root string) []Original {
	return []Original{
		GoListOriginal(root, filepath.Join(root, "tools"), "gopkg.in/yaml.v3\tv3.0.1\t/m\n", "/m"),
		LockOriginal(root, root, []byte(`{"packages": {}}`)),
		ApkOriginal("docker.io/nginx:1.30", []byte("P:nginx\nV:1.30.5-r1\n")),
	}
}

// originalData answers the section Original data of an sbom.md written with the originals.
func originalData(originals []Original) string {
	md := Markdown(SBOM{Repo: "storages", Original: originals})
	_, after, _ := strings.Cut(md, "## Original data\n")
	text, _, _ := strings.Cut(after, "\n## ")
	return text
}

// @implements TC_149@1.2.1
func TestOriginalDataWithoutStore(t *testing.T) {
	data := originalData(originals("/r"))
	for _, want := range []string{
		"| File | Made by | Stored At |",
		"| golist-tools.txt | `cd tools && go list -deps -f '{{with .Module}}",
		"| root-package-lock.json | `npm install` | - |",
		"| apk-docker.io-nginx-1.30.txt | `docker run --rm --entrypoint cat docker.io/nginx:1.30 /lib/apk/db/installed` | - |",
	} {
		if !strings.Contains(data, want) {
			t.Errorf("Original data has no %q:\n%s", want, data)
		}
	}
	if strings.Count(data, " | - |\n") != 3 {
		t.Errorf("not every file stored at -:\n%s", data)
	}
}

// @implements TC_149@1.2.1
func TestOriginalDataEmpty(t *testing.T) {
	if data := originalData(nil); data != "\n`<empty>`\n" {
		t.Errorf("got %q", data)
	}
}

// @implements TC_150@1.2.1
func TestStoreWritesEachFileIntoSbom(t *testing.T) {
	root := t.TempDir()
	stored, err := Store(root, Folder, "20261105T100000_", originals(root))
	if err != nil {
		t.Fatal(err)
	}
	data := originalData(stored)
	for _, name := range []string{"20261105T100000_golist-tools.txt", "20261105T100000_root-package-lock.json", "20261105T100000_apk-docker.io-nginx-1.30.txt"} {
		if _, err := os.Stat(filepath.Join(root, "sbom", name)); err != nil {
			t.Errorf("not written: %v", err)
		}
		if !strings.Contains(data, "| ./sbom/"+name+" |") {
			t.Errorf("Original data does not name ./sbom/%s:\n%s", name, data)
		}
	}
}

// @implements TC_151@1.1.1
func TestOriginalYAMLNamesEveryFileStored(t *testing.T) {
	root := t.TempDir()
	stored, err := Store(root, Folder, "20261105T100000_", originals(root))
	if err != nil {
		t.Fatal(err)
	}
	record, err := OriginalYAML(Name, "1.2.0", "a1b2c3d", "2026-11-05T10:00:00+01:00", stored)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range stored {
		name := strings.TrimPrefix(o.StoredAt, "./sbom/")
		if !strings.Contains(string(record), "- file: "+name+"\n    made_by: ") {
			t.Errorf("%s not named:\n%s", name, record)
		}
	}

	// the record keeps to its schema
	schema, err := os.ReadFile("../../../schemas/original-data.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "schemas"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "schemas", "original-data.schema.json"), schema, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sbom", "20261105T100000.yaml"), record, 0o600); err != nil {
		t.Fatal(err)
	}
	findings, err := check.Check(root, []string{"sbom/20261105T100000.yaml"})
	if err != nil || len(findings) > 0 {
		t.Errorf("%v %v\n%s", err, findings, record)
	}
}

// @implements TC_156@1.1.1
func TestStoredGoListNamesNoLocalPath(t *testing.T) {
	root := t.TempDir()
	list := "gopkg.in/yaml.v3\tv3.0.1\t/home/u/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1\n"
	stored, err := Store(root, Folder, "20261105T100000_", []Original{GoListOriginal(root, filepath.Join(root, "tools"), list, "/home/u/go/pkg/mod")})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "sbom", "20261105T100000_golist-tools.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\t$GOMODCACHE/gopkg.in/yaml.v3@v3.0.1\n") || strings.Contains(string(data), "/home/u") {
		t.Errorf("%s holds %q", stored[0].StoredAt, data)
	}
}
