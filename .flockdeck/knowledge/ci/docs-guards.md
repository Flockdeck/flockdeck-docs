---
type: CI Check
title: Docs guards workflow
description: Pull-request jobs that test the drift tool, check generated pages, and check relay flags and desktop commands are documented.
resource: .github/workflows/docs-guards.yml
tags: [ci, drift, docs]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/docs-guards.yml, tools/docdrift/main.go, tools/docdrift/generated.sh]
---
Trigger: `pull_request`. Read-only token, no secrets, so fork PRs run it. Concurrency cancels in-progress runs per ref. Actions are pinned to commit SHAs. Separate from [CI](/ci/ci.md), whose "Check the site" job is the one branch protection requires; adding a guard to the required checks is a settings change.

| Job (name) | Does |
|---|---|
| `tool-tests` (Drift tool tests) | `gofmt -l tools/docdrift` must be empty, `go vet ./...`, `go test ./...`, `bash -n` on each `tools/docdrift/*.sh` |
| `generated` (Generated files match) | `generated.sh check` with [GOTOOLCHAIN](/config/gotoolchain.md)=auto: regenerate from the stamped flockdeck and compare ([generated.sh](/cli/generated-sh.md)) |
| `relay-flags` (Relay flags are documented) | Fetch depth 2. Builds docdrift, runs `docdrift relay` against the checked-in [snapshot](/data/relay-snapshot.md) and `content/self-hosting/configuration.md`, with the base's copy from `HEAD^1` as `-baseline-docs` when it exists |
| `desktop-cli` (Desktop commands are documented) | Installs GTK 4 and WebKitGTK 6 libraries, [fetches the released desktop help](/cli/fetch-desktop-help-sh.md), runs `docdrift cli` against `app/cli.html` with `HEAD^1` as baseline |

Design: none may go red because flockdeck moved on while an unrelated PR sat open. `generated` compares against the stamp, not the newest release. The two flag checks fail only on findings the PR adds ([docdrift](/cli/docdrift.md)). Scheduled counterpart: [Docs drift](/ci/docs-drift.md). A `pull_request` checkout is a merge commit whose first parent is the base, hence `fetch-depth: 2` and `HEAD^1`.
