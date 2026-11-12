package write

import (
	"bufio"
	"fmt"
	"path"
	"strings"
)

// Original is a file of original data the report is written from.
type Original struct {
	File     string
	MadeBy   string // the command that made it
	StoredAt string // ./<folder>/<file>, or - where it is not stored
}

// Report is what a test report is written from.
type Report struct {
	Level      string     // the test level, such as r5tc
	Name       string     // the name of the test level with its abbreviation, such as subsystem test (SwVal)
	Label      string     // the label of the baseline; empty between baselines
	Commit     string     // the commit the tests ran on
	Date       string     // the day of the run
	System     string     // the operating system
	Runner     string     // the runner of the run
	Runtimes   []string   // each runtime and test runner the tests ran with, with its version
	Originals  []Original // each junit.xml given
	Testcases  []Testcase
	Procedures []Procedure
	Results    []Result
	History    []string // the rows of the Change history, the new one last
}

// row is one row of the tables Test implementation and Test execution: a test procedure with one
// test case it implements, or with none.
type row struct {
	requirement, testcase string
	procedure             Procedure
	index                 int // the place of the procedure in Report.Procedures
}

// Markdown answers the test report: the head of a document, then the sections Test environment,
// Original data, Test specification, Test implementation, Test execution, Deviations, Open and
// Change history.
//
// @implements SR_89@4.1.1
func Markdown(r Report) string {
	var b strings.Builder
	b.WriteString("# Test report for " + r.Name + "\n\n")
	b.WriteString("What the tests of the " + r.Name + " specify, implement and show in a run. It is\n")
	b.WriteString("written by the tool testreport (TOOL_34) and not edited by hand.\n\n")
	environment(&b, r)
	originals(&b, r.Originals)
	specification(&b, r.Testcases)
	rows := rows(r.Testcases, r.Procedures)
	implementation(&b, rows)
	verdicts := verdicts(r.Procedures, r.Results)
	execution(&b, rows, verdicts)
	deviations(&b, r.Procedures, verdicts)
	open(&b, r.Testcases, rows, verdicts)
	b.WriteString("\n## Change history\n\n")
	b.WriteString("| Version | Date | Change | Reason | By |\n|---|---|---|---|---|\n")
	for _, h := range r.History {
		b.WriteString(h + "\n")
	}
	return b.String()
}

// @implements SR_90@1.2.1
func environment(b *strings.Builder, r Report) {
	label := r.Label
	if label == "" {
		label = "-"
	}
	b.WriteString("## Test environment\n\n")
	b.WriteString("- Label: " + label + "\n")
	b.WriteString("- Commit: " + r.Commit + "\n")
	b.WriteString("- Date: " + r.Date + "\n")
	b.WriteString("- Operating system: " + r.System + "\n")
	b.WriteString("- Runner: " + r.Runner + "\n")
	for _, rt := range r.Runtimes {
		b.WriteString("- " + rt + "\n")
	}
}

// @implements SR_91@2.1.1
func originals(b *strings.Builder, originals []Original) {
	b.WriteString("\n## Original data\n\n")
	b.WriteString("| File | Made by | Stored At |\n|---|---|---|\n")
	for _, o := range originals {
		stored := o.StoredAt
		if stored == "" {
			stored = "-"
		}
		b.WriteString("| " + cell(o.File) + " | `" + cell(o.MadeBy) + "` | " + cell(stored) + " |\n")
	}
}

// @implements SR_92@1.2.1
func specification(b *strings.Builder, testcases []Testcase) {
	b.WriteString("\n## Test specification\n\n")
	b.WriteString("| Requirement | Test case | Text | Status |\n|---|---|---|---|\n")
	for _, t := range testcases {
		b.WriteString("| " + cell(strings.Join(t.Targets, ", ")) + " | " + t.Key + "@" + t.Version + " | " + cell(oneLine(t.Text)) + " | " + t.Status + " |\n")
	}
}

// rows answers a row for each test case a test procedure implements, with the requirement of that
// test case; a test procedure that implements none has one row with the requirements it verifies.
//
// @implements SR_93@1.2.1
func rows(testcases []Testcase, procedures []Procedure) []row {
	targets := map[string][]string{}
	for _, t := range testcases {
		targets[t.Key] = t.Targets
	}
	var out []row
	for i, p := range procedures {
		if len(p.Implements) == 0 {
			out = append(out, row{strings.Join(p.Verifies, ", "), "", p, i})
			continue
		}
		for _, tc := range p.Implements {
			key, _, _ := strings.Cut(tc, "@")
			out = append(out, row{strings.Join(targets[key], ", "), tc, p, i})
		}
	}
	return out
}

// @implements SR_93@1.2.1
func implementation(b *strings.Builder, rows []row) {
	b.WriteString("\n## Test implementation\n\n")
	b.WriteString("| Requirement | Test case | Test procedure | File |\n|---|---|---|---|\n")
	for _, r := range rows {
		b.WriteString("| " + cell(r.requirement) + " | " + r.testcase + " | " + cell(r.procedure.Name) + " | " + cell(r.procedure.File) + " |\n")
	}
}

