// Command uplinkcheck checks that every approved requirement names approved uplinks only (TOOL_9),
// and that every requirement of r3, r6 and r7 names a part and an interface of its level (TOOL_39).
//
//	uplinkcheck [-root dir] [-above dir]
//
// -above is the root of the repo above, where an uplink that this repo does not hold is looked up.
// It reads the files of requirements as they stand in the working tree. It prints every uplink of
// an approved requirement that is missing or not approved, and every component and interface that
// names no part of its level, and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treporting/tools/review"
	"github.com/treporting/tools/uplink/check"
)

// @implements SR_28@1.1.1, SR_108@1.1.1, SR_129@1.1.1, SR_130@1.1.1
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

	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	read := func(file string) ([]byte, error) {
		return os.ReadFile(filepath.Join(*root, file)) // #nosec G304 -- the path comes from the user of the tool, not from outside
	}
	findings, approved, err := check.CheckWithAbove(read, index, above)
	if err != nil {
		fail(err)
	}
	findings = append(findings, check.Parts(read, index)...)
	for _, f := range findings {
		fmt.Println(f)
	}
	fmt.Printf("%d approved requirements checked, %d findings\n", approved, len(findings))
	if len(findings) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
