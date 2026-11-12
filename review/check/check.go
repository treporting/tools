// Package check checks the results of reviews against the rules of formats.md (TOOL_03, TOOL_11).
package check

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/treporting/specs-and-deploy/tools/review"
	"gopkg.in/yaml.v3"
)

// Review is one file of review results.
type Review struct {
	Level            string   `yaml:"level"`
	CreatedBy        string   `yaml:"created_by"`
	Status           string   `yaml:"status"`
	CreatedAt        string   `yaml:"created_at"`
	CompletedAt      string   `yaml:"completed_at"`
	Criteria         string   `yaml:"criteria"`
	CriteriaRevision string   `yaml:"criteria_revision"`
	Results          []Result `yaml:"results"`
}

// Result is the review of one requirement or document.
type Result struct {
	ObjectKey string               `yaml:"object_key"`
	TargetIDs []string             `yaml:"target_ids"`
	Revision  string               `yaml:"revision"`
	Criteria  map[string]Judgement `yaml:"criteria"`
	Verdict   string               `yaml:"verdict"`
}

// Judgement is the result of one criterion: `passed`, or a block of result and note.
type Judgement struct {
	Result string `yaml:"result"`
	Note   string `yaml:"note"`
}

// UnmarshalYAML reads a judgement written as a word or as a block.
func (j *Judgement) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		j.Result = node.Value
		return nil
	}
	type plain Judgement
	return node.Decode((*plain)(j))
}

// Finding is one rule a review result breaks.
type Finding struct {
	File    string
	Result  string
	Message string
}

func (f Finding) String() string {
	return fmt.Sprintf("%s %s: %s", f.File, f.Result, f.Message)
}

// Git is what the check asks of the history of the repo.
type Git interface {
	// CommitExists says whether rev names a commit.
	CommitExists(rev string) bool
	// Show answers the content of path in the commit rev.
	Show(rev, path string) ([]byte, error)
	// FileCommits answers the commits up to rev that changed path, the newest first.
	FileCommits(rev, path string) ([]string, error)
	// Parent answers the first parent of a commit, and false for the first commit.
	Parent(commit string) (string, bool)
	// IsMerge says whether a commit has more than one parent.
	IsMerge(commit string) bool
	// Authors answers the author of a commit and the names of its co-authors.
	Authors(commit string) (author string, coAuthors []string, err error)
}

// Check reads the review result in path, below root, and answers every rule it breaks. The index
// tells the file of each target.
//
// @implements SR_01@1.1.1, SR_02@1.1.1, SR_03@1.1.1, SR_04@1.1.1, SR_05@2.1.1, SR_06@1.1.1, SR_07@1.1.1, SR_08@1.1.1, SR_16@1.1.1, SR_84@1.1.1
func Check(root, path string, git Git, index *review.Index) ([]Finding, error) {

	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		return nil, err
	}
	var r Review
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	criteria, err := criteriaKeys(git, r.CriteriaRevision, r.Criteria)
	if err != nil {
		return []Finding{{path, "", err.Error()}}, nil
	}

	var findings []Finding
	for _, message := range checkCompleted(r) {
		findings = append(findings, Finding{path, "", message})
	}
	for _, result := range r.Results {
		for _, message := range checkResult(r, result, criteria, git, index) {
			findings = append(findings, Finding{path, result.ObjectKey, message})
		}
	}
	return findings, nil
}

// checkCompleted answers where a completed review has no completed_at, or one before its
// created_at.
//
// @implements SR_84@1.1.1
func checkCompleted(r Review) []string {

	if r.Status != "completed" {
		return nil
	}
	if r.CompletedAt == "" {
		return []string{"completed, and no completed_at"}
	}
	completed, err := time.Parse(time.RFC3339, r.CompletedAt)
	if err != nil {
		return []string{"the completed_at " + r.CompletedAt + " is no time"}
	}
	created, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return []string{"the created_at " + r.CreatedAt + " is no time"}
	}
	if completed.Before(created) {
		return []string{"the completed_at " + r.CompletedAt + " is before the created_at " + r.CreatedAt}
	}
	return nil
}

// criteriaKeys answers the ids of the criteria in the criteria file, in the revision the review used.
// The revision is a commit that changed the file, so that the criteria the review followed can be
// traced.
//
// @implements SR_07@1.1.1
func criteriaKeys(git Git, revision, file string) ([]string, error) {

	if !git.CommitExists(revision) {
		return nil, fmt.Errorf("the criteria revision %s is no commit", revision)
	}
	commits, err := git.FileCommits(revision, file)
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 || !strings.HasPrefix(commits[0], revision) {
		return nil, fmt.Errorf("the criteria revision %s did not change %s", revision, file)
	}
	data, err := git.Show(revision, file)
	if err != nil {
		return nil, fmt.Errorf("the criteria %s are not in %s", file, revision)
	}
	var c struct {
		Criteria []struct {
			ObjectKey string `yaml:"object_key"`
		} `yaml:"criteria"`
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("%s in %s: %w", file, revision, err)
	}
	keys := make([]string, 0, len(c.Criteria))
	for _, k := range c.Criteria {
		keys = append(keys, k.ObjectKey)
	}
	return keys, nil
}

