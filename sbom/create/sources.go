package create

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ReadLicence answers the text of the licence file in dir, and nothing where there is none.
type ReadLicence func(dir string) string

// licenceFiles are the names a licence file has, in the order they are looked for.
var licenceFiles = []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE", "COPYING"}

// LicenceIn reads the licence file in dir.
func LicenceIn(dir string) string {
	for _, name := range licenceFiles {
		if data, err := os.ReadFile(filepath.Join(dir, name)); err == nil {
			return string(data)
		}
	}
	return ""
}

// moduleFormat makes go list print each module a package belongs to - its replacement where it has
// one - as path, version and directory, and nothing for the module of the subsystem.
const moduleFormat = `{{with .Module}}{{if not .Main}}{{with .Replace}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{else}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{end}}{{end}}{{end}}`

// commandsFormat makes go list print the import path of each command, and nothing for the others.
const commandsFormat = `{{if eq .Name "main"}}{{.ImportPath}}{{end}}`

// commandsCommand is the go list that names the commands, as a shell runs it.
const commandsCommand = "go list -f '" + commandsFormat + "' ./..."

// DirsWith answers every directory of root that holds a file of that name, root included, in lexical
// order - neither below node_modules nor below a directory whose name starts with a dot, as the build
// steps find a go.mod or a package-lock.json.
//
// @implements SR_36@2.1.1
// @implements SR_37@2.1.1
func DirsWith(root, name string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root && (d.Name() == "node_modules" || strings.HasPrefix(d.Name(), ".")) {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == name {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return dirs, err
}

// GoList runs go list in root, the root of the module of the subsystem, and answers a line per
// package of its commands: the module it belongs to, as moduleFormat prints it. It answers nothing
// where the module has no command.
//
// @implements SR_36@2.1.1
func GoList(root string) (string, error) {
	commands, err := run(root, "go", "list", "-f", commandsFormat, "./...")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(commands) == "" {
		return "", nil
	}
	args := append([]string{"list", "-deps", "-f", moduleFormat}, strings.Fields(commands)...)
	return run(root, "go", args...)
}

// GoModules answers the Go modules of the output of GoList, each once, with the licence of the
// licence file in its directory of the module cache.
//
// @implements SR_36@2.1.1
func GoModules(list string, read ReadLicence) []Package {
	seen := map[string]bool{}
	var packages []Package
	for _, line := range strings.Split(list, "\n") {
		fields := strings.Split(strings.TrimSpace(line), "\t")
		if len(fields) != 3 || seen[fields[0]+"@"+fields[1]] {
			continue
		}
		seen[fields[0]+"@"+fields[1]] = true
		packages = append(packages, Package{
			PURL:    "pkg:golang/" + fields[0] + "@" + fields[1],
			Name:    fields[0],
			Version: fields[1],
			Licence: Licence(read(fields[2])),
		})
	}
	return packages
}

// Stdlib answers the standard library of the Go the commands are built with: version as go env
// GOVERSION names it, licence from the licence file of GOROOT.
//
// @implements SR_44@1.1.1
func Stdlib(goVersion, licenceText string) Package {
	return Package{
		PURL:    "pkg:golang/stdlib@" + goVersion,
		Name:    "stdlib",
		Version: goVersion,
		Licence: Licence(licenceText),
	}
}

// GoEnv answers the value of a variable of go env, run in root.
func GoEnv(root, name string) (string, error) {
	out, err := run(root, "go", "env", name)
	return strings.TrimSpace(out), err
}

// NpmPackages answers every package of a package-lock.json that is not marked dev, without the
// package of the subsystem - the entry "" - with the licence it names.
//
// @implements SR_37@2.1.1
func NpmPackages(lock []byte) ([]Package, error) {
	var file struct {
		Packages map[string]struct {
			Name    string          `json:"name"`
			Version string          `json:"version"`
			License json.RawMessage `json:"license"`
			Dev     bool            `json:"dev"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(lock, &file); err != nil {
		return nil, fmt.Errorf("package-lock.json: %w", err)
	}
	var packages []Package
	for key, p := range file.Packages {
		if key == "" || p.Dev {
			continue
		}
		name := p.Name
		if name == "" {
			name = key[strings.LastIndex(key, "node_modules/")+len("node_modules/"):]
		}
		var licence string
		_ = json.Unmarshal(p.License, &licence) // a licence that is no string stays unknown
		packages = append(packages, Package{
			PURL:    "pkg:npm/" + strings.Replace(name, "@", "%40", 1) + "@" + p.Version,
			Name:    name,
			Version: p.Version,
			Licence: Given(licence),
			Named:   licence,
		})
	}
	return packages, nil
}

// ApkInstalled reads /lib/apk/db/installed of an image, the last stage of its build.
//
// @implements SR_38@1.2.1
func ApkInstalled(image string) ([]byte, error) {
	out, err := run("", "docker", "run", "--rm", "--entrypoint", "cat", image, "/lib/apk/db/installed")
	return []byte(out), err
}

// ApkPackages answers every package of an /lib/apk/db/installed, with the licence its field L names.
//
// @implements SR_38@1.2.1
func ApkPackages(installed []byte) []Package {
	var packages []Package
	var name, version, licence string
	flush := func() {
		if name != "" {
			packages = append(packages, Package{
				PURL:    "pkg:apk/alpine/" + name + "@" + version,
				Name:    name,
				Version: version,
				Licence: Given(licence),
				Named:   licence,
			})
		}
		name, version, licence = "", "", ""
	}
	scanner := bufio.NewScanner(strings.NewReader(string(installed)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			flush()
		case strings.HasPrefix(line, "P:"):
			name = line[2:]
		case strings.HasPrefix(line, "V:"):
			version = line[2:]
		case strings.HasPrefix(line, "L:"):
			licence = line[2:]
		}
	}
	flush()
	return packages
}

func run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, exit.Stderr)
		}
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return string(out), nil
}
