package check

import (
	"reflect"
	"testing"
)

const sbomMD = `# SBOM

**Table packages** - every package viewer-trace delivers

| PURL | Name | Version | Licence |
|---|---|---|---|
| pkg:apk/alpine/busybox@1.37.0-r30 | busybox | 1.37.0-r30 | GPL-2.0-only |
| pkg:apk/alpine/nginx@1.30.5-r1 | nginx | 1.30.5-r1 | unknown |
| pkg:github/treporting/viewer-trace@1.0.0 | viewer-trace | 1.0.0 | Apache-2.0 |
| pkg:golang/example.org/gpl@v1.0.0 | example.org/gpl | v1.0.0 | GPL-3.0-only |
| pkg:npm/react@18.3.1 | react | 18.3.1 | MIT |
| pkg:npm/words@1.0.0 | words | 1.0.0 | unknown |

## Open

- ` + "`pkg:apk/alpine/nginx@1.30.5-r1`" + ` - the licence is unknown

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| - | 2026-10-31 | 1.0.0 | - | sbom 249934f |
`

func TestPackagesOfTheTable(t *testing.T) {
	got := Packages([]byte(sbomMD))
	if len(got) != 6 || got[0] != (Package{PURL: "pkg:apk/alpine/busybox@1.37.0-r30", Licence: "GPL-2.0-only"}) {
		t.Errorf("got %+v", got)
	}
}

// @verifies SR_47@1.1.1
// @verifies SR_48@1.1.1
func TestCheckByTheKindOfPackage(t *testing.T) {
	got := Check(Packages([]byte(sbomMD)), nil)
	want := []Failure{
		{PURL: "pkg:apk/alpine/nginx@1.30.5-r1", Licence: "unknown"},
		{PURL: "pkg:golang/example.org/gpl@v1.0.0", Licence: "GPL-3.0-only"},
		{PURL: "pkg:npm/words@1.0.0", Licence: "unknown"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// @verifies SR_47@1.1.1
func TestFits(t *testing.T) {
	for expression, want := range map[string]bool{
		"Apache-2.0":                         true,
		"MIT":                                true,
		"BSD-2-Clause":                       true,
		"BSD-3-Clause":                       true,
		"ISC":                                true,
		"Zlib":                               true,
		"0BSD":                               true,
		"MPL-2.0":                            true,
		"Apache-2.0 AND MIT":                 true,
		"MIT AND GPL-2.0-only":               false,
		"MIT OR GPL-2.0-only":                true,
		"GPL-2.0-only OR LGPL-2.1-only":      false,
		"(MIT OR GPL-3.0-only) AND ISC":      true,
		"MIT OR GPL-3.0-only AND LGPL-3.0":   true,
		"(MIT OR GPL-3.0-only) AND LGPL-3.0": false,
		"Apache-2.0 WITH LLVM-exception":     false,
		"GPL-2.0-only":                       false,
		Unknown:                              false,
		"":                                   false,
		"MIT AND":                            false,
		"(MIT":                               false,
		"MIT)":                               false,
	} {
		if got := Fits(expression); got != want {
			t.Errorf("Fits(%q) = %v, want %v", expression, got, want)
		}
	}
}

// @implements TC_125@1.1.1, TC_130@1.1.1
func TestAssessedLicenceTakesThePlaceOfTheOneInSbomMD(t *testing.T) {
	assessed, err := ReadAssessed([]byte(`
- purl: pkg:apk/alpine/nginx@1.30.5-r1
  licence: BSD-2-Clause
  reason: the text of nginx is the BSD-2-Clause
- purl: pkg:npm/words@1.0.0
  licence: GPL-3.0-only
  reason: its LICENSE.txt is the GPL-3.0
`))
	if err != nil {
		t.Fatal(err)
	}
	got := Check(Packages([]byte(sbomMD)), assessed)
	want := []Failure{
		{PURL: "pkg:golang/example.org/gpl@v1.0.0", Licence: "GPL-3.0-only"},
		{PURL: "pkg:npm/words@1.0.0", Licence: "GPL-3.0-only"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// @verifies SR_49@1.1.1
func TestAssessedEntryNeedsPURLLicenceAndReason(t *testing.T) {
	for _, file := range []string{
		"- purl: pkg:npm/a@1\n  licence: MIT\n",
		"- purl: pkg:npm/a@1\n  reason: read\n",
		"- licence: MIT\n  reason: read\n",
		"purl: x",
	} {
		if _, err := ReadAssessed([]byte(file)); err == nil {
			t.Errorf("no error for %q", file)
		}
	}
}

// @implements TC_132@1.1.1
func TestFailureNamesPURLAndLicence(t *testing.T) {
	if got := (Failure{PURL: "pkg:npm/words@1.0.0", Licence: "unknown"}).String(); got != "pkg:npm/words@1.0.0 unknown" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_124@1.1.1, TC_125@1.1.1, TC_126@1.1.1, TC_127@1.1.1, TC_128@1.1.1, TC_129@1.1.1, TC_131@1.1.1
func TestEachTypeOfPURLWithItsLicence(t *testing.T) {
	assessed := map[string]string{"pkg:apk/alpine/other@1.0.0": "MIT"}
	for _, c := range []struct {
		purl, licence string
		fails         bool
	}{
		{"pkg:golang/example.org/a@v1.0.0", "MIT", false},
		{"pkg:npm/b@1.0.0", "GPL-3.0-only", true},
		{"pkg:github/treporting/c@1.0.0", "MIT AND GPL-3.0-only", true},
		{"pkg:golang/example.org/d@v1.0.0", "MIT OR GPL-3.0-only", false},
		{"pkg:apk/alpine/e@1.0.0", "unknown", true},
		{"pkg:apk/alpine/f@1.0.0", "GPL-2.0-only", false},
	} {
		got := Check([]Package{{PURL: c.purl, Licence: c.licence}}, assessed)
		if fails := len(got) == 1 && got[0] == (Failure{PURL: c.purl, Licence: c.licence}); fails != c.fails || (!c.fails && len(got) != 0) {
			t.Errorf("%s %s: got %+v, want fails %v", c.purl, c.licence, got, c.fails)
		}
	}
}
