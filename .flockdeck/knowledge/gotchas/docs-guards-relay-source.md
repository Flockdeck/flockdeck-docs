---
type: Gotcha
title: The PR relay check reads the snapshot, not the newest release
description: The Docs guards header says relay and desktop checks use the newest release, but the relay job uses the checked-in snapshot.
resource: .github/workflows/docs-guards.yml
tags: [ci, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/docs-guards.yml, tools/docdrift/snapshot/relay.version]
---
The header comment of [Docs guards](/ci/docs-guards.md) says "the relay and desktop checks compare against the newest release". In code, `desktop-cli` does (it downloads the latest release), but `relay-flags` runs `docdrift relay` against `tools/docdrift/snapshot/relay-serve-help.txt`, because the relay image is private and pull requests get no credential. The job's own comment says so.

Consequence: a relay release can change flags without any PR going red until the snapshot is refreshed. The scheduled [Docs drift](/ci/docs-drift.md) `relay-snapshot` issue is what flags a stale [snapshot](/data/relay-snapshot.md).
