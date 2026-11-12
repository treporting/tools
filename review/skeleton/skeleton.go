// Package skeleton writes the skeleton of a review result (TOOL_04): what the reviewer would
// otherwise copy by hand. Every judgement is todo, so an unfilled skeleton fails reviewcheck.
package skeleton

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/treporting/specs-and-deploy/tools/review"
	"gopkg.in/yaml.v3"
)

// resultSchema is the schema of a review result, from the root of the repo, and resultType its
// type (formats.md, section The head of a file).
//
// @implements SR_17@2.1.1
const (
	resultSchema = "schemas/review-result.schema.json"
	resultType   = "review"
)

// recordVersion and recordOpen make the skeleton a record that is still being made (formats.md,
// section Reviews).
//
// @implements SR_83@1.1.1
const (
	recordVersion = "1.1.1"
	recordOpen    = "open"
)

// results is the directory of the review results, apart from the criteria in reviews/.
//
// @implements SR_11@2.1.1
const results = "reviews/results"

// Git is what the skeleton asks of the history of the repo.
type Git interface {
	// Show answers the content of path in the commit rev.
	Show(rev, path string) ([]byte, error)
	// FileCommits answers the commits up to rev that changed path, the newest first.
	FileCommits(rev, path string) ([]string, error)
	// Short answers the short hash of rev.
	Short(rev string) (string, error)
}

// Request is what a skeleton is made for.
type Request struct {
	Reviewer string
	Scope    string
	Targets  []string // ids of requirements, or paths of documents
	Now      time.Time
}

// Make writes the skeleton of a review result for req into reviews/results/ below root, made on the commit
// HEAD, and answers its path relative to root. Where a target is not found, or the targets belong
// to more than one level, it writes nothing.
//
// @implements SR_11@2.1.1, SR_12@1.1.1, SR_13@1.1.1, SR_14@1.1.1, SR_15@1.1.1, SR_17@2.1.1, SR_83@1.1.1
func Make(root string, req Request, git Git, index *review.Index) (string, error) {

	head, err := git.Short("HEAD")
	if err != nil {
		return "", err
	}

	level, err := levelOf(git, index, req.Targets)
	if err != nil {
		return "", err
	}

	criteriaFile, criteriaRevision, criteriaLevel, keys, err := criteriaOf(git, level)
	if err != nil {
		return "", err
	}

	path := filepath.Join(results, req.Now.Format("20060102T150405")+"_"+req.Scope+".yaml")
	if _, err := os.Stat(filepath.Join(root, path)); err == nil {
		return "", fmt.Errorf("%s is there already", path)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Review of %s against the criteria of the level %s.\n", strings.Join(req.Targets, ", "), level)
	fmt.Fprintf(&b, "# Skeleton made by reviewskeleton: every todo is for the reviewer to decide.\n")
	fmt.Fprintf(&b, "# Format: ../../formats.md, section Reviews.\n")
	fmt.Fprintf(&b, "schema: %s\n", resultSchema)
	fmt.Fprintf(&b, "type: %s\n", resultType)
	fmt.Fprintf(&b, "level: %s\n", criteriaLevel)
	fmt.Fprintf(&b, "object_version: '%s'\n", recordVersion)
	fmt.Fprintf(&b, "status: %s\n", recordOpen)
	fmt.Fprintf(&b, "created_by: %s\n", req.Reviewer)
	fmt.Fprintf(&b, "created_at: '%s'\n", req.Now.Format("2006-01-02T15:04:05-07:00"))
	fmt.Fprintf(&b, "criteria: %s\n", criteriaFile)
	fmt.Fprintf(&b, "criteria_revision: '%s'\n", criteriaRevision)
	fmt.Fprintf(&b, "results:\n")
	for i, target := range req.Targets {
		fmt.Fprintf(&b, "- object_key: RR_%02d\n", i+1)
		fmt.Fprintf(&b, "  target_ids: [%s]\n", target)
		fmt.Fprintf(&b, "  revision: '%s'\n", head)
		fmt.Fprintf(&b, "  criteria:\n")
		for _, key := range keys {
			fmt.Fprintf(&b, "    %s: todo\n", key)
		}
		fmt.Fprintf(&b, "  verdict: todo\n")
	}

	if err := os.MkdirAll(filepath.Join(root, results), 0o750); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(root, path), []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return path, nil
}

// levelOf answers the one level of the targets in HEAD: the field level of the YAML that holds a
// requirement, and document for a document.
//
// @implements SR_15@1.1.1
func levelOf(git Git, index *review.Index, targets []string) (string, error) {

	if len(targets) == 0 {
		return "", fmt.Errorf("no id or document given")
	}
	level := ""
	for _, target := range targets {
		file, key := index.FileOf(target)
		if file == "" {
			return "", fmt.Errorf("no file holds %s", target)
		}
		data, err := git.Show("HEAD", file)
		if err != nil || review.ObjectIn(data, key) == nil {
			return "", fmt.Errorf("%s is not in %s", target, file)
		}
		l := "document"
		if key != "" {
			var doc struct {
				Level string `yaml:"level"`
			}
			if err := yaml.Unmarshal(data, &doc); err != nil {
				return "", fmt.Errorf("%s: %w", file, err)
			}
			l = doc.Level
		}
		if level != "" && l != level {
			return "", fmt.Errorf("%s is of the level %s, the others of %s - one review, one level", target, l, level)
		}
		level = l
	}
	return level, nil
}

// criteriaOf answers the criteria file of a level, the last commit that changed it, and the level
// and the ids of its criteria in that commit.
//
// @implements SR_12@1.1.1, SR_17@2.1.1
func criteriaOf(git Git, level string) (file, revision, criteriaLevel string, keys []string, err error) {

	name := "criteria-" + level + ".yaml"
	if level == "document" {
		name = "criteria-documents.yaml"
	}
	file = "reviews/" + name

	commits, err := git.FileCommits("HEAD", file)
	if err != nil || len(commits) == 0 {
		return "", "", "", nil, fmt.Errorf("no commit holds %s", file)
	}
	if revision, err = git.Short(commits[0]); err != nil {
		return "", "", "", nil, err
	}
	data, err := git.Show(revision, file)
	if err != nil {
		return "", "", "", nil, err
	}
	var c struct {
		Level    string `yaml:"level"`
		Criteria []struct {
			ObjectKey string `yaml:"object_key"`
		} `yaml:"criteria"`
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return "", "", "", nil, fmt.Errorf("%s: %w", file, err)
	}
	for _, k := range c.Criteria {
		keys = append(keys, k.ObjectKey)
	}
	return file, revision, c.Level, keys, nil
}
