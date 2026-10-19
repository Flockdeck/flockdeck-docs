---
type: Module
title: Relay flag parser and comparer
description: Parses the flags out of the relay's serve -h, parses the configuration page's flag tables, and reports disagreements.
resource: tools/docdrift/relay.go
tags: [docs, drift, parser]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/relay.go]
---
## Help side: `parseRelayHelp`

Reads Go `flag` package output. Only lines after the exact heading `Flags for serve:` count; without it the parse errors (a changed help format must fail, not pass). A flag starts at a line matching two spaces then `-name`. Usage text is the following tab-indented lines. A blank line ends a flag.

Per flag it extracts:
- `Default` from a trailing `(default ...)`, quotes stripped.
- `Envs` from every `FLOCKDECK_RELAY_*` name inside a `(or ...)` group.

No flags found is also an error.

## Docs side: `parseConfigDoc`

Markdown, not HTML. A flag table is a run of `|` lines whose header's first cell is `Flag`. Optional `Default` and `Env` columns are read. Mentions outside tables count only when a code span is exactly one flag (`-audit` style), or a flag appears in a fenced block that contains `flockdeck-relay serve`. A span right after "no " is skipped, since "no `-x` flag" says the flag is absent.

## Comparison: `compareRelay`

| Finding | Condition |
|---|---|
| undocumented | Flag in help, not mentioned, not allowlisted |
| allowlisted but now documented | Mentioned and in the allowlist |
| wrong default | Table default differs from help default after `canonDefault` |
| env var named by docs but not by help | Env column entry absent from the usage text |
| env var read by relay but missing in docs | Help `Envs` entry absent from the Env column |
| docs mention a removed flag | Mentioned, not in help |
| stale allowlist entry | Allowlisted, not in help |

`canonDefault` treats `off`, `none`, `false`, `0`, `0s` and empty as the same zero value, drops a trailing parenthetical (`720h (30 days)`), and compares durations by value (`720h` equals `720h0m0s`). A docs default that is blank is "not stated" and is not compared.

Allowlist: [allowlist format](/data/allowlist-files.md). Input snapshot: [relay snapshot](/data/relay-snapshot.md). Part of [docdrift](/packages/docdrift.md).
