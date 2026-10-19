---
type: Module
title: docdrift
description: A Go program that checks hand-maintained docs pages against captured -h output of released Flockdeck binaries.
resource: tools/docdrift/main.go
tags: [docs, drift, go]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/main.go, tools/docdrift/go.mod]
---
`tools/docdrift` is its own Go module (`github.com/Flockdeck/flockdeck-docs/tools/docdrift`, go 1.22, no dependencies). It never runs or downloads anything: it reads a captured help text, a docs page and an allowlist, and prints findings.

| Piece | Concept |
|---|---|
| Entry point and modes | [docdrift command](/cli/docdrift.md) |
| `flockdeck-relay serve -h` and the configuration page | [relay parser](/packages/docdrift-relay-parser.md) |
| `flockdeck -h` and the CLI page | [CLI parser](/packages/docdrift-cli-parser.md) |
| Allowlist loader | [allowlist loader](/packages/docdrift-allowlist.md) |
| Tests | [docdrift tests](/tools/docdrift-tests.md) |

The scripts that produce the help text sit beside it: [fetch-desktop-help.sh](/cli/fetch-desktop-help-sh.md), [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md). Workflows that run it: [Docs guards](/ci/docs-guards.md) on pull requests and [Docs drift](/ci/docs-drift.md) on a schedule.

Findings are sorted strings so a head run and a base run can be diffed. `sortedKeys` is a generic helper in `relay.go`, which is why the module needs Go 1.18 or later.
