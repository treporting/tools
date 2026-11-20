package write

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
	"github.com/go-sourcemap/sourcemap"
)

// Bundle is what the head of the section of a JavaScript bundle names.
type Bundle struct {
	// Name is the file of the bundle, Dir the folder of the project that builds it, from the root of
	// the repo.
	Name, Dir string
	// Node is the version of node, such as v24.21.0; Tools are the tools that built the bundle, each
	// with its version.
	Node  string
	Tools []Tool
	// Configuration is the hash of the configuration of the tools.
	Configuration string
	// Commit is the commit it was built from.
	Commit string
}

// Tool is a tool that built a bundle, such as webpack, with its version.
type Tool struct{ Name, Version string }

// ReadTools reads the version of each tool from its package.json below modules, the node_modules of
// the project.
//
// @implements SR_116@1.1.1
func ReadTools(modules string, names []string) ([]Tool, error) {

	var tools []Tool
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(modules, filepath.FromSlash(name), "package.json")) // #nosec G304 -- the path comes from the user of the tool
		if err != nil {
			return nil, err
		}
		var p struct{ Version string }
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		tools = append(tools, Tool{Name: name, Version: p.Version})
	}
	return tools, nil
}

// ConfigurationHash is the hash over the files of the configuration of the tools, each by its name
// and its content, in the order given.
//
// @implements SR_116@1.1.1
func ConfigurationHash(files []string) (string, error) {

	h := sha256.New()
	for _, f := range files {
		data, err := os.ReadFile(f) // #nosec G304 -- the path comes from the user of the tool
		if err != nil {
			return "", err
		}
		_, _ = fmt.Fprintf(h, "%s\n%d\n", filepath.Base(f), len(data))
		_, _ = h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil))[:8], nil
}

// position is a place in a file: its line from 1, its column from 0.
type position struct{ line, column int }

func (p position) before(q position) bool {
	return p.line < q.line || p.line == q.line && p.column < q.column
}

// istanbul is a coverage file of istanbul, coverage-final.json, as Vitest writes it with the reporter
// json: per file of the source, its statements and its functions, each with its count.
type istanbul map[string]struct {
	Path         string
	StatementMap map[string]span
	FnMap        map[string]struct {
		Name string
		Decl span
		Loc  span
	}
	S map[string]int
	F map[string]int
}

// span is a range of a coverage file of istanbul; a column it leaves out, null, is the end of its line.
type span struct {
	Start, End struct {
		Line   int
		Column *int
	}
}

// endOfLine is the column of the end of any line.
const endOfLine = 1 << 30

func columnOf(c *int) int {
	if c == nil {
		return endOfLine
	}
	return *c
}

func (s span) from() position { return position{s.Start.Line, columnOf(s.Start.Column)} }
func (s span) to() position   { return position{s.End.Line, columnOf(s.End.Column)} }

// Istanbul is the coverage of the tests of a JavaScript project, and the functions of its source.
type Istanbul struct {
	functions []Function
	coverage  map[string][3]int
}

// ReadIstanbul reads a coverage file of istanbul. Its functions are those of the source with a name,
// keyed by their file from root, a method also by its class, which istanbul does not name: from the
// file of the source, as Classes finds it; a function of the source that has none, such as a
// callback, is part of the function it lies in.
//
// @implements SR_117@2.1.1
func ReadIstanbul(r io.Reader, root string) (*Istanbul, error) {

	var files istanbul
	if err := json.NewDecoder(r).Decode(&files); err != nil {
		return nil, fmt.Errorf("a coverage file of istanbul: %w", err)
	}
	c := &Istanbul{coverage: map[string][3]int{}}
	for _, file := range files {
		rel, err := filepath.Rel(root, file.Path)
		if err != nil {
			return nil, err
		}
		rel = filepath.ToSlash(rel)
		var classes map[position]string
		for id, fn := range file.FnMap {
			if fn.Name == "" || strings.HasPrefix(fn.Name, "(anonymous") {
				continue
			}
			if classes == nil {
				src, err := os.ReadFile(file.Path) // #nosec G304 -- the file of the source the coverage names
				if err != nil {
					return nil, fmt.Errorf("the source of the coverage: %w", err)
				}
				if classes, err = Classes(file.Path, src); err != nil {
					return nil, err
				}
			}
			key := fn.Name
			if class, ok := classes[fn.Decl.from()]; ok {
				key = class + "." + fn.Name
			}
			//	a function begins with its head, decl, before the body, loc
			from, to := fn.Loc.from(), fn.Loc.to()
			if fn.Decl.Start.Line > 0 && fn.Decl.from().before(from) {
				from = fn.Decl.from()
			}
			f := Function{
				Key: rel + "@" + key, Symbol: fn.Name, File: rel,
				From: from.line, To: to.line, FromColumn: from.column, ToColumn: to.column,
			}
			ran, total := 0, 0
			for sid, s := range file.StatementMap {
				if s.from().before(from) || to.before(s.to()) {
					continue
				}
				total++
				if file.S[sid] > 0 {
					ran++
				}
			}
			called := 0
			if file.F[id] > 0 {
				called = 1
			}
			c.functions = append(c.functions, f)
			c.coverage[f.Key+"\x00"+strconv.Itoa(f.From)] = [3]int{called, ran, total}
		}
	}
	sort.SliceStable(c.functions, func(i, j int) bool {
		a, b := c.functions[i], c.functions[j]
		if a.Key != b.Key {
			return a.Key < b.Key
		}
		return a.From < b.From
	})
	return c, nil
}

