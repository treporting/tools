package check

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/review"
)

// fakeGit holds the files of the base commit in memory.
type fakeGit map[string]string

func (f fakeGit) Show(rev, path string) ([]byte, error) {
	if c, ok := f[path]; ok {
		return []byte(c), nil
	}
	return nil, os.ErrNotExist
}

const requirements = `schema: schemas/software-requirements.schema.json
type: requirements
level: r5
requirements:
- object_key: SR_01
  status: approved
- object_key: SR_02
  status: draft
`

var base = fakeGit{"tools/software-requirements.yaml": requirements}

// index is the index of the files of base, as if the working tree held them.
var index = review.NewIndex(func(file string) ([]byte, error) { return base.Show("", file) }, []string{"tools/software-requirements.yaml"})

func write(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, path), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func messages(findings []Finding) []string {
	var out []string
	for _, f := range findings {
		out = append(out, f.String())
	}
	return out
}

var kinds = Kinds{
	Whitelist: []string{"*.go", "*.js", "*.sh", "Dockerfile", "*.yaml"},
	Blacklist: []string{"*.md", "*.txt", "*.json", "*.png"},
}

func find(t *testing.T, files map[string]string) ([]Annotation, []Finding) {
	t.Helper()
	dir := write(t, files)
	var names []string
	for f := range files {
		names = append(names, f)
	}
	found, unknown, err := Find(dir, names, kinds)
	if err != nil {
		t.Fatal(err)
	}
	return found, unknown
}

func filesOf(found []Annotation) string {
	var got []string
	for _, a := range found {
		got = append(got, a.File)
	}
	sort.Strings(got)
	return strings.Join(got, " ")
}

// @implements TC_71@2.1.1
func TestOnlyCodeIsRead(t *testing.T) {
	found, unknown := find(t, map[string]string{
		"a.go":                  "// @implements SR_01\nfunc a() {}\n",
		"b.js":                  "/** @verifies SR_01 */\n",
		"Dockerfile":            "# @implements SR_01\nFROM scratch\n",
		"run.sh":                "  # @implements SR_01\n",
		"config.yaml":           "# @implements SR_01\nkey: value\n",
		"doc.md":                "// @implements SR_99\n",
		"datei.txt":             "Find(root, files):\n    # @implements SR_99\n",
		"schemas/x.schema.json": `{"$comment": "@implements SR_99"}`,
		"data.yaml":             "schema: schemas/x.schema.json\n# @implements SR_99\n",
		"picture.png":           "# @implements SR_99",
	})
	if got := filesOf(found); got != "Dockerfile a.go b.js config.yaml run.sh" {
		t.Errorf("want the five files of code, got %q", got)
	}
	if len(unknown) != 0 {
		t.Errorf("want no unknown kind, got %q", messages(unknown))
	}
}

// @implements TC_75@1.1.1, TC_76@1.1.1
func TestAKindOnNeitherList(t *testing.T) {
	found, unknown := find(t, map[string]string{
		"query.sql": "-- @implements SR_99\n",
		"a.go":      "// @implements SR_01\n",
	})
	got := messages(unknown)
	if len(got) != 1 || got[0] != "query.sql: its kind is set neither as code nor as no code (tools/annotation/kinds.yaml)" {
		t.Errorf("want query.sql reported, got %q", got)
	}
	if filesOf(found) != "a.go" {
		t.Errorf("want only a.go read, got %q", filesOf(found))
	}
}

// @verifies SR_24@3.1.1
func TestOnlyFirstInALineOfComment(t *testing.T) {
	//	built with \n, not as a string over lines: a line of a raw string that begins with a mark of
	//	comment would be read as an annotation of this file itself
	found, _ := find(t, map[string]string{"a_test.go": "package a\n\n" +
		"const s = \"// @implements SR_99\"\n" +
		"var list = []string{\n\t\"// @implements SR_98\",\n}\n" +
		"func f() {} // @implements SR_97\n\n" +
		"// @verifies SR_01\n" +
		"func TestA() {}\n"})
	if len(found) != 1 || found[0].IDs[0] != "SR_01" || found[0].Line != 9 {
		t.Errorf("want only the comment on line 9, got %+v", found)
	}
}

