// Package check checks the licences of sbom.md against Apache-2.0 (TOOL_22).
package check

import (
	"bufio"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Unknown is the licence sbom.md names where it does not recognize one.
const Unknown = "unknown"

// Fitting are the licences that fit Apache-2.0 for a package linked into the code of a subsystem.
var Fitting = []string{"Apache-2.0", "MIT", "BSD-2-Clause", "BSD-3-Clause", "ISC", "Zlib", "0BSD", "MPL-2.0"}

// Package is one row of the table of sbom.md, as far as the check needs it.
type Package struct {
	PURL    string
	Licence string
}

// Assessed is one entry of licences_assessed.yaml: the licence a human assessed for one package.
type Assessed struct {
	PURL    string `yaml:"purl"`
	Licence string `yaml:"licence"`
	Reason  string `yaml:"reason"`
}

// Failure is a package that fails the check, with the licence it was checked by.
type Failure struct {
	PURL    string
	Licence string
}

func (f Failure) String() string { return f.PURL + " " + f.Licence }

// Packages answers the rows of the table of sbom.md.
func Packages(sbomMD []byte) []Package {
	var packages []Package
	in := false
	scanner := bufio.NewScanner(strings.NewReader(string(sbomMD)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "| PURL |"):
			in = true
		case !in || strings.HasPrefix(line, "|---"):
		case strings.HasPrefix(line, "| "):
			cells := strings.Split(strings.Trim(line, "|"), " | ")
			if len(cells) == 4 {
				packages = append(packages, Package{PURL: strings.TrimSpace(cells[0]), Licence: strings.TrimSpace(cells[3])})
			}
		default:
			in = false
		}
	}
	return packages
}

// ReadAssessed answers the licences of licences_assessed.yaml by PURL. An entry needs a PURL, a
// licence and a reason.
//
// @implements SR_49@1.1.1
func ReadAssessed(file []byte) (map[string]string, error) {
	var entries []Assessed
	if err := yaml.Unmarshal(file, &entries); err != nil {
		return nil, err
	}
	assessed := map[string]string{}
	for i, e := range entries {
		if e.PURL == "" || e.Licence == "" || e.Reason == "" {
			return nil, fmt.Errorf("licences_assessed.yaml: entry %d needs purl, licence and reason", i+1)
		}
		assessed[e.PURL] = e.Licence
	}
	return assessed, nil
}

// Check answers every package that fails: a package linked into the code - of the type golang, npm
// or github - whose licence does not fit, and a package of an image - of the type apk - whose
// licence is unknown. A licence assessed for a package takes the place of the one in sbom.md.
//
// @implements SR_47@1.1.1
// @implements SR_48@1.1.1
// @implements SR_49@1.1.1
func Check(packages []Package, assessed map[string]string) []Failure {
	var failures []Failure
	for _, p := range packages {
		licence := p.Licence
		if a, ok := assessed[p.PURL]; ok {
			licence = a
		}
		var fails bool
		switch kind(p.PURL) {
		case "golang", "npm", "github":
			fails = !Fits(licence)
		case "apk":
			fails = licence == Unknown
		}
		if fails {
			failures = append(failures, Failure{PURL: p.PURL, Licence: licence})
		}
	}
	return failures
}

// kind answers the type of a PURL, such as golang for pkg:golang/gopkg.in/yaml.v3@v3.0.1.
func kind(purl string) string {
	rest, ok := strings.CutPrefix(purl, "pkg:")
	if !ok {
		return ""
	}
	kind, _, _ := strings.Cut(rest, "/")
	return kind
}

// Fits answers whether a licence expression fits Apache-2.0: of licences joined by AND each must be
// one of Fitting, of licences joined by OR one. An exception joined by WITH, and an expression of
// another form, does not fit.
//
// @implements SR_47@1.1.1
func Fits(expression string) bool {
	tokens := strings.Fields(strings.NewReplacer("(", " ( ", ")", " ) ").Replace(expression))
	p := &parser{tokens: tokens}
	fits, ok := p.or()
	return ok && p.at == len(tokens) && fits
}

// parser reads a licence expression: AND binds closer than OR, as in SPDX.
type parser struct {
	tokens []string
	at     int
}

func (p *parser) next() string {
	if p.at < len(p.tokens) {
		return p.tokens[p.at]
	}
	return ""
}

func (p *parser) or() (bool, bool) {
	fits, ok := p.and()
	for ok && p.next() == "OR" {
		p.at++
		var right bool
		right, ok = p.and()
		fits = fits || right
	}
	return fits, ok
}

func (p *parser) and() (bool, bool) {
	fits, ok := p.term()
	for ok && p.next() == "AND" {
		p.at++
		var right bool
		right, ok = p.term()
		fits = fits && right
	}
	return fits, ok
}

func (p *parser) term() (bool, bool) {
	switch token := p.next(); token {
	case "(":
		p.at++
		fits, ok := p.or()
		if !ok || p.next() != ")" {
			return false, false
		}
		p.at++
		return fits, true
	case "", ")", "AND", "OR", "WITH":
		return false, false
	default:
		p.at++
		if p.next() == "WITH" {
			p.at += 2
			return false, p.at <= len(p.tokens)
		}
		for _, f := range Fitting {
			if token == f {
				return true, true
			}
		}
		return false, true
	}
}
