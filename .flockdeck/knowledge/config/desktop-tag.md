---
type: Environment Variable
title: DESKTOP_TAG
description: Pins the desktop release that fetch-desktop-help.sh downloads, for example v0.3.40.
resource: tools/docdrift/fetch-desktop-help.sh
tags: [release, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/fetch-desktop-help.sh]
---
Optional. Default is the latest release tag of `Flockdeck/flockdeck` from `gh release view`; an empty result exits 1. Used by [fetch-desktop-help.sh](/cli/fetch-desktop-help-sh.md). No workflow sets it.