// checkResult answers the rules one result breaks.
//
// @implements SR_01@1.1.1, SR_02@1.1.1, SR_03@1.1.1, SR_04@1.1.1, SR_05@2.1.1, SR_06@1.1.1, SR_08@1.1.1, SR_16@1.1.1
func checkResult(r Review, result Result, criteria []string, git Git, index *review.Index) []string {

	messages, allPassed := checkCriteria(r, result, criteria)
	messages = append(messages, checkVerdict(result.Verdict, allPassed)...)

	if !git.CommitExists(result.Revision) {
		return append(messages, "the revision "+result.Revision+" is no commit")
	}
	return append(messages, checkAuthors(r, result, git, index)...)
}

// checkCriteria answers the rules the criteria of one result break, and whether every criterion
// passed.
func checkCriteria(r Review, result Result, criteria []string) ([]string, bool) {

	var messages []string

	//	every criterion once - a map has every key once, so what is left to
	//	ask is whether one is missing or one is not a criterion of the file;
	//	walked by criterion in the order of the file, those not in it by id
	known := map[string]bool{}
	for _, key := range criteria {
		known[key] = true
	}
	var foreign []string
	for key := range result.Criteria {
		if !known[key] {
			foreign = append(foreign, key)
		}
	}
	sort.Strings(foreign)

	allPassed := true
	for _, key := range append(append([]string(nil), criteria...), foreign...) {
		j, ok := result.Criteria[key]
		if !ok {
			messages = append(messages, "criterion "+key+" is missing")
			continue
		}
		if !known[key] {
			messages = append(messages, "criterion "+key+" is not in "+r.Criteria)
		}
		message, passed := checkJudgement(key, j)
		if message != "" {
			messages = append(messages, message)
		}
		allPassed = allPassed && passed
	}
	return messages, allPassed
}

// checkJudgement answers the rule one judgement breaks, or "", and whether it passed.
func checkJudgement(key string, j Judgement) (string, bool) {
	switch j.Result {
	case "passed":
		return "", true
	case "failed":
		if strings.TrimSpace(j.Note) == "" {
			return "criterion " + key + " failed with no note", false
		}
		return "", false
	default:
		return "criterion " + key + " is neither passed nor failed", false
	}
}

// checkVerdict answers the rule the verdict breaks, where it does not follow from the criteria.
func checkVerdict(verdict string, allPassed bool) []string {
	switch {
	case verdict == "passed" && !allPassed:
		return []string{"the verdict is passed, but not every criterion passed"}
	case verdict == "failed" && allPassed:
		return []string{"the verdict is failed, but every criterion passed"}
	case verdict != "passed" && verdict != "failed":
		return []string{"the verdict is neither passed nor failed"}
	}
	return nil
}

// checkAuthors answers a rule for each target the reviewer is an author of.
func checkAuthors(r Review, result Result, git Git, index *review.Index) []string {
	var messages []string
	for _, target := range result.TargetIDs {
		authors, err := authorsOf(git, index, result.Revision, target)
		if err != nil {
			messages = append(messages, err.Error())
			continue
		}
		for _, name := range authors {
			if name == r.CreatedBy {
				messages = append(messages, "the reviewer "+r.CreatedBy+" is an author of "+target)
				break
			}
		}
	}
	return messages
}

// authorsOf answers the authors of a requirement or document in a revision: the author and the
// co-authors of the last commit, not a merge, that changed it, whatever names stand there. A merge
// joins what others wrote; the one who merged approves, and is no author.
//
// @implements SR_05@2.1.1, SR_08@1.1.1
func authorsOf(git Git, index *review.Index, revision, target string) ([]string, error) {

	file, key := index.FileOf(target)
	if file == "" {
		return nil, fmt.Errorf("no file holds %s", target)
	}

	commits, err := git.FileCommits(revision, file)
	if err != nil {
		return nil, err
	}
	for _, commit := range commits {
		if git.IsMerge(commit) {
			continue
		}
		now := objectAt(git, commit, file, key)
		if now == nil {
			continue
		}
		before := []byte(nil)
		if parent, ok := git.Parent(commit); ok {
			before = objectAt(git, parent, file, key)
		}
		if bytes.Equal(now, before) {
			continue
		}
		author, coAuthors, err := git.Authors(commit)
		if err != nil {
			return nil, err
		}
		return append([]string{author}, coAuthors...), nil
	}
	return nil, fmt.Errorf("%s is not in %s in %s", target, file, revision)
}

// objectAt answers the object with key in file in a commit, as YAML, or the whole file where key
// is empty; nil where it is not there.
func objectAt(git Git, commit, file, key string) []byte {

	data, err := git.Show(commit, file)
	if err != nil {
		return nil
	}
	return review.ObjectIn(data, key)
}

// Results answers the review results to check, below root: with a base commit, those in
// reviews/results/ that differ between it and HEAD and are still there; without one, every one.
//
// @implements SR_103@1.1.1
func Results(root, base string, changed func(base string) ([]string, error)) ([]string, error) {
	if base == "" {
		found, err := filepath.Glob(filepath.Join(root, "reviews", "results", "*.yaml"))
		if err != nil {
			return nil, err
		}
		files := make([]string, 0, len(found))
		for _, f := range found {
			rel, err := filepath.Rel(root, f)
			if err != nil {
				return nil, err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return files, nil
	}
	paths, err := changed(base)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, p := range paths {
		if ok, _ := path.Match("reviews/results/*.yaml", p); !ok {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, p)); err == nil {
			files = append(files, p)
		}
	}
	return files, nil
}
