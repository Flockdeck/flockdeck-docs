---
type: CI Check
title: CI workflow
description: Checks the site files, builds the image, serves and tests it, then pushes it to GHCR.
resource: .github/workflows/ci.yml
tags: [ci, docker, nginx]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/ci.yml, nginx.conf, Dockerfile]
---
Triggers: push to main, tags `v*`, and pull requests. Concurrency group per workflow and ref, cancel in progress. Permissions: `contents: read`, `packages: write`. No configured secret; the push uses `GITHUB_TOKEN`.

## Job `check` ("Check the site")
- Presence: `index.html` and `docs.css` exist, `index.html` has a `<title>` and the word Flockdeck, `app/getting-started.html` and `self-hosting/overview.html` exist.
- Local links resolve: every `href`, `src` and `srcset` candidate in every `*.html` must resolve to a repo file. External, `mailto:`, `data:` and fragment-only refs are skipped; `/x` is repo-root relative, other refs are relative to the page's directory; a trailing `/` means `index.html`.
- HTML parses: `xmllint --html --noout` on every page.

## Job `image` (needs `check`)
1. Buildx; log in to GHCR unless PR.
2. Metadata, see [image tags](/build/image-tags.md).
3. Build into the runner's Docker (`load: true`, `push: false`) with GHA cache.
4. Run the image on port 8080 and assert: `/` contains Flockdeck; `getting-started`, `overview` and `docs.css` return 200; `Cache-Control: no-cache` on `/`, `index.html`, `docs.css`, two pages and both favicons; `favicon.svg`, `favicon.ico`, `apple-touch-icon.png` have an `image/*` content type; four security headers on seven paths including a 404 ([headers](/security/nginx-response-headers.md)); the `Server` header carries no digit; the container log has no request lines ([no access log](/security/no-access-log.md)).
5. Push every tag, only after step 4 and only outside PRs ([gate](/security/push-after-check.md)).

Gotchas: [link check loop](/gotchas/link-check-subshell.md). Config: [nginx](/build/nginx-conf.md), [Dockerfile](/build/dockerfile.md). Triggers [Auto-tag](/ci/autotag.md).
