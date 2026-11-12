// Package export makes the documents of a commit or tag into one package of PDFs (TOOL_01): one PDF
// per document that documents.yaml names, an index of them, in one ZIP.
package export

import (
	"archive/zip"
	"fmt"
	"html"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Document is a document of the package: its path from the root of the repo, its title, and its
// version, revision and status at the commit or tag.
type Document struct {
	Path     string
	Title    string
	Version  string
	Revision string
	Status   string
}

// File is a file of the ZIP: its name in it and its content.
type File struct {
	Name string
	Data []byte
}

// Documents answers every document documents.yaml names, in its order, with its version and status.
//
// @implements SR_68@1.1.1
func Documents(documentsYAML []byte) ([]Document, error) {
	var file struct {
		Documents []struct {
			Name          string `yaml:"name"`
			ObjectVersion string `yaml:"object_version"`
			Status        string `yaml:"status"`
		} `yaml:"documents"`
	}
	if err := yaml.Unmarshal(documentsYAML, &file); err != nil {
		return nil, fmt.Errorf("documents.yaml: %w", err)
	}
	var documents []Document
	for _, d := range file.Documents {
		documents = append(documents, Document{Path: d.Name, Version: d.ObjectVersion, Status: d.Status})
	}
	return documents, nil
}

// PDFName answers the name of the PDF of a document: its path with `.pdf` for `.md`.
//
// @implements SR_68@1.1.1
func PDFName(path string) string {
	return strings.TrimSuffix(path, ".md") + ".pdf"
}

// Title answers the title of a document: its first heading of level 1, and its path where it has none.
func Title(markdown []byte, path string) string {
	for _, line := range strings.Split(string(markdown), "\n") {
		if title, ok := strings.CutPrefix(line, "# "); ok {
			return strings.TrimSpace(title)
		}
	}
	return path
}

// Footer answers what the footer of a document shows of it, such as
// "devops.md · 2.1.1 · a1b2c3d · approved".
//
// @implements SR_71@1.1.1
func Footer(d Document) string {
	return strings.Join([]string{d.Path, d.Version, d.Revision, d.Status}, " · ")
}

// HeaderTemplate answers the header of every page: the title of the document on the left, the name
// of the repo on the right.
//
// @implements SR_70@1.1.1
func HeaderTemplate(title, repo string) string {
	return `<div style="` + band + `"><span>` + html.EscapeString(title) + `</span><span>` +
		html.EscapeString(repo) + `</span></div>`
}

// FooterTemplate answers the footer of every page: the date of the commit or tag on the left, the
// page of the pages in the middle, and what it shows of the document on the right.
//
// @implements SR_70@1.1.1
// @implements SR_71@1.1.1
func FooterTemplate(date, about string) string {
	return `<div style="` + band + `"><span>` + html.EscapeString(date) +
		`</span><span>page <span class="pageNumber"></span> of <span class="totalPages"></span></span><span>` +
		html.EscapeString(about) + `</span></div>`
}

// band is the style of a header and a footer: Chrome prints them without the style of the page.
const band = `font-family: -apple-system, 'Segoe UI', Helvetica, Arial, sans-serif; font-size: 8px; ` +
	`width: 100%; margin: 0 1.5cm; display: flex; justify-content: space-between; gap: 1em; color: #555;`

// Index answers the page of the index: every document with its title, version, revision and status.
//
// @implements SR_72@1.1.1
func Index(documents []Document, repo, ref string) string {
	var b strings.Builder
	b.WriteString("<h1>Documents of " + html.EscapeString(repo) + " at " + html.EscapeString(ref) + "</h1>\n")
	b.WriteString("<table>\n<thead><tr><th>Document</th><th>Title</th><th>Version</th><th>Revision</th><th>Status</th></tr></thead>\n<tbody>\n")
	for _, d := range documents {
		b.WriteString("<tr>")
		for _, cell := range []string{d.Path, d.Title, d.Version, d.Revision, d.Status} {
			b.WriteString("<td>" + html.EscapeString(cell) + "</td>")
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("</tbody>\n</table>\n")
	return b.String()
}

// ZipName answers the name of the ZIP of a repo at a commit or tag.
//
// @implements SR_73@1.1.1
func ZipName(repo, ref string) string {
	return repo + "-" + ref + "-documents.zip"
}

// Zip writes the files into one ZIP, in their order.
//
// @implements SR_73@1.1.1
func Zip(w io.Writer, files []File) error {
	z := zip.NewWriter(w)
	for _, f := range files {
		entry, err := z.Create(f.Name)
		if err != nil {
			return err
		}
		if _, err := entry.Write(f.Data); err != nil {
			return err
		}
	}
	return z.Close()
}
