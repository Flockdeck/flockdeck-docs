---
type: Module
title: Allowlist loader
description: Reads a file of deliberate documentation omissions and rejects entries with no reason or duplicates.
resource: tools/docdrift/allow.go
tags: [docs, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/allow.go]
---
`loadAllow(path)` returns a map of name to reason. See [allowlist format](/data/allowlist-files.md) for the file shape.

Rules enforced:
- Blank lines and lines starting with `#` are skipped.
- A line is the name, a space, then the reason (split on the first space).
- No reason: error `path:line: "name" has no reason`. The file exists to record why something is undocumented.
- A name listed twice: error.

Errors become exit status 2 in the [docdrift command](/cli/docdrift.md). Part of [docdrift](/packages/docdrift.md).
