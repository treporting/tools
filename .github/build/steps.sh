#!/usr/bin/env bash
# The build steps of every repo (TOOL_23 to TOOL_26), run by .github/workflows/build.yml in the root
# of the repo; the same script runs by hand. The baseline is baseline.sh beside it.
#
#	steps.sh build|lint|test|cover|debt|objectcode
#	steps.sh sbom [name=Dockerfile]...
#
# A Go module is every directory with a go.mod, a JavaScript package every directory with a
# package-lock.json - neither below node_modules nor below a directory whose name starts with a dot.
# build, lint and test exit with 1 where a step fails. cover and debt report and exit with 0; they
# exit otherwise only where the measuring itself did not happen. sbom builds each image named from
# its Dockerfile, and exits with 1 where sbom.md does not list the packages delivered. Where JUNIT
# names a directory, test writes the junit.xml of each run there, for the test reports of
# baseline.sh.
#
# For the object code report (TOOL_35): where OUTPUTS names a directory, build also builds each main
# package of each module into it, copies each bundle of each package there - each file with .js
# that npm run build leaves and git does not track, with its .map - and lists them in
# OUTPUTS/outputs.txt. Where COVERAGE names a directory, test writes there the coverage profile of
# each run of the Go tests and the coverage-final.json of the tests of each package, by Vitest.
# objectcode writes from both the section of each build output into SECTIONS, and exits with 2 where
# it cannot.
#
# It needs go, golangci-lint v2 and deadcode on the PATH, gotestsum where JUNIT is set, and node with
# npm where there is a JavaScript package.
set -uo pipefail

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

#	@implements SR_51@1.1.1
modules() {
	find . -name go.mod -not -path '*/node_modules/*' -not -path '*/.*/*' | sed 's#/go.mod$##' | sort
}

#	@implements SR_52@1.1.1
packages() {
	find . -name package-lock.json -not -path '*/node_modules/*' -not -path '*/.*/*' | sed 's#/package-lock.json$##' | sort
}

#	slug <dir>: the name of a module or package in a file name: root, or its path with - for /
slug() {
	if [ "$1" = . ]; then
		echo root
	else
		echo "${1#./}" | tr '/' '-'
	fi
}

#	@implements SR_51@1.1.1, SR_52@1.1.1, SR_120@1.1.1
build() {
	for m in $(modules); do
		run "$m" go build ./...
		[ -z "${OUTPUTS:-}" ] || binaries "$m"
	done
	for p in $(packages); do
		run "$p" npm ci --no-audit --no-fund
		run "$p" npm run build
		[ -z "${OUTPUTS:-}" ] || bundles "$p"
	done
}

