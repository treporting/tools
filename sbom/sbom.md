# Software Bill of Materials (SBOM)

The packages tools delivers, its own and third-party, each with its version, its
licence and its PURL. It is written by the tool sbom (TOOL_21) and not edited by hand.

## Original data

**Table original data** - every file sbom.md is written from

| File | Made by | Stored At |
|---|---|---|
| golist-root.txt | `go list -deps -f '{{with .Module}}{{if not .Main}}{{with .Replace}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{else}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{end}}{{end}}{{end}}' $(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./...)` | - |

## Made with

**Table made with** - the baselines sbom.md is made with

| Baseline of | Tag |
|---|---|
| the repo above | T0471 |

## Packages

**Table packages** - every package tools delivers

| PURL | Name | Version | Licence |
|---|---|---|---|
| pkg:github/treporting/tools@T048.0 | tools | T048.0 | Apache-2.0 |
| pkg:golang/github.com/chromedp/cdproto@v0.0.0-20250724212937-08a3db8b4327 | github.com/chromedp/cdproto | v0.0.0-20250724212937-08a3db8b4327 | MIT |
| pkg:golang/github.com/chromedp/chromedp@v0.14.2 | github.com/chromedp/chromedp | v0.14.2 | MIT |
| pkg:golang/github.com/chromedp/sysutil@v1.1.0 | github.com/chromedp/sysutil | v1.1.0 | MIT |
| pkg:golang/github.com/dop251/goja@v0.0.0-20261009203008-ac6c29549738 | github.com/dop251/goja | v0.0.0-20261009203008-ac6c29549738 | MIT |
| pkg:golang/github.com/evanw/esbuild@v0.28.2 | github.com/evanw/esbuild | v0.28.2 | MIT |
| pkg:golang/github.com/go-json-experiment/json@v0.0.0-20260213210345-44df1a37e875 | github.com/go-json-experiment/json | v0.0.0-20260213210345-44df1a37e875 | BSD-3-Clause |
| pkg:golang/github.com/go-sourcemap/sourcemap@v2.1.4+incompatible | github.com/go-sourcemap/sourcemap | v2.1.4+incompatible | BSD-2-Clause |
| pkg:golang/github.com/gobwas/httphead@v0.1.0 | github.com/gobwas/httphead | v0.1.0 | MIT |
| pkg:golang/github.com/gobwas/pool@v0.2.1 | github.com/gobwas/pool | v0.2.1 | MIT |
| pkg:golang/github.com/gobwas/ws@v1.4.0 | github.com/gobwas/ws | v1.4.0 | MIT |
| pkg:golang/github.com/santhosh-tekuri/jsonschema/v6@v6.0.3 | github.com/santhosh-tekuri/jsonschema/v6 | v6.0.3 | Apache-2.0 |
| pkg:golang/github.com/yuin/goldmark@v1.7.8 | github.com/yuin/goldmark | v1.7.8 | MIT |
| pkg:golang/golang.org/x/sys@v0.34.0 | golang.org/x/sys | v0.34.0 | BSD-3-Clause |
| pkg:golang/golang.org/x/text@v0.14.0 | golang.org/x/text | v0.14.0 | BSD-3-Clause |
| pkg:golang/gopkg.in/yaml.v3@v3.0.1 | gopkg.in/yaml.v3 | v3.0.1 | Apache-2.0 AND MIT |
| pkg:golang/stdlib@go1.25.0 | stdlib | go1.25.0 | BSD-3-Clause |

## Open

`<empty>`

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| 1.1.1 | 2026-11-15 | T046.0 | #15 | sbom 1a45bbb |
| 2.1.1 | 2026-11-15 | T046.1 | #21 | sbom dc56473 |
| 3.1.1 | 2026-11-16 | T047.0 | #26 | sbom ee21d8c |
| 4.1.1 | 2026-11-18 | T0471.0 | #32 | sbom 961bae7 |
| 5.1.1 | 2026-11-24 | T048.0 | #65 | sbom 6162630 |
