---
type: Environment Variable
title: GH_TOKEN
description: Token the gh CLI uses in the drift and release scripts.
resource: tools/docdrift/issue.sh
tags: [github, auth]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/issue.sh, tools/docdrift/generated.sh, tools/docdrift/fetch-desktop-help.sh, .github/workflows/docs-drift.yml, .github/workflows/docs-guards.yml, .github/workflows/autotag.yml]
---
Read implicitly by `gh`. Workflows set it from `${{ github.token }}`.

| Consumer | Needs |
|---|---|
| [issue.sh](/cli/issue-sh.md) | issues write (Docs drift jobs declare `issues: write`) |
| [generated.sh](/cli/generated-sh.md), [fetch-desktop-help.sh](/cli/fetch-desktop-help-sh.md) | read access to release metadata and assets of the public flockdeck repository |
| [autotag.sh](/cli/autotag-sh.md) | `actions: read` to ask whether CI passed |

In [Auto-tag](/ci/autotag.md) it is the workflow token for `gh`, while the tag push uses the separate app token set at checkout ([release bot token](/security/release-bot-token.md)). Value is never stored in the repository.
