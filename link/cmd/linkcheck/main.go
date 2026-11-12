// Command linkcheck makes the link check (TOOL_16): it compares the version each link, entry of
// the allocation or of controls, and annotation names with the version its object has.
//
//	linkcheck [-root dir] [-above dir]
//
// -above is the root of the repo above, where an object that this repo does not hold is looked up.
// It reads the files as they stand in the working tree. It prints every one whose object has
// another x than it names, and exits with 0 all the same: it reports, and a suspect link blocks
// the baseline, not the pull request (PROC_20). It exits with 2 where it cannot read the repo.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	annotation "github.com/treporting/specs-and-deploy/tools/annotation/check"
	"github.com/treporting/specs-and-deploy/tools/link/check"
	"github.com/treporting/specs-and-deploy/tools/review"
)

// @implements SR_32@1.1.1, SR_33@2.1.1, SR_109@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	aboveRoot := flag.String("above", "", "the root of the repo above, empty for none")
	flag.Parse()

	var above *review.Repo
	if *aboveRoot != "" {
		repo, err := review.RepoAt(*aboveRoot)
		if err != nil {
			fail(err)
		}
		above = repo
	}

	read := func(file string) ([]byte, error) {
		return os.ReadFile(filepath.Join(*root, file))
	}
	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	links, comparedLinks, err := check.LinksWithAbove(read, index, above)
	if err != nil {
		fail(err)
	}

	out, err := exec.Command("git", "-C", *root, "ls-files").Output()
	if err != nil {
		fail(err)
	}
	kinds, err := annotation.ReadKinds(filepath.Join(*root, "tools", "annotation", "kinds.yaml"))
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

	for _, f := range append(links, ids...) {
		fmt.Println(f)
	}
	fmt.Printf("%d links and %d ids of annotations compared, %d suspect\n", comparedLinks, comparedIDs, len(links)+len(ids))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
