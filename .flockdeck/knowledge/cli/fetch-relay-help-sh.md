---
type: CLI Command
title: fetch-relay-help.sh
description: Pulls the latest released relay image from GHCR with a read credential and captures serve -h.
resource: tools/docdrift/fetch-relay-help.sh
tags: [drift, ghcr]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh]
---
```
fetch-relay-help.sh OUT_FILE
```

Steps:
1. Requires [RELAY_GHCR_READ_USER](/config/relay-ghcr-read-user.md) and [RELAY_GHCR_READ_TOKEN](/config/relay-ghcr-read-token.md); missing either exits 3.
2. Gets a pull-scoped registry token from `ghcr.io/token` with basic auth.
3. Tag is [RELAY_TAG](/config/relay-tag.md), or else the highest bare-semver tag from the registry's paginated tag list (`Link: rel="next"` followed, `sort -V`). No semver tag: exit 1.
4. Logs in with a throwaway [DOCKER_CONFIG](/config/docker-config.md), pulls quietly, runs `serve -h`, merging stderr.
5. Writes the text to OUT_FILE and the tag to `OUT_FILE.version`.

Exit codes: 3 for any credential problem (missing, refused token, tag list refused, login or pull refused); 1 or 2 for anything else. The [Docs drift](/ci/docs-drift.md) job turns 3 into an issue and leaves the run green, and any other code leaves it red. It never falls back to an older image, since checking against a stale binary is what the tool exists to prevent. Image name: [RELAY_IMAGE](/config/relay-image.md). Isolation of the credential: [relay credential](/security/relay-credential-isolation.md).