// verdicts answers the verdict of each test procedure that ran, by its place in procedures.
func verdicts(procedures []Procedure, results []Result) map[int]string {
	out := map[int]string{}
	for i, p := range procedures {
		for _, r := range results {
			if same(p, r) {
				out[i] = worse(out[i], r.Verdict)
			}
		}
	}
	return out
}

// same says whether a result is the result of a test procedure: of Go by its package and name, of
// Vitest by its file and the name of its test at the end of the name of the result.
func same(p Procedure, r Result) bool {
	if strings.HasSuffix(p.File, ".go") {
		return r.Suite == p.Suite && r.Name == p.Name
	}
	file := r.Suite == p.File || strings.HasSuffix(p.File, "/"+path.Clean(r.Suite))
	return file && (r.Name == p.Name || strings.HasSuffix(r.Name, " > "+p.Name))
}

// worse answers the verdict that weighs more where a test procedure has two results.
func worse(a, b string) string {
	rank := map[string]int{"": 0, "passed": 1, "skipped": 2, "failed": 3}
	if rank[b] > rank[a] {
		return b
	}
	return a
}

// @implements SR_94@1.2.1
func execution(b *strings.Builder, rows []row, verdicts map[int]string) {
	b.WriteString("\n## Test execution\n\n")
	b.WriteString("| Requirement | Test case | Test procedure | Verdict |\n|---|---|---|---|\n")
	for _, r := range rows {
		if v, ok := verdicts[r.index]; ok {
			b.WriteString("| " + cell(r.requirement) + " | " + r.testcase + " | " + cell(r.procedure.Name) + " | " + v + " |\n")
		}
	}
}

// @implements SR_95@2.1.1
func deviations(b *strings.Builder, procedures []Procedure, verdicts map[int]string) {
	b.WriteString("\n## Deviations\n\n")
	n := 0
	for i, p := range procedures {
		if v, ok := verdicts[i]; ok && v != "passed" {
			fmt.Fprintf(b, "- `%s` (%s) - %s\n", p.Name, p.File, v)
			n++
		}
	}
	if n == 0 {
		b.WriteString("`<empty>`\n")
	}
}

// open lists in a table every test case and test procedure without a test result: a test case no
// test procedure implements, with - as test procedure; a test procedure that did not run, in each
// of its rows of Test implementation.
//
// @implements SR_95@2.1.1
func open(b *strings.Builder, testcases []Testcase, rows []row, verdicts map[int]string) {
	b.WriteString("\n## Open\n\n")
	status := map[string]string{}
	for _, t := range testcases {
		status[t.Key] = t.Status
	}
	implemented := map[string]bool{}
	for _, r := range rows {
		implemented[testcaseKey(r)] = true
	}
	var lines []string
	for _, t := range testcases {
		if !implemented[t.Key] {
			lines = append(lines, "| "+cell(strings.Join(t.Targets, ", "))+" | "+t.Key+"@"+t.Version+" | "+t.Status+" | - | - |\n")
		}
	}
	for _, r := range rows {
		if _, ok := verdicts[r.index]; !ok {
			lines = append(lines, "| "+cell(r.requirement)+" | "+r.testcase+" | "+status[testcaseKey(r)]+" | "+cell(r.procedure.Name)+" | no test result |\n")
		}
	}
	if len(lines) == 0 {
		b.WriteString("`<empty>`\n")
		return
	}
	b.WriteString("**Table open** - every test case and test procedure without a test result\n\n")
	b.WriteString("| Requirement | Test case | Test case status | Test procedure | Test procedure status |\n|---|---|---|---|---|\n")
	for _, l := range lines {
		b.WriteString(l)
	}
}

// testcaseKey answers the key of the test case of a row, without its version.
func testcaseKey(r row) string {
	key, _, _ := strings.Cut(r.testcase, "@")
	return key
}

// Row answers the row of the Change history that a run adds: the version of the test report after
// the run; the day it runs; the label; the ticket, as #<number>; the tool with its version. Change
// and Reason are empty where no label is given.
//
// @implements SR_96@2.1.1
func Row(version, date, label, ticket, by string) string {
	reason := ""
	if label != "" {
		reason = "#" + ticket
	}
	return "| " + version + " | " + date + " | " + cell(label) + " | " + cell(reason) + " | " + cell(by) + " |"
}

// History answers the rows of the Change history of a test report as it stands, without its head,
// so that a run keeps them. A file that is not there has none.
//
// @implements SR_96@2.1.1
func History(report []byte) []string {
	var rows []string
	in := false
	scanner := bufio.NewScanner(strings.NewReader(string(report)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "## "):
			in = line == "## Change history"
		case in && strings.HasPrefix(line, "| ") && !strings.HasPrefix(line, "| Version |"):
			rows = append(rows, line)
		}
	}
	return rows
}

// cell keeps a value from breaking the row of a table.
func cell(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
}

// oneLine keeps a text of several lines in one cell.
func oneLine(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// Flags answers an error where a label is given without a ticket: the row of the Change history
// would name a baseline without its reason.
//
// @implements SR_97@1.2.1
func Flags(label, ticket string) error {
	if label != "" && ticket == "" {
		return fmt.Errorf("-label is given without -ticket")
	}
	return nil
}
