package write

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// block is a block of statements of a coverage profile.
type block struct {
	file       string
	from, to   int
	statements int
	count      int
}

// Profile is a coverage profile of `go test -coverprofile`; a block listed more than once, as
// -coverpkg lists it for each test binary, counts with its highest count.
type Profile struct {
	blocks map[string]*block
}

// ReadProfile reads a coverage profile.
func ReadProfile(r io.Reader) (*Profile, error) {

	p := &Profile{blocks: map[string]*block{}}
	lines := bufio.NewScanner(r)
	lines.Buffer(make([]byte, 1<<20), 1<<24)
	for lines.Scan() {
		line := lines.Text()
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		//	file.go:from.col,to.col statements count
		fields := strings.Fields(line)
		file, span, ok := strings.Cut(fields[0], ":")
		if len(fields) != 3 || !ok {
			return nil, fmt.Errorf("a line of a coverage profile: %q", line)
		}
		start, end, _ := strings.Cut(span, ",")
		from, err1 := strconv.Atoi(strings.Split(start, ".")[0])
		to, err2 := strconv.Atoi(strings.Split(end, ".")[0])
		statements, err3 := strconv.Atoi(fields[1])
		count, err4 := strconv.Atoi(fields[2])
		for _, err := range []error{err1, err2, err3, err4} {
			if err != nil {
				return nil, fmt.Errorf("a line of a coverage profile: %q: %w", line, err)
			}
		}
		key := fields[0]
		if b, seen := p.blocks[key]; seen {
			b.count = max(b.count, count)
			continue
		}
		p.blocks[key] = &block{file: file, from: from, to: to, statements: statements, count: count}
	}
	return p, lines.Err()
}

// Coverage is whether the tests call a function, and how many of its statements they ran.
//
// @implements SR_117@1.1.1
func (p *Profile) Coverage(f Function) (called bool, ran, statements int) {

	if p == nil {
		return false, 0, 0
	}
	for _, b := range p.blocks {
		if b.file != f.File || b.from < f.From || b.to > f.To {
			continue
		}
		statements += b.statements
		if b.count > 0 {
			called = true
			ran += b.statements
		}
	}
	return called, ran, statements
}
