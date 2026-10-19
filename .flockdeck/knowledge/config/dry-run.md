---
type: Environment Variable
title: DRY_RUN
description: When set to any non-empty value, autotag.sh prints the tag it would create and pushes nothing.
resource: .github/scripts/autotag.sh
tags: [release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh]
---
Output is `would tag <sha> as <next> (after <last>)`. All the skip conditions still apply first. See [autotag.sh](/cli/autotag-sh.md).
