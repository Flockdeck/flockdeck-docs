#!/usr/bin/env bash
# Refresh tools/docdrift/snapshot from the latest released relay, by hand:
#
#   RELAY_GHCR_READ_USER=... RELAY_GHCR_READ_TOKEN=... tools/docdrift/refresh-relay-snapshot.sh [VERSION]
#
# Then compare the docs (`docdrift relay -help snapshot/relay-serve-help.txt
# -docs content/self-hosting/configuration.md -allow allow-relay-flags.txt`),
# fix or allowlist what it finds, and commit snapshot and docs together.
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
RELAY_TAG=${1:-} "$here/fetch-relay-help.sh" "$work/help.txt"
# Normalize line endings so the committed snapshot is the same from any OS.
tr -d '\r' <"$work/help.txt" >"$here/snapshot/relay-serve-help.txt"
cp "$work/help.txt.version" "$here/snapshot/relay.version"
echo "snapshot is now relay $(cat "$here/snapshot/relay.version")"