#	output <file>...: copy the files of a build output into a folder of its own in $OUTPUTS, and
#	print the path of the first there
output() {
	local n dir
	mkdir -p "$OUTPUTS"
	n="$(find "$OUTPUTS" -mindepth 2 -maxdepth 2 -path "$OUTPUTS/files/*" -type d | wc -l | tr -d ' ')"
	dir="files/$((n + 1))"
	mkdir -p "$OUTPUTS/$dir"
	cp "$@" "$OUTPUTS/$dir/"
	echo "$dir/$(basename "$1")"
}

#	binaries <module>: build each main package of the module, and list it in $OUTPUTS/outputs.txt as
#	binary <file> <module> <main package>
#
#	@implements SR_120@1.1.1
binaries() {
	local m="$1" root dir main bin
	root="$(cd "$m" && pwd)"
	for dir in $(cd "$m" && go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' ./...); do
		main=".${dir#"$root"}"
		bin="$(mktemp -d)/$(basename "$dir")"
		run "$m" go build -o "$bin" "$main"
		[ -f "$bin" ] || continue
		echo "binary $(output "$bin") $m $main" >>"$OUTPUTS/outputs.txt"
	done
}

#	bundles <package>: list each bundle of the package in $OUTPUTS/outputs.txt as bundle <file>
#	<package>
#
#	@implements SR_120@1.1.1
bundles() {
	local p="$1" f files
	for f in $( (cd "$p" && git ls-files --others --exclude-standard -- '*.js' &&
		git ls-files --others --ignored --exclude-standard -- '*.js') | grep -v '\(^\|/\)node_modules/' | sort -u); do
		files=("$p/$f")
		[ ! -f "$p/$f.map" ] || files+=("$p/$f.map")
		echo "bundle $(output "${files[@]}") $p" >>"$OUTPUTS/outputs.txt"
	done
}

#	the packages eslint.config.js imports, installed beside it once
eslint_install() {
	[ -x "$steps/node_modules/.bin/eslint" ] && return
	(cd "$steps" && npm install --no-save --no-audit --no-fund \
		eslint@10.14.0 @eslint/js@10.0.1 globals@17.13.0) || exit 2
}

#	@implements SR_53@2.1.1, SR_54@2.1.1
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

#	@implements SR_55@1.1.1, SR_56@1.1.1, SR_120@1.1.1
test() {
	local unit integration
	[ -z "${COVERAGE:-}" ] || mkdir -p "$COVERAGE"
	for m in $(modules); do
		unit=()
		integration=()
		if [ -n "${COVERAGE:-}" ]; then
			unit=(-coverpkg=./... "-coverprofile=$COVERAGE/$(slug "$m")-unit.out")
			integration=(-coverpkg=./... "-coverprofile=$COVERAGE/$(slug "$m")-integration.out")
		fi
		if [ -n "${JUNIT:-}" ]; then
			junit "$m" unit -race -timeout 30m ${unit[@]+"${unit[@]}"} ./...
			junit "$m" integration -race -timeout 30m -tags integration ${integration[@]+"${integration[@]}"} ./...
		else
			run "$m" go test -race -timeout 30m ${unit[@]+"${unit[@]}"} ./...
			run "$m" go test -race -timeout 30m -tags integration ${integration[@]+"${integration[@]}"} ./...
		fi
	done
	for p in $(packages); do
		run "$p" npm ci --no-audit --no-fund
		if [ -n "${COVERAGE:-}" ]; then
			run "$p" npm test -- --coverage --coverage.reporter=json "--coverage.reportsDirectory=$COVERAGE/$(slug "$p")"
		else
			run "$p" npm test
		fi
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

#	profile <module>: the coverage profiles of the runs of the Go tests of the module in $COVERAGE, as
#	one; nothing where there is none
profile() {
	local f first=1
	for f in "$COVERAGE/$(slug "$1")"-*.out; do
		[ -f "$f" ] || continue
		if [ "$first" = 1 ]; then
			cat "$f"
			first=0
		else
			tail -n +2 "$f"
		fi
	done
}

#	objectcode: write the section of each build output listed in $OUTPUTS/outputs.txt, with the
#	coverage in $COVERAGE, into $SECTIONS as <n>.md, in the order of the list
#
#	@implements SR_120@1.1.1
objectcode() {
	local bin n=0 kind file dir main section cover flags t c
	bin="$(tool objectcode/cmd/objectcode)"
	OUTPUTS="$(cd "$OUTPUTS" && pwd)" || exit 2
	COVERAGE="$(mkdir -p "$COVERAGE" && cd "$COVERAGE" && pwd)" || exit 2
	mkdir -p "$SECTIONS"
	if [ ! -f "$OUTPUTS/outputs.txt" ]; then
		echo "no build outputs"
		return
	fi
	while read -r kind file dir main <&3; do
		n=$((n + 1))
		section="$SECTIONS/$(printf '%03d' "$n").md"
		echo "== $dir: $kind $file"
		case "$kind" in
		binary)
			cover="$(mktemp)"
			profile "$dir" >"$cover"
			flags=()
			[ ! -s "$cover" ] || flags=(-coverprofile "$cover")
			(cd "$dir" && "$bin" -root . -binary "$OUTPUTS/$file" -main "$main" ${flags[@]+"${flags[@]}"}) >"$section" || exit 2
			rm -f "$cover"
			;;
		bundle)
			[ -d "$dir/node_modules" ] || (cd "$dir" && npm ci --no-audit --no-fund) || exit 2
			flags=()
			for t in webpack vite rollup esbuild @babel/core terser; do
				[ ! -f "$dir/node_modules/$t/package.json" ] || flags+=(-tool "$t")
			done
			for c in $(find "$dir" -maxdepth 1 -type f \( -name '*.config.*' -o -name .babelrc \) | sort); do
				flags+=(-config "$c")
			done
			"$bin" -root . -bundle "$OUTPUTS/$file" -dir "$dir" -coverage "$COVERAGE/$(slug "$dir")/coverage-final.json" \
				${flags[@]+"${flags[@]}"} >"$section" || exit 2
			;;
		esac
	done 3<"$OUTPUTS/outputs.txt"
}

#	@implements SR_67@1.1.1
sbom() {
	local bin
	bin="$(tool sbom/cmd/sbom)"
	build_images "$@"
	run . "$bin" -check ${image_flags[@]+"${image_flags[@]}"}
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
cover | debt | objectcode)
	"$1"
	;;
*)
	echo "usage: steps.sh build|lint|test|cover|debt|objectcode|sbom" >&2
	exit 2
	;;
esac
