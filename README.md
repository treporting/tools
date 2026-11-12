# tools

The tools of our work: what checks and supports the requirements, reviews and documents of
treporting. Tools are code (see [devops.md](../devops.md)): every tool realizes an approved tool
requirement of [tool-requirements.yaml](../requirements/tool-requirements.yaml). It is written for
whoever builds, runs or reviews the tools.

## Where the tools live

**Decision**: in `tools/` of this repo, specs-and-deploy.

**Why**: the tools check above all the files of this repo - requirements, reviews, schemas - so what
is checked and what checks lie together. The repo is meant for more than documents already: the
containers and the system tests come here as well. The checks of the code repos run as a reusable
workflow in `treporting/.github`, which takes the tool from here (see Open).

**Rejected**:

- a repo of its own - a clean cut, but the CI of every repo would first have to fetch the tool, and
  it is one repo more
- reporter-annotations - it is meant as a tool beside the product, but a check of reviews has
  nothing to do with annotations

## Language

**Decision**: Go.

**Why**: the same language as the subsystems, and its toolchain stands in devops.md already - no
new language in the chain of our tools, which the threat of an outdated toolchain asks to keep
small. A tool is one program
with no dependency at run time, which the CI runs as it is. It is tested like the code of the
product - a checker with a fault would let faults pass unseen.

**Rejected**:

- Python - quickly written, but a new tool in the chain, with packages of its own, not in
  devops.md, and a setup of its own in the CI
- shell - no dependency, but too error-prone for reading YAML and schemas

## The chain of a tool

**Decision**: a tool requirement leads to software requirements, and those to code and tests - as
for a subsystem. The software requirements of the tools are in
[software-requirements.yaml](software-requirements.yaml), each naming its tool requirement as
uplink. The code names them in `@implements SR_nn`, a test in `@verifies SR_nn`. The tests are on
the levels MT and SwIT.

**Why**: the first tool showed it. TOOL_03 holds five behaviours that can each be tested on their
own; by ISO29148 each is a requirement of its own (singular). Its review noted it, and
each test of reviewcheck verified one of the five while all named only TOOL_03 - which of them a
test verifies could not be told. A tool requirement says what a tool does for our work; a
software requirement says how it behaves.

**Rejected**: a tool requirement leading straight to code and tests - the first decision here. It
held a tool requirement as concrete as a software requirement, which TOOL_03 proved wrong.

## Who changed a requirement

**Decision**: `modified_by` holds the original as git records it: the author of the last commit,
not a merge, that changed the requirement, and its co-authors - `Claude Opus 5.5`, or `Mathias Fischer`.
Nothing is derived from it and kept beside it; there is no `mmodified_by`.

**Why**: git records author and co-authors by itself; a record made by a computer is the evidence an
audit trail asks for (REG_CFR21-11_01), and the original stays original (ARCR_01, REG_FDADI_01).

**Rejected**:

- only the one who wrote it, `claude-code` - a value already interpreted would stand in the field of
  the original, and the author of the commit would be lost
- a field `created_by` or `modified_by` written by hand in the YAML - visible in the file, but a
  claim rather than a record, it can drift from git unseen, and it would be a second source of the
  same truth

## Who the author of a requirement or a document is

**Decision**: the authors of a requirement or a document are the names in its `modified_by` - the
author and the co-authors of the last commit, not a merge, that changed it, whatever names stand
there. The reviewer is none of them (PROC_15; TOOL_03, TOOL_11). Who writes the text is the author of the commit:
Claude Code commits with `--author="Claude Opus 5.5 <noreply@anthropic.com>"`; the Product Owner
stays its committer.

```bash
git commit --author="Claude Opus 5.5 <noreply@anthropic.com>" -m "..."
```

**Why**: git tells author and committer apart by itself. With the writer as author, the check needs
no rule about names: it holds for any assistant or person, and needs no change when the assistant
changes. A commit the Product Owner writes alone names them, and they may not review it - as it
should be.

A merge commit is no authorship: the one who merges approves, and writes nothing. Where two pull
requests change one file from an older `main`, the merge commit equals neither of its parents, and
git counts it as a change of the file - with the one who merged as its author. The merge commit
stays as the record of the approval; it is only not read as authorship. The known limit: a conflict
resolved by hand in a merge is a change no author is found for.

