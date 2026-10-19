# Configuration and environment variables

* [ATTEMPTS](attempts.md) - How many times autotag.sh retries a tag push that lost a race; default 3.
* [AUTOTAG_SOURCED](autotag-sourced.md) - When non-empty, autotag.sh only defines its functions so tests can call them one at a time.
* [BASE_BRANCH](base-branch.md) - Branch autotag.sh tags from and asks about CI for; default main.
* [CI_WORKFLOW](ci-workflow.md) - File name of the CI workflow whose runs autotag.sh checks; default ci.yml.
* [DESKTOP_TAG](desktop-tag.md) - Pins the desktop release that fetch-desktop-help.sh downloads, for example v0.3.40.
* [DOCKER_CONFIG](docker-config.md) - Pointed at a temporary directory so the registry login is not left on the machine.
* [DRY_RUN](dry-run.md) - When set to any non-empty value, autotag.sh prints the tag it would create and pushes nothing.
* [GH_TOKEN](gh-token.md) - Token the gh CLI uses in the drift and release scripts.
* [GOTOOLCHAIN](gotoolchain.md) - Set to auto in the generation steps so the go command can fetch the Go version docgen's module asks for.
* [PR_REVIEW_APP_ID](pr-review-app-id.md) - Secret naming the GitHub App that posts PR reviews; passed to the reusable post workflow.
* [PR_REVIEW_APP_PRIVATE_KEY](pr-review-app-private-key.md) - Secret holding the review GitHub App's private key; passed to the reusable post workflow.
* [PR_REVIEW_LLM_API_KEY](pr-review-llm-api-key.md) - Secret API key for the language model that writes PR reviews; passed to the reusable post workflow.
* [RELAY_GHCR_READ_TOKEN](relay-ghcr-read-token.md) - Secret with read:packages used to pull the private relay image in the scheduled drift job.
* [RELAY_GHCR_READ_USER](relay-ghcr-read-user.md) - GitHub user or bot name used to read the private relay image; also the repository variable that enables the relay drift job.
* [RELAY_IMAGE](relay-image.md) - Overrides the registry repository fetch-relay-help.sh pulls the relay from.
* [RELAY_TAG](relay-tag.md) - Pins the relay image version that fetch-relay-help.sh pulls, as a bare semver.
* [RELEASE_BOT_APP_ID](release-bot-app-id.md) - Secret holding the release bot GitHub App's ID, read by the Auto-tag workflow.
* [RELEASE_BOT_APP_PRIVATE_KEY](release-bot-app-private-key.md) - Secret holding the release bot GitHub App's private key, used only to mint a short-lived tag-push token.
* [REPO](repo.md) - Owner/name of the repository autotag.sh asks about CI runs; required.
* [SHA](sha.md) - The commit autotag.sh considers tagging; defaults to the tip of origin/main.
