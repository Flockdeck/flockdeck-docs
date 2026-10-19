---
type: Gotcha
title: A tag pushed with GITHUB_TOKEN starts no workflow
description: Auto-tag pushes tags with a GitHub App token so the tag's own CI run happens.
resource: .github/workflows/autotag.yml
tags: [release, github-actions]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml, .github/scripts/autotag.sh]
---
GitHub does not start workflows from events caused by `GITHUB_TOKEN`. A tag pushed that way would never trigger [CI](/ci/ci.md), so the bare-semver image would never be built and nothing would deploy. [Auto-tag](/ci/autotag.md) therefore mints an app token ([release bot token](/security/release-bot-token.md)) and checks out with it, so the plain `git push` in [autotag.sh](/cli/autotag-sh.md) uses it.

The converse protects against loops: tagging makes no commit, and the workflow only follows CI runs on main, so a tag's CI run does not retrigger tagging.
