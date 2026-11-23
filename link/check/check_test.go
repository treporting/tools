package check

import (
	"os"
	"reflect"
	"sort"
	"testing"

	annotation "github.com/treporting/tools/annotation/check"
	"github.com/treporting/tools/review"
)

// files holds the files of requirements in memory.
type files map[string]string

func (f files) read(file string) ([]byte, error) {
	if c, ok := f[file]; ok {
		return []byte(c), nil
	}
	return nil, os.ErrNotExist
}

// links makes the link check on the files.
func (f files) links() ([]Finding, int, error) {
	return Links(f.read, f.index())
}

// index makes the index of the files, in the order of their names.
func (f files) index() *review.Index {
	var names []string
	for file := range f {
		names = append(names, file)
	}
	sort.Strings(names)
	return review.NewIndex(f.read, names)
}

func lines(findings []Finding) []string {
	var out []string
	for _, f := range findings {
		out = append(out, f.String())
	}
	return out
}

// TOOL_9 has moved on in x, TOOL_1 in y and z only.
const tools = `schema: schemas/tool-requirements.schema.json
type: requirements
level: tool
requirements:
- object_key: TOOL_1
  object_version: '1.2.3'
- object_key: TOOL_9
  object_version: '4.1.1'
`

// @implements TC_59@1.1.1, TC_60@1.1.1
func TestOnlyAnotherXIsReported(t *testing.T) {
	got, compared, err := files{
		"requirements/tool-requirements.yaml": tools,
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  source_ids: [TOOL_1@1.1.1, TOOL_9@3.1.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SR_1 source_ids names TOOL_9@3.1.1: its object has 4.1.1"}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
	if compared != 2 {
		t.Errorf("compared %d, want 2", compared)
	}
}

// @implements TC_61@1.1.1, TC_62@1.1.1
func TestTheAllocationAndControlsAreCompared(t *testing.T) {
	got, compared, err := files{
		"requirements/system-requirements.yaml": `type: requirements
requirements:
- object_key: SYS_9
  object_version: '4.1.2'
`,
		"requirements/user-needs.yaml": `type: requirements
requirements:
- object_key: UN_PSF_1
  object_version: '2.1.1'
`,
		"architecture/system-architecture.yaml": `type: architecture
elements:
- object_key: ARC_1
  allocated: [SYS_9@3.1.2]
`,
		"safety/hazards.yaml": `type: hazards
controls:
- object_key: HAZC_1
  controls: [UN_PSF_1@1.1.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"ARC_1 allocated names SYS_9@3.1.2: its object has 4.1.2",
		"HAZC_1 controls names UN_PSF_1@1.1.1: its object has 2.1.1",
	}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
	if compared != 2 {
		t.Errorf("compared %d, want 2", compared)
	}
}

// @implements TC_67@1.1.1, TC_68@1.1.1
func TestATargetOfATestCaseIsCompared(t *testing.T) {
	got, compared, err := files{
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  object_version: '2.1.1'
`,
		"tools/testcases/software.yaml": `type: testcases
testcases:
- object_key: TC_1
  target_ids: [SR_1@1.1.1]
- object_key: TC_2
  target_ids: [SR_1@2.2.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"TC_1 target_ids names SR_1@1.1.1: its object has 2.1.1"}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
	if compared != 2 {
		t.Errorf("compared %d, want 2", compared)
	}
}

// A link to an object no file holds is not compared; uplinkcheck reports it.
//
// @verifies SR_32@1.1.1
func TestAMissingObjectIsNotCompared(t *testing.T) {
	got, compared, err := files{
		"requirements/tool-requirements.yaml": tools,
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  source_ids: [TOOL_99@1.1.1, TOOL_9]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 || compared != 0 {
		t.Errorf("findings %q, compared %d, want none", lines(got), compared)
	}
}

// @implements TC_64@1.1.1, TC_65@1.1.1
func TestAnAnnotationWithAnotherXIsReported(t *testing.T) {
	got, compared := Annotations(files{"requirements/tool-requirements.yaml": tools}.read, files{"requirements/tool-requirements.yaml": tools}.index(), []annotation.Annotation{
		{File: "a.go", Line: 3, Kind: "implements", IDs: []string{"TOOL_1", "TOOL_9"}, Versions: []string{"1.1.1", "3.1.1"}},
		{File: "a_test.go", Line: 7, Kind: "verifies", IDs: []string{"TOOL_9"}, Versions: []string{"4.1.3"}},
		{File: "b.go", Line: 1, Kind: "implements", IDs: []string{"TOOL_9"}, Versions: []string{""}},
	})
	want := []string{"a.go:3 @implements names TOOL_9@3.1.1: its object has 4.1.1"}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
	if compared != 3 {
		t.Errorf("compared %d, want 3", compared)
	}
}

// @implements TC_66@1.1.1
func TestARefinesWithAnotherXIsReported(t *testing.T) {
	architecture := files{"architecture/system-architecture.yaml": `type: architecture
elements:
- object_key: ARC_22
  object_version: '3.1.1'
`}
	got, compared := Annotations(architecture.read, architecture.index(), []annotation.Annotation{
		{File: "schemas/a.schema.json", Line: 4, Kind: "refines", IDs: []string{"ARC_22"}, Versions: []string{"2.1.1"}},
	})
	want := []string{"schemas/a.schema.json:4 @refines names ARC_22@2.1.1: its object has 3.1.1"}
	if !reflect.DeepEqual(lines(got), want) || compared != 1 {
		t.Errorf("findings %q, compared %d, want %q and 1", lines(got), compared, want)
	}
}

// A link to a regulatory requirement is compared with the version in the file of its module.
//
// @implements TC_63@1.1.1
func TestALinkToAModuleIsCompared(t *testing.T) {
	got, compared, err := files{
		"regulatory/DO178C.yaml": `type: normative
requirements:
- object_key: REG_DO178C_2
  object_version: '2.1.1'
`,
		"requirements/system-requirements.yaml": `type: requirements
requirements:
- object_key: SYS_9
  source_ids: [REG_DO178C_2@1.1.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SYS_9 source_ids names REG_DO178C_2@1.1.1: its object has 2.1.1"}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
	if compared != 1 {
		t.Errorf("compared %d, want 1", compared)
	}
}

// repoAbove holds TOOL_9 at 4.1.1, as the repo above.
func repoAbove() *review.Repo {
	f := files{"requirements/tool-requirements.yaml": tools}
	return &review.Repo{Read: f.read, Index: f.index()}
}

// @implements TC_168@1.1.1
func TestALinkIntoTheRepoAboveIsCompared(t *testing.T) {
	own := files{"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  source_ids: [TOOL_9@3.1.1]
`}
	got, compared, err := LinksWithAbove(own.read, own.index(), repoAbove())
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SR_1 source_ids names TOOL_9@3.1.1: its object has 4.1.1"}
	if !reflect.DeepEqual(lines(got), want) || compared != 1 {
		t.Errorf("findings %q, compared %d, want %q and 1", lines(got), compared, want)
	}
	got, compared, err = own.links()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 || compared != 0 {
		t.Errorf("without the repo above: findings %q, compared %d, want none", lines(got), compared)
	}
}

// @implements TC_169@1.1.1
func TestAnAnnotationIntoTheRepoAboveIsCompared(t *testing.T) {
	own := files{}
	got, compared := AnnotationsWithAbove(own.read, own.index(), repoAbove(), []annotation.Annotation{
		{File: "a.go", Line: 3, Kind: "implements", IDs: []string{"TOOL_9"}, Versions: []string{"3.1.1"}},
	})
	want := []string{"a.go:3 @implements names TOOL_9@3.1.1: its object has 4.1.1"}
	if !reflect.DeepEqual(lines(got), want) || compared != 1 {
		t.Errorf("findings %q, compared %d, want %q and 1", lines(got), compared, want)
	}
}

// obsoleteTools holds TOOL_9 obsolete and TOOL_1 approved, both at 1.1.1.
const obsoleteTools = `schema: schemas/tool-requirements.schema.json
type: requirements
level: tool
requirements:
- object_key: TOOL_1
  object_version: '1.1.1'
  status: approved
- object_key: TOOL_9
  object_version: '1.1.1'
  status: obsolete
`

// A link, an entry of the allocation and of controls, and an annotation that name an obsolete
// object are reported, though their version matches.
//
// @implements TC_197@1.1.1, TC_198@1.1.1
func TestALinkToAnObsoleteObjectIsReported(t *testing.T) {
	f := files{
		"requirements/tool-requirements.yaml": obsoleteTools,
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  source_ids: [TOOL_1@1.1.1, TOOL_9@1.1.1]
`,
		"architecture/system-architecture.yaml": `type: architecture
elements:
- object_key: ARC_1
  allocated: [TOOL_9@1.1.1]
`,
		"safety/hazards.yaml": `type: hazards
controls:
- object_key: HAZC_1
  controls: [TOOL_9@1.1.1]
`,
		"tools/testcases/software.yaml": `type: testcases
testcases:
- object_key: TC_1
  target_ids: [TOOL_9@1.1.1]
`,
	}
	got, compared, err := f.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"ARC_1 allocated names TOOL_9@1.1.1: its object is obsolete",
		"HAZC_1 controls names TOOL_9@1.1.1: its object is obsolete",
		"SR_1 source_ids names TOOL_9@1.1.1: its object is obsolete",
		"TC_1 target_ids names TOOL_9@1.1.1: its object is obsolete",
	}
	if !reflect.DeepEqual(lines(got), want) || compared != 5 {
		t.Errorf("findings %q, compared %d, want %q and 5", lines(got), compared, want)
	}
	ids, _ := Annotations(f.read, f.index(), []annotation.Annotation{
		{File: "a.go", Line: 3, Kind: "implements", IDs: []string{"TOOL_1", "TOOL_9"}, Versions: []string{"1.1.1", "1.1.1"}},
	})
	want = []string{"a.go:3 @implements names TOOL_9@1.1.1: its object is obsolete"}
	if !reflect.DeepEqual(lines(ids), want) {
		t.Errorf("annotations: findings %q, want %q", lines(ids), want)
	}
}

// A link to an obsolete object in the repo above is reported.
//
// @implements TC_199@1.1.1
func TestALinkToAnObsoleteObjectAboveIsReported(t *testing.T) {
	above := files{"requirements/tool-requirements.yaml": obsoleteTools}
	own := files{"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  source_ids: [TOOL_9@1.1.1]
`}
	got, _, err := LinksWithAbove(own.read, own.index(), &review.Repo{Read: above.read, Index: above.index()})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SR_1 source_ids names TOOL_9@1.1.1: its object is obsolete"}
	if !reflect.DeepEqual(lines(got), want) {
		t.Errorf("findings %q, want %q", lines(got), want)
	}
}

// The links of an obsolete object are not compared: its link to TOOL_9, which has moved on in x,
// is no finding, and the link of SR_2, approved, is.
//
// @implements TC_200@1.1.1
func TestTheLinksOfAnObsoleteObjectAreNotCompared(t *testing.T) {
	got, compared, err := files{
		"requirements/tool-requirements.yaml": tools,
		"tools/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  status: obsolete
  source_ids: [TOOL_9@3.1.1]
- object_key: SR_2
  status: approved
  source_ids: [TOOL_9@3.1.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SR_2 source_ids names TOOL_9@3.1.1: its object has 4.1.1"}
	if !reflect.DeepEqual(lines(got), want) || compared != 1 {
		t.Errorf("findings %q, compared %d, want %q and 1", lines(got), compared, want)
	}
}
