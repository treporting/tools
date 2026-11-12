// Command schemacheck checks every YAML file under version control against the schema it names in
// its field schema (TOOL_05).
//
//	schemacheck [-root dir] [-above dir]
//
// -above is the root of the repo above, where a schema that this repo does not hold is taken from.
// It prints every place where a file does not match, and every file that names no schema or a
// schema that is not there, and exits with 1 where there is one.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/treporting/specs-and-deploy/tools/schema/check"
)

// @implements SR_18@1.1.1, SR_19@1.1.1, SR_107@1.1.1
func main() {

	root := flag.String("root", ".", "the root of the repo")
	above := flag.String("above", "", "the root of the repo above, empty for none")
	flag.Parse()

	files, err := check.Files(*root)
	if err != nil {
		fail(err)
	}
	findings, err := check.CheckWithAbove(*root, *above, files)
	if err != nil {
		fail(err)
	}
	for _, f := range findings {
		fmt.Println(f)
	}
	fmt.Printf("%d YAML files checked, %d findings\n", len(files), len(findings))
	if len(findings) > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
