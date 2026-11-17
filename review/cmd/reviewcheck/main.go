// Command reviewcheck checks the review results of a repo against the rules of formats.md (TOOL_3).
//
//	reviewcheck [-root dir] [-base commit] [file ...]
//
// With no file, it checks every result in reviews/results/; with -base, only those the pull request
// adds or changes against that commit. It prints every
// rule a result breaks and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/treporting/tools/review"
	"github.com/treporting/tools/review/check"
)

// @implements SR_1@1.1.1, SR_2@1.1.1, SR_3@1.1.1, SR_4@1.1.1, SR_5@2.1.1, SR_6@1.1.1, SR_7@1.1.1, SR_8@1.1.1, SR_103@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	base := flag.String("base", "", "the base commit of a pull request; empty for every review result")
	flag.Parse()

	files := flag.Args()
	git := review.ExecGit{Root: *root}
	if len(files) == 0 {
		var err error
		if files, err = check.Results(*root, *base, git.Changed); err != nil {
			fail(err)
		}
	}

	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	broken := 0
	for _, f := range files {
		findings, err := check.Check(*root, f, git, index)
		if err != nil {
			fail(err)
		}
		for _, finding := range findings {
			fmt.Println(finding)
		}
		broken += len(findings)
	}
	fmt.Printf("%d review results checked, %d findings\n", len(files), broken)
	if broken > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
