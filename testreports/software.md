# Test report for subsystem test (SwVal)

What the tests of the subsystem test (SwVal) specify, implement and show in a run. It is
written by the tool testreport (TOOL_34) and not edited by hand.

## Test environment

- Label: T045
- Commit: c811e2f
- Date: 2026-11-05
- Operating system: linux/amd64
- Runner: GitHub Actions 1000005325
- Go go1.25.14

## Original data

| File | Made by | Stored At |
|---|---|---|
| tools-integration.xml | `cd tools && gotestsum --junitfile tools-integration.xml -- -race -timeout 30m -tags integration ./...` | - |
| tools-unit.xml | `cd tools && gotestsum --junitfile tools-unit.xml -- -race -timeout 30m ./...` | - |

## Test specification

| Requirement | Test case | Text | Status |
|---|---|---|---|
| SR_01@1.1.3 | TC_01@1.1.1 | a result that names every criterion of its criteria file once, and no other: reviewcheck reports no criterion for it | approved |
| SR_01@1.1.3 | TC_02@1.1.1 | a result that leaves out a criterion of its criteria file: reviewcheck reports that criterion as left out | approved |
| SR_01@1.1.3 | TC_03@1.1.1 | a result that names a criterion its criteria file does not hold: reviewcheck reports that criterion as not in the file | approved |
| SR_02@1.1.3 | TC_04@1.1.1 | a result with the verdict passed and one criterion failed: reviewcheck reports the verdict | approved |
| SR_02@1.1.3 | TC_05@1.1.1 | a result with the verdict failed and every criterion passed: reviewcheck reports the verdict | approved |
| SR_02@1.1.3 | TC_06@1.1.1 | a result with the verdict passed and every criterion passed: reviewcheck reports no verdict | approved |
| SR_02@1.1.3 | TC_07@1.1.1 | a result with the verdict failed and one criterion failed: reviewcheck reports no verdict | approved |
| SR_03@1.1.3 | TC_08@1.1.1 | a failed criterion without a note: reviewcheck reports the criterion | approved |
| SR_03@1.1.3 | TC_09@1.1.1 | a failed criterion with a note: reviewcheck reports no missing note | approved |
| SR_04@1.1.3 | TC_10@1.1.1 | a result whose revision names no commit of the repo: reviewcheck reports the revision | approved |
| SR_04@1.1.3 | TC_11@1.1.1 | a result whose revision names a commit of the repo: reviewcheck reports no revision | approved |
| SR_05@2.1.1 | TC_12@1.1.1 | the reviewer is the author of the last commit, not a merge, that changed the object reviewed: reviewcheck reports the result | approved |
| SR_05@2.1.1 | TC_13@1.1.1 | the reviewer is a co-author, by a Co-Authored-By line, of the last commit, not a merge, that changed the object reviewed: reviewcheck reports the result | approved |
| SR_05@2.1.1 | TC_14@1.1.1 | the reviewer only merged the commit that changed the object reviewed, and another person wrote it: reviewcheck reports no result | approved |
| SR_05@2.1.1 | TC_15@1.1.1 | the reviewer is the author of a later commit that changed another object of the same file, not the object reviewed: reviewcheck reports no result | approved |
| SR_06@1.1.3 | TC_16@1.1.1 | a result with a criterion todo and the verdict passed: reviewcheck reports the value todo, and the verdict passed, as that criterion did not pass | approved |
| SR_06@1.1.3 | TC_17@1.1.1 | a result with the verdict todo: reviewcheck reports the value of the verdict | approved |
| SR_06@1.1.3 | TC_18@1.1.1 | a result whose criteria and verdict each hold passed or failed: reviewcheck reports no value | approved |
| SR_07@1.1.3 | TC_19@1.1.1 | a review with two results whose criteria revision names a commit that did not change its criteria file: reviewcheck reports the criteria revision once | approved |
| SR_07@1.1.3 | TC_20@1.1.1 | a review whose criteria revision names a commit that changed its criteria file: reviewcheck reports no criteria revision | approved |
| SR_08@1.1.1 | TC_21@1.1.1 | the reviewer is the author of the last commit, not a merge, that changed the document reviewed: reviewcheck reports the result | approved |
| SR_08@1.1.1 | TC_22@1.1.1 | the reviewer merged two pull requests that changed the document reviewed, and another person wrote every line: reviewcheck reports no result | approved |
| SR_11@2.1.3 | TC_23@1.1.1 | reviewskeleton with two ids, the scope sr-objects and a reviewer: it writes one file reviews/results/<yyyymmddThhmmss>_sr-objects.yaml, created_by is the reviewer, and created_at is the time of its name | approved |
| SR_11@2.1.3 | TC_24@1.1.1 | reviewskeleton with a document: it writes a skeleton whose result names the document | approved |
| SR_12@1.1.3 | TC_25@1.1.1 | reviewskeleton with ids of the level r5, and a later commit that changed another file: the skeleton names reviews/criteria-r5.yaml, and as criteria revision the last commit that changed that file | approved |
| SR_13@1.1.3 | TC_26@1.1.1 | reviewskeleton on a commit: the revision of each result is that commit | approved |
| SR_14@1.1.3 | TC_27@1.1.1 | reviewskeleton with two ids: it writes two results, each with every criterion of the criteria file, and each criterion and the verdict todo | approved |
| SR_15@1.1.3 | TC_28@1.1.1 | reviewskeleton with an id that no file holds: it writes nothing and reports the id | approved |
| SR_15@1.1.3 | TC_29@1.1.1 | reviewskeleton with a document that is not there: it writes nothing and reports the document | approved |
| SR_15@1.1.3 | TC_30@1.1.1 | reviewskeleton with two ids of different levels: it writes nothing and reports the levels | approved |
| SR_16@1.1.3 | TC_31@2.1.1 | a result with findings on two criteria in the order of the criteria file, one criterion not in it, the verdict and a revision that is no commit, checked twice: reviewcheck reports them in that order in both runs | approved |
| SR_17@2.1.3 | TC_32@1.1.1 | reviewskeleton with ids of the level r5: the skeleton has schema schemas/review-result.schema.json, type review and level r5 | approved |
| SR_83@1.1.2 | TC_33@1.1.1 | reviewskeleton: the skeleton has object_version 1.1.1 and status open | approved |
| SR_84@1.1.2 | TC_34@1.1.1 | a result with the status completed and no completed_at: reviewcheck reports the result | approved |
| SR_84@1.1.2 | TC_35@1.1.1 | a result with the status completed and a completed_at one second before its created_at: reviewcheck reports the result | approved |
| SR_84@1.1.2 | TC_36@1.1.1 | a result with the status completed and a completed_at equal to its created_at: reviewcheck reports no completed_at | approved |
| SR_84@1.1.2 | TC_37@1.1.1 | a result with the status open and no completed_at: reviewcheck reports no completed_at | approved |
| SR_75@1.1.2 | TC_38@1.1.1 | an approved requirement whose kind has a criteria file, and no review result names it: approvalcheck reports it | approved |
| SR_75@1.1.2 | TC_39@1.1.1 | an approved document with the verification method review, and no review result names it: approvalcheck reports it | approved |
| SR_75@1.1.2 | TC_40@1.1.1 | an approved requirement that a completed review result with the verdict passed names: approvalcheck reports no missing review | approved |
| SR_75@1.1.2 | TC_41@1.1.1 | an approved object whose kind has no criteria file, and no review result names it: approvalcheck reports no missing review | approved |
| SR_76@1.1.2 | TC_42@1.1.1 | an approved requirement whose last review result, by created_at, failed, and an earlier one passed: approvalcheck reports it | approved |
| SR_76@1.1.2 | TC_43@1.1.1 | an approved requirement whose last review result, by created_at, passed, and an earlier one failed: approvalcheck reports no failed review | approved |
| SR_77@1.1.2 | TC_44@1.1.1 | an approved requirement whose object_text changed after the revision of its last review result: approvalcheck reports it | approved |
| SR_77@1.1.2 | TC_45@1.1.1 | an approved requirement whose source_ids alone changed after the revision of its last review result: approvalcheck reports no change of content | approved |
| SR_77@1.1.2 | TC_46@1.1.1 | an approved document whose file changed after the revision of its last review result: approvalcheck reports it | approved |
| SR_78@1.1.2 | TC_47@1.1.1 | an approved requirement whose object_text differs from the first commit in which it was approved with its x: approvalcheck reports it | approved |
| SR_78@1.1.2 | TC_48@1.1.1 | an approved requirement whose source_ids alone differ from the first commit in which it was approved with its x: approvalcheck reports no attribute for x | approved |
| SR_79@2.1.2 | TC_49@1.1.1 | an approved document whose object_hash is not the hash of its file: approvalcheck reports it | approved |
| SR_79@2.1.2 | TC_50@1.1.1 | a completed document whose object_hash is not the hash of its file: approvalcheck reports it | approved |
| SR_79@2.1.2 | TC_51@1.1.1 | an approved document whose object_hash is the hash of its file: approvalcheck reports no hash | approved |
| SR_80@2.1.3 | TC_52@1.1.1 | a base commit given, one file of requirements changed against it, and an approved object without a review in an unchanged file: approvalcheck checks the objects of the changed file and reports nothing of the unchanged one | approved |
| SR_80@2.1.3 | TC_53@1.1.1 | a base commit given, and only the entry of a document in documents.yaml changed against it: approvalcheck checks that document | approved |
| SR_80@2.1.3 | TC_54@1.1.1 | no base commit given: approvalcheck checks every approved object and every completed document | approved |
| SR_85@1.1.1 | TC_55@1.1.1 | an approved requirement whose last review result passed and has the status open: approvalcheck reports it | approved |
| SR_85@1.1.1 | TC_56@1.1.1 | an approved requirement whose last review result passed and has the status completed: approvalcheck reports no open review | approved |
| SR_86@2.1.1 | TC_57@1.1.1 | a file git knows of the type testcases that holds TC_01: the id TC_01 is found in that file | approved |
| SR_86@2.1.1 | TC_58@1.1.1 | a review result that holds RR_01, and a file of criteria without objects: RR_01 is found in no file, and neither file is a file of objects | approved |
| SR_32@1.1.1 | TC_59@1.1.1 | a requirement whose source_ids name an object at another x than it has: linkcheck reports the link | approved |
| SR_32@1.1.1 | TC_60@1.1.1 | a requirement whose source_ids name an object at its x, with another y and z: linkcheck reports no link | approved |
| SR_32@1.1.1 | TC_61@1.1.1 | an element whose allocated names a requirement at another x than it has: linkcheck reports the entry | approved |
| SR_32@1.1.1 | TC_62@1.1.1 | a control whose controls names a user need at another x than it has: linkcheck reports the entry | approved |
| SR_32@1.1.1 | TC_63@1.1.1 | a requirement whose source_ids name a regulatory requirement at another x than the file of its module holds: linkcheck reports the link | approved |
| SR_33@2.1.1 | TC_64@1.1.1 | an @implements whose id names a requirement at another x than it has: linkcheck reports the annotation | approved |
| SR_33@2.1.1 | TC_65@1.1.1 | a @verifies whose id names a requirement at its x, with another z: linkcheck reports no annotation | approved |
| SR_33@2.1.1 | TC_66@1.1.1 | a @refines in the $comment of a schema whose id names an element at another x than it has: linkcheck reports the annotation | approved |
| SR_87@1.1.1 | TC_67@1.1.1 | a test case whose target_ids name a requirement at another x than it has: linkcheck reports the target | approved |
| SR_87@1.1.1 | TC_68@1.1.1 | a test case whose target_ids name a requirement at its x, with another y: linkcheck reports no target | approved |
| SR_24@3.1.1 | TC_69@2.1.1 | an @implements in a file that is code, whose id no object of the base commit holds: annotationcheck reports the id | approved |
| SR_24@3.1.1 | TC_70@2.1.1 | a @verifies in a file that is code, whose id a requirement of the base commit holds: annotationcheck reports no id | approved |
| SR_24@3.1.1 | TC_71@2.1.1 | an @implements in a file that is no code, whose id no object of the base commit holds: annotationcheck reports no id | approved |
| SR_25@2.2.1 | TC_72@1.1.1 | an @implements whose id has the status draft in the base commit: annotationcheck reports the id | approved |
| SR_25@2.2.1 | TC_73@1.1.1 | an @implements whose id is approved in the pull request and draft in the base commit: annotationcheck reports the id | approved |
| SR_25@2.2.1 | TC_74@1.1.1 | a @refines whose id is approved in the base commit: annotationcheck reports no id | approved |
| SR_27@1.2.1 | TC_75@1.1.1 | a file git knows whose kind is set neither as code nor as no code: annotationcheck reports the file | approved |
| SR_27@1.2.1 | TC_76@1.1.1 | a file git knows whose kind is set as code: annotationcheck reports no kind | approved |
| SR_30@2.2.1 | TC_77@1.1.1 | an @implements that names SR_01 without a version: annotationcheck reports the annotation | approved |
| SR_30@2.2.1 | TC_78@1.1.1 | a @verifies that names SR_01@1.1, a version without z: annotationcheck reports the annotation | approved |
| SR_30@2.2.1 | TC_79@1.1.1 | a @refines that names ARC_22@2.1.1: annotationcheck reports no missing version | approved |
| SR_35@1.1.1 | TC_80@1.1.1 | a schema whose $comment holds a @refines with an id not in the requirements of the base commit: annotationcheck reports the id | approved |
| SR_35@1.1.1 | TC_81@1.1.1 | a schema whose $comment holds a @refines with an id in the requirements of the base commit: annotationcheck reports no id | approved |
| SR_21@1.1.1 | TC_82@1.1.1 | a ruleset on GitHub that equals .github/rulesets/main.json in target, enforcement, conditions and rules, with its own id, times, links and bypass_actors: rulesetcheck reports no difference | approved |
| SR_21@1.1.1 | TC_83@1.1.1 | a required check in .github/rulesets/main.json that the ruleset on GitHub lacks: rulesetcheck reports the difference in the rules | approved |
| SR_21@1.1.1 | TC_84@1.1.1 | a rule on GitHub that .github/rulesets/main.json does not hold: rulesetcheck reports the difference in the rules | approved |
| SR_21@1.1.1 | TC_85@1.1.1 | a ruleset on GitHub whose enforcement differs from .github/rulesets/main.json: rulesetcheck reports the difference in the enforcement | approved |
| SR_22@1.1.1 | TC_86@1.1.1 | the ruleset cannot be read from GitHub: rulesetcheck reports it as a finding | approved |
| SR_22@1.1.1 | TC_87@1.1.1 | GitHub has no ruleset with the name in .github/rulesets/main.json: rulesetcheck reports it as a finding | approved |
| SR_28@1.1.1 | TC_88@1.1.1 | an approved requirement whose source_ids name an uplink that is approved: uplinkcheck reports no uplink | approved |
| SR_28@1.1.1 | TC_89@1.1.1 | an approved requirement whose source_ids name an uplink with the status draft: uplinkcheck reports the uplink | approved |
| SR_28@1.1.1 | TC_90@1.1.1 | a requirement with the status draft whose source_ids name an uplink with the status draft: uplinkcheck reports no uplink | approved |
| SR_18@1.1.1 | TC_91@1.1.1 | a YAML file of our data that matches the schema its field schema names: schemacheck reports no place | approved |
| SR_18@1.1.1 | TC_92@2.1.1 | a YAML file of our data with two fields whose values its schema does not allow: schemacheck reports each place, with the file and the path within it | approved |
| SR_18@1.1.1 | TC_93@1.1.1 | a YAML file of our data with a field its schema does not know: schemacheck reports the place | approved |
| SR_18@1.1.1 | TC_94@1.1.1 | a workflow of GitHub in .github/ that names no schema: schemacheck reports nothing of it | approved |
| SR_19@1.1.1 | TC_95@1.1.1 | a YAML file of our data that names no schema: schemacheck reports the file | approved |
| SR_19@1.1.1 | TC_96@1.1.1 | a YAML file of our data whose field schema names a schema that is not there: schemacheck reports the file | approved |
| SR_68@1.1.1 | TC_97@2.1.1 | a tag whose documents.yaml names README.md and requirements/tool-requirements.md: pdf writes README.pdf and requirements/tool-requirements.pdf | approved |
| SR_68@1.1.1 | TC_98@1.1.1 | a Markdown file that documents.yaml of the tag does not name: pdf writes no PDF of it | approved |
| SR_71@1.1.1 | TC_99@1.1.1 | devops.md at 2.1.1, approved in documents.yaml of the tag, and a later commit that changed another file: the footer of its PDF shows the path, 2.1.1, the last commit that changed devops.md, and approved | approved |
| SR_72@1.1.1 | TC_100@1.1.1 | a tag whose documents.yaml names two documents: index.pdf lists both, each with its version, its revision and its status | approved |
| SR_73@1.1.1 | TC_101@1.1.1 | pdf run on the tag v1.0.0 of the repo specs-and-deploy: it packs the PDFs and index.pdf in specs-and-deploy-v1.0.0-documents.zip | approved |
| SR_36@2.1.1 | TC_102@1.1.1 | a Go module that go list -deps names for the commands: sbom lists it with its version, the licence of its LICENSE file in the module cache, and the PURL pkg:golang/<path>@<version> | approved |
| SR_36@2.1.1 | TC_103@1.1.1 | the package of the subsystem itself among what go list -deps names: sbom lists it not as a Go module | approved |
| SR_37@2.1.1 | TC_104@1.1.1 | a package of package-lock.json not marked dev: sbom lists it with its version, its licence and the PURL pkg:npm/<name>@<version> | approved |
| SR_37@2.1.1 | TC_105@1.1.1 | a package of package-lock.json marked dev: sbom does not list it | approved |
| SR_38@1.2.1 | TC_106@1.1.1 | a package of /lib/apk/db/installed in the last stage of an image: sbom lists it with its version, its licence and the PURL pkg:apk/alpine/<name>@<version> | approved |
| SR_38@1.2.1 | TC_107@1.1.1 | a package installed only in an earlier stage of the build of an image: sbom does not list it | obsolete |
| SR_39@2.1.1 | TC_108@2.1.1 | sbom run in the repo of a subsystem: it writes sbom.md in the folder sbom/ of the repo | approved |
| SR_40@3.3.1 | TC_109@1.1.1 | sbom.md at 2.1.1 in documents.yaml, the baseline v1.2.0 and the ticket 42: sbom adds a row with 3.1.1, the day it runs, v1.2.0, #42 and the tool with its version | approved |
| SR_40@3.3.1 | TC_110@1.1.1 | no version of sbom.md in documents.yaml: sbom adds a row with the version 1.1.1 | approved |
| SR_41@2.1.1 | TC_111@2.1.1 | three packages: sbom lists them in the section Packages, in one table with the columns PURL, name, version and licence, sorted by PURL | approved |
| SR_42@1.2.1 | TC_112@1.1.1 | sbom writes sbom.md: it ends with the section Open, then the section Change history | approved |
| SR_43@1.1.1 | TC_113@1.1.1 | a package whose LICENSE file is the MIT licence: sbom names its licence MIT | approved |
| SR_43@1.1.1 | TC_114@1.1.1 | a package with two licences, Apache-2.0 and MIT: sbom names both joined by AND | approved |
| SR_43@1.1.1 | TC_115@1.1.1 | a package whose licence text is not recognized: sbom names its licence unknown | approved |
| SR_44@1.1.1 | TC_116@1.1.1 | commands built with Go 1.23.4: sbom lists the standard library with the version 1.23.4, its licence and the PURL pkg:golang/stdlib@1.23.4 | approved |
| SR_45@1.1.1 | TC_117@1.1.1 | a package whose licence is unknown: sbom lists its PURL in the section Open | approved |
| SR_45@1.1.1 | TC_118@1.1.1 | no package whose licence is unknown: sbom writes <empty> in the section Open | approved |
| SR_46@1.1.1 | TC_119@1.1.1 | the repo specs-and-deploy at the baseline v1.2.0: sbom lists it with the version v1.2.0, the licence of the LICENSE file in its root, and the PURL pkg:github/treporting/specs-and-deploy@v1.2.0 | approved |
| SR_59@2.1.1 | TC_120@1.1.1 | an entry of sbom.md in documents.yaml: sbom sets its version to that of the new row, its object_hash to git hash-object of sbom.md, its status to completed and its verification method to [generated] | approved |
| SR_59@2.1.1 | TC_121@1.1.1 | no entry of sbom.md in documents.yaml: sbom adds one of type document, with the version, the hash, the status completed and the verification method [generated] | approved |
| SR_66@1.1.1 | TC_122@1.1.1 | sbom -check where the packages equal the table of sbom.md: it passes and writes nothing | approved |
| SR_66@1.1.1 | TC_123@1.1.1 | sbom -check where one package has another version than the table of sbom.md: it fails, names that package, and writes nothing | approved |
| SR_47@1.1.1 | TC_124@1.1.1 | a package with a PURL of the type golang and the licence MIT: licencecheck passes it | approved |
| SR_47@1.1.1 | TC_125@1.1.1 | a package with a PURL of the type npm and the licence GPL-3.0: licencecheck fails on it | approved |
| SR_47@1.1.1 | TC_126@1.1.1 | a package with a PURL of the type github and the licence MIT AND GPL-3.0: licencecheck fails on it | approved |
| SR_47@1.1.1 | TC_127@1.1.1 | a package with a PURL of the type golang and the licence MIT OR GPL-3.0: licencecheck passes it | approved |
| SR_48@1.1.1 | TC_128@1.1.1 | a package with a PURL of the type apk and the licence unknown: licencecheck fails on it | approved |
| SR_48@1.1.1 | TC_129@1.1.1 | a package with a PURL of the type apk and the licence GPL-2.0-only: licencecheck passes it | approved |
| SR_49@1.1.1 | TC_130@1.1.1 | a package unknown in sbom.md whose PURL licences_assessed.yaml names with MIT and a reason: licencecheck takes MIT and passes it | approved |
| SR_49@1.1.1 | TC_131@1.1.1 | a package unknown in sbom.md whose PURL licences_assessed.yaml does not name: licencecheck takes unknown | approved |
| SR_50@1.1.1 | TC_132@1.1.1 | a package that fails: licencecheck names it with its PURL and its licence | approved |
| SR_68@1.1.1 | TC_133@1.1.1 | README.md changed in a commit after the tag: pdf run on the tag writes README.pdf with the text of README.md at the tag, not the later one | approved |
| SR_24@3.1.1 | TC_134@1.1.1 | an @implements in a file that is code, whose id is a test case of the base commit: annotationcheck reports no id | approved |
| SR_16@1.1.3 | TC_135@1.1.1 | a result with findings on two criteria in the order of the criteria file, one criterion not in it, the verdict and a reviewer who is the author of the object reviewed, checked twice: reviewcheck reports them in that order in both runs | approved |
| SR_88@2.1.1 | TC_136@1.2.1 | a junit.xml of go test and one of Vitest, a test case in `testcases/`, and a test procedure whose annotation is `@implements` of it: testreport reads the results of both files and links the test procedure to its test case | approved |
| SR_89@4.1.1 | TC_137@2.3.1 | any run of the test level r5tc: `testreports/software.md` has the head of a document and the sections Test environment, Original data, Test specification, Test implementation, Test execution, Deviations, Open and Change history, in this order | approved |
| SR_90@1.2.1 | TC_138@1.2.1 | a run with the label 1.2.0, tests of Go and of Vitest on Node.js: Test environment lists the label, the commit, the date, the operating system, the runner, and Go, Node.js and Vitest each with its version | approved |
| SR_91@2.1.1, SR_104@1.1.1 | TC_139@2.2.1 | two junit.xml given, testreport run without -store: Original data lists both with the command that made them and Stored At `-`; testreport copies neither | approved |
| SR_92@1.2.1 | TC_140@1.2.1 | a test case with the target SR_03@1.1.3: Test specification lists it with SR_03@1.1.3, its key with its version, its text and its status | approved |
| SR_93@1.2.1 | TC_141@1.2.1 | a test procedure with `@implements` of a test case, and one with only `@verifies` of a requirement: Test implementation lists the first with the requirement of its test case, the second with its requirement and the column Test case empty | approved |
| SR_94@1.2.1 | TC_142@1.2.1 | a test procedure in a junit.xml and one not in any: Test execution lists only the first, with its verdict | approved |
| SR_95@2.1.1 | TC_143@2.1.1 | results passed, failed and skipped, and an approved test case without a test procedure: Deviations lists the failed and the skipped one; Open holds the row of the test case, with its requirement and the status approved, and `-` as test procedure and as its status | approved |
| SR_95@2.1.1 | TC_144@1.3.1 | every result passed, and every test case and test procedure with a result: Deviations and Open each hold `<empty>` | approved |
| SR_96@2.1.1 | TC_145@2.1.1 | the label 1.2.0, the ticket 438, and the test report at 3.1.1 in `documents.yaml`: the new row of the Change history has the Version 4.1.1, the day of the run, Change 1.2.0, Reason #438 and By testreport with its version | approved |
| SR_96@2.1.1 | TC_146@2.1.1 | no label, no ticket, and no entry of the test report in `documents.yaml`: the new row has the Version 1.1.1, the day of the run, and Change and Reason empty | approved |
| SR_97@1.2.1 | TC_147@1.2.1 | a label without a ticket: testreport fails and writes nothing | approved |
| SR_98@2.1.1 | TC_148@2.1.1 | no entry of the test report in `documents.yaml`: testreport adds one of type `document` with the version of its new row, the hash of the test report written, the status `completed` and the verification method `[generated]` | approved |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | a go.mod in tools/, a package-lock.json in the root and the image docker.io/nginx:1.30, sbom run without -store: Original data lists three files, each with the command that made it and Stored At `-`; sbom writes nothing into sbom/ but sbom.md | approved |
| SR_101@2.1.1, SR_102@1.1.1 | TC_150@1.2.1 | the same repo, sbom run with -store at 2026-11-05 10:00:00: it writes into sbom/ `20261105T100000_golist-tools.txt`, `20261105T100000_root-package-lock.json` and `20261105T100000_apk-docker.io-nginx-1.30.txt`, and Original data names each with Stored At `./sbom/<file>` | approved |
| SR_102@1.1.1 | TC_151@1.1.1 | sbom run with -store at 2026-11-05 10:00:00: `sbom/20261105T100000.yaml` names every file of original data written, each with the command that made it | approved |
| SR_103@1.1.1 | TC_152@1.1.1 | two review results, one unchanged since the base commit and breaking a rule, one added after it and breaking none, reviewcheck run with -base: it checks only the one added and reports nothing | approved |
| SR_103@1.1.1 | TC_153@1.1.1 | the same two review results, reviewcheck run without -base: it checks both and reports the rule the unchanged one breaks | approved |
| SR_91@2.1.1, SR_104@1.1.1 | TC_154@1.1.1 | junit-go.xml and junit-vitest.xml given, test level r5tc, -dir tools, testreport run with -store at 2026-11-05 10:00:00: it copies them into tools/testreports/ as `20261105T100000_r5tc-junit-go.xml` and `20261105T100000_r5tc-junit-vitest.xml`, and Original data names each with Stored At `./tools/testreports/<file>` | approved |
| SR_104@1.1.1 | TC_155@1.1.1 | the same run with -store: `tools/testreports/20261105T100000_r5tc.yaml` names both files copied, each with the command that made it | approved |
| SR_105@1.1.1 | TC_156@1.1.1 | go list answers the module gopkg.in/yaml.v3 v3.0.1 in the directory `/home/u/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1`, the module cache `/home/u/go/pkg/mod`; sbom run with -store: the stored golist file names `$GOMODCACHE/gopkg.in/yaml.v3@v3.0.1`, and holds no `/home/u` | approved |
| SR_69@1.1.2 | TC_157@1.1.1 | a document naming `SR_01@1.1.3` in a table and `a@1.com` in its text, made into a page as GitHub shows it: `a@1.com` is a link to mailto:, `SR_01@1.1.3` is text and no link | approved |
| SR_95@2.1.1 | TC_158@1.1.1 | a test procedure without a test result that implements an approved test case: Open holds the row of the test case, with its requirement and the status approved, and the test procedure with the status `no test result` | approved |
| SR_95@2.1.1 | TC_159@1.1.1 | an obsolete test case without a test procedure: Open holds its row, with the status obsolete | approved |
| SR_41@2.1.1, SR_101@2.1.1 | TC_160@1.1.1 | sbom run: the sections of sbom.md are Original data, Packages, Open and Change history, in this order | approved |
| SR_89@4.1.1, SR_106@1.1.1 | TC_161@1.2.1 | a run of the test level r5tc, levels.yaml naming r5tc the subsystem test (SwVal): the title of the test report is `Test report for subsystem test (SwVal)`, and its head names neither r5tc nor r5 | approved |
| SR_106@1.1.1 | TC_162@1.1.1 | a run without -levels, `<root>/levels.yaml` naming r5tc: testreport reads it; the same run with -levels naming another levels.yaml: testreport reads that one | approved |
| SR_106@1.1.1 | TC_163@1.1.1 | a run of the test level r5tc, levels.yaml naming no r5tc: testreport exits with an error that names r5tc, and writes no test report | approved |

