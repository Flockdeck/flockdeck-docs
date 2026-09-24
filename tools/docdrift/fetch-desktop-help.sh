#!/usr/bin/env bash
# Capture `flockdeck -h` from the RELEASED desktop binary.
#
#   fetch-desktop-help.sh OUT_FILE
#
# Downloads the linux amd64 tarball of the latest GitHub release of
# Flockdeck/flockdeck (or DESKTOP_TAG, e.g. v0.3.40), checks it against the
# release's checksums.txt, and runs it. Writes the help text to OUT_FILE and
# the tag to OUT_FILE.version. Needs gh (GH_TOKEN set) and, to run the binary,
# the GTK 4 and WebKitGTK 6 libraries it links.
set -euo pipefail

out=${1:?usage: fetch-desktop-help.sh OUT_FILE}
repo=Flockdeck/flockdeck

tag=${DESKTOP_TAG:-$(gh release view --repo "$repo" --json tagName --jq .tagName)}
[ -n "$tag" ] || { echo "no latest release for $repo" >&2; exit 1; }
tarball="flockdeck_${tag}_linux_amd64.tar.gz"

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
gh release download "$tag" --repo "$repo" --dir "$work" --pattern "$tarball" --pattern checksums.txt

# The whole checksums file is not checked, only this artifact's line, and it
# must be there: `sha256sum -c` on an empty selection succeeds.
grep -E "[[:space:]]\*?${tarball}\$" "$work/checksums.txt" >"$work/one.sha256" ||
  { echo "checksums.txt has no entry for ${tarball}" >&2; exit 1; }
(cd "$work" && sha256sum -c one.sha256 >&2)

tar -xzf "$work/$tarball" -C "$work" flockdeck
echo "desktop release: ${tag}" >&2
# -h exits 0 and prints to stderr; both go to the file.
"$work/flockdeck" -h >"$out" 2>&1
printf '%s\n' "$tag" >"${out}.version"