// @implements TC_69@2.1.1
func TestAnIdNotInTheBaseCommit(t *testing.T) {
	got := messages(Check(base, "b1", index, []Annotation{{"a.go", 3, "implements", []string{"SR_01", "SR_07"}, []string{"1.1.1", "1.1.1"}}}))
	if len(got) != 1 || got[0] != "a.go:3 @implements SR_07: no file of objects holds such an id" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_72@1.1.1
func TestAnIdNotApprovedInTheBaseCommit(t *testing.T) {
	got := messages(Check(base, "b1", index, []Annotation{{"a.go", 9, "implements", []string{"SR_02"}, []string{"1.1.1"}}}))
	if len(got) != 1 || got[0] != "a.go:9 @implements SR_02: draft in the base commit b1, not approved" {
		t.Errorf("got %q", got)
	}
}

// @verifies SR_24@3.1.1
func TestSeveralIdsInOneAnnotation(t *testing.T) {
	found, _ := find(t, map[string]string{"a.go": "// @implements SR_01, SR_02,SR_03\n"})
	if len(found) != 1 || strings.Join(found[0].IDs, " ") != "SR_01 SR_02 SR_03" {
		t.Errorf("got %+v", found)
	}
}

// @verifies SR_30@2.2.1
func TestTheVersionIsReadApartFromTheId(t *testing.T) {
	found, _ := find(t, map[string]string{"a.go": "// @implements SR_01@2.1.3, SR_02@1.1.1\n"})
	if len(found) != 1 || strings.Join(found[0].IDs, " ") != "SR_01 SR_02" || strings.Join(found[0].Versions, " ") != "2.1.3 1.1.1" {
		t.Errorf("got %+v", found)
	}
	if got := messages(Unversioned(found)); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_77@1.1.1
func TestAnIdWithoutAVersion(t *testing.T) {
	found, _ := find(t, map[string]string{"a.go": "// @implements SR_01@1.1.1, SR_02, SR_03@1.1\n"})
	got := messages(Unversioned(found))
	want := []string{
		"a.go:1 @implements SR_02: names no version, as SR_02@x.y.z",
		"a.go:1 @implements SR_03: names no version, as SR_03@x.y.z",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got %q, want %q", got, want)
	}
}

// @implements TC_70@2.1.1
func TestAVerifiesOfARequirementIsNotReported(t *testing.T) {
	if got := messages(Check(base, "b1", index, []Annotation{{"a_test.go", 9, "verifies", []string{"SR_01"}, []string{"1.1.1"}}})); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// An @implements of a test procedure names its test case, which no file of requirements holds.
//
// @implements TC_134@1.1.1
func TestAnImplementsOfATestCaseIsNotReported(t *testing.T) {
	git := fakeGit{
		"tools/software-requirements.yaml": requirements,
		"tools/testcases/software.yaml":    "type: testcases\ntestcases:\n- object_key: TC_01\n  status: approved\n",
	}
	files := []string{"tools/software-requirements.yaml", "tools/testcases/software.yaml"}
	ix := review.NewIndex(func(file string) ([]byte, error) { return git.Show("", file) }, files)
	if got := messages(Check(git, "b1", ix, []Annotation{{"a_test.go", 9, "implements", []string{"TC_01"}, []string{"1.1.1"}}})); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_78@1.1.1
func TestAVersionWithoutZ(t *testing.T) {
	found, _ := find(t, map[string]string{"a_test.go": "// @verifies SR_01@1.1\n"})
	got := messages(Unversioned(found))
	if len(got) != 1 || got[0] != "a_test.go:1 @verifies SR_01: names no version, as SR_01@x.y.z" {
		t.Errorf("got %q", got)
	}
}

func refinesIn(t *testing.T, files map[string]string) []Annotation {
	t.Helper()
	dir := write(t, files)
	var names []string
	for f := range files {
		names = append(names, f)
	}
	sort.Strings(names)
	found, err := Refines(dir, names)
	if err != nil {
		t.Fatal(err)
	}
	return found
}

// @implements TC_74@1.1.1, TC_80@1.1.1, TC_81@1.1.1
func TestRefinesIsReadInTheCommentOfASchema(t *testing.T) {
	found := refinesIn(t, map[string]string{
		"schemas/a.schema.json": "{\n  \"$schema\": \"x\",\n  \"$comment\": \"@refines SR_01@1.1.1, SR_07@1.1.1\"\n}\n",
		"schemas/b.schema.json": "{\n  \"$comment\": \"no annotation\"\n}\n",
		"other/c.schema.json":   "{\n  \"$comment\": \"@refines SR_01@1.1.1\"\n}\n",
	})
	if len(found) != 1 || found[0].File != "schemas/a.schema.json" || found[0].Line != 3 || found[0].Kind != "refines" {
		t.Fatalf("got %+v", found)
	}
	got := messages(Check(base, "b1", index, found))
	if len(got) != 1 || got[0] != "schemas/a.schema.json:3 @refines SR_07: no file of objects holds such an id" {
		t.Errorf("got %q", got)
	}
}

// @verifies SR_25@2.2.1
func TestARefinesNotApprovedInTheBaseCommit(t *testing.T) {
	found := refinesIn(t, map[string]string{"schemas/a.schema.json": "{\"$comment\": \"@refines SR_02@1.1.1\"}\n"})
	got := messages(Check(base, "b1", index, found))
	if len(got) != 1 || got[0] != "schemas/a.schema.json:1 @refines SR_02: draft in the base commit b1, not approved" {
		t.Errorf("got %q", got)
	}
}

// @verifies SR_30@2.2.1
func TestARefinesWithoutAVersion(t *testing.T) {
	found := refinesIn(t, map[string]string{"schemas/a.schema.json": "{\"$comment\": \"@refines SR_01\"}\n"})
	got := messages(Unversioned(found))
	if len(got) != 1 || got[0] != "schemas/a.schema.json:1 @refines SR_01: names no version, as SR_01@x.y.z" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_79@1.1.1
func TestARefinesWithAVersion(t *testing.T) {
	found := refinesIn(t, map[string]string{"schemas/a.schema.json": "{\"$comment\": \"@refines ARC_22@2.1.1\"}\n"})
	if got := messages(Unversioned(found)); len(found) != 1 || len(got) != 0 {
		t.Errorf("want the @refines read and no finding, got %+v and %q", found, got)
	}
}

func writeFile(dir, path, content string) error {
	if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), 0o750); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, path), []byte(content), 0o600)
}
