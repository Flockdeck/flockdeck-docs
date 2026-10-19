---
type: Environment Variable
title: GOTOOLCHAIN
description: Set to auto in the generation steps so the go command can fetch the Go version docgen's module asks for.
resource: .github/workflows/docs-guards.yml
tags: [go, ci]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/workflows/docs-guards.yml, .github/workflows/docs-drift.yml]
---
`actions/setup-go` pins `GOTOOLCHAIN=local`. The `generated` jobs in [Docs guards](/ci/docs-guards.md) and [Docs drift](/ci/docs-drift.md) override it with `auto`, because `cmd/docgen` in the flockdeck repository may require a newer Go than the one installed from `tools/docdrift/go.mod`. Without it `go run ./cmd/docgen` in [generated.sh](/cli/generated-sh.md) would fail on a version mismatch.
