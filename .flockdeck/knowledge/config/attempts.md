---
type: Environment Variable
title: ATTEMPTS
description: How many times autotag.sh retries a tag push that lost a race; default 3.
resource: .github/scripts/autotag.sh
tags: [release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh]
---
After a refused `git push`, the script refreshes tags and recomputes the next version, up to this many attempts, then fails with `could not push <tag> after N attempts`. See [autotag.sh](/cli/autotag-sh.md) and [concurrent tagging](/gotchas/concurrent-tagging.md).
