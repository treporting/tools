// Command pdf makes the documents of a commit or tag into one ZIP of PDFs (TOOL_01, TOOL_02,
// TOOL_30, TOOL_31).
//
//	pdf -ref commit-or-tag [-root dir] [-repo name] [-out dir]
//
// It reads the repo at the ref, without checking it out: documents.yaml, and the Markdown of every
// document it names. It writes one PDF per document, named after its path, and index.pdf, which lists
// every document with its version, revision and status, into <repo>-<ref>-documents.zip in out. It
// needs Google Chrome or Chromium, and the network where a document holds a Mermaid diagram.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treporting/specs-and-deploy/tools/pdf/export"
)

func main() {
	root := flag.String("root", ".", "the root of the repo")
	ref := flag.String("ref", "", "the commit or tag, such as T044")
	repo := flag.String("repo", "", "the name of the repo; the name of the directory of root where empty")
	out := flag.String("out", ".", "the directory the ZIP is written to")
	flag.Parse()

	if *ref == "" {
		fail(fmt.Errorf("-ref is missing"))
	}
	if *repo == "" {
		abs, err := filepath.Abs(*root)
		if err != nil {
			fail(err)
		}
		*repo = filepath.Base(abs)
	}

	files, err := build(export.Git{Dir: *root, Ref: *ref}, *repo)
	if err != nil {
		fail(err)
	}
	var zip bytes.Buffer
	if err := export.Zip(&zip, files); err != nil {
		fail(err)
	}
	path := filepath.Join(*out, export.ZipName(*repo, *ref))
	if err := os.WriteFile(path, zip.Bytes(), 0o600); err != nil {
		fail(err)
	}
	fmt.Printf("%s written: %d documents and index.pdf\n", path, len(files)-1)
}

// build answers the PDF of every document documents.yaml names at the ref, and index.pdf last.
//
// @implements SR_68@1.1.1
// @implements SR_72@1.1.1
func build(git export.Git, repo string) ([]export.File, error) {
	if _, err := git.Commit(); err != nil {
		return nil, err
	}
	date, err := git.Date()
	if err != nil {
		return nil, err
	}
	yaml, err := git.Show("documents.yaml")
	if err != nil {
		return nil, err
	}
	documents, err := export.Documents(yaml)
	if err != nil {
		return nil, err
	}

	printer, err := export.NewPrinter()
	if err != nil {
		return nil, err
	}
	defer printer.Close()

	var files []export.File
	for i, d := range documents {
		source, err := git.Show(d.Path)
		if err != nil {
			return nil, err
		}
		if documents[i].Revision, err = git.Revision(d.Path); err != nil {
			return nil, err
		}
		documents[i].Title = export.Title(source, d.Path)
		page, err := export.Page(documents[i].Title, source)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", d.Path, err)
		}
		pdf, err := printer.Print(page,
			export.HeaderTemplate(documents[i].Title, repo), export.FooterTemplate(date, export.Footer(documents[i])))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", d.Path, err)
		}
		files = append(files, export.File{Name: export.PDFName(d.Path), Data: pdf})
	}

	title := "Documents of " + repo + " at " + git.Ref
	index, err := printer.Print(export.Wrap(title, export.Index(documents, repo, git.Ref)+`<script>window.pdfReady = true;</script>`),
		export.HeaderTemplate(title, repo), export.FooterTemplate(date, "index · "+git.Ref))
	if err != nil {
		return nil, fmt.Errorf("index.pdf: %w", err)
	}
	return append(files, export.File{Name: "index.pdf", Data: index}), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
