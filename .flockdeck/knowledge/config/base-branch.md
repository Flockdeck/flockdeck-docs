---
type: Environment Variable
title: BASE_BRANCH
description: Branch autotag.sh tags from and asks about CI for; default main.
resource: .github/scripts/autotag.sh
tags: [release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh]
---
Used for `origin/$BASE_BRANCH` ancestry checks and as the `branch=` filter in the CI run query. No workflow overrides it. See [autotag.sh](/cli/autotag-sh.md).
