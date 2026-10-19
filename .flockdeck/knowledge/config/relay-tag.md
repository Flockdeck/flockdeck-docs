---
type: Environment Variable
title: RELAY_TAG
description: Pins the relay image version that fetch-relay-help.sh pulls, as a bare semver.
resource: tools/docdrift/fetch-relay-help.sh
tags: [ghcr, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh, tools/docdrift/refresh-relay-snapshot.sh]
---
Optional. A bare semver such as `0.2.50`. When unset or empty, [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md) lists registry tags matching `[0-9]+.[0-9]+.[0-9]+` and takes the highest by `sort -V`.

[refresh-relay-snapshot.sh](/cli/refresh-relay-snapshot-sh.md) sets it from its first argument (empty when omitted).
