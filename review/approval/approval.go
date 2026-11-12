// Package approval checks every approval from the approved objects (TOOL_32, TOOL_33): an object
// that needs a review has a passed review of the content it has now, and every approved object has
// the content it was approved with.
package approval

import (
	"crypto/sha1" // #nosec G505 - the blob hash of git, no security
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/treporting/specs-and-deploy/tools/review"
	"gopkg.in/yaml.v3"
)

// Git is what the check asks of the history of the repo.
type Git interface {
	// Show answers the content of path in the commit rev.
	Show(rev, path string) ([]byte, error)
	// FileCommits answers the commits up to rev that changed path, the newest first.
	FileCommits(rev, path string) ([]string, error)
	// Changed answers the files changed between base and HEAD.
	Changed(base string) ([]string, error)
}

// Finding is one approval that does not hold.
type Finding struct {
	File    string
	Object  string
	Message string
}

func (f Finding) String() string {
	return f.File + ": " + f.Object + ": " + f.Message
}

// documents is the file that lists every document, with its version and status.
const documents = "documents.yaml"

// object is one approved requirement or document.
type object struct {
	key     string // the object_key, or the path of a document
	file    string // the file that holds it, or the document itself
	x       string // the first part of its version
	content any    // what x counts: the attributes for x, or the blob hash of a document
	review  bool   // whether it needs a review
	hash    string // object_hash of a document
	entry   entry  // the entry of a document in documents.yaml
}

// Check checks the approved objects of the repo at root. Where base is given, it checks only the
// objects in the files changed against it; else every approved object.
//
// @implements SR_75@1.1.1, SR_76@1.1.1, SR_77@1.1.1, SR_78@1.1.1, SR_79@2.1.2, SR_80@2.1.3, SR_85@1.1.1
func Check(root, base string, git Git, index *review.Index) ([]Finding, error) {

	read := func(file string) ([]byte, error) { return os.ReadFile(filepath.Join(root, file)) }

	requirements, err := approvedRequirements(read, index.Files())
	if err != nil {
		return nil, err
	}
	docs, err := approvedDocuments(read)
	if err != nil {
		return nil, err
	}

	if base != "" {
		if requirements, docs, err = scoped(requirements, docs, base, git); err != nil {
			return nil, err
		}
	}

	last, err := lastReviews(root)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	for _, o := range append(requirements, docs...) {
		if o.review {
			findings = append(findings, checkReview(o, last[o.key], git)...)
		}
	}
	versions, err := checkVersions(requirements, docs, git)
	return append(findings, versions...), err
}

// scoped keeps the requirements and documents in the files changed against base.
//
// @implements SR_80@2.1.3
func scoped(requirements, docs []object, base string, git Git) ([]object, []object, error) {
	changed, err := git.Changed(base)
	if err != nil {
		return nil, nil, err
	}
	docs, err = changedDocuments(docs, changed, base, git)
	return inFiles(requirements, changed), docs, err
}

// checkVersions answers every approved object whose content for x is not the one it was approved
// with: a requirement against its first approval with its x, a document against its object_hash.
//
// @implements SR_78@1.1.1, SR_79@2.1.2
func checkVersions(requirements, docs []object, git Git) ([]Finding, error) {

	firsts, err := firstApprovals(requirements, git)
	if err != nil {
		return nil, err
	}
	var findings []Finding
	for _, o := range requirements {
		if first, ok := firsts[o.file+"\x00"+o.key]; ok && !reflect.DeepEqual(first, o.content) {
			findings = append(findings, Finding{o.file, o.key, "its attributes for x are not those of its approval with x " + o.x})
		}
	}
	for _, o := range docs {
		if o.hash == "" || !strings.HasPrefix(o.content.(string), o.hash) {
			findings = append(findings, Finding{documents, o.key, "its object_hash " + o.hash + " is not the hash of its file"})
		}
	}
	return findings, nil
}

// checkReview answers what is wrong with the last review of an object that needs one.
//
// @implements SR_75@1.1.1, SR_76@1.1.1, SR_77@1.1.1, SR_85@1.1.1
func checkReview(o object, r *result, git Git) []Finding {

	if r == nil {
		return []Finding{{o.file, o.key, "approved, and no review result names it"}}
	}
	if r.verdict != "passed" {
		return []Finding{{o.file, o.key, "the last review result that names it, " + r.path + ", has the verdict " + r.verdict}}
	}
	if r.status == "open" {
		return []Finding{{o.file, o.key, "the last review result that names it, " + r.path + ", is open"}}
	}
	data, err := git.Show(r.revision, o.file)
	if err != nil {
		return []Finding{{o.file, o.key, "the revision " + r.revision + " of " + r.path + " does not hold it"}}
	}
	var reviewed any
	if o.file == o.key {
		reviewed = blobHash(data)
	} else {
		reviewed = xOf(review.ObjectIn(data, o.key), o.attributes())
	}
	if !reflect.DeepEqual(reviewed, o.content) {
		return []Finding{{o.file, o.key, "its content for x is not the one " + r.path + " reviewed at " + r.revision}}
	}
	return nil
}

