package check

import "testing"

// @implements TC_212@1.1.1
func TestAnUplinkIsLookedUpInItsOwnSubsystem(t *testing.T) {
	got, _ := messages(t, files{
		"a-other/requirements/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  status: draft
`,
		"own/requirements/software-requirements.yaml": `type: requirements
requirements:
- object_key: SR_1
  status: approved
- object_key: SR_2
  source_ids: [SR_1@1.1.1]
  status: approved
`,
	})
	if len(got) != 0 {
		t.Errorf("findings %q, want none", got)
	}
}
