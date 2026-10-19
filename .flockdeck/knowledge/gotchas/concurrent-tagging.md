---
type: Gotcha
title: Concurrent tagging races
description: Two auto-tag runs, or a person, can pick the same next version; the loser recomputes.
resource: .github/scripts/autotag.sh
tags: [release, race]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh, .github/workflows/autotag.yml, .github/scripts/autotag_test.sh]
---
Guards:
- Workflow concurrency group `autotag` with `cancel-in-progress: false` serialises runs.
- If the push is still refused, [autotag.sh](/cli/autotag-sh.md) force-fetches tags and recomputes from scratch (including the "already tagged" and "part of a later tag" checks), up to [ATTEMPTS](/config/attempts.md). Exhausting them is a failure, not a quiet skip.
- A commit already contained in the latest tag is skipped: when CI for an earlier commit finishes late, the later commit's tag already released it.
- An in-flight or missing CI run is "not green"; the `workflow_run` trigger and the hourly cron try again.

Covered by [the script's tests](/ci/autotag-test.md).