**Rejected**:

- the co-author whose name begins with `Claude`, else the author of the commit - the first decision;
  it names a product, so another assistant, such as Copilot, would need a change of the rule
- the authors without the committer, where none is left the committer - it keeps the history valid,
  but it is a rule about git again, not the record of git
- the check left to the approver by hand - TOOL_03 asks for it
- rebase merges, so that `main` has no merge commits - git would no longer record who merged, the
  evidence of the approval

*Note*: Every review result made before this decision reviewed commits with the Product Owner as
author.

## The component review

**Decision**: everything of the component `review` lies below `tools/review`:

```
tools/review/                the component: what its tools share - git, and finding an id in its file
  check/                     the module review/check: reviewcheck
  skeleton/                  the module review/skeleton: reviewskeleton
  approval/                  the module review/approval: approvalcheck
  cmd/reviewcheck/           the command
  cmd/reviewskeleton/        the command
  cmd/approvalcheck/         the command
```

**Why**: the hierarchy of the component can be read from the folders. `check`, `skeleton` and
`approval` depend on `review`, not on each other. A command keeps the name of its tool: Go names a
program by the last part of its path. An SR names the component, `review`; the module is set by the code.

**Rejected**:

- the commands in `tools/cmd` and the packages beside them - the component would be spread over
  the tree
- a folder `lib/` - in Go every package without `main` is a library already
- commands named `check` and `skeleton` - too general a name for a program
- the module as the component of an SR - which package realizes a behaviour is a decision of the
  design, on the levels r6 and r7

**Decision**: a tool finds the files of requirements by their `type`. The component `review` reads
the head of every YAML git knows in the working tree, and keeps those whose `type` holds
requirements: `requirements`, `normative`, `architecture`, `hazards` and `rules` (formats.md, table
types). From the `object_key` in them it knows the file of every id. In an old revision a tool reads
an id in the file the working tree names for it.

**Why**: a file says itself what it holds (formats.md, section The head of a file). A new level, a
new file or a changed prefix needs no change of code.

**Rejected**:

- the id tells its file, by a fixed list of prefixes and the rule `REG_<module>_nn` →
  `regulatory/<module>.yaml`, as before - every new level or prefix is a change of code, and the
  list stood a second time beside the table files of the levels in formats.md
- a file `levels.yaml` - see formats.md, section Documents

## reviewcheck

Checks every review result in `reviews/results/` against the rules of formats.md (TOOL_03, TOOL_11; SR_01
to SR_08, SR_16, SR_84 and SR_103 in [software-requirements.yaml](software-requirements.yaml)): every criterion of its
criteria file once, the verdict passed only when every criterion passed, a note to every failed
criterion, no value but passed or failed, the revision reviewed exists, the criteria revision changed
the criteria file, the reviewer is none of the authors of the requirement or document, and a
completed review has a `completed_at` not before its `created_at`. It
reports the findings of a result by criterion, in the order of its criteria file, the same in every run.

```bash
cd tools
go run ./review/cmd/reviewcheck -root ..        # every result in reviews/results/
go run ./review/cmd/reviewcheck -root .. -base origin/main   # those changed against main
go test ./...                            # MT
go test -tags integration ./...          # SwIT, on a history of git made in the test
```

With `-base`, as the CI runs it, it checks only the results a pull request adds or changes: a result
names fixed commits, so its check does not change once it passed. It exits with 1 where a result
breaks a rule. It depends on `gopkg.in/yaml.v3` (MIT and Apache-2.0).

## reviewskeleton

Writes the skeleton of a review result into `reviews/results/` (TOOL_04; SR_11 to SR_15 in
[software-requirements.yaml](software-requirements.yaml)): the name from the time and the scope,
the criteria file of the level and the last commit that changed it, the commit `HEAD` as the
revision of each, and every criterion and the verdict as `todo`. It writes nothing where a target is
not found, or where the targets belong to more than one level.

```bash
cd tools
go run ./review/cmd/reviewskeleton -root .. -reviewer "Mathias Fischer" -scope harness harness.md
go run ./review/cmd/reviewskeleton -root .. -reviewer "Mathias Fischer" -scope sr-11-15 SR_11 SR_12 SR_13
```

