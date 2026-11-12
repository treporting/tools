package create

import (
	"strings"
	"testing"
)

// @implements TC_119@1.1.1
func TestOwnPackage(t *testing.T) {
	got := Own("storages", "1.2.0", apacheText)
	want := Package{PURL: "pkg:github/treporting/storages@1.2.0", Name: "storages", Version: "1.2.0", Licence: "Apache-2.0"}
	if got != want {
		t.Errorf("got %+v", got)
	}
}

// @implements TC_109@1.1.1
func TestRowOfTheChangeHistory(t *testing.T) {
	if got := Row("2.1.1", "2026-10-31", "1.2.0", "343", "sbom a1b2c3d"); got != "| 2.1.1 | 2026-10-31 | 1.2.0 | #343 | sbom a1b2c3d |" {
		t.Errorf("got %q", got)
	}
}

// @verifies SR_40@3.3.1
func TestHistoryKeepsTheRowsOfTheChangeHistoryOnly(t *testing.T) {
	old := Markdown(SBOM{Repo: "storages",
		Packages: []Package{{PURL: "pkg:x", Licence: "MIT"}},
		History:  []string{"| 1.1.1 | 2026-10-01 | 1.0.0 | - | sbom 1111111 |", "| - | 2026-10-31 | 1.1.0 | - | sbom 2222222 |"}})
	got := History([]byte(old))
	if len(got) != 2 || got[0] != "| 1.1.1 | 2026-10-01 | 1.0.0 | - | sbom 1111111 |" || got[1] != "| - | 2026-10-31 | 1.1.0 | - | sbom 2222222 |" {
		t.Errorf("got %q", got)
	}
	if got := History(nil); len(got) != 0 {
		t.Errorf("no file: got %q", got)
	}
}

// @implements TC_111@2.1.1
func TestOneTableSortedByPURL(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", Packages: []Package{
		{PURL: "pkg:golang/stdlib@go1.25.14", Name: "stdlib", Version: "go1.25.14", Licence: "BSD-3-Clause"},
		{PURL: "pkg:github/treporting/storages@1.2.0", Name: "storages", Version: "1.2.0", Licence: "Apache-2.0"},
		{PURL: "pkg:apk/alpine/musl@1.2.6-r2", Name: "musl", Version: "1.2.6-r2", Licence: "MIT"},
	}})
	want := "| PURL | Name | Version | Licence |\n|---|---|---|---|\n" +
		"| pkg:apk/alpine/musl@1.2.6-r2 | musl | 1.2.6-r2 | MIT |\n" +
		"| pkg:github/treporting/storages@1.2.0 | storages | 1.2.0 | Apache-2.0 |\n" +
		"| pkg:golang/stdlib@go1.25.14 | stdlib | go1.25.14 | BSD-3-Clause |\n\n"
	_, packages, _ := strings.Cut(md, "\n## Packages\n")
	if !strings.Contains(packages, want) {
		t.Errorf("table not found in the section Packages of\n%s", md)
	}
	if strings.Count(md, "| PURL |") != 1 {
		t.Errorf("not one table:\n%s", md)
	}
}

// @implements TC_112@1.1.1
func TestEndsWithOpenAndChangeHistory(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", History: []string{"| - | 2026-10-31 | 1.2.0 | - | sbom a1b2c3d |"}})
	open, history := strings.Index(md, "\n## Open\n"), strings.Index(md, "\n## Change history\n")
	if open < 0 || history < open || strings.Index(md, "| PURL |") > open {
		t.Fatalf("sections out of order:\n%s", md)
	}
	if !strings.HasSuffix(md, "| - | 2026-10-31 | 1.2.0 | - | sbom a1b2c3d |\n") || strings.Count(md[open:], "\n## ") != 2 {
		t.Errorf("does not end with Open and Change history:\n%s", md)
	}
}

// @implements TC_160@1.1.1
func TestSectionsInOrder(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", Original: originals("/r")})
	var sections []string
	for _, line := range strings.Split(md, "\n") {
		if heading, ok := strings.CutPrefix(line, "## "); ok {
			sections = append(sections, heading)
		}
	}
	if strings.Join(sections, ", ") != "Original data, Packages, Open, Change history" {
		t.Errorf("sections %v:\n%s", sections, md)
	}
	if !strings.HasPrefix(md, "# Software Bill of Materials (SBOM)\n") {
		t.Errorf("title:\n%s", md)
	}
}

// @implements TC_117@1.1.1
func TestOpenListsEveryPackageOfUnknownLicence(t *testing.T) {
	md := Markdown(SBOM{Repo: "viewer-trace", Packages: []Package{
		{PURL: "pkg:npm/react@18.3.1", Licence: "MIT"},
		{PURL: "pkg:apk/alpine/nginx@1.30.5-r1", Licence: Unknown, Named: "2-clause BSD-like license"},
		{PURL: "pkg:github/treporting/viewer-trace@1.0.0", Licence: Unknown},
	}})
	want := "## Open\n\n" +
		"- `pkg:apk/alpine/nginx@1.30.5-r1` - the licence is unknown; the package names it \"2-clause BSD-like license\"\n" +
		"- `pkg:github/treporting/viewer-trace@1.0.0` - the licence is unknown\n\n"
	if !strings.Contains(md, want) {
		t.Errorf("got\n%s", md)
	}
}

// @implements TC_118@1.1.1
func TestOpenIsEmptyWhereEveryLicenceIsKnown(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", Packages: []Package{{PURL: "pkg:npm/react@18.3.1", Licence: "MIT"}}})
	if !strings.Contains(md, "## Open\n\n`<empty>`\n\n## Change history") {
		t.Errorf("got\n%s", md)
	}
}

// @verifies SR_41@1.2.1
func TestAValueDoesNotBreakTheTable(t *testing.T) {
	md := Markdown(SBOM{Repo: "r", Packages: []Package{{PURL: "pkg:x", Licence: "a|b"}}})
	if !strings.Contains(md, `| a\|b |`) {
		t.Errorf("got\n%s", md)
	}
}

func TestOnceKeepsTheFirstOfEachPURL(t *testing.T) {
	a := Package{PURL: "pkg:npm/a@1.0.0", Name: "a", Version: "1.0.0", Licence: "MIT"}
	b := Package{PURL: "pkg:npm/b@2.0.0", Name: "b", Version: "2.0.0", Licence: "ISC"}
	got := Once([]Package{a, b, a})
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("got %+v", got)
	}
}
