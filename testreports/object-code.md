# Object code report

Each function of the source code of each build output of the repo: whether it is in the
output, the hash of its code there and its coverage by the tests. It is written by the tool
objectcode (TOOL_35) and not edited by hand.

## Build

- Label: T048.0
- Commit: 34ce1b3
- Date: 2026-11-24

## annotationcheck

- Binary: annotationcheck (annotation/cmd/annotationcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| annotation/check/check.go@Check | github.com/treporting/tools/annotation/check.Check | yes | d26748e7 | yes | 6/6 |
| annotation/check/check.go@Find | github.com/treporting/tools/annotation/check.Find | yes | 1b830a85 | yes | 15/16 |
| annotation/check/check.go@Finding.String | github.com/treporting/tools/annotation/check.Finding.String | yes | ab3d8b52 | yes | 3/3 |
| annotation/check/check.go@Kinds.kindOf | github.com/treporting/tools/annotation/check.Kinds.kindOf | yes | 5a8f325c | yes | 10/10 |
| annotation/check/check.go@ReadKinds | github.com/treporting/tools/annotation/check.ReadKinds | yes | 82847ceb | no | 0/5 |
| annotation/check/check.go@Refines | github.com/treporting/tools/annotation/check.Refines | yes | 3e6e021c | yes | 18/20 |
| annotation/check/check.go@Unversioned | github.com/treporting/tools/annotation/check.Unversioned | yes | c80ac058 | yes | 6/6 |
| annotation/check/check.go@fileOf | github.com/treporting/tools/annotation/check.fileOf | yes | 27c79e27 | yes | 6/6 |
| annotation/check/check.go@isData | github.com/treporting/tools/annotation/check.isData | yes | a935bfb1 | yes | 4/4 |
| annotation/check/check.go@named | github.com/treporting/tools/annotation/check.named | yes | eff29ad7 | yes | 6/6 |
| annotation/check/check.go@objectIn | github.com/treporting/tools/annotation/check.objectIn | yes | 22a8cdf4 | yes | 11/13 |
| annotation/check/check.go@sameKey | github.com/treporting/tools/annotation/check.sameKey | yes | 8b76abff | yes | 3/3 |
| annotation/check/check.go@splitKey | github.com/treporting/tools/annotation/check.splitKey | yes | b9b9ed78 | yes | 8/11 |
| annotation/check/check.go@status | github.com/treporting/tools/annotation/check.status | yes | 6863666d | yes | 12/15 |
| annotation/cmd/annotationcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| annotation/cmd/annotationcheck/main.go@main | main.main | yes | ace10c59 | no | 0/38 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | github.com/treporting/tools/review.ExecGit.Show | yes | 1d0edba8 | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | github.com/treporting/tools/review.(*Index).FileFrom | yes | c1764129 | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | github.com/treporting/tools/review.(*Index).Keys | yes | 9012554f | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | github.com/treporting/tools/review.commonFolders | yes | d06ac4c1 | yes | 5/5 |

## licencecheck

- Binary: licencecheck (licence/cmd/licencecheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| licence/check/check.go@Check | github.com/treporting/tools/licence/check.Check | yes | 956b0ed2 | yes | 12/12 |
| licence/check/check.go@Failure.String | github.com/treporting/tools/licence/check.Failure.String | yes | 0e7dae41 | yes | 1/1 |
| licence/check/check.go@Fits | github.com/treporting/tools/licence/check.Fits | yes | 49652093 | yes | 4/4 |
| licence/check/check.go@Packages | github.com/treporting/tools/licence/check.Packages | yes | 6e48e5bf | yes | 12/12 |
| licence/check/check.go@ReadAssessed | github.com/treporting/tools/licence/check.ReadAssessed | yes | eb69d723 | yes | 9/9 |
| licence/check/check.go@kind | github.com/treporting/tools/licence/check.kind | yes | 2f48fe05 | yes | 4/5 |
| licence/check/check.go@parser.and | github.com/treporting/tools/licence/check.(*parser).and | yes | 1bb15649 | yes | 7/7 |
| licence/check/check.go@parser.next | - | no | - | yes | 3/3 |
| licence/check/check.go@parser.or | github.com/treporting/tools/licence/check.(*parser).or | yes | 5a1a10ba | yes | 7/7 |
| licence/check/check.go@parser.term | github.com/treporting/tools/licence/check.(*parser).term | yes | 5b94d860 | yes | 16/16 |
| licence/cmd/licencecheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| licence/cmd/licencecheck/main.go@main | main.main | yes | 1e9f01c8 | no | 0/18 |

## linkcheck

- Binary: linkcheck (link/cmd/linkcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| annotation/check/check.go@Check | - | no | - | yes | 6/6 |
| annotation/check/check.go@Find | github.com/treporting/tools/annotation/check.Find | yes | 1b830a85 | yes | 15/16 |
| annotation/check/check.go@Finding.String | - | no | - | yes | 3/3 |
| annotation/check/check.go@Kinds.kindOf | github.com/treporting/tools/annotation/check.Kinds.kindOf | yes | 5a8f325c | yes | 10/10 |
| annotation/check/check.go@ReadKinds | github.com/treporting/tools/annotation/check.ReadKinds | yes | 82847ceb | no | 0/5 |
| annotation/check/check.go@Refines | github.com/treporting/tools/annotation/check.Refines | yes | 3e6e021c | yes | 18/20 |
| annotation/check/check.go@Unversioned | - | no | - | yes | 6/6 |
| annotation/check/check.go@fileOf | - | no | - | yes | 6/6 |
| annotation/check/check.go@isData | github.com/treporting/tools/annotation/check.isData | yes | a935bfb1 | yes | 4/4 |
| annotation/check/check.go@named | github.com/treporting/tools/annotation/check.named | yes | eff29ad7 | yes | 6/6 |
| annotation/check/check.go@objectIn | - | no | - | yes | 11/13 |
| annotation/check/check.go@sameKey | - | no | - | yes | 3/3 |
| annotation/check/check.go@splitKey | - | no | - | yes | 8/11 |
| annotation/check/check.go@status | - | no | - | yes | 12/15 |
| link/check/check.go@Annotations | - | no | - | yes | 1/1 |
| link/check/check.go@AnnotationsWithAbove | github.com/treporting/tools/link/check.AnnotationsWithAbove | yes | 56e16b5d | yes | 16/16 |
| link/check/check.go@Finding.String | github.com/treporting/tools/link/check.Finding.String | yes | 0a98e934 | yes | 3/3 |
| link/check/check.go@Links | - | no | - | yes | 1/1 |
| link/check/check.go@LinksWithAbove | github.com/treporting/tools/link/check.LinksWithAbove | yes | 9211dbdb | yes | 18/18 |
| link/check/check.go@newVersions | github.com/treporting/tools/link/check.newVersions | yes | c7deac96 | yes | 4/4 |
| link/check/check.go@object.links | - | no | - | yes | 6/6 |
| link/check/check.go@objects | github.com/treporting/tools/link/check.objects | yes | 0850682a | yes | 11/12 |
| link/check/check.go@versions.find | github.com/treporting/tools/link/check.(*versions).find | yes | 9434e92f | yes | 11/13 |
| link/check/check.go@versions.load | github.com/treporting/tools/link/check.(*versions).load | yes | 86be36ec | yes | 6/7 |
| link/check/check.go@versions.obsolete | - | no | - | yes | 2/2 |
| link/check/check.go@versions.suspect | github.com/treporting/tools/link/check.(*versions).suspect | yes | 4c452a07 | yes | 9/9 |
| link/check/check.go@x | github.com/treporting/tools/link/check.x | yes | e748fca1 | yes | 2/2 |
| link/cmd/linkcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| link/cmd/linkcheck/main.go@main | main.main | yes | 1c0dab22 | no | 0/42 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | github.com/treporting/tools/review.(*Index).FileFrom | yes | c1764129 | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | github.com/treporting/tools/review.RepoAt | yes | 77d2e32a | no | 0/5 |
| review/objects.go@commonFolders | github.com/treporting/tools/review.commonFolders | yes | d06ac4c1 | yes | 5/5 |

## objectcode

- Binary: objectcode (objectcode/cmd/objectcode)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| annotation/check/check.go@Check | - | no | - | yes | 6/6 |
| annotation/check/check.go@Find | - | no | - | yes | 15/16 |
| annotation/check/check.go@Finding.String | - | no | - | yes | 3/3 |
| annotation/check/check.go@Kinds.kindOf | - | no | - | yes | 10/10 |
| annotation/check/check.go@ReadKinds | - | no | - | no | 0/5 |
| annotation/check/check.go@Refines | - | no | - | yes | 18/20 |
| annotation/check/check.go@Unversioned | - | no | - | yes | 6/6 |
| annotation/check/check.go@fileOf | - | no | - | yes | 6/6 |
| annotation/check/check.go@isData | - | no | - | yes | 4/4 |
| annotation/check/check.go@named | - | no | - | yes | 6/6 |
| annotation/check/check.go@objectIn | - | no | - | yes | 11/13 |
| annotation/check/check.go@sameKey | - | no | - | yes | 3/3 |
| annotation/check/check.go@splitKey | - | no | - | yes | 8/11 |
| annotation/check/check.go@status | - | no | - | yes | 12/15 |
| objectcode/cmd/objectcode/main.go@exists | - | no | - | no | 0/2 |
| objectcode/cmd/objectcode/main.go@fail | main.fail | yes | 9e70a734 | no | 0/2 |
| objectcode/cmd/objectcode/main.go@goBinary | main.goBinary | yes | 9742cfc6 | no | 0/24 |
| objectcode/cmd/objectcode/main.go@head | main.head | yes | bc438f57 | no | 0/4 |
| objectcode/cmd/objectcode/main.go@jsBundle | main.jsBundle | yes | 13661b91 | no | 0/39 |
| objectcode/cmd/objectcode/main.go@list.Set | main.(*list).Set | yes | 582b006e | no | 0/2 |
| objectcode/cmd/objectcode/main.go@list.String | main.(*list).String | yes | bbac04ab | no | 0/1 |
| objectcode/cmd/objectcode/main.go@main | main.main | yes | 598af42e | no | 0/26 |
| objectcode/cmd/objectcode/main.go@row | main.row | yes | ada09c15 | no | 0/3 |
| objectcode/cmd/objectcode/main.go@version | main.version | yes | 4c6ad5eb | no | 0/15 |
| objectcode/cmd/objectcode/main.go@writeDocument | main.writeDocument | yes | 550e047b | no | 0/27 |
| objectcode/write/bundle.go@Code.Find | github.com/treporting/tools/objectcode/write.(*Code).Find | yes | fddce620 | yes | 10/11 |
| objectcode/write/bundle.go@Code.functionOf | github.com/treporting/tools/objectcode/write.(*Code).functionOf | yes | 607d1458 | yes | 10/13 |
| objectcode/write/bundle.go@Code.index | github.com/treporting/tools/objectcode/write.(*Code).index | yes | 686fc8b9 | yes | 15/17 |
| objectcode/write/bundle.go@Code.normalized | github.com/treporting/tools/objectcode/write.(*Code).normalized | yes | 7dbd8cd0 | yes | 24/25 |
| objectcode/write/bundle.go@Code.source | github.com/treporting/tools/objectcode/write.(*Code).source | yes | 6069177a | yes | 4/4 |
| objectcode/write/bundle.go@ConfigurationHash | github.com/treporting/tools/objectcode/write.ConfigurationHash | yes | a40f969d | yes | 7/8 |
| objectcode/write/bundle.go@Istanbul.Coverage | github.com/treporting/tools/objectcode/write.(*Istanbul).Coverage | yes | f6c09234 | yes | 2/2 |
| objectcode/write/bundle.go@Istanbul.Functions | - | no | - | yes | 1/1 |
| objectcode/write/bundle.go@ReadBundle | github.com/treporting/tools/objectcode/write.ReadBundle | yes | ddb3ab72 | yes | 48/59 |
| objectcode/write/bundle.go@ReadIstanbul | github.com/treporting/tools/objectcode/write.ReadIstanbul | yes | 4621968e | yes | 39/44 |
| objectcode/write/bundle.go@ReadTools | github.com/treporting/tools/objectcode/write.ReadTools | yes | 8620af2a | yes | 8/10 |
| objectcode/write/bundle.go@anchorOf | github.com/treporting/tools/objectcode/write.anchorOf | yes | 26ed8a19 | yes | 5/14 |
| objectcode/write/bundle.go@columnOf | - | no | - | yes | 3/3 |
| objectcode/write/bundle.go@position.before | - | no | - | yes | 1/1 |
| objectcode/write/bundle.go@sourceFile | github.com/treporting/tools/objectcode/write.sourceFile | yes | 47d68b4c | yes | 4/4 |
| objectcode/write/bundle.go@span.from | - | no | - | yes | 1/1 |
| objectcode/write/bundle.go@span.to | - | no | - | yes | 1/1 |
| objectcode/write/class.go@Classes | github.com/treporting/tools/objectcode/write.Classes | yes | 5f83a2fd | yes | 34/37 |
| objectcode/write/cover.go@Profile.Coverage | github.com/treporting/tools/objectcode/write.(*Profile).Coverage | yes | f9e00fe1 | yes | 9/10 |
| objectcode/write/cover.go@ReadProfile | github.com/treporting/tools/objectcode/write.ReadProfile | yes | 8c3c4577 | yes | 23/25 |
| objectcode/write/disassemble.go@Disassembly.Find | github.com/treporting/tools/objectcode/write.Disassembly.Find | yes | dfbc8e5c | yes | 2/2 |
| objectcode/write/disassemble.go@Disassembly.Hash | github.com/treporting/tools/objectcode/write.Disassembly.Hash | yes | ed9ce75b | yes | 13/13 |
| objectcode/write/disassemble.go@Normalize | github.com/treporting/tools/objectcode/write.Normalize | yes | 68a1aa39 | yes | 12/12 |
| objectcode/write/disassemble.go@Objdump | github.com/treporting/tools/objectcode/write.Objdump | yes | 93a571a2 | yes | 12/13 |
| objectcode/write/disassemble.go@ReadObjdump | github.com/treporting/tools/objectcode/write.ReadObjdump | yes | 8242424b | yes | 18/19 |
| objectcode/write/disassemble.go@unnamed | - | no | - | yes | 4/4 |
| objectcode/write/document.go@Markdown | github.com/treporting/tools/objectcode/write.Markdown | yes | 4cafe231 | yes | 19/19 |
| objectcode/write/section.go@Binary.Kind | github.com/treporting/tools/objectcode/write.Binary.Kind | yes | 9b4d408d | yes | 1/1 |
| objectcode/write/section.go@Binary.Lines | github.com/treporting/tools/objectcode/write.Binary.Lines | yes | dfe5065c | yes | 1/1 |
| objectcode/write/section.go@Binary.Title | github.com/treporting/tools/objectcode/write.Binary.Title | yes | b2863530 | yes | 1/1 |
| objectcode/write/section.go@Bundle.Kind | github.com/treporting/tools/objectcode/write.Bundle.Kind | yes | 9b4d408d | no | 0/1 |
| objectcode/write/section.go@Bundle.Lines | github.com/treporting/tools/objectcode/write.Bundle.Lines | yes | d0d6d318 | yes | 4/4 |
| objectcode/write/section.go@Bundle.Title | github.com/treporting/tools/objectcode/write.Bundle.Title | yes | b2863530 | yes | 1/1 |
| objectcode/write/section.go@Note | github.com/treporting/tools/objectcode/write.Note | yes | 0fba92eb | yes | 6/6 |
| objectcode/write/section.go@ReadBinary | github.com/treporting/tools/objectcode/write.ReadBinary | yes | 500315ad | yes | 10/13 |
| objectcode/write/section.go@Rows | github.com/treporting/tools/objectcode/write.Rows | yes | 7993c6ec | yes | 7/7 |
| objectcode/write/section.go@Section | github.com/treporting/tools/objectcode/write.Section | yes | 0f228883 | yes | 15/15 |
| objectcode/write/source.go@Functions | github.com/treporting/tools/objectcode/write.Functions | yes | dfee67c7 | yes | 20/22 |
| objectcode/write/source.go@ListPackages | github.com/treporting/tools/objectcode/write.ListPackages | yes | c1acd933 | yes | 5/6 |
| objectcode/write/source.go@nameOf | - | no | - | yes | 3/3 |
| objectcode/write/source.go@readPackages | github.com/treporting/tools/objectcode/write.readPackages | yes | 46508f0c | yes | 14/15 |
| objectcode/write/source.go@receiverOf | github.com/treporting/tools/objectcode/write.receiverOf | yes | 475e9184 | yes | 10/11 |
| objectcode/write/source.go@symbolOf | github.com/treporting/tools/objectcode/write.symbolOf | yes | 65865d8a | yes | 5/5 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.run | - | no | - | yes | 2/2 |
| review/objects.go@Index.FileFrom | - | no | - | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | - | no | - | yes | 4/5 |
| review/objects.go@NewIndex | - | no | - | yes | 20/22 |
| review/objects.go@ObjectIn | - | no | - | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | - | no | - | yes | 5/5 |
| sbom/create/create.go@Differ | - | no | - | yes | 17/17 |
| sbom/create/create.go@History | - | no | - | yes | 9/9 |
| sbom/create/create.go@Label | - | no | - | yes | 5/6 |
| sbom/create/create.go@Markdown | - | no | - | yes | 44/44 |
| sbom/create/create.go@Once | - | no | - | yes | 7/7 |
| sbom/create/create.go@Own | - | no | - | yes | 1/1 |
| sbom/create/create.go@Row | - | no | - | yes | 1/1 |
| sbom/create/create.go@Table | - | no | - | yes | 12/12 |
| sbom/create/create.go@cell | - | no | - | yes | 1/1 |
| sbom/create/create.go@cells | - | no | - | yes | 5/5 |
| sbom/create/documents.go@Entry | - | no | - | yes | 1/1 |
| sbom/create/documents.go@EntryOf | github.com/treporting/tools/sbom/create.EntryOf | yes | 270bae0b | yes | 19/20 |
| sbom/create/documents.go@Hash | github.com/treporting/tools/sbom/create.Hash | yes | f7be713f | yes | 4/4 |
| sbom/create/documents.go@Next | github.com/treporting/tools/sbom/create.Next | yes | 17f99770 | yes | 9/9 |
| sbom/create/documents.go@Version | - | no | - | yes | 1/1 |
| sbom/create/documents.go@VersionOf | github.com/treporting/tools/sbom/create.VersionOf | yes | a616821d | yes | 6/6 |
| sbom/create/documents.go@entry | github.com/treporting/tools/sbom/create.entry | yes | bad71e16 | yes | 8/8 |
| sbom/create/licence.go@Given | - | no | - | yes | 4/4 |
| sbom/create/licence.go@Licence | - | no | - | yes | 8/8 |
| sbom/create/licence.go@holdsAll | - | no | - | yes | 4/4 |
| sbom/create/licence.go@holdsAny | - | no | - | yes | 4/4 |
| sbom/create/licence.go@identifier | - | no | - | yes | 5/5 |
| sbom/create/licence.go@spdxForm | - | no | - | yes | 8/8 |
| sbom/create/originals.go@ApkOriginal | - | no | - | yes | 1/1 |
| sbom/create/originals.go@GoListOriginal | - | no | - | yes | 3/3 |
| sbom/create/originals.go@LockOriginal | - | no | - | yes | 1/1 |
| sbom/create/originals.go@OriginalYAML | - | no | - | yes | 10/11 |
| sbom/create/originals.go@Store | - | no | - | yes | 8/10 |
| sbom/create/originals.go@dashes | - | no | - | yes | 1/1 |
| sbom/create/originals.go@in | - | no | - | yes | 4/4 |
| sbom/create/originals.go@slug | - | no | - | yes | 4/4 |
| sbom/create/sources.go@ApkInstalled | - | no | - | no | 0/2 |
| sbom/create/sources.go@ApkPackages | - | no | - | yes | 16/16 |
| sbom/create/sources.go@DirsWith | - | no | - | yes | 9/10 |
| sbom/create/sources.go@GoEnv | - | no | - | yes | 2/2 |
| sbom/create/sources.go@GoList | - | no | - | yes | 5/7 |
| sbom/create/sources.go@GoModules | - | no | - | yes | 9/9 |
| sbom/create/sources.go@LicenceIn | - | no | - | yes | 3/4 |
| sbom/create/sources.go@NpmPackages | - | no | - | yes | 14/14 |
| sbom/create/sources.go@Stdlib | - | no | - | yes | 1/1 |
| sbom/create/sources.go@run | - | no | - | yes | 5/8 |
| testreport/write/read.go@ReadJUnit | - | no | - | yes | 17/18 |
| testreport/write/read.go@ReadLevel | - | no | - | no | 0/9 |
| testreport/write/read.go@ReadProcedures | - | no | - | yes | 16/19 |
| testreport/write/read.go@ReadTestcases | - | no | - | no | 0/4 |
| testreport/write/read.go@inFile | - | no | - | yes | 20/20 |
| testreport/write/read.go@matched | - | no | - | yes | 3/4 |
| testreport/write/read.go@suite | - | no | - | yes | 5/6 |
| testreport/write/report.go@Flags | - | no | - | yes | 3/3 |
| testreport/write/report.go@History | github.com/treporting/tools/testreport/write.History | yes | 46bb0129 | no | 0/9 |
| testreport/write/report.go@Markdown | - | no | - | yes | 18/18 |
| testreport/write/report.go@Row | github.com/treporting/tools/testreport/write.Row | yes | 359ffecf | yes | 4/4 |
| testreport/write/report.go@cell | - | no | - | yes | 1/1 |
| testreport/write/report.go@deviations | - | no | - | yes | 8/8 |
| testreport/write/report.go@environment | - | no | - | yes | 11/11 |
| testreport/write/report.go@execution | - | no | - | yes | 5/5 |
| testreport/write/report.go@implementation | - | no | - | yes | 4/4 |
| testreport/write/report.go@oneLine | - | no | - | yes | 1/1 |
| testreport/write/report.go@open | - | no | - | yes | 21/21 |
| testreport/write/report.go@originals | - | no | - | yes | 6/7 |
| testreport/write/report.go@rows | - | no | - | yes | 12/12 |
| testreport/write/report.go@same | - | no | - | yes | 4/4 |
| testreport/write/report.go@specification | - | no | - | yes | 4/4 |
| testreport/write/report.go@testcaseKey | - | no | - | yes | 2/2 |
| testreport/write/report.go@verdicts | - | no | - | yes | 6/6 |
| testreport/write/report.go@worse | - | no | - | yes | 3/4 |

## pdf

- Binary: pdf (pdf/cmd/pdf)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| pdf/cmd/pdf/main.go@build | main.build | yes | a47974ca | no | 0/35 |
| pdf/cmd/pdf/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| pdf/cmd/pdf/main.go@main | main.main | yes | 68cd3737 | no | 0/22 |
| pdf/export/export.go@Documents | github.com/treporting/tools/pdf/export.Documents | yes | d288d9c7 | yes | 6/7 |
| pdf/export/export.go@Footer | - | no | - | yes | 1/1 |
| pdf/export/export.go@FooterTemplate | github.com/treporting/tools/pdf/export.FooterTemplate | yes | 497587f3 | yes | 1/1 |
| pdf/export/export.go@HeaderTemplate | github.com/treporting/tools/pdf/export.HeaderTemplate | yes | b0767333 | yes | 1/1 |
| pdf/export/export.go@Index | github.com/treporting/tools/pdf/export.Index | yes | 83affe21 | yes | 10/10 |
| pdf/export/export.go@PDFName | - | no | - | yes | 1/1 |
| pdf/export/export.go@Title | github.com/treporting/tools/pdf/export.Title | yes | c0b14251 | yes | 4/4 |
| pdf/export/export.go@Zip | github.com/treporting/tools/pdf/export.Zip | yes | b327e659 | yes | 6/8 |
| pdf/export/export.go@ZipName | - | no | - | yes | 1/1 |
| pdf/export/page.go@Page | github.com/treporting/tools/pdf/export.Page | yes | 895d1896 | yes | 6/7 |
| pdf/export/page.go@Wrap | - | no | - | yes | 1/1 |
| pdf/export/page.go@addresses.Transform | github.com/treporting/tools/pdf/export.addresses.Transform | yes | dd647102 | yes | 8/8 |
| pdf/export/page.go@letter | - | no | - | yes | 3/4 |
| pdf/export/sources.go@Git.Commit | github.com/treporting/tools/pdf/export.Git.Commit | yes | 2152645a | yes | 4/4 |
| pdf/export/sources.go@Git.Date | github.com/treporting/tools/pdf/export.Git.Date | yes | 114c2efa | yes | 2/2 |
| pdf/export/sources.go@Git.Revision | github.com/treporting/tools/pdf/export.Git.Revision | yes | a81b28e9 | yes | 2/2 |
| pdf/export/sources.go@Git.Show | - | no | - | yes | 2/2 |
| pdf/export/sources.go@NewPrinter | github.com/treporting/tools/pdf/export.NewPrinter | yes | f1fe5909 | yes | 10/14 |
| pdf/export/sources.go@Printer.Close | github.com/treporting/tools/pdf/export.(*Printer).Close | yes | 0beed62f | yes | 1/1 |
| pdf/export/sources.go@Printer.Print | github.com/treporting/tools/pdf/export.(*Printer).Print | yes | 427fda8b | yes | 12/14 |
| pdf/export/sources.go@git | github.com/treporting/tools/pdf/export.git | yes | 77e3ee40 | yes | 7/8 |

## approvalcheck

- Binary: approvalcheck (review/cmd/approvalcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| review/approval/approval.go@Check | - | no | - | yes | 1/1 |
| review/approval/approval.go@CheckWithAbove | github.com/treporting/tools/review/approval.CheckWithAbove | yes | adee0fce | yes | 22/25 |
| review/approval/approval.go@Finding.String | github.com/treporting/tools/review/approval.Finding.String | yes | 6395c9bd | yes | 1/1 |
| review/approval/approval.go@approvedDocuments | github.com/treporting/tools/review/approval.approvedDocuments | yes | a7a784f0 | yes | 14/18 |
| review/approval/approval.go@approvedRequirements | github.com/treporting/tools/review/approval.approvedRequirements | yes | eff53782 | yes | 14/17 |
| review/approval/approval.go@attributesForX | github.com/treporting/tools/review/approval.attributesForX | yes | 18a34584 | yes | 9/10 |
| review/approval/approval.go@blobHash | github.com/treporting/tools/review/approval.blobHash | yes | df6b7ddf | yes | 2/2 |
| review/approval/approval.go@changedDocuments | github.com/treporting/tools/review/approval.changedDocuments | yes | b6856345 | yes | 15/16 |
| review/approval/approval.go@checkReview | github.com/treporting/tools/review/approval.checkReview | yes | 44433346 | yes | 15/16 |
| review/approval/approval.go@checkVersions | github.com/treporting/tools/review/approval.checkVersions | yes | f3943d42 | yes | 10/11 |
| review/approval/approval.go@checked | - | no | - | yes | 1/1 |
| review/approval/approval.go@criteriaFile | - | no | - | yes | 3/3 |
| review/approval/approval.go@documentEntries | github.com/treporting/tools/review/approval.documentEntries | yes | 201ca8a2 | yes | 3/4 |
| review/approval/approval.go@findX | github.com/treporting/tools/review/approval.findX | yes | 1456f17c | yes | 12/13 |
| review/approval/approval.go@firstApprovals | github.com/treporting/tools/review/approval.firstApprovals | yes | 584c7416 | yes | 22/24 |
| review/approval/approval.go@inFiles | - | no | - | yes | 5/5 |
| review/approval/approval.go@itemsOf | github.com/treporting/tools/review/approval.itemsOf | yes | 1b15b4b4 | yes | 13/15 |
| review/approval/approval.go@keep | github.com/treporting/tools/review/approval.keep | yes | 7337a8c5 | yes | 2/3 |
| review/approval/approval.go@lastReviews | github.com/treporting/tools/review/approval.lastReviews | yes | 7e227b02 | yes | 18/22 |
| review/approval/approval.go@object.attributes | github.com/treporting/tools/review/approval.object.attributes | yes | 6e545349 | yes | 6/6 |
| review/approval/approval.go@reviews.of | github.com/treporting/tools/review/approval.reviews.of | yes | 575666b8 | yes | 5/5 |
| review/approval/approval.go@scoped | github.com/treporting/tools/review/approval.scoped | yes | 4534d42b | yes | 4/5 |
| review/approval/approval.go@xOf | github.com/treporting/tools/review/approval.xOf | yes | 4ec0ba1a | yes | 8/8 |
| review/cmd/approvalcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| review/cmd/approvalcheck/main.go@main | main.main | yes | 9f2ae8f9 | no | 0/15 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | github.com/treporting/tools/review.ExecGit.Changed | yes | e9ea21b5 | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | github.com/treporting/tools/review.ExecGit.FileCommits | yes | 6f94c02c | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | github.com/treporting/tools/review.ExecGit.Show | yes | 1d0edba8 | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | - | no | - | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | - | no | - | yes | 5/5 |

## reviewcheck

- Binary: reviewcheck (review/cmd/reviewcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| review/check/check.go@Check | github.com/treporting/tools/review/check.Check | yes | 5a1d5332 | yes | 14/16 |
| review/check/check.go@Finding.String | github.com/treporting/tools/review/check.Finding.String | yes | dc8d8602 | no | 0/1 |
| review/check/check.go@Judgement.UnmarshalYAML | github.com/treporting/tools/review/check.(*Judgement).UnmarshalYAML | yes | 77f3e879 | yes | 5/5 |
| review/check/check.go@Results | github.com/treporting/tools/review/check.Results | yes | 2b9e9e3d | yes | 17/21 |
| review/check/check.go@authorsOf | github.com/treporting/tools/review/check.authorsOf | yes | 2641d981 | yes | 19/24 |
| review/check/check.go@checkAuthors | github.com/treporting/tools/review/check.checkAuthors | yes | d5a505d2 | yes | 9/11 |
| review/check/check.go@checkCompleted | github.com/treporting/tools/review/check.checkCompleted | yes | 81dd4a4c | yes | 11/13 |
| review/check/check.go@checkCriteria | github.com/treporting/tools/review/check.checkCriteria | yes | 38ef4427 | yes | 22/22 |
| review/check/check.go@checkFile | github.com/treporting/tools/review/check.checkFile | yes | a0457650 | yes | 6/6 |
| review/check/check.go@checkJudgement | github.com/treporting/tools/review/check.checkJudgement | yes | e540ac4c | yes | 6/6 |
| review/check/check.go@checkResult | github.com/treporting/tools/review/check.checkResult | yes | a2810c0a | yes | 7/7 |
| review/check/check.go@checkVerdict | - | no | - | yes | 5/5 |
| review/check/check.go@criteriaKeys | github.com/treporting/tools/review/check.criteriaKeys | yes | d3595140 | yes | 13/17 |
| review/check/check.go@objectAt | github.com/treporting/tools/review/check.objectAt | yes | 5bee92ce | yes | 4/4 |
| review/cmd/reviewcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| review/cmd/reviewcheck/main.go@main | main.main | yes | 2ba5753d | no | 0/23 |
| review/git.go@ExecGit.Authors | github.com/treporting/tools/review.ExecGit.Authors | yes | 5aeddb12 | yes | 8/10 |
| review/git.go@ExecGit.Changed | github.com/treporting/tools/review.ExecGit.Changed | yes | e9ea21b5 | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | github.com/treporting/tools/review.ExecGit.CommitExists | yes | 48c622ee | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | github.com/treporting/tools/review.ExecGit.FileCommits | yes | 6f94c02c | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | github.com/treporting/tools/review.ExecGit.IsMerge | yes | 76b555ab | yes | 2/2 |
| review/git.go@ExecGit.Parent | github.com/treporting/tools/review.ExecGit.Parent | yes | 09dcc0ed | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | github.com/treporting/tools/review.ExecGit.Show | yes | 1d0edba8 | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | - | no | - | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | - | no | - | yes | 5/5 |

## reviewskeleton

- Binary: reviewskeleton (review/cmd/reviewskeleton)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| review/cmd/reviewskeleton/main.go@main | main.main | yes | 4e8d8893 | no | 0/17 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | github.com/treporting/tools/review.ExecGit.FileCommits | yes | 6f94c02c | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | github.com/treporting/tools/review.ExecGit.Short | yes | 5de74606 | yes | 2/2 |
| review/git.go@ExecGit.Show | github.com/treporting/tools/review.ExecGit.Show | yes | 1d0edba8 | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | - | no | - | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | - | no | - | yes | 5/5 |
| review/skeleton/skeleton.go@Make | github.com/treporting/tools/review/skeleton.Make | yes | 6ef0f0d5 | yes | 36/41 |
| review/skeleton/skeleton.go@criteriaOf | github.com/treporting/tools/review/skeleton.criteriaOf | yes | e1b56820 | yes | 14/18 |
| review/skeleton/skeleton.go@fileOf | github.com/treporting/tools/review/skeleton.fileOf | yes | 76f32b32 | yes | 7/8 |
| review/skeleton/skeleton.go@levelOf | github.com/treporting/tools/review/skeleton.levelOf | yes | 6bdc313a | yes | 18/20 |

## rulesetcheck

- Binary: rulesetcheck (ruleset/cmd/rulesetcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| ruleset/check/check.go@Check | github.com/treporting/tools/ruleset/check.Check | yes | 048af0b2 | yes | 20/22 |
| ruleset/check/check.go@Client.get | github.com/treporting/tools/ruleset/check.Client.get | yes | 061a25cf | yes | 16/20 |
| ruleset/check/check.go@Finding.String | github.com/treporting/tools/ruleset/check.Finding.String | yes | 7ac37a39 | yes | 3/3 |
| ruleset/check/check.go@compact | - | no | - | yes | 2/2 |
| ruleset/check/check.go@compareRules | github.com/treporting/tools/ruleset/check.compareRules | yes | 64106ff9 | yes | 21/22 |
| ruleset/cmd/rulesetcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| ruleset/cmd/rulesetcheck/main.go@main | main.main | yes | 0d99b780 | no | 0/15 |

## sbom

- Binary: sbom (sbom/cmd/sbom)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| sbom/cmd/sbom/main.go@apkPackages | main.apkPackages | yes | 25d1264f | no | 0/9 |
| sbom/cmd/sbom/main.go@checkSbom | main.checkSbom | yes | 50ce35c4 | no | 0/14 |
| sbom/cmd/sbom/main.go@commit | main.commit | yes | d56d8638 | no | 0/4 |
| sbom/cmd/sbom/main.go@delivered | main.delivered | yes | d58bcc79 | no | 0/7 |
| sbom/cmd/sbom/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| sbom/cmd/sbom/main.go@goPackages | main.goPackages | yes | a9b5b49e | no | 0/19 |
| sbom/cmd/sbom/main.go@images.Set | main.(*images).Set | yes | 472c57e6 | no | 0/2 |
| sbom/cmd/sbom/main.go@images.String | main.(*images).String | yes | e9932216 | no | 0/1 |
| sbom/cmd/sbom/main.go@main | main.main | yes | 143b393e | no | 0/45 |
| sbom/cmd/sbom/main.go@npmPackages | main.npmPackages | yes | 1c48add1 | no | 0/15 |
| sbom/cmd/sbom/main.go@stdlib | main.stdlib | yes | f047ef81 | no | 0/10 |
| sbom/cmd/sbom/main.go@storeOriginals | main.storeOriginals | yes | e2ce96cb | no | 0/10 |
| sbom/cmd/sbom/main.go@version | main.version | yes | 4c6ad5eb | no | 0/15 |
| sbom/create/create.go@Differ | github.com/treporting/tools/sbom/create.Differ | yes | b3f7ead9 | yes | 17/17 |
| sbom/create/create.go@History | github.com/treporting/tools/sbom/create.History | yes | 22006cc1 | yes | 9/9 |
| sbom/create/create.go@Label | github.com/treporting/tools/sbom/create.Label | yes | 9ebfe912 | yes | 5/6 |
| sbom/create/create.go@Markdown | github.com/treporting/tools/sbom/create.Markdown | yes | d260a5d5 | yes | 44/44 |
| sbom/create/create.go@Once | - | no | - | yes | 7/7 |
| sbom/create/create.go@Own | - | no | - | yes | 1/1 |
| sbom/create/create.go@Row | github.com/treporting/tools/sbom/create.Row | yes | 6021e396 | yes | 1/1 |
| sbom/create/create.go@Table | github.com/treporting/tools/sbom/create.Table | yes | 70a6ae21 | yes | 12/12 |
| sbom/create/create.go@cell | - | no | - | yes | 1/1 |
| sbom/create/create.go@cells | github.com/treporting/tools/sbom/create.cells | yes | 6c39c755 | yes | 5/5 |
| sbom/create/documents.go@Entry | - | no | - | yes | 1/1 |
| sbom/create/documents.go@EntryOf | github.com/treporting/tools/sbom/create.EntryOf | yes | 270bae0b | yes | 19/20 |
| sbom/create/documents.go@Hash | github.com/treporting/tools/sbom/create.Hash | yes | f7be713f | yes | 4/4 |
| sbom/create/documents.go@Next | github.com/treporting/tools/sbom/create.Next | yes | 17f99770 | yes | 9/9 |
| sbom/create/documents.go@Version | - | no | - | yes | 1/1 |
| sbom/create/documents.go@VersionOf | github.com/treporting/tools/sbom/create.VersionOf | yes | a616821d | yes | 6/6 |
| sbom/create/documents.go@entry | github.com/treporting/tools/sbom/create.entry | yes | bad71e16 | yes | 8/8 |
| sbom/create/licence.go@Given | github.com/treporting/tools/sbom/create.Given | yes | daa8c4c3 | yes | 4/4 |
| sbom/create/licence.go@Licence | github.com/treporting/tools/sbom/create.Licence | yes | ce31e5bd | yes | 8/8 |
| sbom/create/licence.go@holdsAll | github.com/treporting/tools/sbom/create.holdsAll | yes | b4bf39d8 | yes | 4/4 |
| sbom/create/licence.go@holdsAny | github.com/treporting/tools/sbom/create.holdsAny | yes | c556eee5 | yes | 4/4 |
| sbom/create/licence.go@identifier | github.com/treporting/tools/sbom/create.identifier | yes | 23040a8c | yes | 5/5 |
| sbom/create/licence.go@spdxForm | github.com/treporting/tools/sbom/create.spdxForm | yes | ec0fdb3e | yes | 8/8 |
| sbom/create/originals.go@ApkOriginal | - | no | - | yes | 1/1 |
| sbom/create/originals.go@GoListOriginal | github.com/treporting/tools/sbom/create.GoListOriginal | yes | be92e41b | yes | 3/3 |
| sbom/create/originals.go@LockOriginal | github.com/treporting/tools/sbom/create.LockOriginal | yes | d0cc7a55 | yes | 1/1 |
| sbom/create/originals.go@OriginalYAML | github.com/treporting/tools/sbom/create.OriginalYAML | yes | 92435c99 | yes | 10/11 |
| sbom/create/originals.go@Store | github.com/treporting/tools/sbom/create.Store | yes | ee378c91 | yes | 8/10 |
| sbom/create/originals.go@dashes | github.com/treporting/tools/sbom/create.dashes | yes | 9c1611b1 | yes | 1/1 |
| sbom/create/originals.go@in | github.com/treporting/tools/sbom/create.in | yes | 5f3e8e07 | yes | 4/4 |
| sbom/create/originals.go@slug | github.com/treporting/tools/sbom/create.slug | yes | 3465d3f0 | yes | 4/4 |
| sbom/create/sources.go@ApkInstalled | - | no | - | no | 0/2 |
| sbom/create/sources.go@ApkPackages | github.com/treporting/tools/sbom/create.ApkPackages | yes | 26f36ff1 | yes | 16/16 |
| sbom/create/sources.go@DirsWith | github.com/treporting/tools/sbom/create.DirsWith | yes | ac068e26 | yes | 9/10 |
| sbom/create/sources.go@GoEnv | github.com/treporting/tools/sbom/create.GoEnv | yes | c9631fae | yes | 2/2 |
| sbom/create/sources.go@GoList | github.com/treporting/tools/sbom/create.GoList | yes | 789d3343 | yes | 5/7 |
| sbom/create/sources.go@GoModules | github.com/treporting/tools/sbom/create.GoModules | yes | b2d42a65 | yes | 9/9 |
| sbom/create/sources.go@LicenceIn | github.com/treporting/tools/sbom/create.LicenceIn | yes | 08847334 | yes | 3/4 |
| sbom/create/sources.go@NpmPackages | github.com/treporting/tools/sbom/create.NpmPackages | yes | 098ba498 | yes | 14/14 |
| sbom/create/sources.go@Stdlib | - | no | - | yes | 1/1 |
| sbom/create/sources.go@run | github.com/treporting/tools/sbom/create.run | yes | 912522ff | yes | 5/8 |

## schemacheck

- Binary: schemacheck (schema/cmd/schemacheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| schema/check/check.go@Check | - | no | - | yes | 1/1 |
| schema/check/check.go@CheckWithAbove | github.com/treporting/tools/schema/check.CheckWithAbove | yes | 1bb7c2cf | yes | 33/40 |
| schema/check/check.go@Files | github.com/treporting/tools/schema/check.Files | yes | 7f820906 | yes | 9/10 |
| schema/check/check.go@Finding.String | github.com/treporting/tools/schema/check.Finding.String | yes | e13a56d2 | yes | 3/3 |
| schema/check/check.go@compile | github.com/treporting/tools/schema/check.compile | yes | 785c3424 | yes | 12/15 |
| schema/cmd/schemacheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| schema/cmd/schemacheck/main.go@main | main.main | yes | 79c63cd9 | no | 0/14 |

## testreport

- Binary: testreport (testreport/cmd/testreport)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| annotation/check/check.go@Check | - | no | - | yes | 6/6 |
| annotation/check/check.go@Find | github.com/treporting/tools/annotation/check.Find | yes | 1b830a85 | yes | 15/16 |
| annotation/check/check.go@Finding.String | - | no | - | yes | 3/3 |
| annotation/check/check.go@Kinds.kindOf | github.com/treporting/tools/annotation/check.Kinds.kindOf | yes | 5a8f325c | yes | 10/10 |
| annotation/check/check.go@ReadKinds | - | no | - | no | 0/5 |
| annotation/check/check.go@Refines | - | no | - | yes | 18/20 |
| annotation/check/check.go@Unversioned | - | no | - | yes | 6/6 |
| annotation/check/check.go@fileOf | - | no | - | yes | 6/6 |
| annotation/check/check.go@isData | github.com/treporting/tools/annotation/check.isData | yes | a935bfb1 | yes | 4/4 |
| annotation/check/check.go@named | github.com/treporting/tools/annotation/check.named | yes | eff29ad7 | yes | 6/6 |
| annotation/check/check.go@objectIn | - | no | - | yes | 11/13 |
| annotation/check/check.go@sameKey | - | no | - | yes | 3/3 |
| annotation/check/check.go@splitKey | - | no | - | yes | 8/11 |
| annotation/check/check.go@status | - | no | - | yes | 12/15 |
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.run | - | no | - | yes | 2/2 |
| review/objects.go@Index.FileFrom | - | no | - | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | - | no | - | yes | 4/5 |
| review/objects.go@NewIndex | - | no | - | yes | 20/22 |
| review/objects.go@ObjectIn | - | no | - | yes | 12/13 |
| review/objects.go@RepoAt | - | no | - | no | 0/5 |
| review/objects.go@commonFolders | - | no | - | yes | 5/5 |
| sbom/create/create.go@Differ | - | no | - | yes | 17/17 |
| sbom/create/create.go@History | - | no | - | yes | 9/9 |
| sbom/create/create.go@Label | - | no | - | yes | 5/6 |
| sbom/create/create.go@Markdown | - | no | - | yes | 44/44 |
| sbom/create/create.go@Once | - | no | - | yes | 7/7 |
| sbom/create/create.go@Own | - | no | - | yes | 1/1 |
| sbom/create/create.go@Row | - | no | - | yes | 1/1 |
| sbom/create/create.go@Table | - | no | - | yes | 12/12 |
| sbom/create/create.go@cell | - | no | - | yes | 1/1 |
| sbom/create/create.go@cells | - | no | - | yes | 5/5 |
| sbom/create/documents.go@Entry | - | no | - | yes | 1/1 |
| sbom/create/documents.go@EntryOf | github.com/treporting/tools/sbom/create.EntryOf | yes | 270bae0b | yes | 19/20 |
| sbom/create/documents.go@Hash | github.com/treporting/tools/sbom/create.Hash | yes | f7be713f | yes | 4/4 |
| sbom/create/documents.go@Next | github.com/treporting/tools/sbom/create.Next | yes | 17f99770 | yes | 9/9 |
| sbom/create/documents.go@Version | - | no | - | yes | 1/1 |
| sbom/create/documents.go@VersionOf | github.com/treporting/tools/sbom/create.VersionOf | yes | a616821d | yes | 6/6 |
| sbom/create/documents.go@entry | github.com/treporting/tools/sbom/create.entry | yes | bad71e16 | yes | 8/8 |
| sbom/create/licence.go@Given | - | no | - | yes | 4/4 |
| sbom/create/licence.go@Licence | - | no | - | yes | 8/8 |
| sbom/create/licence.go@holdsAll | - | no | - | yes | 4/4 |
| sbom/create/licence.go@holdsAny | - | no | - | yes | 4/4 |
| sbom/create/licence.go@identifier | - | no | - | yes | 5/5 |
| sbom/create/licence.go@spdxForm | - | no | - | yes | 8/8 |
| sbom/create/originals.go@ApkOriginal | - | no | - | yes | 1/1 |
| sbom/create/originals.go@GoListOriginal | - | no | - | yes | 3/3 |
| sbom/create/originals.go@LockOriginal | - | no | - | yes | 1/1 |
| sbom/create/originals.go@OriginalYAML | github.com/treporting/tools/sbom/create.OriginalYAML | yes | 92435c99 | yes | 10/11 |
| sbom/create/originals.go@Store | github.com/treporting/tools/sbom/create.Store | yes | ee378c91 | yes | 8/10 |
| sbom/create/originals.go@dashes | - | no | - | yes | 1/1 |
| sbom/create/originals.go@in | - | no | - | yes | 4/4 |
| sbom/create/originals.go@slug | - | no | - | yes | 4/4 |
| sbom/create/sources.go@ApkInstalled | - | no | - | no | 0/2 |
| sbom/create/sources.go@ApkPackages | - | no | - | yes | 16/16 |
| sbom/create/sources.go@DirsWith | - | no | - | yes | 9/10 |
| sbom/create/sources.go@GoEnv | - | no | - | yes | 2/2 |
| sbom/create/sources.go@GoList | - | no | - | yes | 5/7 |
| sbom/create/sources.go@GoModules | - | no | - | yes | 9/9 |
| sbom/create/sources.go@LicenceIn | - | no | - | yes | 3/4 |
| sbom/create/sources.go@NpmPackages | - | no | - | yes | 14/14 |
| sbom/create/sources.go@Stdlib | - | no | - | yes | 1/1 |
| sbom/create/sources.go@run | - | no | - | yes | 5/8 |
| testreport/cmd/testreport/main.go@commit | main.commit | yes | d56d8638 | no | 0/4 |
| testreport/cmd/testreport/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| testreport/cmd/testreport/main.go@gitFiles | main.gitFiles | yes | eb187a1b | no | 0/4 |
| testreport/cmd/testreport/main.go@list.Set | main.(*list).Set | yes | 582b006e | no | 0/2 |
| testreport/cmd/testreport/main.go@list.String | main.(*list).String | yes | bbac04ab | no | 0/1 |
| testreport/cmd/testreport/main.go@main | main.main | yes | 7ad58084 | no | 0/73 |
| testreport/cmd/testreport/main.go@moduleOf | main.moduleOf | yes | 11be33a8 | no | 0/7 |
| testreport/cmd/testreport/main.go@reported | - | no | - | no | 0/4 |
| testreport/cmd/testreport/main.go@runner | - | no | - | no | 0/3 |
| testreport/cmd/testreport/main.go@storeJUnits | main.storeJUnits | yes | fbe65eb4 | no | 0/10 |
| testreport/cmd/testreport/main.go@version | main.version | yes | 4c6ad5eb | no | 0/15 |
| testreport/write/read.go@ReadJUnit | github.com/treporting/tools/testreport/write.ReadJUnit | yes | 3029cf72 | yes | 17/18 |
| testreport/write/read.go@ReadLevel | github.com/treporting/tools/testreport/write.ReadLevel | yes | bd97e9e9 | no | 0/9 |
| testreport/write/read.go@ReadProcedures | github.com/treporting/tools/testreport/write.ReadProcedures | yes | 327289b3 | yes | 16/19 |
| testreport/write/read.go@ReadTestcases | github.com/treporting/tools/testreport/write.ReadTestcases | yes | 29a3ba17 | no | 0/4 |
| testreport/write/read.go@inFile | github.com/treporting/tools/testreport/write.inFile | yes | 5e79ea3a | yes | 20/20 |
| testreport/write/read.go@matched | github.com/treporting/tools/testreport/write.matched | yes | 8a494332 | yes | 3/4 |
| testreport/write/read.go@suite | github.com/treporting/tools/testreport/write.suite | yes | 99138f8e | yes | 5/6 |
| testreport/write/report.go@Flags | - | no | - | yes | 3/3 |
| testreport/write/report.go@History | github.com/treporting/tools/testreport/write.History | yes | 46bb0129 | no | 0/9 |
| testreport/write/report.go@Markdown | github.com/treporting/tools/testreport/write.Markdown | yes | 00d1ec11 | yes | 18/18 |
| testreport/write/report.go@Row | github.com/treporting/tools/testreport/write.Row | yes | 359ffecf | yes | 4/4 |
| testreport/write/report.go@cell | - | no | - | yes | 1/1 |
| testreport/write/report.go@deviations | github.com/treporting/tools/testreport/write.deviations | yes | 9172f356 | yes | 8/8 |
| testreport/write/report.go@environment | github.com/treporting/tools/testreport/write.environment | yes | 9c8e89b6 | yes | 11/11 |
| testreport/write/report.go@execution | github.com/treporting/tools/testreport/write.execution | yes | 1847491e | yes | 5/5 |
| testreport/write/report.go@implementation | github.com/treporting/tools/testreport/write.implementation | yes | 015b5b29 | yes | 4/4 |
| testreport/write/report.go@oneLine | github.com/treporting/tools/testreport/write.oneLine | yes | 1d26c5ad | yes | 1/1 |
| testreport/write/report.go@open | github.com/treporting/tools/testreport/write.open | yes | b29a813c | yes | 21/21 |
| testreport/write/report.go@originals | github.com/treporting/tools/testreport/write.originals | yes | ba598236 | yes | 6/7 |
| testreport/write/report.go@rows | github.com/treporting/tools/testreport/write.rows | yes | f97d110b | yes | 12/12 |
| testreport/write/report.go@same | github.com/treporting/tools/testreport/write.same | yes | 9ba74ec0 | yes | 4/4 |
| testreport/write/report.go@specification | github.com/treporting/tools/testreport/write.specification | yes | fab535d2 | yes | 4/4 |
| testreport/write/report.go@testcaseKey | github.com/treporting/tools/testreport/write.testcaseKey | yes | c141be90 | yes | 2/2 |
| testreport/write/report.go@verdicts | github.com/treporting/tools/testreport/write.verdicts | yes | 2566795a | yes | 6/6 |
| testreport/write/report.go@worse | - | no | - | yes | 3/4 |

## uplinkcheck

- Binary: uplinkcheck (uplink/cmd/uplinkcheck)
- Go go1.25.0, linux/amd64, commit 34ce1b3

| Function | Symbol | In binary | Hash | Called by tests | Statements |
|---|---|---|---|---|---|
| review/git.go@ExecGit.Authors | - | no | - | yes | 8/10 |
| review/git.go@ExecGit.Changed | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.CommitExists | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.FileCommits | - | no | - | yes | 3/4 |
| review/git.go@ExecGit.IsMerge | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Parent | - | no | - | yes | 4/4 |
| review/git.go@ExecGit.Short | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.Show | - | no | - | yes | 2/2 |
| review/git.go@ExecGit.run | github.com/treporting/tools/review.ExecGit.run | yes | 4ffc7119 | yes | 2/2 |
| review/objects.go@Index.FileFrom | github.com/treporting/tools/review.(*Index).FileFrom | yes | c1764129 | yes | 9/9 |
| review/objects.go@Index.FileOf | - | no | - | yes | 5/5 |
| review/objects.go@Index.Files | - | no | - | yes | 1/1 |
| review/objects.go@Index.FilesOf | - | no | - | yes | 3/3 |
| review/objects.go@Index.Keys | - | no | - | yes | 5/5 |
| review/objects.go@IndexAt | github.com/treporting/tools/review.IndexAt | yes | e5d8640d | yes | 4/5 |
| review/objects.go@NewIndex | github.com/treporting/tools/review.NewIndex | yes | d65a25b3 | yes | 20/22 |
| review/objects.go@ObjectIn | github.com/treporting/tools/review.ObjectIn | yes | 2cbf94f7 | yes | 12/13 |
| review/objects.go@RepoAt | github.com/treporting/tools/review.RepoAt | yes | 77d2e32a | no | 0/5 |
| review/objects.go@commonFolders | github.com/treporting/tools/review.commonFolders | yes | d06ac4c1 | yes | 5/5 |
| uplink/check/check.go@Check | - | no | - | yes | 1/1 |
| uplink/check/check.go@CheckWithAbove | github.com/treporting/tools/uplink/check.CheckWithAbove | yes | 667b53dd | yes | 18/18 |
| uplink/check/check.go@Finding.String | github.com/treporting/tools/uplink/check.Finding.String | yes | 11eecd84 | yes | 1/1 |
| uplink/check/check.go@loader | - | no | - | yes | 8/9 |
| uplink/check/check.go@objects | github.com/treporting/tools/uplink/check.objects | yes | aace180d | yes | 11/12 |
| uplink/check/check.go@status | github.com/treporting/tools/uplink/check.status | yes | e526017a | yes | 14/18 |
| uplink/check/parts.go@Parts | github.com/treporting/tools/uplink/check.Parts | yes | 5d40fce6 | yes | 23/26 |
| uplink/check/parts.go@name | github.com/treporting/tools/uplink/check.name | yes | f9f762ff | yes | 1/1 |
| uplink/cmd/uplinkcheck/main.go@fail | main.fail | yes | 0491b040 | no | 0/2 |
| uplink/cmd/uplinkcheck/main.go@main | main.main | yes | 6cf45306 | no | 0/23 |

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| 1.1.1 | 2026-11-20 | - | #37 | objectcode 2b5bcb5 |
| 2.1.1 | 2026-11-24 | T048.0 | #65 | objectcode 34ce1b3 |