// attributes answers the attributes for x a requirement was read with.
func (o object) attributes() []string {
	m, _ := o.content.(map[string]any)
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// result is the review result that names an object last.
type result struct {
	path     string
	at       time.Time
	revision string
	verdict  string
	status   string
}

// lastReviews answers, for every target, the review result that names it with the latest
// created_at of its file.
//
// @implements SR_76@1.1.1
func lastReviews(root string) (map[string]*result, error) {

	files, err := filepath.Glob(filepath.Join(root, "reviews", "results", "*.yaml"))
	if err != nil {
		return nil, err
	}
	last := map[string]*result{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var r struct {
			Status    string `yaml:"status"`
			CreatedAt string `yaml:"created_at"`
			Results   []struct {
				TargetIDs []string `yaml:"target_ids"`
				Revision  string   `yaml:"revision"`
				Verdict   string   `yaml:"verdict"`
			} `yaml:"results"`
		}
		if err := yaml.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		at, err := time.Parse(time.RFC3339, r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: created_at: %w", f, err)
		}
		rel, _ := filepath.Rel(root, f)
		for _, res := range r.Results {
			for _, target := range res.TargetIDs {
				if old, ok := last[target]; ok && !at.After(old.at) {
					continue
				}
				last[target] = &result{path: rel, at: at, revision: res.Revision, verdict: res.Verdict, status: r.Status}
			}
		}
	}
	return last, nil
}

// approvedRequirements answers the approved objects of the files that hold requirements, in the
// order of the files and of the objects in them.
func approvedRequirements(read review.Read, files []string) ([]object, error) {

	var objects []object
	for _, file := range files {
		data, err := read(file)
		if err != nil {
			return nil, err
		}
		var head struct {
			Schema string `yaml:"schema"`
			Level  string `yaml:"level"`
		}
		if err := yaml.Unmarshal(data, &head); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		x, err := attributesForX(read, head.Schema)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		_, noCriteria := read(criteriaFile(head.Level))
		for _, item := range itemsOf(data) {
			if item.status != "approved" {
				continue
			}
			objects = append(objects, object{
				key:     item.key,
				file:    file,
				x:       item.x,
				content: xOf(item.yaml, x[item.list]),
				review:  noCriteria == nil,
			})
		}
	}
	return objects, nil
}

// approvedDocuments answers the approved and completed documents of documents.yaml, in its order,
// with the blob hash of their file; none where the repo has no documents.yaml.
func approvedDocuments(read review.Read) ([]object, error) {

	data, err := read(documents)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	entries, err := documentEntries(data)
	if err != nil {
		return nil, err
	}
	_, noCriteria := read(criteriaFile("document"))
	var objects []object
	for _, e := range entries {
		if !checked(e.Status) {
			continue
		}
		content, err := read(e.Name)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object{
			key:     e.Name,
			file:    e.Name,
			x:       strings.SplitN(e.ObjectVersion, ".", 2)[0],
			content: blobHash(content),
			review:  noCriteria == nil && slices.Contains(e.VerificationMethod, "review"),
			hash:    e.ObjectHash,
			entry:   e,
		})
	}
	return objects, nil
}

// entry is a document as documents.yaml lists it.
type entry struct {
	Name               string   `yaml:"name"`
	ObjectVersion      string   `yaml:"object_version"`
	ObjectHash         string   `yaml:"object_hash"`
	Status             string   `yaml:"status"`
	VerificationMethod []string `yaml:"verification_method"`
}

// checked answers whether a document of this status is checked: an approved one, and a completed
// one, a record a tool wrote (PROC_13).
//
// @implements SR_79@2.1.2
func checked(status string) bool {
	return status == "approved" || status == "completed"
}

func documentEntries(data []byte) ([]entry, error) {
	var d struct {
		Documents []entry `yaml:"documents"`
	}
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("%s: %w", documents, err)
	}
	return d.Documents, nil
}

// criteriaFile answers the criteria file of a level, as reviewskeleton names it.
func criteriaFile(level string) string {
	if level == "document" {
		return "reviews/criteria-documents.yaml"
	}
	return "reviews/criteria-" + level + ".yaml"
}

// inFiles keeps the objects held by one of the files.
//
// @implements SR_80@2.1.3
func inFiles(objects []object, files []string) []object {
	var kept []object
	for _, o := range objects {
		if slices.Contains(files, o.file) {
			kept = append(kept, o)
		}
	}
	return kept
}

