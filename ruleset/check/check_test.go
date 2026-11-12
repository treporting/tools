package check

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const file = `{
  "name": "main is protected",
  "target": "branch",
  "enforcement": "active",
  "conditions": {"ref_name": {"exclude": [], "include": ["~DEFAULT_BRANCH"]}},
  "rules": [
    {"type": "deletion"},
    {"type": "required_status_checks", "parameters": {"required_status_checks": [
      {"context": "reviewcheck", "integration_id": 15368},
      {"context": "schemacheck", "integration_id": 15368}]}}
  ],
  "bypass_actors": []
}`

// github is a stand-in for the API of GitHub, in the test process: it answers the list of the
// rulesets and the one ruleset with what it is given.
func github(t *testing.T, status int, list, ruleset string) Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		switch r.URL.Path {
		case "/repos/o/r/rulesets":
			_, _ = w.Write([]byte(list))
		case "/repos/o/r/rulesets/7":
			_, _ = w.Write([]byte(ruleset))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return Client{API: server.URL, Repo: "o/r"}
}

const list = `[{"id": 7, "name": "main is protected"}]`

// onGitHub is the ruleset as GitHub answers it: the content of the file, with what GitHub sets itself.
func onGitHub(content string) string {
	return `{"id": 7, "source": "o/r", "created_at": "2026-10-23T12:42:42Z", "updated_at": "2026-10-26T18:00:00Z",
	"_links": {"self": {"href": "x"}}, "bypass_actors": [{"actor_id": 5, "actor_type": "RepositoryRole"}], ` + strings.TrimPrefix(strings.TrimSpace(content), "{")
}

func check(t *testing.T, c Client) []string {
	t.Helper()
	findings, err := Check(c, []byte(file))
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, f := range findings {
		out = append(out, f.String())
	}
	return out
}

// @implements TC_82@1.1.1
func TestTheSameRuleset(t *testing.T) {
	if got := check(t, github(t, 200, list, onGitHub(file))); len(got) != 0 {
		t.Errorf("want no finding, got %q", got)
	}
}

// @implements TC_83@1.1.1
func TestARequiredCheckMissingOnGitHub(t *testing.T) {
	changed := strings.Replace(file, `,
      {"context": "schemacheck", "integration_id": 15368}`, "", 1)
	got := check(t, github(t, 200, list, onGitHub(changed)))
	if len(got) != 1 || !strings.HasPrefix(got[0], "rules/required_status_checks: on GitHub") {
		t.Errorf("want one finding on rules/required_status_checks, got %q", got)
	}
}

// @implements TC_85@1.1.1
func TestAnotherEnforcementOnGitHub(t *testing.T) {
	disabled := strings.Replace(file, `"enforcement": "active"`, `"enforcement": "disabled"`, 1)
	got := check(t, github(t, 200, list, onGitHub(disabled)))
	if len(got) != 1 || !strings.HasPrefix(got[0], "enforcement") {
		t.Errorf("want one finding on enforcement, got %q", got)
	}
}

// @implements TC_84@1.1.1
func TestARuleOnlyOnGitHub(t *testing.T) {
	more := strings.Replace(file, `{"type": "deletion"},`, `{"type": "deletion"}, {"type": "non_fast_forward"},`, 1)
	got := check(t, github(t, 200, list, onGitHub(more)))
	if len(got) != 1 || got[0] != "rules/non_fast_forward: on GitHub, not in the file" {
		t.Errorf("want the rule only on GitHub, got %q", got)
	}
}

// @implements TC_86@1.1.1
func TestARulesetThatCannotBeRead(t *testing.T) {
	got := check(t, github(t, 500, list, onGitHub(file)))
	if len(got) != 1 || !strings.Contains(got[0], "cannot be read from GitHub: status 500") {
		t.Errorf("want a finding, not a pass; got %q", got)
	}
}

// @implements TC_87@1.1.1
func TestNoRulesetOfThatName(t *testing.T) {
	got := check(t, github(t, 200, `[{"id": 7, "name": "main protected"}]`, onGitHub(file)))
	if len(got) != 1 || got[0] != `no ruleset on GitHub is named "main is protected"` {
		t.Errorf("want a finding, not a pass; got %q", got)
	}
}
