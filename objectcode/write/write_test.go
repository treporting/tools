package write

import (
	"strings"
	"testing"
)

// objdump is the output of `go tool objdump` for two functions, as it prints it.
const objdump = `TEXT example.org/m/a.F(SB) /src/a/a.go
  a.go:3		0x4a1b20		493b6610		CMPQ SP, 0x10(R14)
  a.go:3		0x4a1b24		7620			JBE 0x4a1b46
  a.go:4		0x4a1b26		4883ec28		SUBQ $0x28, SP
  a.go:4		0x4a1b2a		488d05cf1e0100		LEAQ 0x11ecf(IP), AX
  a.go:5		0x4a1b31		e8aa000000		CALL fmt.Sprintf(SB)
TEXT example.org/m/a.F.func1(SB) /src/a/a.go
  a.go:6		0x4a1b60		c3			RET
TEXT example.org/m/a.G(SB) /src/a/a.go
  a.go:9		0x4a1b80		c3			RET
`

func TestReadObjdump(t *testing.T) {

	d, err := ReadObjdump(strings.NewReader(objdump))
	if err != nil {
		t.Fatal(err)
	}
	if got := d["example.org/m/a.F"]; len(got) != 5 || got[4] != "CALL fmt.Sprintf(SB)" {
		t.Errorf("F has %q", got)
	}
	if got := d["example.org/m/a.G"]; len(got) != 1 || got[0] != "RET" {
		t.Errorf("G has %q", got)
	}
}

func TestNormalize(t *testing.T) {

	for _, tt := range []struct{ in, want string }{
		{"JBE 0x4a1b46", "JBE"},
		{"LEAQ 0x11ecf(IP), AX", "LEAQ ·(IP), AX"},
		{"MOVQ example.org/m/a.list+8(SB), SI", "MOVQ example.org/m/a.list(SB), SI"},
		{"CALL fmt.Sprintf(SB)", "CALL fmt.Sprintf(SB)"},
		{"SUBQ $0x28, SP", "SUBQ $0x28, SP"},
		{"MOVQ AX, 0x60(SP)", "MOVQ AX, 0x60(SP)"},
		{"CMPQ SP, 0x10(R14)", "CMPQ SP, 0x10(R14)"},
		{"LEAQ go:string.*+1234(SB), AX", "LEAQ ·(SB), AX"},
	} {
		if got := Normalize(tt.in); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// hashOf is the hash of a function of the given instructions.
func hashOf(t *testing.T, instructions ...string) string {

	t.Helper()
	h, ok := Disassembly{"p.F": instructions}.Hash("p.F")
	if !ok {
		t.Fatal("p.F is not in the binary")
	}
	return h
}

// @implements TC_184@1.1.1
func TestHashIgnoresTheNamesOfUnnamedData(t *testing.T) {

	one := hashOf(t, "LEAQ type:*+1234(SB), AX", "CALL runtime.newobject(SB)")
	other := hashOf(t, "LEAQ typerel.*+5678(SB), AX", "CALL runtime.newobject(SB)")
	if one != other {
		t.Errorf("the hash is %s with _type:* and %s with typerel.*, want the same", one, other)
	}
}

// @implements TC_185@1.1.1
func TestHashSeesAnotherCallAndAnotherConstant(t *testing.T) {

	sprintf := hashOf(t, "MOVL $0x1, AX", "CALL fmt.Sprintf(SB)")
	if sprint := hashOf(t, "MOVL $0x1, AX", "CALL fmt.Sprint(SB)"); sprint == sprintf {
		t.Errorf("the call of fmt.Sprint has the hash of fmt.Sprintf, %s", sprint)
	}
	if two := hashOf(t, "MOVL $0x2, AX", "CALL fmt.Sprintf(SB)"); two == sprintf {
		t.Errorf("another constant has the same hash, %s", two)
	}
}

func TestHashTakesTheClosures(t *testing.T) {

	d, err := ReadObjdump(strings.NewReader(objdump))
	if err != nil {
		t.Fatal(err)
	}
	with, _ := d.Hash("example.org/m/a.F")
	d["example.org/m/a.F.func1"] = []string{"NOPL", "RET"}
	if changed, _ := d.Hash("example.org/m/a.F"); changed == with {
		t.Error("a change of a closure leaves the hash of its function")
	}
	if _, ok := d.Hash("example.org/m/a.H"); ok {
		t.Error("a function not in the binary has a hash")
	}
}

func TestCoverage(t *testing.T) {

	profile := `mode: set
example.org/m/a/a.go:3.14,5.2 2 1
example.org/m/a/a.go:5.2,6.3 1 0
example.org/m/a/a.go:9.14,10.2 1 0
example.org/m/a/a.go:9.14,10.2 1 1
example.org/m/a/a.go:12.14,13.2 1 0
`
	p, err := ReadProfile(strings.NewReader(profile))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		f          Function
		called     bool
		ran, total int
	}{
		{Function{File: "example.org/m/a/a.go", From: 3, To: 7}, true, 2, 3},
		{Function{File: "example.org/m/a/a.go", From: 9, To: 10}, true, 1, 1},
		{Function{File: "example.org/m/a/a.go", From: 12, To: 13}, false, 0, 1},
	} {
		called, ran, total := p.Coverage(tt.f)
		if called != tt.called || ran != tt.ran || total != tt.total {
			t.Errorf("lines %d to %d: %v %d/%d, want %v %d/%d", tt.f.From, tt.f.To, called, ran, total, tt.called, tt.ran, tt.total)
		}
	}
}

func TestSection(t *testing.T) {

	b := Binary{Name: "x", Main: "cmd/x", GoVersion: "go1.25.14", OS: "linux", Arch: "amd64", Commit: "961bae7"}
	rows := []Row{
		{Function: Function{Key: "a/a.go@F", Symbol: "example.org/m/a.F"}, Name: "example.org/m/a.F", InOutput: true, Hash: "75ebe987", Called: true, Ran: 15, Total: 16},
		{Function: Function{Key: "a/a.go@named", Symbol: "example.org/m/a.named"}, Called: true, Ran: 3, Total: 3},
	}
	want := `## x

- Binary: x (cmd/x)
- Go go1.25.14, linux/amd64, commit 961bae7

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| a/a.go@F | example.org/m/a.F | yes | 75ebe987 | yes | 15/16 |
| a/a.go@named | - | no | - | yes | 3/3 |
`
	if got := Section(b, rows); got != want {
		t.Errorf("the section is\n%s\nwant\n%s", got, want)
	}
}
