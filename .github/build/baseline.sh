#!/usr/bin/env bash
# The baseline of every repo (TOOL_27), run by .github/workflows/baseline.yml in the root of the
# repo; the same script runs by hand.
#
#	baseline.sh <label> <ticket number> [name=Dockerfile]...
#
# It builds the images, writes sbom.md with the label and the ticket, and checks its licences and
# the approvals; it exits with 1 where they do not fit. Where JUNIT names a directory - the junit.xml
# that steps.sh test wrote there - it writes the test report of each test level from them.
#
# ABOVE names the root of the repo above - specs-and-deploy-process for the tools: baseline takes the schemas
# and levels.yaml from there. APPROVAL_BASE, where set, is the commit baseline checks the approvals
# against: in tools, the first commit, as every approval before it has its review in specs-and-deploy-process.
# TOOLS_TAG and ABOVE_TAG, where set, are the tags of the baselines of the tools and of the repo
# above it runs with; sbom names them in sbom.md (SR_113).
#
# It needs go and docker on the PATH.
set -uo pipefail

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

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
	run . "$bin_sbom" -label "$label" -ticket "$ticket" ${TOOLS_TAG:+-tools "$TOOLS_TAG"} \
		${ABOVE_TAG:+-above "$ABOVE_TAG"} ${image_flags[@]+"${image_flags[@]}"}
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

[ $# -ge 2 ] || {
	echo "usage: baseline.sh <label> <ticket number> [name=Dockerfile]..." >&2
	exit 2
}
baseline "$@"
exit "$failed"
