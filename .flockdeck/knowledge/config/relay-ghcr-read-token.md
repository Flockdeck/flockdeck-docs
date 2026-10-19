---
type: Environment Variable
title: RELAY_GHCR_READ_TOKEN
description: Secret with read:packages used to pull the private relay image in the scheduled drift job.
resource: tools/docdrift/fetch-relay-help.sh
tags: [ghcr, auth, secret]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh, .github/workflows/docs-drift.yml]
---
Per the script header: a classic personal access token with `read:packages` and nothing else (fine-grained tokens cannot read packages). Stored as the `RELAY_GHCR_READ_TOKEN` secret of the `relay-release` environment, whose deployment branches are limited to main. Passed in the env of one workflow step only.

Used for the registry token request (basic auth) and for `docker login --password-stdin` into a throwaway [DOCKER_CONFIG](/config/docker-config.md). A refused token exits 3 and opens the `relay-credential` issue. See [relay credential isolation](/security/relay-credential-isolation.md) and [fetch-relay-help.sh](/cli/fetch-relay-help-sh.md). Pair: [RELAY_GHCR_READ_USER](/config/relay-ghcr-read-user.md).
