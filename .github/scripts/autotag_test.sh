#!/bin/sh
# Tests of autotag.sh, against real git repositories made here. Run by
# .github/workflows/autotag-test.yml, and by hand: sh .github/scripts/autotag_test.sh
set -eu

here=$(cd "$(dirname "$0")" && pwd)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0
export GIT_AUTHOR_NAME=t GIT_AUTHOR_EMAIL=t@example.invalid GIT_COMMITTER_NAME=t GIT_COMMITTER_EMAIL=t@example.invalid

failures=0
pass() { printf 'ok   %s\n' "$1"; }
bad() {
	printf 'FAIL %s\n' "$1"
	failures=$((failures + 1))
}
# check NAME WANT GOT
check() {
	if [ "$2" = "$3" ]; then pass "$1"; else
		bad "$1: want '$2', got '$3'"
	fi
}
# holds NAME COMMAND...: the command succeeds.
holds() {
	name=$1
	shift
	if "$@" > /dev/null 2>&1; then pass "$name"; else bad "$name"; fi
}
# refuses NAME COMMAND...: the command fails.
refuses() {
	name=$1
	shift
	if "$@" > /dev/null 2>&1; then bad "$name (was accepted)"; else pass "$name"; fi
}

AUTOTAG_SOURCED=1
export AUTOTAG_SOURCED
# shellcheck source=autotag.sh disable=SC1091
. "$here/autotag.sh"

# --- next_patch_tag ---------------------------------------------------------

check "next patch" v0.9.24 "$(next_patch_tag v0.9.23)"
check "next patch, docs series" v0.1.4 "$(next_patch_tag v0.1.3)"
check "9 is followed by 10" v0.9.10 "$(next_patch_tag v0.9.9)"
check "99 is followed by 100" v1.2.100 "$(next_patch_tag v1.2.99)"
check "0 is followed by 1" v0.9.1 "$(next_patch_tag v0.9.0)"
for bad_tag in "" v0.9 0.9.1 v0.9.1-rc.1 v0.9.x "v0.9.1 v0.9.2" latest; do
	refuses "next patch refuses '$bad_tag'" next_patch_tag "$bad_tag"
done

# --- fixtures ---------------------------------------------------------------

# newrepo DIR: a bare origin (DIR.git), and a clone of it at DIR with a first
# commit on main that has the site, a README, .github and content/, a
# .dockerignore like the real ones, and the tag v0.9.9.
newrepo() {
	git init -q --bare -b main "$1.git"
	git init -q -b main "$1.seed"
	(
		cd "$1.seed"
		mkdir -p .github/workflows content
		echo site > index.html
		echo readme > README.md
		echo w > .github/workflows/ci.yml
		echo c > content/page.md
		echo 'FROM nginx' > Dockerfile
		printf '# not the site\n.git\n.github\n.dockerignore\nDockerfile\nREADME.md\nLICENSE\ncontent\n' > .dockerignore
		git add -A
		git commit -q -m first
		git tag v0.9.9
		git push -q "$1.git" main --tags
	)
	git clone -q "$1.git" "$1"
}

# rnd: something new for every change to hold.
rnd() { od -An -N8 -tu8 /dev/urandom | tr -d " "; echo; }

# commit DIR FILE: changes FILE in the clone and pushes it to main.
commit() {
	mkdir -p "$(dirname "$1/$2")"
	rnd >> "$1/$2"
	git -C "$1" add -A
	git -C "$1" commit -q -m "change $2"
	git -C "$1" push -q origin main
	git -C "$1" rev-parse HEAD
}

# --- latest_tag: numeric, and only versions ---------------------------------

r=$work/latest
newrepo "$r"
cd "$r"
for t in v0.9.2 v0.9.10 v0.9.3 v1.0.0-rc.1 vfoo latest v0.9 nightly-v9.9.9; do git tag "$t"; done
check "latest tag is the highest by number, and only a version" v0.9.10 "$(latest_tag)"

# --- ships ------------------------------------------------------------------

