// Package create creates the SBOM of a subsystem as sbom.md (TOOL_21): every package it delivers,
// its own and third-party, with its name, version, licence and PURL.
package create

import (
	"bufio"
	"sort"
	"strings"
)

// Unknown is the licence of a package whose licence is not recognized.
const Unknown = "unknown"

// Package is one row of the table of sbom.md.
type Package struct {
	PURL    string
	Name    string
	Version string
	Licence string
	Named   string // the licence as the package names it itself, where it does
}

// SBOM is what sbom.md is written from.
type SBOM struct {
	Repo     string     // the repo of the subsystem, such as storages
	Packages []Package  // every package it delivers, its own among them
	Original []Original // every file of original data it is written from
	History  []string   // the rows of the Change history, the new one last
}

// Own answers the package of the subsystem itself: the label of its baseline as its version, the
// licence of the LICENSE file in the root of its repo, and its PURL on GitHub.
//
// @implements SR_46@1.1.1
func Own(repo, label, licenceText string) Package {
	return Package{
		PURL:    "pkg:github/treporting/" + repo + "@" + label,
		Name:    repo,
		Version: label,
		Licence: Licence(licenceText),
	}
}

// Once answers the packages with each PURL once, in the order they come first; two Go modules or
// two package-lock.json of a repo may name the same package.
func Once(packages []Package) []Package {
	seen := map[string]bool{}
	var once []Package
	for _, p := range packages {
		if !seen[p.PURL] {
			seen[p.PURL] = true
			once = append(once, p)
		}
	}
	return once
}

// Row answers the row of the Change history that a run adds: the version of sbom.md after the run;
// the day it runs; the label of the baseline; the ticket, as #<number>; the tool with its version.
//
// @implements SR_40@3.3.1
func Row(version, date, label, ticket, by string) string {
	return "| " + cell(version) + " | " + date + " | " + cell(label) + " | #" + cell(ticket) + " | " + cell(by) + " |"
}

// History answers the rows of the Change history of an sbom.md as it stands, without its head, so
// that a run keeps them. A file that is not there has none.
//
// @implements SR_40@3.3.1
func History(sbomMD []byte) []string {
	var rows []string
	in := false
	scanner := bufio.NewScanner(strings.NewReader(string(sbomMD)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "## "):
			in = line == "## Change history"
		case in && strings.HasPrefix(line, "| ") && !strings.HasPrefix(line, "| Version |"):
			rows = append(rows, line)
		}
	}
	return rows
}

// Markdown answers sbom.md: the sections Original data, Packages - the table of the packages sorted
// by PURL -, Open and Change history.
//
// @implements SR_41@2.1.1
// @implements SR_42@1.2.1
// @implements SR_45@1.1.1
// @implements SR_101@2.1.1
func Markdown(s SBOM) string {
	packages := append([]Package(nil), s.Packages...)
	sort.Slice(packages, func(i, j int) bool { return packages[i].PURL < packages[j].PURL })

	var b strings.Builder
	b.WriteString("# Software Bill of Materials (SBOM)\n\n")
	b.WriteString("The packages " + s.Repo + " delivers, its own and third-party, each with its version, its\n")
	b.WriteString("licence and its PURL. It is written by the tool sbom (TOOL_21) and not edited by hand.\n")
	b.WriteString("\n## Original data\n\n")
	if len(s.Original) == 0 {
		b.WriteString("`<empty>`\n")
	} else {
		b.WriteString("**Table original data** - every file sbom.md is written from\n\n")
		b.WriteString("| File | Made by | Stored At |\n|---|---|---|\n")
	}
	for _, o := range s.Original {
		b.WriteString("| " + cell(o.File) + " | `" + cell(o.MadeBy) + "` | " + cell(o.StoredAt) + " |\n")
	}
	b.WriteString("\n## Packages\n\n")
	b.WriteString("**Table packages** - every package " + s.Repo + " delivers\n\n")
	b.WriteString("| PURL | Name | Version | Licence |\n|---|---|---|---|\n")
	var unknown []Package
	for _, p := range packages {
		b.WriteString("| " + cell(p.PURL) + " | " + cell(p.Name) + " | " + cell(p.Version) + " | " + cell(p.Licence) + " |\n")
		if p.Licence == Unknown {
			unknown = append(unknown, p)
		}
	}
	b.WriteString("\n## Open\n\n")
	if len(unknown) == 0 {
		b.WriteString("`<empty>`\n")
	}
	for _, p := range unknown {
		b.WriteString("- `" + p.PURL + "` - the licence is unknown")
		if p.Named != "" {
			b.WriteString("; the package names it \"" + p.Named + "\"")
		}
		b.WriteString("\n")
	}
	b.WriteString("\n## Change history\n\n")
	b.WriteString("| Version | Date | Change | Reason | By |\n|---|---|---|---|---|\n")
	for _, row := range s.History {
		b.WriteString(row + "\n")
	}
	return b.String()
}

// cell keeps a value from breaking the row of a table.
func cell(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
}

// Table answers the packages of the table of an sbom.md, as it lists them.
//
// @implements SR_66@1.1.1
func Table(sbomMD []byte) []Package {
	var packages []Package
	in := false
	scanner := bufio.NewScanner(strings.NewReader(string(sbomMD)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "| PURL |"):
			in = true
		case in && strings.HasPrefix(line, "|---"):
		case in && strings.HasPrefix(line, "| "):
			fields := cells(line)
			if len(fields) == 4 {
				packages = append(packages, Package{PURL: fields[0], Name: fields[1], Version: fields[2], Licence: fields[3]})
			}
		default:
			in = false
		}
	}
	return packages
}

// Label answers the label of the baseline of an sbom.md: the column Change of the last row of its
// Change history; nothing where it has no row.
//
// @implements SR_66@1.1.1
func Label(sbomMD []byte) string {
	rows := History(sbomMD)
	if len(rows) == 0 {
		return ""
	}
	if fields := cells(rows[len(rows)-1]); len(fields) == 5 {
		return fields[2]
	}
	return ""
}

// Differ answers a line for each package that is in one list and not in the other, or in both with
// another name, version or licence; nothing where both list the same.
//
// @implements SR_66@1.1.1
func Differ(delivered, listed []Package) []string {
	byPURL := func(packages []Package) map[string]Package {
		m := map[string]Package{}
		for _, p := range packages {
			m[p.PURL] = Package{PURL: p.PURL, Name: p.Name, Version: p.Version, Licence: p.Licence}
		}
		return m
	}
	want, have := byPURL(delivered), byPURL(listed)
	var lines []string
	for purl, p := range want {
		q, ok := have[purl]
		switch {
		case !ok:
			lines = append(lines, purl+": delivered, not in sbom.md")
		case p != q:
			lines = append(lines, purl+": sbom.md lists "+q.Name+", "+q.Version+", "+q.Licence+"; delivered is "+p.Name+", "+p.Version+", "+p.Licence)
		}
	}
	for purl := range have {
		if _, ok := want[purl]; !ok {
			lines = append(lines, purl+": in sbom.md, not delivered")
		}
	}
	sort.Strings(lines)
	return lines
}

// cells answers the values of a row of a table, as cell wrote them.
func cells(row string) []string {
	row = strings.TrimSuffix(strings.TrimPrefix(row, "| "), " |")
	parts := strings.Split(strings.ReplaceAll(row, `\|`, "\x00"), " | ")
	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], "\x00", "|")
	}
	return parts
}