// Functions are the functions of the source, by their keys.
func (c *Istanbul) Functions() []Function { return c.functions }

// Coverage is whether the tests call a function, and how many of its statements they ran.
//
// @implements SR_117@2.1.1
func (c *Istanbul) Coverage(f Function) (called bool, ran, statements int) {
	v := c.coverage[f.Key+"\x00"+strconv.Itoa(f.From)]
	return v[0] == 1, v[1], v[2]
}

// Code is the code of each function of the source in a bundle: the functions of the bundle its
// source map maps to it.
type Code struct {
	src    string
	byFunc map[string][]*bundled
	names  []identifier
	smap   *sourcemap.Consumer
	utf16  []int
	starts []int
}

// bundled is a function of a bundle.
type bundled struct {
	from, to int
	name     string
}

// identifier is a name in a bundle; property is whether it names a property, as after a dot, rather
// than a variable.
type identifier struct {
	at       int
	name     string
	property bool
}

// ReadBundle reads a bundle and its source map, and finds for each function of the source the
// functions of the bundle that are its code. dir is the folder of the project that built it, from
// root, to which the source map names its sources.
//
// @implements SR_117@2.1.1, SR_119@1.1.1
func ReadBundle(bundle, sourceMap []byte, dir string, functions []Function) (*Code, error) {

	program, err := parser.ParseFile(nil, "", string(bundle), 0, parser.WithDisableSourceMaps)
	if err != nil {
		return nil, fmt.Errorf("the bundle: %w", err)
	}
	smap, err := sourcemap.Parse("", sourceMap)
	if err != nil {
		return nil, fmt.Errorf("the source map: %w", err)
	}
	c := &Code{src: string(bundle), byFunc: map[string][]*bundled{}, smap: smap}
	c.index()

	byFile := map[string][]Function{}
	for _, f := range functions {
		byFile[f.File] = append(byFile[f.File], f)
	}
	//	the tree names some nodes twice, such as the declarations of a function
	type seen struct {
		at uintptr
		t  reflect.Type
	}
	visited := map[seen]bool{}
	var visit func(v reflect.Value, owner string, property bool)
	visit = func(v reflect.Value, owner string, property bool) {
		if !v.IsValid() {
			return
		}
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if v.IsNil() {
				return
			}
			visit(v.Elem(), owner, property)
			return
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				visit(v.Index(i), owner, false)
			}
			return
		case reflect.Struct:
			if !v.CanAddr() {
				return
			}
		default:
			return
		}
		key := seen{v.Addr().Pointer(), v.Type()}
		if visited[key] {
			return
		}
		visited[key] = true
		switch n := v.Addr().Interface().(type) {
		case *ast.Identifier:
			c.names = append(c.names, identifier{at: int(n.Idx) - 1, name: string(n.Name), property: property})
			return
		case *ast.DotExpression:
			visit(reflect.ValueOf(n.Left), owner, false)
			visit(reflect.ValueOf(&n.Identifier), owner, true)
			return
		case *ast.PropertyKeyed:
			visit(reflect.ValueOf(n.Key), owner, !n.Computed)
			visit(reflect.ValueOf(n.Value), owner, false)
			return
		case *ast.MethodDefinition:
			visit(reflect.ValueOf(n.Key), owner, !n.Computed)
			visit(reflect.ValueOf(n.Body), owner, false)
			return
		case *ast.FunctionLiteral, *ast.ArrowFunctionLiteral:
			node := n.(ast.Node)
			b := &bundled{from: int(node.Idx0()) - 1, to: int(node.Idx1()) - 1}
			if f, ok := n.(*ast.FunctionLiteral); ok && f.Name != nil {
				b.name = string(f.Name.Name)
			}
			if key, ok := c.functionOf(anchorOf(n.(ast.Node)), byFile, dir); ok {
				//	the code of a function of the bundle lies in the function around it where both
				//	are the code of the same function of the source
				if key != owner {
					c.byFunc[key] = append(c.byFunc[key], b)
				}
				owner = key
			}
		}
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				visit(v.Field(i), owner, false)
			}
		}
	}
	visit(reflect.ValueOf(program), "", false)
	sort.Slice(c.names, func(i, j int) bool { return c.names[i].at < c.names[j].at })
	return c, nil
}

