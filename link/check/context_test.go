package check

import (
	"reflect"
	"testing"

	annotation "github.com/treporting/tools/annotation/check"
)

// Two subsystems of one repo both hold SR_1, at another x each; "a" sorts first, so the first file
// that holds SR_1 is the one of the other subsystem.
const otherSR = `type: requirements
requirements:
- object_key: SR_1
  object_version: '1.1.2'
`

const ownSR = `type: requirements
requirements:
- object_key: SR_1
  object_version: '2.1.1'
- object_key: SR_2
  object_version: '1.1.1'
  source_ids: [SR_1@2.1.1]
- object_key: SR_3
  object_version: '1.1.1'
  source_ids: [SR_1@1.1.2]
`

// @implements TC_208@1.1.1, TC_209@1.1.1
func TestALinkIsComparedWithTheObjectOfItsOwnSubsystem(t *testing.T) {
	got, compared, err := files{
		"a-other/requirements/software-requirements.yaml": otherSR,
		"own/requirements/software-requirements.yaml":     ownSR,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"SR_3 source_ids names SR_1@1.1.2: its object has 2.1.1"}
	if !reflect.DeepEqual(lines(got), want) || compared != 2 {
		t.Errorf("findings %q, compared %d, want %q and 2", lines(got), compared, want)
	}
}

// @implements TC_210@1.1.1
func TestATargetIsComparedWithTheObjectOfItsOwnSubsystem(t *testing.T) {
	got, _, err := files{
		"a-other/requirements/software-requirements.yaml": otherSR,
		"own/requirements/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  object_version: '2.1.1'
`,
		"own/testcases/software.yaml": `type: testcases
testcases:
- object_key: TC_1
  target_ids: [SR_1@2.1.1]
`,
	}.links()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("findings %q, want none", lines(got))
	}
}

// @implements TC_211@1.1.1
func TestAnAnnotationIsComparedWithTheObjectOfItsOwnSubsystem(t *testing.T) {
	f := files{
		"a-other/requirements/software-requirements.yaml": otherSR,
		"own/requirements/software-requirements.yaml":     ownSR,
	}
	got, compared := Annotations(f.read, f.index(), []annotation.Annotation{
		{File: "own/trace/gaps.go", Line: 3, Kind: "implements", IDs: []string{"SR_1"}, Versions: []string{"2.1.1"}},
	})
	if len(got) != 0 || compared != 1 {
		t.Errorf("findings %q, compared %d, want none and 1", lines(got), compared)
	}
}
