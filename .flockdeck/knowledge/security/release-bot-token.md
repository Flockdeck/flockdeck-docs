---
type: Security Boundary
title: Release bot token
description: A GitHub App installation token with contents write is the only credential able to push release tags.
resource: .github/workflows/autotag.yml
tags: [github-app, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml, .github/scripts/autotag.sh]
---
Mechanism in [Auto-tag](/ci/autotag.md):
- The token is minted at run time from [RELEASE_BOT_APP_ID](/config/release-bot-app-id.md) and [RELEASE_BOT_APP_PRIVATE_KEY](/config/release-bot-app-private-key.md) with only `permission-contents: write`.
- Those secrets live in environment `release`; the comment says only main may deploy to it, so another branch's workflow cannot read them. The comment adds that the sibling private repository's copy has no `environment:` line, since that plan has no environments.
- Workflow `permissions` are `{}` at top level and `contents: read`, `actions: read` on the job.
- The job runs only after a successful push-triggered CI run whose head repository is this one.
- The checked-out token stays in git config for the push; the job runs nothing but `autotag.sh` from main.
- The commit SHA arrives via `env`, not script interpolation.
- The third-party action is pinned to a full commit SHA.

Why not `GITHUB_TOKEN`: [gotcha](/gotchas/github-token-tags-start-nothing.md).
