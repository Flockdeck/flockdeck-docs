# Packages

* [Allowlist loader](docdrift-allowlist.md) - Reads a file of deliberate documentation omissions and rejects entries with no reason or duplicates.
* [Desktop subcommand parser and comparer](docdrift-cli-parser.md) - Parses the subcommands out of flockdeck -h and checks that the CLI page documents each one.
* [Relay flag parser and comparer](docdrift-relay-parser.md) - Parses the flags out of the relay's serve -h, parses the configuration page's flag tables, and reports disagreements.
* [docdrift](docdrift.md) - A Go program that checks hand-maintained docs pages against captured -h output of released Flockdeck binaries.
