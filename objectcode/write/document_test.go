package write

import (
	"strings"
	"testing"
)

func TestDocument(t *testing.T) {

	md := Markdown(Document{
		Label:    "T047.0",
		Commit:   "7ace229",
		Date:     "2026-11-20",
		Sections: []string{"## a\n\n- Binary: a (cmd/a)\n", "## b.js\n\n- Bundle: b.js (web)\n"},
		History:  []string{"| 1.1.1 | 2026-11-20 | T047.0 | #37 | objectcode 7ace229 |"},
	})
	for _, want := range []string{
		"# Object code report\n",
		"## Build\n\n- Label: T047.0\n- Commit: 7ace229\n- Date: 2026-11-20\n\n## a\n",
		"- Binary: a (cmd/a)\n\n## b.js\n",
		"## Change history\n\n| Version | Date | Change | Reason | By |\n|---|---|---|---|---|\n| 1.1.1 | 2026-11-20 | T047.0 | #37 | objectcode 7ace229 |\n",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("the report has not %q:\n%s", want, md)
		}
	}
	if !strings.Contains(Markdown(Document{Commit: "7ace229"}), "- Label: -\n") {
		t.Error("a report without a label names none")
	}
}
