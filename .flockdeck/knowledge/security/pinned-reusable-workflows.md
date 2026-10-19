---
type: Security Boundary
title: Pinned actions and reusable workflows
description: Reusable workflow calls and most third-party actions are referenced by full commit SHA, not a moving ref.
resource: .github/workflows/pr-review.yml
tags: [supply-chain, github-actions]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review.yml, .github/workflows/pr-review-post.yml, .github/workflows/autotag.yml, .github/workflows/docs-guards.yml, .github/workflows/docs-drift.yml, .github/workflows/ci.yml]
---
- [PR review](/ci/pr-review.md) and [PR review post](/ci/pr-review-post.md) call `Flockdeck/flockdeck` reusable workflows at the same 40-hex commit. The comments say to bump both together and deliberately.
- `actions/checkout` and `actions/setup-go` in [Docs guards](/ci/docs-guards.md) and [Docs drift](/ci/docs-drift.md) are pinned by SHA with a version comment. `create-github-app-token` in [Auto-tag](/ci/autotag.md) is too.
- [CI](/ci/ci.md), [Auto-tag tests](/ci/autotag-test.md) and `pr-review-post`'s `github-script` use version tags (`@v3`, `@v4`, `@v5`, `@v6`, `@v7`) instead.

Reason stated for the reusable pins: the larger repository is active, and a moving ref would let anyone who can merge there change what runs here without a change landing in this repository. The post call also holds the LLM and app keys.
