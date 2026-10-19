---
type: Security Boundary
title: Image is pushed only after it is served and checked
description: CI builds locally, tests the running container, and pushes tags afterwards, because Flux deploys within minutes of a push.
resource: .github/workflows/ci.yml
tags: [deploy, ci]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/ci.yml]
---
The `image` job of [CI](/ci/ci.md) builds with `push: false, load: true`, runs the container and asserts content, cache headers, icons, [security headers](/security/nginx-response-headers.md) and [no access log](/security/no-access-log.md). The `Push` step then pushes the same image under every metadata tag, skipped for pull requests.

The comment's reason: a new tag is deployed within minutes, so pushing before the check would deploy whether or not the check passed. Pull requests never log in to GHCR.

Related: [image tags](/build/image-tags.md), [release by tag](/build/release-by-tag.md).