// changedDocuments keeps the documents whose file changed against base, or whose entry in
// documents.yaml did.
//
// @implements SR_80@2.1.3
func changedDocuments(docs []object, changed []string, base string, git Git) ([]object, error) {

	before := map[string]entry{}
	if slices.Contains(changed, documents) {
		data, err := git.Show(base, documents)
		if err == nil {
			entries, err := documentEntries(data)
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				before[e.Name] = e
			}
		}
	}
	var kept []object
	for _, o := range docs {
		b, ok := before[o.key]
		entryChanged := slices.Contains(changed, documents) && (!ok || !reflect.DeepEqual(b, o.entry))
		if slices.Contains(changed, o.file) || entryChanged {
			kept = append(kept, o)
		}
	}
	return kept, nil
}

// firstApprovals answers, for each requirement, the attributes for x of the first commit in which
// it was approved with its x, keyed by file and key. A requirement no commit holds so has none.
//
// @implements SR_78@1.1.1
func firstApprovals(objects []object, git Git) (map[string]any, error) {

	open := map[string]map[string]object{} // file -> key -> object not found yet
	for _, o := range objects {
		if open[o.file] == nil {
			open[o.file] = map[string]object{}
		}
		open[o.file][o.key] = o
	}
	firsts := map[string]any{}
	for file, wanted := range open {
		commits, err := git.FileCommits("HEAD", file)
		if err != nil {
			return nil, err
		}
		slices.Reverse(commits)
		for _, commit := range commits {
			if len(wanted) == 0 {
				break
			}
			data, err := git.Show(commit, file)
			if err != nil {
				continue
			}
			for _, item := range itemsOf(data) {
				o, ok := wanted[item.key]
				if !ok || item.status != "approved" || item.x != o.x {
					continue
				}
				firsts[file+"\x00"+item.key] = xOf(item.yaml, o.attributes())
				delete(wanted, item.key)
			}
		}
	}
	return firsts, nil
}

// item is an object of a YAML file that holds objects.
type item struct {
	list   string
	key    string
	x      string
	status string
	yaml   map[string]any
}

// itemsOf answers the objects of a YAML file, in the order of its sections and of the objects in
// them; a file it cannot parse holds none.
func itemsOf(data []byte) []item {
	var doc map[string]any
	if yaml.Unmarshal(data, &doc) != nil {
		return nil
	}
	var items []item
	for _, list := range review.Lists {
		entries, _ := doc[list].([]any)
		for _, e := range entries {
			m, _ := e.(map[string]any)
			key, _ := m["object_key"].(string)
			if key == "" {
				continue
			}
			version, _ := m["object_version"].(string)
			status, _ := m["status"].(string)
			items = append(items, item{list, key, strings.SplitN(version, ".", 2)[0], status, m})
		}
	}
	return items
}

// xOf answers the attributes for x of an object, as YAML or as its map; an attribute it lacks is
// nil.
func xOf(obj any, attributes []string) map[string]any {
	m, ok := obj.(map[string]any)
	if !ok {
		data, _ := obj.([]byte)
		_ = yaml.Unmarshal(data, &m)
	}
	content := map[string]any{}
	for _, a := range attributes {
		content[a] = m[a]
	}
	return content
}

// attributesForX answers, for each section of the files of a schema, the attributes its items
// name for x (traceability.md, section The version in the schema).
func attributesForX(read review.Read, schema string) (map[string][]string, error) {

	data, err := read(schema)
	if err != nil {
		return nil, err
	}
	var s struct {
		Properties map[string]struct {
			Items map[string]any `json:"items" yaml:"items"`
		} `yaml:"properties"`
	}
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("%s: %w", schema, err)
	}
	x := map[string][]string{}
	for list, p := range s.Properties {
		x[list] = findX(p.Items)
	}
	return x, nil
}

// findX answers the first attributes_for_version_x in a schema of items, also inside oneOf.
func findX(schema map[string]any) []string {
	if names, ok := schema["attributes_for_version_x"].([]any); ok {
		var x []string
		for _, n := range names {
			if s, ok := n.(string); ok {
				x = append(x, s)
			}
		}
		return x
	}
	for _, key := range []string{"oneOf", "anyOf"} {
		variants, _ := schema[key].([]any)
		for _, v := range variants {
			if m, ok := v.(map[string]any); ok {
				if x := findX(m); x != nil {
					return x
				}
			}
		}
	}
	return nil
}

// blobHash answers the hash git gives the content of a file.
func blobHash(data []byte) string {
	sum := sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(data))), data...)) // #nosec G401 - the blob hash of git, no security
	return hex.EncodeToString(sum[:])
}
