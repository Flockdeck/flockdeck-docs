---
type: Build Step
title: Image tags
description: Which container tags CI produces for pull requests, main, and version tags.
resource: .github/workflows/ci.yml
tags: [docker, ghcr, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/ci.yml]
---
Image: `ghcr.io/<repository_owner>/flockdeck-docs`. Tags come from `docker/metadata-action@v5`:

| Rule | Produces | When |
|---|---|---|
| `type=raw,value=latest,enable={{is_default_branch}}` | `latest` | push to main |
| `type=sha,format=short` | `sha-<short>` | pushes |
| `type=semver,pattern={{version}}` | bare `0.9.24` | push of tag `v0.9.24` |
| `type=ref,event=pr` | `pr-<n>` | pull requests (built, never pushed) |

The bare semver is the one that deploys: the file states a Flux image policy elsewhere matches `>=0.0.1`, so a `v`-prefixed tag would be published and never rolled out. `latest` and `sha-*` are ignored by that policy. Pushed only on non-PR events, after the checks in [CI](/ci/ci.md).

See [release by tag](/build/release-by-tag.md) and [tag gotcha](/gotchas/bare-semver-tag.md).
