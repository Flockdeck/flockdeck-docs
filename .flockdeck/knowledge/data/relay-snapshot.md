---
type: File Format
title: Relay flag snapshot
description: A checked-in copy of the released relay's serve -h plus its version, so pull requests can run the relay check without the private image.
resource: tools/docdrift/snapshot/relay.version
tags: [docs, drift, snapshot]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/snapshot/relay.version, tools/docdrift/snapshot/relay-serve-help.txt, tools/docdrift/refresh-relay-snapshot.sh, tools/docdrift/relay_test.go]
---
| File | Content |
|---|---|
| `snapshot/relay-serve-help.txt` | `serve -h` output with LF line endings |
| `snapshot/relay.version` | Bare semver of the relay it came from (no leading v) |

Why it exists: the relay image is private, so pull requests, including fork PRs, get no credential. The [relay-flags job](/ci/docs-guards.md) compares the page with this file. The scheduled [Docs drift](/ci/docs-drift.md) relay job compares the snapshot with the latest release and opens a tracking issue when it is behind.

Refresh by hand with [refresh-relay-snapshot.sh](/cli/refresh-relay-snapshot-sh.md).

`TestRelaySnapshotIsWellFormed` requires the help file to parse, to hold at least 20 flags, and the version file to match `^\d+\.\d+\.\d+\s*$`. See [docdrift tests](/tools/docdrift-tests.md).
