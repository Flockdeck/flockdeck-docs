# Known gotchas

* [The deployed image tag has no v prefix](bare-semver-tag.md) - Git tags are vX.Y.Z but the image tag Flux matches is bare X.Y.Z.
* [Concurrent tagging races](concurrent-tagging.md) - Two auto-tag runs, or a person, can pick the same next version; the loser recomputes.
* [The PR relay check reads the snapshot, not the newest release](docs-guards-relay-source.md) - The Docs guards header says relay and desktop checks use the newest release, but the relay job uses the checked-in snapshot.
* [Generated pages must not be edited by hand](generated-pages-hand-edit.md) - app/ and self-hosting/ HTML is overwritten by docgen, and the PR guard fails on any manual edit.
* [A tag pushed with GITHUB_TOKEN starts no workflow](github-token-tags-start-nothing.md) - Auto-tag pushes tags with a GitHub App token so the tag's own CI run happens.
* [A changed help format fails the check instead of passing](help-format-fails-loudly.md) - Both parsers require a known heading and return an error when it or all entries are missing.
* [Link check loops must not run in a pipe subshell](link-check-subshell.md) - The CI link check feeds its loops by process substitution so a failure sets the exit status.
* [A moved flockdeck release tag fails the generated check](moved-release-tag.md) - The stamp records a commit as well as a tag so a retagged release is detected.
* [A job that needs a skipped job is skipped](skipped-needs-job.md) - PR review post spells out its own condition so trusted authors still get reviews when the gate is skipped.
