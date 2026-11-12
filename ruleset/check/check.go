// Package check compares the ruleset of main on GitHub with the file in the repo (TOOL_13).
package check

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

// Finding is one difference between GitHub and the file, or a ruleset that cannot be checked.
type Finding struct {
	Field   string
	Message string
}

func (f Finding) String() string {
	if f.Field == "" {
		return f.Message
	}
	return f.Field + ": " + f.Message
}

// compared are the fields of a ruleset that are compared. What GitHub sets itself - id, times,
// links - is left out, and bypass_actors, which cannot be read without the rights of an admin.
var compared = []string{"target", "enforcement", "conditions"}

// Client reads a repo of GitHub through its API.
type Client struct {
	API   string // such as https://api.github.com
	Repo  string // such as treporting/specs-and-deploy
	Token string // may be empty; with it the limit of requests is higher
	HTTP  *http.Client
}

// Check reads the ruleset of main from GitHub, found by the name in want, the content of the file,
// and answers each difference. A ruleset that cannot be read, or no ruleset of that name, is a
// finding too, not a pass.
//
// @implements SR_21@1.1.1, SR_22@1.1.1
func Check(c Client, want []byte) ([]Finding, error) {

	var file map[string]any
	if err := json.Unmarshal(want, &file); err != nil {
		return nil, fmt.Errorf("the file: %w", err)
	}
	name, _ := file["name"].(string)

	var list []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := c.get("/repos/"+c.Repo+"/rulesets", &list); err != nil {
		return []Finding{{"", "the rulesets cannot be read from GitHub: " + err.Error()}}, nil
	}
	id := int64(-1)
	for _, r := range list {
		if r.Name == name {
			id = r.ID
		}
	}
	if id < 0 {
		return []Finding{{"", fmt.Sprintf("no ruleset on GitHub is named %q", name)}}, nil
	}

	var github map[string]any
	if err := c.get(fmt.Sprintf("/repos/%s/rulesets/%d", c.Repo, id), &github); err != nil {
		return []Finding{{"", "the ruleset cannot be read from GitHub: " + err.Error()}}, nil
	}

	var findings []Finding
	for _, field := range compared {
		if !reflect.DeepEqual(github[field], file[field]) {
			findings = append(findings, Finding{field, "on GitHub " + compact(github[field]) + ", in the file " + compact(file[field])})
		}
	}
	findings = append(findings, compareRules(github["rules"], file["rules"])...)
	return findings, nil
}

// compareRules answers each rule that is only on GitHub, only in the file, or different in both,
// by the type of the rule.
//
// @implements SR_21@1.1.1
func compareRules(github, file any) []Finding {

	byType := func(rules any) (map[string]any, []string) {
		m := map[string]any{}
		var order []string
		list, _ := rules.([]any)
		for _, r := range list {
			rule, _ := r.(map[string]any)
			t, _ := rule["type"].(string)
			m[t] = rule
			order = append(order, t)
		}
		return m, order
	}
	gh, ghOrder := byType(github)
	fl, flOrder := byType(file)

	var findings []Finding
	for _, t := range flOrder {
		g, ok := gh[t]
		switch {
		case !ok:
			findings = append(findings, Finding{"rules/" + t, "in the file, not on GitHub"})
		case !reflect.DeepEqual(g, fl[t]):
			findings = append(findings, Finding{"rules/" + t, "on GitHub " + compact(g) + ", in the file " + compact(fl[t])})
		}
	}
	for _, t := range ghOrder {
		if _, ok := fl[t]; !ok {
			findings = append(findings, Finding{"rules/" + t, "on GitHub, not in the file"})
		}
	}
	return findings
}

// get reads path from the API into v.
func (c Client) get(path string, v any) error {

	client := c.HTTP
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	req, err := http.NewRequest(http.MethodGet, c.API+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return json.Unmarshal(body, v)
}

func compact(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
