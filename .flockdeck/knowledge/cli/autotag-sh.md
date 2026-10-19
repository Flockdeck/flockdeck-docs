---
type: CLI Command
title: autotag.sh
description: Tags a green main commit with the next patch version when something that ships changed.
resource: .github/scripts/autotag.sh
tags: [release, git, shell]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh]
---
POSIX `sh`, `set -eu`. Run by [Auto-tag](/ci/autotag.md) in a full-history checkout. Tags are what deploy, see [release by tag](/build/release-by-tag.md).

Tags a commit only if all hold, otherwise prints `not tagging: <why>` (also to the step summary) and exits 0:
1. The commit is an ancestor of `origin/$BASE_BRANCH`.
2. No version tag is on it.
3. It is not an ancestor of the latest version tag (already released by a later tag).
4. `ci_green`: the latest `CI_WORKFLOW` run for that SHA, branch main, event `push`, is `completed/success` (asked through `gh api`).
5. `ships`: something in the image changed since the latest tag.

Functions:

| Function | Behaviour |
|---|---|
| `next_patch_tag` | `v0.9.23` to `v0.9.24`, numeric so 9 is followed by 10; anything not `vN.N.N` fails |
| `version_tags` / `latest_tag` | Only `^v[0-9]+\.[0-9]+\.[0-9]+$` tags, sorted numerically; `v1.2.3-rc.1` is ignored |
| `dockerignored_files` | `git diff --name-only` between two commits, minus paths matched by the target commit's [.dockerignore](/build/dockerignore.md) (pattern or any folder above it). `Dockerfile` and `.dockerignore` always count. A negation line (`!`) makes every changed file count |
| `shipping_files` | Above, minus `.gitattributes`, `.editorconfig`, `.gitignore`, `CODEOWNERS` in any folder |
| `ships` | `shipping_files` is non-empty |

Push: `git push origin <sha>:refs/tags/<next>`. If refused (tag taken by a concurrent run), it refreshes tags and recomputes, up to `ATTEMPTS` (default 3), then fails with `::error::`. No version tag at all is also an error. `DRY_RUN` prints `would tag` and tags nothing. When `AUTOTAG_SOURCED` is set the script only defines functions, which [its tests](/ci/autotag-test.md) use.

Environment: [REPO](/config/repo.md), [SHA](/config/sha.md), [GH_TOKEN](/config/gh-token.md), [BASE_BRANCH](/config/base-branch.md), [CI_WORKFLOW](/config/ci-workflow.md), [ATTEMPTS](/config/attempts.md), [DRY_RUN](/config/dry-run.md), [AUTOTAG_SOURCED](/config/autotag-sourced.md).