In its head the skeleton names its schema, the type `review` and the level of its criteria file
(SR_17), and is a record still being made: `object_version` 1.1.1 and `status` open (SR_83). The
decision of the reviewer sets `completed` and `completed_at`. It prints the path of the file it
wrote. The code is the package `review/skeleton`, in the component `review`.

## approvalcheck

Checks every approval from the approved objects (TOOL_32, TOOL_33; SR_75 to SR_80 and SR_85 in
[software-requirements.yaml](software-requirements.yaml)). It reports:

- an approved object that needs a review, and that no review result names
- one whose last review result, by its `created_at`, has the verdict failed
- one whose last review result is still open: a proposal, not a decision
- one whose content x counts in the revision of that result is not its content now
- an approved requirement whose attributes for x are not those of the first commit in which it was
  approved with its x
- an approved or completed document whose `object_hash` is not the hash of its file

An object needs a review where its kind has a criteria file in `reviews/`, as reviewskeleton finds
it; a document where it has the verification method `review` as well. What x counts is read from
the schema of its file, `attributes_for_version_x`; of a document it is the hash of its file.

```bash
cd tools
go run ./review/cmd/approvalcheck -root ..              # every approved object
go run ./review/cmd/approvalcheck -root .. -base main   # the objects in the files changed against main
```

With `-base`, it checks the objects in the files changed against that commit - a requirement in its
file, a document in its file or its entry in `documents.yaml`. It runs so on every pull request
(SR_81), and without `-base` in the baseline (SR_82). It reads the history of git, so a checkout
needs all of it. It exits with 1 where an approval does not hold. The code is the package
`review/approval`, in the component `review`.

## schemacheck

Checks every YAML file of our own data under version control against the schema it names in its
field `schema` (TOOL_05; SR_18 and SR_19 in [software-requirements.yaml](software-requirements.yaml)),
and reports each place where one does not match, with the file and the path within it. A YAML in
`.github/` has a form GitHub owns, and is not checked. The CI runs it on every pull request (SR_20).

```bash
cd tools
go run ./schema/cmd/schemacheck -root ..
```

With `-above`, the root of the repo above, a schema this repo does not hold is taken from there
(SR_107). Without it, such a schema is reported as not there.

A finding reads `requirements/tool-requirements.yaml /requirements/6/status: value must be one of
'draft', 'review', 'approved', 'obsolete'`. The code is the component `schema`: the package
`schema/check` and the command `schema/cmd/schemacheck`. Its dependencies are in devops.md.

## rulesetcheck

Compares the ruleset of main on GitHub with `.github/rulesets/main.json` (TOOL_13; SR_21 and SR_22 in
[software-requirements.yaml](software-requirements.yaml)), and reports each difference in its
target, enforcement, conditions and rules. What GitHub sets itself is left out, and `bypass_actors`,
which cannot be read without the rights of an admin; that main has no bypass, its demonstration
shows. A ruleset that cannot be read, or no ruleset of the name in the file, is a finding, not a
pass. The CI runs it on every pull request (SR_23), with the token of the workflow, read only.

```bash
cd tools
GITHUB_TOKEN=$(gh auth token) go run ./ruleset/cmd/rulesetcheck -root ..
```

The code is the component `ruleset`: the package `ruleset/check` and the command
`ruleset/cmd/rulesetcheck`. Its tests stand in for GitHub with a server in the test process
(`net/http/httptest`), so they run without a network.

## annotationcheck

Checks that every id named in `@implements`, `@verifies` or `@refines` is in the requirements of the
base commit of the pull request, and approved there (TOOL_06; SR_24, SR_25 and SR_35 in
[software-requirements.yaml](software-requirements.yaml)), and that it names its version, as
`ID@x.y.z` (TOOL_20, SR_30). The id is read apart from its version and checked as before. A requirement approved in the same pull
request does not count: it has to be through before the code (PROC_05, PROC_12).

