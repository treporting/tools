#!/usr/bin/env bash
# The helpers of steps.sh and baseline.sh, read by both with source. It sets steps, the root of
# tools, and failed, which run sets to 1 where a command fails.

#	the root of tools; .github/build holds the configurations of the linters
steps="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

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