r=$work/ships
newrepo "$r"
cd "$r"
base=$(git rev-parse HEAD)
shipping() { # FILE...: does a commit changing them ship?
	git reset -q --hard "$base"
	for f in "$@"; do
		mkdir -p "$(dirname "$f")"
		rnd >> "$f"
	done
	git add -A
	git commit -q -m x
	ships "$base" HEAD
}
holds "a page ships" shipping index.html
holds "a new page ships" shipping about.html
holds "a stylesheet in a folder ships" shipping assets/site.css
holds "nginx.conf ships" shipping nginx.conf
holds "the Dockerfile ships: its base image is the deployment" shipping Dockerfile
refuses ".dockerignore does not ship: no build reads it into the image" shipping .dockerignore
refuses ".gitattributes does not ship" shipping .gitattributes
refuses ".editorconfig does not ship" shipping .editorconfig
refuses ".gitignore does not ship" shipping .gitignore
refuses "CODEOWNERS does not ship" shipping CODEOWNERS
refuses "a CODEOWNERS under .github does not ship" shipping .github/CODEOWNERS
refuses "repository config in a folder does not ship" shipping assets/.gitignore
holds "look-alikes of the config names ship" shipping .gitattributes.bak CODEOWNERS.html
refuses "README.md does not ship" shipping README.md
refuses "a workflow does not ship" shipping .github/workflows/ci.yml
refuses "a new workflow does not ship" shipping .github/workflows/autotag.yml
refuses "a script under .github does not ship" shipping .github/scripts/autotag.sh
refuses "content/ does not ship (docgen's generated pages are what does)" shipping content/page.md
refuses "a nested file of content/ does not ship" shipping content/deep/er/page.md
refuses "several things that do not ship, together, do not ship" shipping README.md .github/x.yml content/a.md
holds "one thing that ships among things that do not" shipping README.md .github/x.yml index.html
holds "config files do not hide a page that ships beside them" shipping .gitattributes .editorconfig index.html
refuses "repository config, together, does not ship" shipping .gitattributes .editorconfig .gitignore .dockerignore CODEOWNERS
refuses "no change at all does not ship" ships "$base" "$base"
git reset -q --hard "$base"
git rm -q .dockerignore
git commit -q -m "no .dockerignore"
noignore=$(git rev-parse HEAD)
echo r >> README.md
git add -A
git commit -q -m readme
holds "with no .dockerignore, everything counts" ships "$noignore" HEAD
git reset -q --hard "$base"
printf '.github\n!.github/keep\n' > .dockerignore
git add -A
git commit -q -m negation
neg=$(git rev-parse HEAD)
echo r >> .github/x.yml
git add -A
git commit -q -m x
holds "a negation .dockerignore this cannot follow makes everything count" ships "$neg" HEAD
echo r >> .gitattributes
git add -A
git commit -q -m attrs
refuses "even then, repository config does not ship" ships "$(git rev-parse HEAD~1)" HEAD
git reset -q --hard "$base"
git mv index.html .gitignore
git commit -q -m "rename a page onto a config name"
holds "a page renamed onto a config name is still seen leaving" ships "$base" HEAD

# --- the whole run ----------------------------------------------------------

REPO=Flockdeck/flockdeck-docs
export REPO
green=1
ci_green() { [ "$green" = 1 ]; }
tagsat() { git -C "$1.git" tag --points-at "$2" | tr '\n' ' '; }
alltags() { git -C "$1.git" tag | sort -t. -k3,3n | tr '\n' ' '; }
# runmain runs main in a subshell, as a run of the script is a process of its
# own, and says how it ended.
runmain() {
	rc=0
	(main) 2>&1 || rc=$?
	echo "exit=$rc"
}
# tagrun DIR [SHA]: runs main in a clone, as the workflow does (a fresh
# checkout of main's history and all the tags), and says what it printed.
tagrun() {
	(
		cd "$1"
		if [ -n "${2:-}" ]; then SHA=$2; else unset SHA; fi
		export SHA
		runmain
	)
}

r=$work/e2e
newrepo "$r"
c1=$(commit "$r" index.html)
out=$(tagrun "$r" "$c1")
check "a shipping green commit is tagged with the next patch" "v0.9.10" "$(tagsat "$r" "$c1" | tr -d ' ')"

out=$(tagrun "$r" "$c1")
check "the same commit is not tagged twice" "v0.9.10" "$(tagsat "$r" "$c1" | tr -d ' ')"
check "and no other tag appeared" "v0.9.9 v0.9.10 " "$(alltags "$r")"
case $out in *"tagged already"*) pass "and it says why" ;; *) bad "did not say it was tagged already: $out" ;; esac

c2=$(commit "$r" README.md)
out=$(tagrun "$r" "$c2")
check "a README-only change is not tagged" "" "$(tagsat "$r" "$c2")"
case $out in *"nothing that ships"*) pass "and it says why" ;; *) bad "did not say nothing shipped: $out" ;; esac

