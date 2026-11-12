// Command annotationcheck checks that every id named in @implements, @verifies or @refines is
// approved in the base commit of the pull request (TOOL_06), and names its version (TOOL_20).
//
//	annotationcheck [-root dir] [-base rev]
//
// Without -base, the base commit is where HEAD branched from main. It prints every id that is
// missing or not approved there, or names no version, and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/treporting/specs-and-deploy/tools/annotation/check"
	"github.com/treporting/specs-and-deploy/tools/review"
)

// @implements SR_24@3.1.1, SR_25@2.2.1, SR_27@1.1.1, SR_30@2.2.1, SR_35@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	base := flag.String("base", "", "the base commit of the pull request; where HEAD branched from main, if empty")
	flag.Parse()

	if *base == "" {
		out, err := exec.Command("git", "-C", *root, "merge-base", "HEAD", "main").Output()
		if err != nil {
			fail(fmt.Errorf("git merge-base HEAD main: %w", err))
		}
		*base = strings.TrimSpace(string(out))
	}

	out, err := exec.Command("git", "-C", *root, "ls-files").Output()
	if err != nil {
		fail(err)
	}
	kinds, err := check.ReadKinds(filepath.Join(*root, "tools", "annotation", "kinds.yaml"))
	if err != nil {
		fail(err)
	}
	files := strings.Fields(string(out))
	annotations, unknown, err := check.Find(*root, files, kinds)
	if err != nil {
		fail(err)
	}
	refines, err := check.Refines(*root, files)
	if err != nil {
		fail(err)
	}
	annotations = append(annotations, refines...)
	findings := append(unknown, check.Unversioned(annotations)...)
	index, err := review.IndexAt(*root)
	if err != nil {
		fail(err)
	}
	findings = append(findings, check.Check(review.ExecGit{Root: *root}, *base, index, annotations)...)
	for _, f := range findings {
		fmt.Println(f)
	}
	ids := 0
	for _, a := range annotations {
		ids += len(a.IDs)
	}
	fmt.Printf("%d annotations, %d ids, checked against the base commit %.7s, %d findings\n", len(annotations), ids, *base, len(findings))
	if len(findings) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
