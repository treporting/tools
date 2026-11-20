// Command objectcode writes the section of the object code report of a build output (TOOL_35): its
// head, and a row per function of the source code it is built from.
//
//	objectcode -binary file -main ./path/of/main [-coverprofile file] [-commit sha] [-root dir]
//	objectcode -bundle file -dir folder -coverage coverage-final.json [-sourcemap file]
//	           [-tool name]... [-config file]... [-modules dir] [-node version] [-commit sha] [-root dir]
//	objectcode -document [-label x.y.z -ticket number] [-root dir] section.md...
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
// Where the coverage file or the source map is not there, the section of the bundle is its head and a
// note that says which. It writes the section to the standard output.
//
// With -document it writes the object code report as testreports/object-code.md in root from the
// sections given, in their order, with the label and the ticket, or the ticket alone. It keeps the rows of the Change
// history of a report that is there, and adds one. It sets the entry of the report in documents.yaml
// of the repo: the version of the new row, the hash of the file written, completed and generated.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/treporting/tools/objectcode/write"
	"github.com/treporting/tools/sbom/create"
	report "github.com/treporting/tools/testreport/write"
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
	document := flag.Bool("document", false, "write testreports/object-code.md from the sections given as arguments")
	label := flag.String("label", "", "the label of the baseline, such as 1.2.0")
	ticket := flag.String("ticket", "", "the number of the ticket of the baseline, such as 438")
	flag.Parse()

	abs, err := filepath.Abs(*root)
	if err != nil {
		fail(err)
	}
	switch {
	case *document && *binary == "" && *bundle == "":
		writeDocument(abs, *label, *ticket, flag.Args())
	case *binary != "" && *mainPkg != "" && *bundle == "":
		fmt.Print(goBinary(abs, *binary, *mainPkg, *coverprofile, *commit))
	case *bundle != "" && *dir != "" && *coverage != "" && *binary == "":
		fmt.Print(jsBundle(abs, *bundle, *sourceMap, *dir, *coverage, *modules, *node, *commit, tools, configs))
	default:
		fail(fmt.Errorf("give -binary and -main, or -bundle, -dir and -coverage, or -document"))
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

// jsBundle is the section of a JavaScript bundle; its head alone, with a note, where the coverage
// file or the source map is not there.
func jsBundle(root, bundle, sourceMap, dir, coverage, modules, node, commit string, tools, configs []string) string {

	dir = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(dir)), "./")
	if modules == "" {
		modules = filepath.Join(root, dir, "node_modules")
	}
	b := write.Bundle{Name: filepath.Base(bundle), Dir: dir, Node: node, Commit: commit}
	var err error
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
		b.Commit = head(root)
	}
	if sourceMap == "" {
		sourceMap = bundle + ".map"
	}
	if !exists(coverage) {
		return write.Note(b, "No coverage of the tests: "+filepath.Base(coverage)+" is not there.")
	}
	if !exists(sourceMap) {
		return write.Note(b, "No source map: "+filepath.Base(sourceMap)+" is not there.")
	}
	f, err := os.Open(coverage) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	c, err := write.ReadIstanbul(f, root)
	_ = f.Close()
	if err != nil {
		fail(fmt.Errorf("%s: %w", coverage, err))
	}
	src, err := os.ReadFile(bundle) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	smap, err := os.ReadFile(sourceMap) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	code, err := write.ReadBundle(src, smap, dir, c.Functions())
	if err != nil {
		fail(err)
	}
	return write.Section(b, write.Rows(c.Functions(), code, c))
}

// exists is whether there is a file at path.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// writeDocument writes the object code report from the files of the sections, and its entry in
// documents.yaml.
//
// @implements SR_120@1.1.1, SR_121@1.1.1
func writeDocument(root, label, ticket string, files []string) {

	if err := report.Flags(label, ticket); err != nil {
		fail(err)
	}
	var sections []string
	for _, f := range files {
		data, err := os.ReadFile(f) // #nosec G304 -- the path comes from the user of the tool
		if err != nil {
			fail(err)
		}
		sections = append(sections, string(data))
	}
	const name = "testreports/object-code.md"
	path := filepath.Join(root, filepath.FromSlash(name))
	documentsPath := filepath.Join(root, "documents.yaml")
	documents, err := os.ReadFile(documentsPath) // #nosec G304 -- the path comes from the user of the tool
	if err != nil {
		fail(err)
	}
	next, err := create.Next(create.VersionOf(documents, name))
	if err != nil {
		fail(err)
	}
	date := time.Now().Format("2006-01-02")
	old, _ := os.ReadFile(path) // #nosec G304 -- the path comes from the user of the tool; a report that is not there has no rows to keep
	md := []byte(write.Markdown(write.Document{
		Label:    label,
		Commit:   head(root),
		Date:     date,
		Sections: sections,
		History:  append(report.History(old), row(next, date, label, ticket)),
	}))
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		fail(err)
	}
	if err := os.WriteFile(path, md, 0o600); err != nil {
		fail(err)
	}
	if err := os.WriteFile(documentsPath, create.EntryOf(documents, name, next, create.Hash(md)), 0o600); err != nil { // #nosec G703 -- the path comes from the user of the tool
		fail(err)
	}
	fmt.Fprintf(os.Stderr, "%s written: %d sections, version %s\n", path, len(sections), next)
}

// row is the row of the Change history this run adds; with a ticket but without a label, it names the
// ticket as its reason, with - as its change.
func row(next, date, label, ticket string) string {
	if label == "" && ticket != "" {
		label = "-"
	}
	return report.Row(next, date, label, ticket, "objectcode "+version())
}

// head is the commit of the head of the repo.
func head(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short=7", "HEAD").Output() // #nosec G204 -- root comes from the user of the tool
	if err != nil {
		fail(fmt.Errorf("the commit: %w", err))
	}
	return strings.TrimSpace(string(out))
}

// version is the version of this tool: the commit of tools it was built from, with -dirty where the
// working tree had changes; (devel) where Go recorded none, as with go run.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	var revision, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value
		}
	}
	if revision == "" {
		return info.Main.Version
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if modified == "true" {
		revision += "-dirty"
	}
	return revision
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "objectcode:", err)
	os.Exit(1)
}
