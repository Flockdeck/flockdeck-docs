---
type: Gotcha
title: Link check loops must not run in a pipe subshell
description: The CI link check feeds its loops by process substitution so a failure sets the exit status.
resource: .github/workflows/ci.yml
tags: [shell, ci]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/ci.yml]
---
In the "Local links resolve" step of [CI](/ci/ci.md), the page list is read with `done < <(find ...)` and each `srcset` with a here-string. A `find | while read` pipe would run the loop in a subshell, so `missing=1` would never reach the final `exit "$missing"` and the check could never fail.

Resolution rules that also trip people up: a ref without a leading slash is relative to the page's directory, as in a browser (not the repository root); a ref ending in `/` is checked as `index.html`; `srcset` candidates are comma-split and the descriptor after the URL is dropped.
