---
type: Build Step
title: .dockerignore
description: Decides which repository files stay out of the image, and doubles as the rule autotag uses for "ships".
resource: .dockerignore
tags: [docker, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.dockerignore, .github/scripts/autotag.sh]
---
Excluded from the image: `.git`, `.github`, `.gitignore`, `.dockerignore`, `Dockerfile`, `README.md`, `LICENSE`, `content`, `tools`, `.docs-generated-from`, `go.mod`.

Reasons given in the file: `content/` is the Markdown docgen reads, and the site is the HTML it wrote; `tools/` and the stamp are repository, not website; `go.mod` only carries a version retraction (see [go.mod](/build/go-mod-retraction.md)).

Second role: [autotag.sh](/cli/autotag-sh.md) reads this file from the target commit to decide whether a change ships. A change entirely under an excluded path does not trigger a tag, with two exceptions: `Dockerfile` and `.dockerignore` themselves count. A pattern the script cannot follow (a `!` negation) makes everything count. Adding a path here therefore also stops changes to it from deploying.

Used by [Dockerfile](/build/dockerfile.md). Layout: [site layout](/data/site-layout.md).
