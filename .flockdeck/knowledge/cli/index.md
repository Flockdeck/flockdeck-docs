# CLI commands

* [autotag.sh](autotag-sh.md) - Tags a green main commit with the next patch version when something that ships changed.
* [docdrift relay | cli](docdrift.md) - Compares a captured help text with a docs page and an allowlist, and exits 1 when they disagree.
* [fetch-desktop-help.sh](fetch-desktop-help-sh.md) - Downloads the released Linux desktop binary, verifies its checksum, and captures flockdeck -h.
* [fetch-relay-help.sh](fetch-relay-help-sh.md) - Pulls the latest released relay image from GHCR with a read credential and captures serve -h.
* [generated.sh](generated-sh.md) - Regenerates, verifies, or reports drift of the docgen-produced pages from the flockdeck repository.
* [issue.sh](issue-sh.md) - Keeps one open tracking issue per drift check: opens, edits, or closes it.
* [refresh-relay-snapshot.sh](refresh-relay-snapshot-sh.md) - Maintainer script that replaces the checked-in relay snapshot with the latest released relay's help.
