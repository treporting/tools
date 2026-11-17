package review

import (
	"os"
	"reflect"
	"testing"
)

// A file is indexed when it holds an object, whatever its type; a review result and a file without
// objects are not.
//
// @implements TC_57@2.1.1, TC_58@2.1.1
func TestEveryFileThatHoldsObjectsIsIndexed(t *testing.T) {
	files := map[string]string{
		"tools/testcases/software.yaml": `type: testcases
testcases:
- object_key: TC_1
`,
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
`,
		"reviews/results/r.yaml": `type: review
results:
- object_key: RR_1
`,
		"reviews/criteria-r5.yaml": `type: criteria
`,
		"README.md": "# Readme\n",
	}
	read := func(file string) ([]byte, error) {
		if c, ok := files[file]; ok {
			return []byte(c), nil
		}
		return nil, os.ErrNotExist
	}
	ix := NewIndex(read, []string{"README.md", "reviews/criteria-r5.yaml", "reviews/results/r.yaml", "tools/software-requirements.yaml", "tools/testcases/software.yaml"})
	want := []string{"tools/software-requirements.yaml", "tools/testcases/software.yaml"}
	if !reflect.DeepEqual(ix.Files(), want) {
		t.Errorf("files %q, want %q", ix.Files(), want)
	}
	if file, key := ix.FileOf("TC_1"); file != "tools/testcases/software.yaml" || key != "TC_1" {
		t.Errorf("TC_1 in %q as %q", file, key)
	}
	if file, _ := ix.FileOf("RR_1"); file != "" {
		t.Errorf("RR_1 in %q, want none", file)
	}
}
