---
type: Environment Variable
title: REPO
description: Owner/name of the repository autotag.sh asks about CI runs; required.
resource: .github/scripts/autotag.sh
tags: [release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh, .github/workflows/autotag.yml]
---
`main` fails immediately with `REPO is required` if unset. Used in `gh api repos/$REPO/actions/workflows/...` by `ci_green`. The workflow sets it to `github.repository`. See [autotag.sh](/cli/autotag-sh.md).
