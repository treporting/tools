// Package write writes the object code report of a build output (TOOL_35): per function of the
// source code whether it is in the output, the hash of its code there, and the coverage of the tests.
package write

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// Disassembly is the machine code of each symbol of a Go binary, as `go tool objdump` prints it: the
// instructions only, without file, line, address and bytes.
type Disassembly map[string][]string

// Objdump reads the machine code of the symbols of the packages from a binary, by `go tool objdump`.
func Objdump(binary string, packages []Package) (Disassembly, error) {

	var names []string
	for _, p := range packages {
		name := regexp.QuoteMeta(p.ImportPath)
		if p.Name == "main" {
			name = "main"
		}
		names = append(names, name)
	}
	cmd := exec.Command("go", "tool", "objdump", "-s", `^(`+strings.Join(names, "|")+`)\.`, binary) // #nosec G204 -- the binary comes from the user of the tool
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go tool objdump %s: %w", binary, err)
	}
	return ReadObjdump(&out)
}

// ReadObjdump reads the output of `go tool objdump`.
func ReadObjdump(r io.Reader) (Disassembly, error) {

	d := Disassembly{}
	symbol := ""
	lines := bufio.NewScanner(r)
	lines.Buffer(make([]byte, 1<<20), 1<<24)
	for lines.Scan() {
		line := lines.Text()
		if name, ok := strings.CutPrefix(line, "TEXT "); ok {
			symbol, _, _ = strings.Cut(name, "(SB)")
			d[symbol] = d[symbol][:0:0]
			continue
		}
		if symbol == "" {
			continue
		}
		var fields []string
		for _, f := range strings.Split(line, "\t") {
			if f = strings.TrimSpace(f); f != "" {
				fields = append(fields, f)
			}
		}
		//	file:line, address, bytes, instruction
		if len(fields) >= 4 {
			d[symbol] = append(d[symbol], strings.Join(fields[3:], " "))
		}
	}
	return d, lines.Err()
}

var (
	// a load relative to the instruction pointer: unnamed data, such as a string or a type
	relative = regexp.MustCompile(`-?0x[0-9a-f]+\(IP\)`)
	// a reference to a symbol, with an offset where it has one
	reference = regexp.MustCompile(`(\$?)([^\s,()$]+?)(?:[+-](?:0x)?[0-9a-f]+)?\(SB\)`)
	// an address as an operand: a jump target, without $ and not an offset of a register
	address = regexp.MustCompile(`(^|[ ,])0x[0-9a-f]+([ ,]|$)`)
)

// placeholder stands for every reference to unnamed data.
const placeholder = "·"

// unnamed is whether a symbol names data the linker makes and names on its own, such as a type
// descriptor or string data, rather than a function or a variable of the source.
func unnamed(symbol string) bool {
	for _, prefix := range []string{"type:", "typerel.", "go:", "runtime.rodata", "$f32.", "$f64.", "$i64."} {
		if strings.HasPrefix(symbol, prefix) {
			return true
		}
	}
	return false
}

// Normalize is an instruction without what the linker moves alone: every address is removed, a
// reference to a named symbol keeps the name and loses the offset, a reference to unnamed data
// becomes one placeholder; instruction, registers, immediates and stack offsets stay.
//
// @implements SR_118@1.1.1
func Normalize(instruction string) string {

	s := relative.ReplaceAllString(instruction, placeholder+"(IP)")
	s = reference.ReplaceAllStringFunc(s, func(m string) string {
		parts := reference.FindStringSubmatch(m)
		if unnamed(parts[2]) {
			return parts[1] + placeholder + "(SB)"
		}
		return parts[1] + parts[2] + "(SB)"
	})
	for {
		next := address.ReplaceAllString(s, "$1$2")
		if next == s {
			break
		}
		s = next
	}
	return strings.Trim(strings.Join(strings.Fields(s), " "), " ,")
}

// Hash is the hash of the machine code of a function: over its normalized instructions, and over
// those of its closures, in the order of their names.
//
// @implements SR_118@1.1.1
func (d Disassembly) Hash(symbol string) (string, bool) {

	var symbols []string
	for s := range d {
		if s == symbol || strings.HasPrefix(s, symbol+".func") || strings.HasPrefix(s, symbol+"[") {
			symbols = append(symbols, s)
		}
	}
	if len(symbols) == 0 {
		return "", false
	}
	sort.Strings(symbols)
	h := sha256.New()
	for _, s := range symbols {
		for _, instruction := range d[s] {
			_, _ = io.WriteString(h, Normalize(instruction))
			_, _ = io.WriteString(h, "\n")
		}
	}
	return hex.EncodeToString(h.Sum(nil))[:8], true
}
