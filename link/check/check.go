// Package check makes the link check (TOOL_16): it compares the version a link, a target, an entry
// of the allocation or of controls, or an annotation names with the version its object has, and reports
// each whose object has another x. It reports and does not decide (definitions.md, section Link
// check and impact assessment).
package check

import (
	"fmt"
	"regexp"
	"strings"

	annotation "github.com/treporting/specs-and-deploy/tools/annotation/check"
	"github.com/treporting/specs-and-deploy/tools/review"
	"gopkg.in/yaml.v3"
)

// Read answers the content of a file of the repo, by its path from the root.
type Read func(file string) ([]byte, error)

// Finding is one link whose object has another x than the version it names.
type Finding struct {
	Where string // the object and its field, or the file and line of an annotation
	Link  string // the link as it names it, with its version
	Has   string // the version its object has
}

func (f Finding) String() string {
	return f.Where + " names " + f.Link + ": its object has " + f.Has
}

// fields are the fields of an object that name other objects with their version: the uplinks,
// the targets of a test case, the allocation (r4) and the down-link of a control to its user needs.
var fields = []string{"source_ids", "target_ids", "allocated", "controls"}

// Links reads every file of objects and answers each entry of source_ids, target_ids, allocated
// and controls whose object has another x than the version it names. It answers too how many entries
// it compared. An entry whose object no file holds, or that names no version, is not compared:
// uplinkcheck and schemacheck report it.
//
// @implements SR_32@1.1.1
// @implements SR_87@1.1.1
func Links(read Read, index *review.Index) ([]Finding, int, error) {
	return LinksWithAbove(read, index, nil)
}

// LinksWithAbove is Links, and compares an entry whose object no file of the repo holds with its
// object in the repo above, where one is given.
//
// @implements SR_109@1.1.1
func LinksWithAbove(read Read, index *review.Index, above *review.Repo) ([]Finding, int, error) {

	versions := newVersions(read, index, above)
	var findings []Finding
	compared := 0
	for _, file := range index.Files() {
		for _, object := range objects(versions.load(file)) {
			for _, field := range fields {
				for _, link := range object.links(field) {
					has, ok := versions.suspect(link)
					if !ok {
						continue
					}
					compared++
					if has != "" {
						findings = append(findings, Finding{object.Key + " " + field, link, has})
					}
				}
			}
		}
	}
	return findings, compared, nil
}

// Annotations answers each id an annotation names whose object has another x than the version it
// names, and how many ids it compared. An id whose object no file holds, or that names no version,
// is not compared: annotationcheck reports it.
//
// @implements SR_33@2.1.1
func Annotations(read Read, index *review.Index, annotations []annotation.Annotation) ([]Finding, int) {
	return AnnotationsWithAbove(read, index, nil, annotations)
}

// AnnotationsWithAbove is Annotations, and compares an id whose object no file of the repo holds
// with its object in the repo above, where one is given.
//
// @implements SR_109@1.1.1
func AnnotationsWithAbove(read Read, index *review.Index, above *review.Repo, annotations []annotation.Annotation) ([]Finding, int) {

	versions := newVersions(read, index, above)
	var findings []Finding
	compared := 0
	for _, a := range annotations {
		for i, id := range a.IDs {
			link := id + "@" + a.Versions[i]
			has, ok := versions.suspect(link)
			if !ok {
				continue
			}
			compared++
			if has != "" {
				findings = append(findings, Finding{fmt.Sprintf("%s:%d @%s", a.File, a.Line, a.Kind), link, has})
			}
		}
	}
	return findings, compared
}

// versions finds the version of an object in its file, and reads each file once. An object its
// repo does not hold is looked for in above, where there is one.
type versions struct {
	read  Read
	index *review.Index
	files map[string][]byte
	above *versions
}

func newVersions(read Read, index *review.Index, above *review.Repo) *versions {
	v := &versions{read, index, map[string][]byte{}, nil}
	if above != nil {
		v.above = newVersions(Read(above.Read), above.Index, nil)
	}
	return v
}

func (v *versions) load(file string) []byte {
	if data, ok := v.files[file]; ok {
		return data
	}
	data, err := v.read(file)
	if err != nil {
		data = nil
	}
	v.files[file] = data
	return data
}

// suspect answers the version the object of link has where its x differs from the x the link
// names, "" where it is the same, and false where the two cannot be compared.
func (v *versions) suspect(link string) (string, bool) {

	id, named, _ := strings.Cut(link, "@")
	if !version.MatchString(named) {
		return "", false
	}
	file, key := v.index.FileOf(id)
	if file == "" || key == "" {
		if v.above != nil {
			return v.above.suspect(link)
		}
		return "", false
	}
	found := review.ObjectIn(v.load(file), key)
	if found == nil {
		return "", false
	}
	var o struct {
		Version string `yaml:"object_version"`
	}
	if yaml.Unmarshal(found, &o) != nil || !version.MatchString(o.Version) {
		return "", false
	}
	if x(o.Version) == x(named) {
		return "", true
	}
	return o.Version, true
}

// version is the form of a version: x.y.z.
var version = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// x answers the first part of a version.
func x(version string) string {
	first, _, _ := strings.Cut(version, ".")
	return first
}

// object is an object of a file of requirements, with the fields that name other objects.
type object struct {
	Key    string         `yaml:"object_key"`
	Fields map[string]any `yaml:",inline"`
}

// links answers the entries of a field that holds a list of links, none where it holds another.
func (o object) links(field string) []string {
	list, _ := o.Fields[field].([]any)
	var out []string
	for _, entry := range list {
		if s, ok := entry.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// objects answers the objects of every section of a file of requirements; none where it cannot be
// read.
func objects(data []byte) []object {

	var doc yaml.Node
	if yaml.Unmarshal(data, &doc) != nil || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil
	}
	//	the sections in the order of the file: keys and values alternate
	sections := doc.Content[0].Content
	var out []object
	for i := 1; i < len(sections); i += 2 {
		var list []object
		if sections[i].Decode(&list) == nil {
			for _, o := range list {
				if o.Key != "" {
					out = append(out, o)
				}
			}
		}
	}
	return out
}
