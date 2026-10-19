---
type: Setting
title: RELEASE_BOT_APP_ID
description: Secret holding the release bot GitHub App's ID, read by the Auto-tag workflow.
resource: .github/workflows/autotag.yml
tags: [github-app, secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml]
---
Referenced as `secrets.RELEASE_BOT_APP_ID` and passed to `actions/create-github-app-token` as `app-id`. Stored in this repository's `release` environment, which only main may deploy to. See [release bot token](/security/release-bot-token.md) and [Auto-tag](/ci/autotag.md). Pair: [RELEASE_BOT_APP_PRIVATE_KEY](/config/release-bot-app-private-key.md).
