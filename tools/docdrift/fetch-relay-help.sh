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
# The image is private on purpose (the built relay comes with an Enterprise
# licence), so this needs a credential, and pull requests never run it: the
# PR check reads the checked-in snapshot instead (tools/docdrift/snapshot).
# Only the scheduled job, in the main-only `relay-release` environment, and a
# maintainer refreshing the snapshot by hand get here. It needs:
#
#   RELAY_GHCR_READ_USER    a GitHub user or bot that can read the package
#   RELAY_GHCR_READ_TOKEN   a token for it with the read:packages scope (a classic
#                      personal access token: fine-grained ones cannot read
#                      packages) and nothing else
#
# Missing either, or refused by the registry, this exits 3 (a credential
# problem, which the scheduled job reports as an issue rather than a red run);
# any other failure exits 1 or 2. It fails and says so; it never falls back to an older
# image, because checking the docs against a stale binary is the failure this
# tool exists to prevent.
set -euo pipefail

out=${1:?usage: fetch-relay-help.sh OUT_FILE}
image=${RELAY_IMAGE:-ghcr.io/flockdeck/flockdeck-relay}
repo=${image#ghcr.io/}

[ -n "${RELAY_GHCR_READ_USER:-}" ] && [ -n "${RELAY_GHCR_READ_TOKEN:-}" ] || {
  echo "RELAY_GHCR_READ_USER and RELAY_GHCR_READ_TOKEN are not set: reading the private relay image needs a read:packages token." >&2
  echo "See 'The relay snapshot' in the README for what to create and where it is stored." >&2
  exit 3
}

denied() {
  echo "cannot read ${image} with the given credential: it may be expired, or lack read:packages or access to the package." >&2
  exit 3
}

token=$(curl -fsS -u "${RELAY_GHCR_READ_USER}:${RELAY_GHCR_READ_TOKEN}" \
  "https://ghcr.io/token?service=ghcr.io&scope=repository:${repo}:pull" |
  sed -n 's/.*"token":"\([^"]*\)".*/\1/p') || denied

tag=${RELAY_TAG:-}
if [ -z "$tag" ]; then
  tags=""
  url="https://ghcr.io/v2/${repo}/tags/list?n=1000"
  while [ -n "$url" ]; do
    page=$(curl -fsS -D "${TMPDIR:-/tmp}/relay-tags.headers" -H "Authorization: Bearer ${token}" "$url") || denied
    tags+=$(printf '%s' "$page" | grep -oE '"[0-9]+\.[0-9]+\.[0-9]+"' | tr -d '"')$'\n'
    next=$(tr -d '\r' <"${TMPDIR:-/tmp}/relay-tags.headers" | sed -n 's/^[Ll]ink: <\([^>]*\)>; rel="next".*/\1/p')
    if [ -n "$next" ]; then url="https://ghcr.io${next}"; else url=""; fi
  done
  tag=$(printf '%s' "$tags" | sort -V | tail -n1)
  [ -n "$tag" ] || { echo "${image} lists no semver tags" >&2; exit 1; }
fi

echo "relay image: ${image}:${tag}" >&2
# A private config dir, so the login is not left behind on a shared machine.
export DOCKER_CONFIG
DOCKER_CONFIG=$(mktemp -d)
trap 'rm -rf "$DOCKER_CONFIG"' EXIT
printf '%s' "$RELAY_GHCR_READ_TOKEN" | docker login ghcr.io -u "$RELAY_GHCR_READ_USER" --password-stdin >&2 || denied
# Pulled apart from the run so the pull's progress is not in the help text.
docker pull -q "${image}:${tag}" >&2 || denied
docker run --rm "${image}:${tag}" serve -h >"$out" 2>&1
printf '%s\n' "$tag" >"${out}.version"
