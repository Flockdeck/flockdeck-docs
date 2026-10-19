# CI checks

* [Auto-tag tests workflow](autotag-test.md) - Runs the shell test suite and shellcheck for the auto-tag script on pull requests and main.
* [Auto-tag workflow](autotag.md) - Mints a release-bot token and runs autotag.sh to tag green, shipping commits on main.
* [CI workflow](ci.md) - Checks the site files, builds the image, serves and tests it, then pushes it to GHCR.
* [Docs drift workflow](docs-drift.md) - Daily check of the site against the latest released flockdeck, relay and desktop, keeping one tracking issue per finding.
* [Docs guards workflow](docs-guards.md) - Pull-request jobs that test the drift tool, check generated pages, and check relay flags and desktop commands are documented.
* [PR review post workflow (trusted half)](pr-review-post.md) - After the untrusted diff job succeeds, optionally waits for approval, then calls a pinned reusable workflow that posts the review.
* [PR review workflow (untrusted half)](pr-review.md) - Computes a pull request's diff with a read-only token by calling a pinned reusable workflow in the flockdeck repository.
