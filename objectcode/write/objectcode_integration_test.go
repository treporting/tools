//go:build integration

package write

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

// module writes a module example.org/m with the files given, by their paths from its root.
func module(t *testing.T, files map[string]string) string {

	t.Helper()
	dir := t.TempDir()
	files["go.mod"] = "module example.org/m\n\ngo 1.25\n"
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// run runs the go command in dir, for linux/amd64, so that the machine code is the same on every host.
func run(t *testing.T, dir string, args ...string) {

	t.Helper()
	runOn(t, dir, []string{"GOOS=linux", "GOARCH=amd64"}, args...)
}

// runOn runs the go command in dir with the environment given beside that of the test.
func runOn(t *testing.T, dir string, env []string, args ...string) {

	t.Helper()
	cmd := exec.Command("go", args...) // #nosec G204 -- the arguments come from the test
	cmd.Dir = dir
	cmd.Env = append(append(os.Environ(), "GOFLAGS=-buildvcs=false", "CGO_ENABLED=0"), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go %v: %v\n%s", args, err, out)
	}
}

// build builds the main package of dir as dir/x and reads its functions and machine code.
func build(t *testing.T, dir string) ([]Function, Disassembly) {

	t.Helper()
	run(t, dir, "build", "-o", "x", ".")
	packages, err := ListPackages(dir, ".")
	if err != nil {
		t.Fatal(err)
	}
	functions, err := Functions(dir, packages)
	if err != nil {
		t.Fatal(err)
	}
	d, err := Objdump(filepath.Join(dir, "x"), packages)
	if err != nil {
		t.Fatal(err)
	}
	return functions, d
}

// @implements TC_178@1.1.1
func TestTheSectionNamesTheBinary(t *testing.T) {

	dir := module(t, map[string]string{"main.go": "package main\n\nfunc main() {}\n"})
	run(t, dir, "build", "-o", "x", ".")
	b, err := ReadBinary(filepath.Join(dir, "x"), "961bae7c0ffee")
	if err != nil {
		t.Fatal(err)
	}
	if b.GoVersion != runtime.Version() || b.OS != "linux" || b.Arch != "amd64" || b.Commit != "961bae7" {
		t.Errorf("the binary is %+v, want %s, linux/amd64, commit 961bae7", b, runtime.Version())
	}
}

const called = `package main

import "fmt"

//go:noinline
func Called(n int) string { return fmt.Sprint(n * 3) }

//go:noinline
func NotCalled(n int) string { return fmt.Sprint(n + 1) }

func inlined(n int) int { return n + 2 }

func main() {
	fmt.Println(Called(inlined(1)), NotCalled(2))
}
`

const calledTest = `package main

import "testing"

func TestCalled(t *testing.T) {
	if Called(1) != "3" {
		t.Fail()
	}
}
`

// @implements TC_181@1.1.1
func TestTheTableNamesEachFunction(t *testing.T) {

	dir := module(t, map[string]string{"main.go": called, "main_test.go": calledTest})
	functions, d := build(t, dir)
	//	the tests run on the host; the coverage profile names the same files there
	runOn(t, dir, nil, "test", "-coverprofile", "cover.out", ".")
	f, err := os.Open(filepath.Join(dir, "cover.out")) // #nosec G304 -- a file of the test
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	p, err := ReadProfile(f)
	if err != nil {
		t.Fatal(err)
	}
	rows := map[string]Row{}
	for _, r := range Rows(functions, d, p) {
		rows[r.Key] = r
	}
	for _, tt := range []struct {
		key        string
		in, called bool
		ran, total int
	}{
		{key: "main.go@Called", in: true, called: true, ran: 1, total: 1},
		{key: "main.go@NotCalled", in: true, called: false, ran: 0, total: 1},
		{key: "main.go@inlined", in: false, called: false, ran: 0, total: 1},
	} {
		r, ok := rows[tt.key]
		if !ok {
			t.Errorf("no row for %s", tt.key)
			continue
		}
		if r.InOutput != tt.in || (r.Hash != "") != tt.in || r.Called != tt.called || r.Ran != tt.ran || r.Total != tt.total {
			t.Errorf("%s: in binary %v, hash %q, called %v, %d/%d", tt.key, r.InOutput, r.Hash, r.Called, r.Ran, r.Total)
		}
	}
}

// @implements TC_183@1.1.1
func TestTheHashStaysWhereOnlyAddressesMove(t *testing.T) {

	const f = `
//go:noinline
func F(n int) string { return fmt.Sprintf("%d items", n*7) }
`
	const before = `
//go:noinline
func Before(n int) int { return n*n + 11 }
`
	const main = `
func main() { fmt.Println(F(len(os.Args))) }
`
	const main2 = `
func main() { fmt.Println(F(len(os.Args)), Before(3)) }
`
	head := "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n)\n"
	one := module(t, map[string]string{"main.go": head + f + main})
	other := module(t, map[string]string{"main.go": head + before + f + main2})
	_, d1 := build(t, one)
	_, d2 := build(t, other)
	h1, ok1 := d1.Hash("main.F")
	h2, ok2 := d2.Hash("main.F")
	if !ok1 || !ok2 {
		t.Fatal("F is not in the binary")
	}
	if slices.Equal(d1["main.F"], d2["main.F"]) {
		t.Fatal("the machine code of F is the same in both, with no address moved; the test shows nothing")
	}
	if h1 != h2 {
		t.Errorf("F has the hash %s, and %s with a function before it, want the same", h1, h2)
	}
}
