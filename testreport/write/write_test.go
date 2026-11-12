package write

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/treporting/specs-and-deploy/tools/sbom/create"
)

const goJUnit = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite name="example.org/m/a" tests="3">
		<testcase classname="example.org/m/a" name="TestOne" time="0.01"></testcase>
		<testcase classname="example.org/m/a" name="TestOne/sub" time="0.01"><failure message="sub"></failure></testcase>
		<testcase classname="example.org/m/a" name="TestTwo" time="0.01"><failure message="no"></failure></testcase>
		<testcase classname="example.org/m/a" name="TestThree" time="0.01"><skipped></skipped></testcase>
	</testsuite>
</testsuites>`

const vitestJUnit = `<?xml version="1.0" encoding="UTF-8" ?>
<testsuites name="vitest tests" tests="1">
	<testsuite name="src/app.test.ts" tests="1">
		<testcase classname="src/app.test.ts" name="App &gt; shows the rail" time="0.002"></testcase>
	</testsuite>
</testsuites>`

// The files of tests below write AT for the @ of an annotation, so that their annotations are read
// only by the test, not by annotationcheck or by a test report of this repo.
const goTest = `package a

import "testing"

// AT implements TC_01@1.1.1
func TestOne(t *testing.T) {}

// AT verifies SR_02@1.1.1
func TestTwo(t *testing.T) {}

// AT implements TC_01@1.1.1
func TestThree(t *testing.T) {}

// AT verifies SR_02@1.1.1
func TestNotRun(t *testing.T) {}

func TestWithoutAnnotation(t *testing.T) {}
`

const vitestTest = `import { it } from 'vitest'