It reads a file by its kind, from the two lists in [kinds.yaml](annotation/kinds.yaml): the
whitelist of the kinds that are code, the blacklist of those that are no code. A file of a kind on
neither list is a finding (SR_27), until its kind is set; a YAML that names a schema is our data and
is not read. In a file it reads, an annotation stands first in a line of comment: after spaces, a
comment mark of one to three signs other than a quote, then `@implements` or `@verifies` - in every
language, with no list of languages. Why both lists: formats.md, "Where an annotation is read".
`@refines` it reads in the `$comment` of each schema in [schemas/](../schemas/) (SR_35).

**Decision**: a line of comment is told by its start alone, with no knowledge of the strings of a
language. Its known limit: a line of a string over lines - a raw string of Go, a template of
JavaScript - that begins with a comment mark and an annotation is read as one. A test writes such
content with `\n`, not over lines.

**Why**: the limit does the gate no harm. With an id that is approved, no breach gets through; with
one that is not, the false alarm is loud, and the remedy is a line in a test. It arises in practice
only in tests that need annotations as their content.

**Rejected**:

- a lexer for each kind, set in `kinds.yaml` - comment marks, strings, escapes - a large effort,
  with traps of its own (escapes that differ, a `'` that is no string, strings with no quotes), for
  an error that does the gate no harm
- a parser for each language - right where an error harms: in reporter-annotations, which counts an
  annotation as a realization; not for the gate
- `go/scanner` for Go alone - exact for Go, but a switch by language again

The CI runs it on every pull request with the base commit GitHub gives (SR_26, SR_31).

```bash
cd tools
go run ./annotation/cmd/annotationcheck -root ..                 # base: where HEAD branched from main
go run ./annotation/cmd/annotationcheck -root .. -base <commit>
```

The code is the component `annotation`: the package `annotation/check` and the command
`annotation/cmd/annotationcheck`.

## uplinkcheck

Checks that every approved requirement names approved uplinks only (TOOL_09; SR_28 and SR_29 in
[software-requirements.yaml](software-requirements.yaml)). It reads every file of requirements as it
stands - the levels, the tools, the rules - and reports each uplink in `source_ids` of an approved
requirement that no file holds, or that is not approved. The version after `@` is not compared here;
that is PROC_20.

It checks the state, not a change: on `main` a requirement stays approved until its change passes
its review, and a change whose review failed does not reach `main` (harness.md, step 10). So an
approved requirement never names a draft there, unless someone merges one.

The CI runs it on every pull request (SR_29).

```bash
cd tools
go run ./uplink/cmd/uplinkcheck -root ..
```

With `-above`, the root of the repo above, an uplink no file of this repo holds is looked up there,
and checked there to be approved (SR_108).

The code is the component `uplink`: the package `uplink/check` and the command
`uplink/cmd/uplinkcheck`. It finds an id in its file through the component `review`.

## linkcheck

Makes the link check (TOOL_16; SR_32 to SR_34 and SR_87 in [software-requirements.yaml](software-requirements.yaml)),
as definitions.md sets it out in the section Link check and impact assessment. It compares the
version each entry names with the version its object has, and reports each whose object has
another x (traceability.md, section When a link is suspect):

- the links in `source_ids`, the entries of `allocated` and of `controls` in every file of
  objects (SR_32)
- the targets in `target_ids` of every test case (SR_87)
- the ids of `@implements`, `@verifies` and `@refines`, found as annotationcheck finds them (SR_33)

It reads the files as they stand. An entry whose object no file holds, or that names no version, is
not compared: uplinkcheck, schemacheck and annotationcheck report it. It reports and does not
decide whether the child still fits; that is the impact assessment.

The CI runs it on every pull request in a job of its own (SR_34). The job does not fail and is no
required check: a suspect link blocks the baseline, not the pull request (PROC_20).

```bash
cd tools
go run ./link/cmd/linkcheck -root ..
```

With `-above`, the root of the repo above, an entry or an id whose object this repo does not hold
is compared with its object there (SR_109).

The code is the component `link`: the package `link/check` and the command `link/cmd/linkcheck`. It
finds an id in its file through the component `review`, and the annotations through the component
`annotation`.

## sbom

Creates the SBOM of a subsystem as `sbom/sbom.md` of its repo (TOOL_21; SR_36 to SR_46, SR_59, SR_66, SR_101 and SR_102 in
[software-requirements.yaml](software-requirements.yaml)). It runs in the repo checked out at the
baseline and lists every package the subsystem delivers, each once, sorted by PURL:

- the package of the subsystem itself, with the label of the baseline as its version (SR_46)
- where a directory of the repo holds a `go.mod`: the Go modules of its commands (SR_36) and the
  standard library of Go (SR_44)
- for each `package-lock.json` of the repo: its packages that are not marked `dev` (SR_37)
- for each image named: the Alpine packages of `/lib/apk/db/installed` (SR_38)

A run keeps the rows of the Change history of an `sbom.md` that is there and adds one, with the
version of `sbom.md` and the ticket given (SR_40). It sets the entry of `sbom.md` in `documents.yaml`
of the repo: x of its version raised by one, the hash of the file written, `completed` and
`[generated]` - a generated document is a record and passes no review (PROC_13, SR_59). The repo
needs a `documents.yaml`. A package whose licence is unknown is listed in Open (SR_45).

Original data names the files `sbom.md` is written from (SR_101). With `-store` they are kept in
`sbom/`, beside `sbom.md`, and named in `sbom/<ts>.yaml` (SR_102).

```bash
cd tools
go build -o sbom ./sbom/cmd/sbom
./sbom -root ../../storages -label 1.2.0 -ticket 343 -image storages:1.2.0
```

With `-check` it writes nothing (TOOL_28, SR_66). It compares the packages delivered with the table
of `sbom.md`, with the label of the last row of its Change history as the version of the package of
the subsystem. It exits with 1 and names each package that differs: delivered and not listed,
listed and not delivered, or listed with another name, version or licence.

The code is the component `sbom`: the package `sbom/create` and the command `sbom/cmd/sbom`.

**Decision**: the licence of a Go module is recognized by sentences that only its licence text
holds, such as "Permission is hereby granted, free of charge" for MIT, in the licence file of its
directory in the module cache. A text that holds those of several licences names them all, joined
by AND; a text that holds none is unknown.

**Why**: Go keeps the licence file beside the code, but names no licence. The sentences need no
tool from outside and no list to keep up to date.

**Rejected**:

- the whole text compared with the texts of SPDX - the copyright lines and the spacing differ from
  module to module, and the texts would have to be kept here
- a licence scanner from outside - a new tool in the chain, which the threat of an outdated
  toolchain asks to keep small

**Decision**: a licence that a package names itself - the field `license` of `package-lock.json`,
the field `L` of an Alpine package - is taken as it is where it has the form of an SPDX expression:
identifiers joined by AND, OR or WITH. A licence named in words, such as "2-clause BSD-like
license", is unknown; Open shows the words.

**Why**: npm and Alpine name the licence in SPDX already. Words are no identifier: the check of the
licences (PROC_08) cannot compare them, so a human has to look.

**Rejected**: the identifiers checked against the list of SPDX - the list would have to be kept
here. The known limit: an identifier of the right form that SPDX does not list, such as
`Public-Domain`, passes.

**Decision**: the packages of an image are read by `docker run --rm --entrypoint cat <image>
/lib/apk/db/installed`, on the image built at the baseline.

**Why**: the image is what is delivered; its database of packages names what it holds, the
Dockerfile only what it asks for.

**Rejected**: unpacking the layers of `docker save` - the same content, with much more code.

**Decision**: the tool is the version of the commit of specs-and-deploy it was built from, as Go
records it in the build, with `-dirty` where the working tree had changes. A run with `go run`
records none and shows `(devel)`, so a release builds the tool first.

**Why**: the version needs no care by hand, and it names the code that wrote the SBOM.

**Rejected**: a version in the code - it would have to be raised by hand and can drift.

## licence check

Checks the licences of `sbom.md` against Apache-2.0 (TOOL_22; SR_47 to SR_50 in
[software-requirements.yaml](software-requirements.yaml)). It reads `sbom/sbom.md` of the
repo of a subsystem, and `licences_assessed.yaml` where the root holds one. It prints every package
that fails with its PURL and its licence, and exits with 1 where one fails.

```bash
cd tools
go build -o licencecheck ./licence/cmd/licencecheck
./licencecheck -root ../../storages
```

An entry of `licences_assessed.yaml`:

```yaml
- purl: pkg:apk/alpine/nginx@1.30.5-r1
  licence: BSD-2-Clause
  reason: the text of nginx is the BSD-2-Clause
```

