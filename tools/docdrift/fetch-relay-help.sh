#!/usr/bin/env bash
# Capture `flockdeck-relay serve -h` from the RELEASED relay image.
#
#   fetch-relay-help.sh OUT_FILE
#
# Writes the help text to OUT_FILE and the version it came from to
# OUT_FILE.version. RELAY_TAG pins a version (a bare semver like 0.2.50);
# without it, the highest semver tag the registry lists is used. RELAY_IMAGE
# overrides the repository.
#
# The image has to be pullable anonymously, like everything else this repo's
# CI reads: no credential is given to a public repository's workflows. If it
# is not, this fails and says so, rather than falling back to something stale
# — the public ghcr.io/jmwri/flockdeck-relay froze at 0.2.44 when the relay
# moved to the Flockdeck organization, and quietly checking the docs against
# a months-old binary is the failure this whole tool exists to prevent.
set -euo pipefail

out=${1:?usage: fetch-relay-help.sh OUT_FILE}
image=${RELAY_IMAGE:-ghcr.io/flockdeck/flockdeck-relay}
repo=${image#ghcr.io/}

private() {
  echo "cannot read ${image} anonymously: the package is private or missing." >&2
  echo "Make it public (Package settings > Change visibility), or set RELAY_IMAGE to a public repository." >&2
  exit 1
}

token=$(curl -fsS "https://ghcr.io/token?service=ghcr.io&scope=repository:${repo}:pull" |
  sed -n 's/.*"token":"\([^"]*\)".*/\1/p') || private

tag=${RELAY_TAG:-}
if [ -z "$tag" ]; then
  tags=""
  url="https://ghcr.io/v2/${repo}/tags/list?n=1000"
  while [ -n "$url" ]; do
    page=$(curl -fsS -D /tmp/relay-tags.headers -H "Authorization: Bearer ${token}" "$url") || private
    tags+=$(printf '%s' "$page" | grep -oE '"[0-9]+\.[0-9]+\.[0-9]+"' | tr -d '"')$'\n'
    next=$(tr -d '\r' </tmp/relay-tags.headers | sed -n 's/^[Ll]ink: <\([^>]*\)>; rel="next".*/\1/p')
    if [ -n "$next" ]; then url="https://ghcr.io${next}"; else url=""; fi
  done
  tag=$(printf '%s' "$tags" | sort -V | tail -n1)
  [ -n "$tag" ] || { echo "${image} lists no semver tags" >&2; exit 1; }
fi

echo "relay image: ${image}:${tag}" >&2
# Anonymous, on purpose: a runner that happened to be logged in would hide a
# package that stopped being public.
export DOCKER_CONFIG
DOCKER_CONFIG=$(mktemp -d)
# Pulled apart from the run so the pull's progress is not in the help text.
docker pull -q "${image}:${tag}" >&2 || private
docker run --rm "${image}:${tag}" serve -h >"$out" 2>&1
printf '%s\n' "$tag" >"${out}.version"
