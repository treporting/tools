package write

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// segment maps a column of the one line of a bundle to a line and column of the source, from 0, and
// where name is not empty, to that name.
type segment struct {
	column, line, sourceColumn int
	name                       string
}

// vlq is a number as a source map writes it.
func vlq(n int) string {

	const digits = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	v := n << 1
	if n < 0 {
		v = (-n << 1) | 1
	}
	var s strings.Builder
	for {
		d := v & 31
		v >>= 5
		if v > 0 {
			d |= 32
		}
		s.WriteByte(digits[d])
		if v == 0 {
			return s.String()
		}
	}
}

// sourceMapOf is a source map of a bundle of one line from one source, as webpack names it.
func sourceMapOf(t *testing.T, segments []segment) []byte {

	t.Helper()
	var names []string
	index := map[string]int{}
	var parts []string
	column, line, sourceColumn, name := 0, 0, 0, 0
	for _, s := range segments {
		p := vlq(s.column-column) + vlq(0) + vlq(s.line-line) + vlq(s.sourceColumn-sourceColumn)
		column, line, sourceColumn = s.column, s.line, s.sourceColumn
		if s.name != "" {
			i, ok := index[s.name]
			if !ok {
				i = len(names)
				index[s.name] = i
				names = append(names, s.name)
			}
			p += vlq(i - name)
			name = i
		}
		parts = append(parts, p)
	}
	data, err := json.Marshal(map[string]any{
		"version": 3, "sources": []string{"webpack://web/./src/a.js"}, "names": names, "mappings": strings.Join(parts, ","),
	})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// load is a function of the source src/a.js of the project web, lines 1 to 3, as istanbul names it.
var load = Function{Key: "web/src/a.js@load", Symbol: "load", File: "web/src/a.js", From: 1, To: 3, FromColumn: 16, ToColumn: endOfLine}

// at is the column of the first part in a bundle.
func at(t *testing.T, bundle, part string) int {

	t.Helper()
	i := strings.Index(bundle, part)
	if i < 0 {
		t.Fatalf("%q is not in %q", part, bundle)
	}
	return i
}

// hashOfBundle is the hash of the function load in a bundle of one function: its body maps to line
// 2 of the source; the property after `r.` maps to the name given.
func hashOfBundle(t *testing.T, bundle, property, called string) string {

	t.Helper()
	segments := []segment{
		{column: 0, line: 0, sourceColumn: 7},
		{column: at(t, bundle, "return"), line: 1, sourceColumn: 2},
		{column: at(t, bundle, property), line: 1, sourceColumn: 9, name: called},
		{column: at(t, bundle, property) + len(property), line: 1, sourceColumn: 16},
	}
	code, err := ReadBundle([]byte(bundle), sourceMapOf(t, segments), "web", []Function{load})
	if err != nil {
		t.Fatal(err)
	}
	_, hash, ok := code.Find(load)
	if !ok {
		t.Fatalf("load is not found in %q", bundle)
	}
	return hash
}

// @implements TC_186@1.1.1
func TestTheHashIgnoresTheNamesTheBundlerChose(t *testing.T) {

	one := hashOfBundle(t, "function a(n){return(0,r.Xy)(n+1)}", "Xy", "fetchIt")
	other := hashOfBundle(t, "function b(t){return(0,o.Qm)(t+1)}", "Qm", "fetchIt")
	if one != other {
		t.Errorf("the hash is %s and %s with other names of the bundler, want the same", one, other)
	}
}

// @implements TC_187@1.1.1
func TestTheHashSeesAnotherCallAndAnotherLiteral(t *testing.T) {

	fetchIt := hashOfBundle(t, "function a(n){return(0,r.Xy)(n+1)}", "Xy", "fetchIt")
	if other := hashOfBundle(t, "function a(n){return(0,r.Xy)(n+1)}", "Xy", "fetchOther"); other == fetchIt {
		t.Errorf("the call of another function of another module has the same hash, %s", other)
	}
	if two := hashOfBundle(t, "function a(n){return(0,r.Xy)(n+2)}", "Xy", "fetchIt"); two == fetchIt {
		t.Errorf("another literal has the same hash, %s", two)
	}
}

// readIstanbul reads a coverage file of istanbul of the source src/a.js of the project web, written
// below a folder of the test; the coverage names the source by ${path}.
func readIstanbul(t *testing.T, src, coverage string) *Istanbul {

	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, "web", "src", "a.js")
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	name, err := json.Marshal(path)
	if err != nil {
		t.Fatal(err)
	}
	c, err := ReadIstanbul(strings.NewReader(strings.ReplaceAll(coverage, `"${path}"`, string(name))), root)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// loader is the source of the async method load of the class Loader, lines 2 to 4, and its coverage.
const loader = `export class Loader {
  async load(n) {
    return g(n);
  }
}
`

const loaderCoverage = `{"${path}": {"path": "${path}",
	"statementMap": {"0": {"start": {"line": 3, "column": 4}, "end": {"line": 3, "column": 16}}},
	"fnMap": {"0": {"name": "load", "decl": {"start": {"line": 2, "column": 8}, "end": {"line": 2, "column": 12}},
		"loc": {"start": {"line": 2, "column": 2}, "end": {"line": 4, "column": 3}}}},
	"s": {"0": 1}, "f": {"0": 1}}}`

// split is a bundle in which the async method load of the class Loader became a wrapper and a
// function _callee, both mapped into it; the body of each maps to line 2 or 3 of the source.
func split(t *testing.T, callee string) (*Code, Function) {

	t.Helper()
	fs := readIstanbul(t, loader, loaderCoverage).Functions()
	if len(fs) != 1 {
		t.Fatalf("the functions are %+v, want load alone", fs)
	}
	bundle := "function w(){return c.apply(this,arguments)}function c(n){" + callee + "}"
	segments := []segment{
		{column: at(t, bundle, "return c"), line: 1, sourceColumn: 2},
		{column: at(t, bundle, "function c"), line: 1, sourceColumn: 2},
		{column: at(t, bundle, callee), line: 2, sourceColumn: 4},
	}
	code, err := ReadBundle([]byte(bundle), sourceMapOf(t, segments), "web", fs)
	if err != nil {
		t.Fatal(err)
	}
	return code, fs[0]
}

// @implements TC_182@2.1.1
func TestAMethodSplitByTheBundlerIsOneRow(t *testing.T) {

	code, load := split(t, "return g(n)")
	if parts := code.byFunc[load.Key+"\x002"]; len(parts) != 2 {
		t.Fatalf("load has %d functions of the bundle, want the wrapper and the _callee", len(parts))
	}
	rows := Rows([]Function{load}, code, &Istanbul{coverage: map[string][3]int{}})
	if len(rows) != 1 || rows[0].Key != "web/src/a.js@Loader.load" || !rows[0].InOutput || rows[0].Name != "w" {
		t.Errorf("the rows are %+v, want one, web/src/a.js@Loader.load, in the bundle as w", rows)
	}
}

// @implements TC_188@1.1.1
func TestTheHashOfASplitMethodSeesItsCallee(t *testing.T) {

	one, load := split(t, "return g(n)")
	other, _ := split(t, "return g(n,1)")
	_, h1, _ := one.Find(load)
	_, h2, _ := other.Find(load)
	if h1 == h2 {
		t.Errorf("a change in the _callee only leaves the hash of load, %s", h1)
	}
}

// project writes node_modules with webpack, @babel/core and terser, and a configuration of webpack
// with the option of terser given.
func project(t *testing.T, option string) (modules, config string) {

	t.Helper()
	dir := t.TempDir()
	for name, version := range map[string]string{"webpack": "5.111.1", "@babel/core": "7.29.7", "terser": "5.51.2"} {
		p := filepath.Join(dir, "node_modules", filepath.FromSlash(name))
		if err := os.MkdirAll(p, 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(p, "package.json"), []byte(`{"name":"`+name+`","version":"`+version+`"}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	config = filepath.Join(dir, "webpack.config.js")
	if err := os.WriteFile(config, []byte("module.exports = {optimization: {minimizer: [{terserOptions: {"+option+"}}]}}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, "node_modules"), config
}

// head is the head of the section of a bundle built with the option of terser given.
func head(t *testing.T, option string) Bundle {

	t.Helper()
	modules, config := project(t, option)
	tools, err := ReadTools(modules, []string{"webpack", "@babel/core", "terser"})
	if err != nil {
		t.Fatal(err)
	}
	hash, err := ConfigurationHash([]string{config})
	if err != nil {
		t.Fatal(err)
	}
	return Bundle{Name: "maturity.js", Dir: "web", Node: "v24.21.0", Tools: tools, Configuration: hash, Commit: "e543ec4"}
}

// @implements TC_179@1.1.1
func TestTheSectionNamesTheBundle(t *testing.T) {

	b := head(t, "mangle: true")
	lines := b.Lines()
	want := "node v24.21.0, webpack 5.111.1, @babel/core 7.29.7, terser 5.51.2, configuration " + b.Configuration + ", commit e543ec4"
	if len(lines) != 2 || lines[0] != "Bundle: maturity.js (web)" || lines[1] != want {
		t.Errorf("the head is %q, want the bundle, then %q", lines, want)
	}
	if len(b.Configuration) != 8 {
		t.Errorf("the hash of the configuration is %q", b.Configuration)
	}
}

// @implements TC_180@1.1.1
func TestTheHashOfTheConfigurationSeesAnOption(t *testing.T) {

	one, other := head(t, "mangle: true"), head(t, "mangle: false")
	if one.Configuration == other.Configuration {
		t.Errorf("another option of terser has the same hash of the configuration, %s", one.Configuration)
	}
	if strings.Join(one.Lines()[1:], "") == strings.Join(other.Lines()[1:], "") {
		t.Error("the heads are the same")
	}
	for i := range one.Tools {
		if one.Tools[i] != other.Tools[i] {
			t.Errorf("the versions differ: %v and %v", one.Tools[i], other.Tools[i])
		}
	}
}

func TestReadIstanbul(t *testing.T) {

	coverage := `{"${path}": {"path": "${path}",
		"statementMap": {"0": {"start": {"line": 2, "column": 2}, "end": {"line": 2, "column": 20}},
			"1": {"start": {"line": 6, "column": 2}, "end": {"line": 6, "column": 9}}},
		"fnMap": {"0": {"name": "load", "decl": {"start": {"line": 1, "column": 16}, "end": {"line": 1, "column": 20}},
				"loc": {"start": {"line": 1, "column": 25}, "end": {"line": 3, "column": null}}},
			"1": {"name": "(anonymous_1)", "decl": {"start": {"line": 5, "column": 0}, "end": {"line": 5, "column": 1}},
				"loc": {"start": {"line": 5, "column": 5}, "end": {"line": 7, "column": null}}}},
		"s": {"0": 3, "1": 0}, "f": {"0": 1, "1": 0}}}`
	src := "export function load(n) {\n  return fetchIt(n + 1);\n}\n\nf(() => {\n  return 1;\n});\n"
	c := readIstanbul(t, src, coverage)
	fs := c.Functions()
	if len(fs) != 1 || fs[0].Key != "web/src/a.js@load" || fs[0].From != 1 || fs[0].FromColumn != 16 || fs[0].ToColumn != endOfLine {
		t.Fatalf("the functions are %+v, want load alone, from its head to the end of line 3", fs)
	}
	if called, ran, total := c.Coverage(fs[0]); !called || ran != 1 || total != 1 {
		t.Errorf("load: %v %d/%d, want called, 1/1", called, ran, total)
	}
}

// twoClasses is a source with JSX of two classes, each with a method render: A's on line 4, B's
// on line 10.
const twoClasses = `import React from 'react';

export class A extends React.Component {
  render() {
    return <p>Don't {this.props.x}</p>;
  }
}

class B {
  async render(n) {
    return ` + "`${n} }`" + `;
  }
}
`

// @implements TC_190@1.1.1
func TestMethodsOfTheSameNameInTwoClassesAreTwoRows(t *testing.T) {

	coverage := `{"${path}": {"path": "${path}",
		"statementMap": {"0": {"start": {"line": 5, "column": 4}, "end": {"line": 5, "column": 39}},
			"1": {"start": {"line": 11, "column": 4}, "end": {"line": 11, "column": 20}}},
		"fnMap": {"0": {"name": "render", "decl": {"start": {"line": 4, "column": 2}, "end": {"line": 4, "column": 8}},
				"loc": {"start": {"line": 4, "column": 2}, "end": {"line": 6, "column": 3}}},
			"1": {"name": "render", "decl": {"start": {"line": 10, "column": 8}, "end": {"line": 10, "column": 14}},
				"loc": {"start": {"line": 10, "column": 2}, "end": {"line": 12, "column": 3}}}},
		"s": {"0": 2, "1": 0}, "f": {"0": 2, "1": 0}}}`
	c := readIstanbul(t, twoClasses, coverage)
	fs := c.Functions()
	if len(fs) != 2 || fs[0].Key != "web/src/a.js@A.render" || fs[1].Key != "web/src/a.js@B.render" {
		t.Fatalf("the functions are %+v, want web/src/a.js@A.render and web/src/a.js@B.render", fs)
	}
	if called, ran, total := c.Coverage(fs[0]); !called || ran != 1 || total != 1 {
		t.Errorf("A.render: %v %d/%d, want called, 1/1", called, ran, total)
	}
	if called, ran, total := c.Coverage(fs[1]); called || ran != 0 || total != 1 {
		t.Errorf("B.render: %v %d/%d, want not called, 0/1", called, ran, total)
	}
}

func TestClassesOfTypeScript(t *testing.T) {

	classes, err := Classes("a.ts", []byte("class C<T> {\n  private n: number = 1;\n  get(x: T): T { return x; }\n}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := classes[position{3, 2}]; got != "C" {
		t.Errorf("the classes are %v, want C at line 3, column 2", classes)
	}
}

func TestSourceFile(t *testing.T) {

	for in, want := range map[string]string{
		"webpack://viewer/./src/a.js": "web/src/a.js",
		"./src/a.js":                  "web/src/a.js",
		"src/a.js":                    "web/src/a.js",
	} {
		if got := sourceFile(in, "web"); got != want {
			t.Errorf("sourceFile(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAFunctionInsideIsPartOfItsFunction(t *testing.T) {

	bundle := "function a(n){return n.map(function(t){return t+1})}"
	segments := []segment{
		{column: at(t, bundle, "return n"), line: 1, sourceColumn: 2},
		{column: at(t, bundle, "return t"), line: 1, sourceColumn: 20},
	}
	code, err := ReadBundle([]byte(bundle), sourceMapOf(t, segments), "web", []Function{load})
	if err != nil {
		t.Fatal(err)
	}
	if parts := code.byFunc[load.Key+"\x001"]; len(parts) != 1 {
		t.Errorf("load has %d functions of the bundle, want the outer one alone", len(parts))
	}
}
