#!/bin/sh
# Tags a commit of main with the next patch version, which is what deploys it.
#
# Merging to main here deploys nothing. A tag like v0.9.24 is what CI builds
# into the bare-semver image (0.9.24) that Flux's ImagePolicy in k8s-infra
# watches for, and Flux rolls out. A merge that nobody tagged is live in the
# repository and not on the site, which has been the state of this site for
# days at a time.
#
# It is run by .github/workflows/autotag.yml, in a checkout with every tag and
# all of main's history, holding a token that can push a tag and is not
# GITHUB_TOKEN: a tag pushed with that starts no workflow, so the tag's own
# image build would never run.
#
# It tags a commit only if all of these hold, and otherwise says why not and
# succeeds:
#
#   - it is on main;
#   - no version tag is on it already, and none is on a commit that has it in
#     its history (so a later commit tagged first, and this one is already
#     released by it);
#   - CI passed on it (ci_green);
#   - something that ships changed since the latest tag (ships).
#
# Tags are chosen as the latest version tag with its patch number one higher,
# so each repository carries on its own series. A push that finds the tag
# taken -- two runs, or a person, tagging together -- is tried again from the
# tags as they now are.
#
# Environment: REPO (owner/name, for asking about CI), SHA (the commit; default
# the tip of main), GH_TOKEN (to ask about CI), BASE_BRANCH (default main),
# CI_WORKFLOW (default ci.yml), DRY_RUN (set: say what would be tagged, tag
# nothing).
set -eu

BASE_BRANCH=${BASE_BRANCH:-main}
CI_WORKFLOW=${CI_WORKFLOW:-ci.yml}
ATTEMPTS=${ATTEMPTS:-3}

say() { printf '%s\n' "$*"; }
skip() {
	say "not tagging: $*"
	if [ -n "${GITHUB_STEP_SUMMARY:-}" ]; then
		printf 'Not tagged: %s\n' "$*" >> "$GITHUB_STEP_SUMMARY"
	fi
	exit 0
}
fail() {
	printf '::error::%s\n' "$*" >&2
	exit 1
}

# next_patch_tag v0.9.23 -> v0.9.24. Numeric, so 9 is followed by 10.
next_patch_tag() {
	printf '%s\n' "$1" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$' || return 1
	base=${1%.*}
	patch=${1##*.}
	printf '%s.%s\n' "$base" "$((patch + 1))"
}

# version_tags lists the version tags, oldest first: v1.2.3 and nothing else,
# so a candidate (v1.2.3-rc.1) or a tag of some other kind is not the latest.
version_tags() {
	git tag -l 'v*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -t. -k1.2,1n -k2,2n -k3,3n
}

latest_tag() { version_tags | tail -n 1; }

# tags_at lists the version tags on a commit.
tags_at() {
	git tag --points-at "$1" -l 'v*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' || true
}

# shipping_files lists what changed between two commits that ends up in the
# image. The Dockerfile copies the whole directory into it, less what
# .dockerignore names, and that file is the authority: a change to README.md,
# to .github, or to content/ (which docs' generator reads, and whose output
# arrives as a change of its own) is not a change to the site. Dockerfile and
# .dockerignore are kept out of the image by .dockerignore but change how it
# is built, so they count. Patterns are used as shell patterns matched against
# the path or any folder above it; one that this cannot follow, a negation,
# makes everything count, since a tag that did not need to be pushed costs a
# deploy and one that was needed and missed costs a stale site.
shipping_files() {
	ignore=$(git show "$2:.dockerignore" 2> /dev/null || true)
	if printf '%s\n' "$ignore" | grep -q '^[[:space:]]*!'; then
		git -c core.quotepath=off diff --name-only "$1" "$2"
		return
	fi
	git -c core.quotepath=off diff --name-only "$1" "$2" | while IFS= read -r path; do
		case $path in Dockerfile | .dockerignore) printf '%s\n' "$path"; continue ;; esac
		ignored=""
		while IFS= read -r pat; do
			pat=${pat%%#*}
			pat=$(printf '%s' "$pat" | sed 's/^[[:space:]]*//; s/[[:space:]]*$//; s|^/||; s|/$||')
			[ -n "$pat" ] || continue
			# shellcheck disable=SC2254 # a pattern is what this is for
			case $path in $pat | $pat/*) ignored=1; break ;; esac
		done << EOF
$ignore
EOF
		[ -n "$ignored" ] || printf '%s\n' "$path"
	done
}

ships() { [ -n "$(shipping_files "$1" "$2")" ]; }

# ci_green: the latest run of the CI workflow for the commit, on a push to
# main, finished and succeeded. A run still going, or none yet, is not green;
# the workflow_run trigger comes again when it finishes.
ci_green() {
	state=$(gh api "repos/$REPO/actions/workflows/$CI_WORKFLOW/runs?head_sha=$1&branch=$BASE_BRANCH&event=push&per_page=20" \
		--jq '[.workflow_runs[]] | sort_by(.created_at) | last | (.status + "/" + (.conclusion // "none"))') || return 1
	[ "$state" = "completed/success" ]
}

refresh_tags() { git fetch -q --tags --force origin; }

main() {
	: "${REPO:?REPO is required}"
	refresh_tags
	sha=${SHA:-$(git rev-parse "origin/$BASE_BRANCH")}
	git cat-file -e "$sha^{commit}" 2> /dev/null || fail "$sha is not a commit here"
	sha=$(git rev-parse "$sha^{commit}")
	git merge-base --is-ancestor "$sha" "origin/$BASE_BRANCH" || skip "$sha is not on $BASE_BRANCH"

	attempt=1
	while :; do
		existing=$(tags_at "$sha" | tr '\n' ' ')
		[ -z "$existing" ] || skip "$sha is tagged already ($existing)"
		last=$(latest_tag)
		[ -n "$last" ] || fail "there is no version tag to count on from"
		if git merge-base --is-ancestor "$sha" "$last"; then
			skip "$sha is already part of $last"
		fi
		ci_green "$sha" || skip "CI has not passed on $sha"
		ships "$last" "$sha" || skip "nothing that ships has changed since $last"
		next=$(next_patch_tag "$last") || fail "'$last' is not a version"

		if [ -n "${DRY_RUN:-}" ]; then
			say "would tag $sha as $next (after $last)"
			return 0
		fi
		if git push -q origin "$sha:refs/tags/$next"; then
			say "tagged $sha as $next"
			if [ -n "${GITHUB_STEP_SUMMARY:-}" ]; then
				printf 'Tagged %s as **%s**; its CI builds and pushes the image.\n' "$sha" "$next" >> "$GITHUB_STEP_SUMMARY"
			fi
			return 0
		fi
		[ "$attempt" -lt "$ATTEMPTS" ] || fail "could not push $next after $attempt attempts"
		say "$next was taken first; working it out again from the tags as they are now"
		attempt=$((attempt + 1))
		refresh_tags
	done
}

# Sourced by the tests, which call the functions one at a time.
if [ -z "${AUTOTAG_SOURCED:-}" ]; then
	main "$@"
fi
