---
type: Security Boundary
title: Relay registry credential isolation
description: How the read credential for the private relay image is kept out of pull requests and shared machines.
resource: .github/workflows/docs-drift.yml
tags: [secret-handling, ghcr]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/docs-drift.yml, tools/docdrift/fetch-relay-help.sh, .github/workflows/docs-guards.yml]
---
Controls, all visible in code:
- The image is private, so pull requests never receive a credential. [Docs guards](/ci/docs-guards.md) compares against the checked-in [snapshot](/data/relay-snapshot.md) instead.
- Only the scheduled `relay` job in [Docs drift](/ci/docs-drift.md) reads it, inside environment `relay-release`; the comments state its deployment branches are restricted to main.
- The job is skipped unless the repository is `Flockdeck/flockdeck-docs` and the user variable is set, so forks cannot reach it.
- The secret is placed in the env of the single step that runs [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md).
- The script logs in with `--password-stdin` into a temporary [DOCKER_CONFIG](/config/docker-config.md) removed on exit.
- Scope per the script header: a classic token with `read:packages` only.
- A refused credential exits 3: an issue is opened, the run stays green, and no older image is substituted.

Variables: [RELAY_GHCR_READ_USER](/config/relay-ghcr-read-user.md), [RELAY_GHCR_READ_TOKEN](/config/relay-ghcr-read-token.md).
