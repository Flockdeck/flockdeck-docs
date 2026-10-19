---
type: Security Boundary
title: Fork PR review split
description: Untrusted pull request code and secrets are kept in two workflows, joined only by a diff artifact.
resource: .github/workflows/pr-review-post.yml
tags: [github-actions, forks]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review.yml, .github/workflows/pr-review-post.yml]
---
| Half | Trigger | Token | Secrets | Sees PR code |
|---|---|---|---|---|
| [PR review](/ci/pr-review.md) | `pull_request` | read-only, repo-scoped | none | yes, to compute a diff |
| [PR review post](/ci/pr-review-post.md) | `workflow_run` | `pull-requests: write` | app id, app key, LLM key | no, only the artifact |

Layers on the trusted half:
1. `check-permission` asks the API for the PR author's repo permission. Only `admin` or `write` is trusted; any error counts as untrusted (fail closed).
2. Untrusted authors hit `gate`, which needs approval through environment `pr-review-external` (required reviewers are configured outside the repo; the file says that setup is manual).
3. `post` runs only if the diff job succeeded and either gate passed or the author is trusted.
4. Both reusable-workflow calls are [pinned to a commit](/security/pinned-reusable-workflows.md).

Gotcha: [skipped needs](/gotchas/skipped-needs-job.md). Secrets: [PR_REVIEW_APP_ID](/config/pr-review-app-id.md), [PR_REVIEW_APP_PRIVATE_KEY](/config/pr-review-app-private-key.md), [PR_REVIEW_LLM_API_KEY](/config/pr-review-llm-api-key.md).
