---
type: Build Step
title: Line-ending policy
description: .gitattributes forces LF for every text file and marks known binary types.
resource: .gitattributes
tags: [git, platform]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.gitattributes]
---
`* text=auto eol=lf` makes every text file LF in every checkout regardless of `core.autocrlf`. The listed binary extensions (png, jpg, jpeg, gif, webp, ico, icns, woff, woff2, ttf, otf, eot, zip, gz, tgz, tar, xz, pdf, exe, dll, so, dylib, syso) are marked `binary` so git never rewrites one it mistakes for text.

Why it matters here: a CRLF checkout makes `gofmt -l` report formatted files, makes the [generated-files check](/ci/docs-guards.md) report every page as changed, and breaks shell scripts on Linux. `generated.sh` also reads from `git archive HEAD` with `core.autocrlf=false` for the same reason ([generated.sh](/cli/generated-sh.md)). A `.gitattributes`-only commit does not ship ([autotag.sh](/cli/autotag-sh.md)).
