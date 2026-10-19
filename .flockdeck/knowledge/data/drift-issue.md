---
type: Data Type
title: Drift tracking issue
description: The single GitHub issue per drift check, identified by a label and a body marker.
resource: tools/docdrift/issue.sh
tags: [github, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/issue.sh, .github/workflows/docs-drift.yml]
---
| Field | Value |
|---|---|
| Label | `docs-drift` (colour d4c5f9, "The docs disagree with a released Flockdeck artifact") |
| First body line | `<!-- docdrift-check: KEY -->` |
| Title | Free text, may contain a version; changing it does not lose the issue |

KEYs and the title each sets in [Docs drift](/ci/docs-drift.md):

| KEY | Opened when |
|---|---|
| `generated` | Regenerating from the latest flockdeck release would change pages |
| `relay` | `content/self-hosting/configuration.md` disagrees with the released relay |
| `relay-credential` | The relay registry credential is refused |
| `relay-snapshot` | The checked-in [snapshot](/data/relay-snapshot.md) is not the latest release |
| `desktop` | `app/cli.html` disagrees with the released desktop |

Lifecycle is managed by [issue.sh](/cli/issue-sh.md): edited only when title or body changes, closed with a comment when the drift goes away.
