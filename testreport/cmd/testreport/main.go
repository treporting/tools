// Command testreport writes the test report of a test level as testreports/<test level>.md beside
// testcases/ (TOOL_34).
//
//	testreport -level r5tc -junit file -made-by command ... [-runtime "name version"]...
//	           [-store] [-label x.y.z -ticket number] [-root dir] [-dir dir] [-levels file]
//
// It reads the result of each test from each junit.xml given, the test cases of the level from
// testcases/, the name of the level from levels.yaml - at -levels, or in the root of the repo - and
// the test procedures from the annotations of the files of tests under -dir. Each
// -junit takes the -made-by given in the same place among them. With -store it copies each junit.xml
// into testreports/ as <ts>_<test level>-<file>, and names them in testreports/<ts>_<test level>.yaml.
// It keeps the rows of the Change history of a test report that is there, and adds
// one. It sets the entry of the test report in documents.yaml of the repo: the version of the new
// row, the hash of the file written, completed and generated (PROC_13). A test that fails does not
// fail it.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/treporting/specs-and-deploy/tools/sbom/create"
	"github.com/treporting/specs-and-deploy/tools/testreport/write"
)

// list collects a flag that may be given more than once.
type list []string

func (l *list) String() string     { return strings.Join(*l, ",") }
func (l *list) Set(v string) error { *l = append(*l, v); return nil }

// @implements SR_88@2.1.1, SR_89@4.1.1, SR_106@1.1.1, SR_96@2.1.1, SR_98@2.1.1, SR_91@2.1.1, SR_104@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	dir := flag.String("dir", ".", "the directory of testcases/ and of the tests, relative to root, such as tools")
	level := flag.String("level", "", "the test level, such as r5tc")
	label := flag.String("label", "", "the label of the baseline, such as 1.2.0")
	ticket := flag.String("ticket", "", "the number of the ticket of the baseline, such as 438")
	levels := flag.String("levels", "", "levels.yaml, which names each test level; levels.yaml in root where empty")
	store := flag.Bool("store", false, "copy each junit.xml into testreports/, named in testreports/<ts>_<test level>.yaml")
	var junits, madeBy, runtimes list
	flag.Var(&junits, "junit", "a junit.xml of a run of the tests; may be given more than once")
	flag.Var(&madeBy, "made-by", "the command that made the junit.xml given in the same place")
	flag.Var(&runtimes, "runtime", "a runtime or test runner the tests ran with, with its version, such as \"Go go1.25.3\"")
	flag.Parse()

	if err := write.Flags(*label, *ticket); err != nil {
		fail(err)
	}
	file, ok := write.Files[*level]
	if !ok {
		fail(fmt.Errorf("-level %q is no test level", *level))
	}
	if *levels == "" {
		*levels = filepath.Join(*root, "levels.yaml")
	}
	levelsYAML, err := os.ReadFile(*levels)
	if err != nil {
		fail(err)
	}
	levelName, err := write.ReadLevel(levelsYAML, *level)
	if err != nil {
		fail(fmt.Errorf("%s: %w", *levels, err))
	}
	if len(junits) == 0 || len(madeBy) != len(junits) {
		fail(fmt.Errorf("give each -junit with its -made-by"))
	}

	testcases, err := os.ReadFile(filepath.Join(*root, *dir, "testcases", file+".yaml"))
	if err != nil && !os.IsNotExist(err) {
		fail(err)
	}
	cases, err := write.ReadTestcases(testcases)
	if err != nil {
		fail(err)
	}

	var results []write.Result
	var originals []create.Original
	for i, j := range junits {
		data, err := os.ReadFile(j)
		if err != nil {
			fail(err)
		}
		r, err := write.ReadJUnit(data)
		if err != nil {
			fail(fmt.Errorf("%s: %w", j, err))
		}
		results = append(results, r...)
		originals = append(originals, create.Original{File: filepath.Base(j), MadeBy: madeBy[i], StoredAt: "-", Data: data})
	}

	files := gitFiles(*root, *dir)
	module := moduleOf(filepath.Join(*root, *dir, "go.mod"))
	procedures, err := write.ReadProcedures(*root, files, module, filepath.Clean(*dir))
	if err != nil {
		fail(err)
	}

	name := filepath.ToSlash(filepath.Join(*dir, "testreports", file+".md"))
	path := filepath.Join(*root, name)
	documentsPath := filepath.Join(*root, "documents.yaml")
	documents, err := os.ReadFile(documentsPath)
	if err != nil {
		fail(err)
	}
	next, err := create.Next(create.VersionOf(documents, name))
	if err != nil {
		fail(err)
	}
	now := time.Now()
	date := now.Format("2006-01-02")
	head := commit(*root)
	if *store {
		originals = storeJUnits(*root, filepath.Dir(name), name, *level, *label, head, now, originals)
	}
	old, _ := os.ReadFile(path) // a test report that is not there has no rows to keep
	history := append(write.History(old), write.Row(next, date, *label, *ticket, "testreport "+version()))

	md := []byte(write.Markdown(write.Report{
		Level:      *level,
		Name:       levelName,
		Label:      *label,
		Commit:     head,
		Date:       date,
		System:     runtime.GOOS + "/" + runtime.GOARCH,
		Runner:     runner(),
		Runtimes:   runtimes,
		Originals:  reported(originals),
		Testcases:  cases,
		Procedures: procedures,
		Results:    results,
		History:    history,
	}))
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		fail(err)
	}
	if err := os.WriteFile(path, md, 0o600); err != nil {
		fail(err)
	}
	if err := os.WriteFile(documentsPath, create.EntryOf(documents, name, next, create.Hash(md)), 0o600); err != nil {
		fail(err)
	}
	fmt.Printf("%s written: %d test cases, %d test procedures, %d results, version %s\n", path, len(cases), len(procedures), len(results), next)
}

