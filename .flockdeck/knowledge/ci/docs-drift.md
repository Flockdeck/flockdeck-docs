---
type: CI Check
title: Docs drift workflow
description: Daily check of the site against the latest released flockdeck, relay and desktop, keeping one tracking issue per finding.
resource: .github/workflows/docs-drift.yml
tags: [ci, drift, schedule]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/docs-drift.yml, tools/docdrift/issue.sh, tools/docdrift/generated.sh]
---
Triggers: cron `17 6 * * *` and `workflow_dispatch`. Concurrency group per workflow, not cancelled. Not a PR gate. Each job opens, edits or closes one issue through [issue.sh](/cli/issue-sh.md) (see [tracking issue](/data/drift-issue.md)). A job that cannot run fails red instead of opening an issue, so a broken guard cannot look like "all current".

| Job | Condition | Does |
|---|---|---|
| `generated` | always | `generated.sh upstream`, then sync or clear issue key `generated` with the pages that would change and the regenerate command |
| `relay` | repository is `Flockdeck/flockdeck-docs` and variable `RELAY_GHCR_READ_USER` is non-empty; environment `relay-release` | See below |
| `desktop` | always | Installs GTK 4 and WebKitGTK 6, [fetches desktop help](/cli/fetch-desktop-help-sh.md), runs `docdrift cli` against `app/cli.html`; key `desktop` |

`relay` job steps:
1. [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md) with the credential given to that step only. Exit 3 sets `ok=false` and syncs issue `relay-credential`; the job stays green. Other non-zero codes fail the job.
2. If fetched: clear `relay-credential`.
3. Compare `snapshot/relay.version` and the normalised snapshot text (CR and trailing whitespace stripped) with the fetched help. Different: sync `relay-snapshot` with the first 150 lines of the diff. Same: clear it.
4. Build docdrift and run `docdrift relay` on the fetched help against `content/self-hosting/configuration.md`; key `relay`. Exit 1 syncs, exit 0 clears, exit 2 fails the job.

Permissions: `contents: read`, `issues: write`. The relay job is dormant until its credential exists. Related: [relay credential isolation](/security/relay-credential-isolation.md), [docs guards](/ci/docs-guards.md). Env: [GH_TOKEN](/config/gh-token.md), [GOTOOLCHAIN](/config/gotoolchain.md).
