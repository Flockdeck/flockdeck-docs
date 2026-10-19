---
type: Environment Variable
title: RELAY_IMAGE
description: Overrides the registry repository fetch-relay-help.sh pulls the relay from.
resource: tools/docdrift/fetch-relay-help.sh
tags: [ghcr, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh]
---
Default `ghcr.io/flockdeck/flockdeck-relay`. The script strips the `ghcr.io/` prefix to build the token scope and tag-list URL, so the value must be a GHCR path. The image is private by design. Used by [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md).
