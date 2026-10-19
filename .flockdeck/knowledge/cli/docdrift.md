---
type: CLI Command
title: docdrift relay | cli
description: Compares a captured help text with a docs page and an allowlist, and exits 1 when they disagree.
resource: tools/docdrift/main.go
tags: [docs, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/main.go]
---
```
docdrift relay|cli -help FILE -docs FILE -allow FILE [-baseline-docs FILE]
```

| Mode | Help text | Docs page checked in CI |
|---|---|---|
| `relay` | `flockdeck-relay serve -h`, see [relay parser](/packages/docdrift-relay-parser.md) | `content/self-hosting/configuration.md` |
| `cli` | `flockdeck -h`, see [CLI parser](/packages/docdrift-cli-parser.md) | `app/cli.html` |

| Flag | Meaning |
|---|---|
| `-help` | Captured `-h` output of a released binary |
| `-docs` | The page to check |
| `-allow` | Allowlist of deliberate omissions, see [allowlist format](/data/allowlist-files.md) |
| `-baseline-docs` | The page as the base branch has it. Only findings absent from the base fail |

`-help`, `-docs` and `-allow` are all required (`load` errors otherwise).

With `-baseline-docs`, the check runs twice. Findings the head has and the base lacks are printed as `- <finding>` and fail. Findings already on the base are printed as `note (already so on the base branch, not failing this PR)`.

Exit status: 0 nothing (new) to report, 1 findings, 2 could not run (bad usage, unreadable file, a help format with no recognised heading, an allowlist entry with no reason). With no arguments it prints a usage line and exits 2. An unknown mode also exits 2.

Used by [Docs guards](/ci/docs-guards.md) and [Docs drift](/ci/docs-drift.md).