// AT implements TC_03@1.1.1
it('shows the rail', () => {})
`

var testcases = []Testcase{
	{Key: "TC_01", Version: "1.1.1", Text: "one", Targets: []string{"SR_01@1.1.1"}, Status: "approved"},
	{Key: "TC_03", Version: "1.1.1", Text: "three", Targets: []string{"SR_03@2.1.1"}, Status: "approved"},
	{Key: "TC_04", Version: "1.1.1", Text: "four", Targets: []string{"SR_04@1.1.1"}, Status: "approved"},
}

// repo writes the files of tests into a directory and answers the test procedures read from it.
func repo(t *testing.T) []Procedure {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{"m/a/a_test.go": goTest, "m/viewer/src/app.test.ts": vitestTest}
	for name, content := range files {
		content = strings.ReplaceAll(content, "AT ", "@")
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	procedures, err := ReadProcedures(root, []string{"m/a/a_test.go", "m/viewer/src/app.test.ts"}, "example.org/m", "m")
	if err != nil {
		t.Fatal(err)
	}
	return procedures
}

// results answers the results of both junit.xml.
func results(t *testing.T) []Result {
	t.Helper()
	var all []Result
	for _, x := range []string{goJUnit, vitestJUnit} {
		r, err := ReadJUnit([]byte(x))
		if err != nil {
			t.Fatal(err)
		}
		all = append(all, r...)
	}
	return all
}

// report answers a test report of the test cases, the procedures and the results.
func report(t *testing.T) string {
	t.Helper()
	return Markdown(Report{
		Level: "r5tc", Name: "subsystem test (SwVal)", Label: "1.2.0", Commit: "a1b2c3d", Date: "2026-11-04", System: "linux/amd64", Runner: "GitHub Actions 2",
		Runtimes:  []string{"Go go1.25.3", "Node.js v22.11.0", "Vitest 5.0.3"},
		Originals: []Original{{File: "junit-go.xml", MadeBy: "gotestsum --junitfile junit-go.xml -- ./...", StoredAt: "-"}, {File: "junit-vitest.xml", MadeBy: "vitest run --reporter=junit", StoredAt: "-"}},
		Testcases: testcases, Procedures: repo(t), Results: results(t),
		History: []string{Row("1.1.1", "2026-11-04", "1.2.0", "438", "testreport a1b2c3d")},
	})
}

// section answers the text of a section of a test report, up to the next one.
func section(md, heading string) string {
	_, after, ok := strings.Cut(md, "## "+heading+"\n")
	if !ok {
		return ""
	}
	text, _, _ := strings.Cut(after, "\n## ")
	return text
}

// @implements TC_136@1.2.1
func TestReadsBothJUnitAndLinksTheProcedureToItsTestCase(t *testing.T) {
	r := results(t)
	if len(r) != 4 {
		t.Fatalf("got %d results, want 4 (a subtest has none of its own): %+v", len(r), r)
	}
	execution := section(report(t), "Test execution")
	for _, want := range []string{
		"| SR_01@1.1.1 | TC_01@1.1.1 | TestOne | passed |",
		"| SR_03@2.1.1 | TC_03@1.1.1 | shows the rail | passed |",
	} {
		if !strings.Contains(execution, want) {
			t.Errorf("Test execution has no %q:\n%s", want, execution)
		}
	}
}

// @implements TC_137@2.3.1
func TestSectionsInTheirOrder(t *testing.T) {
	md := report(t)
	if !strings.HasPrefix(md, "# Test report for subsystem test (SwVal)\n") {
		t.Errorf("no head:\n%s", md)
	}
	last := -1
	for _, h := range []string{"Test environment", "Original data", "Test specification", "Test implementation", "Test execution", "Deviations", "Open", "Change history"} {
		i := strings.Index(md, "\n## "+h+"\n")
		if i <= last {
			t.Errorf("section %s missing or out of order", h)
		}
		last = i
	}
}

// @implements TC_138@1.2.1
func TestEnvironmentNamesEachRuntimeWithItsVersion(t *testing.T) {
	env := section(report(t), "Test environment")
	for _, want := range []string{"Label: 1.2.0", "Commit: a1b2c3d", "Date: 2026-11-04", "Operating system: linux/amd64", "Runner: GitHub Actions 2", "- Go go1.25.3", "- Node.js v22.11.0", "- Vitest 5.0.3"} {
		if !strings.Contains(env, want) {
			t.Errorf("Test environment has no %q:\n%s", want, env)
		}
	}
}

// @implements TC_139@2.2.1
func TestOriginalDataWithoutStore(t *testing.T) {
	data := section(report(t), "Original data")
	for _, want := range []string{
		"| File | Made by | Stored At |",
		"| junit-go.xml | `gotestsum --junitfile junit-go.xml -- ./...` | - |",
		"| junit-vitest.xml | `vitest run --reporter=junit` | - |",
	} {
		if !strings.Contains(data, want) {
			t.Errorf("Original data has no %q:\n%s", want, data)
		}
	}
}

// @implements TC_140@1.2.1
func TestSpecificationListsTheTestCaseWithItsRequirement(t *testing.T) {
	spec := section(report(t), "Test specification")
	if !strings.Contains(spec, "| Requirement | Test case | Text | Status |") || !strings.Contains(spec, "| SR_01@1.1.1 | TC_01@1.1.1 | one | approved |") {
		t.Errorf("got:\n%s", spec)
	}
}

// @implements TC_141@1.2.1
func TestImplementationWithAndWithoutATestCase(t *testing.T) {
	impl := section(report(t), "Test implementation")
	for _, want := range []string{
		"| SR_01@1.1.1 | TC_01@1.1.1 | TestOne | m/a/a_test.go |",
		"| SR_02@1.1.1 |  | TestTwo | m/a/a_test.go |",
	} {
		if !strings.Contains(impl, want) {
			t.Errorf("Test implementation has no %q:\n%s", want, impl)
		}
	}
	if strings.Contains(impl, "TestWithoutAnnotation") {
		t.Errorf("a test without annotation is listed:\n%s", impl)
	}
}

// @implements TC_142@1.2.1
func TestExecutionListsOnlyWhatRan(t *testing.T) {
	execution := section(report(t), "Test execution")
	if !strings.Contains(execution, "| SR_02@1.1.1 |  | TestTwo | failed |") {
		t.Errorf("no verdict of TestTwo:\n%s", execution)
	}
	if strings.Contains(execution, "TestNotRun") {
		t.Errorf("a test that did not run is listed:\n%s", execution)
	}
}

// @implements TC_143@2.1.1
func TestDeviationsAndOpen(t *testing.T) {
	md := report(t)
	deviations, open := section(md, "Deviations"), section(md, "Open")
	for _, want := range []string{"`TestTwo` (m/a/a_test.go) - failed", "`TestThree` (m/a/a_test.go) - skipped"} {
		if !strings.Contains(deviations, want) {
			t.Errorf("Deviations has no %q:\n%s", want, deviations)
		}
	}
	if strings.Contains(deviations, "TestOne") {
		t.Errorf("a passed test is a deviation:\n%s", deviations)
	}
	for _, want := range []string{
		"**Table open** - every test case and test procedure without a test result",
		"| Requirement | Test case | Test case status | Test procedure | Test procedure status |",
		"| SR_04@1.1.1 | TC_04@1.1.1 | approved | - | - |",
	} {
		if !strings.Contains(open, want) {
			t.Errorf("Open has no %q:\n%s", want, open)
		}
	}
}

// @implements TC_144@1.3.1
func TestDeviationsAndOpenEmpty(t *testing.T) {
	procedures := []Procedure{{Name: "TestOne", File: "a/a_test.go", Suite: "m/a", Implements: []string{"TC_01@1.1.1"}}}
	md := Markdown(Report{
		Level: "r5tc", Testcases: testcases[:1], Procedures: procedures,
		Results: []Result{{Suite: "m/a", Name: "TestOne", Verdict: "passed"}},
	})
	if section(md, "Deviations") != "\n`<empty>`\n" || section(md, "Open") != "\n`<empty>`\n" {
		t.Errorf("got:\n%s", md)
	}
}

// @implements TC_158@1.1.1
func TestOpenTestProcedureWithoutResult(t *testing.T) {
	procedures := []Procedure{{Name: "TestOne", File: "a/a_test.go", Suite: "m/a", Implements: []string{"TC_01@1.1.1"}}}
	open := section(Markdown(Report{Level: "r5tc", Testcases: testcases[:1], Procedures: procedures}), "Open")
	if !strings.Contains(open, "| SR_01@1.1.1 | TC_01@1.1.1 | approved | TestOne | no test result |") {
		t.Errorf("got:\n%s", open)
	}
	if strings.Contains(open, "| - | - |") {
		t.Errorf("the test case is listed without its test procedure as well:\n%s", open)
	}
}

// @implements TC_159@1.1.1
func TestOpenObsoleteTestCase(t *testing.T) {
	obsolete := []Testcase{{Key: "TC_05", Version: "1.1.1", Text: "five", Targets: []string{"SR_05@1.1.1"}, Status: "obsolete"}}
	open := section(Markdown(Report{Level: "r5tc", Testcases: obsolete}), "Open")
	if !strings.Contains(open, "| SR_05@1.1.1 | TC_05@1.1.1 | obsolete | - | - |") {
		t.Errorf("got:\n%s", open)
	}
}

// @implements TC_145@2.1.1
func TestRowWithLabelAndTicket(t *testing.T) {
	documents := []byte("- name: tools/testreports/software.md\n  type: document\n  object_version: '3.1.1'\n")
	next, err := create.Next(create.VersionOf(documents, "tools/testreports/software.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got := Row(next, "2026-11-04", "1.2.0", "438", "testreport a1b2c3d"); got != "| 4.1.1 | 2026-11-04 | 1.2.0 | #438 | testreport a1b2c3d |" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_146@2.1.1
func TestRowWithoutLabel(t *testing.T) {
	next, err := create.Next(create.VersionOf([]byte("- name: sbom.md\n"), "tools/testreports/software.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got := Row(next, "2026-11-04", "", "", "testreport a1b2c3d"); got != "| 1.1.1 | 2026-11-04 |  |  | testreport a1b2c3d |" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_147@1.2.1
func TestLabelWithoutTicketFails(t *testing.T) {
	if Flags("1.2.0", "") == nil {
		t.Error("no error for a label without a ticket")
	}
	if Flags("", "") != nil || Flags("1.2.0", "438") != nil {
		t.Error("an error where none is due")
	}
}

// @implements TC_148@2.1.1
func TestEntryAddedToDocuments(t *testing.T) {
	md := []byte("# Test report r5tc\n")
	got := string(create.EntryOf([]byte("- name: sbom.md\n  type: document\n"), "tools/testreports/software.md", "1.1.1", create.Hash(md)))
	want := "- name: tools/testreports/software.md\n  type: document\n  object_version: '1.1.1'\n  object_hash: '" + create.Hash(md) + "'\n  status: completed\n  verification_method: [generated]\n"
	if !strings.HasSuffix(got, want) {
		t.Errorf("got\n%s", got)
	}
}
