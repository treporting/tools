package write

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/go-sourcemap/sourcemap"
)

// loaders are the loaders of esbuild by the extension of a file of the source; any other file is
// JavaScript, with or without JSX.
var loaders = map[string]api.Loader{".ts": api.LoaderTS, ".mts": api.LoaderTS, ".cts": api.LoaderTS, ".tsx": api.LoaderTSX}

// Classes answers the class of each method of a file of the JavaScript source, by the place of the
// name of the method in the file. esbuild turns the file, with JSX or TypeScript, into JavaScript
// that goja parses, with a source map back to the file. A class without a name has no entry.
//
// @implements SR_117@2.1.1
func Classes(name string, src []byte) (map[position]string, error) {

	loader, ok := loaders[filepath.Ext(name)]
	if !ok {
		loader = api.LoaderJSX
	}
	out := api.Transform(string(src), api.TransformOptions{
		Loader: loader, Format: api.FormatCommonJS, Target: api.ES2020, Sourcemap: api.SourceMapExternal, Sourcefile: name,
	})
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("%s: %s", name, out.Errors[0].Text)
	}
	program, err := parser.ParseFile(nil, name, string(out.Code), 0, parser.WithDisableSourceMaps)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	smap, err := sourcemap.Parse("", out.Map)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	c := &Code{src: string(out.Code), smap: smap}
	c.index()

	classes := map[position]string{}
	//	the tree names some nodes twice, such as the declarations of a class
	type seen struct {
		at uintptr
		t  reflect.Type
	}
	visited := map[seen]bool{}
	var visit func(v reflect.Value)
	visit = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if !v.IsNil() {
				visit(v.Elem())
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				visit(v.Index(i))
			}
		case reflect.Struct:
			if !v.CanAddr() || visited[seen{v.Addr().Pointer(), v.Type()}] {
				return
			}
			visited[seen{v.Addr().Pointer(), v.Type()}] = true
			if class, ok := v.Addr().Interface().(*ast.ClassLiteral); ok && class.Name != nil {
				for _, e := range class.Body {
					if m, ok := e.(*ast.MethodDefinition); ok && !m.Computed {
						if _, p, _, ok := c.source(int(m.Key.Idx0()) - 1); ok {
							classes[p] = string(class.Name.Name)
						}
					}
				}
			}
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() {
					visit(v.Field(i))
				}
			}
		}
	}
	visit(reflect.ValueOf(program))
	return classes, nil
}
