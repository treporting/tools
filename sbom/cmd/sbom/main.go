// Command sbom creates the SBOM of a subsystem as sbom.md in the folder sbom/ of its repo (TOOL_21).
//
//	sbom -label x.y.z -ticket number [-store] [-root dir] [-repo name] [-image name]...
//	sbom -check [-root dir] [-repo name] [-image name]...
//
// It runs in the repo checked out at the baseline the label names. It lists the package of the
// subsystem itself; the Go modules of its commands and the standard library of Go, where a directory
// of the repo holds a go.mod; the packages of each package-lock.json of the repo; and the Alpine
// packages of each image named. It lists each package once. It keeps the rows of the Change history of an sbom.md that is there,
// and adds one. It sets the entry of sbom.md in documents.yaml of the repo: the version of the new
// row, the hash of the file written, completed and generated (PROC_13). It lists in sbom.md each
// file of original data it read: the output of go list, each package-lock.json, the Alpine packages
// of each image. With -store it writes them into sbom/ as <ts>_<file>, and names them in
// sbom/<ts>.yaml; the output of go list with its module cache as $GOMODCACHE.
//
// With -check it writes nothing: it compares the packages with the table of sbom.md, and exits with
// 1 where they differ (TOOL_28).
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

	"github.com/treporting/specs-and-deploy/tools/sbom/create"
)

// images collects the flag -image, which may be given more than once.
type images []string

func (i *images) String() string     { return strings.Join(*i, ",") }
func (i *images) Set(v string) error { *i = append(*i, v); return nil }

// @implements SR_39@2.1.1
// @implements SR_59@2.1.1
// @implements SR_101@1.1.1, SR_102@1.1.1, SR_105@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo of the subsystem")
	label := flag.String("label", "", "the label of the baseline, such as 1.2.0")
	ticket := flag.String("ticket", "", "the number of the ticket of the baseline, such as 343")
	repo := flag.String("repo", "", "the name of the repo; the name of the directory of root where empty")
	check := flag.Bool("check", false, "compare the packages with the table of sbom.md, and write nothing")
	store := flag.Bool("store", false, "store the original data in sbom/, named in sbom/<ts>.yaml")
	var named images
	flag.Var(&named, "image", "an image of the subsystem, built at the baseline; may be given more than once")
	flag.Parse()

	if *repo == "" {
		abs, err := filepath.Abs(*root)
		if err != nil {
			fail(err)
		}
		*repo = filepath.Base(abs)
	}
	path := filepath.Join(*root, create.Name)

	if *check {
		checkSbom(path, *root, *repo, named)
		return
	}

	if *label == "" {
		fail(fmt.Errorf("-label is missing"))
	}
	if *ticket == "" {
		fail(fmt.Errorf("-ticket is missing"))
	}
	documentsPath := filepath.Join(*root, "documents.yaml")
	documents, err := os.ReadFile(documentsPath)
	if err != nil {
		fail(err)
	}
	next, err := create.Next(create.Version(documents))
	if err != nil {
		fail(err)
	}

	packages, originals := delivered(*root, *repo, *label, named)

	now := time.Now()
	old, _ := os.ReadFile(path) // an sbom.md that is not there has no rows to keep
	history := append(create.History(old), create.Row(next, now.Format("2006-01-02"), *label, *ticket, "sbom "+version()))

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		fail(err)
	}
	if *store && len(originals) > 0 {
		originals = storeOriginals(*root, *label, now, originals)
	}
	md := []byte(create.Markdown(create.SBOM{Repo: *repo, Packages: packages, Original: originals, History: history}))
	if err := os.WriteFile(path, md, 0o600); err != nil {
		fail(err)
	}
	if err := os.WriteFile(documentsPath, create.Entry(documents, next, create.Hash(md)), 0o600); err != nil {
		fail(err)
	}
	fmt.Printf("%s written: %d packages, version %s\n", path, len(packages), next)
}

// checkSbom compares the packages the subsystem delivers with the table of its sbom.md, with the
// label of its last row as the version of its own package. It writes nothing, and exits with 1
// naming each package that differs.
//
// @implements SR_66@1.1.1
func checkSbom(path, root, repo string, named []string) {
	md, err := os.ReadFile(path)
	if err != nil {
		fail(err)
	}
	label := create.Label(md)
	if label == "" {
		fail(fmt.Errorf("%s: no row in its Change history, so no label", path))
	}
	packages, _ := delivered(root, repo, label, named)
	differ := create.Differ(packages, create.Table(md))
	for _, line := range differ {
		fmt.Println(line)
	}
	if len(differ) > 0 {
		fmt.Printf("%s does not list the packages delivered: %d differ\n", path, len(differ))
		os.Exit(1)
	}
	fmt.Printf("%s lists the packages delivered\n", path)
}

