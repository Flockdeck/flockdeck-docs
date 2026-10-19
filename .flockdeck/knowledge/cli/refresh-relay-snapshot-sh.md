---
type: CLI Command
title: refresh-relay-snapshot.sh
description: Maintainer script that replaces the checked-in relay snapshot with the latest released relay's help.
resource: tools/docdrift/refresh-relay-snapshot.sh
tags: [drift, snapshot]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/refresh-relay-snapshot.sh]
---
```
RELAY_GHCR_READ_USER=... RELAY_GHCR_READ_TOKEN=... tools/docdrift/refresh-relay-snapshot.sh [VERSION]
```

Runs [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md) with `RELAY_TAG` set from the optional VERSION into a temp dir, strips `\r` from the help and writes `snapshot/relay-serve-help.txt`, copies the version into `snapshot/relay.version`, and prints the new version. The CR stripping keeps the committed [snapshot](/data/relay-snapshot.md) the same from any OS.

After running, compare the docs with `docdrift relay` against the snapshot, fix or allowlist what it reports, and commit the snapshot with the docs.
