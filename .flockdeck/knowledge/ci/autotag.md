---
type: CI Check
title: Auto-tag workflow
description: Mints a release-bot token and runs autotag.sh to tag green, shipping commits on main.
resource: .github/workflows/autotag.yml
tags: [ci, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml, .github/scripts/autotag.sh]
---
Triggers: `workflow_run` of "CI" completed on main; hourly cron `17 * * * *` as a safety net for missed events or re-run CI; `workflow_dispatch`. Concurrency group `autotag`, never cancelled (a waiting run is replaced by a newer one, which is harmless because the newer commit contains the older). Top-level `permissions: {}`.

Job `tag` (10-minute timeout, environment `release`, permissions `contents: read`, `actions: read`). It runs only if the event is not `workflow_run`, or the run concluded `success`, was a `push`, and its head repository is this repository.

Steps:
1. `actions/create-github-app-token` (pinned to a commit SHA) with [RELEASE_BOT_APP_ID](/config/release-bot-app-id.md) and [RELEASE_BOT_APP_PRIVATE_KEY](/config/release-bot-app-private-key.md), `permission-contents: write`. See [release bot token](/security/release-bot-token.md).
2. Checkout with `fetch-depth: 0` and the app token.
3. `sh .github/scripts/autotag.sh` with `REPO`, `SHA` (from the triggering run, empty otherwise meaning main's tip) and `GH_TOKEN`. SHA is passed through env, never spliced into the script.

The job runs nothing but the script in main. Tagging makes no commit and the tag's own CI run is for a tag, so it cannot retrigger itself. Tests: [Auto-tag tests](/ci/autotag-test.md). Script: [autotag.sh](/cli/autotag-sh.md). Overview: [release by tag](/build/release-by-tag.md).
