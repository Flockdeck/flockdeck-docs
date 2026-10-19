---
type: File Format
title: .docs-generated-from stamp
description: One line, "<tag> <commit>", naming the flockdeck release the checked-in generated pages came from.
resource: .docs-generated-from
tags: [generator, stamp]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.docs-generated-from, tools/docdrift/generated.sh, .dockerignore]
---
Format: `v<semver> <40-hex commit>`, for example `v0.8.1 <sha>`. Written by `generated.sh regenerate`, read by `check` and `upstream` after stripping control characters (`tr -d "[:cntrl:]"`). A missing file, or one without both fields, fails `check` with an explanatory message.

Role:
- `check` regenerates from exactly this tag, so a newer flockdeck release never turns an unrelated docs PR red.
- `upstream` reports the stamp as the "last generated from" version in the tracking issue.
- The tag must still resolve to the stamped commit, otherwise `check` fails (moved tag).

Not served: [.dockerignore](/build/dockerignore.md) excludes it from the image. See [generated.sh](/cli/generated-sh.md).
