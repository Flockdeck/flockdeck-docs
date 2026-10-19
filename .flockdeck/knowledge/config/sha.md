---
type: Environment Variable
title: SHA
description: The commit autotag.sh considers tagging; defaults to the tip of origin/main.
resource: .github/scripts/autotag.sh
tags: [release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh, .github/workflows/autotag.yml]
---
The workflow sets it to `github.event.workflow_run.head_sha`, which is empty for schedule and manual runs, so the script falls back to `origin/$BASE_BRANCH`. The script requires it to name a commit present in the checkout (`::error::` otherwise), resolves it to a full SHA, and skips if it is not an ancestor of main. The workflow passes it through `env`, never inlined into script text. See [autotag.sh](/cli/autotag-sh.md).
