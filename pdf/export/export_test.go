package export

import (
	"archive/zip"
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

const documentsYAML = `schema: schemas/documents.schema.json
type: documents
documents:
- name: README.md
  type: document
  object_version: '9.1.1'
  object_hash: 'f2685eb'
  status: approved
  verification_method: [review]
- name: requirements/tool-requirements.md
  type: explanation
  level: tool
  object_version: '16.1.1'
  object_hash: '1234567'
  status: review
  verification_method: [review]
`

// @implements TC_97@2.1.1, TC_98@1.1.1
func TestDocumentsOfDocumentsYAML(t *testing.T) {
	got, err := Documents([]byte(documentsYAML))
	if err != nil {
		t.Fatal(err)
	}
	want := []Document{
		{Path: "README.md", Version: "9.1.1", Status: "approved"},
		{Path: "requirements/tool-requirements.md", Version: "16.1.1", Status: "review"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

// @implements TC_97@2.1.1
func TestPDFNameIsThePathWithPDF(t *testing.T) {
	for path, want := range map[string]string{"README.md": "README.pdf", "requirements/tool-requirements.md": "requirements/tool-requirements.pdf"} {
		if got := PDFName(path); got != want {
			t.Errorf("PDFName(%q) = %q", path, got)
		}
	}
}

func TestTitleIsTheFirstHeading(t *testing.T) {
	if got := Title([]byte("text\n# Tool requirements\n\n## Launches\n"), "x.md"); got != "Tool requirements" {
		t.Errorf("got %q", got)
	}
	if got := Title([]byte("## Only a second level\n"), "x.md"); got != "x.md" {
		t.Errorf("no heading: got %q", got)
	}
}

// @implements TC_99@1.1.1
func TestFooterNamesPathVersionRevisionStatus(t *testing.T) {
	got := Footer(Document{Path: "devops.md", Version: "2.1.1", Revision: "a1b2c3d", Status: "approved"})
	if got != "devops.md · 2.1.1 · a1b2c3d · approved" {
		t.Errorf("got %q", got)
	}
}

// @verifies SR_71@1.1.1
func TestFooterTemplateHoldsDatePagesAndDocument(t *testing.T) {
	got := FooterTemplate("2026-11-01", "devops.md · 2.1.1 · a1b2c3d · approved")
	for _, want := range []string{"2026-11-01", `class="pageNumber"`, `class="totalPages"`, "devops.md · 2.1.1 · a1b2c3d · approved"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %q in %s", want, got)
		}
	}
}

func TestHeaderTemplateEscapes(t *testing.T) {
	got := HeaderTemplate("A <b> & C", "specs-and-deploy")
	if !strings.Contains(got, "A &lt;b&gt; &amp; C") || !strings.Contains(got, "specs-and-deploy") {
		t.Errorf("got %s", got)
	}
}

// @implements TC_100@1.1.1
func TestIndexListsEveryDocument(t *testing.T) {
	got := Index([]Document{
		{Path: "README.md", Title: "specs-and-deploy", Version: "9.1.1", Revision: "ab57a8b", Status: "approved"},
		{Path: "sbom.md", Title: "SBOM", Version: "1.1.1", Revision: "15c222f", Status: "approved"},
	}, "specs-and-deploy", "T044")
	for _, want := range []string{
		"<td>README.md</td><td>specs-and-deploy</td><td>9.1.1</td><td>ab57a8b</td><td>approved</td>",
		"<td>sbom.md</td><td>SBOM</td><td>1.1.1</td><td>15c222f</td><td>approved</td>",
		"specs-and-deploy at T044",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %q in\n%s", want, got)
		}
	}
}

// @implements TC_101@1.1.1
func TestZipNameNamesRepoAndRef(t *testing.T) {
	if got := ZipName("specs-and-deploy", "T044"); got != "specs-and-deploy-T044-documents.zip" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_101@1.1.1
func TestZipHoldsEveryFile(t *testing.T) {
	var buf bytes.Buffer
	files := []File{{Name: "README.pdf", Data: []byte("%PDF-1")}, {Name: "requirements/tool-requirements.pdf", Data: []byte("%PDF-2")}, {Name: "index.pdf", Data: []byte("%PDF-3")}}
	if err := Zip(&buf, files); err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(z.File) != len(files) {
		t.Fatalf("%d files", len(z.File))
	}
	for i, f := range z.File {
		r, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(r)
		if f.Name != files[i].Name || string(data) != string(files[i].Data) {
			t.Errorf("file %d: %s %q", i, f.Name, data)
		}
	}
}

func TestPageLoadsMermaidOnlyWhereThereIsADiagram(t *testing.T) {
	plain, err := Page("T", []byte("# T\n\n| a | b |\n|---|---|\n| 1 | 2 |\n"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(plain, "mermaid.min.js") || !strings.Contains(plain, "<table>") || !strings.Contains(plain, "window.pdfReady = true") {
		t.Errorf("plain page:\n%s", plain)
	}
	diagram, err := Page("T", []byte("# T\n\n```mermaid\nflowchart LR\n  a --> b\n```\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(diagram, `integrity="sha384-`) || !strings.Contains(diagram, `class="language-mermaid"`) {
		t.Errorf("page with a diagram:\n%s", diagram)
	}
}

// @implements TC_157@1.1.1
func TestAddressLinkedAsGitHubDoes(t *testing.T) {
	page, err := Page("t", []byte("| Requirement |\n|---|\n| SR_01@1.1.3 |\n\nWrite to a@1.com.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(page, `<a href="mailto:a@1.com">a@1.com</a>`) {
		t.Errorf("a@1.com is no link to mailto:\n%s", page)
	}
	if !strings.Contains(page, "<td>SR_01@1.1.3</td>") {
		t.Errorf("SR_01@1.1.3 is not text alone\n%s", page)
	}
}
