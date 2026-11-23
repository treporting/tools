package check

import (
	"testing"

	"github.com/treporting/tools/review"
)

// @implements TC_213@1.1.1
func TestAnIdIsLookedUpInItsOwnSubsystem(t *testing.T) {
	git := fakeGit{
		"a-other/requirements/software-requirements.yaml": "type: requirements\nrequirements:\n- object_key: SR_1\n  status: draft\n",
		"own/requirements/software-requirements.yaml":     "type: requirements\nrequirements:\n- object_key: SR_1\n  status: approved\n",
	}
	ix := review.NewIndex(func(file string) ([]byte, error) { return git.Show("", file) },
		[]string{"a-other/requirements/software-requirements.yaml", "own/requirements/software-requirements.yaml"})
	got := messages(Check(git, "b1", ix, []Annotation{{"own/trace/gaps.go", 3, "implements", []string{"SR_1"}, []string{"1.1.1"}}}))
	if len(got) != 0 {
		t.Errorf("got %q, want none", got)
	}
}
