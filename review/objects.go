// Package review is the component review: what its tools share - the history of git, and how an
// id or a document is found in its file. The check is in review/check, the skeleton in
// review/skeleton, their commands in review/cmd.
package review

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read answers the content of a file of the repo, by its path from the root.
type Read func(file string) ([]byte, error)

// Index tells the file that holds an id: the YAML files that hold objects, each with the ids in
// it (tools/README.md, section The component review).
type Index struct {
	files []string
	of    map[string]string
}

// NewIndex reads the files and keeps, in their order, each that holds an object in one of the
// Lists. A file it cannot read or parse holds none.
//
// @implements SR_86@2.1.1
func NewIndex(read Read, files []string) *Index {

	ix := &Index{of: map[string]string{}}
	for _, file := range files {
		if !strings.HasSuffix(file, ".yaml") && !strings.HasSuffix(file, ".yml") {
			continue
		}
		data, err := read(file)
		if err != nil {
			continue
		}
		var doc map[string]any
		if yaml.Unmarshal(data, &doc) != nil {
			continue
		}
		holds := false
		for _, list := range Lists {
			items, _ := doc[list].([]any)
			for _, item := range items {
				object, _ := item.(map[string]any)
				if key, ok := object["object_key"].(string); ok {
					ix.of[key] = file
					holds = true
				}
			}
		}
		if holds {
			ix.files = append(ix.files, file)
		}
	}
	return ix
}

// IndexAt makes the index of the files git knows in the working tree at root.
func IndexAt(root string) (*Index, error) {

	out, err := ExecGit{Root: root}.run("ls-files")
	if err != nil {
		return nil, err
	}
	return NewIndex(func(file string) ([]byte, error) {
		return os.ReadFile(filepath.Join(root, file))
	}, strings.Fields(out)), nil
}

// Repo is a repo as a check reads it: its files, and the index of the ids they hold.
type Repo struct {
	Read  Read
	Index *Index
}

// RepoAt answers the repo at root, as git knows its files in the working tree.
func RepoAt(root string) (*Repo, error) {

	index, err := IndexAt(root)
	if err != nil {
		return nil, err
	}
	return &Repo{func(file string) ([]byte, error) {
		return os.ReadFile(filepath.Join(root, file))
	}, index}, nil
}

// Files answers every file that holds objects, each once, in the order git lists them.
func (ix *Index) Files() []string {
	return ix.files
}

// Lists are the sections of a YAML that hold objects with an object_key.
var Lists = []string{"requirements", "elements", "decisions", "hazards", "controls", "rules", "testcases"}

// FileOf answers the file that holds a target, and the key to find it there; a document is its
// own file and has no key. An id no file holds answers none.
func (ix *Index) FileOf(target string) (file, key string) {

	if strings.HasSuffix(target, ".md") {
		return target, ""
	}
	if file, ok := ix.of[target]; ok {
		return file, target
	}
	return "", ""
}

// ObjectIn answers the object with key in the YAML data, as YAML, or the whole data where key is
// empty; nil where it is not there.
func ObjectIn(data []byte, key string) []byte {

	if key == "" {
		return data
	}
	var doc map[string]any
	if yaml.Unmarshal(data, &doc) != nil {
		return nil
	}
	for _, list := range Lists {
		items, _ := doc[list].([]any)
		for _, item := range items {
			object, _ := item.(map[string]any)
			if object["object_key"] == key {
				out, _ := yaml.Marshal(object)
				return out
			}
		}
	}
	return nil
}
