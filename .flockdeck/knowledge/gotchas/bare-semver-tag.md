---
type: Gotcha
title: The deployed image tag has no v prefix
description: Git tags are vX.Y.Z but the image tag Flux matches is bare X.Y.Z.
resource: .github/workflows/ci.yml
tags: [release, docker]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/ci.yml, .github/scripts/autotag.sh]
---
CI's metadata rule `type=semver,pattern={{version}}` strips the `v`. The workflow comment says the image policy matches `>=0.0.1`, so a tag carrying the `v` would be published and never deployed. Git tags still need the `v`: CI triggers on `v*`, and [autotag.sh](/cli/autotag-sh.md) counts only `^v[0-9]+\.[0-9]+\.[0-9]+$` tags, ignoring candidates such as `v1.2.3-rc.1`. See [image tags](/build/image-tags.md).
