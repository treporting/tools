// Command linkcheck makes the link check (TOOL_16): it compares the version each link, entry of
// the allocation or of controls, and annotation names with the version its object has. It makes the
// status check too (TOOL_36): it reports each whose object is obsolete, and compares no link of an
// obsolete object (TOOL_38).
//
//	linkcheck [-root dir] [-above dir] [-kinds file]
//
// -above is the root of the repo above, where an object that this repo does not hold is looked up.
// -kinds is the list of the kinds of files of the repo; annotation/kinds.yaml below the root where it
// is not given.
// It reads the files as they stand in the working tree. It prints every one whose object has
// another x than it names or is obsolete, and exits with 0 all the same: it reports, and such a link
// blocks the baseline, not the pull request (PROC_20, PROC_23). It exits with 2 where it cannot read the repo.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	annotation "github.com/treporting/tools/annotation/check"
	"github.com/treporting/tools/link/check"
	"github.com/treporting/tools/review"
)

// @implements SR_32@1.1.1, SR_33@2.1.1, SR_109@1.1.1, SR_127@1.1.1, SR_128@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	aboveRoot := flag.String("above", "", "the root of the repo above, empty for none")
	kindsFile := flag.String("kinds", "", "the kinds of files of the repo; annotation/kinds.yaml below the root, if empty")
	flag.Parse()
	if *kindsFile == "" {
		*kindsFile = filepath.Join(*root, "annotation", "kinds.yaml")
	}

	var above *review.Repo
	if *aboveRoot != "" {
		repo, err := review.RepoAt(*aboveRoot)
		if err != nil {
			fail(err)
		}
		above = repo
	}

	read := func(file string) ([]byte, error) {
		return os.ReadFile(filepath.Join(*root, file)) // #nosec G304 -- the path comes from the user of the tool, not from outside
	}
	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	links, comparedLinks, err := check.LinksWithAbove(read, index, above)
	if err != nil {
		fail(err)
	}

	out, err := exec.Command("git", "-C", *root, "ls-files").Output() // #nosec G204 -- git or a command the user of the tool names, not input from outside
	if err != nil {
		fail(err)
	}
	kinds, err := annotation.ReadKinds(*kindsFile)
	if err != nil {
		fail(err)
	}
	//	a file of unknown kind is annotationcheck's to report
	files := strings.Fields(string(out))
	annotations, _, err := annotation.Find(*root, files, kinds)
	if err != nil {
		fail(err)
	}
	refines, err := annotation.Refines(*root, files)
	if err != nil {
		fail(err)
	}
	annotations = append(annotations, refines...)
	ids, comparedIDs := check.AnnotationsWithAbove(read, index, above, annotations)

	suspect, obsolete := 0, 0
	for _, f := range append(links, ids...) {
		fmt.Println(f)
		if f.Obsolete {
			obsolete++
		} else {
			suspect++
		}
	}
	fmt.Printf("%d links and %d ids of annotations compared, %d suspect, %d obsolete\n", comparedLinks, comparedIDs, suspect, obsolete)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
