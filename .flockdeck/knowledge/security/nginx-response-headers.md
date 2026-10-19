---
type: Security Boundary
title: Security response headers
description: The four headers nginx adds to every answer, and how CI asserts them.
resource: nginx.conf
tags: [nginx, headers]
timestamp: 2026-10-19T00:00:00Z
code_paths: [nginx.conf, .github/workflows/ci.yml]
---
Set with `always` in `location /` of [nginx.conf](/build/nginx-conf.md):

| Header | Value |
|---|---|
| `Strict-Transport-Security` | `max-age=31536000` |
| `X-Content-Type-Options` | `nosniff` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `X-Frame-Options` | `DENY` |

`server_tokens off` hides the nginx version. No Content-Security-Policy is set.

The `image` job of [CI](/ci/ci.md) requests `/`, `index.html`, two pages, `docs.css`, `favicon.ico` and a nonexistent path, and fails if any lacks an exact header line (case-insensitive, whole line) or if the `Server` header contains a digit.