The code is the component `licence`: the package `licence/check` and the command
`licence/cmd/licencecheck`.

**Decision**: the check tells two kinds of package apart. A package linked into the code of the
subsystem - a Go module, an npm package - needs a licence that fits Apache-2.0. A package that lies
beside it in an image - an Alpine package - may have any licence that is known. An unknown licence
fails both.

**Why**: a licence such as the GPL binds the code it is linked with, not a program that only lies
beside it. The image delivers such programs: busybox, the shell of Alpine, is GPL-2.0-only.

**Rejected**:

- one list for every package - every image on Alpine would fail
- any known licence for every package - a Go module under the GPL would pass

**Decision**: a licence that a human has assessed for one package is kept in
`licences_assessed.yaml` in the root of the repo of the subsystem: the PURL with its version, the
licence, and the reason. The check takes that licence in place of the one in `sbom.md`.

**Why**: some packages name their licence in words, such as nginx with "2-clause BSD-like
license". A human can read the text and name the licence; the check cannot. The version is part of
the PURL because a licence can change from one version to the next.

**Rejected**:

- the PURL without its version - an update would keep an assessment made for another text
- a list in specs-and-deploy for every subsystem - the packages differ per subsystem, and so does
  its baseline
- the name `licences_whitelist.yaml` - the file allows no licence; it names the licence of a
  package

## build steps

The build steps of every repo (TOOL_23 to TOOL_26; SR_51 to SR_58 in
[software-requirements.yaml](software-requirements.yaml)): build, linters and tests fail; the
report on test coverage and structural debt blocks nothing.

Where a pull request changes `sbom.md`, the job `build` also runs `steps.sh sbom` (TOOL_28, SR_67):
it builds each image the repo names in the input `images`, one `name=Dockerfile` per line, and runs
`sbom -check` with them. So a pull request with an `sbom.md` that no longer fits its head cannot be
merged.

The workflow `.github/workflows/build.yml` has a job for each: `build`, `lint`, `test` and
`report`. Each runs `.github/build/steps.sh` in the root of the repo. A Go module is every directory
with a `go.mod`, a JavaScript package every directory with a `package-lock.json`. A package names
`build` and `test` in the scripts of its `package.json`, and `coverage` where it reports its test
coverage. The script runs by hand too, with go, golangci-lint v2, deadcode and node on the `PATH`:

```bash
../specs-and-deploy/.github/build/steps.sh lint
../specs-and-deploy/.github/build/steps.sh sbom storages=Dockerfile
```

`.github/build` holds the configurations: `.golangci.yml` for the gate, `debt.golangci.yml` for the
structural debt, and `eslint.config.js`. They lie there, not in the root, because a YAML file outside
`.github` is our data and names a schema.

The code is the component `ci`: the workflow and the files of `.github/build`.

**Decision**: the build steps are one reusable workflow, `.github/workflows/build.yml` of
specs-and-deploy, which every repo calls:

```yaml
jobs:
  build:
    uses: treporting/specs-and-deploy/.github/workflows/build.yml@main
    with:
      images: |
        storages=Dockerfile
```

**Why**: a reusable workflow may lie in any public repo. Here the requirement, the code and its
`@implements` stay in one place, beside the tools it uses.

**Rejected**: a repo `treporting/.github` - one more repo with a baseline of its own, its SRs apart
from its code.

**Decision**: the linters are configured once, by `.golangci.yml` and `eslint.config.js` of
specs-and-deploy, for every repo.

**Why**: every repo is held to the same rules, and a change of them is a change of an SR with its
review.

**Rejected**: a configuration in each repo - it can drift from repo to repo unnoticed. The cost of
the decision: an exception for one file of one repo is made here too.

**Decision**: the structural debt is what nothing uses - unused and unparam with `--tests=false`,
and deadcode - and every function of a cyclomatic complexity of 15 and more (gocyclo). The test
coverage of Go is measured with `-coverpkg=./...`.

**Why**: code nothing uses, and a function too complex to test, are what grows unseen. With
`--tests=false` a function only the tests still call is named. With `-coverpkg` a test counts for
the package it reaches, not only for its own. 15 is the example gocyclo gives; a report blocks
nothing, so the stricter bound costs nothing.

