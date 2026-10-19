---
type: Environment Variable
title: RELAY_GHCR_READ_USER
description: GitHub user or bot name used to read the private relay image; also the repository variable that enables the relay drift job.
resource: tools/docdrift/fetch-relay-help.sh
tags: [ghcr, auth, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh, .github/workflows/docs-drift.yml]
---
Required by [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md); if empty (or the token is), the script exits 3.

In [Docs drift](/ci/docs-drift.md) it is a repository variable (`vars.RELAY_GHCR_READ_USER`), not a secret. The `relay` job's `if:` requires it to be non-empty and the repository to be `Flockdeck/flockdeck-docs`, so the job is skipped until set and on forks. Paired with [RELAY_GHCR_READ_TOKEN](/config/relay-ghcr-read-token.md).
