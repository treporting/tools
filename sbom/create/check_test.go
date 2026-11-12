package create

import (
	"strings"
	"testing"
)

var listed = []Package{
	{PURL: "pkg:github/treporting/storages@1.2.0", Name: "storages", Version: "1.2.0", Licence: "Apache-2.0"},
	{PURL: "pkg:npm/a%7Cb@1.0.0", Name: "a|b", Version: "1.0.0", Licence: "MIT"},
	{PURL: "pkg:npm/left-pad@1.3.0", Name: "left-pad", Version: "1.3.0", Licence: Unknown, Named: "WTFPL"},
}

// @verifies SR_66@1.1.1
func TestTableAsSbomMDListsIt(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", Packages: listed, History: []string{"| 1.1.1 | 2026-11-01 | 1.2.0 | #343 | sbom a1b2c3d |"}})
	got := Table([]byte(md))
	if len(got) != len(listed) {
		t.Fatalf("got %+v", got)
	}
	for i, p := range listed {
		p.Named = ""
		if got[i] != p {
			t.Errorf("row %d: got %+v, want %+v", i, got[i], p)
		}
	}
}

// @verifies SR_66@1.1.1
func TestLabelOfTheLastRow(t *testing.T) {
	md := Markdown(SBOM{Repo: "storages", History: []string{
		"| 1.1.1 | 2026-10-01 | 1.1.0 | #300 | sbom 1111111 |",
		"| 2.1.1 | 2026-11-01 | 1.2.0 | #343 | sbom 2222222 |",
	}})
	if got := Label([]byte(md)); got != "1.2.0" {
		t.Errorf("got %q", got)
	}
	if got := Label([]byte(Markdown(SBOM{Repo: "storages"}))); got != "" {
		t.Errorf("no row: got %q", got)
	}
}

// @implements TC_123@1.1.1
func TestDifferNamesEachPackageThatDiffers(t *testing.T) {
	if got := Differ(listed, listed); len(got) != 0 {
		t.Errorf("same lists: %q", got)
	}
	delivered := []Package{
		listed[0],
		{PURL: "pkg:npm/a%7Cb@1.0.0", Name: "a|b", Version: "1.0.0", Licence: "ISC"},
		{PURL: "pkg:golang/new@v1.0.0", Name: "new", Version: "v1.0.0", Licence: "MIT"},
	}
	got := strings.Join(Differ(delivered, listed), "\n")
	for _, want := range []string{
		"pkg:golang/new@v1.0.0: delivered, not in sbom.md",
		"pkg:npm/a%7Cb@1.0.0: sbom.md lists a|b, 1.0.0, MIT; delivered is a|b, 1.0.0, ISC",
		"pkg:npm/left-pad@1.3.0: in sbom.md, not delivered",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "storages") {
		t.Errorf("names a package that fits:\n%s", got)
	}
}