**Rejected**:

- unused and unparam in the linters of the gate - they report on the line of a declaration, which
  does not change when its last caller goes
- the bound 30 of golangci-lint - it names hardly any function

## baseline

Prepares the baseline of a repo (TOOL_27; SR_60 to SR_65 in
[software-requirements.yaml](software-requirements.yaml)). Whoever wants one opens a ticket with the
template Baseline and names the label in it. Whoever may label tickets of the repo sets the label
`baseline` on it. The caller of the repo, `.github/workflows/ticket-baseline.yml`, then runs
`.github/workflows/baseline.yml` of specs-and-deploy:

1. `start`: reads the label from the field Label of the ticket - letters, digits, `.`, `_` and `-`
   only, as it becomes a tag. It fails where the field is empty, or where the branch
   `<ticket number>-baseline-<label>`, the tag `<label>` or another run of the caller exists.
2. `build`: the jobs of `build.yml` on the head of `main` as the ticket event names it.
3. `sbom`: `steps.sh baseline <label> <ticket number> [name=Dockerfile]...` builds the images,
   writes `sbom.md` and its entry in `documents.yaml`, checks the licences, and checks the
   approvals of every approved object (approvalcheck, SR_82). Where it passes, it
   pushes the branch with one commit and opens the pull request, with a token of the GitHub App
   `treporting-baseline` (devops.md, section Baseline).
4. `result`: a comment on the ticket - the pull request, or the step that failed, with the link to
   the run.

To start it again, take the label `baseline` off the ticket and set it again. The step runs by hand
too, in the root of the repo; it writes `sbom/sbom.md` and `documents.yaml` and opens nothing:

```bash
../specs-and-deploy/.github/build/steps.sh baseline 1.2.0 343 storages=Dockerfile
```

The code is the component `baseline`: the two workflows, the template, and the step `baseline` of
`steps.sh`.

**Decision**: every job of a run checks out `github.sha`, the head of `main` when the label was set.

**Why**: `main` may move while the run lasts. Each job checks out on its own; with the head of each
moment, sbom.md could describe another commit than the one built and tested.

**Rejected**: the head of `main` in each job - shorter, and wrong as soon as a pull request is merged
during the run.

## pdf

Makes the documents of a commit or tag into one ZIP of PDFs (TOOL_01, TOOL_02, TOOL_30, TOOL_31;
SR_68 to SR_74 in [software-requirements.yaml](software-requirements.yaml)). It reads the repo at the
ref without checking it out: `documents.yaml`, and the Markdown of every document it names. The ZIP
`<repo>-<ref>-documents.zip` holds one PDF per document, named after its path, and `index.pdf`,
which lists every document with its title, version, revision and status.

- the header of a page: the title of the document, the name of the repo
- the footer: the date of the commit of the ref, page x of y, and the path, version, revision and
  status of the document, such as `devops.md · 2.1.1 · a1b2c3d · approved` - the version and status
  from `documents.yaml` at the ref, the revision the last commit up to the ref that changed the file
- A4, margins of 2 cm at top and bottom, 1.5 cm at the sides

Markdown is rendered as GitHub renders it (goldmark with its extension GFM), and printed by headless
Chrome (chromedp). A document with a Mermaid diagram loads Mermaid 11.4.1 from jsDelivr, with the
hash of the script: Chrome refuses a script that differs, and the run fails.

The workflow `.github/workflows/pdf.yml` runs it, started by hand: Actions, pdf, Run workflow, with
the commit or tag. The ZIP is an artifact of the run, to download from there; GitHub keeps it 90
days, and the ZIP of a tag can be made again. By hand, from `tools/`, with Google Chrome or Chromium
installed:

```bash
go run ./pdf/cmd/pdf -root .. -ref T044 -out ..
```

## testreport

Writes the test report of a test level as `testreports/<test level>.md` beside `testcases/`
(TOOL_34; SR_88 to SR_98 and SR_104 in [software-requirements.yaml](software-requirements.yaml)),
from the junit.xml of a run of the tests. A test that fails does not fail it; the run of the tests
does. With `-store` it copies each junit.xml into `testreports/` beside the test report, named in
`testreports/<ts>_<test level>.yaml`, as sbom does with its original data.

