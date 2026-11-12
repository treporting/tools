// Command licencecheck checks the licences of sbom.md against Apache-2.0 (TOOL_22).
//
//	licencecheck [-root dir]
//
// It reads sbom.md in the folder sbom/ of the repo of a subsystem, and licences_assessed.yaml where the root
// holds one. It prints every package that fails with its PURL and its licence, and exits with 1
// where there is one.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/treporting/specs-and-deploy/tools/licence/check"
)

// @implements SR_50@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo of the subsystem")
	flag.Parse()

	sbomMD, err := os.ReadFile(filepath.Join(*root, "sbom", "sbom.md"))
	if err != nil {
		fail(err)
	}
	assessed := map[string]string{}
	file, err := os.ReadFile(filepath.Join(*root, "licences_assessed.yaml"))
	switch {
	case err == nil:
		if assessed, err = check.ReadAssessed(file); err != nil {
			fail(err)
		}
	case !errors.Is(err, fs.ErrNotExist):
		fail(err)
	}

	packages := check.Packages(sbomMD)
	failures := check.Check(packages, assessed)
	for _, f := range failures {
		fmt.Println(f)
	}
	fmt.Printf("%d packages checked, %d fail\n", len(packages), len(failures))
	if len(failures) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
