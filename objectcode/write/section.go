package write

import (
	"debug/buildinfo"
	"fmt"
	"strings"
)

// Binary is what the head of the section of a Go binary names.
type Binary struct {
	// Name is the file of the binary, Main the folder of its main package from the root of the repo.
	Name, Main string
	// GoVersion, such as go1.25.14; OS and Arch, such as linux and amd64.
	GoVersion, OS, Arch string
	// Commit is the commit it was built from.
	Commit string
}

// ReadBinary reads the Go version, the OS/arch and the commit from the build information of a Go
// binary; commit is taken where the binary names none.
//
// @implements SR_115@1.1.1
func ReadBinary(path, commit string) (Binary, error) {

	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return Binary{}, err
	}
	b := Binary{GoVersion: info.GoVersion, Commit: commit}
	for _, s := range info.Settings {
		switch s.Key {
		case "GOOS":
			b.OS = s.Value
		case "GOARCH":
			b.Arch = s.Value
		case "vcs.revision":
			if b.Commit == "" {
				b.Commit = s.Value
			}
		}
	}
	if len(b.Commit) > 7 {
		b.Commit = b.Commit[:7]
	}
	return b, nil
}

// Head is the head of the section of a build output.
type Head interface {
	// Title is the name of the output; Kind what it is, binary or bundle; Lines what the head names.
	Title() string
	Kind() string
	Lines() []string
}

// Title is the file of the binary.
func (b Binary) Title() string { return b.Name }

// Kind is binary.
func (b Binary) Kind() string { return "binary" }

// Lines name the binary and its main package, the Go version, the OS/arch and the commit.
//
// @implements SR_115@1.1.1
func (b Binary) Lines() []string {
	return []string{
		fmt.Sprintf("Binary: %s (%s)", b.Name, b.Main),
		fmt.Sprintf("Go %s, %s/%s, commit %s", b.GoVersion, b.OS, b.Arch, b.Commit),
	}
}

// Title is the file of the bundle.
func (b Bundle) Title() string { return b.Name }

// Kind is bundle.
func (b Bundle) Kind() string { return "bundle" }

// Lines name the bundle and its project, node and each tool with its version, the hash of their
// configuration and the commit.
//
// @implements SR_116@1.1.1
func (b Bundle) Lines() []string {
	versions := []string{"node " + b.Node}
	for _, t := range b.Tools {
		versions = append(versions, t.Name+" "+t.Version)
	}
	return []string{
		fmt.Sprintf("Bundle: %s (%s)", b.Name, b.Dir),
		fmt.Sprintf("%s, configuration %s, commit %s", strings.Join(versions, ", "), b.Configuration, b.Commit),
	}
}

// Output is a build output: it finds a function of the source in it.
type Output interface {
	Find(f Function) (name, hash string, ok bool)
}

// Coverage is the coverage of the tests of the source.
type Coverage interface {
	Coverage(f Function) (called bool, ran, statements int)
}

// Row is the row of a function in the table of a section.
type Row struct {
	Function
	// Name is its name in the output.
	Name       string
	InOutput   bool
	Hash       string
	Called     bool
	Ran, Total int
}

// Rows are the rows of the functions: each with its name in the output, whether it is there, the
// hash of its code there, and its coverage.
//
// @implements SR_117@2.1.1
func Rows(functions []Function, out Output, c Coverage) []Row {

	var rows []Row
	for _, f := range functions {
		r := Row{Function: f}
		r.Name, r.Hash, r.InOutput = out.Find(f)
		r.Called, r.Ran, r.Total = c.Coverage(f)
		rows = append(rows, r)
	}
	return rows
}

// Section is the section of a build output: its head, then the table of its functions.
//
// @implements SR_115@1.1.1, SR_116@1.1.1, SR_117@2.1.1
func Section(h Head, rows []Row) string {

	var s strings.Builder
	fmt.Fprintf(&s, "## %s\n\n", h.Title())
	for _, line := range h.Lines() {
		fmt.Fprintf(&s, "- %s\n", line)
	}
	fmt.Fprintf(&s, "\n| Function | Symbol | In %s | Hash | Called by tests | Statements |\n", h.Kind())
	s.WriteString("|---|---|---|---|---|---|\n")
	for _, r := range rows {
		name, in, hash := "-", "no", "-"
		if r.InOutput {
			name, in, hash = r.Name, "yes", r.Hash
		}
		called := "no"
		if r.Called {
			called = "yes"
		}
		fmt.Fprintf(&s, "| %s | %s | %s | %s | %s | %d/%d |\n", r.Key, name, in, hash, called, r.Ran, r.Total)
	}
	return s.String()
}

// Note is the section of a build output whose functions cannot be listed: its head, then the note
// that says why.
//
// @implements SR_116@1.1.1
func Note(h Head, note string) string {

	var s strings.Builder
	fmt.Fprintf(&s, "## %s\n\n", h.Title())
	for _, line := range h.Lines() {
		fmt.Fprintf(&s, "- %s\n", line)
	}
	fmt.Fprintf(&s, "\n%s\n", note)
	return s.String()
}
