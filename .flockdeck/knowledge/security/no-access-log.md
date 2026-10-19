---
type: Security Boundary
title: No access log
description: The site keeps no request log, and CI fails if the container logs requests.
resource: nginx.conf
tags: [privacy, nginx]
timestamp: 2026-10-19T00:00:00Z
code_paths: [nginx.conf, .github/workflows/ci.yml]
---
`access_log off` and `log_not_found off` in [nginx.conf](/build/nginx-conf.md). The base image sends logs to the container's output, so each request would otherwise be a line naming the client address.

CI makes many requests to the running image, then greps `docker logs` for `"(GET|HEAD) /` or `client: `; any match fails the job ([CI](/ci/ci.md)). The check runs before the image is pushed ([push gate](/security/push-after-check.md)).
