---
type: CLI Command
title: issue.sh
description: Keeps one open tracking issue per drift check: opens, edits, or closes it.
resource: tools/docdrift/issue.sh
tags: [github, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/issue.sh]
---
```
issue.sh sync KEY TITLE BODY_FILE
issue.sh clear KEY COMMENT
```

| Command | Behaviour |
|---|---|
| `sync` | Finds the open issue by label and marker. None: creates the `docs-drift` label (errors ignored) and the issue. Exists with identical title and body: prints "already up to date" and changes nothing. Otherwise edits title and body |
| `clear` | Closes the open issue for KEY with COMMENT, or prints that there is none |

The issue is found by label `docs-drift` plus a marker line `<!-- docdrift-check: KEY -->` at the top of its body, so the title can carry a version and still match. The lookup lists up to 100 open issues. KEYs used by [Docs drift](/ci/docs-drift.md): `generated`, `relay`, `relay-credential`, `relay-snapshot`, `desktop`. See [tracking issue](/data/drift-issue.md).

Needs `gh`, [GH_TOKEN](/config/gh-token.md) with issues write, and `GH_REPO` or a checkout. An unknown command exits 2.
