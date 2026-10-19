---
type: CI Check
title: Auto-tag tests workflow
description: Runs the shell test suite and shellcheck for the auto-tag script on pull requests and main.
resource: .github/workflows/autotag-test.yml
tags: [ci, shell, tests]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag-test.yml, .github/scripts/autotag_test.sh, .github/scripts/autotag.sh]
---
Triggers: push to main and pull requests, only when `.github/scripts/**` or `.github/workflows/autotag*.yml` change. `contents: read`, checkout with `persist-credentials: false`, no secrets, so fork PRs may run it.

Steps: `sh .github/scripts/autotag_test.sh`, then `shellcheck -s sh` on both scripts.

`autotag_test.sh` sources `autotag.sh` with `AUTOTAG_SOURCED` set and builds real git repositories (a bare origin plus clones, tag `v0.9.9`). It covers:
- `next_patch_tag` numeric increments, `latest_tag` ignoring non-version tags.
- `ships`: pages, stylesheets, `nginx.conf`, `Dockerfile` and `.dockerignore` ship; `.gitattributes`, `.editorconfig`, `.gitignore`, `CODEOWNERS`, README, workflows, `.github` scripts and `content/` do not; look-alike names ship; no `.dockerignore` or a negation line makes everything count; a rename onto a config name still counts.
- The whole run: already tagged, nothing ships, CI not passed, already part of a later tag, default SHA is main's tip, commit off main, two runners racing for one tag name (retries), a refused push failing the run, dry run tagging nothing.

Subject: [autotag.sh](/cli/autotag-sh.md).
