package check

import (
	"reflect"
	"testing"
)

func parts(f files) []string {
	var out []string
	for _, finding := range Parts(f.read, f.index()) {
		out = append(out, finding.String())
	}
	return out
}

const r3 = `type: architecture
level: r3
elements:
- object_key: ARC_1
  kind: subsystem
  repo: trace/reporter
  object_text: 'REST API: answers the trace'
- object_key: ARC_2
  kind: interface
  object_text: '` + "`StorageSource`" + `'
requirements:
- object_key: ARC_10
  component: trace/reporter
  interface: StorageSource
`

// @implements TC_201@1.1.1, TC_205@1.1.1
func TestAPartAndAnInterfaceOfR3Pass(t *testing.T) {
	if got := parts(files{"architecture/system-architecture.yaml": r3}); len(got) != 0 {
		t.Errorf("findings %q, want none", got)
	}
}

// @implements TC_202@1.1.1, TC_206@1.1.1
func TestAComponentAndAnInterfaceR3DoesNotHaveAreReported(t *testing.T) {
	got := parts(files{"architecture/system-architecture.yaml": r3 + `- object_key: ARC_11
  component: reporter-trace
- object_key: ARC_12
  component: trace/reporter
  interface: REST API
`})
	want := []string{
		"ARC_11 names component reporter-trace: no subsystem of r3 has this repo",
		"ARC_12 names interface REST API: no interface of r3 has this name",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findings %q, want %q", got, want)
	}
}

const r6 = `type: architecture
level: r6
elements:
- object_key: SA_1
  kind: component
  object_text: trace
- object_key: SA_2
  kind: interface
  object_text: Gaps
requirements:
- object_key: SA_10
  component: trace
`

const r7 = `type: design
level: r7
elements:
- object_key: SD_1
  kind: module
  object_text: trace/gaps
requirements:
- object_key: SD_10
  component: trace/gaps
`

// @implements TC_203@1.1.1
func TestAPartOfR6AndOfR7Pass(t *testing.T) {
	got := parts(files{
		"reporter/architecture/software-architecture.yaml": r6,
		"reporter/design/software-design.yaml":             r7,
	})
	if len(got) != 0 {
		t.Errorf("findings %q, want none", got)
	}
}

// @implements TC_204@1.1.1
func TestAComponentR6AndR7DoNotHaveIsReported(t *testing.T) {
	got := parts(files{
		"reporter/architecture/software-architecture.yaml": r6 + `- object_key: SA_11
  component: maturity
`,
		"reporter/design/software-design.yaml": r7 + `- object_key: SD_11
  component: trace/cycles
`,
	})
	want := []string{
		"SA_11 names component maturity: no component of r6 has this name",
		"SD_11 names component trace/cycles: no module of r7 has this name",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findings %q, want %q", got, want)
	}
}

// @implements TC_207@1.1.1
func TestAnInterfaceR6DoesNotHaveIsReported(t *testing.T) {
	got := parts(files{"reporter/architecture/software-architecture.yaml": r6 + `- object_key: SA_12
  component: trace
  interface: Cycles
- object_key: SA_13
  component: trace
`})
	want := []string{"SA_12 names interface Cycles: no interface of r6 has this name"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findings %q, want %q", got, want)
	}
}
