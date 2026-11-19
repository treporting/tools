package write

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Function is a function of the source code.
type Function struct {
	// Key is folder/file@function, the file from the root of the repo, as the source code (SC)
	// names it; a method by its name alone.
	Key string
	// Symbol is its name in the binary, such as example.org/m/a.(*T).M.
	Symbol string
	// File is the file as a coverage profile names it, the import path of its package and its name.
	File string
	// From and To are its first and its last line.
	From, To int
}

// Package is a package of the source code, as `go list -json` prints it.
type Package struct {
	ImportPath string
	Name       string
	Dir        string
	GoFiles    []string
	Module     *struct{ Path string }
}

// ListPackages lists the packages of the module that the main package at main is built from, by
// `go list -deps -json` in dir.
func ListPackages(dir, main string) ([]Package, error) {

	cmd := exec.Command("go", "list", "-deps", "-json", main) // #nosec G204 -- main comes from the user of the tool
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list %s: %w", main, err)
	}
	return readPackages(strings.NewReader(string(out)))
}

func readPackages(r io.Reader) ([]Package, error) {

	var all []Package
	decoder := json.NewDecoder(r)
	var module string
	for decoder.More() {
		var p Package
		if err := decoder.Decode(&p); err != nil {
			return nil, err
		}
		all = append(all, p)
	}
	//	the main package is listed last; its module is the module of the source
	if len(all) > 0 && all[len(all)-1].Module != nil {
		module = all[len(all)-1].Module.Path
	}
	var own []Package
	for _, p := range all {
		if p.Module != nil && p.Module.Path == module {
			own = append(own, p)
		}
	}
	return own, nil
}

// Functions lists the functions of the packages, with their keys from root.
//
// @implements SR_117@1.1.1
func Functions(root string, packages []Package) ([]Function, error) {

	var functions []Function
	fset := token.NewFileSet()
	for _, p := range packages {
		//	the compiler names the symbols of the main package main, not by its import path
		prefix := p.ImportPath
		if p.Name == "main" {
			prefix = "main"
		}
		for _, name := range p.GoFiles {
			path := filepath.Join(p.Dir, name)
			file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
			if err != nil {
				return nil, err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil, err
			}
			for _, decl := range file.Decls {
				f, ok := decl.(*ast.FuncDecl)
				if !ok || f.Body == nil {
					continue
				}
				functions = append(functions, Function{
					Key:    filepath.ToSlash(rel) + "@" + f.Name.Name,
					Symbol: symbolOf(prefix, f),
					File:   p.ImportPath + "/" + name,
					From:   fset.Position(f.Pos()).Line,
					To:     fset.Position(f.End()).Line,
				})
			}
		}
	}
	sort.SliceStable(functions, func(i, j int) bool { return functions[i].Key < functions[j].Key })
	return functions, nil
}

// symbolOf is the name the compiler gives a function: pkg.F, pkg.T.M or pkg.(*T).M.
func symbolOf(prefix string, f *ast.FuncDecl) string {

	if f.Recv == nil || len(f.Recv.List) == 0 {
		return prefix + "." + f.Name.Name
	}
	t := f.Recv.List[0].Type
	pointer := false
	if star, ok := t.(*ast.StarExpr); ok {
		pointer, t = true, star.X
	}
	//	a generic receiver, T[K]
	switch g := t.(type) {
	case *ast.IndexExpr:
		t = g.X
	case *ast.IndexListExpr:
		t = g.X
	}
	name := fmt.Sprint(t)
	if id, ok := t.(*ast.Ident); ok {
		name = id.Name
	}
	if pointer {
		return prefix + ".(*" + name + ")." + f.Name.Name
	}
	return prefix + "." + name + "." + f.Name.Name
}
