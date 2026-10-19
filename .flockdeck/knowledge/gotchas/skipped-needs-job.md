---
type: Gotcha
title: A job that needs a skipped job is skipped
description: PR review post spells out its own condition so trusted authors still get reviews when the gate is skipped.
resource: .github/workflows/pr-review-post.yml
tags: [github-actions]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review-post.yml]
---
Two GitHub rules shape [PR review post](/ci/pr-review-post.md): `environment:` cannot be combined with `uses:` on one job, so approval (`gate`) and posting (`post`) are separate jobs; and a job whose `needs` includes a skipped job is skipped by default.

`gate` is deliberately skipped for trusted authors, so `post` uses `always()` plus an explicit expression: the upstream run succeeded and (`gate` succeeded, or `gate` was skipped and `trusted == 'true'`). A skipped `gate` for any other reason does not pass. See [fork review split](/security/fork-pr-review-split.md).
