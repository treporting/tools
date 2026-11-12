package approval

import (
	"reflect"
	"testing"
)

// @verifies SR_79@2.1.2
func TestBlobHashIsTheHashOfGit(t *testing.T) {
	//	git hash-object of a file holding "hello\n"
	if got := blobHash([]byte("hello\n")); got != "ce013625030ba8dba906f756967f9e9ca394464a" {
		t.Errorf("blobHash = %s", got)
	}
}

// @verifies SR_78@1.1.1
func TestFindXReadsTheAttributesForXAlsoInsideOneOf(t *testing.T) {
	plain := map[string]any{"attributes_for_version_x": []any{"object_text", "safety"}}
	if got := findX(plain); !reflect.DeepEqual(got, []string{"object_text", "safety"}) {
		t.Errorf("plain: %v", got)
	}
	nested := map[string]any{"oneOf": []any{map[string]any{"attributes_for_version_x": []any{"object_text"}}}}
	if got := findX(nested); !reflect.DeepEqual(got, []string{"object_text"}) {
		t.Errorf("oneOf: %v", got)
	}
}

// @verifies SR_78@1.1.1
func TestXOfKeepsOnlyTheAttributesForX(t *testing.T) {
	got := xOf([]byte("object_key: TOOL_01\nobject_text: a text\npriority: high\n"), []string{"object_text", "safety"})
	want := map[string]any{"object_text": "a text", "safety": nil}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("xOf = %v, want %v", got, want)
	}
}

// @verifies SR_80@2.1.3
func TestInFilesKeepsTheObjectsOfTheChangedFiles(t *testing.T) {
	objects := []object{{key: "TOOL_01", file: "a.yaml"}, {key: "SR_01", file: "b.yaml"}}
	got := inFiles(objects, []string{"b.yaml"})
	if len(got) != 1 || got[0].key != "SR_01" {
		t.Errorf("inFiles = %v", got)
	}
}
