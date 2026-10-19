---
type: CI Check
title: PR review workflow (untrusted half)
description: Computes a pull request's diff with a read-only token by calling a pinned reusable workflow in the flockdeck repository.
resource: .github/workflows/pr-review.yml
tags: [ci, review, security]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review.yml]
---
Trigger: `pull_request` (opened, synchronize, reopened); the plain trigger, not `pull_request_target`. Permissions: `contents: read`.

Single job `diff` calls `Flockdeck/flockdeck/.github/workflows/reusable-pr-diff.yml` pinned to a full commit SHA, passing `pr-number`, `head-sha` and `base-ref`. The comment says the pin matters because this call is what a fork PR's untrusted job runs, and a moving ref would let anyone who can merge in the larger repo change what runs here.

It has no secrets and a fork PR gets a read-only, repo-scoped token. It only produces a diff artifact; the trusted half is [PR review post](/ci/pr-review-post.md). Boundary description: [fork review split](/security/fork-pr-review-split.md).
