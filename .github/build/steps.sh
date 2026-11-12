#!/usr/bin/env bash
# The build steps of every repo (TOOL_23 to TOOL_26), run by .github/workflows/build.yml in the root
# of the repo; the same script runs by hand.
#
#	steps.sh build|lint|test|cover|debt
#	steps.sh sbom [name=Dockerfile]...
#	steps.sh baseline <label> <ticket number> [name=Dockerfile]...
#
# A Go module is every directory with a go.mod, a JavaScript package every directory with a
# package-lock.json - neither below node_modules nor below a directory whose name starts with a dot.
# build, lint and test exit with 1 where a step fails. cover and debt report and exit with 0; they
# exit otherwise only where the measuring itself did not happen. sbom builds each image named from
# its Dockerfile, and exits with 1 where sbom.md does not list the packages delivered. baseline
# builds the images, writes sbom.md with the label and the ticket, and checks its licences; it exits
# with 1 where the licences do not fit. Where JUNIT names a directory, test writes the junit.xml of
# each run there, and baseline writes the test report of each test level from them.
#
# ABOVE names the root of the repo above - specs-and-deploy for the tools: baseline takes the schemas
# and levels.yaml from there. APPROVAL_BASE, where set, is the commit baseline checks the approvals
# against: in tools, the first commit, as every approval before it has its review in specs-and-deploy.
#
# It needs go, golangci-lint v2 and deadcode on the PATH, gotestsum where JUNIT is set, and node with
# npm where there is a JavaScript package.
set -uo pipefail

#	the root of tools; .github/build holds the configurations of the linters
steps="$(cd "$(dirname "$0")/../.." && pwd)"

#	@implements SR_51@1.1.1
modules() {
	find . -name go.mod -not -path '*/node_modules/*' -not -path '*/.*/*' | sed 's#/go.mod$##' | sort
}

#	@implements SR_52@1.1.1
packages() {
	find . -name package-lock.json -not -path '*/node_modules/*' -not -path '*/.*/*' | sed 's#/package-lock.json$##' | sort
}

failed=0

#	run <dir> <command>...: run a command in a directory, and note when it fails
run() {
	local dir="$1"
	shift
	echo "== $dir: $*"
	if ! (cd "$dir" && "$@"); then
		echo "FAIL $dir: $*"
		failed=1
	fi
}

#	@implements SR_51@1.1.1, SR_52@1.1.1
build() {
	for m in $(modules); do
		run "$m" go build ./...
	done
	for p in $(packages); do
		run "$p" npm ci --no-audit --no-fund
		run "$p" npm run build
	done
}

#	the packages eslint.config.js imports, installed beside it once
eslint_install() {
	[ -x "$steps/node_modules/.bin/eslint" ] && return
	(cd "$steps" && npm install --no-save --no-audit --no-fund \
		eslint@10.14.0 @eslint/js@10.0.1 globals@17.13.0) || exit 2
}

#	@implements SR_53@1.1.1, SR_54@1.1.1
lint() {
	for m in $(modules); do
		run "$m" golangci-lint run -c "$steps/.github/build/.golangci.yml" ./...
		run "$m" golangci-lint run -c "$steps/.github/build/.golangci.yml" --build-tags integration ./...
	done
	if [ -n "$(packages)" ]; then
		eslint_install
	fi
	for p in $(packages); do
		run "$p" "$steps/node_modules/.bin/eslint" -c "$steps/.github/build/eslint.config.js" .
	done
}

