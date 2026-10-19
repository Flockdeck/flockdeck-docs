---
type: Environment Variable
title: DOCKER_CONFIG
description: Pointed at a temporary directory so the registry login is not left on the machine.
resource: tools/docdrift/fetch-relay-help.sh
tags: [docker, secret-handling]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-relay-help.sh]
---
[fetch-relay-help.sh](/cli/fetch-relay-help-sh.md) exports `DOCKER_CONFIG=$(mktemp -d)` before `docker login`, and removes the directory in an `EXIT` trap. The login stores the credential there, so a shared machine keeps none after the script ends. The token is passed with `--password-stdin`, not on the command line. See [relay credential isolation](/security/relay-credential-isolation.md).