// delivered answers every package the subsystem delivers, each once: its own, with the label as its
// version; the Go modules of its commands and the standard library of Go, where a directory of the
// repo holds a go.mod; the packages of each package-lock.json of the repo; and the Alpine packages of
// each image. It answers as well each file of original data it read.
func delivered(root, repo, label string, named []string) ([]create.Package, []create.Original) {
	packages := []create.Package{create.Own(repo, label, create.LicenceIn(root))}
	var originals []create.Original
	for _, source := range []func(string, []string) ([]create.Package, []create.Original){goPackages, npmPackages, apkPackages} {
		p, o := source(root, named)
		packages = append(packages, p...)
		originals = append(originals, o...)
	}
	return create.Once(packages), originals
}

// storeOriginals writes the originals into sbom/ and names them in sbom/<ts>.yaml, ts the time of
// the run, and answers them with where each is stored.
func storeOriginals(root, label string, now time.Time, originals []create.Original) []create.Original {
	ts := now.Format("20060102T150405")
	stored, err := create.Store(root, create.Folder, ts+"_", originals)
	if err != nil {
		fail(err)
	}
	record, err := create.OriginalYAML(create.Name, label, commit(root), now.Format("2006-01-02T15:04:05-07:00"), stored)
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile(filepath.Join(root, create.Folder, ts+".yaml"), record, 0o600); err != nil {
		fail(err)
	}
	return stored
}

// commit answers the commit of the head of the repo in root.
func commit(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short=7", "HEAD").Output() // #nosec G204 -- root is the caller's own flag
	if err != nil {
		fail(fmt.Errorf("git rev-parse: %w", err))
	}
	return strings.TrimSpace(string(out))
}

// goPackages answers the Go modules of the commands of every go.mod of the repo, and the standard
// library of Go where there is a command; nothing where there is none.
func goPackages(root string, _ []string) ([]create.Package, []create.Original) {
	dirs, err := create.DirsWith(root, "go.mod")
	if err != nil {
		fail(err)
	}
	var lists, withCommands []string
	var originals []create.Original
	for _, dir := range dirs {
		list, err := create.GoList(dir)
		if err != nil {
			fail(err)
		}
		if list != "" {
			lists = append(lists, list)
			withCommands = append(withCommands, dir)
			modcache, err := create.GoEnv(dir, "GOMODCACHE")
			if err != nil {
				fail(err)
			}
			originals = append(originals, create.GoListOriginal(root, dir, list, modcache))
		}
	}
	if len(lists) == 0 {
		return nil, nil
	}
	return append(create.GoModules(strings.Join(lists, "\n"), create.LicenceIn), stdlib(withCommands[0])), originals
}

// stdlib answers the standard library of the Go that builds the module in dir.
func stdlib(dir string) create.Package {
	version, err := create.GoEnv(dir, "GOVERSION")
	if err != nil {
		fail(err)
	}
	goroot, err := create.GoEnv(dir, "GOROOT")
	if err != nil {
		fail(err)
	}
	licence := create.LicenceIn(goroot)
	if licence == "" {
		licence = create.LicenceIn(filepath.Dir(goroot)) // Homebrew keeps it beside GOROOT
	}
	return create.Stdlib(version, licence)
}

// npmPackages answers the packages of every package-lock.json of the repo.
func npmPackages(root string, _ []string) ([]create.Package, []create.Original) {
	dirs, err := create.DirsWith(root, "package-lock.json")
	if err != nil {
		fail(err)
	}
	var packages []create.Package
	var originals []create.Original
	for _, dir := range dirs {
		lock, err := os.ReadFile(filepath.Join(dir, "package-lock.json"))
		if err != nil {
			fail(err)
		}
		npm, err := create.NpmPackages(lock)
		if err != nil {
			fail(err)
		}
		packages = append(packages, npm...)
		originals = append(originals, create.LockOriginal(root, dir, lock))
	}
	return packages, originals
}

// apkPackages answers the Alpine packages of each image named.
func apkPackages(_ string, named []string) ([]create.Package, []create.Original) {
	var packages []create.Package
	var originals []create.Original
	for _, image := range named {
		installed, err := create.ApkInstalled(image)
		if err != nil {
			fail(err)
		}
		packages = append(packages, create.ApkPackages(installed)...)
		originals = append(originals, create.ApkOriginal(image, installed))
	}
	return packages, originals
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
