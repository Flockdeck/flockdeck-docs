---
type: Gotcha
title: A moved flockdeck release tag fails the generated check
description: The stamp records a commit as well as a tag so a retagged release is detected.
resource: tools/docdrift/generated.sh
tags: [generator, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/generated.sh, .docs-generated-from]
---
`fetch_flockdeck TAG SHA` clones the tag shallowly and compares `HEAD` with the stamped SHA. A mismatch exits 1 with "A release tag was moved. Regenerate with: tools/docdrift/generated.sh regenerate <tag>". The `upstream` mode passes no SHA, since it follows the latest release by definition. See [stamp](/data/docs-generated-from.md) and [generated.sh](/cli/generated-sh.md).
