package check

import (
	"strings"

	"github.com/treporting/tools/review"
	"gopkg.in/yaml.v3"
)

// partKind is the kind of the elements a requirement of a level names in component; a level not
// here has no parts to name.
var partKind = map[string]string{"r3": "subsystem", "r6": "component", "r7": "module"}

// interfaceLevel tells the levels whose requirements name an interface in interface.
var interfaceLevel = map[string]bool{"r3": true, "r6": true}

// level is a file of r3, r6 or r7 as the check reads it: its elements and the parts its
// requirements name.
type level struct {
	Level    string `yaml:"level"`
	Elements []struct {
		Kind string `yaml:"kind"`
		Repo string `yaml:"repo"`
		Text string `yaml:"object_text"`
	} `yaml:"elements"`
	Requirements []struct {
		Key       string `yaml:"object_key"`
		Component string `yaml:"component"`
		Interface string `yaml:"interface"`
	} `yaml:"requirements"`
}

// Parts reads every file of r3, r6 and r7 the index knows and answers each requirement whose
// component names no part of its level, and each whose interface names no interface of its level.
// The parts and interfaces of a level are the elements of its file: on r3 a subsystem by its repo,
// else an element by its name, its object_text without backticks. A part no requirement names is
// no finding: a missing link blocks nothing.
//
// @implements SR_129@1.1.1, SR_130@1.1.1
func Parts(read Read, index *review.Index) []Finding {

	var findings []Finding
	for _, file := range index.Files() {
		data, err := read(file)
		if err != nil {
			continue
		}
		var l level
		if yaml.Unmarshal(data, &l) != nil {
			continue
		}
		kind, ok := partKind[l.Level]
		if !ok {
			continue
		}
		parts, interfaces := map[string]bool{}, map[string]bool{}
		for _, e := range l.Elements {
			switch {
			case e.Kind == kind && l.Level == "r3":
				parts[e.Repo] = true
			case e.Kind == kind:
				parts[name(e.Text)] = true
			case e.Kind == "interface":
				interfaces[name(e.Text)] = true
			}
		}
		by := "name"
		if l.Level == "r3" {
			by = "repo"
		}
		for _, r := range l.Requirements {
			//	a requirement without component: schemacheck reports it
			if r.Component != "" && !parts[r.Component] {
				findings = append(findings, Finding{r.Key, "component " + r.Component, "no " + kind + " of " + l.Level + " has this " + by})
			}
			if r.Interface != "" && interfaceLevel[l.Level] && !interfaces[r.Interface] {
				findings = append(findings, Finding{r.Key, "interface " + r.Interface, "no interface of " + l.Level + " has this name"})
			}
		}
	}
	return findings
}

// name answers the name of an element: its object_text without backticks.
func name(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "`", ""))
}
