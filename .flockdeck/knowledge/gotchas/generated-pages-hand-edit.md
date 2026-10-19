---
type: Gotcha
title: Generated pages must not be edited by hand
description: app/ and self-hosting/ HTML is overwritten by docgen, and the PR guard fails on any manual edit.
resource: tools/docdrift/generated.sh
tags: [generator]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/generated.sh, .github/workflows/docs-guards.yml, .github/workflows/docs-drift.yml]
---
`generated.sh check` regenerates from the stamped flockdeck and fails with "Do not edit app/*.html or self-hosting/*.html directly." The same failure occurs when `content/self-hosting/*.md` changes without regeneration, because those files are docgen inputs. The `desktop` drift issue text says `app/cli.html` is generated from a help page in the flockdeck repository: fix it there, then regenerate here.

The Markdown for self-hosting pages is the exception: it lives in this repository and is edited here, then regenerated. See [page generation](/build/regenerate-pages.md), [generated.sh](/cli/generated-sh.md), [layout](/data/site-layout.md).