## Test implementation

| Requirement | Test case | Test procedure | File |
|---|---|---|---|
| SR_25@2.2.1 | TC_73@1.1.1 | TestApprovedInTheSamePullRequestDoesNotCount | tools/annotation/check/check_integration_test.go |
| SR_24@3.1.1 | TC_71@2.1.1 | TestOnlyCodeIsRead | tools/annotation/check/check_test.go |
| SR_27@1.2.1 | TC_75@1.1.1 | TestAKindOnNeitherList | tools/annotation/check/check_test.go |
| SR_27@1.2.1 | TC_76@1.1.1 | TestAKindOnNeitherList | tools/annotation/check/check_test.go |
| SR_24@3.1.1 |  | TestOnlyFirstInALineOfComment | tools/annotation/check/check_test.go |
| SR_24@3.1.1 | TC_69@2.1.1 | TestAnIdNotInTheBaseCommit | tools/annotation/check/check_test.go |
| SR_25@2.2.1 | TC_72@1.1.1 | TestAnIdNotApprovedInTheBaseCommit | tools/annotation/check/check_test.go |
| SR_24@3.1.1 |  | TestSeveralIdsInOneAnnotation | tools/annotation/check/check_test.go |
| SR_30@2.2.1 |  | TestTheVersionIsReadApartFromTheId | tools/annotation/check/check_test.go |
| SR_30@2.2.1 | TC_77@1.1.1 | TestAnIdWithoutAVersion | tools/annotation/check/check_test.go |
| SR_24@3.1.1 | TC_70@2.1.1 | TestAVerifiesOfARequirementIsNotReported | tools/annotation/check/check_test.go |
| SR_24@3.1.1 | TC_134@1.1.1 | TestAnImplementsOfATestCaseIsNotReported | tools/annotation/check/check_test.go |
| SR_30@2.2.1 | TC_78@1.1.1 | TestAVersionWithoutZ | tools/annotation/check/check_test.go |
| SR_25@2.2.1 | TC_74@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | tools/annotation/check/check_test.go |
| SR_35@1.1.1 | TC_80@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | tools/annotation/check/check_test.go |
| SR_35@1.1.1 | TC_81@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | tools/annotation/check/check_test.go |
| SR_25@2.2.1 |  | TestARefinesNotApprovedInTheBaseCommit | tools/annotation/check/check_test.go |
| SR_30@2.2.1 |  | TestARefinesWithoutAVersion | tools/annotation/check/check_test.go |
| SR_30@2.2.1 | TC_79@1.1.1 | TestARefinesWithAVersion | tools/annotation/check/check_test.go |
| SR_49@1.1.1, SR_50@1.1.1 |  | TestCommandNamesWhatFailsAndExitsWithOne | tools/licence/check/check_integration_test.go |
| SR_47@1.1.1, SR_48@1.1.1 |  | TestCheckByTheKindOfPackage | tools/licence/check/check_test.go |
| SR_47@1.1.1 |  | TestFits | tools/licence/check/check_test.go |
| SR_47@1.1.1 | TC_125@1.1.1 | TestAssessedLicenceTakesThePlaceOfTheOneInSbomMD | tools/licence/check/check_test.go |
| SR_49@1.1.1 | TC_130@1.1.1 | TestAssessedLicenceTakesThePlaceOfTheOneInSbomMD | tools/licence/check/check_test.go |
| SR_49@1.1.1 |  | TestAssessedEntryNeedsPURLLicenceAndReason | tools/licence/check/check_test.go |
| SR_50@1.1.1 | TC_132@1.1.1 | TestFailureNamesPURLAndLicence | tools/licence/check/check_test.go |
| SR_47@1.1.1 | TC_124@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_47@1.1.1 | TC_125@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_47@1.1.1 | TC_126@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_47@1.1.1 | TC_127@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_48@1.1.1 | TC_128@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_48@1.1.1 | TC_129@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_49@1.1.1 | TC_131@1.1.1 | TestEachTypeOfPURLWithItsLicence | tools/licence/check/check_test.go |
| SR_32@1.1.1 | TC_59@1.1.1 | TestOnlyAnotherXIsReported | tools/link/check/check_test.go |
| SR_32@1.1.1 | TC_60@1.1.1 | TestOnlyAnotherXIsReported | tools/link/check/check_test.go |
| SR_32@1.1.1 | TC_61@1.1.1 | TestTheAllocationAndControlsAreCompared | tools/link/check/check_test.go |
| SR_32@1.1.1 | TC_62@1.1.1 | TestTheAllocationAndControlsAreCompared | tools/link/check/check_test.go |
| SR_87@1.1.1 | TC_67@1.1.1 | TestATargetOfATestCaseIsCompared | tools/link/check/check_test.go |
| SR_87@1.1.1 | TC_68@1.1.1 | TestATargetOfATestCaseIsCompared | tools/link/check/check_test.go |
| SR_32@1.1.1 |  | TestAMissingObjectIsNotCompared | tools/link/check/check_test.go |
| SR_33@2.1.1 | TC_64@1.1.1 | TestAnAnnotationWithAnotherXIsReported | tools/link/check/check_test.go |
| SR_33@2.1.1 | TC_65@1.1.1 | TestAnAnnotationWithAnotherXIsReported | tools/link/check/check_test.go |
| SR_33@2.1.1 | TC_66@1.1.1 | TestARefinesWithAnotherXIsReported | tools/link/check/check_test.go |
| SR_32@1.1.1 | TC_63@1.1.1 | TestALinkToAModuleIsCompared | tools/link/check/check_test.go |
| SR_71@1.1.1 | TC_99@1.1.1 | TestGitReadsTheRepoAtTheTag | tools/pdf/export/export_integration_test.go |
| SR_68@1.1.1 | TC_133@1.1.1 | TestGitReadsTheRepoAtTheTag | tools/pdf/export/export_integration_test.go |
| SR_74@1.1.1 |  | TestGitFailsWhereThereIsNoSuchRef | tools/pdf/export/export_integration_test.go |
| SR_68@1.1.1 |  | TestPrinterWritesAPDF | tools/pdf/export/export_integration_test.go |
| SR_68@1.1.1 | TC_97@2.1.1 | TestDocumentsOfDocumentsYAML | tools/pdf/export/export_test.go |
| SR_68@1.1.1 | TC_98@1.1.1 | TestDocumentsOfDocumentsYAML | tools/pdf/export/export_test.go |
| SR_68@1.1.1 | TC_97@2.1.1 | TestPDFNameIsThePathWithPDF | tools/pdf/export/export_test.go |
| SR_71@1.1.1 | TC_99@1.1.1 | TestFooterNamesPathVersionRevisionStatus | tools/pdf/export/export_test.go |
| SR_71@1.1.1 |  | TestFooterTemplateHoldsDatePagesAndDocument | tools/pdf/export/export_test.go |
| SR_72@1.1.1 | TC_100@1.1.1 | TestIndexListsEveryDocument | tools/pdf/export/export_test.go |
| SR_73@1.1.1 | TC_101@1.1.1 | TestZipNameNamesRepoAndRef | tools/pdf/export/export_test.go |
| SR_73@1.1.1 | TC_101@1.1.1 | TestZipHoldsEveryFile | tools/pdf/export/export_test.go |
| SR_69@1.1.2 | TC_157@1.1.1 | TestAddressLinkedAsGitHubDoes | tools/pdf/export/export_test.go |
| SR_75@1.1.2 | TC_38@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_75@1.1.2 | TC_40@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_76@1.1.2 | TC_42@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_77@1.1.2 | TC_44@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_77@1.1.2 | TC_46@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_78@1.1.2 | TC_47@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_79@2.1.2 | TC_49@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_80@2.1.3 | TC_54@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_85@1.1.1 | TC_56@1.1.1 | TestEveryApprovedObject | tools/review/approval/approval_integration_test.go |
| SR_80@2.1.3 | TC_52@1.1.1 | TestTheFilesChangedAgainstTheBase | tools/review/approval/approval_integration_test.go |
| SR_85@1.1.1 | TC_55@1.1.1 | TestAnOpenResultIsNoApproval | tools/review/approval/approval_integration_test.go |
| SR_79@2.1.2 | TC_50@1.1.1 | TestACompletedDocumentHasItsHash | tools/review/approval/approval_integration_test.go |
| SR_75@1.1.2 | TC_39@1.1.1 | TestADocumentWithoutAReview | tools/review/approval/approval_integration_test.go |
| SR_79@2.1.2 | TC_51@1.1.1 | TestADocumentWithoutAReview | tools/review/approval/approval_integration_test.go |
| SR_75@1.1.2 | TC_41@1.1.1 | TestAKindWithoutCriteria | tools/review/approval/approval_integration_test.go |
| SR_76@1.1.2 | TC_43@1.1.1 | TestALaterPassedResult | tools/review/approval/approval_integration_test.go |
| SR_77@1.1.2 | TC_45@1.1.1 | TestOnlyALinkChanged | tools/review/approval/approval_integration_test.go |
| SR_78@1.1.2 | TC_48@1.1.1 | TestOnlyALinkChanged | tools/review/approval/approval_integration_test.go |
| SR_80@2.1.3 | TC_53@1.1.1 | TestOnlyTheEntryOfADocumentChanged | tools/review/approval/approval_integration_test.go |
| SR_79@2.1.2 |  | TestBlobHashIsTheHashOfGit | tools/review/approval/approval_test.go |
| SR_78@1.1.1 |  | TestFindXReadsTheAttributesForXAlsoInsideOneOf | tools/review/approval/approval_test.go |
| SR_78@1.1.1 |  | TestXOfKeepsOnlyTheAttributesForX | tools/review/approval/approval_test.go |
| SR_80@2.1.3 |  | TestInFilesKeepsTheObjectsOfTheChangedFiles | tools/review/approval/approval_test.go |
| SR_05@2.1.1 | TC_15@1.1.1 | TestTheCheckOnAHistoryOfGit | tools/review/check/check_integration_test.go |
| SR_07@1.1.1 |  | TestTheCriteriaRevisionOnAHistoryOfGit | tools/review/check/check_integration_test.go |
| SR_08@1.1.1 | TC_22@1.1.1 | TestTwoMergesOfOneDocument | tools/review/check/check_integration_test.go |
| SR_103@1.1.1 | TC_152@1.1.1 | TestWithBaseOnlyTheResultsAdded | tools/review/check/check_integration_test.go |
| SR_103@1.1.1 | TC_153@1.1.1 | TestWithoutBaseEveryResult | tools/review/check/check_integration_test.go |
| SR_01@1.1.3 | TC_01@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_02@1.1.3 | TC_06@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_04@1.1.3 | TC_11@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_06@1.1.3 | TC_18@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_07@1.1.3 | TC_20@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_84@1.1.2 | TC_37@1.1.1 | TestAGoodResultHasNoFinding | tools/review/check/check_test.go |
| SR_01@1.1.3 | TC_02@1.1.1 | TestEveryCriterionOnce | tools/review/check/check_test.go |
| SR_01@1.1.3 | TC_03@1.1.1 | TestEveryCriterionOnce | tools/review/check/check_test.go |
| SR_02@1.1.3 | TC_04@1.1.1 | TestTheVerdictFollowsTheCriteria | tools/review/check/check_test.go |
| SR_02@1.1.3 | TC_05@1.1.1 | TestTheVerdictFollowsTheCriteria | tools/review/check/check_test.go |
| SR_02@1.1.3 | TC_07@1.1.1 | TestAFailedVerdictWithANote | tools/review/check/check_test.go |
| SR_03@1.1.3 | TC_09@1.1.1 | TestAFailedVerdictWithANote | tools/review/check/check_test.go |
| SR_03@1.1.3 | TC_08@1.1.1 | TestAFailedCriterionHasANote | tools/review/check/check_test.go |
| SR_04@1.1.3 | TC_10@1.1.1 | TestTheRevisionExists | tools/review/check/check_test.go |
| SR_05@2.1.1 | TC_12@1.1.1 | TestTheReviewerIsNotTheAuthor | tools/review/check/check_test.go |
| SR_05@2.1.1 | TC_13@1.1.1 | TestTheReviewerIsNotTheAuthor | tools/review/check/check_test.go |
| SR_05@2.1.1 | TC_12@1.1.1 | TestASoftwareRequirementOfTheToolsHasAnAuthor | tools/review/check/check_test.go |
| SR_06@1.1.3 | TC_16@1.1.1 | TestAValueOtherThanPassedOrFailed | tools/review/check/check_test.go |
| SR_07@1.1.3 | TC_19@1.1.1 | TestTheCriteriaRevisionChangedTheCriteria | tools/review/check/check_test.go |
| SR_08@1.1.1 | TC_21@1.1.1 | TestTheReviewerIsNotTheAuthorOfADocument | tools/review/check/check_test.go |
| SR_05@2.1.1 | TC_14@1.1.1 | TestAMergeIsNoAuthorship | tools/review/check/check_test.go |
| SR_06@1.1.3 | TC_17@1.1.1 | TestAnUnfilledSkeletonIsReported | tools/review/check/check_test.go |
| SR_16@1.1.3 | TC_31@2.1.1 | TestTheFindingsInTheSameOrderInEveryRun | tools/review/check/check_test.go |
| SR_16@1.1.3 | TC_135@1.1.1 | TestTheFindingsUpToTheReviewerInTheSameOrder | tools/review/check/check_test.go |
| SR_84@1.1.2 | TC_34@1.1.1 | TestACompletedReviewHasItsCompletedAt | tools/review/check/check_test.go |
| SR_84@1.1.2 | TC_35@1.1.1 | TestACompletedReviewHasItsCompletedAt | tools/review/check/check_test.go |
| SR_84@1.1.2 | TC_36@1.1.1 | TestACompletedReviewHasItsCompletedAt | tools/review/check/check_test.go |
| SR_86@2.1.1 | TC_57@1.1.1 | TestEveryFileThatHoldsObjectsIsIndexed | tools/review/objects_test.go |
| SR_86@2.1.1 | TC_58@1.1.1 | TestEveryFileThatHoldsObjectsIsIndexed | tools/review/objects_test.go |
| SR_11@2.1.1, SR_12@1.1.1, SR_13@1.1.1, SR_14@1.1.1 |  | TestASkeletonOnAHistoryOfGit | tools/review/skeleton/skeleton_integration_test.go |
| SR_11@2.1.3 | TC_23@1.1.1 | TestTheNameTheReviewerAndTheTime | tools/review/skeleton/skeleton_test.go |
| SR_11@2.1.3 | TC_24@1.1.1 | TestTheCriteriaOfTheLevel | tools/review/skeleton/skeleton_test.go |
| SR_12@1.1.3 | TC_25@1.1.1 | TestTheCriteriaOfTheLevel | tools/review/skeleton/skeleton_test.go |
| SR_13@1.1.3 | TC_26@1.1.1 | TestTheRevisionIsHead | tools/review/skeleton/skeleton_test.go |
| SR_14@1.1.3 | TC_27@1.1.1 | TestEveryCriterionAndTheVerdictAreTodo | tools/review/skeleton/skeleton_test.go |
| SR_15@1.1.3 | TC_28@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | tools/review/skeleton/skeleton_test.go |
| SR_15@1.1.3 | TC_29@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | tools/review/skeleton/skeleton_test.go |
| SR_15@1.1.3 | TC_30@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | tools/review/skeleton/skeleton_test.go |
| SR_17@2.1.3 | TC_32@1.1.1 | TestTheSkeletonNamesItsSchemaTypeAndLevel | tools/review/skeleton/skeleton_test.go |
| SR_83@1.1.2 | TC_33@1.1.1 | TestTheSkeletonIsAnOpenRecord | tools/review/skeleton/skeleton_test.go |
| SR_21@1.1.1 | TC_82@1.1.1 | TestTheSameRuleset | tools/ruleset/check/check_test.go |
| SR_21@1.1.1 | TC_83@1.1.1 | TestARequiredCheckMissingOnGitHub | tools/ruleset/check/check_test.go |
| SR_21@1.1.1 | TC_85@1.1.1 | TestAnotherEnforcementOnGitHub | tools/ruleset/check/check_test.go |
| SR_21@1.1.1 | TC_84@1.1.1 | TestARuleOnlyOnGitHub | tools/ruleset/check/check_test.go |
| SR_22@1.1.1 | TC_86@1.1.1 | TestARulesetThatCannotBeRead | tools/ruleset/check/check_test.go |
| SR_22@1.1.1 | TC_87@1.1.1 | TestNoRulesetOfThatName | tools/ruleset/check/check_test.go |
| SR_66@1.1.1 |  | TestTableAsSbomMDListsIt | tools/sbom/create/check_test.go |
| SR_66@1.1.1 |  | TestLabelOfTheLastRow | tools/sbom/create/check_test.go |
| SR_66@1.1.1 | TC_123@1.1.1 | TestDifferNamesEachPackageThatDiffers | tools/sbom/create/check_test.go |
| SR_36@2.1.1 |  | TestGoListOfTheTools | tools/sbom/create/create_integration_test.go |
| SR_44@1.1.1 |  | TestGoEnvNamesTheVersionOfGo | tools/sbom/create/create_integration_test.go |
| SR_39@2.1.1 | TC_108@2.1.1 | TestCommandWritesSbomMDInTheRootAndKeepsTheHistory | tools/sbom/create/create_integration_test.go |
| SR_66@1.1.1 | TC_122@1.1.1 | TestCheckFailsWhereSbomMDDoesNotListWhatIsDelivered | tools/sbom/create/create_integration_test.go |
| SR_66@1.1.1 | TC_123@1.1.1 | TestCheckFailsWhereSbomMDDoesNotListWhatIsDelivered | tools/sbom/create/create_integration_test.go |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestCommandWithoutStoreWritesOnlySbomMD | tools/sbom/create/create_integration_test.go |
| SR_101@2.1.1, SR_102@1.1.1 | TC_150@1.2.1 | TestCommandWithStoreWritesTheOriginalsAndTheirRecord | tools/sbom/create/create_integration_test.go |
| SR_102@1.1.1 | TC_151@1.1.1 | TestCommandWithStoreWritesTheOriginalsAndTheirRecord | tools/sbom/create/create_integration_test.go |
| SR_46@1.1.1 | TC_119@1.1.1 | TestOwnPackage | tools/sbom/create/create_test.go |
| SR_40@3.3.1 | TC_109@1.1.1 | TestRowOfTheChangeHistory | tools/sbom/create/create_test.go |
| SR_40@3.3.1 |  | TestHistoryKeepsTheRowsOfTheChangeHistoryOnly | tools/sbom/create/create_test.go |
| SR_41@2.1.1 | TC_111@2.1.1 | TestOneTableSortedByPURL | tools/sbom/create/create_test.go |
| SR_42@1.2.1 | TC_112@1.1.1 | TestEndsWithOpenAndChangeHistory | tools/sbom/create/create_test.go |
| SR_41@2.1.1, SR_101@2.1.1 | TC_160@1.1.1 | TestSectionsInOrder | tools/sbom/create/create_test.go |
| SR_45@1.1.1 | TC_117@1.1.1 | TestOpenListsEveryPackageOfUnknownLicence | tools/sbom/create/create_test.go |
| SR_45@1.1.1 | TC_118@1.1.1 | TestOpenIsEmptyWhereEveryLicenceIsKnown | tools/sbom/create/create_test.go |
| SR_41@1.2.1 |  | TestAValueDoesNotBreakTheTable | tools/sbom/create/create_test.go |
| SR_40@3.3.1 | TC_109@1.1.1 | TestVersionOfSbomMD | tools/sbom/create/documents_test.go |
| SR_40@3.3.1 | TC_109@1.1.1 | TestNextRaisesX | tools/sbom/create/documents_test.go |
| SR_40@3.3.1 | TC_110@1.1.1 | TestNextRaisesX | tools/sbom/create/documents_test.go |
| SR_59@2.1.1 | TC_120@1.1.1 | TestHashIsThatOfGitHashObject | tools/sbom/create/documents_test.go |
| SR_59@2.1.1 | TC_121@1.1.1 | TestEntryAddedWhereThereIsNone | tools/sbom/create/documents_test.go |
| SR_59@2.1.1 | TC_120@1.1.1 | TestEntrySetInPlace | tools/sbom/create/documents_test.go |
| SR_59@2.1.1 | TC_120@1.1.1 | TestEntryCompletedAndGenerated | tools/sbom/create/documents_test.go |
| SR_43@1.1.1 | TC_113@1.1.1 | TestLicenceByItsSPDXIdentifier | tools/sbom/create/licence_test.go |
| SR_43@1.1.1 | TC_114@1.1.1 | TestSeveralLicencesJoinedByAnd | tools/sbom/create/licence_test.go |
| SR_43@1.1.1 | TC_115@1.1.1 | TestLicenceNotRecognizedIsUnknown | tools/sbom/create/licence_test.go |
| SR_43@1.1.1 |  | TestGivenLicenceInTheFormOfSPDX | tools/sbom/create/licence_test.go |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestOriginalDataWithoutStore | tools/sbom/create/originals_test.go |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestOriginalDataEmpty | tools/sbom/create/originals_test.go |
| SR_101@2.1.1, SR_102@1.1.1 | TC_150@1.2.1 | TestStoreWritesEachFileIntoSbom | tools/sbom/create/originals_test.go |
| SR_102@1.1.1 | TC_151@1.1.1 | TestOriginalYAMLNamesEveryFileStored | tools/sbom/create/originals_test.go |
| SR_105@1.1.1 | TC_156@1.1.1 | TestStoredGoListNamesNoLocalPath | tools/sbom/create/originals_test.go |
| SR_36@2.1.1 | TC_102@1.1.1 | TestGoModulesWithVersionLicenceAndPURL | tools/sbom/create/sources_test.go |
| SR_36@2.1.1 | TC_103@1.1.1 | TestGoModulesWithVersionLicenceAndPURL | tools/sbom/create/sources_test.go |
| SR_36@2.1.1, SR_43@1.1.1 |  | TestGoModuleWithoutLicenceFileIsUnknown | tools/sbom/create/sources_test.go |
| SR_44@1.1.1 | TC_116@1.1.1 | TestStdlib | tools/sbom/create/sources_test.go |
| SR_37@2.1.1 | TC_104@1.1.1 | TestNpmPackagesWithoutDevAndWithoutTheSubsystem | tools/sbom/create/sources_test.go |
| SR_37@2.1.1 | TC_105@1.1.1 | TestNpmPackagesWithoutDevAndWithoutTheSubsystem | tools/sbom/create/sources_test.go |
| SR_37@2.1.1 |  | TestNpmLockThatIsNoJSONIsAnError | tools/sbom/create/sources_test.go |
| SR_38@1.2.1 | TC_106@1.1.1 | TestApkPackages | tools/sbom/create/sources_test.go |
| SR_36@2.1.1, SR_37@2.1.1 |  | TestDirsWithAsTheBuildStepsFindThem | tools/sbom/create/sources_test.go |
| SR_18@1.1.1 | TC_94@1.1.1 | TestTheFilesOfAHistoryOfGit | tools/schema/check/check_integration_test.go |
| SR_18@1.1.1 | TC_91@1.1.1 | TestEachPlaceThatDoesNotMatch | tools/schema/check/check_test.go |
| SR_18@1.1.1 | TC_92@2.1.1 | TestEachPlaceThatDoesNotMatch | tools/schema/check/check_test.go |
| SR_18@1.1.1 | TC_93@1.1.1 | TestAFieldTheSchemaDoesNotKnow | tools/schema/check/check_test.go |
| SR_19@1.1.1 | TC_95@1.1.1 | TestNoSchemaOrASchemaThatIsNotThere | tools/schema/check/check_test.go |
| SR_19@1.1.1 | TC_96@1.1.1 | TestNoSchemaOrASchemaThatIsNotThere | tools/schema/check/check_test.go |
| SR_89@4.1.1, SR_106@1.1.1 | TC_161@1.2.1 | TestTitleNamesTheTestLevel | tools/testreport/write/testreport_integration_test.go |
| SR_106@1.1.1 | TC_162@1.1.1 | TestLevelsFromTheRootOrFromLevels | tools/testreport/write/testreport_integration_test.go |
| SR_106@1.1.1 | TC_163@1.1.1 | TestLevelNotNamedFails | tools/testreport/write/testreport_integration_test.go |
| SR_91@2.1.1, SR_104@1.1.1 | TC_139@2.2.1 | TestCommandWithoutStoreCopiesNothing | tools/testreport/write/testreport_integration_test.go |
| SR_91@2.1.1, SR_104@1.1.1 | TC_154@1.1.1 | TestCommandWithStoreCopiesEachJUnit | tools/testreport/write/testreport_integration_test.go |
| SR_104@1.1.1 | TC_155@1.1.1 | TestCommandWithStoreCopiesEachJUnit | tools/testreport/write/testreport_integration_test.go |
| SR_88@2.1.1 | TC_136@1.2.1 | TestReadsBothJUnitAndLinksTheProcedureToItsTestCase | tools/testreport/write/write_test.go |
| SR_89@4.1.1 | TC_137@2.3.1 | TestSectionsInTheirOrder | tools/testreport/write/write_test.go |
| SR_90@1.2.1 | TC_138@1.2.1 | TestEnvironmentNamesEachRuntimeWithItsVersion | tools/testreport/write/write_test.go |
| SR_91@2.1.1, SR_104@1.1.1 | TC_139@2.2.1 | TestOriginalDataWithoutStore | tools/testreport/write/write_test.go |
| SR_92@1.2.1 | TC_140@1.2.1 | TestSpecificationListsTheTestCaseWithItsRequirement | tools/testreport/write/write_test.go |
| SR_93@1.2.1 | TC_141@1.2.1 | TestImplementationWithAndWithoutATestCase | tools/testreport/write/write_test.go |
| SR_94@1.2.1 | TC_142@1.2.1 | TestExecutionListsOnlyWhatRan | tools/testreport/write/write_test.go |
| SR_95@2.1.1 | TC_143@2.1.1 | TestDeviationsAndOpen | tools/testreport/write/write_test.go |
| SR_95@2.1.1 | TC_144@1.3.1 | TestDeviationsAndOpenEmpty | tools/testreport/write/write_test.go |
| SR_95@2.1.1 | TC_158@1.1.1 | TestOpenTestProcedureWithoutResult | tools/testreport/write/write_test.go |
| SR_95@2.1.1 | TC_159@1.1.1 | TestOpenObsoleteTestCase | tools/testreport/write/write_test.go |
| SR_96@2.1.1 | TC_145@2.1.1 | TestRowWithLabelAndTicket | tools/testreport/write/write_test.go |
| SR_96@2.1.1 | TC_146@2.1.1 | TestRowWithoutLabel | tools/testreport/write/write_test.go |
| SR_97@1.2.1 | TC_147@1.2.1 | TestLabelWithoutTicketFails | tools/testreport/write/write_test.go |
| SR_98@2.1.1 | TC_148@2.1.1 | TestEntryAddedToDocuments | tools/testreport/write/write_test.go |
| SR_28@1.1.1 | TC_88@1.1.1 | TestApprovedUplinksPass | tools/uplink/check/check_test.go |
| SR_28@1.1.1 | TC_89@1.1.1 | TestUplinkNotApprovedIsAFinding | tools/uplink/check/check_test.go |
| SR_28@1.1.1 | TC_90@1.1.1 | TestDraftIsNotChecked | tools/uplink/check/check_test.go |
| SR_28@1.1.1 |  | TestAFileIsReadByItsSections | tools/uplink/check/check_test.go |
| SR_28@1.1.1 |  | TestARegulatoryRequirementIsFoundInItsModule | tools/uplink/check/check_test.go |

