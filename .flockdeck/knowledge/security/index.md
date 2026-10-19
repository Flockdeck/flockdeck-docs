# Security boundaries

* [Checksum-verified desktop download](checksum-verified-download.md) - The released desktop tarball is verified against checksums.txt before it is executed.
* [Fork PR review split](fork-pr-review-split.md) - Untrusted pull request code and secrets are kept in two workflows, joined only by a diff artifact.
* [Security response headers](nginx-response-headers.md) - The four headers nginx adds to every answer, and how CI asserts them.
* [No access log](no-access-log.md) - The site keeps no request log, and CI fails if the container logs requests.
* [Pinned actions and reusable workflows](pinned-reusable-workflows.md) - Reusable workflow calls and most third-party actions are referenced by full commit SHA, not a moving ref.
* [Image is pushed only after it is served and checked](push-after-check.md) - CI builds locally, tests the running container, and pushes tags afterwards, because Flux deploys within minutes of a push.
* [Relay registry credential isolation](relay-credential-isolation.md) - How the read credential for the private relay image is kept out of pull requests and shared machines.
* [Release bot token](release-bot-token.md) - A GitHub App installation token with contents write is the only credential able to push release tags.
