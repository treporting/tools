// Package check checks every YAML file against the schema it names in its field schema (TOOL_05).
package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

// Finding is one place where a file does not match its schema, or a file that cannot be checked.
type Finding struct {
	File    string
	Path    string // the place within the file, empty for the file as a whole
	Message string
}

func (f Finding) String() string {
	if f.Path == "" {
		return f.File + ": " + f.Message
	}
	return f.File + " " + f.Path + ": " + f.Message
}

// notOurs are the folders whose YAML files have a form another owns, and are not our data.
var notOurs = []string{".github/"}

// Files answers the YAML files of our own data under version control below root, relative to it.
//
// @implements SR_18@1.1.1
func Files(root string) ([]string, error) {
	out, err := exec.Command("git", "-C", root, "ls-files", "--", "*.yaml", "*.yml").Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	var files []string
next:
	for _, f := range strings.Fields(string(out)) {
		for _, folder := range notOurs {
			if strings.HasPrefix(f, folder) {
				continue next
			}
		}
		files = append(files, f)
	}
	return files, nil
}

// Check checks each of the files, relative to root, against the schema its field schema names, and
// answers every place where one does not match, and every file that names no schema or a schema
// that is not there.
//
// @implements SR_18@1.1.1, SR_19@1.1.1
func Check(root string, files []string) ([]Finding, error) {
	return CheckWithAbove(root, "", files)
}

// CheckWithAbove is Check, and takes a schema that root does not hold from above, the root of
// the repo above, where it is not empty.
//
// @implements SR_107@1.1.1
func CheckWithAbove(root, above string, files []string) ([]Finding, error) {

	schemas := map[string]*jsonschema.Schema{}
	var findings []Finding
	for _, file := range files {

		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			return nil, err
		}
		var doc any
		if err := yaml.Unmarshal(data, &doc); err != nil {
			findings = append(findings, Finding{file, "", "is no YAML: " + err.Error()})
			continue
		}

		object, _ := doc.(map[string]any)
		name, _ := object["schema"].(string)
		if name == "" {
			findings = append(findings, Finding{file, "", "names no schema"})
			continue
		}
		schema, ok := schemas[name]
		if !ok {
			schema, err = compile(root, name)
			if errors.Is(err, fs.ErrNotExist) && above != "" {
				schema, err = compile(above, name)
			}
			if errors.Is(err, fs.ErrNotExist) {
				findings = append(findings, Finding{file, "", "names the schema " + name + ", which is not there"})
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			schemas[name] = schema
		}

		//	the schema checks JSON: the YAML is read as JSON, so that a number keeps its kind
		raw, err := json.Marshal(doc)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		var invalid *jsonschema.ValidationError
		if err := schema.Validate(instance); errors.As(err, &invalid) {
			for _, unit := range invalid.BasicOutput().Errors {
				if unit.Error != nil {
					findings = append(findings, Finding{file, "/" + strings.TrimPrefix(unit.InstanceLocation, "/"), unit.Error.String()})
				}
			}
		} else if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
	}
	return findings, nil
}

// compile reads the schema name, relative to root, and compiles it.
func compile(root, name string) (*jsonschema.Schema, error) {

	path, err := filepath.Abs(filepath.Join(root, name))
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	doc, err := jsonschema.UnmarshalJSON(f)
	if err != nil {
		return nil, err
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(path, doc); err != nil {
		return nil, err
	}
	return c.Compile(path)
}
