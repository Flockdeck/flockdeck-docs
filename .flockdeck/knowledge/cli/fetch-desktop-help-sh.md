---
type: CLI Command
title: fetch-desktop-help.sh
description: Downloads the released Linux desktop binary, verifies its checksum, and captures flockdeck -h.
resource: tools/docdrift/fetch-desktop-help.sh
tags: [drift, release]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-desktop-help.sh]
---
```
fetch-desktop-help.sh OUT_FILE
```

Steps:
1. Tag is [DESKTOP_TAG](/config/desktop-tag.md) or the latest release of `Flockdeck/flockdeck` from `gh release view`.
2. `gh release download` of `flockdeck_<tag>_linux_amd64.tar.gz` and `checksums.txt`.
3. Take only that tarball's line from `checksums.txt`; fail if there is none, then `sha256sum -c` on it.
4. Extract the `flockdeck` binary and run `flockdeck -h`, merging stderr into stdout (help prints to stderr, exit 0).
5. Write the text to OUT_FILE and the tag to `OUT_FILE.version`.

The binary links GTK 4 and WebKitGTK 6 and needs them to print `-h`, so workflows install `libgtk-4-1` and `libwebkitgtk-6.0-4` first. Needs [GH_TOKEN](/config/gh-token.md).

See [security note on checksum](/security/checksum-verified-download.md). Output feeds the [CLI parser](/packages/docdrift-cli-parser.md).
