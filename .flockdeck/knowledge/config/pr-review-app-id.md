---
type: Setting
title: PR_REVIEW_APP_ID
description: Secret naming the GitHub App that posts PR reviews; passed to the reusable post workflow.
resource: .github/workflows/pr-review-post.yml
tags: [github-app, secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review-post.yml]
---
Mapped to the reusable workflow's `review-app-id` secret input, in the trusted `post` job of [PR review post](/ci/pr-review-post.md) only. Never available to the [untrusted diff job](/ci/pr-review.md). See [fork review split](/security/fork-pr-review-split.md). Siblings: [PR_REVIEW_APP_PRIVATE_KEY](/config/pr-review-app-private-key.md), [PR_REVIEW_LLM_API_KEY](/config/pr-review-llm-api-key.md).
