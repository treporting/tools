// Package write writes the test report of a test level (TOOL_34): what was specified, implemented
// and run, from the junit.xml of a run, the test cases of the level and the annotations of its tests.
package write

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/treporting/specs-and-deploy/tools/annotation/check"
)

// Files names the file of the test cases and of the test report of each test level, without .yaml
// or .md (formats.md, section Test cases).
var Files = map[string]string{
	"r1tc": "acceptance",
	"r2tc": "system",
	"r3tc": "system-integration",
	"r5tc": "software",
	"r6tc": "software-integration",
	"r7tc": "unit",
}

// ReadLevel answers the name of the test level key as levels.yaml names it, with its abbreviation,
// such as subsystem test (SwVal); an error where levels.yaml names no such test level.
//
// @implements SR_106@1.1.1
func ReadLevel(levelsYAML []byte, key string) (string, error) {
	var file struct {
		Levels []struct {
			Key          string `yaml:"key"`
			Name         string `yaml:"name"`
			Abbreviation string `yaml:"abbreviation"`
		} `yaml:"levels"`
	}
	if err := yaml.Unmarshal(levelsYAML, &file); err != nil {
		return "", err
	}
	for _, l := range file.Levels {
		if l.Key == key {
			if l.Abbreviation == "" {
				return l.Name, nil
			}
			return l.Name + " (" + l.Abbreviation + ")", nil
		}
	}
	return "", fmt.Errorf("levels.yaml names no test level %s", key)
}

// Testcase is one test case of a test level.
type Testcase struct {
	Key     string   `yaml:"object_key"`
	Version string   `yaml:"object_version"`
	Text    string   `yaml:"object_text"`
	Targets []string `yaml:"target_ids"`
	Status  string   `yaml:"status"`
}

// Procedure is one test procedure: a test function with its annotations.
type Procedure struct {
	Name       string
	File       string   // relative to the root of the repo
	Suite      string   // the classname a junit.xml gives its result: the package of Go, the file of Vitest
	Implements []string // test cases, each as ID@x.y.z
	Verifies   []string // requirements, each as ID@x.y.z
}

// Result is the result of one test in a junit.xml.
type Result struct {
	Suite   string
	Name    string
	Verdict string // passed, failed or skipped
}

// ReadTestcases answers the test cases of a file of test cases.
//
// @implements SR_88@2.1.1
func ReadTestcases(data []byte) ([]Testcase, error) {
	var file struct {
		Testcases []Testcase `yaml:"testcases"`
	}
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	return file.Testcases, nil
}

// ReadJUnit answers the result of each test of a junit.xml. A subtest of Go, named with a /, is part
// of its test and has no result of its own.
//
// @implements SR_88@2.1.1
func ReadJUnit(data []byte) ([]Result, error) {
	type testcase struct {
		Classname string    `xml:"classname,attr"`
		Name      string    `xml:"name,attr"`
		Failure   *struct{} `xml:"failure"`
		Error     *struct{} `xml:"error"`
		Skipped   *struct{} `xml:"skipped"`
	}
	type suite struct {
		Testcases []testcase `xml:"testcase"`
	}
	var doc struct {
		Suites    []suite    `xml:"testsuite"`
		Testcases []testcase `xml:"testcase"` // a junit.xml whose root is one testsuite
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("no junit.xml: %w", err)
	}
	all := doc.Testcases
	for _, s := range doc.Suites {
		all = append(all, s.Testcases...)
	}
	var results []Result
	for _, t := range all {
		if strings.Contains(t.Name, "/") {
			continue
		}
		verdict := "passed"
		switch {
		case t.Failure != nil || t.Error != nil:
			verdict = "failed"
		case t.Skipped != nil:
			verdict = "skipped"
		}
		results = append(results, Result{Suite: t.Classname, Name: t.Name, Verdict: verdict})
	}
	return results, nil
}

// testFiles is the kind of a file of tests: of Go and of Vitest.
var testFiles = check.Kinds{Whitelist: []string{"*_test.go", "*.test.js", "*.test.jsx", "*.test.ts", "*.test.tsx"}}

// declaration is the line that names a test: a test function of Go, or it or test of Vitest.
var declaration = regexp.MustCompile("^\\s*(?:func (Test\\w+)\\(|(?:it|test)\\(\\s*['\"`](.+?)['\"`])")

// ReadProcedures answers the test procedures of the files, relative to root: every test whose
// comment above it holds an @implements or an @verifies. module is the module of Go of the files,
// with modRoot its directory relative to root, so that a test of Go is matched with its package.
//
// @implements SR_88@2.1.1
func ReadProcedures(root string, files []string, module, modRoot string) ([]Procedure, error) {
	var tests []string
	for _, f := range files {
		if matched(f) {
			tests = append(tests, f)
		}
	}
	annotations, _, err := check.Find(root, tests, testFiles)
	if err != nil {
		return nil, err
	}
	byFile := map[string][]check.Annotation{}
	for _, a := range annotations {
		byFile[a.File] = append(byFile[a.File], a)
	}
	var procedures []Procedure
	for _, file := range tests {
		if len(byFile[file]) == 0 {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			return nil, err
		}
		procedures = append(procedures, inFile(file, strings.Split(string(data), "\n"), byFile[file], suite(file, module, modRoot))...)
	}
	return procedures, nil
}

// inFile answers the test procedures of one file: each annotation belongs to the first test
// declared after it.
func inFile(file string, lines []string, annotations []check.Annotation, suite string) []Procedure {
	var procedures []Procedure
	var open Procedure
	next := 0
	for i, line := range lines {
		for next < len(annotations) && annotations[next].Line == i+1 {
			a := annotations[next]
			for j, id := range a.IDs {
				named := id + "@" + a.Versions[j]
				if a.Kind == "implements" {
					open.Implements = append(open.Implements, named)
				} else {
					open.Verifies = append(open.Verifies, named)
				}
			}
			next++
		}
		m := declaration.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if len(open.Implements)+len(open.Verifies) > 0 {
			open.Name, open.File, open.Suite = m[1]+m[2], file, suite
			procedures = append(procedures, open)
		}
		open = Procedure{}
	}
	return procedures
}

// matched says whether a file is a file of tests.
func matched(file string) bool {
	for _, p := range testFiles.Whitelist {
		if ok, _ := filepath.Match(p, filepath.Base(file)); ok {
			return true
		}
	}
	return false
}

// suite answers the classname a junit.xml gives the tests of a file: the import path of its package
// for Go, the file itself otherwise.
func suite(file, module, modRoot string) string {
	if !strings.HasSuffix(file, ".go") {
		return file
	}
	rel, err := filepath.Rel(modRoot, filepath.Dir(file))
	if err != nil || rel == "." {
		return module
	}
	return module + "/" + filepath.ToSlash(rel)
}
