#!/usr/bin/env bash
# Keep ONE tracking issue per drift check: open it when there is drift, edit
# it when the drift changes, close it when it is gone, and do nothing (no
# comment, no edit) when it has not changed.
#
#   issue.sh sync KEY TITLE BODY_FILE   drift: open or update the issue
#   issue.sh clear KEY COMMENT          no drift: close the issue, if any
#
# KEY names the check (generated, relay, desktop). The issue is found by the
# `docs-drift` label plus a marker line for KEY in its body, so its title can
# change (it names a version) without losing it. Needs gh, GH_TOKEN with
# issues: write, and GH_REPO or a checkout.
set -euo pipefail

label=docs-drift
cmd=${1:?usage: issue.sh sync|clear KEY ...}
key=${2:?KEY}
marker="<!-- docdrift-check: ${key} -->"

existing() {
  gh issue list --label "$label" --state open --limit 100 --json number,body \
    --jq ".[] | select(.body | contains(\"${marker}\")) | .number" | head -n1
}

case "$cmd" in
sync)
  title=${3:?TITLE}
  body="$marker
$(cat "${4:?BODY_FILE}")"
  n=$(existing)
  if [ -z "$n" ]; then
    gh label create "$label" --color d4c5f9 --description "The docs disagree with a released Flockdeck artifact" 2>/dev/null || true
    gh issue create --title "$title" --body "$body" --label "$label"
    exit 0
  fi
  have=$(gh issue view "$n" --json title,body --jq '.title + "\n" + .body')
  if [ "$have" = "$(printf '%s\n%s' "$title" "$body")" ]; then
    echo "issue #$n is already up to date"
  else
    gh issue edit "$n" --title "$title" --body "$body"
    echo "updated issue #$n"
  fi
  ;;
clear)
  n=$(existing)
  if [ -n "$n" ]; then
    gh issue close "$n" --comment "${3:?COMMENT}"
    echo "closed issue #$n"
  else
    echo "no open ${key} drift issue"
  fi
  ;;
*)
  echo "usage: issue.sh sync|clear KEY ..." >&2
  exit 2
  ;;
esac
