---
type: Security Boundary
title: Checksum-verified desktop download
description: The released desktop tarball is verified against checksums.txt before it is executed.
resource: tools/docdrift/fetch-desktop-help.sh
tags: [supply-chain]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-desktop-help.sh]
---
[fetch-desktop-help.sh](/cli/fetch-desktop-help-sh.md) downloads the tarball and `checksums.txt` together, selects only the line for that exact tarball name (`grep -E "[[:space:]]\*?<tarball>$"`), fails if no line matches, and only then runs `sha256sum -c` on it. The explicit empty-selection check exists because `sha256sum -c` on empty input succeeds.

The binary is then executed with `-h` only, in jobs that hold a read-only or issues-write token and no relay credential.
