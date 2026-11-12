// Command approvalcheck checks every approval from the approved objects (TOOL_32, TOOL_33).
//
//	approvalcheck [-root dir] [-base commit]
//
// With -base, it checks the objects in the files changed against that commit; without, every
// approved object. It prints every approval that does not hold and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/treporting/specs-and-deploy/tools/review"
	"github.com/treporting/specs-and-deploy/tools/review/approval"
)

// @implements SR_80@2.1.3
func main() {

	root := flag.String("root", ".", "the root of the repo")
	base := flag.String("base", "", "the base commit of a pull request; empty for every approved object")
	flag.Parse()

	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	findings, err := approval.Check(*root, *base, review.ExecGit{Root: *root}, index)
	if err != nil {
		fail(err)
	}
	for _, f := range findings {
		fmt.Println(f)
	}
	fmt.Printf("approvals checked, %d findings\n", len(findings))
	if len(findings) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
