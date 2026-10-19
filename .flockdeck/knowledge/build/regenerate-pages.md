---
type: Build Step
title: Page generation by docgen
description: How app/ and self-hosting/ are produced from the flockdeck repository's cmd/docgen and this repo's content directory.
resource: tools/docdrift/generated.sh
tags: [generator, docs]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/generated.sh, .docs-generated-from, .github/workflows/docs-drift.yml, .github/workflows/docs-guards.yml]
---
The pages are built in the public `Flockdeck/flockdeck` repository, by `cmd/docgen` (run as `go run ./cmd/docgen -out <dir>`). Its inputs are help pages kept in that repository plus `content/self-hosting/*.md` read from `-out`. This repository holds the output and no generator code.

Flow:
1. A flockdeck release exists.
2. A maintainer runs `tools/docdrift/generated.sh regenerate [TAG]` ([generated.sh](/cli/generated-sh.md)); `app/` and `self-hosting/` are replaced and the [stamp](/data/docs-generated-from.md) updated.
3. The diff is reviewed and committed as a PR.
4. On pull requests the [Generated files match](/ci/docs-guards.md) job regenerates from the stamped tag and fails on any difference.
5. Daily, [Docs drift](/ci/docs-drift.md) regenerates from the latest release and opens an issue if pages would change. It does not touch `content/`.

A change to `content/self-hosting/*.md` alone also fails the check until regenerated, because those pages are inputs to docgen.

Go toolchain: [GOTOOLCHAIN](/config/gotoolchain.md) is set to `auto` so docgen's own `go.mod` can pull a newer Go than the one set up.

Related: [layout](/data/site-layout.md), [hand-edit gotcha](/gotchas/generated-pages-hand-edit.md), [moved tag gotcha](/gotchas/moved-release-tag.md).
