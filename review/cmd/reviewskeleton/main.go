// Command reviewskeleton writes the skeleton of a review result into reviews/results/ (TOOL_04).
//
//	reviewskeleton [-root dir] -reviewer name -scope scope id-or-document ...
//
// It prints the path of the file it wrote. Every criterion and the verdict are todo, for the
// reviewer to decide; reviewcheck reports them until then.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/treporting/specs-and-deploy/tools/review"
	"github.com/treporting/specs-and-deploy/tools/review/skeleton"
)

// @implements SR_11@2.1.1, SR_12@1.1.1, SR_13@1.1.1, SR_14@1.1.1, SR_15@1.1.1, SR_17@2.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	reviewer := flag.String("reviewer", "", "the name of the reviewer, as created_by")
	scope := flag.String("scope", "", "what is reviewed, for the name of the file, such as harness")
	flag.Parse()

	if *reviewer == "" || *scope == "" || flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: reviewskeleton [-root dir] -reviewer name -scope scope id-or-document ...")
		os.Exit(2)
	}

	index, err := review.IndexAt(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	path, err := skeleton.Make(*root, skeleton.Request{
		Reviewer: *reviewer,
		Scope:    *scope,
		Targets:  flag.Args(),
		Now:      time.Now(),
	}, review.ExecGit{Root: *root}, index)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(path)
}