// storeJUnits copies each junit.xml into folder as <ts>_<test level>-<file> and names them in
// <ts>_<test level>.yaml there, ts the time of the run, and answers them with where each is stored.
func storeJUnits(root, folder, document, level, label, head string, now time.Time, originals []create.Original) []create.Original {
	ts := now.Format("20060102T150405")
	stored, err := create.Store(root, folder, ts+"_"+level+"-", originals)
	if err != nil {
		fail(err)
	}
	record, err := create.OriginalYAML(document, label, head, now.Format("2006-01-02T15:04:05-07:00"), stored)
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile(filepath.Join(root, folder, ts+"_"+level+".yaml"), record, 0o600); err != nil {
		fail(err)
	}
	return stored
}

// reported answers the originals as the test report lists them.
func reported(originals []create.Original) []write.Original {
	var rows []write.Original
	for _, o := range originals {
		rows = append(rows, write.Original{File: o.File, MadeBy: o.MadeBy, StoredAt: o.StoredAt})
	}
	return rows
}

// gitFiles answers the files git knows under dir, relative to root.
func gitFiles(root, dir string) []string {
	out, err := exec.Command("git", "-C", root, "ls-files", "--", dir).Output() // #nosec G204 -- root and dir are the caller's own flags
	if err != nil {
		fail(fmt.Errorf("git ls-files: %w", err))
	}
	return strings.Fields(string(out))
}

// moduleOf answers the module a go.mod names, and nothing where there is none.
func moduleOf(goMod string) string {
	data, err := os.ReadFile(goMod) // #nosec G304 -- the go.mod of the caller's own directory
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if m, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(m)
		}
	}
	return ""
}

// commit answers the commit of the head of the repo.
func commit(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short=7", "HEAD").Output() // #nosec G204 -- root is the caller's own flag
	if err != nil {
		fail(fmt.Errorf("git rev-parse: %w", err))
	}
	return strings.TrimSpace(string(out))
}

// runner answers the runner of GitHub Actions the run is on, and local outside it.
func runner() string {
	if name := os.Getenv("RUNNER_NAME"); name != "" {
		return name
	}
	return "local"
}

// version answers the version of this tool: the commit of specs-and-deploy it was built from, as Go
// records it in the build, with -dirty where the working tree had changes; (devel) where Go recorded
// none, as with go run.
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
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
