// Package check checks that every id an annotation names is approved in the base commit of the
// pull request (TOOL_06), and that it names its version (TOOL_20). It reads @refines in the
// $comment of the schemas as well.
package check

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/treporting/specs-and-deploy/tools/review"
	"gopkg.in/yaml.v3"
)

// Git is what the check asks of the history of the repo.
type Git interface {
	// Show answers the content of path in the commit rev.
	Show(rev, path string) ([]byte, error)
}

// Annotation is one @implements, @verifies or @refines, and the ids it names.
type Annotation struct {
	File     string
	Line     int
	Kind     string   // implements, verifies or refines
	IDs      []string // the ids alone
	Versions []string // the version each id names after its @, "" where it names none
}

// Finding is one id an annotation names that is missing or not approved in the base commit, or
// that names no version, or a file whose kind is set neither as code nor as no code.
type Finding struct {
	Annotation Annotation
	ID         string
	Message    string
}

func (f Finding) String() string {
	if f.Annotation.Kind == "" {
		return f.Annotation.File + ": " + f.Message
	}
	return fmt.Sprintf("%s:%d @%s %s: %s", f.Annotation.File, f.Annotation.Line, f.Annotation.Kind, f.ID, f.Message)
}

// Kinds are the two lists of tools/annotation/kinds.yaml: the kinds of files that are code, read
// for annotations, and those that are no code, never read (formats.md, "Where an annotation is
// read").
type Kinds struct {
	Whitelist []string `yaml:"whitelist"`
	Blacklist []string `yaml:"blacklist"`
}

// ReadKinds reads the two lists from the file path.
func ReadKinds(path string) (Kinds, error) {
	var k Kinds
	data, err := os.ReadFile(path)
	if err != nil {
		return k, err
	}
	return k, yaml.Unmarshal(data, &k)
}

// kindOf answers whether a file is code to read, no code, or of a kind on neither list - by its
// name alone.
func (k Kinds) kindOf(file string) string {
	name := filepath.Base(file)
	matches := func(patterns []string) bool {
		for _, p := range patterns {
			if ok, _ := filepath.Match(p, name); ok {
				return true
			}
		}
		return false
	}
	switch {
	case matches(k.Blacklist):
		return "no code"
	case matches(k.Whitelist):
		return "code"
	}
	return "unknown"
}

// annotation is a line of comment that begins with an annotation: after spaces, a comment mark of
// one to three signs other than a quote, then @implements or @verifies and its ids, each with what
// follows its @, if anything.
var annotation = regexp.MustCompile("^\\s*[^\\w\\s\"'`]{1,3}\\s*@(implements|verifies)\\s+(" + ids + ")")

// refines is the $comment of a schema that begins with @refines and its ids.
var refines = regexp.MustCompile("^\\s*@(refines)\\s+(" + ids + ")")

// ids is a list of ids, each with what follows its @, if anything.
const ids = "[A-Z][A-Za-z0-9]*_[A-Za-z0-9_-]*[0-9](?:@[^\\s,]*)?(?:\\s*,\\s*[A-Z][A-Za-z0-9]*_[A-Za-z0-9_-]*[0-9](?:@[^\\s,]*)?)*"

// version is the form of a version: x.y.z.
var version = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// Find answers every annotation in the files, relative to root, whose kind is code, and a finding
// for every file whose kind is on neither list. A YAML that names a schema is our data and is not
// read. A file is read only once its name has told it is code.
//
// @implements SR_24@3.1.1, SR_27@1.1.1
func Find(root string, files []string, kinds Kinds) ([]Annotation, []Finding, error) {

	var found []Annotation
	var unknown []Finding
	for _, file := range files {
		switch kinds.kindOf(file) {
		case "no code":
			continue
		case "unknown":
			unknown = append(unknown, Finding{Annotation{File: file}, "", "its kind is set neither as code nor as no code (tools/annotation/kinds.yaml)"})
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			return nil, nil, err
		}
		if isData(file, data) {
			continue
		}
		for i, line := range strings.Split(string(data), "\n") {
			if m := annotation.FindStringSubmatch(line); m != nil {
				found = append(found, named(file, i+1, m))
			}
		}
	}
	return found, unknown, nil
}

// Refines answers the @refines in the $comment of each schema among the files, relative to root: a
// JSON in schemas/. Its line is the line of the $comment.
//
// @implements SR_35@1.1.1
func Refines(root string, files []string) ([]Annotation, error) {

	var found []Annotation
	for _, file := range files {
		if filepath.Dir(file) != "schemas" || filepath.Ext(file) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			return nil, err
		}
		var schema struct {
			Comment string `json:"$comment"`
		}
		if json.Unmarshal(data, &schema) != nil {
			continue
		}
		m := refines.FindStringSubmatch(schema.Comment)
		if m == nil {
			continue
		}
		line := 0
		for i, l := range strings.Split(string(data), "\n") {
			if strings.Contains(l, `"$comment"`) {
				line = i + 1
				break
			}
		}
		found = append(found, named(file, line, m))
	}
	return found, nil
}

// named answers the annotation of a match: its kind, then its ids, each with its version.
func named(file string, line int, m []string) Annotation {
	var ids, versions []string
	for _, n := range strings.Split(m[2], ",") {
		id, v, _ := strings.Cut(strings.TrimSpace(n), "@")
		ids = append(ids, id)
		versions = append(versions, v)
	}
	return Annotation{file, line, m[1], ids, versions}
}

// isData says whether a YAML is our data: one that names its schema in its field schema.
func isData(file string, content []byte) bool {
	if !strings.HasSuffix(file, ".yaml") && !strings.HasSuffix(file, ".yml") {
		return false
	}
	var doc struct {
		Schema string `yaml:"schema"`
	}
	return yaml.Unmarshal(content, &doc) == nil && doc.Schema != ""
}

// Unversioned answers each id the annotations name without a version as x.y.z after its @.
//
// @implements SR_30@2.2.1
func Unversioned(annotations []Annotation) []Finding {

	var findings []Finding
	for _, a := range annotations {
		for i, id := range a.IDs {
			if !version.MatchString(a.Versions[i]) {
				findings = append(findings, Finding{a, id, fmt.Sprintf("names no version, as %s@x.y.z", id)})
			}
		}
	}
	return findings
}

// Check answers each id the annotations name that is not in the requirements of the base commit,
// or not approved there. The index tells the file of an id, as the working tree holds it.
//
// @implements SR_24@3.1.1, SR_25@2.2.1
func Check(git Git, base string, index *review.Index, annotations []Annotation) []Finding {

	var findings []Finding
	for _, a := range annotations {
		for _, id := range a.IDs {
			if message := status(git, base, index, id); message != "" {
				findings = append(findings, Finding{a, id, message})
			}
		}
	}
	return findings
}

// status answers why id does not count in the base commit, or "" where it is approved there.
func status(git Git, base string, index *review.Index, id string) string {

	file, key := index.FileOf(id)
	if file == "" || key == "" {
		return "no file of objects holds such an id"
	}
	data, err := git.Show(base, file)
	if err != nil {
		return fmt.Sprintf("%s is not in the base commit %.7s", file, base)
	}
	object := review.ObjectIn(data, key)
	if object == nil {
		return fmt.Sprintf("not in %s in the base commit %.7s", file, base)
	}
	var o struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal(object, &o); err != nil {
		return err.Error()
	}
	if o.Status != "approved" {
		return fmt.Sprintf("%s in the base commit %.7s, not approved", o.Status, base)
	}
	return ""
}