## Test execution

| Requirement | Test case | Test procedure | Verdict |
|---|---|---|---|
| SR_25@2.2.1 | TC_73@1.1.1 | TestApprovedInTheSamePullRequestDoesNotCount | passed |
| SR_24@3.1.1 | TC_71@2.1.1 | TestOnlyCodeIsRead | passed |
| SR_27@1.2.1 | TC_75@1.1.1 | TestAKindOnNeitherList | passed |
| SR_27@1.2.1 | TC_76@1.1.1 | TestAKindOnNeitherList | passed |
| SR_24@3.1.1 |  | TestOnlyFirstInALineOfComment | passed |
| SR_24@3.1.1 | TC_69@2.1.1 | TestAnIdNotInTheBaseCommit | passed |
| SR_25@2.2.1 | TC_72@1.1.1 | TestAnIdNotApprovedInTheBaseCommit | passed |
| SR_24@3.1.1 |  | TestSeveralIdsInOneAnnotation | passed |
| SR_30@2.2.1 |  | TestTheVersionIsReadApartFromTheId | passed |
| SR_30@2.2.1 | TC_77@1.1.1 | TestAnIdWithoutAVersion | passed |
| SR_24@3.1.1 | TC_70@2.1.1 | TestAVerifiesOfARequirementIsNotReported | passed |
| SR_24@3.1.1 | TC_134@1.1.1 | TestAnImplementsOfATestCaseIsNotReported | passed |
| SR_30@2.2.1 | TC_78@1.1.1 | TestAVersionWithoutZ | passed |
| SR_25@2.2.1 | TC_74@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | passed |
| SR_35@1.1.1 | TC_80@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | passed |
| SR_35@1.1.1 | TC_81@1.1.1 | TestRefinesIsReadInTheCommentOfASchema | passed |
| SR_25@2.2.1 |  | TestARefinesNotApprovedInTheBaseCommit | passed |
| SR_30@2.2.1 |  | TestARefinesWithoutAVersion | passed |
| SR_30@2.2.1 | TC_79@1.1.1 | TestARefinesWithAVersion | passed |
| SR_49@1.1.1, SR_50@1.1.1 |  | TestCommandNamesWhatFailsAndExitsWithOne | passed |
| SR_47@1.1.1, SR_48@1.1.1 |  | TestCheckByTheKindOfPackage | passed |
| SR_47@1.1.1 |  | TestFits | passed |
| SR_47@1.1.1 | TC_125@1.1.1 | TestAssessedLicenceTakesThePlaceOfTheOneInSbomMD | passed |
| SR_49@1.1.1 | TC_130@1.1.1 | TestAssessedLicenceTakesThePlaceOfTheOneInSbomMD | passed |
| SR_49@1.1.1 |  | TestAssessedEntryNeedsPURLLicenceAndReason | passed |
| SR_50@1.1.1 | TC_132@1.1.1 | TestFailureNamesPURLAndLicence | passed |
| SR_47@1.1.1 | TC_124@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_47@1.1.1 | TC_125@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_47@1.1.1 | TC_126@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_47@1.1.1 | TC_127@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_48@1.1.1 | TC_128@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_48@1.1.1 | TC_129@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_49@1.1.1 | TC_131@1.1.1 | TestEachTypeOfPURLWithItsLicence | passed |
| SR_32@1.1.1 | TC_59@1.1.1 | TestOnlyAnotherXIsReported | passed |
| SR_32@1.1.1 | TC_60@1.1.1 | TestOnlyAnotherXIsReported | passed |
| SR_32@1.1.1 | TC_61@1.1.1 | TestTheAllocationAndControlsAreCompared | passed |
| SR_32@1.1.1 | TC_62@1.1.1 | TestTheAllocationAndControlsAreCompared | passed |
| SR_87@1.1.1 | TC_67@1.1.1 | TestATargetOfATestCaseIsCompared | passed |
| SR_87@1.1.1 | TC_68@1.1.1 | TestATargetOfATestCaseIsCompared | passed |
| SR_32@1.1.1 |  | TestAMissingObjectIsNotCompared | passed |
| SR_33@2.1.1 | TC_64@1.1.1 | TestAnAnnotationWithAnotherXIsReported | passed |
| SR_33@2.1.1 | TC_65@1.1.1 | TestAnAnnotationWithAnotherXIsReported | passed |
| SR_33@2.1.1 | TC_66@1.1.1 | TestARefinesWithAnotherXIsReported | passed |
| SR_32@1.1.1 | TC_63@1.1.1 | TestALinkToAModuleIsCompared | passed |
| SR_71@1.1.1 | TC_99@1.1.1 | TestGitReadsTheRepoAtTheTag | passed |
| SR_68@1.1.1 | TC_133@1.1.1 | TestGitReadsTheRepoAtTheTag | passed |
| SR_74@1.1.1 |  | TestGitFailsWhereThereIsNoSuchRef | passed |
| SR_68@1.1.1 |  | TestPrinterWritesAPDF | passed |
| SR_68@1.1.1 | TC_97@2.1.1 | TestDocumentsOfDocumentsYAML | passed |
| SR_68@1.1.1 | TC_98@1.1.1 | TestDocumentsOfDocumentsYAML | passed |
| SR_68@1.1.1 | TC_97@2.1.1 | TestPDFNameIsThePathWithPDF | passed |
| SR_71@1.1.1 | TC_99@1.1.1 | TestFooterNamesPathVersionRevisionStatus | passed |
| SR_71@1.1.1 |  | TestFooterTemplateHoldsDatePagesAndDocument | passed |
| SR_72@1.1.1 | TC_100@1.1.1 | TestIndexListsEveryDocument | passed |
| SR_73@1.1.1 | TC_101@1.1.1 | TestZipNameNamesRepoAndRef | passed |
| SR_73@1.1.1 | TC_101@1.1.1 | TestZipHoldsEveryFile | passed |
| SR_69@1.1.2 | TC_157@1.1.1 | TestAddressLinkedAsGitHubDoes | passed |
| SR_75@1.1.2 | TC_38@1.1.1 | TestEveryApprovedObject | passed |
| SR_75@1.1.2 | TC_40@1.1.1 | TestEveryApprovedObject | passed |
| SR_76@1.1.2 | TC_42@1.1.1 | TestEveryApprovedObject | passed |
| SR_77@1.1.2 | TC_44@1.1.1 | TestEveryApprovedObject | passed |
| SR_77@1.1.2 | TC_46@1.1.1 | TestEveryApprovedObject | passed |
| SR_78@1.1.2 | TC_47@1.1.1 | TestEveryApprovedObject | passed |
| SR_79@2.1.2 | TC_49@1.1.1 | TestEveryApprovedObject | passed |
| SR_80@2.1.3 | TC_54@1.1.1 | TestEveryApprovedObject | passed |
| SR_85@1.1.1 | TC_56@1.1.1 | TestEveryApprovedObject | passed |
| SR_80@2.1.3 | TC_52@1.1.1 | TestTheFilesChangedAgainstTheBase | passed |
| SR_85@1.1.1 | TC_55@1.1.1 | TestAnOpenResultIsNoApproval | passed |
| SR_79@2.1.2 | TC_50@1.1.1 | TestACompletedDocumentHasItsHash | passed |
| SR_75@1.1.2 | TC_39@1.1.1 | TestADocumentWithoutAReview | passed |
| SR_79@2.1.2 | TC_51@1.1.1 | TestADocumentWithoutAReview | passed |
| SR_75@1.1.2 | TC_41@1.1.1 | TestAKindWithoutCriteria | passed |
| SR_76@1.1.2 | TC_43@1.1.1 | TestALaterPassedResult | passed |
| SR_77@1.1.2 | TC_45@1.1.1 | TestOnlyALinkChanged | passed |
| SR_78@1.1.2 | TC_48@1.1.1 | TestOnlyALinkChanged | passed |
| SR_80@2.1.3 | TC_53@1.1.1 | TestOnlyTheEntryOfADocumentChanged | passed |
| SR_79@2.1.2 |  | TestBlobHashIsTheHashOfGit | passed |
| SR_78@1.1.1 |  | TestFindXReadsTheAttributesForXAlsoInsideOneOf | passed |
| SR_78@1.1.1 |  | TestXOfKeepsOnlyTheAttributesForX | passed |
| SR_80@2.1.3 |  | TestInFilesKeepsTheObjectsOfTheChangedFiles | passed |
| SR_05@2.1.1 | TC_15@1.1.1 | TestTheCheckOnAHistoryOfGit | passed |
| SR_07@1.1.1 |  | TestTheCriteriaRevisionOnAHistoryOfGit | passed |
| SR_08@1.1.1 | TC_22@1.1.1 | TestTwoMergesOfOneDocument | passed |
| SR_103@1.1.1 | TC_152@1.1.1 | TestWithBaseOnlyTheResultsAdded | passed |
| SR_103@1.1.1 | TC_153@1.1.1 | TestWithoutBaseEveryResult | passed |
| SR_01@1.1.3 | TC_01@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_02@1.1.3 | TC_06@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_04@1.1.3 | TC_11@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_06@1.1.3 | TC_18@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_07@1.1.3 | TC_20@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_84@1.1.2 | TC_37@1.1.1 | TestAGoodResultHasNoFinding | passed |
| SR_01@1.1.3 | TC_02@1.1.1 | TestEveryCriterionOnce | passed |
| SR_01@1.1.3 | TC_03@1.1.1 | TestEveryCriterionOnce | passed |
| SR_02@1.1.3 | TC_04@1.1.1 | TestTheVerdictFollowsTheCriteria | passed |
| SR_02@1.1.3 | TC_05@1.1.1 | TestTheVerdictFollowsTheCriteria | passed |
| SR_02@1.1.3 | TC_07@1.1.1 | TestAFailedVerdictWithANote | passed |
| SR_03@1.1.3 | TC_09@1.1.1 | TestAFailedVerdictWithANote | passed |
| SR_03@1.1.3 | TC_08@1.1.1 | TestAFailedCriterionHasANote | passed |
| SR_04@1.1.3 | TC_10@1.1.1 | TestTheRevisionExists | passed |
| SR_05@2.1.1 | TC_12@1.1.1 | TestTheReviewerIsNotTheAuthor | passed |
| SR_05@2.1.1 | TC_13@1.1.1 | TestTheReviewerIsNotTheAuthor | passed |
| SR_05@2.1.1 | TC_12@1.1.1 | TestASoftwareRequirementOfTheToolsHasAnAuthor | passed |
| SR_06@1.1.3 | TC_16@1.1.1 | TestAValueOtherThanPassedOrFailed | passed |
| SR_07@1.1.3 | TC_19@1.1.1 | TestTheCriteriaRevisionChangedTheCriteria | passed |
| SR_08@1.1.1 | TC_21@1.1.1 | TestTheReviewerIsNotTheAuthorOfADocument | passed |
| SR_05@2.1.1 | TC_14@1.1.1 | TestAMergeIsNoAuthorship | passed |
| SR_06@1.1.3 | TC_17@1.1.1 | TestAnUnfilledSkeletonIsReported | passed |
| SR_16@1.1.3 | TC_31@2.1.1 | TestTheFindingsInTheSameOrderInEveryRun | passed |
| SR_16@1.1.3 | TC_135@1.1.1 | TestTheFindingsUpToTheReviewerInTheSameOrder | passed |
| SR_84@1.1.2 | TC_34@1.1.1 | TestACompletedReviewHasItsCompletedAt | passed |
| SR_84@1.1.2 | TC_35@1.1.1 | TestACompletedReviewHasItsCompletedAt | passed |
| SR_84@1.1.2 | TC_36@1.1.1 | TestACompletedReviewHasItsCompletedAt | passed |
| SR_86@2.1.1 | TC_57@1.1.1 | TestEveryFileThatHoldsObjectsIsIndexed | passed |
| SR_86@2.1.1 | TC_58@1.1.1 | TestEveryFileThatHoldsObjectsIsIndexed | passed |
| SR_11@2.1.1, SR_12@1.1.1, SR_13@1.1.1, SR_14@1.1.1 |  | TestASkeletonOnAHistoryOfGit | passed |
| SR_11@2.1.3 | TC_23@1.1.1 | TestTheNameTheReviewerAndTheTime | passed |
| SR_11@2.1.3 | TC_24@1.1.1 | TestTheCriteriaOfTheLevel | passed |
| SR_12@1.1.3 | TC_25@1.1.1 | TestTheCriteriaOfTheLevel | passed |
| SR_13@1.1.3 | TC_26@1.1.1 | TestTheRevisionIsHead | passed |
| SR_14@1.1.3 | TC_27@1.1.1 | TestEveryCriterionAndTheVerdictAreTodo | passed |
| SR_15@1.1.3 | TC_28@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | passed |
| SR_15@1.1.3 | TC_29@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | passed |
| SR_15@1.1.3 | TC_30@1.1.1 | TestNothingWhereATargetIsMissingOrTheLevelsDiffer | passed |
| SR_17@2.1.3 | TC_32@1.1.1 | TestTheSkeletonNamesItsSchemaTypeAndLevel | passed |
| SR_83@1.1.2 | TC_33@1.1.1 | TestTheSkeletonIsAnOpenRecord | passed |
| SR_21@1.1.1 | TC_82@1.1.1 | TestTheSameRuleset | passed |
| SR_21@1.1.1 | TC_83@1.1.1 | TestARequiredCheckMissingOnGitHub | passed |
| SR_21@1.1.1 | TC_85@1.1.1 | TestAnotherEnforcementOnGitHub | passed |
| SR_21@1.1.1 | TC_84@1.1.1 | TestARuleOnlyOnGitHub | passed |
| SR_22@1.1.1 | TC_86@1.1.1 | TestARulesetThatCannotBeRead | passed |
| SR_22@1.1.1 | TC_87@1.1.1 | TestNoRulesetOfThatName | passed |
| SR_66@1.1.1 |  | TestTableAsSbomMDListsIt | passed |
| SR_66@1.1.1 |  | TestLabelOfTheLastRow | passed |
| SR_66@1.1.1 | TC_123@1.1.1 | TestDifferNamesEachPackageThatDiffers | passed |
| SR_36@2.1.1 |  | TestGoListOfTheTools | passed |
| SR_44@1.1.1 |  | TestGoEnvNamesTheVersionOfGo | passed |
| SR_39@2.1.1 | TC_108@2.1.1 | TestCommandWritesSbomMDInTheRootAndKeepsTheHistory | passed |
| SR_66@1.1.1 | TC_122@1.1.1 | TestCheckFailsWhereSbomMDDoesNotListWhatIsDelivered | passed |
| SR_66@1.1.1 | TC_123@1.1.1 | TestCheckFailsWhereSbomMDDoesNotListWhatIsDelivered | passed |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestCommandWithoutStoreWritesOnlySbomMD | passed |
| SR_101@2.1.1, SR_102@1.1.1 | TC_150@1.2.1 | TestCommandWithStoreWritesTheOriginalsAndTheirRecord | passed |
| SR_102@1.1.1 | TC_151@1.1.1 | TestCommandWithStoreWritesTheOriginalsAndTheirRecord | passed |
| SR_46@1.1.1 | TC_119@1.1.1 | TestOwnPackage | passed |
| SR_40@3.3.1 | TC_109@1.1.1 | TestRowOfTheChangeHistory | passed |
| SR_40@3.3.1 |  | TestHistoryKeepsTheRowsOfTheChangeHistoryOnly | passed |
| SR_41@2.1.1 | TC_111@2.1.1 | TestOneTableSortedByPURL | passed |
| SR_42@1.2.1 | TC_112@1.1.1 | TestEndsWithOpenAndChangeHistory | passed |
| SR_41@2.1.1, SR_101@2.1.1 | TC_160@1.1.1 | TestSectionsInOrder | passed |
| SR_45@1.1.1 | TC_117@1.1.1 | TestOpenListsEveryPackageOfUnknownLicence | passed |
| SR_45@1.1.1 | TC_118@1.1.1 | TestOpenIsEmptyWhereEveryLicenceIsKnown | passed |
| SR_41@1.2.1 |  | TestAValueDoesNotBreakTheTable | passed |
| SR_40@3.3.1 | TC_109@1.1.1 | TestVersionOfSbomMD | passed |
| SR_40@3.3.1 | TC_109@1.1.1 | TestNextRaisesX | passed |
| SR_40@3.3.1 | TC_110@1.1.1 | TestNextRaisesX | passed |
| SR_59@2.1.1 | TC_120@1.1.1 | TestHashIsThatOfGitHashObject | passed |
| SR_59@2.1.1 | TC_121@1.1.1 | TestEntryAddedWhereThereIsNone | passed |
| SR_59@2.1.1 | TC_120@1.1.1 | TestEntrySetInPlace | passed |
| SR_59@2.1.1 | TC_120@1.1.1 | TestEntryCompletedAndGenerated | passed |
| SR_43@1.1.1 | TC_113@1.1.1 | TestLicenceByItsSPDXIdentifier | passed |
| SR_43@1.1.1 | TC_114@1.1.1 | TestSeveralLicencesJoinedByAnd | passed |
| SR_43@1.1.1 | TC_115@1.1.1 | TestLicenceNotRecognizedIsUnknown | passed |
| SR_43@1.1.1 |  | TestGivenLicenceInTheFormOfSPDX | passed |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestOriginalDataWithoutStore | passed |
| SR_101@2.1.1, SR_102@1.1.1 | TC_149@1.2.1 | TestOriginalDataEmpty | passed |
| SR_101@2.1.1, SR_102@1.1.1 | TC_150@1.2.1 | TestStoreWritesEachFileIntoSbom | passed |
| SR_102@1.1.1 | TC_151@1.1.1 | TestOriginalYAMLNamesEveryFileStored | passed |
| SR_105@1.1.1 | TC_156@1.1.1 | TestStoredGoListNamesNoLocalPath | passed |
| SR_36@2.1.1 | TC_102@1.1.1 | TestGoModulesWithVersionLicenceAndPURL | passed |
| SR_36@2.1.1 | TC_103@1.1.1 | TestGoModulesWithVersionLicenceAndPURL | passed |
| SR_36@2.1.1, SR_43@1.1.1 |  | TestGoModuleWithoutLicenceFileIsUnknown | passed |
| SR_44@1.1.1 | TC_116@1.1.1 | TestStdlib | passed |
| SR_37@2.1.1 | TC_104@1.1.1 | TestNpmPackagesWithoutDevAndWithoutTheSubsystem | passed |
| SR_37@2.1.1 | TC_105@1.1.1 | TestNpmPackagesWithoutDevAndWithoutTheSubsystem | passed |
| SR_37@2.1.1 |  | TestNpmLockThatIsNoJSONIsAnError | passed |
| SR_38@1.2.1 | TC_106@1.1.1 | TestApkPackages | passed |
| SR_36@2.1.1, SR_37@2.1.1 |  | TestDirsWithAsTheBuildStepsFindThem | passed |
| SR_18@1.1.1 | TC_94@1.1.1 | TestTheFilesOfAHistoryOfGit | passed |
| SR_18@1.1.1 | TC_91@1.1.1 | TestEachPlaceThatDoesNotMatch | passed |
| SR_18@1.1.1 | TC_92@2.1.1 | TestEachPlaceThatDoesNotMatch | passed |
| SR_18@1.1.1 | TC_93@1.1.1 | TestAFieldTheSchemaDoesNotKnow | passed |
| SR_19@1.1.1 | TC_95@1.1.1 | TestNoSchemaOrASchemaThatIsNotThere | passed |
| SR_19@1.1.1 | TC_96@1.1.1 | TestNoSchemaOrASchemaThatIsNotThere | passed |
| SR_89@4.1.1, SR_106@1.1.1 | TC_161@1.2.1 | TestTitleNamesTheTestLevel | passed |
| SR_106@1.1.1 | TC_162@1.1.1 | TestLevelsFromTheRootOrFromLevels | passed |
| SR_106@1.1.1 | TC_163@1.1.1 | TestLevelNotNamedFails | passed |
| SR_91@2.1.1, SR_104@1.1.1 | TC_139@2.2.1 | TestCommandWithoutStoreCopiesNothing | passed |
| SR_91@2.1.1, SR_104@1.1.1 | TC_154@1.1.1 | TestCommandWithStoreCopiesEachJUnit | passed |
| SR_104@1.1.1 | TC_155@1.1.1 | TestCommandWithStoreCopiesEachJUnit | passed |
| SR_88@2.1.1 | TC_136@1.2.1 | TestReadsBothJUnitAndLinksTheProcedureToItsTestCase | passed |
| SR_89@4.1.1 | TC_137@2.3.1 | TestSectionsInTheirOrder | passed |
| SR_90@1.2.1 | TC_138@1.2.1 | TestEnvironmentNamesEachRuntimeWithItsVersion | passed |
| SR_91@2.1.1, SR_104@1.1.1 | TC_139@2.2.1 | TestOriginalDataWithoutStore | passed |
| SR_92@1.2.1 | TC_140@1.2.1 | TestSpecificationListsTheTestCaseWithItsRequirement | passed |
| SR_93@1.2.1 | TC_141@1.2.1 | TestImplementationWithAndWithoutATestCase | passed |
| SR_94@1.2.1 | TC_142@1.2.1 | TestExecutionListsOnlyWhatRan | passed |
| SR_95@2.1.1 | TC_143@2.1.1 | TestDeviationsAndOpen | passed |
| SR_95@2.1.1 | TC_144@1.3.1 | TestDeviationsAndOpenEmpty | passed |
| SR_95@2.1.1 | TC_158@1.1.1 | TestOpenTestProcedureWithoutResult | passed |
| SR_95@2.1.1 | TC_159@1.1.1 | TestOpenObsoleteTestCase | passed |
| SR_96@2.1.1 | TC_145@2.1.1 | TestRowWithLabelAndTicket | passed |
| SR_96@2.1.1 | TC_146@2.1.1 | TestRowWithoutLabel | passed |
| SR_97@1.2.1 | TC_147@1.2.1 | TestLabelWithoutTicketFails | passed |
| SR_98@2.1.1 | TC_148@2.1.1 | TestEntryAddedToDocuments | passed |
| SR_28@1.1.1 | TC_88@1.1.1 | TestApprovedUplinksPass | passed |
| SR_28@1.1.1 | TC_89@1.1.1 | TestUplinkNotApprovedIsAFinding | passed |
| SR_28@1.1.1 | TC_90@1.1.1 | TestDraftIsNotChecked | passed |
| SR_28@1.1.1 |  | TestAFileIsReadByItsSections | passed |
| SR_28@1.1.1 |  | TestARegulatoryRequirementIsFoundInItsModule | passed |

## Deviations

`<empty>`

## Open

**Table open** - every test case and test procedure without a test result

| Requirement | Test case | Test case status | Test procedure | Test procedure status |
|---|---|---|---|---|
| SR_38@1.2.1 | TC_107@1.1.1 | obsolete | - | - |

## Change history

| Version | Date | Change | Reason | By |
|---|---|---|---|---|
| 1.1.1 | 2026-11-05 | T045 | #556 | testreport c811e2f-dirty |
