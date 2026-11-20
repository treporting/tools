package write

import "strings"

// Document is what the object code report is written from.
type Document struct {
	// Label is the label of the baseline, empty outside one; Commit the commit of the repo; Date the
	// day it is written.
	Label, Commit, Date string
	// Sections are the sections of the build outputs, as Section writes them.
	Sections []string
	// History are the rows of its Change history, the row of this run last.
	History []string
}

// Markdown answers the object code report: the head of a document, the section Build, a section per
// build output, and the Change history.
//
// @implements SR_120@1.1.1, SR_121@1.1.1
func Markdown(d Document) string {

	var b strings.Builder
	b.WriteString("# Object code report\n\n")
	b.WriteString("Each function of the source code of each build output of the repo: whether it is in the\n")
	b.WriteString("output, the hash of its code there and its coverage by the tests. It is written by the tool\n")
	b.WriteString("objectcode (TOOL_35) and not edited by hand.\n\n")
	label := d.Label
	if label == "" {
		label = "-"
	}
	b.WriteString("## Build\n\n")
	b.WriteString("- Label: " + label + "\n")
	b.WriteString("- Commit: " + d.Commit + "\n")
	b.WriteString("- Date: " + d.Date + "\n")
	for _, s := range d.Sections {
		b.WriteString("\n" + strings.TrimSpace(s) + "\n")
	}
	b.WriteString("\n## Change history\n\n")
	b.WriteString("| Version | Date | Change | Reason | By |\n|---|---|---|---|---|\n")
	for _, h := range d.History {
		b.WriteString(h + "\n")
	}
	return b.String()
}
