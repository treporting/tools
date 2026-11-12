// Command rulesetcheck compares the ruleset of main on GitHub with .github/rulesets/main.json (TOOL_13).
//
//	rulesetcheck [-root dir] [-repo owner/name] [-api url]
//
// It reads the token from GITHUB_TOKEN where it is set. It prints every difference, and every
// ruleset that cannot be read, and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treporting/specs-and-deploy/tools/ruleset/check"
)

// @implements SR_21@1.1.1, SR_22@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	repo := flag.String("repo", "treporting/specs-and-deploy", "the repo on GitHub")
	api := flag.String("api", "https://api.github.com", "the API of GitHub")
	flag.Parse()

	want, err := os.ReadFile(filepath.Join(*root, ".github", "rulesets", "main.json"))
	if err != nil {
		fail(err)
	}
	findings, err := check.Check(check.Client{API: *api, Repo: *repo, Token: os.Getenv("GITHUB_TOKEN")}, want)
	if err != nil {
		fail(err)
	}
	for _, f := range findings {
		fmt.Println(f)
	}
	fmt.Printf("the ruleset of main checked against .github/rulesets/main.json, %d findings\n", len(findings))
	if len(findings) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
