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

// Row is the row of a function in the table of a section.
type Row struct {
	Function
	InOutput   bool
	Hash       string
	Called     bool
	Ran, Total int
}

// Rows are the rows of the functions: each with whether it is in the binary, the hash of its machine
// code, and its coverage.
//
// @implements SR_117@1.1.1
func Rows(functions []Function, d Disassembly, p *Profile) []Row {

	var rows []Row
	for _, f := range functions {
		r := Row{Function: f}
		r.Hash, r.InOutput = d.Hash(f.Symbol)
		r.Called, r.Ran, r.Total = p.Coverage(f)
		rows = append(rows, r)
	}
	return rows
}

// Section is the section of a Go binary: its head, then the table of its functions.
//
// @implements SR_115@1.1.1, SR_117@1.1.1
func Section(b Binary, rows []Row) string {

	var s strings.Builder
	fmt.Fprintf(&s, "## %s\n\n", b.Name)
	fmt.Fprintf(&s, "- Binary: %s (%s)\n", b.Name, b.Main)
	fmt.Fprintf(&s, "- Go %s, %s/%s, commit %s\n\n", b.GoVersion, b.OS, b.Arch, b.Commit)
	s.WriteString("| Function | Symbol | In binary | Hash | Called by tests | Statements |\n")
	s.WriteString("|---|---|---|---|---|---|\n")
	for _, r := range rows {
		symbol, in, hash := "-", "no", "-"
		if r.InOutput {
			symbol, in, hash = r.Symbol, "yes", r.Hash
		}
		called := "no"
		if r.Called {
			called = "yes"
		}
		fmt.Fprintf(&s, "| %s | %s | %s | %s | %s | %d/%d |\n", r.Key, symbol, in, hash, called, r.Ran, r.Total)
	}
	return s.String()
}
