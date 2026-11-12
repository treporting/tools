package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Git reads a repo at a commit or tag, without checking it out.
type Git struct {
	Dir string
	Ref string
}

// Commit answers the commit of the ref, and an error where there is none.
//
// @implements SR_74@1.1.1
func (g Git) Commit() (string, error) {
	out, err := git(g.Dir, "rev-parse", "--verify", "--quiet", g.Ref+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("no commit or tag %q in %s", g.Ref, g.Dir)
	}
	return strings.TrimSpace(out), nil
}

// Show answers the content of a file at the ref.
//
// @implements SR_68@1.1.1
func (g Git) Show(path string) ([]byte, error) {
	out, err := git(g.Dir, "show", g.Ref+":"+path)
	return []byte(out), err
}

// Revision answers the last commit up to the ref that changed the file, 7 characters.
//
// @implements SR_71@1.1.1
func (g Git) Revision(path string) (string, error) {
	out, err := git(g.Dir, "log", "-1", "--format=%h", "--abbrev=7", g.Ref, "--", path)
	return strings.TrimSpace(out), err
}

// Date answers the day of the commit of the ref, such as 2026-11-01.
//
// @implements SR_70@1.1.1
func (g Git) Date() (string, error) {
	out, err := git(g.Dir, "log", "-1", "--format=%cs", g.Ref)
	return strings.TrimSpace(out), err
}

// Printer prints HTML pages to PDF in one headless Chrome.
type Printer struct {
	ctx    context.Context
	cancel func()
	dir    string
}

// NewPrinter starts headless Chrome.
func NewPrinter() (*Printer, error) {
	dir, err := os.MkdirTemp("", "pdf")
	if err != nil {
		return nil, err
	}
	//	60 seconds for Chrome to start, not the 20 of chromedp - a busy
	//	runner can take longer
	options := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WSURLReadTimeout(60*time.Second))
	alloc, cancelAlloc := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancelCtx := chromedp.NewContext(alloc)
	if err := chromedp.Run(ctx); err != nil {
		cancelCtx()
		cancelAlloc()
		return nil, fmt.Errorf("headless Chrome: %w", err)
	}
	return &Printer{ctx: ctx, cancel: func() { cancelCtx(); cancelAlloc(); _ = os.RemoveAll(dir) }, dir: dir}, nil
}

// Close ends Chrome.
func (p *Printer) Close() { p.cancel() }

// A4 in inches, and 2 cm at top and bottom, 1.5 cm at the sides.
const (
	paperWidth  = 8.27
	paperHeight = 11.69
	topBottom   = 2 / 2.54
	sides       = 1.5 / 2.54
)

// Print answers the PDF of an HTML page with the header and footer given, once the page is done -
// its Mermaid diagrams drawn.
//
// @implements SR_70@1.1.1
func (p *Printer) Print(htmlPage, header, footer string) ([]byte, error) {
	file := filepath.Join(p.dir, "page.html")
	if err := os.WriteFile(file, []byte(htmlPage), 0o600); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(p.ctx, 2*time.Minute)
	defer cancel()
	var pdf []byte
	var failed string
	err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+file),
		chromedp.Poll(`window.pdfReady === true || typeof window.pdfError === 'string'`, nil),
		chromedp.Evaluate(`window.pdfError || ''`, &failed),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if failed != "" {
				return errors.New(failed)
			}
			var err error
			pdf, _, err = page.PrintToPDF().
				WithPaperWidth(paperWidth).WithPaperHeight(paperHeight).
				WithMarginTop(topBottom).WithMarginBottom(topBottom).
				WithMarginLeft(sides).WithMarginRight(sides).
				WithDisplayHeaderFooter(true).WithHeaderTemplate(header).WithFooterTemplate(footer).
				WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	return pdf, err
}

// git runs git in dir and answers its output.
func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, exit.Stderr)
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
}
