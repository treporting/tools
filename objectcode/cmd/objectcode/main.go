// Command objectcode writes the section of the object code report of a build output (TOOL_35): its
// head, and a row per function of the source code it is built from.
//
//	objectcode -binary file -main ./path/of/main [-coverprofile file] [-commit sha] [-root dir]
//	objectcode -bundle file -dir folder -coverage coverage-final.json [-sourcemap file]
//	           [-tool name]... [-config file]... [-modules dir] [-node version] [-commit sha] [-root dir]
//
// For a Go binary it lists the functions of the packages of the module the main package is built
// from, reads their machine code from the binary by `go tool objdump`, and their coverage from the
// coverage profile of the run of the tests, as `go test -coverprofile` writes it.
//
// For a JavaScript bundle it takes the functions of the source and their coverage from the coverage
// file of istanbul, as Vitest writes it with the reporter json, and finds the code of each in the
// bundle by its source map, <bundle>.map where -sourcemap is empty. -dir is the folder of the project
// that built the bundle, from root. It names each -tool with the version of its package.json in
// -modules, <dir>/node_modules where empty, and the hash of the files of -config.
//
// It writes the section to the standard output.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/treporting/tools/objectcode/write"
)

// list collects a flag that may be given more than once.
type list []string

func (l *list) String() string     { return strings.Join(*l, ",") }
func (l *list) Set(v string) error { *l = append(*l, v); return nil }

// @implements SR_115@1.1.1, SR_116@1.1.1, SR_117@1.1.1, SR_118@1.1.1, SR_119@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	commit := flag.String("commit", "", "the commit the output was built from, where it names none")
	binary := flag.String("binary", "", "the Go binary")
	mainPkg := flag.String("main", "", "the main package of the binary, from root, such as ./annotation/cmd/annotationcheck")
	coverprofile := flag.String("coverprofile", "", "the coverage profile of the run of the Go tests; none where empty")
	bundle := flag.String("bundle", "", "the JavaScript bundle")
	sourceMap := flag.String("sourcemap", "", "the source map of the bundle; <bundle>.map where empty")
	dir := flag.String("dir", "", "the folder of the project that built the bundle, from root")
	coverage := flag.String("coverage", "", "the coverage file of istanbul of the run of the JavaScript tests")
	modules := flag.String("modules", "", "the node_modules of the project; <dir>/node_modules where empty")
	node := flag.String("node", "", "the version of node; that of node on the PATH where empty")
	var tools, configs list
	flag.Var(&tools, "tool", "a tool that built the bundle, by the name of its package, such as webpack; may be given more than once")
	flag.Var(&configs, "config", "a file of the configuration of the tools; may be given more than once")
	flag.Parse()

	abs, err := filepath.Abs(*root)
	if err != nil {
		fail(err)
	}
	switch {
	case *binary != "" && *mainPkg != "" && *bundle == "":
		fmt.Print(goBinary(abs, *binary, *mainPkg, *coverprofile, *commit))
	case *bundle != "" && *dir != "" && *coverage != "" && *binary == "":
		fmt.Print(jsBundle(abs, *bundle, *sourceMap, *dir, *coverage, *modules, *node, *commit, tools, configs))
	default:
		fail(fmt.Errorf("give -binary and -main, or -bundle, -dir and -coverage"))
	}
}

// goBinary is the section of a Go binary.
func goBinary(root, binary, mainPkg, coverprofile, commit string) string {

	packages, err := write.ListPackages(root, mainPkg)
	if err != nil {
		fail(err)
	}
	functions, err := write.Functions(root, packages)
	if err != nil {
		fail(err)
	}
	d, err := write.Objdump(binary, packages)
	if err != nil {
		fail(err)
	}
	var profile *write.Profile
	if coverprofile != "" {
		f, err := os.Open(coverprofile) // #nosec G304 -- the path comes from the user of the tool
		if err != nil {
			fail(err)
		}
		profile, err = write.ReadProfile(f)
		_ = f.Close()
		if err != nil {
			fail(fmt.Errorf("%s: %w", coverprofile, err))
		}
	}
	b, err := write.ReadBinary(binary, commit)
	if err != nil {
		fail(err)
	}
	b.Name = filepath.Base(binary)
	b.Main = strings.TrimPrefix(filepath.ToSlash(mainPkg), "./")
	return write.Section(b, write.Rows(functions, d, profile))
}

// jsBundle is the section of a JavaScript bundle.
func jsBundle(root, bundle, sourceMap, dir, coverage, modules, node, commit string, tools, configs []string) string {

	f, err := os.Open(coverage) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	c, err := write.ReadIstanbul(f, root)
	_ = f.Close()
	if err != nil {
		fail(fmt.Errorf("%s: %w", coverage, err))
	}
	if sourceMap == "" {
		sourceMap = bundle + ".map"
	}
	src, err := os.ReadFile(bundle) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	smap, err := os.ReadFile(sourceMap) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	dir = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(dir)), "./")
	code, err := write.ReadBundle(src, smap, dir, c.Functions())
	if err != nil {
		fail(err)
	}
	if modules == "" {
		modules = filepath.Join(root, dir, "node_modules")
	}
	b := write.Bundle{Name: filepath.Base(bundle), Dir: dir, Node: node, Commit: commit}
	if b.Tools, err = write.ReadTools(modules, tools); err != nil {
		fail(err)
	}
	if b.Configuration, err = write.ConfigurationHash(configs); err != nil {
		fail(err)
	}
	if b.Node == "" {
		out, err := exec.Command("node", "--version").Output()
		if err != nil {
			fail(fmt.Errorf("node --version: %w", err))
		}
		b.Node = strings.TrimSpace(string(out))
	}
	if b.Commit == "" {
		out, err := exec.Command("git", "-C", root, "rev-parse", "--short=7", "HEAD").Output() // #nosec G204 -- root comes from the user of the tool
		if err != nil {
			fail(fmt.Errorf("the commit: %w", err))
		}
		b.Commit = strings.TrimSpace(string(out))
	}
	return write.Section(b, write.Rows(c.Functions(), code, c))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "objectcode:", err)
	os.Exit(1)
}
