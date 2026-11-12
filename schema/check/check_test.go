package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const schema = `{
  "$comment": "@refines ARC_22",
  "type": "object",
  "required": ["schema", "level"],
  "additionalProperties": false,
  "properties": {
    "schema": {"const": "schemas/t.schema.json"},
    "level": {"enum": ["r1", "r2"]},
    "items": {"type": "array", "items": {"type": "object", "properties": {"status": {"enum": ["draft", "approved"]}}}}
  }
}`

// root writes the schema and the files into a temporary root.
func root(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	files["schemas/t.schema.json"] = schema
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

func expect(t *testing.T, got []string, want string) {
	t.Helper()
	for _, g := range got {
		if strings.Contains(g, want) {
			return
		}
	}
	t.Errorf("no finding %q in %q", want, got)
}

// @implements TC_91@1.1.1, TC_92@2.1.1
func TestEachPlaceThatDoesNotMatch(t *testing.T) {
	dir := root(t, map[string]string{
		"good.yaml": "schema: schemas/t.schema.json\nlevel: r1\n",
		"bad.yaml":  "schema: schemas/t.schema.json\nlevel: r9\nitems:\n- status: done\n",
	})
	findings, err := Check(dir, []string{"good.yaml", "bad.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	got := messages(findings)
	if len(got) != 2 {
		t.Errorf("want two findings, both in bad.yaml; got %q", got)
	}
	expect(t, got, "bad.yaml /level:")
	expect(t, got, "bad.yaml /items/0/status:")
}

// @implements TC_93@1.1.1
func TestAFieldTheSchemaDoesNotKnow(t *testing.T) {
	dir := root(t, map[string]string{
		"other.yaml": "schema: schemas/t.schema.json\nlevel: r1\nextra: 1\n",
	})
	findings, err := Check(dir, []string{"other.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, messages(findings), "other.yaml /:")
}

// @implements TC_95@1.1.1, TC_96@1.1.1
func TestNoSchemaOrASchemaThatIsNotThere(t *testing.T) {
	dir := root(t, map[string]string{
		"none.yaml":    "level: r1\n",
		"missing.yaml": "schema: schemas/gone.schema.json\nlevel: r1\n",
	})
	findings, err := Check(dir, []string{"none.yaml", "missing.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	got := messages(findings)
	expect(t, got, "none.yaml: names no schema")
	expect(t, got, "missing.yaml: names the schema schemas/gone.schema.json, which is not there")
}

// a file of this repo that names the schema only the repo above holds, and does not match it
var namesSchemaAbove = map[string]string{"bad.yaml": "schema: schemas/t.schema.json\nlevel: r9\n"}

// own writes the files into a temporary root that holds no schema.
func own(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		if err := os.WriteFile(filepath.Join(dir, path), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// @implements TC_164@1.1.1
func TestASchemaIsTakenFromTheRepoAbove(t *testing.T) {
	above := root(t, map[string]string{})
	findings, err := CheckWithAbove(own(t, namesSchemaAbove), above, []string{"bad.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, messages(findings), "bad.yaml /level")
}

// @implements TC_165@1.1.1
func TestWithoutTheRepoAboveItsSchemaIsNotThere(t *testing.T) {
	findings, err := Check(own(t, namesSchemaAbove), []string{"bad.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, messages(findings), "names the schema schemas/t.schema.json, which is not there")
}
