---
type: Gotcha
title: A changed help format fails the check instead of passing
description: Both parsers require a known heading and return an error when it or all entries are missing.
resource: tools/docdrift/relay.go
tags: [drift, parser]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/relay.go, tools/docdrift/cli.go, tools/docdrift/relay_test.go, tools/docdrift/cli_test.go]
---
The relay parser needs the line `Flags for serve:`; the desktop parser needs `Subcommands:` and stops at the first blank line. Missing heading, or zero entries, returns an error and [docdrift](/cli/docdrift.md) exits 2 (job red, no issue). An empty parse would otherwise report no findings and look like "all documented". Other parsing details that depend on the binaries' current output:
- Relay usage lines are tab-indented after a two-space `-name` line (Go `flag` format); the help must reach the file with stderr merged (`2>&1`) in both fetch scripts.
- Desktop commands must be two-space indented with deeper-indented descriptions.
- Env var names are taken only from `FLOCKDECK_RELAY_*` names inside `(or ...)` groups.

See [relay parser](/packages/docdrift-relay-parser.md), [CLI parser](/packages/docdrift-cli-parser.md), [tests](/tools/docdrift-tests.md).
