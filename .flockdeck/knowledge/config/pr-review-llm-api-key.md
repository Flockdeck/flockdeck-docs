---
type: Setting
title: PR_REVIEW_LLM_API_KEY
description: Secret API key for the language model that writes PR reviews; passed to the reusable post workflow.
resource: .github/workflows/pr-review-post.yml
tags: [secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review-post.yml]
---
Mapped to the `llm-api-key` secret input in the `post` job of [PR review post](/ci/pr-review-post.md). Reachable only after the permission check or environment approval described there. See [fork review split](/security/fork-pr-review-split.md).
