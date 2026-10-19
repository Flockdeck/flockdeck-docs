---
type: Build Step
title: nginx.conf
description: The site's single nginx server block: compression, security headers, no cache, no access log.
resource: nginx.conf
tags: [nginx, http]
timestamp: 2026-10-19T00:00:00Z
code_paths: [nginx.conf]
---
| Setting | Value |
|---|---|
| Listen | 80 on IPv4 and IPv6, `server_name localhost` |
| `server_tokens` | `off` |
| `access_log`, `log_not_found` | `off` |
| Root, index | `/usr/share/nginx/html`, `index.html` |
| gzip | on, `gzip_vary on`, min length 512; types css, plain text, svg, javascript |
| `location /` headers | `Cache-Control: no-cache`, plus the four in [response headers](/security/nginx-response-headers.md) |
| `error_page` | 500, 502, 503, 504 to `/50x.html` |

Caching: nothing is fingerprinted, because the docs regenerate on every push to main and not per release. Every file is revalidated each use (`no-cache`; nginx's default ETag turns an unchanged answer into a 304), so a visit after a deploy sees what was just published.

`add_header` is set inside `location /`, so the [CI image check](/ci/ci.md) fetches a nonexistent path to confirm error answers carry the headers too. See [no access log](/security/no-access-log.md).

Installed by [Dockerfile](/build/dockerfile.md).
