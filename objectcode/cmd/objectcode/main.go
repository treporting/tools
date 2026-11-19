// Command objectcode writes the section of the object code report of a Go binary (TOOL_35): its
// head, and a row per function of the source code it is built from.
//
//	objectcode -binary file -main ./path/of/main [-coverprofile file] [-commit sha] [-root dir]
//
// It lists the functions of the packages of the module the main package is built from, reads their
// machine code from the binary by `go tool objdump`, and their coverage from the coverage profile of
// the run of the tests, as `go test -coverprofile` writes it. It writes the section to the standard
// output.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/treporting/tools/objectcode/write"
)

// @implements SR_115@1.1.1, SR_117@1.1.1, SR_118@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	binary := flag.String("binary", "", "the Go binary")
	mainPkg := flag.String("main", "", "the main package of the binary, from root, such as ./annotation/cmd/annotationcheck")
	coverprofile := flag.String("coverprofile", "", "the coverage profile of the run of the tests; none where empty")
	commit := flag.String("commit", "", "the commit the binary was built from, where the binary names none")
	flag.Parse()

	if *binary == "" || *mainPkg == "" {
		fail(fmt.Errorf("give -binary and -main"))
	}
	packages, err := write.ListPackages(*root, *mainPkg)
	if err != nil {
		fail(err)
	}
	abs, err := filepath.Abs(*root)
	if err != nil {
		fail(err)
	}
	functions, err := write.Functions(abs, packages)
	if err != nil {
		fail(err)
	}
	d, err := write.Objdump(*binary, packages)
	if err != nil {
		fail(err)
	}
	var profile *write.Profile
	if *coverprofile != "" {
		f, err := os.Open(*coverprofile) // #nosec G304 -- the path comes from the user of the tool
		if err != nil {
			fail(err)
		}
		profile, err = write.ReadProfile(f)
		_ = f.Close()
		if err != nil {
			fail(fmt.Errorf("%s: %w", *coverprofile, err))
		}
	}
	b, err := write.ReadBinary(*binary, *commit)
	if err != nil {
		fail(err)
	}
	b.Name = filepath.Base(*binary)
	b.Main = strings.TrimPrefix(filepath.ToSlash(*mainPkg), "./")
	fmt.Print(write.Section(b, write.Rows(functions, d, profile)))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "objectcode:", err)
	os.Exit(1)
}