// anchorOf is the place in a function of the bundle its source map maps most surely into the
// function of the source: its first statement. The keyword function maps to what stands before the
// function in the source, such as export.
func anchorOf(n ast.Node) int {

	var body ast.Node
	switch f := n.(type) {
	case *ast.FunctionLiteral:
		if f.Body != nil {
			if len(f.Body.List) > 0 {
				return int(f.Body.List[0].Idx0()) - 1
			}
			body = f.Body
		}
	case *ast.ArrowFunctionLiteral:
		switch b := f.Body.(type) {
		case *ast.BlockStatement:
			if len(b.List) > 0 {
				return int(b.List[0].Idx0()) - 1
			}
			body = b
		case *ast.ExpressionBody:
			return int(b.Expression.Idx0()) - 1
		}
	}
	if body != nil {
		return int(body.Idx0()) - 1
	}
	return int(n.Idx0()) - 1
}

// index notes where each line of the bundle starts, and the column of each byte in UTF-16 units, as
// a source map counts it.
func (c *Code) index() {

	c.starts = []int{0}
	c.utf16 = make([]int, len(c.src)+1)
	column := 0
	for i := 0; i < len(c.src); {
		c.utf16[i] = column
		r, size := utf8.DecodeRuneInString(c.src[i:])
		for k := 1; k < size; k++ {
			c.utf16[i+k] = column
		}
		i += size
		if r == '\n' {
			c.starts = append(c.starts, i)
			column = 0
			continue
		}
		column++
		if r >= 0x10000 {
			column++
		}
	}
	c.utf16[len(c.src)] = column
}

// source is the place in the source a place in the bundle comes from, and the name it has there.
func (c *Code) source(at int) (file string, p position, name string, ok bool) {

	line := sort.Search(len(c.starts), func(i int) bool { return c.starts[i] > at })
	file, name, l, col, ok := c.smap.Source(line, c.utf16[at])
	return file, position{l, col}, name, ok
}

// sourceFile is the file a source map names, from the root of the repo: webpack names it
// webpack://<project>/./src/a.js, from the folder of the project.
func sourceFile(name, dir string) string {

	if rest, ok := strings.CutPrefix(name, "webpack://"); ok {
		if _, after, found := strings.Cut(rest, "/"); found {
			name = after
		}
	}
	return path.Join(dir, path.Clean(strings.TrimPrefix(name, "./")))
}

// functionOf is the innermost function of the source with a name that a place in the bundle maps into.
func (c *Code) functionOf(at int, byFile map[string][]Function, dir string) (string, bool) {

	name, p, _, ok := c.source(at)
	if !ok {
		return "", false
	}
	var inner *Function
	for i, f := range byFile[sourceFile(name, dir)] {
		from, to := position{f.From, f.FromColumn}, position{f.To, f.ToColumn}
		if p.before(from) || to.before(p) {
			continue
		}
		if inner == nil || (position{inner.From, inner.FromColumn}).before(from) {
			inner = &byFile[sourceFile(name, dir)][i]
		}
	}
	if inner == nil {
		return "", false
	}
	return inner.Key + "\x00" + strconv.Itoa(inner.From), true
}

// Find is the name of a function of the source in the bundle, and the hash of its code there: with
// every name the bundler chose replaced by its order of appearance, a property - as the call of a
// function of another module - by its name in the source as the source map gives it, and statements,
// operators and literals as they are; over all the functions of the bundle that are its code.
//
// @implements SR_119@1.1.1
func (c *Code) Find(f Function) (name, hash string, ok bool) {

	parts := c.byFunc[f.Key+"\x00"+strconv.Itoa(f.From)]
	if len(parts) == 0 {
		return "", "", false
	}
	name = "-"
	if parts[0].name != "" {
		name = parts[0].name
	}
	h := sha256.New()
	for _, b := range parts {
		_, _ = io.WriteString(h, c.normalized(b))
		_, _ = io.WriteString(h, "\n")
	}
	return name, hex.EncodeToString(h.Sum(nil))[:8], true
}

// normalized is the text of a function of the bundle with its names replaced.
func (c *Code) normalized(b *bundled) string {

	var s strings.Builder
	order := map[string]int{}
	at := b.from
	first := sort.Search(len(c.names), func(i int) bool { return c.names[i].at >= b.from })
	for _, id := range c.names[first:] {
		if id.at >= b.to {
			break
		}
		if !strings.HasPrefix(c.src[id.at:], id.name) {
			//	a name written with an escape: left as it is
			continue
		}
		s.WriteString(c.src[at:id.at])
		at = id.at + len(id.name)
		if id.property {
			if _, _, original, ok := c.source(id.at); ok && original != "" {
				s.WriteString(original)
				continue
			}
			s.WriteString(id.name)
			continue
		}
		n, seen := order[id.name]
		if !seen {
			n = len(order) + 1
			order[id.name] = n
		}
		s.WriteString("$" + strconv.Itoa(n))
	}
	s.WriteString(c.src[at:b.to])
	return strings.Join(strings.Fields(s.String()), " ")
}