#	junit <module> <kind> <flags of go test>...: run the tests of a module with gotestsum, which writes
#	their junit.xml into $JUNIT as <module>-<kind>.xml, and the command that made it beside it, as
#	<module>-<kind>.cmd (SR_99)
junit() {
	local m="$1" kind="$2" name cd=""
	shift 2
	name="$(echo "${m#./}" | tr '/' '-')"
	if [ "$m" = . ]; then
		name=root
	else
		cd="cd ${m#./} && "
	fi
	name="$name-$kind.xml"
	mkdir -p "$JUNIT"
	echo "${cd}gotestsum --junitfile $name -- $*" >"$JUNIT/${name%.xml}.cmd"
	run "$m" gotestsum --junitfile "$JUNIT/$name" -- "$@"
}

#	@implements SR_55@1.1.1, SR_56@1.1.1
test() {
	for m in $(modules); do
		if [ -n "${JUNIT:-}" ]; then
			junit "$m" unit -race -timeout 30m ./...
			junit "$m" integration -race -timeout 30m -tags integration ./...
		else
			run "$m" go test -race -timeout 30m ./...
			run "$m" go test -race -timeout 30m -tags integration ./...
		fi
	done
	for p in $(packages); do
		run "$p" npm ci --no-audit --no-fund
		run "$p" npm test
	done
}

#	@implements SR_57@2.1.2
cover() {
	local profile
	profile="$(mktemp)"
	for m in $(modules); do
		echo "== $m: test coverage of MT and SwIT, with -coverpkg=./..."
		if ! (cd "$m" && go test -tags integration -coverpkg=./... -coverprofile="$profile" ./... >/dev/null); then
			echo "a test failed; the coverage is of the tests that ran"
		fi
		(cd "$m" && go tool cover -func="$profile") || exit 2
	done
	for p in $(packages); do
		echo "== $p: test coverage"
		if (cd "$p" && npm pkg get scripts.coverage | grep -qv '{}'); then
			(cd "$p" && npm ci --no-audit --no-fund && npm run coverage) || exit 2
		else
			echo "no script coverage in package.json"
		fi
	done
	rm -f "$profile"
}

#	@implements SR_58@2.1.2
debt() {
	for m in $(modules); do
		echo "== $m: what nothing uses, and a cyclomatic complexity of 15 and more"
		(cd "$m" && golangci-lint run -c "$steps/.github/build/debt.golangci.yml" --tests=false ./...)
		[ $? -le 1 ] || exit 2
		echo "== $m: what no program reaches (deadcode)"
		(cd "$m" && deadcode ./...) || exit 2
	done
}

#	tool <path>: build a command of the tools, such as sbom/cmd/sbom, and print
#	where it lies
tool() {
	local bin
	bin="$(mktemp -d)/$(basename "$1")"
	(cd "$steps" && GOTOOLCHAIN=auto go build -o "$bin" "./$1") >&2 || exit 2
	echo "$bin"
}

#	build_images <name=Dockerfile>...: build each image, and set image_flags to -image <name> each
build_images() {
	image_flags=()
	for image in "$@"; do
		run . docker build -t "${image%%=*}" -f "${image#*=}" .
		image_flags+=(-image "${image%%=*}")
	done
	[ "$failed" = 0 ] || exit 2
}

#	@implements SR_67@1.1.1
sbom() {
	local bin
	bin="$(tool sbom/cmd/sbom)"
	build_images "$@"
	run . "$bin" -check ${image_flags[@]+"${image_flags[@]}"}
}

#	@implements SR_63@1.1.1, SR_82@2.1.3
baseline() {
	local label="$1" ticket="$2" bin_sbom bin_licence bin_approval bin_testreport
	shift 2
	#	every tool before the first of them writes, so that none is built from a changed working
	#	tree and records its revision with -dirty
	bin_sbom="$(tool sbom/cmd/sbom)"
	bin_licence="$(tool licence/cmd/licencecheck)"
	bin_approval="$(tool review/cmd/approvalcheck)"
	bin_testreport="$(tool testreport/cmd/testreport)"
	build_images "$@"
	run . "$bin_sbom" -label "$label" -ticket "$ticket" ${image_flags[@]+"${image_flags[@]}"}
	[ "$failed" = 0 ] || return
	run . "$bin_licence" -root .
	#	every approved object of the head, after the licences - or those changed since APPROVAL_BASE
	run . "$bin_approval" -root . ${ABOVE:+-above "$ABOVE"} ${APPROVAL_BASE:+-base "$APPROVAL_BASE"}
	[ "$failed" = 0 ] || return
	if [ -n "${JUNIT:-}" ]; then
		testreports "$bin_testreport" "$label" "$ticket"
	fi
}

#	testreports <bin> <label> <ticket>: write with the testreport at <bin> the test report of each
#	test level that has a file of test cases - testcases/<file>.yaml with its level in its head -
#	from every junit.xml in $JUNIT
#
#	@implements SR_99@2.1.1
testreports() {
	local bin="$1" label="$2" ticket="$3" file dir level xml
	local flags=()
	for xml in "$JUNIT"/*.xml; do
		[ -e "$xml" ] || continue
		flags+=(-junit "$xml" -made-by "$(cat "${xml%.xml}.cmd")")
	done
	[ ${#flags[@]} -gt 0 ] || return
	flags+=(-runtime "Go $(go env GOVERSION)")
	for file in $(find . -path '*/testcases/*.yaml' -not -path '*/node_modules/*' -not -path '*/.*/*' | sort); do
		dir="$(dirname "$(dirname "$file")")"
		level="$(awk '/^level:/ {print $2; exit}' "$file")"
		run . "$bin" -root . -dir "${dir#./}" -level "$level" -label "$label" -ticket "$ticket" \
			-levels "${ABOVE:-.}/levels.yaml" "${flags[@]}"
	done
}

case "${1:-}" in
build | lint | test)
	"$1"
	exit "$failed"
	;;
sbom)
	shift
	sbom "$@"
	exit "$failed"
	;;
baseline)
	[ $# -ge 3 ] || {
		echo "usage: steps.sh baseline <label> <ticket number> [name=Dockerfile]..." >&2
		exit 2
	}
	shift
	baseline "$@"
	exit "$failed"
	;;
cover | debt)
	"$1"
	;;
*)
	echo "usage: steps.sh build|lint|test|cover|debt|sbom|baseline" >&2
	exit 2
	;;
esac
