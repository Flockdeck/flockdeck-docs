---
type: Setting
title: RELEASE_BOT_APP_PRIVATE_KEY
description: Secret holding the release bot GitHub App's private key, used only to mint a short-lived tag-push token.
resource: .github/workflows/autotag.yml
tags: [github-app, secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml]
---
Referenced as `secrets.RELEASE_BOT_APP_PRIVATE_KEY`, passed as `private-key` to `actions/create-github-app-token`. The workflow requests `permission-contents: write` only; the file's comment says that is for this repository and for an hour. The key sits in the `release` environment, so a workflow on another branch cannot read it. See [release bot token](/security/release-bot-token.md). Pair: [RELEASE_BOT_APP_ID](/config/release-bot-app-id.md).
