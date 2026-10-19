---
type: Setting
title: PR_REVIEW_APP_PRIVATE_KEY
description: Secret holding the review GitHub App's private key; passed to the reusable post workflow.
resource: .github/workflows/pr-review-post.yml
tags: [github-app, secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review-post.yml]
---
Mapped to `review-app-private-key` for the `post` job of [PR review post](/ci/pr-review-post.md). The comment on that job says the call is pinned to a commit because this job holds the app key and the LLM key, and what it runs should not change when the called repository's main moves. See [pinned reusable workflows](/security/pinned-reusable-workflows.md). Related: [PR_REVIEW_APP_ID](/config/pr-review-app-id.md).
