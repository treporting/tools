package create

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// licences reads the licence text of a module directory from memory.
type licences map[string]string

func (l licences) read(dir string) string { return l[dir] }

// @implements TC_102@1.1.1, TC_103@1.1.1
func TestGoModulesWithVersionLicenceAndPURL(t *testing.T) {
	list := "gopkg.in/yaml.v3\tv3.0.1\t/cache/yaml\n" +
		"\n" + // a package of the standard library or of the subsystem prints nothing
		"golang.org/x/text\tv0.14.0\t/cache/text\n" +
		"gopkg.in/yaml.v3\tv3.0.1\t/cache/yaml\n" // a second package of the same module
	got := GoModules(list, licences{"/cache/yaml": mitText, "/cache/text": bsd3Text}.read)
	want := []Package{
		{PURL: "pkg:golang/gopkg.in/yaml.v3@v3.0.1", Name: "gopkg.in/yaml.v3", Version: "v3.0.1", Licence: "MIT"},
		{PURL: "pkg:golang/golang.org/x/text@v0.14.0", Name: "golang.org/x/text", Version: "v0.14.0", Licence: "BSD-3-Clause"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

// @verifies SR_36@2.1.1
// @verifies SR_43@1.1.1
func TestGoModuleWithoutLicenceFileIsUnknown(t *testing.T) {
	got := GoModules("example.org/m\tv1.0.0\t/cache/m\n", licences{}.read)
	if len(got) != 1 || got[0].Licence != Unknown {
		t.Errorf("got %+v", got)
	}
}

// @implements TC_116@1.1.1
func TestStdlib(t *testing.T) {
	got := Stdlib("go1.25.14", bsd3Text)
	want := Package{PURL: "pkg:golang/stdlib@go1.25.14", Name: "stdlib", Version: "go1.25.14", Licence: "BSD-3-Clause"}
	if got != want {
		t.Errorf("got %+v", got)
	}
}

const lock = `{
  "name": "viewer-trace",
  "lockfileVersion": 3,
  "packages": {
    "": {"name": "viewer-trace", "version": "0.1.0", "license": "Apache-2.0"},
    "node_modules/react": {"version": "18.3.1", "license": "MIT"},
    "node_modules/@babel/runtime": {"version": "7.29.10", "license": "MIT"},
    "node_modules/a/node_modules/b": {"version": "2.0.0", "license": "ISC"},
    "node_modules/eslint": {"version": "9.0.0", "license": "MIT", "dev": true},
    "node_modules/old": {"version": "1.0.0", "license": {"type": "MIT"}},
    "node_modules/words": {"version": "1.0.0", "license": "SEE LICENSE IN LICENSE.txt"}
  }
}`

// @implements TC_104@1.1.1, TC_105@1.1.1
func TestNpmPackagesWithoutDevAndWithoutTheSubsystem(t *testing.T) {
	got, err := NpmPackages([]byte(lock))
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i].PURL < got[j].PURL })
	want := []Package{
		{PURL: "pkg:npm/%40babel/runtime@7.29.10", Name: "@babel/runtime", Version: "7.29.10", Licence: "MIT", Named: "MIT"},
		{PURL: "pkg:npm/b@2.0.0", Name: "b", Version: "2.0.0", Licence: "ISC", Named: "ISC"},
		{PURL: "pkg:npm/old@1.0.0", Name: "old", Version: "1.0.0", Licence: Unknown},
		{PURL: "pkg:npm/react@18.3.1", Name: "react", Version: "18.3.1", Licence: "MIT", Named: "MIT"},
		{PURL: "pkg:npm/words@1.0.0", Name: "words", Version: "1.0.0", Licence: Unknown, Named: "SEE LICENSE IN LICENSE.txt"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// @verifies SR_37@2.1.1
func TestNpmLockThatIsNoJSONIsAnError(t *testing.T) {
	if _, err := NpmPackages([]byte("{")); err == nil {
		t.Error("no error")
	}
}

const installed = `C:Q1abc=
P:musl
V:1.2.6-r2
A:x86_64
L:MIT
o:musl

P:nginx
V:1.30.5-r1
L:2-clause BSD-like license

P:ca-certificates-bundle
V:20260909-r0
L:MPL-2.0 AND MIT
`

// @implements TC_106@1.1.1
func TestApkPackages(t *testing.T) {
	got := ApkPackages([]byte(installed))
	want := []Package{
		{PURL: "pkg:apk/alpine/musl@1.2.6-r2", Name: "musl", Version: "1.2.6-r2", Licence: "MIT", Named: "MIT"},
		{PURL: "pkg:apk/alpine/nginx@1.30.5-r1", Name: "nginx", Version: "1.30.5-r1", Licence: Unknown, Named: "2-clause BSD-like license"},
		{PURL: "pkg:apk/alpine/ca-certificates-bundle@20260909-r0", Name: "ca-certificates-bundle", Version: "20260909-r0", Licence: "MPL-2.0 AND MIT", Named: "MPL-2.0 AND MIT"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// @verifies SR_36@2.1.1
// @verifies SR_37@2.1.1
func TestDirsWithAsTheBuildStepsFindThem(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{".", "tools", "web/app", "web/node_modules/m", ".github/x", "docs"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o750); err != nil {
			t.Fatal(err)
		}
		if dir == "docs" {
			continue
		}
		if err := os.WriteFile(filepath.Join(root, dir, "go.mod"), []byte("module m\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := DirsWith(root, "go.mod")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{root, filepath.Join(root, "tools"), filepath.Join(root, "web/app")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}
