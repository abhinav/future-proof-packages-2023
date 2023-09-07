#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

# HACK:
# Always use pikchr in dark mode.

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

args=("--dark-mode")
while [[ $# -gt 0 ]]; do
	case "$1" in
	--dark-mode)
		shift # in case it was provided
		;;
	*)
		args+=("$1")
		shift
		;;
	esac
done

exec "$SCRIPTDIR"/../third_party/bin/pikchr "${args[@]}"
