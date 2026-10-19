---
type: CLI Command
title: generated.sh
description: Regenerates, verifies, or reports drift of the docgen-produced pages from the flockdeck repository.
resource: tools/docdrift/generated.sh
tags: [generator, docs]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/generated.sh, .docs-generated-from]
---
```
generated.sh regenerate [TAG]
generated.sh check
generated.sh upstream REPORT
```

| Mode | What it does |
|---|---|
| `regenerate [TAG]` | Clones `Flockdeck/flockdeck` at TAG (default: latest GitHub release via `gh release view`), deletes `app/` and `self-hosting/`, runs `go run ./cmd/docgen -out <repo root>`, writes `<tag> <sha>` to the [stamp](/data/docs-generated-from.md) |
| `check` | Clones the stamped tag, requires the tag to still point at the stamped SHA, regenerates into a scratch copy and fails if anything differs |
| `upstream REPORT` | Same against the latest release. Writes changed pages to REPORT, and the tags to `REPORT.latest` and `REPORT.stamped`. Exits 0 whether or not there is drift |

Details:
- `check` and `upstream` work from `git archive HEAD`, not the working copy, so they see what CI sees and not CRLF or uncommitted edits. The scratch copy has `app/` and `self-hosting/` removed first, so a page docgen no longer makes shows up as "only in the repo".
- `changed_pages` uses `diff -rq` and rewrites its lines as `changed:`, `only in the repo (docgen no longer makes it):` or `missing from the repo (docgen makes it):`.
- Moved tag: `check` fails with a message to regenerate if the tag no longer resolves to the stamped SHA.
- Any other argument prints the header lines and exits 2.

Needs `gh`, `git`, `go`, `tar`. The generator is `cmd/docgen` in the flockdeck repository; see [page generation](/build/regenerate-pages.md). Used by [Docs guards](/ci/docs-guards.md) (`check`) and [Docs drift](/ci/docs-drift.md) (`upstream`). Env: [GOTOOLCHAIN](/config/gotoolchain.md), [GH_TOKEN](/config/gh-token.md).