```bash
cd tools
go install gotest.tools/gotestsum@v1.12.3
gotestsum --junitfile junit-go.xml -- -tags integration ./...
go build -o testreport ./testreport/cmd/testreport
./testreport -root .. -dir tools -level r5tc -junit junit-go.xml \
  -made-by "gotestsum --junitfile junit-go.xml -- -tags integration ./..." -runtime "Go $(go env GOVERSION)"
```

The code is the component `testreport`: the package `testreport/write` and the command
`testreport/cmd/testreport`. It uses the entry in `documents.yaml` of the component `sbom`.

**Decision**: a test procedure is a test function of Go, or an `it` or `test` of Vitest, with an
`@implements` or an `@verifies` in the comment above it. Its result in the junit.xml is found by its
package and name for Go, by its file and the end of its name for Vitest.

**Why**: gotestsum names a test by the import path of its package; Vitest by its file, relative to
its own root, and by the names of its describe blocks before its own.

**Rejected**: `go test -json` as input - Vitest does not write it; JUnit XML is written by both.

## Open

- `treporting = 1_tools` - **Checks of the code repos**: not set up yet. They are to call
  `build.yml` of specs-and-deploy.

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| 1.1.1 | 2026-10-28 | - | #272 | Claude Opus 5.5 and Mathias Fischer, authors; Mathias Fischer, reviewer |
| 2.1.1 | 2026-10-28 | annotationcheck | #279 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 3.1.1 | 2026-10-28 | uplinkcheck; linkcheck | #280 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 4.1.1 | 2026-10-28 | The component review | #290 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 5.1.1 | 2026-10-28 | Who changed a requirement | #40 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 6.1.1 | 2026-10-28 | The chain of a tool; Who changed a requirement; Who the author of a requirement or a document is; reviewcheck; reviewskeleton; schemacheck; rulesetcheck; annotationcheck; uplinkcheck; linkcheck | #291 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 7.1.1 | 2026-10-29 | annotationcheck; linkcheck | #283 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 8.1.1 | 2026-10-29 | reviewskeleton | #304 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 9.1.1 | 2026-10-29 | reviewcheck; reviewskeleton | #307 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 10.1.1 | 2026-10-29 | The component review | #303 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 11.1.1 | 2026-10-29 | Where the tools live; linkcheck; Open | #322 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 12.1.1 | 2026-10-30 | Open | #322 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 13.1.1 | 2026-10-30 | Change history | #323 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 14.1.1 | 2026-10-30 | Who the author of a requirement or a document is | #320 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 15.1.1 | 2026-10-31 | sbom | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 16.1.1 | 2026-10-31 | licence check | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 17.1.1 | 2026-10-31 | licence check | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 18.1.1 | 2026-10-31 | build steps; Open | #342 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 19.1.1 | 2026-10-31 | build steps | #342 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 20.1.1 | 2026-11-01 | sbom | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 21.1.1 | 2026-11-01 | sbom; build steps | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 22.1.1 | 2026-11-01 | baseline | #343 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 23.1.1 | 2026-11-01 | sbom; Open | #370 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 24.1.1 | 2026-11-01 | pdf | #392 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 25.1.1 | 2026-11-02 | The component review; approvalcheck; baseline | #399 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 26.1.1 | 2026-11-02 | reviewcheck; reviewskeleton | #153 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 27.1.1 | 2026-11-03 | approvalcheck: an open review result | #429 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 28.1.1 | 2026-11-03 | sbom; approvalcheck: a generated document is a completed record | #423 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 29.1.1 | 2026-11-03 | linkcheck: the targets of a test case | #449 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 30.1.1 | 2026-11-04 | testreport | #438 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 31.1.1 | 2026-11-04 | sbom; licence check; baseline: sbom.md in sbom/ | #491 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 32.1.1 | 2026-11-05 | sbom: Original data and -store | #480 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 33.1.1 | 2026-11-05 | reviewcheck: -base | #497 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 34.1.1 | 2026-11-05 | testreport: -store | #492 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
| 35.1.1 | 2026-11-12 | schemacheck, uplinkcheck, linkcheck: -above | #638 | Claude Opus 5.5, author; Mathias Fischer, reviewer |
