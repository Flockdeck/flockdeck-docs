---
type: CI Check
title: PR review post workflow (trusted half)
description: After the untrusted diff job succeeds, optionally waits for approval, then calls a pinned reusable workflow that posts the review.
resource: .github/workflows/pr-review-post.yml
tags: [ci, review, security]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/pr-review-post.yml]
---
Trigger: `workflow_run` of "PR review" completed. Permissions: `contents: read`, `pull-requests: write`. It never checks out the PR's code, only the diff artifact from [the untrusted job](/ci/pr-review.md).

Jobs:

| Job | Runs when | Does |
|---|---|---|
| `check-permission` | the triggering run concluded `success` | `github-script` calls `getCollaboratorPermissionLevel` for `workflow_run.actor.login`. `admin` or `write` sets output `trusted=true`. An API error is caught, logged as a warning and leaves `trusted=false` |
| `gate` | success, and `trusted != 'true'` | Holds environment `pr-review-external`; its body only echoes "approved". Required reviewers on that environment provide the human approval |
| `post` | `always()`, run succeeded, and (`gate` succeeded, or `gate` was skipped and `trusted == 'true'`) | Calls `Flockdeck/flockdeck/.github/workflows/reusable-pr-review-post.yml` pinned to a commit SHA, with `artifact-run-id` and secrets `PR_REVIEW_APP_ID`, `PR_REVIEW_APP_PRIVATE_KEY`, `PR_REVIEW_LLM_API_KEY` mapped to review-app-id, review-app-private-key and llm-api-key |

Why two jobs: GitHub rejects `environment:` on a job that also uses `uses:`. Why the explicit `if:` on `post`: a job whose `needs` includes a skipped job is skipped by default, which would stop reviews for trusted authors. Secrets: [PR_REVIEW_APP_ID](/config/pr-review-app-id.md), [PR_REVIEW_APP_PRIVATE_KEY](/config/pr-review-app-private-key.md), [PR_REVIEW_LLM_API_KEY](/config/pr-review-llm-api-key.md). Boundary: [fork review split](/security/fork-pr-review-split.md); [skipped-needs gotcha](/gotchas/skipped-needs-job.md).

The header says the `pr-review-external` environment needs manual setup, and that until it exists the gate runs ungated.
