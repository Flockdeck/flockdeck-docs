# Data types and file formats

* [Drift allowlists](allowlist-files.md) - Two text files listing relay flags and desktop subcommands the docs leave out on purpose, each with a reason.
* [.docs-generated-from stamp](docs-generated-from.md) - One line, "<tag> <commit>", naming the flockdeck release the checked-in generated pages came from.
* [Drift tracking issue](drift-issue.md) - The single GitHub issue per drift check, identified by a label and a body marker.
* [Relay flag snapshot](relay-snapshot.md) - A checked-in copy of the released relay's serve -h plus its version, so pull requests can run the relay check without the private image.
* [Repository and site layout](site-layout.md) - Which top-level paths are generated, which are hand-kept, and which end up in the served image.
