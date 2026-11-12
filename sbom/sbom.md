# Software Bill of Materials (SBOM)

The packages specs-and-deploy delivers, its own and third-party, each with its version, its
licence and its PURL. It is written by the tool sbom (TOOL_21) and not edited by hand.

## Original data

**Table original data** - every file sbom.md is written from

| File | Made by | Stored At |
|---|---|---|
| golist-tools.txt | `cd tools && go list -deps -f '{{with .Module}}{{if not .Main}}{{with .Replace}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{else}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{end}}{{end}}{{end}}' $(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./...)` | - |

## Packages

**Table packages** - every package specs-and-deploy delivers

| PURL | Name | Version | Licence |
|---|---|---|---|
| pkg:github/treporting/specs-and-deploy@T045 | specs-and-deploy | T045 | Apache-2.0 |
| pkg:golang/github.com/chromedp/cdproto@v0.0.0-20250724212937-08a3db8b4327 | github.com/chromedp/cdproto | v0.0.0-20250724212937-08a3db8b4327 | MIT |
| pkg:golang/github.com/chromedp/chromedp@v0.14.2 | github.com/chromedp/chromedp | v0.14.2 | MIT |
| pkg:golang/github.com/chromedp/sysutil@v1.1.0 | github.com/chromedp/sysutil | v1.1.0 | MIT |
| pkg:golang/github.com/go-json-experiment/json@v0.0.0-20260213210345-44df1a37e875 | github.com/go-json-experiment/json | v0.0.0-20260213210345-44df1a37e875 | BSD-3-Clause |
| pkg:golang/github.com/gobwas/httphead@v0.1.0 | github.com/gobwas/httphead | v0.1.0 | MIT |
| pkg:golang/github.com/gobwas/pool@v0.2.1 | github.com/gobwas/pool | v0.2.1 | MIT |
| pkg:golang/github.com/gobwas/ws@v1.4.0 | github.com/gobwas/ws | v1.4.0 | MIT |
| pkg:golang/github.com/santhosh-tekuri/jsonschema/v6@v6.0.3 | github.com/santhosh-tekuri/jsonschema/v6 | v6.0.3 | Apache-2.0 |
| pkg:golang/github.com/yuin/goldmark@v1.7.8 | github.com/yuin/goldmark | v1.7.8 | MIT |
| pkg:golang/golang.org/x/text@v0.14.0 | golang.org/x/text | v0.14.0 | BSD-3-Clause |
| pkg:golang/gopkg.in/yaml.v3@v3.0.1 | gopkg.in/yaml.v3 | v3.0.1 | Apache-2.0 AND MIT |
| pkg:golang/stdlib@go1.25.14 | stdlib | go1.25.14 | BSD-3-Clause |

## Open

`<empty>`

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| 1.1.1 | 2026-11-01 | T044 | #377 | sbom 8788514 |
| 2.1.1 | 2026-11-02 | T0441 | #418 | sbom eef8e1e |
| 3.1.1 | 2026-11-05 | T045 | #556 | sbom c811e2f |
