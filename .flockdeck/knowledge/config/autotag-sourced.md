---
type: Environment Variable
title: AUTOTAG_SOURCED
description: When non-empty, autotag.sh only defines its functions so tests can call them one at a time.
resource: .github/scripts/autotag.sh
tags: [tests]
timestamp: 2026-10-19T00:00:00Z
code_paths: [.github/scripts/autotag.sh, .github/scripts/autotag_test.sh]
---
The last lines of the script call `main` only when this is empty. [autotag_test.sh](/ci/autotag-test.md) sets it before sourcing. Subject: [autotag.sh](/cli/autotag-sh.md).
