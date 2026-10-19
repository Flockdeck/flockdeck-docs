---
type: Developer Tool
title: docdrift tests
description: Go unit tests for the help parsers and comparers, plus fixtures of real help output.
resource: tools/docdrift/relay_test.go
tags: [tests, go]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/relay_test.go, tools/docdrift/cli_test.go, tools/docdrift/testdata/desktop-v0.3.40-help.txt, tools/docdrift/testdata/relay-v0.2.50-serve-help.txt]
---
Run in CI by the `tool-tests` job of [Docs guards](/ci/docs-guards.md): `gofmt -l`, `go vet ./...`, `go test ./...`, and `bash -n` on each `tools/docdrift/*.sh`.

| File | Covers |
|---|---|
| `relay_test.go` | help parsing, config page parsing, `canonDefault`, undocumented / removed / wrong-default / env findings, allowlist cases, `newFindings`, `loadAllow` requiring a reason, the checked-in [snapshot](/data/relay-snapshot.md) being well formed |
| `cli_test.go` | `parseDesktopHelp` including nested commands, failing loudly on a missing `Subcommands:` heading, missing and removed commands and subcommands, prose and flags being ignored, allowlist cases |
| `testdata/` | Real captured help text from desktop v0.3.40 and relay v0.2.50, used by the "real help against a complete page" tests |

The "fails loudly" tests are the guard against the check silently passing when a help format changes. Subjects: [relay parser](/packages/docdrift-relay-parser.md), [CLI parser](/packages/docdrift-cli-parser.md).