c3=$(commit "$r" .github/workflows/ci.yml)
out=$(tagrun "$r" "$c3")
check "a workflow-only change is not tagged" "" "$(tagsat "$r" "$c3")"

green=0
c4=$(commit "$r" index.html)
out=$(tagrun "$r" "$c4")
check "a commit whose CI has not passed is never tagged" "" "$(tagsat "$r" "$c4")"
case $out in *"CI has not passed"*) pass "and it says why" ;; *) bad "did not say CI had not passed: $out" ;; esac

green=1
out=$(tagrun "$r" "$c4")
check "it is tagged when CI has passed, counting from the last tag, with the commits between it" "v0.9.11" "$(tagsat "$r" "$c4" | tr -d ' ')"

# The commit before, if its own run finishes late, is already in the tag.
c5=$(commit "$r" index.html)
c6=$(commit "$r" index.html)
tagrun "$r" "$c6" > /dev/null
out=$(tagrun "$r" "$c5")
check "an older commit than the tagged one is not tagged after it" "" "$(tagsat "$r" "$c5")"
case $out in *"already part of"*) pass "and it says why" ;; *) bad "did not say it was part of a tag already: $out" ;; esac

# With no SHA, it is the tip of main.
c7=$(commit "$r" index.html)
tagrun "$r" > /dev/null
check "with no commit named, main's tip is tagged" "v0.9.13" "$(tagsat "$r" "$c7" | tr -d ' ')"

# Something that is not on main is not tagged, whatever it is.
git -C "$r" checkout -q -b topic
echo t >> "$r/index.html"
git -C "$r" commit -q -am topic
topic=$(git -C "$r" rev-parse HEAD)
git -C "$r" push -q origin topic
git -C "$r" checkout -q main
out=$(tagrun "$r" "$topic")
check "a commit that is not on main is not tagged" "" "$(tagsat "$r" "$topic")"

# Two runs, or a person, tagging together: the first push wins and the other
# works it out again from the tags as they are.
r=$work/race
newrepo "$r"
ca=$(commit "$r" index.html)
# Another clone's run tags ca as v0.9.10 first.
git clone -q "$r.git" "$r.other"
(cd "$r.other" && SHA=$ca && export SHA && runmain > /dev/null)
check "the other run tagged the commit" "v0.9.10" "$(tagsat "$r" "$ca" | tr -d ' ')"
cb=$(commit "$r" index.html)
# This one has not heard: its first look at the tags is out of date.
out=$(
	cd "$r"
	stale=1
	# shellcheck disable=SC2317 # called by main
	refresh_tags() {
		if [ "$stale" = 1 ]; then stale=0; else git fetch -q --tags --force origin; fi
	}
	# ca was tagged v0.9.10 by the other run; this one thinks v0.9.9 is the
	# latest, and that ca is not tagged, and picks v0.9.10 for cb.
	SHA=$cb runmain
)
check "a tag that was taken is taken again from the tags as they are now" "v0.9.11" "$(tagsat "$r" "$cb" | tr -d ' ')"
check "and the tag that was taken is where it was" "v0.9.10" "$(tagsat "$r" "$ca" | tr -d ' ')"
case $out in *"was taken first"*) pass "and it says so" ;; *) bad "did not say the tag was taken: $out" ;; esac

# A tag that will not be won is an error, not a quiet success.
r=$work/lost
newrepo "$r"
cl=$(commit "$r" index.html)
out=$(
	cd "$r"
	# The remote refuses every push, as a token that cannot write would have it.
	printf '#!/bin/sh\nexit 1\n' > "$r.git/hooks/pre-receive"
	chmod +x "$r.git/hooks/pre-receive"
	ATTEMPTS=2 SHA=$cl runmain
)
case $out in *"exit=1"*) pass "a push that is refused fails the run" ;; *) bad "a refused push was not an error: $out" ;; esac
check "and nothing is tagged" "" "$(tagsat "$r" "$cl")"

# Dry run tags nothing.
r=$work/dry
newrepo "$r"
cd_=$(commit "$r" index.html)
out=$(cd "$r" && DRY_RUN=1 SHA=$cd_ runmain)
check "a dry run tags nothing" "" "$(tagsat "$r" "$cd_")"
case $out in *"would tag $cd_ as v0.9.10"*) pass "and says what it would tag" ;; *) bad "dry run said: $out" ;; esac

if [ "$failures" -ne 0 ]; then
	printf '\n%s failed\n' "$failures"
	exit 1
fi
printf '\nall passed\n'
