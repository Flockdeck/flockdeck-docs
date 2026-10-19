---
type: Build Step
title: Root go.mod retraction
description: A go.mod with no code that exists only to retract two MIT-licensed versions.
resource: go.mod
tags: [go, licence]
timestamp: 2026-10-19T00:00:00Z
code_paths: [go.mod]
---
Module `github.com/Flockdeck/flockdeck-docs`, `go 1.16`, with `retract [v0.1.0, v0.1.1]`. The repository has no Go code at its root. A retraction is only read from the `go.mod` of the highest version, so the file must exist in every release.

The comment gives the reason: v0.1.0 and v0.1.1 carried the MIT licence, and later versions are PolyForm Noncommercial 1.0.0 (see `LICENSE`).

Excluded from the image by [.dockerignore](/build/dockerignore.md). The drift tool has its own separate module, see [docdrift](/packages/docdrift.md).
