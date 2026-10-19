---
type: Deploy Step
title: Release by tag
description: End-to-end path from a merge to main to a deployed site, which only happens through a v* tag.
resource: .github/workflows/autotag.yml
tags: [release, deploy]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/autotag.yml, .github/workflows/ci.yml, .github/scripts/autotag.sh, Dockerfile]
---
1. A PR merges to main. [CI](/ci/ci.md) runs on the push: checks, image build, served-image checks, then pushes `latest` and `sha-<short>`. Neither deploys.
2. [Auto-tag](/ci/autotag.md) runs when CI completes on main (also hourly at minute 17, and on manual dispatch). [autotag.sh](/cli/autotag-sh.md) tags the commit `vX.Y.(Z+1)` if CI passed and something that ships changed.
3. The tag push starts CI again, which pushes the bare semver image tag ([image tags](/build/image-tags.md)).
4. A Flux image policy in another repository (named in comments only) watches for the new semver tag and rolls it out. Nothing in this repository performs that step, so it is not verified here.

Merging alone deploys nothing. A merge that nobody tagged is live in the repository and not on the site. Auto-tag exists for that reason. The tag must be pushed with a GitHub App token because a tag pushed with `GITHUB_TOKEN` starts no workflow; see [release bot token](/security/release-bot-token.md) and [gotcha](/gotchas/github-token-tags-start-nothing.md).

Version series: each repository counts on from its own latest `vN.N.N` tag, patch only. Changes limited to README, `.github`, `content/` and the like do not tag ([.dockerignore](/build/dockerignore.md)); regenerated pages arriving from docgen do.
