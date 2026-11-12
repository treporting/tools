package create

import (
	"crypto/sha1" // #nosec G505 -- the hash of a git blob is SHA-1; it names content, it secures nothing
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Name is the name of sbom.md in documents.yaml: its path from the root of the repo.
const Name = "sbom/sbom.md"

// Version answers the version of sbom.md in a documents.yaml, and nothing where it has no entry.
func Version(documentsYAML []byte) string {
	return VersionOf(documentsYAML, Name)
}

// VersionOf answers the version of the document name in a documents.yaml, and nothing where it has
// no entry.
func VersionOf(documentsYAML []byte, name string) string {
	lines := strings.Split(string(documentsYAML), "\n")
	start, end := entry(lines, name)
	for _, line := range lines[start:end] {
		if value, ok := strings.CutPrefix(line, "  object_version: "); ok {
			return strings.Trim(value, `'"`)
		}
	}
	return ""
}

// Next answers the version of sbom.md after a run: x raised by one, y and z as they are, and 1.1.1
// where it has none.
//
// @implements SR_40@3.3.1
func Next(version string) (string, error) {
	if version == "" {
		return "1.1.1", nil
	}
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("version %q is not x.y.z", version)
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("version %q is not x.y.z", version)
	}
	return strconv.Itoa(x+1) + "." + parts[1] + "." + parts[2], nil
}

// Hash answers what git hash-object answers for a file of this content, 7 characters.
//
// @implements SR_59@2.1.1
func Hash(content []byte) string {
	h := sha1.New() // #nosec G401 -- see the import
	h.Write([]byte("blob " + strconv.Itoa(len(content)) + "\x00"))
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))[:7]
}

// Entry answers the documents.yaml with the entry of sbom.md set to version and hash, completed and
// generated (PROC_13), and the entry added at its end where there is none. The other lines stay as
// they are.
//
// @implements SR_59@2.1.1
func Entry(documentsYAML []byte, version, hash string) []byte {
	return EntryOf(documentsYAML, Name, version, hash)
}

// EntryOf answers the documents.yaml with the entry of the document name set as Entry sets the one
// of sbom.md.
//
// @implements SR_98@2.1.1
func EntryOf(documentsYAML []byte, name, version, hash string) []byte {
	block := []string{
		"- name: " + name,
		"  type: document",
		"  object_version: '" + version + "'",
		"  object_hash: '" + hash + "'",
		"  status: completed",
		"  verification_method: [generated]",
	}
	text := strings.TrimRight(string(documentsYAML), "\n")
	lines := strings.Split(text, "\n")
	start, end := entry(lines, name)
	if start == end {
		return []byte(strings.Join(append(lines, block...), "\n") + "\n")
	}
	set := map[string]string{"object_version": block[2], "object_hash": block[3], "status": block[4], "verification_method": block[5]}
	var changed []string
	for _, line := range lines[start:end] {
		key, _, _ := strings.Cut(strings.TrimSpace(line), ":")
		if replacement, ok := set[key]; ok {
			line = replacement
			delete(set, key)
		}
		changed = append(changed, line)
	}
	for _, line := range block[2:] { // a field the entry lacks, in the order of the block
		key, _, _ := strings.Cut(strings.TrimSpace(line), ":")
		if _, ok := set[key]; ok {
			changed = append(changed, line)
		}
	}
	out := append(append(append([]string{}, lines[:start]...), changed...), lines[end:]...)
	return []byte(strings.Join(out, "\n") + "\n")
}

// entry answers the lines of the entry of the document name, from its line "- name:" to the line
// before the next entry; start equals end where there is none.
func entry(lines []string, name string) (start, end int) {
	for i, line := range lines {
		if strings.TrimSpace(line) != "- name: "+name {
			continue
		}
		end = i + 1
		for end < len(lines) && strings.HasPrefix(lines[end], "  ") {
			end++
		}
		return i, end
	}
	return len(lines), len(lines)
}
